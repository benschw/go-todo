package main

import (
	"fmt"
	"github.com/benschw/go-todo/api"
	"github.com/benschw/go-todo/client"
	"log"
	"reflect"
	"testing"
)

var _ = fmt.Print // For debugging; delete when done.
var _ = log.Print // For debugging; delete when done.

func TestUpdateTodo(t *testing.T) {

	// given
	client := client.TodoClient{Host: "http://localhost:8080"}
	todo, _ := client.CreateTodo("foo", "bar")

	// when
	todo.Status = "doing"
	todo.Title = "baz"
	todo.Description = "bing"
	_, err := client.UpdateTodo(todo)

	// then
	if err != nil {
		t.Error(err)
	}

	todoResult, _ := client.GetTodo(todo.Id)

	if !reflect.DeepEqual(todo, todoResult) {
		t.Error("returned todo not updated")
	}

	// cleanup
	_ = client.DeleteTodo(todo.Id)
}

func TestUpdateNonExistantTodo(t *testing.T) {

	// given
	client := client.TodoClient{Host: "http://localhost:8080"}
	todo, _ := client.CreateTodo("foo", "bar")
	_ = client.DeleteTodo(todo.Id)

	// when
	todo.Status = "doing"
	todo.Title = "baz"
	todo.Description = "bing"
	_, err := client.UpdateTodo(todo)

	// then
	if err == nil {
		t.Error(err)
	}

}

func TestUpdateTodosStatus(t *testing.T) {

	// given
	client := client.TodoClient{Host: "http://localhost:8080"}
	todo, _ := client.CreateTodo("foo", "bar")

	// when
	_, err := client.UpdateTodoStatus(todo.Id, api.DoingStatus)

	// then
	if err != nil {
		t.Error(err)
	}

	todoResult, _ := client.GetTodo(todo.Id)

	if todoResult.Status != "doing" {
		t.Error("returned todo status not updated")
	}

	// cleanup
	_ = client.DeleteTodo(todo.Id)
}

func TestUpdateNotFoundTodosStatus(t *testing.T) {

	// given
	client := client.TodoClient{Host: "http://localhost:8080"}
	id := int32(3)
	// when
	_, err := client.UpdateTodoStatus(id, api.DoingStatus)

	// then
	if err == nil {
		t.Error(err)
	}
}
