package events

import (
	"fmt"
	"testing"
	"time"
	"todo-app/dto"
	"todo-app/events/admin"
	"todo-app/events/consumer"
	"todo-app/events/producer"
)

func TestProduceConsume(t *testing.T) {
	topic := fmt.Sprintf("todo%v", time.Now().Format("_2006_01_02_15_04_05"))
	t.Log(topic)
	if err := producer.Connect(topic); err != nil {
		t.Fatal(err)
	}
	defer producer.Close()
	if err := consumer.Connect(topic, topic); err != nil {
		t.Fatal(err)
	}
	defer consumer.Close()
	defer admin.DeleteTopic(topic)
	if err := producer.Produce(&dto.ToDoEvent{
		Action: dto.ActionAdd,
		SavedToDo: dto.SavedToDo{
			Id: "123",
			ToDo: dto.ToDo{
				Description: "test",
				Done:        false,
			},
		},
	}); err != nil {
		t.Fatal(err)
	}
	todo, err := consumer.Consume()
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
