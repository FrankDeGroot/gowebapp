package kafka

import (
	"fmt"
	"testing"
	"time"
	"todo-app/act"
	"todo-app/dto"
	"todo-app/kafka/admin"
	"todo-app/kafka/consumer"
	"todo-app/kafka/producer"

	"github.com/stretchr/testify/assert"
)

func TestProduceConsume(t *testing.T) {
	topic := fmt.Sprintf("todo%v", time.Now().Format("_2006_01_02_15_04_05"))
	t.Log(topic)
	if err := producer.Connect(topic); err != nil {
		t.Fatal(err)
	}
	defer producer.Close()
	c, err := consumer.Connect(topic, topic)
	assert.NoError(t, err)
	defer c.Close()
	defer admin.DeleteTopic(topic)
	if err := producer.Produce(act.Make(act.Add, &dto.SavedTodo{
		Id: "123",
		Todo: dto.Todo{
			Description: "test",
			Done:        false,
		},
	})); err != nil {
		t.Fatal(err)
	}
	todo, err := c.Consume()
	assert.NoError(t, err)
	if todo.Id != "123" {
		t.Fatalf("Wanted %v got %v", "123", todo.Id)
	}
	if todo.Description != "test" {
		t.Fatalf("Wanted %v got %v", "test", todo.Description)
	}
}
