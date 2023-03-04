package main

import (
	entity "chat/entities"
	handler "chat/handlers"
)

func main() {
	err := entity.DbTablesInit()
	if err != nil {
		panic(err)
	}
	handler.StartServer()
}
