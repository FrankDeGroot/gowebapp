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
	db.Connect()
	defer db.Close()

	registerHandlers()
	srv := httptest.NewServer(nil)
	defer srv.Close()

	toDo := dto.ToDo{Description: "test" + time.Now().Format(time.RFC3339), Done: false}
	toDoBytes, err := json.Marshal(toDo)
	if err != nil {
		t.Fatal(err)
	}
	post, err := http.Post(srv.URL+TODO_PATH, CONTENT_TYPE_JSON, bytes.NewBuffer(toDoBytes))
	if err != nil {
		t.Fatal(err)
	}
	if post.StatusCode != http.StatusOK {
		t.Fatal(err)
	}
	postToDo := dto.SavedToDo{}
	json.NewDecoder(post.Body).Decode(&postToDo)
	if postToDo.Description != toDo.Description {
		t.Fatalf("Wanted %v got %v", toDo.Description, postToDo.Description)
	}
	if postToDo.Done != toDo.Done {
		t.Fatalf("Wanted %v got %v", toDo.Done, postToDo.Done)
	}

	getOne, err := http.Get(srv.URL + TODO_PATH + "/" + postToDo.Id)
	if err != nil {
		t.Fatal(err)
	}
	getToDo := dto.SavedToDo{}
	json.NewDecoder(getOne.Body).Decode(&getToDo)
	if getToDo.Description != toDo.Description {
		t.Fatalf("Wanted %v got %v", toDo.Description, getToDo.Description)
	}
	if getToDo.Done != toDo.Done {
		t.Fatalf("Wanted %v got %v", toDo.Done, getToDo.Done)
	}

	getAll, err := http.Get(srv.URL + TODO_PATH)
	if err != nil {
		t.Fatal(err)
	}
	getToDos := []dto.SavedToDo{}
	json.NewDecoder(getAll.Body).Decode(&getToDos)
	found := false
	for _, getToDo := range getToDos {
		if getToDo.Id == postToDo.Id {
			found = true
			if getToDo.Description != toDo.Description {
				t.Fatalf("Wanted %v got %v", toDo.Description, getToDo.Description)
			}
			if getToDo.Done != toDo.Done {
				t.Fatalf("Wanted %v got %v", toDo.Done, getToDo.Done)
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
