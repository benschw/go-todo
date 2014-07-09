package main

import (
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
	server := TodoService{SvcHost: svcHost, DbUser: dbUser, DbPassword: dbPassword, DbHost: dbHost, DbName: dbName}

	err := server.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}
