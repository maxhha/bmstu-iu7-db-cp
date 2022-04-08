package database

import (
	"auction-back/models"
	"auction-back/ports"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type productDB struct{ *Database }

func (d *Database) Product() ports.ProductDB { return &productDB{d} }

type Product struct {
	ID          string              `gorm:"default:generated();"`
	State       models.ProductState `gorm:"default:'CREATED';"`
	Title       string
	Description string
	CreatorID   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
}

func (p *Product) into() models.Product {
	return models.Product{
		ID:          p.ID,
		State:       p.State,
		Title:       p.Title,
		Description: p.Description,
		CreatorID:   p.CreatorID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		DeletedAt:   p.DeletedAt,
	}
}

func (p *Product) copy(product *models.Product) {
	if product == nil {
		return
	}
	p.ID = product.ID
	p.State = product.State
	p.Title = product.Title
	p.Description = product.Description
	p.CreatorID = product.CreatorID
	p.CreatedAt = product.CreatedAt
	p.UpdatedAt = product.UpdatedAt
	p.DeletedAt = product.DeletedAt
}

func (d *productDB) Create(product *models.Product) error {
	if product == nil {
		return fmt.Errorf("product is nil")
	}
	p := Product{}
	p.copy(product)
	if err := d.db.Create(&p).Error; err != nil {
		return fmt.Errorf("create: %w", err)
	}

	*product = p.into()
	return nil
}

func (d *productDB) Pagination(config ports.ProductPaginationConfig) (models.ProductsConnection, error) {
	query := d.db.Model(&Product{})

	query, err := paginationQueryByCreatedAtDesc(query, config.First, config.After)

	if err != nil {
		return models.ProductsConnection{}, fmt.Errorf("pagination: %w", err)
	}

	var products []models.Product
	result := query.Find(&products)

	if result.Error != nil {
		return models.ProductsConnection{}, result.Error
	}

	if len(products) == 0 {
		return models.ProductsConnection{
			PageInfo: &models.PageInfo{},
			Edges:    make([]*models.ProductsConnectionEdge, 0),
		}, nil
	}

	hasNextPage := false

	if config.First != nil {
		hasNextPage = len(products) > *config.First
		products = products[:len(products)-1]
	}

	edges := make([]*models.ProductsConnectionEdge, 0, len(products))

	for _, node := range products {
		edges = append(edges, &models.ProductsConnectionEdge{
			Cursor: node.ID,
			Node:   &node,
		})
	}

	return models.ProductsConnection{
		PageInfo: &models.PageInfo{
			HasNextPage: hasNextPage,
			StartCursor: &products[0].ID,
			EndCursor:   &products[len(products)-1].ID,
		},
		Edges: edges,
	}, nil
}

func (d *productDB) GetCreator(p models.Product) (models.User, error) {
	return d.User().Get(p.CreatorID)
}

func (d *productDB) GetOwner(p models.Product) (models.User, error) {
	lastAuctionFinishedAtQuery := d.db.Model(&Auction{}).
		Select("MAX(finished_at)").
		Where("product_id = ?", p.ID)

	lastBuyerQuery := d.db.Model(&Auction{}).
		Select("buyer_id").
		Where("finished_at = ?", lastAuctionFinishedAtQuery)

	owner := User{}
	err := d.db.Take(&owner, "id IN ?", lastBuyerQuery).Error
	if err == nil {
		return owner.into(), err
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, fmt.Errorf("take owner: %w", err)
	}

	creator, err := d.User().Get(p.CreatorID)

	if err != nil {
		return creator, fmt.Errorf("ensure creator: %w", err)
	}

	return creator, nil
}
