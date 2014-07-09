package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

type DB struct {
	*sql.DB
}
type Tx struct {
	*sql.Tx
}

// Open returns a DB reference for a data source.
func Open(dataSourceName string) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// Begin starts an returns a new transaction.
func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

func main() {
	var (
		dbUser     = "root"
		dbPassword = ""
		dbHost     = "localhost"
		dbName     = "Todo"
	)
	db, err := Open(dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":3306)/" + dbName)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	todoResource := &TodoResource{db: db}

	r := gin.Default()

	r.GET("/todo", todoResource.GetAllTodos)
	r.GET("/todo/:id", todoResource.GetTodo)
	r.POST("/todo", todoResource.CreateTodo)

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
