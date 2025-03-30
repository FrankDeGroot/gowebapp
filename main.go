package main

import (
	"todo-app/db"
	"todo-app/events/consumer"
	"todo-app/events/producer"
	"todo-app/web"
)

const topic = "todo"

func main() {

	db.Connect()
	defer db.Close()

	producer.Connect(topic)
	defer producer.Close()

	consumer.Connect(topic, topic)
	defer consumer.Close()

	web.Serve()
}
