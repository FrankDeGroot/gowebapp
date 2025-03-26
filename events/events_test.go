package events

import (
	"testing"
	"time"
	"todo-app/dto"
)

func TestProduceConsume(t *testing.T) {
	p, err := NewToDoProducer()
	if err != nil {
		t.Fail()
	}
	defer p.Close()
	c, err := NewToDoConsumer()
	if err != nil {
		t.Fail()
	}
	defer c.Stop()
	p.Produce(dto.SavedToDo{
		Id: "123",
		ToDo: dto.ToDo{
			Description: "test",
			Done:        false,
		},
	})
	time.Sleep(time.Second)
	todo := c.Receive()
	if todo.Description != "test" {
		t.Errorf("Wanted %v got %v", "test", todo.Description)
	}
}
