package database

import (
	"cend/database/collection"
)

type DB struct {
	collections map[string]collection.Collection
}

func New() *DB {
	collections := make(map[string]collection.Collection)
	return &DB{collections: collections}
}