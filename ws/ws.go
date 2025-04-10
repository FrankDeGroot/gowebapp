package ws

import (
	"context"
	"log"
	"net/http"
	"time"
	"todo-app/dto"
	"todo-app/events/consumer"
	"todo-app/events/producer"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

var wsConn = make([]*websocket.Conn, 0)

type Notify struct{}

func Serve() {
	http.HandleFunc("GET /ws/todos", getEvents)
	go consume()
}

func consume() {
	for {
		if len(wsConn) == 0 {
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
		newWsConn := make([]*websocket.Conn, 0)
		for _, c := range wsConn {
			log.Printf("Sending todo event")
			err = wsjson.Write(context.Background(), c, todo)
			if err != nil {
				log.Printf("Error writing to socket: %v", err)
			} else {
				newWsConn = append(newWsConn, c)
			}
		}
		wsConn = newWsConn
		log.Printf("Sent events successfully to %v connections", len(wsConn))
	}
}

func getEvents(w http.ResponseWriter, r *http.Request) {
	log.Printf("Getting events")
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error accepting websocket: %v\n", err)
	}
	wsConn = append(wsConn, c)
	log.Printf("Websocket connected %v connections", len(wsConn))
}

func (n Notify) Add(savedToDo *dto.SavedToDo) {
	err := producer.Produce(&dto.ToDoEvent{
		Action:    dto.ActionAdd,
		SavedToDo: *savedToDo,
	})
	if err != nil {
		log.Printf("Error producing todo: %v\n", err)
	}
}

func (n Notify) Change(savedToDo *dto.SavedToDo) {
	if err := producer.Produce(&dto.ToDoEvent{
		Action:    dto.ActionChg,
		SavedToDo: *savedToDo,
	}); err != nil {
		log.Printf("Error producing todo: %v\n", err)
	}
}

func (n Notify) Delete(savedToDo *dto.SavedToDo) {
	if err := producer.Produce(&dto.ToDoEvent{
		Action:    dto.ActionDel,
		SavedToDo: *savedToDo,
	}); err != nil {
		log.Printf("Error producing todo: %v\n", err)
	}
}
