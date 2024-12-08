package main

import (
	"os"
	"todo-app/db"
	"todo-app/web"
)

func main() {

	db.Connect()
	defer db.Close()

	var args = os.Args
	if len(args) > 1 && args[1] == "init" {
		db.Init()
	} else {
		web.Serve()
	}
}
