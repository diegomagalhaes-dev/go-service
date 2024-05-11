// Package order provides support for describing the ordering of data.
package order

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/diegomagalhaes-dev/go-service/foundation/validate"
)

const (
	ASC  = "ASC"
	DESC = "DESC"
)

var directions = map[string]string{
	ASC:  "ASC",
	DESC: "DESC",
}

type By struct {
	Field     string
	Direction string
}

func NewBy(field string, direction string) By {
	return By{
		Field:     field,
		Direction: direction,
	}
}

func Parse(r *http.Request, defaultOrder By) (By, error) {
	v := r.URL.Query().Get("orderBy")

	if v == "" {
		return defaultOrder, nil
	}

	orderParts := strings.Split(v, ",")

	var by By
	switch len(orderParts) {
	case 1:
		by = NewBy(strings.Trim(orderParts[0], " "), ASC)
	case 2:
		by = NewBy(strings.Trim(orderParts[0], " "), strings.Trim(orderParts[1], " "))
	default:
		return By{}, validate.NewFieldsError(v, errors.New("unknown order field"))
	}

	if _, exists := directions[by.Direction]; !exists {
		return By{}, validate.NewFieldsError(v, fmt.Errorf("unknown direction: %s", by.Direction))
	}

	return by, nil
}
