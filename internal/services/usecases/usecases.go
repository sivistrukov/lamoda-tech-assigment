package usecases

import (
	"lamoda-tech-assigment/internal/domain"
	"lamoda-tech-assigment/internal/services/uow"
)

type IUseCases interface {
	ReserveProductsInWarehouse(warehouseID uint, products ...domain.Product) error
	CancelReservationInWarehouse(warehouseID uint, products ...domain.Product) error
	AddWarehouse(name string) (domain.Warehouse, error)
	AddProduct(code string, name string, size string) (domain.Product, error)
	AddProductToWarehouse(warehouseID uint, products ...domain.Product) ([]map[string]any, error)
	GetWarehouseProducts(warehouseID uint) ([]domain.Product, error)
	GetProductQuantityInWarehouse(warehouseID uint, code string) (uint, error)
}

type UseCases struct {
	uow *uow.UnitOfWork
}

// New returns pointer to new UseCases structure
func New(uow *uow.UnitOfWork) *UseCases {

	return &UseCases{uow: uow}
}

func (uc *UseCases) ReserveProductsInWarehouse(warehouseID uint, products ...domain.Product) error {
	warehouse, err := uc.uow.Warehouses.Get(warehouseID)
	if err != nil {
		return err
	}

	for _, product := range products {
		err = warehouse.ReserveProduct(product.Code, product.Quantity)
		if err != nil {
			_ = uc.uow.Rollback()
			return err
		}
	}

	err = uc.uow.Warehouses.Update(warehouse)
	if err != nil {
		_ = uc.uow.Rollback()
		return err
	}

	err = uc.uow.Commit()
	if err != nil {
		_ = uc.uow.Rollback()
		return err
	}

	return nil
}

func (uc *UseCases) CancelReservationInWarehouse(warehouseID uint, products ...domain.Product) error {
	warehouse, err := uc.uow.Warehouses.Get(warehouseID)
	if err != nil {
		return err
	}

	for _, product := range products {
		err = warehouse.CancelProductReservation(product.Code, product.Quantity)
		if err != nil {
			_ = uc.uow.Rollback()
			return err
		}
	}

	err = uc.uow.Warehouses.Update(warehouse)
	if err != nil {
		_ = uc.uow.Rollback()
		return err
	}

	err = uc.uow.Commit()
	if err != nil {
		_ = uc.uow.Rollback()
		return err
	}

	return nil
}

// AddWarehouse create new warehouse and save it in database
func (uc *UseCases) AddWarehouse(name string) (domain.Warehouse, error) {
	warehouse := domain.Warehouse{
		Name:        name,
		IsAvailable: true,
	}

	err := uc.uow.Warehouses.Add(&warehouse)
	if err != nil {
		_ = uc.uow.Rollback()
		return domain.Warehouse{}, err
	}

	err = uc.uow.Commit()
	if err != nil {
		_ = uc.uow.Rollback()
		return domain.Warehouse{}, err
	}
	return warehouse, nil
}

// AddProduct create new product and save it in database
func (uc *UseCases) AddProduct(code string, name string, size string) (domain.Product, error) {
	product := domain.Product{
		Code: code,
		Name: name,
		Size: size,
	}

	err := uc.uow.Products.Add(&product)
	if err != nil {
		_ = uc.uow.Rollback()
		return domain.Product{}, err
	}

	err = uc.uow.Commit()
	if err != nil {
		_ = uc.uow.Rollback()
		return domain.Product{}, err
	}
	return product, nil
}

// AddProductToWarehouse add product to warehouse and save it in database.
func (uc *UseCases) AddProductToWarehouse(
	warehouseID uint, products ...domain.Product,
) ([]map[string]any, error) {
	warehouse, err := uc.uow.Warehouses.Get(warehouseID)
	if err != nil {
		_ = uc.uow.Rollback()
		return make([]map[string]any, 0), err
	}

	warehouse.AddProducts(products...)

	err = uc.uow.Warehouses.Update(warehouse)
	if err != nil {
		_ = uc.uow.Rollback()
		return make([]map[string]any, 0), err
	}

	err = uc.uow.Commit()
	if err != nil {
		_ = uc.uow.Rollback()
		return make([]map[string]any, 0), err
	}

	updated := warehouse.Products()
	result := make([]map[string]any, 0, len(updated))
	for _, p := range updated {
		result = append(result, map[string]any{"code": p.Code, "quantity": p.Quantity})
	}

	return result, err
}

// GetWarehouseProducts returns list of available products in warehouse.
func (uc *UseCases) GetWarehouseProducts(warehouseID uint) ([]domain.Product, error) {
	warehouse, err := uc.uow.Warehouses.Get(warehouseID)
	if err != nil {
		return []domain.Product{}, err
	}

	result := productMapToSlice(warehouse.Products())

	return result, nil
}

// GetProductQuantityInWarehouse returns quantity of product in warehouse.
func (uc *UseCases) GetProductQuantityInWarehouse(warehouseID uint, code string) (uint, error) {
	warehouse, err := uc.uow.Warehouses.Get(warehouseID)
	if err != nil {
		return 0, err
	}

	product, err := warehouse.GetProduct(code)
	if err != nil {
		return 0, err
	}

	return product.Quantity, nil
}

func productMapToSlice(products map[string]domain.Product) []domain.Product {
	result := make([]domain.Product, 0, len(products))

	for _, product := range products {
		result = append(result, product)
	}

	return result
}
