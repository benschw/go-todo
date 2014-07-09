package main

import (
	"github.com/gin-gonic/gin"
)

type TodoService struct {
	SvcHost    string
	DbUser     string
	DbPassword string
	DbHost     string
	DbName     string
}

func (s *TodoService) Run() error {
	connectionString := s.DbUser + ":" + s.DbPassword + "@tcp(" + s.DbHost + ":3306)/" + s.DbName

	db, err := Open(connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	todoResource := &TodoResource{db: db}

	r := gin.Default()

	r.GET("/todo", todoResource.GetAllTodos)
	r.GET("/todo/:id", todoResource.GetTodo)
	r.POST("/todo", todoResource.CreateTodo)
	r.DELETE("/todo/:id", todoResource.DeleteTodo)

	r.Run(s.SvcHost)

	return nil
}
