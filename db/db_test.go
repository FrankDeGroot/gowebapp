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

	todo := dto.Todo{Description: "Test" + time.Now().Format(time.RFC3339), Done: true}
	savedTodo, err := Insert(&todo)
	if err != nil {
		t.Fatalf("Error Inserting To Do %v", err)
	}
	if savedTodo.Description != todo.Description {
		t.Fatalf("Wanted %v got %v", todo.Description, savedTodo.Description)
	}
	if savedTodo.Done != todo.Done {
		t.Fatalf("Wanted %v got %v", todo.Done, savedTodo.Done)
	}

	foundToDo, err := GetOne(savedTodo.Id)
	if err != nil {
		t.Fatalf("Error finding To Do %v", err)
	}
	if foundToDo.Description != todo.Description {
		t.Fatalf("Wanted %v got %v", todo.Description, savedTodo.Description)
	}
	if foundToDo.Done != todo.Done {
		t.Fatalf("Wanted %v got %v", todo.Done, savedTodo.Done)
	}

	todos, err := GetAll()
	found := false
	if err != nil {
		t.Fatalf("Error finding To Dos %v", err)
	}
	for _, todo := range *todos {
		if todo.Id == savedTodo.Id {
			found = true
			if todo.Description != todo.Description {
				t.Fatalf("Wanted %v got %v", todo.Description, todo.Description)
			}
			if todo.Done != todo.Done {
				t.Fatalf("Wanted %v got %v", todo.Done, todo.Done)
			}
			break
		}
	}
	if !found {
		t.Fatalf("To Do not found in GetAll")
	}

	err = Delete(savedTodo.Id)
	if err != nil {
		t.Fatalf("Error Deleting To Do %v", err)
	}
}
