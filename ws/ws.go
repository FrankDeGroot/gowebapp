package ws

import (
	"context"
	"log"
	"net/http"
	"todo-app/act"
	"todo-app/db"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const WS_PATH = "/ws"

var (
	prod     Producer
	connChan = make(chan *websocket.Conn)
)

func init() {
	http.HandleFunc("GET "+WS_PATH, connect)
}

func Open(p Producer, c Consumer, r db.TaskDber) act.Notifier {
	prod = p
	go broadcast(c, r)
	return notify
}

func notify(taskAction *act.TaskAction) {
	err := prod.Produce(taskAction)
	if err != nil {
		log.Printf("Error producing task: %v\n", err)
	}
}

func broadcast(cons Consumer, repo db.TaskDber) {
	conns := make(map[*websocket.Conn]struct{}, 0)
	consContChan := make(chan bool)
	defer close(consContChan)
	consChan := make(chan *act.TaskAction)
	defer (func() { consContChan <- false })()
	for {
		select {
		case conn, ok := <-connChan:
			if !ok {
				return
			}
			conns[conn] = struct{}{}
			readContChan := make(chan struct{})
			defer close(readContChan)
			go read(conn, readContChan, repo)
			if len(conns) == 1 {
				go consume(cons, consChan, consContChan)
			}
		case task, ok := <-consChan:
			if !ok {
				return
			}
			log.Printf("Broadcasting %v to %v conns", task, len(conns))
			for conn := range conns {
				err := wsjson.Write(context.Background(), conn, task)
				if err != nil {
					conn.CloseNow()
					delete(conns, conn)
					log.Printf("Error writing to socket broadcasting %v", err)
					if len(conns) == 0 {
						consContChan <- false
					}
				}
			}
		}
	}
}

func connect(w http.ResponseWriter, r *http.Request) {
	log.Printf("Connect")
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error accepting websocket: %v", err)
	}
	connChan <- conn
}
