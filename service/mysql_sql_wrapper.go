package service

/**
 * Wrapping sql interface & mysql driver implementation to keep db driver decision encapsulated
 *
 * explanation: https://medium.com/@benbjohnson/structuring-applications-in-go-3b04be4ff091
 */

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

// Open returns a DB reference for a data source.
func Open(dataSourceName string) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
