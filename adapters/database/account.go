package database

import (
	"auction-back/models"
	"auction-back/ports"
	"database/sql"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type accountDB struct{ *Database }

func (d *Database) Account() ports.AccountDB { return &accountDB{d} }

type Account struct {
	ID        string `gorm:"default:generated();"`
	Type      models.AccountType
	UserID    string
	BankID    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (a *Account) into() models.Account {
	return models.Account{
		ID:        a.ID,
		Type:      a.Type,
		UserID:    a.UserID,
		BankID:    a.BankID,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		DeletedAt: sql.NullTime(a.DeletedAt),
	}
}

func (a *Account) copy(account *models.Account) {
	if account == nil {
		return
	}
	a.ID = account.ID
	a.Type = account.Type
	a.UserID = account.UserID
	a.BankID = account.BankID
	a.CreatedAt = account.CreatedAt
	a.UpdatedAt = account.UpdatedAt
	a.DeletedAt = gorm.DeletedAt(account.DeletedAt)
}

func (a *Account) String() string {
	return fmt.Sprintf("Account[id=%s]", a.ID)
}

func (d *accountDB) Get(id string) (models.Account, error) {
	obj := Account{}
	if err := d.db.Take(&obj, "id = ?", id).Error; err != nil {
		return models.Account{}, fmt.Errorf("take: %w", convertError(err))
	}

	return obj.into(), nil
}

func (d *Database) findAccountsForPagination(config ports.AccountPaginationConfig) ([]models.Account, error) {
	query := d.db.Model(&models.Account{})

	if len(config.UserIDs) > 0 {
		query = query.Where("user_id IN ?", config.UserIDs)
	}

	query, err := paginationQueryByCreatedAtDesc(query, config.First, config.After)
	if err != nil {
		return nil, fmt.Errorf("pagination: %w", err)
	}

	var objs []Account
	if err := query.Find(&objs).Error; err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}

	converted := make([]models.Account, 0, len(objs))
	for _, obj := range objs {
		converted = append(converted, obj.into())
	}

	return converted, nil
}

// Creates pagination for accounts
func (d *Database) Pagination(config ports.AccountPaginationConfig) (models.AccountsConnection, error) {
	objs, err := d.findAccountsForPagination(config)

	if err != nil {
		return models.AccountsConnection{}, fmt.Errorf("find for pagination: %w", err)
	}

	if len(objs) == 0 {
		return models.AccountsConnection{
			PageInfo: &models.PageInfo{},
			Edges:    make([]*models.AccountsConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if config.First != nil {
		hasNextPage = len(objs) > *config.First
		objs = objs[:len(objs)-1]
	}

	edges := make([]*models.AccountsConnectionEdge, 0, len(objs))
	var errors error

	for _, obj := range objs {
		node, err := obj.ConcreteType()

		if err == nil {
			edges = append(edges, &models.AccountsConnectionEdge{
				Cursor: obj.ID,
				Node:   node,
			})
		} else {
			errors = multierror.Append(
				errors,
				fmt.Errorf("%v concrete type: %w", obj, err),
			)
		}
	}

	return models.AccountsConnection{
		PageInfo: &models.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: false,
			StartCursor:     &objs[0].ID,
			EndCursor:       &objs[len(objs)-1].ID,
		},
		Edges: edges,
	}, errors
}

// Creates pagination for accounts
func (d *Database) UserPagination(config ports.AccountPaginationConfig) (models.UserAccountsConnection, error) {
	objs, err := d.findAccountsForPagination(config)

	if err != nil {
		return models.UserAccountsConnection{}, err
	}

	if len(objs) == 0 {
		return models.UserAccountsConnection{
			PageInfo: &models.PageInfo{},
			Edges:    make([]*models.UserAccountsConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if config.First != nil {
		hasNextPage = len(objs) > *config.First
		objs = objs[:len(objs)-1]
	}

	edges := make([]*models.UserAccountsConnectionEdge, 0, len(objs))
	var errors error

	for _, obj := range objs {
		account, err := obj.ConcreteType()
		if err != nil {
			errors = multierror.Append(errors, err)
		}

		switch account := account.(type) {
		case models.UserAccount:
			edges = append(edges, &models.UserAccountsConnectionEdge{
				Cursor: obj.ID,
				Node:   &account,
			})
		default:
			errors = multierror.Append(
				errors,
				fmt.Errorf("unexpected user account type: %s", obj.Type))
		}
	}

	return models.UserAccountsConnection{
		PageInfo: &models.PageInfo{
			HasNextPage: hasNextPage,
			StartCursor: &objs[0].ID,
			EndCursor:   &objs[len(objs)-1].ID,
		},
		Edges: edges,
	}, errors
}

func (d *accountDB) Create(account *models.Account) error {
	if account == nil {
		return fmt.Errorf("account is nil")
	}
	a := Account{}
	a.copy(account)
	if err := d.db.Create(&a).Error; err != nil {
		return fmt.Errorf("create: %w", err)
	}

	*account = a.into()
	return nil
}

func (d *accountDB) Take(config ports.AccountTakeConfig) (models.Account, error) {
	query := d.db

	if len(config.UserIDs) > 0 {
		query = query.Where("user_id IN ?", config.UserIDs)
	}

	account := Account{}
	if err := query.Take(&account).Error; err != nil {
		return models.Account{}, fmt.Errorf("take: %w", convertError(err))
	}

	return account.into(), nil
}

func (d *accountDB) LockFull(account *models.Account) error {
	if account == nil {
		return ports.ErrAccountIsNil
	}
	obj := Account{}
	err := d.db.Clauses(clause.Locking{
		Strength: "UPDATE",
		Table:    clause.Table{Name: clause.CurrentTable},
	}).
		Take(&obj, "id = ?", account.ID).
		Error
	if err != nil {
		return convertError(err)
	}

	*account = obj.into()
	return nil
}

// var countingTransactionStates = []models.TransactionState{
// 	models.TransactionStateSucceeded,
// 	models.TransactionStateProcessing,
// 	models.TransactionStateError,
// }

// func (d *accountDB) availableMoneyQuery(query *gorm.DB) *gorm.DB {
// 	toTrs := d.db.Model(&Transaction{}).
// 		Select("currency, account_to_id as account_id, amount").
// 		Joins(
// 			"JOIN ( ? ) a ON account_to_id = a.id AND transaction.state IN ?",
// 			query.Session(&gorm.Session{Initialized: true}).Model(&Account{}),
// 			countingTransactionStates,
// 		)

// 	fromTrs := d.db.Model(&Transaction{}).
// 		Select("currency, account_from_id as account_id, -amount").
// 		Joins(
// 			"JOIN ( ? ) a ON account_from_id = a.id AND ( state IN ? OR ( type IN ? AND state IN ? ) )",
// 			query.Session(&gorm.Session{Initialized: true}).Model(&Account{}),
// 			countingTransactionStates,
// 			[]models.TransactionType{
// 				models.TransactionTypeBuy,
// 				// TODO: Add models.TransactionTypeBuyFee
// 			},
// 			[]models.TransactionState{
// 				models.TransactionStateCreated,
// 			},
// 		)

// 	allTrs := d.db.Raw("? UNION ALL ?", fromTrs, toTrs)

// 	moneyQuery := d.db.
// 		Select("a.currency, a.account_id, SUM(a.amount) as amount").
// 		Table("? a", allTrs).
// 		Group("a.currency, a.account_id")

// 	return moneyQuery
// }

func (d *accountDB) GetAvailableMoney(account models.Account) (map[models.CurrencyEnum]models.Money, error) {
	// TODO: Test this
	// query := d.availableMoneyQuery(d.db.Model(&Account{}).Where("id = ?", account.ID))

	// var moneys []models.Money
	// if err := query.Scan(&moneys).Error; err != nil {
	// 	return nil, fmt.Errorf("scan: %w", err)
	// }

	// moneysMap := make(map[models.CurrencyEnum]models.Money, len(moneys))
	// for _, m := range moneys {
	// 	moneysMap[m.Currency] = m
	// }

	// return moneysMap, nil
	moneysMap := make(map[models.CurrencyEnum]models.Money, 1)
	moneysMap[models.CurrencyEnumRub] = models.Money{
		Currency: models.CurrencyEnumRub,
		Amount:   decimal.NewFromFloatWithExponent(100, -2),
	}

	return moneysMap, nil
}
