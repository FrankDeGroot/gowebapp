package ws

import (
	"testing"
	"todo-app/ws/mocks"
)

var (
	mockProducer = new(mocks.MockProducer)
	mockConsumer = new(mocks.MockConsumer)
)

func TestInit(t *testing.T) {
	if Init(mockProducer, mockConsumer) == nil {
		t.Fatal("Init should not return nil")
	}
}
