package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"todo-app/db"
	"todo-app/dto"
)

func Serve() {
	fmt.Println("Starting web server")
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("GET /api/todos", getAll)
	http.HandleFunc("GET /api/todos/{id}", getOne)
	http.HandleFunc("POST /api/todos", post)
	http.HandleFunc("PUT /api/todos", put)
	http.HandleFunc("DELETE /api/todos/{id}", delete)
	http.ListenAndServe(":8000", nil)
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
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	switch err {
	case nil:
		break
	case strconv.ErrSyntax:
	case strconv.ErrRange:
		http.Error(w, "id not an integer", http.StatusBadRequest)
	default:
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error parsing id: %v\n", err)
	}

	toDo, err := db.Get(int(id))
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
	encode(w, toDo)
}

func delete(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Error", http.StatusInternalServerError)
}

func encode(w http.ResponseWriter, a any) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(a)
}
