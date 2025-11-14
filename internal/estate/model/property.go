package model

import (
	"time"
)

// Property (Объекты недвижимости)
type Property struct {
	ID        int       `json:"id" db:"id"`
	Address   string    `json:"address" db:"address"`
	Type      string    `json:"type" db:"type"` // Квартира, Дом, Земля
	Price     float64   `json:"price" db:"price"`
	AgentID   int       `json:"agent_id" db:"agent_id"` // Агент, ответственный за объект
	Status    string    `json:"status" db:"status"`     // available, sold, purchased
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Sale (Продажи)
type Sale struct {
	ID         int       `json:"id" db:"id"`
	PropertyID int       `json:"property_id" db:"property_id"`
	BuyerID    int       `json:"buyer_id" db:"buyer_id"`
	SaleDate   time.Time `json:"sale_date" db:"sale_date"`
	FinalPrice float64   `json:"final_price" db:"final_price"`
	AgentID    int       `json:"agent_id" db:"agent_id"`
}

// Purchase (Покупки)
type Purchase struct {
	ID           int       `json:"id" db:"id"`
	PropertyID   int       `json:"property_id" db:"property_id"`
	SellerID     int       `json:"seller_id" db:"seller_id"`
	PurchaseDate time.Time `json:"purchase_date" db:"purchase_date"`
	InitialPrice float64   `json:"initial_price" db:"initial_price"`
	AgentID      int       `json:"agent_id" db:"agent_id"`
}
