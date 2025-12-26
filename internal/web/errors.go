package web

import (
	"errors"
	"net/http"
	"strings"

	"shopping/internal/domain/products"
	"shopping/internal/domain/shoppinglist"
)

func (s *Server) writeDBError(w http.ResponseWriter, err error) {
	msg := err.Error()
	if strings.Contains(msg, "no such table") {
		http.Error(w, "database schema missing: run migrations in ./migrations (see README.md)", http.StatusInternalServerError)
		return
	}
	http.Error(w, msg, http.StatusInternalServerError)
}

func (s *Server) writeUserError(w http.ResponseWriter, err error) {
	http.Error(w, userErrorMessage(err), http.StatusBadRequest)
}

func userErrorMessage(err error) string {
	switch {
	case errors.Is(err, products.ErrNameRequired), errors.Is(err, shoppinglist.ErrNameRequired):
		return "Nazwa jest wymagana."
	case errors.Is(err, products.ErrInvalidUnit):
		return "Nieprawidłowa jednostka."
	case errors.Is(err, products.ErrQuantityMustBeNonNegative):
		return "Ilość musi być większa lub równa 0."
	case errors.Is(err, products.ErrMinQuantityMustBeNonNegative):
		return "Minimalna ilość musi być większa lub równa 0."
	case errors.Is(err, shoppinglist.ErrQuantityMustBePositive):
		return "Ilość musi być większa od 0."
	default:
		// Last resort: do not lose context in the UI.
		return err.Error()
	}
}
