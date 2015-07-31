package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mingderwang/go-todo/api"
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

func (s *TodoService) getDb(cfg Config) (gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "/tmp/Todo.db")
	//db.CreateTable(new(api.Todo))
	//db.LogMode(true)
	return db, err
}

func (s *TodoService) Migrate(cfg Config) error {
	db, err := s.getDb(cfg)
	if err != nil {
		return err
	}
	db.SingularTable(true)

	db.AutoMigrate(&api.Todo{})
	return nil
}
func (s *TodoService) Run(cfg Config) error {
	db, err := s.getDb(cfg)
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
