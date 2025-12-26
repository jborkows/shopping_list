package shoppinglist

import "errors"

var (
	ErrNameRequired           = errors.New("name required")
	ErrQuantityMustBePositive = errors.New("quantity must be > 0")
)
