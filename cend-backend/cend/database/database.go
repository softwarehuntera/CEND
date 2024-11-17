package database

import (
	"cend/database/collection"
	"fmt"
	"os"
	"path/filepath"
)

type DB struct {
	name string
	path string
	collections map[string]*collection.Collection
} 

func New(name string) *DB {
	// TODO: Make this a folder of collections
	collections := make(map[string]*collection.Collection)

	var dbPath string
	dbPath, exists := os.LookupEnv("DB_PATH")
	if !exists {
		dbPath = "./cend-db"
	}
	return &DB{name: name, path: dbPath, collections: collections}
}

func (db *DB) AddCollection(name string) {
	collectionPath := filepath.Join(db.path, name)
	c := collection.New(name, collectionPath)
	db.collections[name] = c
}

func (db *DB) GetCollection(name string)  (*collection.Collection, error) {
	if collection, exists := db.collections[name]; exists {
		return collection, nil
	}
	return nil, fmt.Errorf("collection %s not found", name)
}
