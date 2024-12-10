package main

import (
	"todo-app/db"
	"todo-app/web"
)

func main() {

	db.Connect()
	defer db.Close()

	web.Serve()
}
