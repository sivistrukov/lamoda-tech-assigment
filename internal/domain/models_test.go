package domain

import (
	"errors"
	"testing"
)

func TestWarehouse_GetProduct_Success(t *testing.T) {
	warehouse := Warehouse{
		products: map[string]Product{
			"A": {
				Code:     "A",
				Name:     "Product A",
				Quantity: 10,
			},
		},
		productsForDelivery: make(map[string]Product),
	}

	product, err := warehouse.GetProduct("A")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if product != warehouse.products["A"] {
		t.Errorf(
			"the wrong product was returned: expected %v, got %v",
			warehouse.products["A"], product,
		)
	}
}

func TestWarehouse_GetProduct_ProductNotFound(t *testing.T) {
	warehouse := Warehouse{
		products: map[string]Product{
			"A": {
				Code:     "A",
				Name:     "Product A",
				Quantity: 10,
			},
		},
		productsForDelivery: make(map[string]Product),
	}

	_, err := warehouse.GetProduct("B")
	if !errors.Is(err, ProductNotFoundError) {
		t.Errorf("unexpected error: expect %v, got %v", ProductNotFoundError, err)
	}
}

func TestWarehouse_GetReservedForDelivery_Success(t *testing.T) {
	warehouse := Warehouse{
		productsForDelivery: map[string]Product{
			"A": {
				Code:     "A",
				Name:     "Product A",
				Quantity: 10,
			},
		},
		products: make(map[string]Product),
	}

	reserved, err := warehouse.GetReservedForDelivery("A")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if reserved != warehouse.productsForDelivery["A"] {
		t.Errorf(
			"the wrong reserved was returned: expected %v, got %v",
			warehouse.products["A"], reserved,
		)
	}
}

func TestWarehouse_GetReservedForDelivery_ProductNotFound(t *testing.T) {
	warehouse := Warehouse{
		productsForDelivery: map[string]Product{
			"A": {
				Code:     "A",
				Name:     "Product A",
				Quantity: 10,
			},
		},
		products: make(map[string]Product),
	}

	_, err := warehouse.GetReservedForDelivery("B")
	if !errors.Is(err, ProductNotFoundError) {
		t.Errorf("unexpected error: expect %v, got %v", ProductNotFoundError, err)
	}
}

func TestWarehouse_AddProducts(t *testing.T) {
	warehouse := Warehouse{
		products: map[string]Product{
			"A": {
				Code:     "A",
				Name:     "Product A",
				Quantity: 10,
			},
		},
		productsForDelivery: make(map[string]Product),
	}
	productsCopy := make(map[string]Product, len(warehouse.products))
	for code, product := range warehouse.products {
		productsCopy[code] = product
	}
	newProducts := []Product{
		{
			Code:     "A",
			Name:     "Product A",
			Quantity: 5,
		},
		{
			Code:     "B",
			Name:     "Product B",
			Quantity: 10,
		},
	}

	warehouse.AddProducts(newProducts...)

	for _, newProduct := range newProducts {
		product, exist := warehouse.products[newProduct.Code]
		if !exist {
			t.Errorf("newProduct %s not found in stock", newProduct.Code)
		}

		productCopy, _ := productsCopy[newProduct.Code]
		if product.Quantity != productCopy.Quantity+newProduct.Quantity {
			t.Errorf(
				"quantity stored differs from added: %v - %v",
				product.Quantity, newProduct.Quantity,
			)
		}
	}
}

func TestWarehouse_ReserveProduct_SuccessCase(t *testing.T) {
	warehouse := Warehouse{
		products: map[string]Product{
			"A": {
				Code:     "A",
				Name:     "Product A",
				Quantity: 10,
			},
		},
		productsForDelivery: make(map[string]Product),
	}

	err := warehouse.ReserveProduct("A", 10)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if warehouse.products["A"].Quantity != 0 {
		t.Errorf("the quantity of available product in the warehouse did not decrease")
	}

	reserved, exist := warehouse.productsForDelivery["A"]
	if !exist {
		t.Errorf("the product is not listed as reserved")
	}

	if reserved.Quantity != 10 {
		t.Errorf("discrepancy in the quantity of ordered goods")
	}
}

func TestWarehouse_ReserveProduct_NotEnoughInStock(t *testing.T) {
	warehouse := Warehouse{
		products: map[string]Product{
			"A": {
				Code:     "A",
				Name:     "Product A",
				Quantity: 10,
			},
		},
		productsForDelivery: make(map[string]Product),
	}

	err := warehouse.ReserveProduct("A", 11)

	var notEnoughErr *NotEnoughProductQuantityError
	if !errors.As(err, &notEnoughErr) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWarehouse_ReserveProduct_ProductNotFound(t *testing.T) {
	warehouse := Warehouse{
		products: map[string]Product{
			"A": {
				Code:     "A",
				Name:     "Product A",
				Quantity: 10,
			},
		},
		productsForDelivery: make(map[string]Product),
	}

	err := warehouse.ReserveProduct("B", 1)
	if !errors.Is(err, ProductNotFoundError) {
		t.Errorf("unexpected error: expect %v, got %v", ProductNotFoundError, err)
	}
}

func TestWarehouse_CancelProductReservation_WithoutQuantity(t *testing.T) {
	warehouse := Warehouse{
		productsForDelivery: map[string]Product{
			"A": {
				Code:     "A",
				Name:     "Product A",
				Quantity: 10,
			},
		},
		products: make(map[string]Product),
	}

	err := warehouse.CancelProductReservation("A")

	if err != nil {
		t.Errorf("unexpected error:  %v", err)
	}

	if warehouse.productsForDelivery["A"].Quantity != 0 {
		t.Errorf("the quantity of reserved product did not decrease")
	}

	product, exits := warehouse.products["A"]
	if !exits {
		t.Errorf("the product is not listed as available in stock")
	}

	if product.Quantity != 10 {
		t.Errorf("the quantity of available product in the warehouse did not increase")
	}
}

func TestWarehouse_CancelProductReservation_WithQuantity(t *testing.T) {
	warehouse := Warehouse{
		productsForDelivery: map[string]Product{
			"A": {
				Code:     "A",
				Name:     "Product A",
				Quantity: 10,
			},
		},
		products: make(map[string]Product),
	}

	err := warehouse.CancelProductReservation("A", 5)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if warehouse.productsForDelivery["A"].Quantity != 5 {
		t.Errorf("the quantity of reserved product did not decrease")
	}

	product, exits := warehouse.products["A"]
	if !exits {
		t.Errorf("the product is not listed as available in stock")
	}

	if product.Quantity != 5 {
		t.Errorf("the quantity of available product in the warehouse did not increase")
	}
}

func TestWarehouse_CancelProductReservation_ProductNotReserved(t *testing.T) {
	warehouse := Warehouse{
		productsForDelivery: map[string]Product{
			"A": {
				Code:     "A",
				Name:     "Product A",
				Quantity: 10,
			},
		},
		products: make(map[string]Product),
	}

	err := warehouse.CancelProductReservation("B", 5)

	if !errors.Is(err, ProductNotFoundError) {
		t.Errorf("unexpected error: expect %v, got %v", ProductNotFoundError, err)
	}
}

func TestWarehouse_CancelProductReservation_CanceledMoreThenReserved(t *testing.T) {
	warehouse := Warehouse{
		productsForDelivery: map[string]Product{
			"A": {
				Code:     "A",
				Name:     "Product A",
				Quantity: 10,
			},
		},
		products: make(map[string]Product),
	}

	err := warehouse.CancelProductReservation("A", 15)

	var notEnoughErr *NotEnoughProductQuantityError
	if !errors.As(err, &notEnoughErr) {
		t.Errorf("unexpected error: %v", err)
	}
}
