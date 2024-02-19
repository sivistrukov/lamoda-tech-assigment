package uow

import (
	"database/sql"
	"lamoda-tech-assigment/internal/adapters/postgresql"
	"lamoda-tech-assigment/internal/domain"
)

type WarehouseRepo interface {
	Add(warehouse *domain.Warehouse) error
	Get(id uint) (domain.Warehouse, error)
	Update(warehouse domain.Warehouse) error
}

type ProductRepo interface {
	Add(product *domain.Product) error
}

type UnitOfWork struct {
	Products   ProductRepo
	Warehouses WarehouseRepo
	db         *sql.DB
	tx         *sql.Tx
}

func New(db *sql.DB) *UnitOfWork {
	tx, _ := db.Begin()

	warehouse := postgresql.NewWarehouseRepository(tx)
	product := postgresql.NewProductRepository(tx)

	return &UnitOfWork{
		Products:   product,
		Warehouses: warehouse,
		db:         db,
		tx:         tx,
	}
}

func (uow *UnitOfWork) Commit() error {
	err := uow.tx.Commit()

	uow.updateTx()

	return err
}

func (uow *UnitOfWork) Rollback() error {
	err := uow.tx.Rollback()

	uow.updateTx()

	return err
}

func (uow *UnitOfWork) updateTx() {
	tx, _ := uow.db.Begin()
	uow.tx = tx

	uow.Warehouses = postgresql.NewWarehouseRepository(tx)
	uow.Products = postgresql.NewProductRepository(tx)
}
