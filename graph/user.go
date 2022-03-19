package graph

import (
	"auction-back/db"
	"encoding/base64"
	"fmt"
	"os"

	"crypto/sha256"

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

func (r *mutationResolver) getOrCreateUserForm(viewer *db.User) (db.UserForm, error) {
	form := db.UserForm{}

	err := r.DB.Order("created_at desc").Take(&form, "user_id = ?", viewer.ID).Error

	if err == nil {
		if form.State == db.UserFormStateCreated || form.State == db.UserFormStateDeclained {
			return form, nil
		} else if form.State == db.UserFormStateApproved {
			// create new form duplicating previos one
			form.ID = ""
			form.State = db.UserFormStateCreated
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
	if err = r.DB.Create(&form).Error; err != nil {
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
