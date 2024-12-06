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
	http.HandleFunc("GET /api/todos/{id}", get)
	http.HandleFunc("PUT /api/todos", put)
	http.HandleFunc("POST /api/todos", put)
	http.ListenAndServe(":8000", nil)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	toDos, err := db.GetAll()
	switch err {
	case nil:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(toDos)
	default:
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error getting todos: %v\n", err)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(toDo)
	case db.ErrNotFound:
		http.Error(w, "id not found", http.StatusNotFound)
	default:
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error getting todo: %v\n", err)
	}
}

func put(w http.ResponseWriter, r *http.Request) {
	var toDo = dto.ToDo{}
	if err := json.NewDecoder(r.Body).Decode(&toDo); err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
	}
	if err := db.Upsert(&toDo); err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error putting todo: %v\n", err)
	}
}
