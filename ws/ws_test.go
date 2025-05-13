package ws

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"todo-app/act"
	dbm "todo-app/db/mocks"
	"todo-app/dto"
	wsm "todo-app/ws/mocks"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/stretchr/testify/assert"
)

func TestConsumeReceive(t *testing.T) {
	_, mockConsumer, _, _ := setup()
	defer Close()
	srv := httptest.NewServer(nil)
	defer srv.Close()

	task := taskAction(act.Post)
	consumeWait := make(chan time.Time)
	defer close(consumeWait)
	mockConsumer.On("Consume").WaitUntil(consumeWait).Return(task, nil)

	c := openConn(t, srv)
	defer c.Close(websocket.StatusGoingAway, t.Name())

	consumeWait <- time.Now()
	recvTask := &act.TaskAction{}
	err := wsjson.Read(t.Context(), c, recvTask)
	assert.NoError(t, err)
	assert.NotNil(t, recvTask)
	assert.Equal(t, task.Id, recvTask.Id)
	assert.Equal(t, task.Description, recvTask.Description)
}

func TestReadApplyPost(t *testing.T) {
	mockProducer, mockConsumer, mockTaskDber, _ := setup()
	defer Close()
	srv := httptest.NewServer(nil)
	defer srv.Close()

	taskAction := taskAction(act.Post)
	consumeWait := make(chan time.Time)
	defer close(consumeWait)
	mockConsumer.On("Consume").WaitUntil(consumeWait).Return(taskAction, nil)

	c := openConn(t, srv)
	defer c.Close(websocket.StatusGoingAway, t.Name())

	mockTaskDber.On("Insert", &taskAction.Task).Return(&taskAction.SavedTask, nil)
	mockProducer.On("Produce", act.Make(act.Post, &taskAction.SavedTask)).Return(nil)

	err := wsjson.Write(t.Context(), c, taskAction)
	assert.NoError(t, err)
}

func TestReadApplyPut(t *testing.T) {
	mockProducer, mockConsumer, mockTaskDber, _ := setup()
	defer Close()
	srv := httptest.NewServer(nil)
	defer srv.Close()

	taskAction := taskAction(act.Put)
	consumeWait := make(chan time.Time)
	defer close(consumeWait)
	mockConsumer.On("Consume").WaitUntil(consumeWait).Return(taskAction, nil)

	mockTaskDber.On("Update", &taskAction.SavedTask).Return(nil)
	mockProducer.On("Produce", act.Make(act.Put, &taskAction.SavedTask)).Return(nil)

	c := openConn(t, srv)
	defer c.Close(websocket.StatusGoingAway, t.Name())

	err := wsjson.Write(t.Context(), c, taskAction)
	assert.NoError(t, err)
}

func TestReadApplyDelete(t *testing.T) {
	mockProducer, mockConsumer, mockTaskDber, _ := setup()
	defer Close()
	srv := httptest.NewServer(nil)
	defer srv.Close()

	taskAction := taskAction(act.Delete)
	consumeWait := make(chan time.Time)
	defer close(consumeWait)
	mockConsumer.On("Consume").WaitUntil(consumeWait).Return(taskAction, nil)

	mockTaskDber.On("Delete", taskAction.Id).Return(nil)
	mockProducer.On("Produce", act.Make(act.Delete, &taskAction.SavedTask)).Return(nil)

	c := openConn(t, srv)
	defer c.Close(websocket.StatusGoingAway, t.Name())

	err := wsjson.Write(t.Context(), c, taskAction)
	assert.NoError(t, err)
}

func setup() (*wsm.MockProducer, *wsm.MockConsumer, *dbm.MockTaskDb, func(*act.TaskAction)) {
	mockProducer := new(wsm.MockProducer)
	mockConsumer := new(wsm.MockConsumer)
	mockTaskDber := new(dbm.MockTaskDb)
	ntfy := Open(mockProducer, mockConsumer, mockTaskDber)
	return mockProducer, mockConsumer, mockTaskDber, ntfy
}

func openConn(t *testing.T, srv *httptest.Server) *websocket.Conn {
	url := "ws" + strings.TrimPrefix(srv.URL, "http") + WS_PATH
	t.Log(url)
	c, _, err := websocket.Dial(t.Context(), url, nil)
	assert.NoError(t, err)
	return c
}

func taskAction(verb act.Verb) *act.TaskAction {
	return &act.TaskAction{Verb: verb,
		SavedTask: dto.SavedTask{Id: "123",
			Task: dto.Task{Description: "123", Done: false}}}
}
