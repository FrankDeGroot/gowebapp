package events

import (
	"encoding/json"
	"log"
	"os"
	"time"
	"todo-app/dto"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ToDoConsumer struct {
	consumer *kafka.Consumer
	topic    string
	name     string
}

func NewToDoConsumer(topic string, name string) (*ToDoConsumer, error) {
	tc := ToDoConsumer{nil, topic, name}
	var err error
	tc.consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKER"),
		"group.id":          name,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	err = tc.consumer.Subscribe(topic, nil)
	if err != nil {
		return nil, err
	}
	return &tc, nil
}

func (tc *ToDoConsumer) Consume() (*dto.SavedToDo, error) {
	for {
		msg, err := tc.consumer.ReadMessage(time.Second)
		if err == nil {
			todo := dto.SavedToDo{}
			err := json.Unmarshal(msg.Value, &todo)
			if err != nil {
				log.Fatalf("Unable to unmarshal %v", msg.Value)
			}
			return &todo, nil
		} else if !err.(kafka.Error).IsTimeout() {
			return nil, err
		}
	}
}

func (tc *ToDoConsumer) Close() {
	tc.consumer.Close()
}
