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
	db, err := Open()
	assert.NoError(t, err)
	defer db.Close()

	task := dto.Task{Description: "Test" + time.Now().Format(time.RFC3339), Done: true}
	savedTask, err := db.Insert(&task)
	assert.NoError(t, err)
	assert.Equal(t, task.Description, savedTask.Description)
	assert.Equal(t, task.Done, savedTask.Done)

	foundTask, err := db.GetOne(savedTask.Id)
	assert.NoError(t, err)
	assert.Equal(t, task.Description, foundTask.Description)
	assert.Equal(t, task.Done, foundTask.Done)

	tasks, err := db.GetAll()
	assert.NoError(t, err)

	found := false
	for _, task := range *tasks {
		if task.Id == savedTask.Id {
			found = true
			assert.Equal(t, savedTask.Description, task.Description)
			assert.Equal(t, savedTask.Done, task.Done)
			break
		}
	}
	assert.True(t, found, "To Do not found in GetAll")

	err = db.Delete(savedTask.Id)
	assert.NoError(t, err)
}
