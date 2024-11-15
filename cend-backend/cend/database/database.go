package database

import (
	"cend/database/collection"
	"fmt"
)

type DB struct {
	name string
	collections map[string]*collection.Collection
}

func New(name string) *DB {
	// TODO: Make this a folder of collections
	collections := make(map[string]*collection.Collection)
	return &DB{name: name, collections: collections}
}

func (db *DB) AddCollection(name string, c *collection.Collection) {
	db.collections[name] = c
}

func (db *DB) GetCollection(name string)  (*collection.Collection, error) {
	if collection, exists := db.collections[name]; exists {
		return collection, nil
	}
	return nil, fmt.Errorf("collection %s not found", name)
}
