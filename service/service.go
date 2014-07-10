package service

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Config struct {
	SvcHost    string
	DbUser     string
	DbPassword string
	DbHost     string
	DbName     string
}

type TodoService struct {
}

func (s *TodoService) Run(cfg Config) error {
	connectionString := cfg.DbUser + ":" + cfg.DbPassword + "@tcp(" + cfg.DbHost + ":3306)/" + cfg.DbName + "?charset=utf8&parseTime=True"

	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	db.SingularTable(true)

	todoResource := &TodoResource{db: db}

	r := gin.Default()

	r.GET("/todo", todoResource.GetAllTodos)
	r.GET("/todo/:id", todoResource.GetTodo)
	r.POST("/todo", todoResource.CreateTodo)
	r.PUT("/todo/:id", todoResource.UpdateTodo)
	r.PATCH("/todo/:id", todoResource.PatchTodo)
	r.DELETE("/todo/:id", todoResource.DeleteTodo)

	r.Run(cfg.SvcHost)

	return nil
}
