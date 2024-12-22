package db

import (
	"testing"
	"todo-app/dto"
)

func TestInsertDelete(t *testing.T) {
	Connect()
	defer Close()

	toDo := dto.ToDo{Description: "Test", Done: true}
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
		t.Errorf("Error finding To Do %v", savedToDo)
	}
	if foundToDo.Description != toDo.Description {
		t.Errorf("Wanted %v got %v", toDo.Description, savedToDo.Description)
	}
	if foundToDo.Done != toDo.Done {
		t.Errorf("Wanted %v got %v", toDo.Done, savedToDo.Done)
	}

	err = Delete(savedToDo.Id)
	if err != nil {
		t.Errorf("Error Deleting To Do %v", err)
	}
}
