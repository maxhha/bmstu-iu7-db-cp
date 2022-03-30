package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"auction-back/auth"
	"auction-back/db"
	"auction-back/graph/generated"
	"auction-back/graph/model"
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

	_, err = (&model.UserFormFilled{}).From(form)

	if form.Password == nil {
		err = multierror.Append(err, fmt.Errorf("no password"))
	}

	if err != nil {
		return false, err
	}

	data := map[string]interface{}{}
	if err := r.Token.Create(db.TokenActionModerateUserForm, viewer, data); err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) ApproveModerateUserForm(ctx context.Context, input model.TokenInput) (*model.UserResult, error) {
	viewer := auth.ForViewer(ctx)
	token, err := r.Token.Activate(db.TokenActionModerateUserForm, input.Token, viewer)

	if err != nil {
		return nil, fmt.Errorf("token activate: %w", err)
	}

	form := db.UserForm{}
	err = r.DB.
		Order("created_at desc").
		Take(
			&form,
			"user_id = ? AND state IN ?",
			token.UserID,
			[]db.UserFormState{db.UserFormStateCreated, db.UserFormStateDeclained},
		).
		Error

	if err != nil {
		return nil, fmt.Errorf("take: %w", err)
	}

	err = r.DB.Model(&form).Update("state", db.UserFormStateModerating).Error
	if err != nil {
		return nil, fmt.Errorf("update: %w", err)
	}

	return &model.UserResult{
		User: viewer,
	}, nil
}

func (r *mutationResolver) ApproveUserForm(ctx context.Context, input model.ApproveUserFormInput) (*model.UserFormResult, error) {
	form := db.UserForm{}

	if err := r.DB.Take(&form, "id = ?", input.UserFormID).Error; err != nil {
		return nil, fmt.Errorf("take: %w", err)
	}

	if form.State != db.UserFormStateModerating {
		return nil, fmt.Errorf("state is not %s", db.UserFormStateModerating)
	}

	form.State = db.UserFormStateApproved

	if err := r.Bank.UserFormApproved(form); err != nil {
		return nil, fmt.Errorf("bank: %w", err)
	}

	if err := r.DB.Save(&form).Error; err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}

	return &model.UserFormResult{
		UserForm: &form,
	}, nil
}

func (r *mutationResolver) DeclineUserForm(ctx context.Context, input model.DeclineUserFormInput) (*model.UserFormResult, error) {
	form := db.UserForm{}

	if err := r.DB.Take(&form, "id = ?", input.UserFormID).Error; err != nil {
		return nil, fmt.Errorf("take: %w", err)
	}

	if form.State != db.UserFormStateModerating {
		return nil, fmt.Errorf("state is not %s", db.UserFormStateModerating)
	}

	form.State = db.UserFormStateDeclained
	form.DeclainReason = input.DeclainReason

	if err := r.DB.Save(&form).Error; err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}

	return &model.UserFormResult{
		UserForm: &form,
	}, nil
}

func (r *queryResolver) UserForms(ctx context.Context, first *int, after *string, filter *model.UserFormsFilter) (*model.UserFormsConnection, error) {
	query := r.DB.Model(&db.UserForm{}).Order("created_at desc")

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

func (r *userFormResolver) State(ctx context.Context, obj *db.UserForm) (model.UserFormStateEnum, error) {
	return model.UserFormStateEnum(obj.State), nil
}

// UserForm returns generated.UserFormResolver implementation.
func (r *Resolver) UserForm() generated.UserFormResolver { return &userFormResolver{r} }

type userFormResolver struct{ *Resolver }
