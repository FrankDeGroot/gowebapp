package ws

import (
	"testing"
	"todo-app/ws/mocks"
)

var mockConsumer = new(mocks.MockConsumer)

func TestInit(t *testing.T) {
	if Init(mockConsumer) == nil {
		t.Fatal("Init should not return nil")
	}
}
