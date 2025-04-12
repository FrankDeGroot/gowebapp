package main

import (
	"log"
	"todo-app/db"
	"todo-app/kafka/consumer"
	"todo-app/kafka/producer"
	"todo-app/web"
	"todo-app/ws"
)

const topic = "todo"

func main() {

	repo, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	p, err := producer.Connect(topic)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	c, err := consumer.Connect(topic, topic)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	web.Serve(ws.Init(p, c), repo)
}
