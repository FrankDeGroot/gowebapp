package ws

import (
	"context"
	"log"
	"net/http"
	"todo-app/act"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const WS_PATH = "/ws"

var (
	prod     Producer
	connChan = make(chan *websocket.Conn)
)

func Open(p Producer, c Consumer) act.Notifier {
	prod = p
	http.HandleFunc("GET "+WS_PATH, wsConnect)
	go consume(c)
	return notify
}

func Close() {
	close(connChan)
}

func notify(taskAction *act.TaskAction) {
	err := prod.Produce(taskAction)
	if err != nil {
		log.Printf("Error producing task: %v\n", err)
	}
}

func consume(cons Consumer) {
	conns := make(map[*websocket.Conn]struct{}, 0)
	for conn := range connChan {
		conns[conn] = struct{}{}

		for len(conns) != 0 {
			task, err := cons.Consume()
			if err != nil {
				log.Printf("Error consuming task: %v\n", err)
			}
			if task == nil {
				continue
			}

		addConns:
			for {
				select {
				case conn, ok := <-connChan:
					if !ok {
						return
					}
					conns[conn] = struct{}{}
				default:
					break addConns
				}
			}
			for conn := range conns {
				log.Printf("Sending task event")
				err = wsjson.Write(context.Background(), conn, task)
				if err != nil {
					delete(conns, conn)
					log.Printf("Error writing to socket: %v", err)
				}
			}
			log.Printf("Sent events successfully to %v connections", len(conns))
		}
	}
}

func wsConnect(w http.ResponseWriter, r *http.Request) {
	log.Printf("Getting events")
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error accepting websocket: %v\n", err)
	}
	connChan <- c
}
