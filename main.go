package main

import (
	"todo-app/db"
	"todo-app/web"
)

func main() {

	db.Init()
	defer db.Close()

	web.Serve()
}
