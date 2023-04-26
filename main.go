package main

import (
	entity "chat/entities"
	handler "chat/handlers"
	"log"
)

func init() {

	err := entity.MySQLTablesInit()
	if err != nil {
		log.Fatal(err)
	}
	_, err = entity.RedisConnect()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	handler.StartServer()
}
