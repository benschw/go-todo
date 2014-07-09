package main

import (
	"fmt"
	"log"
	"testing"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.

func TestCreateTodo(t *testing.T) {

	// given
	client := TodoClient{host: "localhost:8080"}

	// when
	todo, err := client.CreateTodo("foo", "bar")

	//then
	if err != nil {
		t.Error(err)
	}

	if todo.Title != "foo" && todo.Description != "bar" {
		t.Error("returned todo not right")
	}

	// cleanup
	_ = client.DeleteTodo(todo.Id)
}

func TestGetTodo(t *testing.T) {

	// given
	client := TodoClient{host: "localhost:8080"}
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

func TestGetAllTodos(t *testing.T) {

	// given
	client := TodoClient{host: "localhost:8080"}
	client.CreateTodo("foo", "bar")
	client.CreateTodo("baz", "bing")

	// when
	todos, err := client.GetAllTodos()

	// then
	if err != nil {
		t.Error(err)
	}

	if len(todos) != 2 {
		t.Error("wrong number of todos")
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

func TestDeleteTodo(t *testing.T) {

	// given
	client := TodoClient{host: "localhost:8080"}
	todo, _ := client.CreateTodo("foo", "bar")
	id := todo.Id

	// when
	err := client.DeleteTodo(id)

	// then
	if err != nil {
		t.Error(err)
	}

	_, err = client.GetTodo(id)
	if err == nil {
		t.Error(err)
	}
}
