package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"time"
)

type TodoResource struct {
	db *DB
}

func (tr *TodoResource) CreateTodo(c *gin.Context) {
	var json TodoJSON

	if c.EnsureBody(&json) {
		tx, err := tr.db.Begin()
		if err != nil {
			log.Print(err)
			c.JSON(500, gin.H{"error": "database error"})
			return
		}

		json.Status = "todo"
		json.Created = int32(time.Now().Unix())
		_, err = tx.Exec(
			`INSERT INTO Todo (created, status, title, description) VALUES (?, ?, ?, ?)`,
			json.Created, json.Status, json.Title, json.Description,
		)
		if err != nil {
			log.Print(err)
			c.JSON(500, gin.H{"error": "database error"})
			return
		}
		tx.Commit()

		c.JSON(200, gin.H{"title": json.Title, "description": json.Description})
	}

}
func (tr *TodoResource) GetAllTodos(c *gin.Context) {
	var (
		id          int32
		created     int32
		status      string
		title       string
		description string
	)

	rows, err := tr.db.Query("SELECT id, created, status, title, description FROM Todo ORDER BY created DESC")
	if err != nil {
		log.Print(err)
		c.JSON(500, gin.H{"error": "database error"})
		return
	}

	var todos = make([]TodoJSON, 0)

	for rows.Next() {
		rows.Scan(&id, &created, &status, &title, &description)

		todos = append(todos, TodoJSON{Id: id, Created: created, Status: status, Title: title, Description: description})
	}

	c.JSON(200, todos)

}

func (tr *TodoResource) GetTodo(c *gin.Context) {
	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		c.JSON(500, gin.H{"error": "database error"})
		return
	}

	var (
		created     int32
		status      string
		title       string
		description string
	)

	err = tr.db.
		QueryRow("SELECT created, status, title, description FROM Todo WHERE id = ?", id).
		Scan(&created, &status, &title, &description)

	todo := &TodoJSON{Id: int32(id), Created: created, Status: status, Title: title, Description: description}

	c.JSON(200, todo)

}
