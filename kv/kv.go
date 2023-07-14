package kv

import (
	"errors"
	"fmt"
)

type DB struct {
	Database map[string]string
}

func NewDB() *DB {
	return &DB{
		Database: make(map[string]string),
	}
}

func (db *DB) Get(key string) (string, error) {
	val, ok := db.Database[key]
	if ok {
		return val, nil
	}
	return "", fmt.Errorf("not found")
}

func (db *DB) Delete(key string) error {
	if _, ok := db.Database[key]; !ok {
		return errors.New("key not found")
	}
	delete(db.Database, key)
	return nil
}

func (db *DB) Put(key string, val string) (bool, error) {
	_, ok := db.Database[key]
	db.Database[key] = val
	if ok {
		return true, nil
	}
	return false, nil
}

func (db *DB) Clear() {
	db.Database = make(map[string]string)
}
