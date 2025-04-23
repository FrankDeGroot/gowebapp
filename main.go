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

	repo, err := db.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	p, err := producer.Open(topic)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	c, err := consumer.Open(topic, topic)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	n := ws.Open(p, c)
	defer ws.Close()
	web.Serve(n, repo)
}
