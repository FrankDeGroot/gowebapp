package ws

import (
	"context"
	"log"
	"todo-app/act"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func read(conn *websocket.Conn, cont chan struct{}) {
	for {
		select {
		case _, ok := <-cont:
			if !ok {
				return
			}
		default:
			action := act.TaskAction{}
			err := wsjson.Read(context.Background(), conn, &action)
			if err != nil {
				log.Printf("Error reading websocket %v", err)
				conn.CloseNow()
				return
			}
			log.Printf("Read %v", action)
		}
	}
}
