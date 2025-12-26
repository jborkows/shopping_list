package shoppinglist

import (
	"regexp"
	"strings"
)

var whitespaceRE = regexp.MustCompile(`\s+`)

func NormalizeItemName(name string) (string, error) {
	name = strings.TrimSpace(name)
	name = whitespaceRE.ReplaceAllString(name, " ")
	if name == "" {
		return "", ErrNameRequired
	}
	return name, nil
}
