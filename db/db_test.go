package db

import (
	"os"
	"testing"
	"time"
	"todo-app/dto"

	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/stretchr/testify/assert"
)

func TestInsertDelete(t *testing.T) {
	ts, err := testserver.NewTestServer()
	assert.NoError(t, err)
	defer ts.Stop()

	os.Setenv("DATABASE_URL", ts.PGURL().String())
	db, err := Connect()
	assert.NoError(t, err)
	defer db.Close()

	todo := dto.Todo{Description: "Test" + time.Now().Format(time.RFC3339), Done: true}
	savedTodo, err := db.Insert(&todo)
	assert.NoError(t, err)
	assert.Equal(t, todo.Description, savedTodo.Description)
	assert.Equal(t, todo.Done, savedTodo.Done)

	foundToDo, err := db.GetOne(savedTodo.Id)
	assert.NoError(t, err)
	assert.Equal(t, todo.Description, foundToDo.Description)
	assert.Equal(t, todo.Done, foundToDo.Done)

	todos, err := db.GetAll()
	assert.NoError(t, err)

	found := false
	for _, todo := range *todos {
		if todo.Id == savedTodo.Id {
			found = true
			assert.Equal(t, savedTodo.Description, todo.Description)
			assert.Equal(t, savedTodo.Done, todo.Done)
			break
		}
	}
	assert.True(t, found, "To Do not found in GetAll")

	err = db.Delete(savedTodo.Id)
	assert.NoError(t, err)
}
