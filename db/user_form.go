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
	ID            string `json:"id" gorm:"default:generated();"`
	UserID        string
	User          User
	State         UserFormState `gorm:"default:'CREATED';"`
	Name          *string       `json:"name"`
	Password      *string
	Phone         *string `json:"phone"`
	Email         *string `json:"email"`
	DeclainReason *string
}
