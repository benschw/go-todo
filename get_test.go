package main

import (
	"fmt"
	"github.com/benschw/go-todo/client"
	"log"
	"testing"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.

func TestGetTodo(t *testing.T) {

	// given
	client := client.TodoClient{Host: "http://localhost:8080"}
	todo, _ := client.CreateTodo("foo", "bar")
	id := todo.Id

	// when
	todo, err := client.GetTodo(id)

	// then
	if err != nil {
		t.Error(err)
	}

	if todo.Title != "foo" && todo.Description != "bar" {
		t.Error("returned todo not right")
	}

	// cleanup
	_ = client.DeleteTodo(todo.Id)
}

func TestGetNotFoundTodo(t *testing.T) {

	// given
	client := client.TodoClient{Host: "http://localhost:8080"}
	id := int32(3)

	// when
	_, err := client.GetTodo(id)

	// then
	if err == nil {
		t.Error(err)
	}
}

func TestGetAllTodos(t *testing.T) {

	// given
	client := client.TodoClient{Host: "http://localhost:8080"}
	client.CreateTodo("foo", "bar")
	client.CreateTodo("baz", "bing")

	// when
	todos, err := client.GetAllTodos()

	// then
	if err != nil {
		t.Error(err)
	}

	if len(todos) != 2 {
		t.Errorf("wrong number of todos: %d", len(todos))
	}
	if todos[0].Title != "foo" && todos[0].Description != "bar" {
		t.Error("returned todo not right")
	}
	if todos[1].Title != "baz" && todos[1].Description != "bing" {
		t.Error("returned todo not right")
	}

	// cleanup
	_ = client.DeleteTodo(todos[0].Id)
	_ = client.DeleteTodo(todos[1].Id)
}
