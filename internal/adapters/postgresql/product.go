package postgresql

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"lamoda-tech-assigment/internal/domain"
)

type ProductRepository struct {
	tx      *sql.Tx
	builder sq.StatementBuilderType
}

func NewProductRepository(tx *sql.Tx) *ProductRepository {

	return &ProductRepository{
		tx:      tx,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *ProductRepository) Add(product *domain.Product) error {
	stmt, args, err := r.builder.Insert("products").
		Columns("code", "name", "size").
		Values(product.Code, product.Name, product.Size).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.tx.Exec(stmt, args...)
	if err != nil {
		return err
	}

	return nil
}
