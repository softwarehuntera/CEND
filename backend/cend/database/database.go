package database

import (
	"cend/database/collection"
)

type DB struct {
	collections map[string]collection.Collection
}
func New() {
	collection.New()
}