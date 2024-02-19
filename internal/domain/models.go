package domain

// Product is representation of product in warehouse.
type Product struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Size     string `json:"size"`
	Quantity uint   `json:"quantity"`
}

// Warehouse is representation of warehouse with stored products.
type Warehouse struct {
	ID                  uint               `json:"id"`
	Name                string             `json:"name"`
	IsAvailable         bool               `json:"isAvailable"`
	products            map[string]Product // key is product code
	productsForDelivery map[string]Product // key is product code
}

// Products returns map of Product where key is product code.
func (w *Warehouse) Products() map[string]Product {
	product := make(map[string]Product, len(w.products))
	for k, v := range w.products {
		product[k] = v
	}
	return product
}

// SetProducts is setter for Warehouse products
func (w *Warehouse) SetProducts(products map[string]Product) {
	p := make(map[string]Product, len(products))
	for k, v := range products {
		p[k] = v
	}

	w.products = p
}

// ProductsForDelivery returns map of reserved for delivery
// Product where key is product code.
func (w *Warehouse) ProductsForDelivery() map[string]Product {
	product := make(map[string]Product, len(w.productsForDelivery))
	for k, v := range w.productsForDelivery {
		product[k] = v
	}
	return product
}

// SetProductsForDelivery is setter for Warehouse products for delivery
func (w *Warehouse) SetProductsForDelivery(products map[string]Product) {
	p := make(map[string]Product, len(products))
	for k, v := range products {
		p[k] = v
	}

	w.productsForDelivery = p
}

// GetProduct returns product or error if Product not found.
func (w *Warehouse) GetProduct(code string) (Product, error) {
	product, exist := w.products[code]

	if !exist {
		return Product{}, ProductNotFoundError
	}

	return product, nil
}

// GetReservedForDelivery returns reserved for delivery Product
// or error if product not found.
func (w *Warehouse) GetReservedForDelivery(code string) (Product, error) {
	reserved, exist := w.productsForDelivery[code]

	if !exist {
		return Product{}, ProductNotFoundError
	}

	return reserved, nil
}

// AddProducts adds a new products to the warehouse or increase
// quantity if product already exists in stock.
func (w *Warehouse) AddProducts(products ...Product) {
	for _, newProduct := range products {
		product, exist := w.products[newProduct.Code]
		if exist {
			product.Quantity += newProduct.Quantity
		} else {
			product = newProduct
		}
		w.products[newProduct.Code] = product
	}
}

// ReserveProduct reserves products in the warehouse for delivery,
// returns an error if the operation failed.
func (w *Warehouse) ReserveProduct(code string, quantity uint) error {
	product, exist := w.products[code]
	if !exist || product.Quantity == 0 {
		return ProductNotFoundError
	}

	if product.Quantity < quantity {
		return &NotEnoughProductQuantityError{
			ProductCode: code,
			Shortage:    quantity - product.Quantity,
		}
	}
	product.Quantity -= quantity
	w.products[code] = product

	reserved, exist := w.productsForDelivery[code]
	if exist {
		reserved.Quantity += quantity
	} else {
		reserved = Product{
			Code:     product.Code,
			Name:     product.Name,
			Size:     product.Size,
			Quantity: quantity,
		}
	}
	w.productsForDelivery[code] = reserved

	return nil
}

// CancelProductReservation cancels a product reservation for
// a specific quantity of product or for all if the quantity is not specified.
func (w *Warehouse) CancelProductReservation(code string, quantity ...uint) error {
	reserved, exist := w.productsForDelivery[code]
	if !exist || reserved.Quantity == 0 {
		return ProductNotFoundError
	}

	var cancelQuantity uint
	if len(quantity) > 0 {
		cancelQuantity = quantity[0]

		if cancelQuantity > reserved.Quantity {
			return &NotEnoughProductQuantityError{
				ProductCode: code,
				Shortage:    cancelQuantity - reserved.Quantity,
			}
		}

	} else {
		cancelQuantity = reserved.Quantity
	}

	reserved.Quantity -= cancelQuantity
	w.productsForDelivery[code] = reserved

	product, exist := w.products[code]
	if exist {
		product.Quantity += cancelQuantity
	} else {
		product = Product{
			Code:     reserved.Code,
			Name:     reserved.Name,
			Quantity: cancelQuantity,
			Size:     product.Size,
		}
	}
	w.products[code] = product

	return nil
}
