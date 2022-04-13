package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/auth"
	"auction-back/models"
	"auction-back/ports"
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

func (r *mutationResolver) RequestModerateUserForm(ctx context.Context) (bool, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return false, err
	}

	form, err := getOrCreateUserDraftForm(r.DB, viewer)
	if err != nil {
		return false, fmt.Errorf("draft form: %w", err)
	}

	_, err = (&models.UserFormFilled{}).From(&form)
	if form.Password == nil {
		err = multierror.Append(err, fmt.Errorf("no password"))
	}

	if err != nil {
		return false, err
	}

	data := map[string]interface{}{}
	if err := r.TokenPort.Create(models.TokenActionModerateUserForm, viewer, data); err != nil {
		return false, fmt.Errorf("token create: %w", err)
	}

	return true, nil
}

func (r *mutationResolver) ApproveModerateUserForm(ctx context.Context, input models.TokenInput) (*models.UserResult, error) {
	viewer, err := auth.ForViewer(ctx)
	if err != nil {
		return nil, err
	}

	token, err := r.TokenPort.Activate(models.TokenActionModerateUserForm, input.Token, viewer)
	if err != nil {
		return nil, fmt.Errorf("token activate: %w", err)
	}

	form, err := r.DB.UserForm().Take(ports.UserFormTakeConfig{
		OrderBy:   ports.UserFormFieldCreatedAt,
		OrderDesc: true,
		UserFormsFilter: models.UserFormsFilter{
			UserID: []string{token.UserID},
			State: []models.UserFormState{
				models.UserFormStateCreated,
				models.UserFormStateDeclained,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("take: %w", err)
	}

	form.State = models.UserFormStateModerating
	err = r.DB.UserForm().Update(&form)
	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}

	return &models.UserResult{
		User: &viewer,
	}, nil
}

func (r *mutationResolver) ApproveUserForm(ctx context.Context, input models.ApproveUserFormInput) (*models.UserFormResult, error) {
	form, err := r.DB.UserForm().Get(input.UserFormID)

	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	if form.State != models.UserFormStateModerating {
		return nil, fmt.Errorf("state is not %s", models.UserFormStateModerating)
	}

	form.State = models.UserFormStateApproved
	if err := r.BankPort.UserFormApproved(form); err != nil {
		return nil, fmt.Errorf("bank: %w", err)
	}

	if err := r.DB.UserForm().Update(&form); err != nil {
		return nil, fmt.Errorf("db update: %w", err)
	}

	return &models.UserFormResult{
		UserForm: &form,
	}, nil
}

func (r *mutationResolver) DeclineUserForm(ctx context.Context, input models.DeclineUserFormInput) (*models.UserFormResult, error) {
	form, err := r.DB.UserForm().Get(input.UserFormID)

	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	if form.State != models.UserFormStateModerating && form.State != models.UserFormStateApproved {
		return nil, fmt.Errorf("state is not %s", models.UserFormStateModerating)
	}

	form.State = models.UserFormStateDeclained
	form.DeclainReason = input.DeclainReason
	if err := r.DB.UserForm().Update(&form); err != nil {
		return nil, fmt.Errorf("db update: %w", err)
	}

	return &models.UserFormResult{
		UserForm: &form,
	}, nil
}

func (r *queryResolver) UserForms(ctx context.Context, first *int, after *string, filter *models.UserFormsFilter) (*models.UserFormsConnection, error) {
	config := ports.UserFormPaginationConfig{
		First: first,
		After: after,
	}

	if filter != nil {
		config.UserFormsFilter = *filter
	}

	connection, err := r.DB.UserForm().Pagination(config)
	if err != nil {
		return nil, fmt.Errorf("db pagination: %w", err)
	}

	return &connection, nil
}
