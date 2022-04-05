package graph

import (
	"auction-back/models"
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

func getOrCreateUserDraftForm(DB *gorm.DB, viewer *models.User) (models.UserForm, error) {
	form := models.UserForm{}

	err := DB.Order("created_at desc").Take(&form, "user_id = ?", viewer.ID).Error

	if err == nil {
		if form.State == models.UserFormStateCreated || form.State == models.UserFormStateDeclained {
			return form, nil
		} else if form.State == models.UserFormStateApproved {
			// clone form
			form.ID = ""
			form.State = models.UserFormStateCreated
			form.CreatedAt = time.Time{}
			form.UpdatedAt = time.Time{}
		} else if form.State == models.UserFormStateModerating {
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

func (r *userResolver) isOwnerOrManager(viewer *models.User, obj *models.User) error {
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

	if err := r.RolePort.HasRole(models.RoleTypeManager, viewer); err != nil {
		errors = multierror.Append(errors, err)
	} else {
		return nil
	}

	return errors
}

// Creates pagination for users
func UserPagination(query *gorm.DB, first *int, after *string) (*models.UsersConnection, error) {
	query, err := PaginationQueryByCreatedAtDesc(query, first, after)

	if err != nil {
		return nil, fmt.Errorf("pagination: %w", err)
	}

	var objs []models.User
	if err := query.Find(&objs).Error; err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}

	if len(objs) == 0 {
		return &models.UsersConnection{
			PageInfo: &models.PageInfo{},
			Edges:    make([]*models.UsersConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(objs) > *first
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

	return &models.UsersConnection{
		PageInfo: &models.PageInfo{
			HasNextPage: hasNextPage,
			StartCursor: &objs[0].ID,
			EndCursor:   &objs[len(objs)-1].ID,
		},
		Edges: edges,
	}, nil
}
