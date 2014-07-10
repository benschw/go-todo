package main

import (
	"fmt"
	"github.com/benschw/go-todo/client"
	"log"
	"testing"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.

func TestDeleteTodo(t *testing.T) {

	// given
	client := client.TodoClient{Host: "http://localhost:8080"}
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

func TestDeleteNotFoundTodo(t *testing.T) {

	// given
	client := client.TodoClient{Host: "http://localhost:8080"}
	id := int32(3)
	// when
	err := client.DeleteTodo(id)

	// then
	if err == nil {
		t.Error(err)
	}

}
