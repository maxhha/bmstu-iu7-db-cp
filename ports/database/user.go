package database

import (
	"auction-back/models"
	"auction-back/ports"
	"database/sql"
	"fmt"
	"time"
)

type userDB struct{ *Database }

func (d *Database) User() ports.UserDB { return &userDB{d} }

type User struct {
	ID           string    `gorm:"default:generated();"`
	CreatedAt    time.Time `gorm:"default:now();"`
	DeletedAt    sql.NullTime
	BlockedUntil sql.NullTime
}

func (u *User) into() models.User {
	return models.User{
		ID:           u.ID,
		CreatedAt:    u.CreatedAt,
		DeletedAt:    u.DeletedAt,
		BlockedUntil: u.BlockedUntil,
	}
}

func (u *User) copy(user *models.User) {
	if user == nil {
		return
	}
	u.ID = user.ID
	u.CreatedAt = user.CreatedAt
	u.DeletedAt = user.DeletedAt
	u.BlockedUntil = user.BlockedUntil
}

func (d *userDB) Get(id string) (models.User, error) {
	obj := User{}
	if err := d.db.Take(&obj, "id = ?", id).Error; err != nil {
		return models.User{}, err
	}

	return obj.into(), nil
}

func (d *userDB) LastApprovedUserForm(user models.User) (models.UserForm, error) {
	form := UserForm{}
	err := d.db.
		Order("created_at desc").
		Take(&form, "user_id = ? AND state = ?", user.ID, models.UserFormStateApproved).
		Error

	if err != nil {
		return models.UserForm{}, fmt.Errorf("take: %w", convertError(err))
	}

	return form.into(), nil
}

func (d *userDB) Pagination(config ports.UserPaginationConfig) (models.UsersConnection, error) {
	query := d.db.Model(&User{})
	query, err := paginationQueryByCreatedAtDesc(query, config.First, config.After)
	if err != nil {
		return models.UsersConnection{}, fmt.Errorf("pagination: %w", err)
	}

	var objs []models.User
	if err := query.Find(&objs).Error; err != nil {
		return models.UsersConnection{}, fmt.Errorf("find: %w", err)
	}

	if len(objs) == 0 {
		return models.UsersConnection{
			PageInfo: &models.PageInfo{},
			Edges:    make([]*models.UsersConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false
	if config.First != nil {
		hasNextPage = len(objs) > *config.First
		objs = objs[:len(objs)-1]
	}

	edges := make([]*models.UsersConnectionEdge, 0, len(objs))
	for _, obj := range objs {
		node := obj

		edges = append(edges, &models.UsersConnectionEdge{
			Cursor: obj.ID,
			Node:   &node,
		})
	}

	return models.UsersConnection{
		PageInfo: &models.PageInfo{
			HasNextPage: hasNextPage,
			StartCursor: &objs[0].ID,
			EndCursor:   &objs[len(objs)-1].ID,
		},
		Edges: edges,
	}, nil
}

func (d *userDB) Create(user *models.User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}
	u := User{}
	u.copy(user)
	if err := d.db.Create(&u).Error; err != nil {
		return fmt.Errorf("create: %w", err)
	}

	*user = u.into()
	return nil
}

func (d *userDB) MostRelevantUserForm(user models.User) (models.UserForm, error) {
	form := UserForm{}
	err := approvedOrFirstUserFormFilter(&d.db).Take(&form, "user_id = ?", user.ID).Error
	if err != nil {
		return models.UserForm{}, fmt.Errorf("take: %w", convertError(err))
	}

	return form.into(), nil
}
