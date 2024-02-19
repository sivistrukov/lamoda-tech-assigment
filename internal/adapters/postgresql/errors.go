package postgresql

import "errors"

var (
	EntityNotFoundError = errors.New("entity not found")
	ForeignKeyError     = errors.New("violates foreign key constraint")
	EntityAlreadyExist  = errors.New("entity already exist")
)
