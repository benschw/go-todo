package service

import (
	"github.com/benschw/go-todo/api"
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
	var json api.Todo

	if c.EnsureBody(&json) {
		tx, err := tr.db.Begin()
		if err != nil {
			log.Print(err)
			c.JSON(500, gin.H{"error": "database error"})
			return
		}

		json.Status = "todo"
		json.Created = int32(time.Now().Unix())
		resp, err := tx.Exec(
			`INSERT INTO Todo (created, status, title, description) VALUES (?, ?, ?, ?)`,
			json.Created, json.Status, json.Title, json.Description,
		)
		if err != nil {
			log.Print(err)
			c.JSON(500, gin.H{"error": "database error"})
			return
		}
		tx.Commit()

		id, err := resp.LastInsertId()
		if err != nil {
			log.Print(err)
			c.JSON(500, gin.H{"error": "database error"})
			return
		}

		todo, err := tr.queryForTodo(int(id))
		if err != nil {
			log.Print(err)
			c.JSON(500, gin.H{"error": "database error"})
			return
		}

		c.JSON(201, todo)
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

	var todos = make([]api.Todo, 0)

	for rows.Next() {
		rows.Scan(&id, &created, &status, &title, &description)

		todos = append(todos, api.Todo{Id: id, Created: created, Status: status, Title: title, Description: description})
	}

	c.JSON(200, todos)
}

func (tr *TodoResource) GetTodo(c *gin.Context) {
	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		c.JSON(500, gin.H{"error": "input error"})
		return
	}
	todo, err := tr.queryForTodo(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	c.JSON(200, todo)
}

func (tr *TodoResource) UpdateTodo(c *gin.Context) {
	var json api.Todo

	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		c.JSON(500, gin.H{"error": "input error"})
		return
	}

	if c.EnsureBody(&json) {
		tx, err := tr.db.Begin()
		if err != nil {
			log.Print(err)
			c.JSON(500, gin.H{"error": "database error"})
			return
		}

		_, err = tx.Exec(
			`UPDATE Todo SET status = ?, title = ?, description = ? WHERE id = ?`,
			json.Status, json.Title, json.Description, id,
		)
		if err != nil {
			log.Print(err)
			c.JSON(500, gin.H{"error": "database error"})
			return
		}
		tx.Commit()

		todo, err := tr.queryForTodo(id)
		if err != nil {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}

		c.JSON(200, todo)
	}
}

func (tr *TodoResource) PatchTodo(c *gin.Context) {
	var json []api.Patch

	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		c.JSON(500, gin.H{"error": "input error"})
		return
	}

	// this is a hack because Gin falsely claims my unmarshalled obj is invalid.
	// recovering from the panic and using my object that already has the json body bound to it.
	defer func() {
		if r := recover(); r != nil {
			log.Printf("%+v", json)

			if json[0].Op != "replace" && json[0].Path != "/status" {
				c.JSON(422, gin.H{"error": "PATCH support is limited and can only replace the /status path"})
				return
			}
			tx, err := tr.db.Begin()
			if err != nil {
				log.Print(err)
				c.JSON(500, gin.H{"error": "database error"})
				return
			}
			_, err = tx.Exec(
				`UPDATE Todo SET status = ? WHERE id = ?`,
				json[0].Value, id,
			)
			if err != nil {
				log.Print(err)
				c.JSON(500, gin.H{"error": "database error"})
				return
			}
			tx.Commit()

			todo, err := tr.queryForTodo(id)
			if err != nil {
				c.JSON(404, gin.H{"error": "not found"})
				return
			}

			c.JSON(200, todo)
		}
	}()
	c.Bind(&json)

}

func (tr *TodoResource) DeleteTodo(c *gin.Context) {
	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		c.JSON(500, gin.H{"error": "input error"})
		return
	}

	tx, err := tr.db.Begin()
	if err != nil {
		log.Print(err)
		c.JSON(500, gin.H{"error": "database error"})
		return
	}

	_, err = tx.Exec(`DELETE FROM Todo WHERE id = ?`, id)

	if err != nil {
		log.Print(err)
		c.JSON(500, gin.H{"error": "database error"})
		return
	}
	tx.Commit()

	c.Data(204, "application/json", make([]byte, 0))
}

func (tr *TodoResource) queryForTodo(id int) (api.Todo, error) {
	var (
		created     int32
		status      string
		title       string
		description string
	)

	err := tr.db.
		QueryRow("SELECT created, status, title, description FROM Todo WHERE id = ?", id).
		Scan(&created, &status, &title, &description)

	if err != nil {
		return api.Todo{}, err
	}

	return api.Todo{Id: int32(id), Created: created, Status: status, Title: title, Description: description}, nil

}
