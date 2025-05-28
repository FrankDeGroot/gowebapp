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
	"github.com/stretchr/testify/mock"
)

func TestConsumeReceive(t *testing.T) {
	_, mockConsumer, _, _ := setup()
	srv := httptest.NewServer(nil)
	defer srv.Close()

	task := taskAction(act.Post, t.Name())
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

	mockConsumer.AssertExpectations(t)
}

func TestReadApplyPost(t *testing.T) {
	mockProducer, mockConsumer, mockTaskDber, _ := setup()
	srv := httptest.NewServer(nil)
	defer srv.Close()

	taskAction := taskAction(act.Post, t.Name())
	consumeWait := make(chan time.Time)
	defer close(consumeWait)
	mockConsumer.On("Consume").WaitUntil(consumeWait).Return(taskAction, nil)

	c := openConn(t, srv)
	defer c.Close(websocket.StatusGoingAway, t.Name())

	insertCalledChan := make(chan struct{})
	mockTaskDber.On("Insert", &taskAction.Task).
		Return(&taskAction.SavedTask, nil).
		Run(func(args mock.Arguments) {
			close(insertCalledChan)
		}).
		Once()
	produceCalledChan := make(chan struct{})
	mockProducer.On("Produce", act.Make(act.Post, &taskAction.SavedTask)).
		Return(nil).
		Run(func(args mock.Arguments) {
			close(produceCalledChan)
		}).
		Once()

	err := wsjson.Write(t.Context(), c, taskAction)
	assert.NoError(t, err)

	<-insertCalledChan
	<-produceCalledChan

	mockConsumer.AssertExpectations(t)
}

func TestReadApplyPut(t *testing.T) {
	mockProducer, mockConsumer, mockTaskDber, _ := setup()
	srv := httptest.NewServer(nil)
	defer srv.Close()

	taskAction := taskAction(act.Put, t.Name())
	consumeWait := make(chan time.Time)
	defer close(consumeWait)
	mockConsumer.On("Consume").WaitUntil(consumeWait).Return(taskAction, nil)

	updateCalledChan := make(chan struct{})
	mockTaskDber.On("Update", &taskAction.SavedTask).
		Return(nil).
		Run(func(args mock.Arguments) {
			close(updateCalledChan)
		}).
		Once()
	produceCalledChan := make(chan struct{})
	mockProducer.On("Produce", act.Make(act.Put, &taskAction.SavedTask)).
		Return(nil).
		Run(func(args mock.Arguments) {
			close(produceCalledChan)
		}).
		Once()

	c := openConn(t, srv)
	defer c.Close(websocket.StatusGoingAway, t.Name())

	err := wsjson.Write(t.Context(), c, taskAction)
	assert.NoError(t, err)

	<-updateCalledChan
	<-produceCalledChan

	mockConsumer.AssertExpectations(t)
}

func TestReadApplyDelete(t *testing.T) {
	mockProducer, mockConsumer, mockTaskDber, _ := setup()
	srv := httptest.NewServer(nil)
	defer srv.Close()

	taskAction := taskAction(act.Delete, t.Name())
	consumeWait := make(chan time.Time)
	defer close(consumeWait)
	mockConsumer.On("Consume").WaitUntil(consumeWait).Return(taskAction, nil)

	deleteCalledChan := make(chan struct{})
	mockTaskDber.On("Delete", taskAction.Id).
		Return(nil).
		Run(func(args mock.Arguments) {
			close(deleteCalledChan)
		}).
		Once()
	produceCalledChan := make(chan struct{})
	mockProducer.On("Produce", act.Make(act.Delete, &taskAction.SavedTask)).
		Return(nil).
		Run(func(args mock.Arguments) {
			close(produceCalledChan)
		}).
		Once()

	c := openConn(t, srv)
	defer c.Close(websocket.StatusGoingAway, t.Name())

	err := wsjson.Write(t.Context(), c, taskAction)
	assert.NoError(t, err)

	<-deleteCalledChan
	<-produceCalledChan

	mockConsumer.AssertExpectations(t)
}

func TestReadApplyGetAll(t *testing.T) {
	_, mockConsumer, mockTaskDber, _ := setup()
	srv := httptest.NewServer(nil)
	defer srv.Close()

	consumeWait := make(chan time.Time)
	mockConsumer.On("Consume").WaitUntil(consumeWait).
		Return(taskAction(act.Post, "dummy"), nil)
	defer close(consumeWait)

	c := openConn(t, srv)
	defer c.Close(websocket.StatusGoingAway, t.Name())

	getAllAction := &act.TaskAction{Verb: act.Get}

	expectedTasks := []dto.SavedTask{
		{Id: "task-id-1", Task: dto.Task{Description: "First test task", Done: false}},
		{Id: "task-id-2", Task: dto.Task{Description: "Second test task", Done: true}},
	}
	mockTaskDber.On("GetAll").Return(&expectedTasks, nil).Once()

	err := wsjson.Write(t.Context(), c, getAllAction)
	assert.NoError(t, err, "Writing Get All action to websocket should not fail")

	var receivedTasks []dto.SavedTask
	err = wsjson.Read(t.Context(), c, &receivedTasks)
	assert.NoError(t, err, "Reading response from websocket should not fail")
	assert.Equal(t, expectedTasks, receivedTasks, "Received tasks should match the expected tasks")

	mockTaskDber.AssertExpectations(t)
}

func TestReadApplyGetOne(t *testing.T) {
	_, mockConsumer, mockTaskDber, _ := setup()
	srv := httptest.NewServer(nil)
	defer srv.Close()

	consumeWait := make(chan time.Time)
	mockConsumer.On("Consume").WaitUntil(consumeWait).
		Return(taskAction(act.Post, "dummy"), nil)
	defer close(consumeWait)

	c := openConn(t, srv)
	defer c.Close(websocket.StatusGoingAway, t.Name())

	getOneAction := &act.TaskAction{Verb: act.Get, SavedTask: dto.SavedTask{Id: "123"}}

	expectedTask := dto.SavedTask{Id: "task-id-1", Task: dto.Task{Description: "First test task", Done: false}}
	mockTaskDber.On("GetOne", getOneAction.Id).Return(&expectedTask, nil).Once()

	err := wsjson.Write(t.Context(), c, getOneAction)
	assert.NoError(t, err, "Writing Get All action to websocket should not fail")

	var receivedTask dto.SavedTask
	err = wsjson.Read(t.Context(), c, &receivedTask)
	assert.NoError(t, err, "Reading response from websocket should not fail")
	assert.Equal(t, expectedTask, receivedTask, "Received tasks should match the expected tasks")

	mockTaskDber.AssertExpectations(t)
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

func taskAction(verb act.Verb, descr string) *act.TaskAction {
	return &act.TaskAction{Verb: verb,
		SavedTask: dto.SavedTask{Id: "123",
			Task: dto.Task{Description: descr, Done: false}}}
}
