package postgresql

import "errors"

var EntityNotFoundError = errors.New("entity not found")
var ForeignKeyError = errors.New("violates foreign key constraint")
