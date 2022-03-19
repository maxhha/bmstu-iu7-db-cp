package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/auth"
	"auction-back/db"
	"auction-back/graph/generated"
	"auction-back/graph/model"
	"auction-back/jwt"
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func (r *mutationResolver) Register(ctx context.Context) (*model.TokenResult, error) {
	user := db.User{}

	if err := r.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	token, err := jwt.NewUser(user.ID)

	if err != nil {
		return nil, err
	}

	return &model.TokenResult{
		Token: token,
	}, nil
}

func (r *mutationResolver) RequestSetUserEmail(ctx context.Context, input *model.RequestSetUserEmailInput) (*bool, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	data := map[string]interface{}{"email": input.Email}
	if err := r.Token.Create(db.TokenActionSetUserEmail, viewer, data); err != nil {
		return nil, err
	}

	res := true
	return &res, nil
}

func (r *mutationResolver) RequestSetUserPhone(ctx context.Context, input *model.RequestSetUserPhoneInput) (*bool, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	data := map[string]interface{}{"phone": input.Phone}
	if err := r.Token.Create(db.TokenActionSetUserPhone, viewer, data); err != nil {
		return nil, err
	}

	res := true
	return &res, nil
}

func (r *mutationResolver) ApproveSetUserEmail(ctx context.Context, input *model.TokenInput) (*model.UserResult, error) {
	viewer := auth.ForViewer(ctx)
	token, err := r.Token.Activate(db.TokenActionSetUserEmail, input.Token, viewer)

	if err != nil {
		return nil, err
	}

	email, ok := token.Data["email"].(string)

	if !ok {
		return nil, fmt.Errorf("no email in token")
	}

	form, err := r.getOrCreateUserForm(viewer)

	if err != nil {
		return nil, err
	}

	if err = r.DB.Model(&form).Update("email", email).Error; err != nil {
		return nil, err
	}

	return &model.UserResult{
		User: viewer,
	}, nil
}

func (r *mutationResolver) ApproveSetUserPhone(ctx context.Context, input *model.TokenInput) (*model.UserResult, error) {
	viewer := auth.ForViewer(ctx)
	token, err := r.Token.Activate(db.TokenActionSetUserEmail, input.Token, viewer)

	if err != nil {
		return nil, err
	}

	phone, ok := token.Data["phone"].(string)

	if !ok {
		return nil, fmt.Errorf("no phone in token")
	}

	form, err := r.getOrCreateUserForm(viewer)

	if err != nil {
		return nil, err
	}

	if err = r.DB.Model(&form).Update("phone", phone).Error; err != nil {
		return nil, fmt.Errorf("update: err")
	}

	return &model.UserResult{
		User: viewer,
	}, nil
}

func (r *mutationResolver) UpdateUserPassword(ctx context.Context, input *model.UpdateUserPasswordInput) (*model.UserResult, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	form, err := r.getOrCreateUserForm(viewer)

	if err != nil {
		return nil, err
	}

	old_password_equal := false

	if form.Password == nil && input.OldPassword == nil {
		old_password_equal = true
	}
	if form.Password != nil && input.OldPassword != nil {
		old_password_equal = checkHashAndPassword(*form.Password, *input.OldPassword)
	}

	if !old_password_equal {
		return nil, fmt.Errorf("old password not match")
	}

	password, err := hashPassword(input.Password)

	if err != nil {
		return nil, fmt.Errorf("hash: %w", err)
	}

	if err = r.DB.Model(&form).Update("password", password).Error; err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}

	return &model.UserResult{
		User: viewer,
	}, nil
}

func (r *mutationResolver) Login(ctx context.Context, input *model.LoginInput) (*model.TokenResult, error) {
	form := db.UserForm{}

	err := r.DB.
		Where(`(
			state = 'APPROVED' 
			OR (SELECT COUNT(1) FROM user_forms u WHERE "user_forms"."user_id" = u.user_id) = 1
		)`).
		Where(
			"name = @username OR email = @username OR phone = @username",
			sql.Named("username", input.Username),
		).
		Where(
			"password IS NOT NULL",
		).
		Take(
			&form,
		).Error

	if err != nil {
		return nil, fmt.Errorf("take: %w", err)
	}

	if form.Password == nil {
		return nil, fmt.Errorf("password not set")
	}

	if !checkHashAndPassword(*form.Password, input.Password) {
		return nil, fmt.Errorf("password mismatch")
	}

	token, err := jwt.NewUser(form.UserID)

	if err != nil {
		return nil, err
	}

	return &model.TokenResult{
		Token: token,
	}, nil
}

func (r *queryResolver) Viewer(ctx context.Context) (*db.User, error) {
	viewer := auth.ForViewer(ctx)
	return viewer, nil
}

func (r *userResolver) Form(ctx context.Context, obj *db.User) (*model.UserFormFilled, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) DraftForm(ctx context.Context, obj *db.User) (*db.UserForm, error) {
	viewer := auth.ForViewer(ctx)

	if viewer.ID != obj.ID {
		return nil, fmt.Errorf("denied")
	}

	form := db.UserForm{}
	query := "user_id = ? AND state in ('CREATED', 'MODERATING', 'DECLAINED')"
	err := r.DB.Order("created_at desc").Take(&form, query, obj.ID).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &form, nil
}

func (r *userResolver) FormHistory(ctx context.Context, obj *db.User, first *int, after *string) (*model.UserFormsConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) BlockedUntil(ctx context.Context, obj *db.User) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Available(ctx context.Context, obj *db.User) ([]*model.Money, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Blocked(ctx context.Context, obj *db.User) ([]*model.Money, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Accounts(ctx context.Context, obj *db.User, first *int, after *string) (*model.UserAccountsConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Offers(ctx context.Context, obj *db.User, first *int, after *string) (*model.OffersConnection, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	if viewer.ID != obj.ID {
		return nil, fmt.Errorf("denied")
	}

	query := db.DB.Where("consumer_id = ?", obj.ID).Order("id")

	return OfferPagination(query, first, after)
}

func (r *userResolver) Products(ctx context.Context, obj *db.User, first *int, after *string) (*model.ProductsConnection, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	if viewer.ID != obj.ID {
		return nil, fmt.Errorf("denied")
	}

	query := db.DB.Where("owner_id = ?", obj.ID).Order("id")

	return ProductPagination(query, first, after)
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
