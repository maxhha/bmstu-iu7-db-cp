package database

import (
	"auction-back/models"
	"auction-back/ports"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userFormDB struct{ *Database }

func (d *Database) UserForm() ports.UserFormDB { return &userFormDB{d} }

type UserForm struct {
	ID            string `gorm:"default:generated();"`
	UserID        string
	State         models.UserFormState `gorm:"default:'CREATED';"`
	Name          *string
	Password      *string
	Phone         *string
	Email         *string
	DeclainReason *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     sql.NullTime
}

func (u *UserForm) into() models.UserForm {
	return models.UserForm{
		ID:            u.ID,
		UserID:        u.UserID,
		State:         u.State,
		Name:          u.Name,
		Password:      u.Password,
		Phone:         u.Phone,
		Email:         u.Email,
		DeclainReason: u.DeclainReason,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		DeletedAt:     u.DeletedAt,
	}
}

func (f *UserForm) copy(form *models.UserForm) {
	if form == nil {
		return
	}
	f.ID = form.ID
	f.UserID = form.UserID
	f.State = form.State
	f.Name = form.Name
	f.Password = form.Password
	f.Phone = form.Phone
	f.Email = form.Email
	f.DeclainReason = form.DeclainReason
	f.CreatedAt = form.CreatedAt
	f.UpdatedAt = form.UpdatedAt
	f.DeletedAt = form.DeletedAt
}

var userFormFieldToColumn = map[ports.UserFormField]string{
	ports.UserFormFieldCreatedAt: "created_at",
}

func (d *userFormDB) Get(id string) (models.UserForm, error) {
	form := UserForm{}
	if err := d.db.Take(&form, "id = ?", id).Error; err != nil {
		return models.UserForm{}, fmt.Errorf("take: %w", convertError(err))
	}

	return form.into(), nil
}

func (d *userFormDB) filter(query *gorm.DB, config models.UserFormsFilter) *gorm.DB {
	if len(config.ID) > 0 {
		query = query.Where("id IN ?", config.ID)
	}

	if len(config.UserID) > 0 {
		query = query.Where("user_id IN ?", config.UserID)
	}

	if len(config.State) > 0 {
		query = query.Where("state IN ?", config.State)
	}

	return query
}

func (d *userFormDB) Pagination(config ports.UserFormPaginationConfig) (models.UserFormsConnection, error) {
	query := d.filter(d.db.Model(&UserForm{}), config.UserFormsFilter)
	query, err := paginationQueryByCreatedAtDesc(query, config.First, config.After)
	if err != nil {
		return models.UserFormsConnection{}, fmt.Errorf("pagination: %w", err)
	}

	var objs []models.UserForm
	if err := query.Find(&objs).Error; err != nil {
		return models.UserFormsConnection{}, fmt.Errorf("find: %w", err)
	}

	if len(objs) == 0 {
		return models.UserFormsConnection{
			PageInfo: &models.PageInfo{},
			Edges:    make([]*models.UserFormsConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if config.First != nil {
		hasNextPage = len(objs) > *config.First
		objs = objs[:len(objs)-1]
	}

	edges := make([]*models.UserFormsConnectionEdge, 0, len(objs))

	for _, obj := range objs {
		node := obj

		edges = append(edges, &models.UserFormsConnectionEdge{
			Cursor: obj.ID,
			Node:   &node,
		})
	}

	return models.UserFormsConnection{
		PageInfo: &models.PageInfo{
			HasNextPage: hasNextPage,
			StartCursor: &objs[0].ID,
			EndCursor:   &objs[len(objs)-1].ID,
		},
		Edges: edges,
	}, nil
}

func (d *userFormDB) Take(config ports.UserFormTakeConfig) (models.UserForm, error) {
	query := d.filter(&d.db, config.UserFormsFilter)

	if config.OrderBy != "" {
		column, ok := userFormFieldToColumn[config.OrderBy]
		if !ok {
			return models.UserForm{}, fmt.Errorf("unknown field '%s'", config.OrderBy)
		}

		query = query.Order(clause.OrderByColumn{
			Column: clause.Column{Name: column},
			Desc:   config.OrderDesc,
		})
	}

	userForm := UserForm{}
	if err := query.Take(&userForm).Error; err != nil {
		return models.UserForm{}, fmt.Errorf("take: %w", convertError(err))
	}

	return userForm.into(), nil
}

// TODO: check if gorm.DB.Update updates objects field UpdatedAt
func (d *userFormDB) Update(form *models.UserForm) error {
	if form == nil {
		return fmt.Errorf("form is nil")
	}

	f := UserForm{}
	f.copy(form)

	if err := d.db.Save(&f).Error; err != nil {
		return fmt.Errorf("save: %w", err)
	}

	return nil
}

func (d *userFormDB) Create(form *models.UserForm) error {
	if form == nil {
		return fmt.Errorf("form is nil")
	}
	f := UserForm{}
	f.copy(form)
	if err := d.db.Create(&f).Error; err != nil {
		return fmt.Errorf("create: %w", err)
	}

	*form = f.into()
	return nil
}

func approvedOrFirstUserFormFilter(query *gorm.DB) *gorm.DB {
	return query.
		Where(`(
			state = 'APPROVED'
			OR (SELECT COUNT(1) FROM user_forms u WHERE "user_forms"."user_id" = u.user_id) = 1
		)`)
}

func (d *userFormDB) GetLoginForm(input models.LoginInput) (models.UserForm, error) {
	form := UserForm{}
	query := approvedOrFirstUserFormFilter(&d.db)
	err := query.
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
		return models.UserForm{}, fmt.Errorf("take: %w", convertError(err))
	}

	return form.into(), nil
}
