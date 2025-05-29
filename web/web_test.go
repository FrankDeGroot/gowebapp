package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"todo-app/db/mocks"
	"todo-app/dto"

	"github.com/stretchr/testify/assert"
)

var (
	task      dto.Task       = dto.Task{Description: "test" + time.Now().Format(time.RFC3339), Done: false}
	savedTask *dto.SavedTask = &dto.SavedTask{Id: "1", Task: task}
)

func TestPostTask(t *testing.T) {
	mockRepo := new(mocks.MockTaskDb)
	repo = mockRepo
	srv := httptest.NewServer(nil)
	defer srv.Close()

	mockRepo.On("Insert", &task).Return(savedTask, nil)

	taskBytes, err := json.Marshal(task)
	assert.NoError(t, err)
	post, err := http.Post(srv.URL+TASKS_PATH, CONTENT_TYPE_JSON, bytes.NewBuffer(taskBytes))
	assert.NoError(t, err)
	postTask := dto.SavedTask{}
	json.NewDecoder(post.Body).Decode(&postTask)
	assert.Equal(t, task.Description, postTask.Description)
	assert.Equal(t, task.Done, postTask.Done)
}

func TestGetOneTask(t *testing.T) {
	mockRepo := new(mocks.MockTaskDb)
	repo = mockRepo
	srv := httptest.NewServer(nil)
	defer srv.Close()

	mockRepo.On("GetOne", "1").Return(savedTask, nil)

	getOne, err := http.Get(srv.URL + TASKS_PATH + "/" + savedTask.Id)
	assert.NoError(t, err)
	getTask := dto.SavedTask{}
	json.NewDecoder(getOne.Body).Decode(&getTask)
	assert.NoError(t, err)
	assert.Equal(t, task.Description, getTask.Description)
	assert.Equal(t, task.Done, getTask.Done)
}

func TestGetAllTasks(t *testing.T) {
	mockRepo := new(mocks.MockTaskDb)
	repo = mockRepo
	srv := httptest.NewServer(nil)
	defer srv.Close()

	mockRepo.On("GetAll").Return(&[]dto.SavedTask{*savedTask}, nil)

	getAll, err := http.Get(srv.URL + TASKS_PATH)
	assert.NoError(t, err)
	getTasks := []dto.SavedTask{}
	json.NewDecoder(getAll.Body).Decode(&getTasks)
	found := false
	for _, getTask := range getTasks {
		if getTask.Id == savedTask.Id {
			found = true
			assert.Equal(t, task.Description, getTask.Description)
			assert.Equal(t, task.Done, getTask.Done)
			break
		}
	}
	assert.True(t, found, "Did not find the task")
}

func TestDeleteTask(t *testing.T) {
	mockRepo := new(mocks.MockTaskDb)
	repo = mockRepo
	srv := httptest.NewServer(nil)
	defer srv.Close()

	mockRepo.On("Delete", savedTask.Id).Return(nil)

	delReq, err := http.NewRequest(http.MethodDelete, srv.URL+TASKS_PATH+"/"+savedTask.Id, nil)
	assert.NoError(t, err)
	delRes, err := http.DefaultClient.Do(delReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, delRes.StatusCode)

	mockRepo.AssertCalled(t, "Delete", savedTask.Id)
}
