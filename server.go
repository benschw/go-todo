package main

import (
	"errors"
	"github.com/benschw/go-todo/service"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"log"
	"os"
)

func getConfig(yamlPath string) (service.Config, error) {
	config := service.Config{}

	if _, err := os.Stat(yamlPath); err != nil {
		return config, errors.New("config path not valid")
	}

	ymlData, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal([]byte(ymlData), &config)
	return config, err
}

func main() {
	yamlPath := "config.yaml"

	cfg, err := getConfig(yamlPath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	svc := service.TodoService{}

	if err = svc.Run(cfg); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}
