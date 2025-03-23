package events

import (
	"testing"
	"todo-app/dto"
)

func TestProduceConsume(t *testing.T) {
	prodCh, err := TodoProducer()
	if err != nil {
		t.Fail()
	}
	defer close(prodCh)
	quitCh := make(chan bool)
	defer close(quitCh)
	consCh, err := TodoConsumer(quitCh)
	if err != nil {
		t.Fail()
	}
	defer close(consCh)
	prodCh <- dto.SavedToDo{
		Id: "123",
		ToDo: dto.ToDo{
			Description: "test",
			Done:        false,
		},
	}
	todo := <-consCh
	if todo.Description != "test" {
		t.Errorf("Wanted %v got %v", "test", todo.Description)
	}
	quitCh <- true
}
