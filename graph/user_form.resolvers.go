package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/auth"
	"auction-back/models"
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

func (r *mutationResolver) RequestModerateUserForm(ctx context.Context) (bool, error) {
	viewer := auth.ForViewer(ctx)

	if viewer == nil {
		return false, fmt.Errorf("unauthorized")
	}

	form, err := r.User().DraftForm(ctx, viewer)

	if err != nil {
		return false, fmt.Errorf("draft form: %w", err)
	}

	_, err = (&models.UserFormFilled{}).From(form)

	if form.Password == nil {
		err = multierror.Append(err, fmt.Errorf("no password"))
	}

	if err != nil {
		return false, err
	}

	data := map[string]interface{}{}
	if err := r.TokenPort.Create(models.TokenActionModerateUserForm, viewer, data); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ApproveModerateUserForm(ctx context.Context, input models.TokenInput) (*models.UserResult, error) {
	viewer := auth.ForViewer(ctx)
	token, err := r.TokenPort.Activate(models.TokenActionModerateUserForm, input.Token, viewer)

	if err != nil {
		return nil, fmt.Errorf("token activate: %w", err)
	}

	form := models.UserForm{}
	err = r.DB.
		Order("created_at desc").
		Take(
			&form,
			"user_id = ? AND state IN ?",
			token.UserID,
			[]models.UserFormState{models.UserFormStateCreated, models.UserFormStateDeclained},
		).
		Error

	if err != nil {
		return nil, fmt.Errorf("take: %w", err)
	}

	err = r.DB.Model(&form).Update("state", models.UserFormStateModerating).Error
	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}

	return &models.UserResult{
		User: viewer,
	}, nil
}

func (r *mutationResolver) ApproveUserForm(ctx context.Context, input models.ApproveUserFormInput) (*models.UserFormResult, error) {
	form := models.UserForm{}

	if err := r.DB.Take(&form, "id = ?", input.UserFormID).Error; err != nil {
		return nil, fmt.Errorf("take: %w", err)
	}

	if form.State != models.UserFormStateModerating {
		return nil, fmt.Errorf("state is not %s", models.UserFormStateModerating)
	}

	form.State = models.UserFormStateApproved

	if err := r.BankPort.UserFormApproved(form); err != nil {
		return nil, fmt.Errorf("bank: %w", err)
	}

	if err := r.DB.Save(&form).Error; err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}

	return &models.UserFormResult{
		UserForm: &form,
	}, nil
}

func (r *mutationResolver) DeclineUserForm(ctx context.Context, input models.DeclineUserFormInput) (*models.UserFormResult, error) {
	form := models.UserForm{}

	if err := r.DB.Take(&form, "id = ?", input.UserFormID).Error; err != nil {
		return nil, fmt.Errorf("take: %w", err)
	}

	if form.State != models.UserFormStateModerating {
		return nil, fmt.Errorf("state is not %s", models.UserFormStateModerating)
	}

	form.State = models.UserFormStateDeclained
	form.DeclainReason = input.DeclainReason

	if err := r.DB.Save(&form).Error; err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}

	return &models.UserFormResult{
		UserForm: &form,
	}, nil
}

func (r *queryResolver) UserForms(ctx context.Context, first *int, after *string, filter *models.UserFormsFilter) (*models.UserFormsConnection, error) {
	query := r.DB.Model(&models.UserForm{})

	if filter != nil {
		if len(filter.ID) > 0 {
			query = query.Where("id in ?", filter.ID)
		}

		if len(filter.State) > 0 {
			query = query.Where("state in ?", filter.State)
		}

		if len(filter.UserID) > 0 {
			query = query.Where("user_id in ?", filter.UserID)
		}
	}

	return UserFormPagination(query, first, after)
}
