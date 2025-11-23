package estate

import (
	"example-app/pkg/store"
	"time"
)

type Helper interface {
	GetNameTable() string
	GetNameColumns() string
	GetPlaceholder() string
	GetValues() []interface{}
}
type Property struct {
	ID        int       `json:"id" db:"id"`
	Address   string    `json:"address" db:"address"`
	Type      string    `json:"type" db:"type"`
	Price     float64   `json:"price" db:"price"`
	OwnerID   int       `json:"owner_id" db:"owner_id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}
type Purchase struct {
	ID           int       `json:"id" db:"id"`
	PropertyID   int       `json:"property_id" db:"property_id"`
	SellerID     int       `json:"seller_id" db:"seller_id"`
	PurchaseDate time.Time `json:"purchase_date" db:"purchase_date"`
	InitialPrice float64   `json:"initial_price" db:"initial_price"`
	OwnerID      int       `json:"owner_id" db:"owner_id"`
}

type Sale struct {
	ID         int       `json:"id" db:"id"`
	PropertyID int       `json:"property_id" db:"property_id"`
	BuyerID    int       `json:"buyer_id" db:"buyer_id"`
	SaleDate   time.Time `json:"sale_date" db:"sale_date"`
	FinalPrice float64   `json:"final_price" db:"final_price"`
	OwnerID    int       `json:"owner_id" db:"owner_id"`
}
type User struct {
	store.User
}

func (p Property) GetNameTable() string {
	return "properties"
}
func (p Property) GetNameColumns() string {
	return "address, type, price, status"
}
func (p Property) GetPlaceholder() string {
	return "$1, $2, $3, $4"
}
func (p Property) GetValues() []interface{} {
	return []interface{}{
		p.Address, p.Type, p.Price, p.Status,
	}
}
func (p Purchase) GetNameTable() string {
	return "purchases"
}
func (p Purchase) GetNameColumns() string {
	return "property_id, seller_id, initial_price"
}
func (p Purchase) GetPlaceholder() string {
	return "$1, $2, $3"
}
func (p Purchase) GetValues() []interface{} {
	return []interface{}{
		p.PropertyID, p.SellerID, p.InitialPrice,
	}
}
func (s Sale) GetNameTable() string {
	return "sales"
}
func (s Sale) GetNameColumns() string {
	return "property_id, buyer_id, final_price"
}
func (s Sale) GetPlaceholder() string {
	return "$1, $2, $3"
}
func (s Sale) GetValues() []interface{} {
	return []interface{}{
		s.PropertyID, s.BuyerID, s.FinalPrice,
	}
}
func (u User) GetNameTable() string {
	return "users"
}
func (u User) GetNameColumns() string {
	return "role_id"
}
func (u User) GetPlaceholder() string {
	return "$1"
}
func (u User) GetValues() []interface{} {
	return []interface{}{
		u.RoleID,
	}
}
