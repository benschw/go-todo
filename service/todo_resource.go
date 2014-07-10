package service

import (
	"github.com/benschw/go-todo/api"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"time"
)

type TodoResource struct {
	db gorm.DB
}

func (tr *TodoResource) CreateTodo(c *gin.Context) {
	var todo api.Todo

	if !c.Bind(&todo) {
		c.JSON(400, api.NewError("problem decoding body"))
		return
	}
	todo.Status = "todo"
	todo.Created = int32(time.Now().Unix())

	tr.db.Save(&todo)

	c.JSON(201, todo)
}

func (tr *TodoResource) GetAllTodos(c *gin.Context) {
	var todos []api.Todo

	tr.db.Order("created desc").Find(&todos)

	c.JSON(200, todos)
}

func (tr *TodoResource) GetTodo(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(400, api.NewError("problem decoding id sent"))
		return
	}

	var todo api.Todo

	tr.db.First(&todo, id)

	if todo.Id == 0 {
		c.JSON(404, api.NewError("not found"))
	} else {
		c.JSON(200, todo)
	}
}

func (tr *TodoResource) UpdateTodo(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(400, api.NewError("problem decoding id sent"))
		return
	}

	var todo api.Todo

	if !c.Bind(&todo) {
		c.JSON(400, api.NewError("problem decoding body"))
		return
	}
	todo.Id = int32(id)

	var existing api.Todo
	tr.db.First(&existing, id)

	if existing.Id == 0 {
		c.JSON(404, api.NewError("not found"))
	} else {
		tr.db.Save(&todo)
		c.JSON(200, todo)
	}
}

func (tr *TodoResource) PatchTodo(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(400, api.NewError("problem decoding id sent"))
		return
	}

	// this is a hack because Gin falsely claims my unmarshalled obj is invalid.
	// recovering from the panic and using my object that already has the json body bound to it.
	var json []api.Patch
	defer func() {
		if r := recover(); r != nil {
			if json[0].Op != "replace" && json[0].Path != "/status" {
				c.JSON(400, api.NewError("PATCH support is limited and can only replace the /status path"))
				return
			}

			var todo api.Todo

			tr.db.First(&todo, id)
			if todo.Id == 0 {
				c.JSON(404, api.NewError("not found"))
				return
			}
			todo.Status = json[0].Value

			tr.db.Save(&todo)
			c.JSON(200, todo)
		}
	}()
	c.Bind(&json)
}

func (tr *TodoResource) DeleteTodo(c *gin.Context) {
	id, err := tr.getId(c)
	if err != nil {
		c.JSON(400, api.NewError("problem decoding id sent"))
		return
	}

	var todo api.Todo

	tr.db.First(&todo, id)

	if todo.Id == 0 {
		c.JSON(404, api.NewError("not found"))
	} else {
		tr.db.Delete(&todo)
		c.Data(204, "application/json", make([]byte, 0))
	}
}

func (tr *TodoResource) getId(c *gin.Context) (int32, error) {
	idStr := c.Params.ByName("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	return int32(id), nil
}
