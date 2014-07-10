package service

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type TodoService struct {
	SvcHost    string
	DbUser     string
	DbPassword string
	DbHost     string
	DbName     string
}

func (s *TodoService) Run() error {
	connectionString := s.DbUser + ":" + s.DbPassword + "@tcp(" + s.DbHost + ":3306)/" + s.DbName + "?charset=utf8&parseTime=True"

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

	r.Run(s.SvcHost)

	return nil
}
