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
const TASKS_PATH = "/api/tasks"

var (
	ntfy = func(*act.TaskAction) {}
	repo db.TaskDber
)

func Serve(n act.Notifier, r db.TaskDber) {
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
	http.HandleFunc("GET "+TASKS_PATH, getAll)
	http.HandleFunc("GET "+TASKS_PATH+"/{id}", getOne)
	http.HandleFunc("POST "+TASKS_PATH, post)
	http.HandleFunc("PUT "+TASKS_PATH, put)
	http.HandleFunc("DELETE "+TASKS_PATH+"/{id}", delete)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := repo.GetAll()
	switch err {
	case nil:
		encode(w, tasks)
	default:
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error getting tasks: %v\n", err)
	}
}

func getOne(w http.ResponseWriter, r *http.Request) {
	task, err := repo.GetOne(r.PathValue("id"))
	switch err {
	case nil:
		encode(w, task)
	case db.ErrNotFound:
		http.Error(w, "id not found", http.StatusNotFound)
	default:
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error getting task: %v\n", err)
	}
}

func post(w http.ResponseWriter, r *http.Request) {
	var task = dto.Task{}
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
	}
	savedTask, err := repo.Insert(&task)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error posting task: %v\n", err)
	}
	ntfy(act.Make(act.Post, savedTask))
	encode(w, savedTask)
}

func put(w http.ResponseWriter, r *http.Request) {
	var savedTask = dto.SavedTask{}
	if err := json.NewDecoder(r.Body).Decode(&savedTask); err != nil {
		http.Error(w, "Error", http.StatusBadRequest)
	}
	if err := repo.Update(&savedTask); err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error putting task: %v\n", err)
	}
	ntfy(act.Make(act.Put, &savedTask))
	w.WriteHeader(http.StatusNoContent)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := repo.Delete(id); err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		log.Printf("Error deleting task: %v\n", err)
	}
	ntfy(act.Make(act.Delete, &dto.SavedTask{Id: id}))
	w.WriteHeader(http.StatusNoContent)
}

func encode(w http.ResponseWriter, a any) error {
	w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
	return json.NewEncoder(w).Encode(a)
}
