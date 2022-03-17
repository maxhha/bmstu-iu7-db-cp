package db

import "gorm.io/gorm"

type UserFormState string

const (
	UserFormStateCreated    UserFormState = "CREATED"
	UserFormStateModerating UserFormState = "MODERATING"
	UserFormStateDeclained  UserFormState = "DECLAINED"
	UserFormStateApproved   UserFormState = "APPROVED"
)

type UserForm struct {
	gorm.Model
	ID            string `gorm:"default:generated();"`
	State         UserFormState
	Name          *string
	Password      *string
	Phone         *string
	Email         *string
	DeclainReason *string
}
