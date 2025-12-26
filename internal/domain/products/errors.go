package products

import "errors"

var (
	ErrNameRequired                 = errors.New("name required")
	ErrInvalidUnit                  = errors.New("invalid unit")
	ErrQuantityMustBeNonNegative    = errors.New("quantity must be >= 0")
	ErrMinQuantityMustBeNonNegative = errors.New("min quantity must be >= 0")
)
