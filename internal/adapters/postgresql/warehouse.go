package postgresql

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"lamoda-tech-assigment/internal/domain"
	"log"
	"strings"
)

type ProductStatus string

const (
	AvailableStatus   = "available"
	ReservationStatus = "reservation"
)

type WarehouseRepository struct {
	tx      *sql.Tx
	builder sq.StatementBuilderType
}

type difference struct {
	Added   []domain.Product
	Removed []string
	Updated map[string]uint
}

func NewWarehouseRepository(tx *sql.Tx) *WarehouseRepository {

	return &WarehouseRepository{
		tx:      tx,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *WarehouseRepository) Add(new *domain.Warehouse) error {
	stmt, args, err := r.builder.Insert("warehouses").
		Columns("name", "is_available").
		Values(new.Name, new.IsAvailable).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return err
	}

	row := r.tx.QueryRow(stmt, args...)

	err = row.Scan(&new.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *WarehouseRepository) Get(id uint) (domain.Warehouse, error) {
	stmt, args, err := r.builder.
		Select(
			"warehouses.id as wr_id",
			"warehouses.name as wr_name",
			"warehouses.is_available as wr_is_available",
			"products.code",
			"products.name as p_name",
			"products.size",
			"wp.quantity",
			"wp.status",
		).
		From("warehouses").
		InnerJoin("warehouse_products as wp ON warehouses.id = wp.warehouse_id").
		InnerJoin("products ON products.code = wp.product_code").
		Where(sq.Eq{"warehouse_id": id}).
		ToSql()

	var warehouse domain.Warehouse

	rows, err := r.tx.Query(stmt, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	counter := 0
	products := make(map[string]domain.Product)
	forDelivery := make(map[string]domain.Product)
	for rows.Next() {
		var status ProductStatus
		counter++

		var product domain.Product

		err = rows.Scan(
			&warehouse.ID, &warehouse.Name, &warehouse.IsAvailable,
			&product.Code, &product.Name, &product.Size, &product.Quantity, &status,
		)
		if err != nil {
			return domain.Warehouse{}, err
		}

		switch status {
		case AvailableStatus:
			products[product.Code] = product
			break
		case ReservationStatus:
			forDelivery[product.Code] = product
			break
		}
	}

	if counter < 1 {
		return domain.Warehouse{}, EntityNotFoundError
	}

	warehouse.SetProducts(products)
	warehouse.SetProductsForDelivery(forDelivery)

	return warehouse, nil
}

func (r *WarehouseRepository) Update(warehouse domain.Warehouse) error {

	stmt, args, err := r.builder.Update("warehouses").
		Set("name", warehouse.Name).
		Set("is_available", warehouse.IsAvailable).
		Where(sq.Eq{"id": warehouse.ID}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.tx.Exec(stmt, args...)
	if err != nil {
		return err
	}

	old, err := r.Get(warehouse.ID)

	err = r.updateWarehouseProducts(warehouse.ID, AvailableStatus, old.Products(), warehouse.Products())
	if err != nil {
		return err
	}

	err = r.updateWarehouseProducts(warehouse.ID, ReservationStatus, old.ProductsForDelivery(), warehouse.ProductsForDelivery())
	if err != nil {
		return err
	}

	return nil
}

func (r *WarehouseRepository) updateWarehouseProducts(
	warehouseID uint, productStatus string, old map[string]domain.Product, updated map[string]domain.Product,
) error {
	diff := getProductsDiff(old, updated)

	// Updated
	for code, quantity := range diff.Updated {
		stmt, args, err := r.builder.Update("warehouse_products").
			Set("quantity", quantity).
			Where(sq.Eq{"warehouse_id": warehouseID}).
			Where(sq.Eq{"product_code": code}).
			Where(sq.Eq{"status": productStatus}).
			ToSql()
		if err != nil {
			return err
		}

		_, err = r.tx.Exec(stmt, args...)
		if err != nil {
			return err
		}
	}

	// Removed
	for _, code := range diff.Removed {
		stmt, args, err := r.builder.Update("warehouse_products").
			Set("quantity", 0).
			Where(sq.Eq{"warehouse_id": warehouseID}).
			Where(sq.Eq{"product_code": code}).
			Where(sq.Eq{"status": productStatus}).
			ToSql()
		if err != nil {
			return err
		}

		_, err = r.tx.Exec(stmt, args...)
		if err != nil {
			return err
		}
	}

	// Added
	if len(diff.Added) > 0 {
		values := make([]any, 0, len(diff.Added)*3)
		for _, product := range diff.Added {
			values = append(values, warehouseID)
			values = append(values, product.Code)
			values = append(values, product.Quantity)
			values = append(values, productStatus)
		}

		stmt, args, err := r.builder.Insert("warehouse_products").
			Columns("warehouse_id", "product_code", "quantity", "status").
			Values(values...).
			ToSql()
		if err != nil {
			return err
		}

		_, err = r.tx.Exec(stmt, args...)
		if err != nil {
			if strings.Contains(err.Error(), "violates foreign key constraint") {
				return ForeignKeyError
			}
			return err
		}
	}

	return nil
}

func getProductsDiff(old map[string]domain.Product, updated map[string]domain.Product) difference {
	diff := difference{
		Added:   make([]domain.Product, 0),
		Removed: make([]string, 0),
		Updated: make(map[string]uint),
	}

	// Added
	for code, product := range updated {
		if _, exist := old[code]; !exist {
			diff.Added = append(diff.Added, product)
		}
	}

	// Removed
	for code := range old {
		if _, exist := updated[code]; !exist {
			diff.Removed = append(diff.Removed, code)
		}
	}

	// Updated
	for code, product := range updated {
		if oldProduct, exist := old[code]; exist {
			if oldProduct != product {
				diff.Updated[code] = product.Quantity
			}
		}
	}

	return diff
}
