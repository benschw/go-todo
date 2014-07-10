package main

import (
	"fmt"
	"github.com/benschw/go-todo/client"
	"log"
	"testing"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.

func TestCreateTodo(t *testing.T) {

	// given
	client := client.TodoClient{Host: "http://localhost:8080"}

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
