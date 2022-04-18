package database

import (
	"auction-back/models"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

//go:generate go run ../../codegen/gormdbops/main.go --out product_gen.go --model Product --methods Get,Update,Create,Pagination

type Product struct {
	ID            string              `gorm:"default:generated();"`
	State         models.ProductState `gorm:"default:'CREATED';"`
	Title         string
	Description   string
	CreatorID     string
	DeclainReason *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

func (p *Product) into() models.Product {
	return models.Product{
		ID:            p.ID,
		State:         p.State,
		Title:         p.Title,
		Description:   p.Description,
		CreatorID:     p.CreatorID,
		DeclainReason: p.DeclainReason,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
		DeletedAt:     sql.NullTime(p.DeletedAt),
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
	p.DeclainReason = product.DeclainReason
	p.CreatedAt = product.CreatedAt
	p.UpdatedAt = product.UpdatedAt
	p.DeletedAt = gorm.DeletedAt(product.DeletedAt)
}

func (d *productDB) filter(query *gorm.DB, config *models.ProductsFilter) *gorm.DB {
	if config == nil {
		return query
	}

	if len(config.OwnerIDs) > 0 {
		query = query.Joins(
			"JOIN ( ? ) ofd ON products.id = ofd.product_id AND ofd.owner_n = 1 AND ofd.owner_id IN ?",
			d.ownersNumberedQuery(query),
			config.OwnerIDs,
		)
	}
	return query
}

func (d *productDB) GetCreator(p models.Product) (models.User, error) {
	return d.User().Get(p.CreatorID)
}

// Query all owners of queried products.
// Returns select query with fields: product_id, owner_id, from_date
func (d *productDB) ownersQuery(query *gorm.DB) *gorm.DB {
	creators := query.Session(&gorm.Session{Initialized: true}).
		Model(&Product{}).
		Select("id as product_id, creator_id as owner_id, created_at as from_date")

	buyers := query.Session(&gorm.Session{Initialized: true}).
		Model(&Product{}).
		Select("products.id as product_id, auctions.buyer_id as owner_id, auctions.finished_at as from_date").
		Joins("JOIN auctions ON auctions.product_id = products.id AND auctions.state = ?", models.AuctionStateSucceeded)

	return d.db.Raw("? UNION ALL ?", creators, buyers)
}

// Query current owner of products
// Returns select query with fields: product_id, owner_id, from_date, owner_n
func (d *productDB) ownersNumberedQuery(query *gorm.DB) *gorm.DB {
	query = d.ownersQuery(query)

	ownersNumbered := d.db.Table("( ? ) as ofd", query)

	ownersNumbered = ownersNumbered.Select(
		"*, ROW_NUMBER() OVER(PARTITION BY ofd.product_id ORDER BY ofd.from_date DESC) as owner_n",
	)

	return ownersNumbered
}

func (d *productDB) GetOwner(p models.Product) (models.User, error) {
	query := d.db.Model(&Product{}).Where("products.id = ?", p.ID)
	query = d.ownersNumberedQuery(query)

	owner := User{}
	if err := d.db.Model(&owner).
		Joins("JOIN ( ? ) ofd ON users.id = ofd.owner_id AND ofd.owner_n = 1", query).
		Take(&owner).Error; err != nil {
		return models.User{}, fmt.Errorf("take owner: %w", convertError(err))
	}

	return owner.into(), nil
}
