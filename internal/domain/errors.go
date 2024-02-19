package domain

import (
	"errors"
	"fmt"
)

var (
	//ProductNotFoundError product not found in warehouse.
	ProductNotFoundError = errors.New("product not found in warehouse")
)

// NotEnoughProductQuantityError represents an error when
// there are not enough products less than requested.
type NotEnoughProductQuantityError struct {
	ProductCode string
	Shortage    uint
}

func (e *NotEnoughProductQuantityError) Error() string {

	return fmt.Sprintf(
		"product (code: %s) less than requested (shortage: %v)",
		e.ProductCode, e.Shortage,
	)
}
