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

var (
	ntfy = func(*act.TodoAction) {}
	repo TodoRepo
)

func Serve(n act.Notifier, r TodoRepo) {
	ntfy = n
	repo = r
	log.Println("Starting web server")
	setHandlers()
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatalf("Error serving http %v", err)
	}
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
	todos, err := repo.GetAll()
	switch err {
	case nil:
		encode(w, todos)
	default:
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error getting todos: %v\n", err)
	}
}

func getOne(w http.ResponseWriter, r *http.Request) {
	todo, err := repo.GetOne(r.PathValue("id"))
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
	savedTodo, err := repo.Insert(&todo)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error posting todo: %v\n", err)
	}
	ntfy(act.Make(act.Post, savedTodo))
	encode(w, savedTodo)
}

func put(w http.ResponseWriter, r *http.Request) {
	var savedTodo = dto.SavedTodo{}
	if err := json.NewDecoder(r.Body).Decode(&savedTodo); err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
	}
	if err := repo.Update(&savedTodo); err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error putting todo: %v\n", err)
	}
	ntfy(act.Make(act.Put, &savedTodo))
	w.WriteHeader(http.StatusNoContent)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := repo.Delete(id); err != nil {
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
