package service

import (
	"errors"
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

	if !c.Bind(&json) {
		tr.err(c, errors.New("bad json, cannont create"), 422, "problem decoding json")
		return
	}
	json.Status = "todo"
	json.Created = int32(time.Now().Unix())
	resp, err := tr.db.Exec(
		`INSERT INTO Todo (created, status, title, description) VALUES (?, ?, ?, ?)`,
		json.Created, json.Status, json.Title, json.Description,
	)
	if err != nil {
		tr.err(c, err, 500, "database error")
		return
	}

	id, err := resp.LastInsertId()
	if err != nil {
		tr.err(c, err, 500, "database error")
		return
	}

	todo, err := tr.queryForTodo(int(id))
	if err != nil {
		tr.err(c, err, 500, "inserted todo doesn't seem to exist")
		return
	}

	c.JSON(201, todo)
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
		tr.err(c, err, 500, "database error")
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
		tr.err(c, err, 422, "problem decoding id sent")
		return
	}
	todo, err := tr.queryForTodo(id)
	if err != nil {
		c.JSON(404, api.NewError("not found"))
		return
	}

	c.JSON(200, todo)
}

func (tr *TodoResource) UpdateTodo(c *gin.Context) {
	var json api.Todo

	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		tr.err(c, err, 422, "problem decoding id sent")
		return
	}

	if !c.Bind(&json) {
		tr.err(c, errors.New("bad json, cannont update"), 422, "problem decoding json")
		return
	}
	json.Id = int32(id)

	if err = tr.execTodoUpdate(json); err != nil {
		tr.err(c, err, 404, "problem updating todo, doesn't exist?")
		return
	}

	c.JSON(200, json)
}

func (tr *TodoResource) PatchTodo(c *gin.Context) {
	var json []api.Patch

	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		tr.err(c, err, 422, "problem decoding id sent")
		return
	}

	todo, err := tr.queryForTodo(id)
	if err != nil {
		c.JSON(404, api.NewError("not found"))
		return
	}

	// this is a hack because Gin falsely claims my unmarshalled obj is invalid.
	// recovering from the panic and using my object that already has the json body bound to it.
	defer func() {
		if r := recover(); r != nil {
			if json[0].Op != "replace" && json[0].Path != "/status" {
				c.JSON(422, api.NewError("PATCH support is limited and can only replace the /status path"))
				return
			}
			todo.Status = json[0].Value

			if err = tr.execTodoUpdate(todo); err != nil {
				tr.err(c, err, 404, "problem updating todo, doesn't exist?")
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
		tr.err(c, err, 422, "problem decoding id sent")
		return
	}

	if _, err = tr.db.Exec(`DELETE FROM Todo WHERE id = ?`, id); err != nil {
		tr.err(c, err, 500, "database error")
		return
	}

	c.Data(204, "application/json", make([]byte, 0))
}

func (tr *TodoResource) err(c *gin.Context, err error, status int, msg string) {
	log.Print(err)
	c.JSON(status, api.NewError(msg))
	return
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

func (tr *TodoResource) execTodoUpdate(todo api.Todo) error {
	_, err := tr.db.Exec(
		`UPDATE Todo SET status = ?, title = ?, description = ? WHERE id = ?`,
		todo.Status, todo.Title, todo.Description, todo.Id,
	)
	return err
}
