package ws

import (
	"context"
	"log"
	"net/http"
	"time"
	"todo-app/act"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

var (
	conns = make(map[*websocket.Conn]struct{}, 0)
	prod  Producer
)

func Init(p Producer, c Consumer) act.Notifier {
	prod = p
	http.HandleFunc("GET /ws/todos", getToDoActions)
	go consume(c)
	return notify
}

func notify(todoAction *act.TodoAction) {
	err := prod.Produce(todoAction)
	if err != nil {
		log.Printf("Error producing todo: %v\n", err)
	}
}

func consume(cons Consumer) {
	for {
		if len(conns) == 0 {
			time.Sleep(time.Second)
			continue
		}
		todo, err := cons.Consume()
		if err != nil {
			log.Printf("Error consuming todo: %v\n", err)
		}
		if todo == nil {
			continue
		}
		for conn := range conns {
			log.Printf("Sending todo event")
			err = wsjson.Write(context.Background(), conn, todo)
			if err != nil {
				delete(conns, conn)
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
