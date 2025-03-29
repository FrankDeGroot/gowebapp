package main

import (
	"todo-app/db"
	"todo-app/events/producer"
	"todo-app/web"
)

func main() {

	db.Connect()
	defer db.Close()

	producer.Connect("todo")
	defer producer.Close()

	web.Serve()
}
