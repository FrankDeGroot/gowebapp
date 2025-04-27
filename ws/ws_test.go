package ws

import (
	"net/http/httptest"
	"strings"
	"testing"
	"todo-app/act"
	"todo-app/dto"
	"todo-app/ws/mocks"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/stretchr/testify/assert"
)

var (
	mockProducer = new(mocks.MockProducer)
	mockConsumer = new(mocks.MockConsumer)
)

func TestWebSocketConnection(t *testing.T) {
	assert.NotNil(t, Open(mockProducer, mockConsumer))
	defer Close()
	srv := httptest.NewServer(nil)
	defer srv.Close()

	task := &act.TaskAction{Verb: act.Post,
		SavedTask: dto.SavedTask{Id: "123",
			Task: dto.Task{Description: "123", Done: false}}}
	mockConsumer.On("Consume").Return(task, nil)

	url := "ws" + strings.TrimPrefix(srv.URL, "http") + WS_PATH
	t.Log(url)
	c, _, err := websocket.Dial(t.Context(), url, nil)
	assert.NoError(t, err)
	defer c.CloseNow()

	recvTask := &act.TaskAction{}
	err = wsjson.Read(t.Context(), c, recvTask)
	assert.NoError(t, err)
	assert.NotNil(t, recvTask)
	assert.Equal(t, task.Id, recvTask.Id)
	assert.Equal(t, task.Description, recvTask.Description)
}
