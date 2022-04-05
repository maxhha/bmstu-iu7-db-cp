package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/auth"
	"auction-back/graph/generated"
	"auction-back/jwt"
	"auction-back/models"
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func (r *mutationResolver) Register(ctx context.Context) (*models.TokenResult, error) {
	user := models.User{}

	if err := r.DB.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	token, err := jwt.NewUser(user.ID)

	if err != nil {
		return nil, err
	}

	return &models.TokenResult{
		Token: token,
	}, nil
}

func (r *mutationResolver) Login(ctx context.Context, input models.LoginInput) (*models.TokenResult, error) {
	form := models.UserForm{}

	err := form.MostRelevantFilter(r.DB).
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

	return &models.TokenResult{
		Token: token,
	}, nil
}

func (r *mutationResolver) RequestSetUserEmail(ctx context.Context, input models.RequestSetUserEmailInput) (bool, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return false, fmt.Errorf("unauthorized")
	}

	data := map[string]interface{}{"email": input.Email}
	if err := r.TokenPort.Create(models.TokenActionSetUserEmail, viewer, data); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) RequestSetUserPhone(ctx context.Context, input models.RequestSetUserPhoneInput) (bool, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return false, fmt.Errorf("unauthorized")
	}

	data := map[string]interface{}{"phone": input.Phone}
	if err := r.TokenPort.Create(models.TokenActionSetUserPhone, viewer, data); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ApproveSetUserEmail(ctx context.Context, input models.TokenInput) (*models.UserResult, error) {
	viewer := auth.ForViewer(ctx)
	token, err := r.TokenPort.Activate(models.TokenActionSetUserEmail, input.Token, viewer)

	if err != nil {
		return nil, fmt.Errorf("token activate: %w", err)
	}

	email, ok := token.Data["email"].(string)

	if !ok {
		return nil, fmt.Errorf("no email in token")
	}

	form, err := getOrCreateUserDraftForm(r.DB, viewer)

	if err != nil {
		return nil, err
	}

	if err = r.DB.Model(&form).Update("email", email).Error; err != nil {
		return nil, err
	}

	return &models.UserResult{
		User: viewer,
	}, nil
}

func (r *mutationResolver) ApproveSetUserPhone(ctx context.Context, input models.TokenInput) (*models.UserResult, error) {
	viewer := auth.ForViewer(ctx)
	token, err := r.TokenPort.Activate(models.TokenActionSetUserPhone, input.Token, viewer)

	if err != nil {
		return nil, fmt.Errorf("token activate: %w", err)
	}

	phone, ok := token.Data["phone"].(string)

	if !ok {
		return nil, fmt.Errorf("no phone in token")
	}

	form, err := getOrCreateUserDraftForm(r.DB, viewer)

	if err != nil {
		return nil, err
	}

	if err = r.DB.Model(&form).Update("phone", phone).Error; err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}

	return &models.UserResult{
		User: viewer,
	}, nil
}

func (r *mutationResolver) UpdateUserPassword(ctx context.Context, input models.UpdateUserPasswordInput) (*models.UserResult, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	form, err := getOrCreateUserDraftForm(r.DB, viewer)

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

	return &models.UserResult{
		User: viewer,
	}, nil
}

func (r *mutationResolver) UpdateUserDraftForm(ctx context.Context, input models.UpdateUserDraftFormInput) (*models.UserResult, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	form, err := getOrCreateUserDraftForm(r.DB, viewer)

	if err != nil {
		return nil, err
	}

	form.Name = input.Name

	if err = r.DB.Save(&form).Error; err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}

	return &models.UserResult{
		User: viewer,
	}, nil
}

func (r *queryResolver) Viewer(ctx context.Context) (*models.User, error) {
	viewer := auth.ForViewer(ctx)
	return viewer, nil
}

func (r *queryResolver) Users(ctx context.Context, first *int, after *string, filter *models.UsersFilter) (*models.UsersConnection, error) {
	query := r.DB.Model(&models.User{})

	if filter != nil {
		if len(filter.ID) > 0 {
			query = query.Where("id in ?", filter.ID)
		}
	}

	return UserPagination(query, first, after)
}

func (r *userResolver) Form(ctx context.Context, obj *models.User) (*models.UserFormFilled, error) {
	viewer := auth.ForViewer(ctx)

	if err := r.isOwnerOrManager(viewer, obj); err != nil {
		return nil, fmt.Errorf("denied: %w", err)
	}

	form, err := viewer.LastApprovedUserForm(r.DB)

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return (&models.UserFormFilled{}).From(&form)
}

func (r *userResolver) DraftForm(ctx context.Context, obj *models.User) (*models.UserForm, error) {
	viewer := auth.ForViewer(ctx)

	if err := r.isOwnerOrManager(viewer, obj); err != nil {
		return nil, fmt.Errorf("denied: %w", err)
	}

	form := models.UserForm{}
	err := r.DB.Order("created_at desc").Take(&form, "user_id = ?", obj.ID).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if form.State == "CREATED" || form.State == "MODERATING" || form.State == "DECLAINED" {
		return &form, nil
	}

	return nil, nil
}

func (r *userResolver) FormHistory(ctx context.Context, obj *models.User, first *int, after *string, filter *models.UserFormHistoryFilter) (*models.UserFormsConnection, error) {
	if obj == nil {
		return nil, fmt.Errorf("user is nil")
	}

	query := r.DB.Model(&models.UserForm{}).Where("user_id = ?", obj.ID)

	if filter != nil {
		if len(filter.ID) > 0 {
			query = query.Where("id in ?", filter.ID)
		}

		if len(filter.State) > 0 {
			query = query.Where("state in ?", filter.State)
		}
	}

	return UserFormPagination(query, first, after)
}

func (r *userResolver) BlockedUntil(ctx context.Context, obj *models.User) (*time.Time, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Available(ctx context.Context, obj *models.User) ([]*models.Money, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Blocked(ctx context.Context, obj *models.User) ([]*models.Money, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) Accounts(ctx context.Context, obj *models.User, first *int, after *string) (*models.UserAccountsConnection, error) {
	viewer := auth.ForViewer(ctx)

	if err := r.isOwnerOrManager(viewer, obj); err != nil {
		return nil, fmt.Errorf("denied: %w", err)
	}

	query := r.DB.Model(&models.Account{}).Where("user_id = ?", obj.ID)

	return UserAccountPagination(query, first, after)
}

func (r *userResolver) Offers(ctx context.Context, obj *models.User, first *int, after *string) (*models.OffersConnection, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	if viewer.ID != obj.ID {
		return nil, fmt.Errorf("denied")
	}

	query := r.DB.Where("consumer_id = ?", obj.ID).Order("id")

	return OfferPagination(query, first, after)
}

func (r *userResolver) Products(ctx context.Context, obj *models.User, first *int, after *string) (*models.ProductsConnection, error) {
	viewer := auth.ForViewer(ctx)

	if err := r.isOwnerOrManager(viewer, obj); err != nil {
		return nil, err
	}

	query := r.DB.Where("owner_id = ?", obj.ID)

	return ProductPagination(query, first, after)
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
