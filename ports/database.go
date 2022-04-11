package ports

import "auction-back/models"

//go:generate go run ../codegen/portmocks/main.go --config ../portmocksgen.yml --in database.go --out database_mock.go --outpkg ports

type UserFormField string

const (
	UserFormFieldCreatedAt     UserFormField = "created_at"
	UserFormFieldState         UserFormField = "state"
	UserFormFieldEmail         UserFormField = "email"
	UserFormFieldPhone         UserFormField = "phone"
	UserFormFieldPassword      UserFormField = "password"
	UserFormFieldDeclainReason UserFormField = "declain_reason"
)

type AccountPaginationConfig struct {
	First   *int
	After   *string
	UserIDs []string
}

type ProductPaginationConfig struct {
	Filter models.ProductsFilter
	First  *int
	After  *string
}

type UserFormPaginationConfig struct {
	models.UserFormsFilter
	First *int
	After *string
}

type UserFormTakeConfig struct {
	models.UserFormsFilter
	OrderBy   UserFormField
	OrderDesc bool
}

type UserPaginationConfig struct {
	models.UsersFilter
	First *int
	After *string
}

type BankTakeConfig struct {
	IDs   []string
	Names []string
}

type AccountTakeConfig struct {
	UserIDs []string
}

type RoleFindConfig struct {
	UserIDs []string
	Types   []models.RoleType
	Limit   int
}

type TokenTakeConfig struct {
	UserIDs []string
	IDs     []string
}

type AccountDB interface {
	Take(config AccountTakeConfig) (models.Account, error)
	Create(account *models.Account) error
	Pagination(config AccountPaginationConfig) (models.AccountsConnection, error)
	UserPagination(config AccountPaginationConfig) (models.UserAccountsConnection, error)
}

type BankDB interface {
	Get(id string) (models.Bank, error)
	Take(config BankTakeConfig) (models.Bank, error)
	GetAccount(bank models.Bank) (models.BankAccount, error)
}

type UserDB interface {
	Get(id string) (models.User, error)
	Create(user *models.User) error
	Pagination(config UserPaginationConfig) (models.UsersConnection, error)
	LastApprovedUserForm(user models.User) (models.UserForm, error)
	MostRelevantUserForm(user models.User) (models.UserForm, error)
}

type UserFormDB interface {
	Get(id string) (models.UserForm, error)
	Take(config UserFormTakeConfig) (models.UserForm, error)
	Create(form *models.UserForm) error
	Update(form *models.UserForm) error
	Pagination(config UserFormPaginationConfig) (models.UserFormsConnection, error)
	GetLoginForm(input models.LoginInput) (models.UserForm, error)
}

type ProductDB interface {
	Get(id string) (models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	Pagination(config ProductPaginationConfig) (models.ProductsConnection, error)
	GetOwner(product models.Product) (models.User, error)
	GetCreator(product models.Product) (models.User, error)
}

type RoleDB interface {
	Find(config RoleFindConfig) ([]models.Role, error)
}

type TokenDB interface {
	Take(config TokenTakeConfig) (models.Token, error)
	Create(token *models.Token) error
	Update(token *models.Token) error
	GetUser(token models.Token) (models.User, error)
}

type DB interface {
	Account() AccountDB
	Bank() BankDB
	User() UserDB
	Product() ProductDB
	UserForm() UserFormDB
	Role() RoleDB
	Token() TokenDB
}
