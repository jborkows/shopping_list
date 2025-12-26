package shoppinglist

import (
	"time"

	"shopping/internal/domain/products"
)

type ItemID int64

type Item struct {
	ID        ItemID
	Name      string
	ProductID *products.ProductID
	IconKey   string
	GroupName string
	Quantity  float64
	Unit      products.Unit
	Done      bool
	CreatedAt time.Time
}
