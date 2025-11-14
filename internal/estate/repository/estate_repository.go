package repository

import (
	"auth-service/internal/estate/model"
	"context"

	"github.com/jmoiron/sqlx"
)

type EstateRepository interface {
	CreateProperty(ctx context.Context, p *model.Property) (int, error)
	GetPropertyByID(ctx context.Context, id int) (*model.Property, error)
	ListProperties(ctx context.Context) ([]model.Property, error)
	CreateSale(ctx context.Context, s *model.Sale) (int, error)
	CreatePurchase(ctx context.Context, p *model.Purchase) (int, error)
}

type SQLEstateRepository struct {
	db *sqlx.DB
}

func NewSQLEstateRepository(db *sqlx.DB) *SQLEstateRepository {
	return &SQLEstateRepository{db: db}
}

func (r *SQLEstateRepository) CreateProperty(ctx context.Context, p *model.Property) (int, error) {
	var propertyID int
	query := `INSERT INTO properties (address, type, price, agent_id, status)
              VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, p.Address, p.Type, p.Price, p.AgentID, p.Status).Scan(&propertyID)
	if err != nil {
		return 0, err
	}
	return propertyID, nil
}

func (r *SQLEstateRepository) GetPropertyByID(ctx context.Context, id int) (*model.Property, error) {
	var p model.Property
	query := `SELECT id, address, type, price, agent_id, status, created_at, updated_at 
              FROM properties WHERE id = $1`

	err := r.db.GetContext(ctx, &p, query, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *SQLEstateRepository) ListProperties(ctx context.Context) ([]model.Property, error) {
	var properties []model.Property
	query := `SELECT id, address, type, price, agent_id, status, created_at, updated_at 
              FROM properties ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &properties, query)
	if err != nil {
		return nil, err
	}
	return properties, nil
}

func (r *SQLEstateRepository) CreateSale(ctx context.Context, s *model.Sale) (int, error) {
	var saleID int
	query := `INSERT INTO sales (property_id, buyer_id, final_price, agent_id)
              VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, s.PropertyID, s.BuyerID, s.FinalPrice, s.AgentID).Scan(&saleID)
	if err != nil {
		return 0, err
	}
	return saleID, nil
}

func (r *SQLEstateRepository) CreatePurchase(ctx context.Context, p *model.Purchase) (int, error) {
	var purchaseID int
	query := `INSERT INTO purchases (property_id, seller_id, initial_price, agent_id)
              VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRowContext(ctx, query, p.PropertyID, p.SellerID, p.InitialPrice, p.AgentID).Scan(&purchaseID)
	if err != nil {
		return 0, err
	}
	return purchaseID, nil
}
