package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/benschw/go-todo/api"
	"io/ioutil"
	"net/http"
	"strconv"
)

type TodoClient struct {
	Host string
}

func (tc *TodoClient) CreateTodo(title string, description string) (api.Todo, error) {
	var respTodo api.Todo
	todo := api.Todo{Title: title, Description: description}

	b, err := json.Marshal(todo)
	if err != nil {
		return respTodo, err
	}

	body := bytes.NewBuffer(b)
	r, err := http.Post("http://"+tc.Host+"/todo", "application/json", body)
	if err != nil {
		return respTodo, err
	}
	if r.StatusCode != 201 {
		return respTodo, errors.New("response status of " + r.Status)
	}
	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return respTodo, err
	}

	if err = json.Unmarshal(respBody, &respTodo); err != nil {
		return respTodo, err
	}

	return respTodo, nil
}

func (tc *TodoClient) GetAllTodos() ([]api.Todo, error) {
	var respTodos []api.Todo

	r, err := http.Get("http://" + tc.Host + "/todo")
	if err != nil {
		return respTodos, err
	}
	if r.StatusCode != 200 {
		return respTodos, errors.New("response status of " + r.Status)
	}
	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return respTodos, err
	}

	if err = json.Unmarshal(respBody, &respTodos); err != nil {
		return respTodos, err
	}

	return respTodos, nil
}

func (tc *TodoClient) GetTodo(id int32) (api.Todo, error) {
	var respTodo api.Todo

	r, err := http.Get("http://" + tc.Host + "/todo/" + strconv.FormatInt(int64(id), 10))
	if err != nil {
		return respTodo, err
	}
	if r.StatusCode != 200 {
		return respTodo, errors.New("response status of " + r.Status)
	}
	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return respTodo, err
	}

	if err = json.Unmarshal(respBody, &respTodo); err != nil {
		return respTodo, err
	}

	return respTodo, nil
}

func (tc *TodoClient) DeleteTodo(id int32) error {
	url := "http://" + tc.Host + "/todo/" + strconv.FormatInt(int64(id), 10)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if r.StatusCode != 204 {
		return errors.New("response status of " + r.Status)
	}
	return nil
}
