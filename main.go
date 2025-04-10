package main

import (
	"todo-app/db"
	"todo-app/events/consumer"
	"todo-app/events/producer"
	"todo-app/web"
	"todo-app/ws"
)

const topic = "todo"

func main() {

	db.Connect()
	defer db.Close()

	producer.Connect(topic)
	defer producer.Close()

	consumer.Connect(topic, topic)
	defer consumer.Close()

	ws.Serve()
	web.Serve()
}
