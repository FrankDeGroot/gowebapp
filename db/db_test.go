package db

import (
	"os"
	"testing"
	"time"
	"todo-app/dto"

	"github.com/cockroachdb/cockroach-go/v2/testserver"
)

func TestInsertDelete(t *testing.T) {
	ts, err := testserver.NewTestServer()
	if err != nil {
		t.Fatal(err)
	}
	defer ts.Stop()
	os.Setenv("DATABASE_URL", ts.PGURL().String())
	Connect()
	defer Close()

	toDo := dto.ToDo{Description: "Test" + time.Now().Format(time.RFC3339), Done: true}
	savedToDo, err := Insert(&toDo)
	if err != nil {
		t.Errorf("Error Inserting To Do %v", err)
	}
	if savedToDo.Description != toDo.Description {
		t.Errorf("Wanted %v got %v", toDo.Description, savedToDo.Description)
	}
	if savedToDo.Done != toDo.Done {
		t.Errorf("Wanted %v got %v", toDo.Done, savedToDo.Done)
	}

	foundToDo, err := GetOne(savedToDo.Id)
	if err != nil {
		t.Errorf("Error finding To Do %v", err)
	}
	if foundToDo.Description != toDo.Description {
		t.Errorf("Wanted %v got %v", toDo.Description, savedToDo.Description)
	}
	if foundToDo.Done != toDo.Done {
		t.Errorf("Wanted %v got %v", toDo.Done, savedToDo.Done)
	}

	toDos, err := GetAll()
	found := false
	if err != nil {
		t.Errorf("Error finding To Dos %v", err)
	}
	for _, todo := range *toDos {
		if todo.Id == savedToDo.Id {
			found = true
			if todo.Description != toDo.Description {
				t.Errorf("Wanted %v got %v", toDo.Description, todo.Description)
			}
			if todo.Done != toDo.Done {
				t.Errorf("Wanted %v got %v", toDo.Done, todo.Done)
			}
			break
		}
	}
	if !found {
		t.Errorf("To Do not found in GetAll")
	}

	err = Delete(savedToDo.Id)
	if err != nil {
		t.Errorf("Error Deleting To Do %v", err)
	}
}
