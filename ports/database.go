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

type OfferField string

const (
	OfferFieldCreatedAt OfferField = "created_at"
)

type TokenField string

const (
	TokenFieldCreatedAt TokenField = "created_at"
)

type TransactionField string

const (
	TransactionFieldID TransactionField = "id"
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
	UserIDs   []string
	IDs       []string
	Actions   []models.TokenAction
	OrderBy   TokenField
	OrderDesc bool
}

type OfferTakeConfig struct {
	Filter    *models.OffersFilter
	OrderBy   OfferField
	OrderDesc bool
}

type TransactionTakeConfig struct {
	Filter    *models.TransactionsFilter
	OrderBy   TransactionField
	OrderDesc bool
}

type TransactionFindConfig struct {
	Filter    *models.TransactionsFilter
	OrderBy   TransactionField
	OrderDesc bool
	Limit     int
}

type AccountDB interface {
	Get(id string) (models.Account, error)
	Take(config AccountTakeConfig) (models.Account, error)
	Create(account *models.Account) error
	Pagination(config AccountPaginationConfig) (models.AccountsConnection, error)
	UserPagination(config AccountPaginationConfig) (models.UserAccountsConnection, error)
	LockFull(account *models.Account) error
	GetAvailableMoney(account models.Account) (map[models.CurrencyEnum]models.Money, error)
}

type AuctionDB interface {
	Get(id string) (models.Auction, error)
	Take(filter *models.AuctionsFilter) (models.Auction, error)
	Create(auction *models.Auction) error
	Update(auction *models.Auction) error
	Pagination(first *int, after *string, filter *models.AuctionsFilter) (models.AuctionsConnection, error)
	LockShare(auction *models.Auction) error
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

type OfferDB interface {
	Get(id string) (models.Offer, error)
	Take(config OfferTakeConfig) (models.Offer, error)
	Create(offer *models.Offer) error
	Update(offer *models.Offer) error
	Pagination(first *int, after *string, filter *models.OffersFilter) (models.OffersConnection, error)
}

type TransactionDB interface {
	Get(id int) (models.Transaction, error)
	Take(config TransactionTakeConfig) (models.Transaction, error)
	Find(config TransactionFindConfig) ([]models.Transaction, error)
	Create(tr *models.Transaction) error
	Update(tr *models.Transaction) error
	Pagination(first *int, after *string, filter *models.TransactionsFilter) (models.TransactionsConnection, error)
}

type DB interface {
	Account() AccountDB
	Auction() AuctionDB
	Bank() BankDB
	User() UserDB
	Product() ProductDB
	UserForm() UserFormDB
	Role() RoleDB
	Token() TokenDB
	Offer() OfferDB
	Transaction() TransactionDB
	Tx() TXDB
}

type TXDB interface {
	DB() DB
	Rollback()
	Commit() error
}
