package ws

import (
	"context"
	"log"
	"net/http"
	"todo-app/act"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const WS_TODO_PATH = "/ws/todos"

var (
	prod     Producer
	connChan = make(chan *websocket.Conn)
)

func Init(p Producer, c Consumer) act.Notifier {
	prod = p
	http.HandleFunc("GET "+WS_TODO_PATH, getToDoActions)
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
	conns := make(map[*websocket.Conn]struct{}, 0)
	for {
		conn := <-connChan
		conns[conn] = struct{}{}

		for len(conns) != 0 {
			todo, err := cons.Consume()
			if err != nil {
				log.Printf("Error consuming todo: %v\n", err)
			}
			if todo == nil {
				continue
			}

		AddConns:
			for {
				select {
				case conn := <-connChan:
					conns[conn] = struct{}{}
				default:
					break AddConns
				}
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
}

func getToDoActions(w http.ResponseWriter, r *http.Request) {
	log.Printf("Getting events")
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error accepting websocket: %v\n", err)
	}
	connChan <- c
}
