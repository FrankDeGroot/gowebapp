package web

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"todo-app/db"
	"todo-app/dto"
	"todo-app/events/consumer"
	"todo-app/events/producer"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

const CONTENT_TYPE_JSON = "application/json;charset=utf-8"
const TODO_PATH = "/api/todos"

var wsConn = make([]*websocket.Conn, 0)

func Serve() {
	log.Println("Starting web server")
	registerHandlers()
	go func() {
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
	}()
	http.ListenAndServe(":8000", nil)
}

func registerHandlers() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("GET "+TODO_PATH, getAll)
	http.HandleFunc("GET "+TODO_PATH+"/{id}", getOne)
	http.HandleFunc("POST "+TODO_PATH, post)
	http.HandleFunc("PUT "+TODO_PATH, put)
	http.HandleFunc("DELETE "+TODO_PATH+"/{id}", delete)
	http.HandleFunc("GET /ws/todos", getEvents)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	toDos, err := db.GetAll()
	switch err {
	case nil:
		encode(w, toDos)
	default:
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error getting todos: %v\n", err)
	}
}

func getOne(w http.ResponseWriter, r *http.Request) {
	toDo, err := db.GetOne(r.PathValue("id"))
	switch err {
	case nil:
		encode(w, toDo)
	case db.ErrNotFound:
		http.Error(w, "id not found", http.StatusNotFound)
	default:
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error getting todo: %v\n", err)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	var toDo = dto.ToDo{}
	if err := json.NewDecoder(r.Body).Decode(&toDo); err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
	}
	savedToDo, err := db.Insert(&toDo)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error posting todo: %v\n", err)
	}
	err = producer.Produce(savedToDo)
	if err != nil {
		log.Printf("Error producing todo: %v\n", err)
	}
	encode(w, savedToDo)
}

func put(w http.ResponseWriter, r *http.Request) {
	var toDo = dto.SavedToDo{}
	if err := json.NewDecoder(r.Body).Decode(&toDo); err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
	}
	if err := db.Update(&toDo); err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error putting todo: %v\n", err)
	}
	if err := producer.Produce(&toDo); err != nil {
		log.Printf("Error producing todo: %v\n", err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := db.Delete(id); err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error deleting todo: %v\n", err)
	}
	if err := producer.Produce(&dto.SavedToDo{Id: id}); err != nil {
		log.Printf("Error producing todo: %v\n", err)
	}
	w.WriteHeader(http.StatusNoContent)
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

func encode(w http.ResponseWriter, a any) error {
	w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
	return json.NewEncoder(w).Encode(a)
}
