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
	http.HandleFunc("GET "+WS_PATH, connect)
	go broadcast(c)
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

func broadcast(cons Consumer) {
	conns := make(map[*websocket.Conn]struct{}, 0)
	consChan := make(chan *act.TaskAction)
	defer close(consChan)
	consContChan := make(chan bool)
	defer close(consContChan)
	for {
		select {
		case conn, ok := <-connChan:
			if !ok {
				return
			}
			conns[conn] = struct{}{}
			readContChan := make(chan struct{})
			defer close(readContChan)
			go read(conn, readContChan)
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
					log.Printf("Error writing to socket: %v", err)
					if len(conns) == 0 {
						consContChan <- false
					}
				}
			}
		}
	}
}

func consume(cons Consumer, consChan chan *act.TaskAction, contChan chan bool) {
	cont, ok := true, true
	for {
		select {
		case cont, ok = <-contChan:
			if !cont || !ok {
				return
			}
		default:
			task, err := cons.Consume()
			if err != nil {
				log.Printf("Error consuming task: %v\n", err)
			}
			if task == nil {
				continue
			}
			log.Printf("Consumed %v", task)
			consChan <- task
		}
	}
}

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

func connect(w http.ResponseWriter, r *http.Request) {
	log.Printf("Connect")
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error accepting websocket: %v", err)
	}
	connChan <- conn
}
