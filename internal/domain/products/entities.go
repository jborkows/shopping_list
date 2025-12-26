package products

import "time"

type ProductID int64
type GroupID int64
type Unit string

const (
	UnitKG    Unit = "kg"
	UnitLiter Unit = "litr"
	UnitPiece Unit = "sztuk"
	UnitGram  Unit = "gramy"
)

const MaxProductsPageSize int64 = 30

type Group struct {
	ID   GroupID
	Name string
}

type Product struct {
	ID          ProductID
	Name        string
	IconKey     string
	GroupID     *GroupID
	GroupName   string
	Quantity    float64
	Unit        Unit
	MinQuantity float64
	Missing     bool
	UpdatedAt   time.Time
}

type ProductFilter struct {
	OnlyMissingOrLow bool
	NameQuery        string
	GroupIDs         []GroupID
	Limit            int64
	Offset           int64
}

type NewProduct struct {
	Name        string
	IconKey     string
	GroupID     *GroupID
	Quantity    float64
	Unit        Unit
	MinQuantity float64
}
