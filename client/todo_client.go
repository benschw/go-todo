package client

import (
	"github.com/benschw/go-todo/api"
	"log"
	"strconv"
)

var _ = log.Print

type TodoClient struct {
	Host string
}

func (tc *TodoClient) CreateTodo(title string, description string) (api.Todo, error) {
	var respTodo api.Todo
	todo := api.Todo{Title: title, Description: description}

	url := tc.Host + "/todo"
	r, err := makeRequest("POST", url, todo)
	if err != nil {
		return respTodo, err
	}
	err = processResponseEntity(r, &respTodo, 201)
	return respTodo, err
}

func (tc *TodoClient) GetAllTodos() ([]api.Todo, error) {
	var respTodos []api.Todo

	url := tc.Host + "/todo"
	r, err := makeRequest("GET", url, nil)
	if err != nil {
		return respTodos, err
	}
	err = processResponseEntity(r, &respTodos, 200)
	return respTodos, err
}

func (tc *TodoClient) GetTodo(id int32) (api.Todo, error) {
	var respTodo api.Todo

	url := tc.Host + "/todo/" + strconv.FormatInt(int64(id), 10)
	r, err := makeRequest("GET", url, nil)
	if err != nil {
		return respTodo, err
	}
	err = processResponseEntity(r, &respTodo, 200)
	return respTodo, err
}

func (tc *TodoClient) UpdateTodo(todo api.Todo) (api.Todo, error) {
	var respTodo api.Todo

	url := tc.Host + "/todo/" + strconv.FormatInt(int64(todo.Id), 10)
	r, err := makeRequest("PUT", url, todo)
	if err != nil {
		return respTodo, err
	}
	err = processResponseEntity(r, &respTodo, 200)
	return respTodo, err
}

func (tc *TodoClient) UpdateTodoStatus(id int32, status string) (api.Todo, error) {
	var respTodo api.Todo

	patchArr := make([]api.Patch, 1)
	patchArr[0] = api.Patch{Op: "replace", Path: "/status", Value: string(status)}

	url := tc.Host + "/todo/" + strconv.FormatInt(int64(id), 10)
	r, err := makeRequest("PATCH", url, patchArr)
	if err != nil {
		return respTodo, err
	}
	err = processResponseEntity(r, &respTodo, 200)
	return respTodo, err
}

func (tc *TodoClient) DeleteTodo(id int32) error {
	url := tc.Host + "/todo/" + strconv.FormatInt(int64(id), 10)
	r, err := makeRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	return processResponse(r, 204)
}
