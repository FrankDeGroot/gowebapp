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

	todo := &act.TodoAction{Action: act.Add,
		SavedTodo: dto.SavedTodo{Id: "123",
			Todo: dto.Todo{Description: "123", Done: false}}}
	mockConsumer.On("Consume").Return(todo, nil)

	url := "ws" + strings.TrimPrefix(srv.URL, "http") + WS_TODO_PATH
	t.Log(url)
	c, _, err := websocket.Dial(t.Context(), url, nil)
	assert.NoError(t, err)
	defer c.CloseNow()

	recvTodo := &act.TodoAction{}
	err = wsjson.Read(t.Context(), c, recvTodo)
	assert.NoError(t, err)
	assert.NotNil(t, recvTodo)
	assert.Equal(t, todo.Id, recvTodo.Id)
	assert.Equal(t, todo.Description, recvTodo.Description)
}
