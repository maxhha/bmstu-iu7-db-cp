package graph

import (
	"auction-back/models"
	"auction-back/ports"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"

	"crypto/sha256"

	"github.com/hashicorp/go-multierror"
)

var ErrNoPassword = errors.New("password not set")
var ErrPasswordMissmatch = errors.New("password mismatch")

var passwordHashSalt []byte

func InitPasswordHashSalt() {
	key, ok := os.LookupEnv("PASSWORD_HASH_SALT")
	if !ok {
		panic("PASSWORD_HASH_SALT does not exist in environment variables!")
	}

	passwordHashSalt = []byte(key)
}

var ErrUserFormModerating = errors.New("moderating")
var ErrUserNotOwner = errors.New("viewer is not owner")
var ErrUserIsNil = errors.New("user is nil")

func getOrCreateUserDraftForm(DB ports.DB, viewer models.User) (models.UserForm, error) {
	form, err := DB.UserForm().Take(ports.UserFormTakeConfig{
		OrderBy:   ports.UserFormFieldCreatedAt,
		OrderDesc: true,
		UserFormsFilter: models.UserFormsFilter{
			UserID: []string{viewer.ID},
		},
	})

	if err == nil {
		if form.IsEditable() {
			return form, nil
		} else if form.State == models.UserFormStateApproved {
			// clone form
			form.ID = ""
			form.State = models.UserFormStateCreated
			form.CreatedAt = time.Time{}
			form.UpdatedAt = time.Time{}
		} else if form.State == models.UserFormStateModerating {
			return form, ErrUserFormModerating
		} else {
			return form, fmt.Errorf("unknown form state: %s", form.State)
		}
	} else if !errors.Is(err, ports.ErrRecordNotFound) {
		return form, fmt.Errorf("take: %w", err)
	}

	// Create new form only if no form exists
	// or there is approved form for copy
	form.UserID = viewer.ID
	if err = DB.UserForm().Create(&form); err != nil {
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

func (r *userResolver) isOwnerOrManager(viewer models.User, obj *models.User) error {
	if obj == nil {
		return ErrUserIsNil
	}

	var errors error

	if viewer.ID != obj.ID {
		errors = multierror.Append(errors, ErrUserNotOwner)
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
