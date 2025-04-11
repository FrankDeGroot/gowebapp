package web

import (
	"encoding/json"
	"log"
	"net/http"
	"todo-app/act"
	"todo-app/db"
	"todo-app/dto"
)

const CONTENT_TYPE_JSON = "application/json;charset=utf-8"
const TODO_PATH = "/api/todos"

var ntfy = func(*act.TodoAction) {}

func Serve(n act.Notifier) {
	ntfy = n
	log.Println("Starting web server")
	setHandlers()
	http.ListenAndServe(":8000", nil)
}

func setHandlers() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("GET "+TODO_PATH, getAll)
	http.HandleFunc("GET "+TODO_PATH+"/{id}", getOne)
	http.HandleFunc("POST "+TODO_PATH, post)
	http.HandleFunc("PUT "+TODO_PATH, put)
	http.HandleFunc("DELETE "+TODO_PATH+"/{id}", delete)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	todos, err := db.GetAll()
	switch err {
	case nil:
		encode(w, todos)
	default:
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error getting todos: %v\n", err)
	}
}

func getOne(w http.ResponseWriter, r *http.Request) {
	todo, err := db.GetOne(r.PathValue("id"))
	switch err {
	case nil:
		encode(w, todo)
	case db.ErrNotFound:
		http.Error(w, "id not found", http.StatusNotFound)
	default:
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error getting todo: %v\n", err)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	var todo = dto.Todo{}
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
	}
	savedTodo, err := db.Insert(&todo)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error posting todo: %v\n", err)
	}
	ntfy(act.Make(act.Add, savedTodo))
	encode(w, savedTodo)
}

func put(w http.ResponseWriter, r *http.Request) {
	var savedTodo = dto.SavedTodo{}
	if err := json.NewDecoder(r.Body).Decode(&savedTodo); err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
	}
	if err := db.Update(&savedTodo); err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error putting todo: %v\n", err)
	}
	ntfy(act.Make(act.Change, &savedTodo))
	w.WriteHeader(http.StatusNoContent)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := db.Delete(id); err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error deleting todo: %v\n", err)
	}
	ntfy(act.Make(act.Delete, &dto.SavedTodo{Id: id}))
	w.WriteHeader(http.StatusNoContent)
}

func encode(w http.ResponseWriter, a any) error {
	w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
	return json.NewEncoder(w).Encode(a)
}
