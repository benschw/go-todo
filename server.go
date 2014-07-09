package main

import (
	"github.com/benschw/go-todo/service"
	"log"
	"os"
)

func main() {
	var (
		svcHost    = ":8080" // 0.0.0.0:8080
		dbUser     = "root"
		dbPassword = ""
		dbHost     = "localhost"
		dbName     = "Todo"
	)
	svc := service.TodoService{SvcHost: svcHost, DbUser: dbUser, DbPassword: dbPassword, DbHost: dbHost, DbName: dbName}

	err := svc.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}
