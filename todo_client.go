package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type TodoClient struct {
	host string
}

func (tc *TodoClient) CreateTodo(title string, description string) (Todo, error) {
	var respTodo Todo
	todo := Todo{Title: title, Description: description}

	b, err := json.Marshal(todo)
	if err != nil {
		return respTodo, err
	}

	body := bytes.NewBuffer(b)
	r, err := http.Post("http://"+tc.host+"/todo", "application/json", body)
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

func (tc *TodoClient) GetAllTodos() ([]Todo, error) {
	var respTodos []Todo

	r, err := http.Get("http://" + tc.host + "/todo")
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

func (tc *TodoClient) GetTodo(id int32) (Todo, error) {
	var respTodo Todo

	r, err := http.Get("http://" + tc.host + "/todo/" + strconv.FormatInt(int64(id), 10))
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
	url := "http://" + tc.host + "/todo/" + strconv.FormatInt(int64(id), 10)
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
