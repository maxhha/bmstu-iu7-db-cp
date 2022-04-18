package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/auth"
	"auction-back/graph/generated"
	"auction-back/jwt"
	"auction-back/models"
	"auction-back/ports"
	"context"
	"errors"
	"fmt"
)

func (r *mutationResolver) Register(ctx context.Context) (*models.TokenResult, error) {
	user := models.User{}

	if err := r.DB.User().Create(&user); err != nil {
		return nil, fmt.Errorf("db create: %w", err)
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
	form, err := r.DB.UserForm().GetLoginForm(input)
	if err != nil {
		return nil, fmt.Errorf("get login form: %w", err)
	}

	if form.Password == nil {
		return nil, ErrNoPassword
	}

	if !checkHashAndPassword(*form.Password, input.Password) {
		return nil, ErrPasswordMissmatch
	}

	token, err := jwt.NewUser(form.UserID)
	if err != nil {
		return nil, fmt.Errorf("jwt new user: %w", err)
	}

	return &models.TokenResult{
		Token: token,
	}, nil
}

func (r *mutationResolver) RequestSetUserEmail(ctx context.Context, input models.RequestSetUserEmailInput) (bool, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return false, err
	}

	data := map[string]interface{}{"email": input.Email}
	if err := r.TokenPort.Create(models.TokenActionSetUserEmail, viewer, data); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) RequestSetUserPhone(ctx context.Context, input models.RequestSetUserPhoneInput) (bool, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return false, err
	}

	data := map[string]interface{}{"phone": input.Phone}
	if err := r.TokenPort.Create(models.TokenActionSetUserPhone, viewer, data); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ApproveSetUserEmail(ctx context.Context, input models.TokenInput) (*models.UserResult, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return nil, err
	}

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

	form.Email = &email
	if err = r.DB.UserForm().Update(&form); err != nil {
		return nil, fmt.Errorf("db update: %w", err)
	}

	return &models.UserResult{
		User: &viewer,
	}, nil
}

func (r *mutationResolver) ApproveSetUserPhone(ctx context.Context, input models.TokenInput) (*models.UserResult, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return nil, err
	}

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

	form.Phone = &phone
	if err = r.DB.UserForm().Update(&form); err != nil {
		return nil, fmt.Errorf("db update: %w", err)
	}

	return &models.UserResult{
		User: &viewer,
	}, nil
}

func (r *mutationResolver) UpdateUserPassword(ctx context.Context, input models.UpdateUserPasswordInput) (*models.UserResult, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return nil, err
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

	form.Password = &password
	if err = r.DB.UserForm().Update(&form); err != nil {
		return nil, fmt.Errorf("db update: %w", err)
	}

	return &models.UserResult{
		User: &viewer,
	}, nil
}

func (r *mutationResolver) UpdateUserDraftForm(ctx context.Context, input models.UpdateUserDraftFormInput) (*models.UserResult, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return nil, err
	}

	form, err := getOrCreateUserDraftForm(r.DB, viewer)

	if err != nil {
		return nil, err
	}

	form.Name = input.Name
	form.Currency = input.Currency
	if err = r.DB.UserForm().Update(&form); err != nil {
		return nil, fmt.Errorf("db update: %w", err)
	}

	return &models.UserResult{
		User: &viewer,
	}, nil
}

func (r *queryResolver) Viewer(ctx context.Context) (*models.User, error) {
	viewer, err := auth.ForViewer(ctx)
	if errors.Is(err, auth.ErrUnauthorized) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &viewer, nil
}

func (r *queryResolver) Users(ctx context.Context, first *int, after *string, filter *models.UsersFilter) (*models.UsersConnection, error) {
	connection, err := r.DB.User().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("db pagination: %w", err)
	}

	return &connection, nil
}

func (r *userResolver) Form(ctx context.Context, obj *models.User) (*models.UserFormFilled, error) {
	form, err := r.DB.User().LastApprovedUserForm(*obj)
	if errors.Is(err, ports.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return (&models.UserFormFilled{}).From(&form)
}

func (r *userResolver) DraftForm(ctx context.Context, obj *models.User) (*models.UserForm, error) {
	form, err := r.DB.UserForm().Take(ports.UserFormTakeConfig{
		OrderBy:   ports.UserFormFieldCreatedAt,
		OrderDesc: true,
		UserFormsFilter: models.UserFormsFilter{
			UserID: []string{obj.ID},
		},
	})

	if errors.Is(err, ports.ErrRecordNotFound) {
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
	userFormsFilter := models.UserFormsFilter{
		UserID: []string{obj.ID},
	}

	if filter != nil {
		userFormsFilter.ID = filter.ID
		userFormsFilter.State = filter.State
	}

	connection, err := r.DB.UserForm().Pagination(first, after, &userFormsFilter)
	if err != nil {
		return nil, fmt.Errorf("db pagination: %w", err)
	}

	return &connection, nil
}

func (r *userResolver) Accounts(ctx context.Context, obj *models.User, first *int, after *string, filter *models.AccountsFilter) (*models.AccountsConnection, error) {
	if filter == nil {
		filter = &models.AccountsFilter{}
	}

	if len(filter.UserIDs) > 0 {
		return nil, fmt.Errorf("userIDs filter must be empty")
	}

	filter.UserIDs = []string{obj.ID}

	connection, err := r.DB.Account().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("db account pagination: %w", err)
	}

	return &connection, nil
}

func (r *userResolver) Auctions(ctx context.Context, obj *models.User, first *int, after *string) (*models.AuctionsConnection, error) {
	filter := models.AuctionsFilter{
		SellerIDs: []string{obj.ID},
	}

	conn, err := r.DB.Auction().Pagination(first, after, &filter)
	if err != nil {
		return nil, fmt.Errorf("db auction pagination: %w", err)
	}

	return &conn, nil
}

func (r *userResolver) Offers(ctx context.Context, obj *models.User, first *int, after *string, filter *models.OffersFilter) (*models.OffersConnection, error) {
	if filter == nil {
		filter = &models.OffersFilter{}
	}

	if len(filter.UserIDs) > 0 {
		return nil, fmt.Errorf("userIDs must be empty")
	}

	filter.UserIDs = []string{obj.ID}
	connection, err := r.DB.Offer().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("db offer pagination: %w", err)
	}

	return &connection, nil
}

func (r *userResolver) Products(ctx context.Context, obj *models.User, first *int, after *string) (*models.ProductsConnection, error) {
	filter := &models.ProductsFilter{
		OwnerIDs: []string{obj.ID},
	}

	connection, err := r.DB.Product().Pagination(first, after, filter)
	if err != nil {
		return nil, fmt.Errorf("db pagination: %w", err)
	}
	return &connection, nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
