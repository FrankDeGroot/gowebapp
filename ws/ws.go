package ws

import (
	"context"
	"log"
	"net/http"
	"time"
	"todo-app/act"
	"todo-app/kafka/consumer"
	"todo-app/kafka/producer"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

var conns = make(map[*websocket.Conn]struct{}, 0)

func Init() act.Notifier {
	http.HandleFunc("GET /ws/todos", getToDoActions)
	go consume()
	return Notify
}

func Notify(todoAction *act.TodoAction) {
	err := producer.Produce(todoAction)
	if err != nil {
		log.Printf("Error producing todo: %v\n", err)
	}
	if err != nil {
		log.Printf("Error producing todo: %v\n", err)
	}
}

func consume() {
	for {
		if len(conns) == 0 {
			time.Sleep(time.Second)
			continue
		}
		todo, err := consumer.Consume()
		if err != nil {
			log.Printf("Error consuming todo: %v\n", err)
		}
		if todo == nil {
			continue
		}
		for c := range conns {
			log.Printf("Sending todo event")
			err = wsjson.Write(context.Background(), c, todo)
			if err != nil {
				delete(conns, c)
				log.Printf("Error writing to socket: %v", err)
			}
		}
		log.Printf("Sent events successfully to %v connections", len(conns))
	}
}

func getToDoActions(w http.ResponseWriter, r *http.Request) {
	log.Printf("Getting events")
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error accepting websocket: %v\n", err)
	}
	conns[c] = struct{}{}
	log.Printf("Websocket connected %v connections", len(conns))
}
