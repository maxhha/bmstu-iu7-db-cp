package graph

import (
	"auction-back/db"
	"auction-back/graph/model"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"crypto/sha256"

	"github.com/hashicorp/go-multierror"
	"gorm.io/gorm"
)

var passwordHashSalt []byte

func InitPasswordHashSalt() {
	key, ok := os.LookupEnv("PASSWORD_HASH_SALT")
	if !ok {
		panic("PASSWORD_HASH_SALT does not exist in environment variables!")
	}

	passwordHashSalt = []byte(key)
}

func getOrCreateUserDraftForm(DB *gorm.DB, viewer *db.User) (db.UserForm, error) {
	form := db.UserForm{}

	err := DB.Order("created_at desc").Take(&form, "user_id = ?", viewer.ID).Error

	if err == nil {
		if form.State == db.UserFormStateCreated || form.State == db.UserFormStateDeclained {
			return form, nil
		} else if form.State == db.UserFormStateApproved {
			// clone form
			form.ID = ""
			form.State = db.UserFormStateCreated
			form.CreatedAt = time.Time{}
			form.UpdatedAt = time.Time{}
		} else if form.State == db.UserFormStateModerating {
			return form, fmt.Errorf("moderating")
		} else {
			return form, fmt.Errorf("unknown form state: %s", form.State)
		}
	} else if err != gorm.ErrRecordNotFound {
		return form, fmt.Errorf("take: %w", err)
	}

	// Create new form only if no form exists
	// or there is approved form for copy
	form.UserID = viewer.ID
	if err = DB.Create(&form).Error; err != nil {
		return form, fmt.Errorf("create: %w", err)
	}

	return form, nil
}

func hashPassword(password string) (string, error) {
	h := sha256.New()

	if _, err := h.Write(passwordHashSalt); err != nil {
		return "", err
	}

	if _, err := h.Write([]byte(password)); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(h.Sum(nil)), nil
}

func checkHashAndPassword(hash string, password string) bool {
	hash2, _ := hashPassword(password)
	return hash == hash2
}

func (r *userResolver) isOwnerOrManager(viewer *db.User, obj *db.User) error {
	if viewer == nil {
		return fmt.Errorf("unauthorized")
	}

	if obj == nil {
		return fmt.Errorf("user is nil")
	}

	var errors error

	if viewer.ID != obj.ID {
		errors = multierror.Append(errors, fmt.Errorf("viewer is not owner"))
	} else {
		return nil
	}

	if err := r.RolePort.HasRole(db.RoleTypeManager, viewer); err != nil {
		errors = multierror.Append(errors, err)
	} else {
		return nil
	}

	return errors
}

// Creates pagination for users
func UserPagination(query *gorm.DB, first *int, after *string) (*model.UsersConnection, error) {
	query, err := PaginationByCreatedAtDesc(query, first, after)

	if err != nil {
		return nil, fmt.Errorf("pagination: %w", err)
	}

	var objs []db.User
	if err := query.Find(&objs).Error; err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}

	if len(objs) == 0 {
		return &model.UsersConnection{
			PageInfo: &model.PageInfo{},
			Edges:    make([]*model.UsersConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(objs) > *first
		objs = objs[:len(objs)-1]
	}

	edges := make([]*model.UsersConnectionEdge, 0, len(objs))

	for _, obj := range objs {
		node := obj

		edges = append(edges, &model.UsersConnectionEdge{
			Cursor: obj.ID,
			Node:   &node,
		})
	}

	return &model.UsersConnection{
		PageInfo: &model.PageInfo{
			HasNextPage: hasNextPage,
			StartCursor: &objs[0].ID,
			EndCursor:   &objs[len(objs)-1].ID,
		},
		Edges: edges,
	}, nil
}
