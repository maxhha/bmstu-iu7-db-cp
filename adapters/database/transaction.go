package database

import (
	"auction-back/models"
	"auction-back/ports"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type transactionDB struct{ *Database }

func (d *Database) Transaction() ports.TransactionDB { return &transactionDB{d} }

type Transaction struct {
	ID            int `gorm:"default:serial();"`
	Date          *time.Time
	State         models.TransactionState `gorm:"default:CREATED"`
	Type          models.TransactionType
	Currency      models.CurrencyEnum
	Amount        decimal.Decimal
	Error         *string
	AccountFromID string
	AccountToID   string
	OfferID       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func (t *Transaction) into() models.Transaction {
	transaction := models.Transaction{
		ID:            t.ID,
		Date:          t.Date,
		State:         t.State,
		Type:          t.Type,
		Currency:      t.Currency,
		Amount:        t.Amount,
		Error:         t.Error,
		AccountFromID: t.AccountFromID,
		AccountToID:   t.AccountToID,
		OfferID:       t.OfferID,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}

	if t.DeletedAt.Valid {
		transaction.DeletedAt = &t.DeletedAt.Time
	}

	return transaction
}

func (t *Transaction) copy(transaction *models.Transaction) {
	if transaction == nil {
		return
	}

	t.ID = transaction.ID
	t.Date = transaction.Date
	t.State = transaction.State
	t.Type = transaction.Type
	t.Currency = transaction.Currency
	t.Amount = transaction.Amount
	t.Error = transaction.Error
	t.AccountFromID = transaction.AccountFromID
	t.AccountToID = transaction.AccountToID
	t.OfferID = transaction.OfferID
	t.CreatedAt = transaction.CreatedAt
	t.UpdatedAt = transaction.UpdatedAt

	if transaction.DeletedAt == nil {
		t.DeletedAt.Valid = false
	} else {
		t.DeletedAt.Time = *transaction.DeletedAt
		t.DeletedAt.Valid = true
	}
}

var transactionFieldToColumn = map[ports.TransactionField]string{
	ports.TransactionFieldID: "id",
}

func (d *transactionDB) Get(id int) (models.Transaction, error) {
	tr := Transaction{}
	if err := d.db.Take(&tr, "id = ?", id).Error; err != nil {
		return models.Transaction{}, fmt.Errorf("take: %w", convertError(err))
	}

	return tr.into(), nil
}

func (d *transactionDB) filter(query *gorm.DB, config *models.TransactionsFilter) *gorm.DB {
	if config == nil {
		return query
	}

	if len(config.IDs) > 0 {
		query = query.Where("id IN ?", config.IDs)
	}

	if config.DateRange != nil {
		if config.DateRange.From != nil {
			query = query.Where("date >= ?", config.DateRange.From)
		}

		if config.DateRange.To != nil {
			query = query.Where("date < ?", config.DateRange.To)
		}
	}

	if len(config.States) > 0 {
		query = query.Where("state IN ?", config.States)
	}

	if len(config.Types) > 0 {
		query = query.Where("type IN ?", config.Types)
	}

	if len(config.Currencies) > 0 {
		query = query.Where("currency IN ?", config.Currencies)
	}

	if len(config.AccountFormIDs) > 0 {
		query = query.Where("account_from_id IN ?", config.AccountFormIDs)
	}

	if len(config.AccountToIDs) > 0 {
		query = query.Where("account_to_id IN ?", config.AccountToIDs)
	}

	if len(config.OfferIDs) > 0 {
		query = query.Where("offer_id IN ?", config.OfferIDs)
	}

	return query
}

func (d *transactionDB) order(query *gorm.DB, orderBy ports.TransactionField, orderDesc bool) (*gorm.DB, error) {
	if orderBy == "" {
		return query, nil
	}

	column, ok := transactionFieldToColumn[orderBy]
	if !ok {
		return nil, fmt.Errorf("unknown field '%s'", orderBy)
	}

	query = query.Order(clause.OrderByColumn{
		Column: clause.Column{Name: column},
		Desc:   orderDesc,
	})

	return query, nil
}

func (d *transactionDB) Take(config ports.TransactionTakeConfig) (models.Transaction, error) {
	query := d.filter(d.db, config.Filter)
	query, err := d.order(query, config.OrderBy, config.OrderDesc)
	if err != nil {
		return models.Transaction{}, fmt.Errorf("order: %w", err)
	}

	tr := Transaction{}
	if err := query.Take(&tr).Error; err != nil {
		return models.Transaction{}, fmt.Errorf("take: %w", convertError(err))
	}

	return tr.into(), nil
}

func (d *transactionDB) Find(config ports.TransactionFindConfig) ([]models.Transaction, error) {
	query := d.filter(d.db, config.Filter)
	query, err := d.order(query, config.OrderBy, config.OrderDesc)
	if err != nil {
		return nil, fmt.Errorf("order: %w", err)
	}

	if config.Limit > 0 {
		query = query.Limit(config.Limit)
	}

	var objs []Transaction
	if err := query.Find(&objs).Error; err != nil {
		return nil, fmt.Errorf("find: %w", convertError(err))
	}

	arr := make([]models.Transaction, 0, len(objs))
	for _, obj := range objs {
		arr = append(arr, obj.into())
	}

	return arr, nil
}

func (d *transactionDB) Create(transaction *models.Transaction) error {
	if transaction == nil {
		return ports.ErrTransactionIsNil
	}
	t := Transaction{}
	t.copy(transaction)
	if err := d.db.Create(&t).Error; err != nil {
		return fmt.Errorf("create: %w", convertError(err))
	}

	*transaction = t.into()
	return nil
}

func (d *transactionDB) Update(form *models.Transaction) error {
	if form == nil {
		return ports.ErrTransactionIsNil
	}

	f := Transaction{}
	f.copy(form)

	if err := d.db.Save(&f).Error; err != nil {
		return fmt.Errorf("save: %w", convertError(err))
	}
	*form = f.into()

	return nil
}

func (d *transactionDB) Pagination(first *int, after *string, filter *models.TransactionsFilter) (models.TransactionsConnection, error) {
	query := d.filter(d.db.Model(&Transaction{}), filter)
	query, err := paginationQueryByCreatedAtDesc(query, first, after)
	if err != nil {
		return models.TransactionsConnection{}, fmt.Errorf("pagination: %w", err)
	}

	var objs []Transaction
	if err := query.Find(&objs).Error; err != nil {
		return models.TransactionsConnection{}, fmt.Errorf("find: %w", err)
	}

	if len(objs) == 0 {
		return models.TransactionsConnection{
			PageInfo: &models.PageInfo{},
			Edges:    make([]*models.TransactionsConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if first != nil {
		hasNextPage = len(objs) > *first
		objs = objs[:len(objs)-1]
	}

	edges := make([]*models.TransactionsConnectionEdge, 0, len(objs))

	for _, obj := range objs {
		node := obj.into()

		edges = append(edges, &models.TransactionsConnectionEdge{
			Cursor: fmt.Sprintf("%d", node.ID),
			Node:   &node,
		})
	}

	startCursor := fmt.Sprintf("%d", objs[0].ID)
	endCursor := fmt.Sprintf("%d", objs[len(objs)-1].ID)

	return models.TransactionsConnection{
		PageInfo: &models.PageInfo{
			HasNextPage: hasNextPage,
			StartCursor: &startCursor,
			EndCursor:   &endCursor,
		},
		Edges: edges,
	}, nil
}
