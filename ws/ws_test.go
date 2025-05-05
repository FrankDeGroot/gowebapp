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

var (
	mockProducer = new(wsm.MockProducer)
	mockConsumer = new(wsm.MockConsumer)
	mockTaskDber = new(dbm.MockTaskDb)
	ntfy         func(*act.TaskAction)
)

func init() {
	ntfy = Open(mockProducer, mockConsumer, mockTaskDber)
}

func TestConsumeReceive(t *testing.T) {
	srv := httptest.NewServer(nil)
	defer srv.Close()

	task := postTaskAction()
	mockConsumer.On("Consume").Return(task, nil)

	c := openConn(t, srv)
	defer c.Close(websocket.StatusGoingAway, t.Name())

	recvTask := &act.TaskAction{}
	err := wsjson.Read(t.Context(), c, recvTask)
	assert.NoError(t, err)
	assert.NotNil(t, recvTask)
	assert.Equal(t, task.Id, recvTask.Id)
	assert.Equal(t, task.Description, recvTask.Description)
}

func TestReadApply(t *testing.T) {
	srv := httptest.NewServer(nil)
	defer srv.Close()

	taskAction := postTaskAction()
	wait := make(chan time.Time)
	mockConsumer.On("Consume").WaitUntil(wait).Return(taskAction, nil)

	c := openConn(t, srv)
	defer c.Close(websocket.StatusGoingAway, t.Name())

	mockTaskDber.On("Insert", &taskAction.Task).Return(&taskAction.SavedTask, nil)
	mockProducer.On("Produce", act.Make(act.Post, &taskAction.SavedTask)).Return(nil)

	err := wsjson.Write(t.Context(), c, taskAction)
	assert.NoError(t, err)

	wait <- time.Now()
}

func openConn(t *testing.T, srv *httptest.Server) *websocket.Conn {
	url := "ws" + strings.TrimPrefix(srv.URL, "http") + WS_PATH
	t.Log(url)
	c, _, err := websocket.Dial(t.Context(), url, nil)
	assert.NoError(t, err)
	return c
}

func postTaskAction() *act.TaskAction {
	return &act.TaskAction{Verb: act.Post,
		SavedTask: dto.SavedTask{Id: "123",
			Task: dto.Task{Description: "123", Done: false}}}
}
