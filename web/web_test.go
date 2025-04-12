package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"todo-app/dto"
	"todo-app/web/mocks"

	"github.com/stretchr/testify/assert"
)

var (
	todo      dto.Todo            = dto.Todo{Description: "test" + time.Now().Format(time.RFC3339), Done: false}
	savedTodo *dto.SavedTodo      = &dto.SavedTodo{Id: "1", Todo: todo}
	mockRepo  *mocks.MockTodoRepo = new(mocks.MockTodoRepo)
)

func init() {
	repo = mockRepo
	setHandlers()
}

func TestPostTodo(t *testing.T) {
	srv := httptest.NewServer(nil)
	defer srv.Close()

	mockRepo.On("Insert", &todo).Return(savedTodo, nil)

	todoBytes, err := json.Marshal(todo)
	assert.NoError(t, err)
	post, err := http.Post(srv.URL+TODO_PATH, CONTENT_TYPE_JSON, bytes.NewBuffer(todoBytes))
	assert.NoError(t, err)
	postToDo := dto.SavedTodo{}
	json.NewDecoder(post.Body).Decode(&postToDo)
	assert.Equal(t, todo.Description, postToDo.Description)
	assert.Equal(t, todo.Done, postToDo.Done)
}

func TestGetOneTodo(t *testing.T) {
	srv := httptest.NewServer(nil)
	defer srv.Close()

	mockRepo.On("GetOne", "1").Return(savedTodo, nil)

	getOne, err := http.Get(srv.URL + TODO_PATH + "/" + savedTodo.Id)
	assert.NoError(t, err)
	getToDo := dto.SavedTodo{}
	json.NewDecoder(getOne.Body).Decode(&getToDo)
	assert.NoError(t, err)
	assert.Equal(t, todo.Description, getToDo.Description)
	assert.Equal(t, todo.Done, getToDo.Done)
}

func TestGetAllTodos(t *testing.T) {
	srv := httptest.NewServer(nil)
	defer srv.Close()

	mockRepo.On("GetAll").Return(&[]dto.SavedTodo{*savedTodo}, nil)

	getAll, err := http.Get(srv.URL + TODO_PATH)
	assert.NoError(t, err)
	getToDos := []dto.SavedTodo{}
	json.NewDecoder(getAll.Body).Decode(&getToDos)
	found := false
	for _, getToDo := range getToDos {
		if getToDo.Id == savedTodo.Id {
			found = true
			assert.Equal(t, todo.Description, getToDo.Description)
			assert.Equal(t, todo.Done, getToDo.Done)
			break
		}
	}
	assert.True(t, found, "Did not find the todo")
}

func TestDeleteTodo(t *testing.T) {
	srv := httptest.NewServer(nil)
	defer srv.Close()

	mockRepo.On("Delete", savedTodo.Id).Return(nil)

	delReq, err := http.NewRequest(http.MethodDelete, srv.URL+TODO_PATH+"/"+savedTodo.Id, nil)
	assert.NoError(t, err)
	delRes, err := http.DefaultClient.Do(delReq)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, delRes.StatusCode)

	mockRepo.AssertCalled(t, "Delete", savedTodo.Id)
}
