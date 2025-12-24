package products

import "time"

type ProductID int64
type GroupID int64

type Group struct {
	ID   GroupID
	Name string
}

type Product struct {
	ID          ProductID
	Name        string
	GroupID     *GroupID
	GroupName   string
	Quantity    int
	MinQuantity int
	Missing     bool
	UpdatedAt   time.Time
}

type ProductFilter struct {
	OnlyMissingOrLow bool
}

type NewProduct struct {
	Name    string
	GroupID *GroupID
}
