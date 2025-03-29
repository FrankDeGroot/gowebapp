package events

import (
	"fmt"
	"testing"
	"time"
	"todo-app/dto"
)

func TestProduceConsume(t *testing.T) {
	topic := fmt.Sprintf("todo%v", time.Now().Format("_2006_01_02_15_04_05"))
	t.Log(topic)
	p, err := NewToDoProducer(topic)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	c, err := NewToDoConsumer(topic, topic)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	defer DeleteTopic(topic)
	err = p.Produce(dto.SavedToDo{
		Id: "123",
		ToDo: dto.ToDo{
			Description: "test",
			Done:        false,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	todo, err := c.Consume()
	if err != nil {
		t.Fatal(err)
	}
	if todo.Id != "123" {
		t.Fatalf("Wanted %v got %v", "123", todo.Id)
	}
	if todo.Description != "test" {
		t.Fatalf("Wanted %v got %v", "test", todo.Description)
	}
}
