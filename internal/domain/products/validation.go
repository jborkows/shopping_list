package products

import (
	"strings"
)

func NormalizeGroupName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrNameRequired
	}
	return name, nil
}

func NormalizeProductName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ErrNameRequired
	}
	return name, nil
}

func NormalizeUnit(u Unit) (Unit, error) {
	u = Unit(strings.TrimSpace(string(u)))
	switch u {
	case UnitKG, UnitLiter, UnitPiece, UnitGram:
		return u, nil
	default:
		return "", ErrInvalidUnit
	}
}
