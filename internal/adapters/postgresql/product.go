package postgresql

import (
	"database/sql"
	"lamoda-tech-assigment/internal/domain"
	"strings"

	sq "github.com/Masterminds/squirrel"
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
		if strings.Contains(err.Error(), "unique constraint") {
			return EntityAlreadyExist
		}
		return err
	}

	return nil
}
