package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"todo-app/db"
	"todo-app/dto"
)

func TestGetOne(t *testing.T) {
	// TODO Use test database server like db_test
	db.Connect()
	defer db.Close()

	setHandlers()
	srv := httptest.NewServer(nil)
	defer srv.Close()

	todo := dto.Todo{Description: "test" + time.Now().Format(time.RFC3339), Done: false}
	todoBytes, err := json.Marshal(todo)
	if err != nil {
		t.Fatal(err)
	}
	post, err := http.Post(srv.URL+TODO_PATH, CONTENT_TYPE_JSON, bytes.NewBuffer(todoBytes))
	if err != nil {
		t.Fatal(err)
	}
	if post.StatusCode != http.StatusOK {
		t.Fatal(err)
	}
	postToDo := dto.SavedTodo{}
	json.NewDecoder(post.Body).Decode(&postToDo)
	if postToDo.Description != todo.Description {
		t.Fatalf("Wanted %v got %v", todo.Description, postToDo.Description)
	}
	if postToDo.Done != todo.Done {
		t.Fatalf("Wanted %v got %v", todo.Done, postToDo.Done)
	}

	getOne, err := http.Get(srv.URL + TODO_PATH + "/" + postToDo.Id)
	if err != nil {
		t.Fatal(err)
	}
	getToDo := dto.SavedTodo{}
	json.NewDecoder(getOne.Body).Decode(&getToDo)
	if getToDo.Description != todo.Description {
		t.Fatalf("Wanted %v got %v", todo.Description, getToDo.Description)
	}
	if getToDo.Done != todo.Done {
		t.Fatalf("Wanted %v got %v", todo.Done, getToDo.Done)
	}

	getAll, err := http.Get(srv.URL + TODO_PATH)
	if err != nil {
		t.Fatal(err)
	}
	getToDos := []dto.SavedTodo{}
	json.NewDecoder(getAll.Body).Decode(&getToDos)
	found := false
	for _, getToDo := range getToDos {
		if getToDo.Id == postToDo.Id {
			found = true
			if getToDo.Description != todo.Description {
				t.Fatalf("Wanted %v got %v", todo.Description, getToDo.Description)
			}
			if getToDo.Done != todo.Done {
				t.Fatalf("Wanted %v got %v", todo.Done, getToDo.Done)
			}
			break
		}
	}
	if !found {
		t.Fatal("Did not find the todo")
	}

	delReq, err := http.NewRequest(http.MethodDelete, srv.URL+TODO_PATH+"/"+postToDo.Id, nil)
	if err != nil {
		t.Fatal(err)
	}
	delRes, err := http.DefaultClient.Do(delReq)
	if err != nil {
		t.Fatal(err)
	}
	if delRes.StatusCode != http.StatusNoContent {
		t.Fatal(err)
	}
}
