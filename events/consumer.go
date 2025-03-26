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
	c    *kafka.Consumer
	quit chan struct{}
	ch   chan dto.SavedToDo
}

func NewToDoConsumer() (*ToDoConsumer, error) {
	tc := ToDoConsumer{nil, make(chan struct{}), make(chan dto.SavedToDo)}
	var err error
	tc.c, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKER"),
		"group.id":          toDoTopic,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	err = tc.c.Subscribe(toDoTopic, nil)
	if err != nil {
		return nil, err
	}
	go func() {
		defer tc.c.Close()
		defer close(tc.ch)
		for {
			select {
			case <-tc.quit:
				return
			default:
			}
			msg, err := tc.c.ReadMessage(time.Second)
			if err == nil {
				todo := dto.SavedToDo{}
				err := json.Unmarshal(msg.Value, &todo)
				if err != nil {
					log.Fatalf("Unable to unmarshal %v", msg.Value)
				}
				tc.ch <- todo
			} else if !err.(kafka.Error).IsTimeout() {
				log.Fatalf("Consumer error: %v\n", err)
			}
		}
	}()
	return &tc, nil
}

func (tc *ToDoConsumer) Receive() dto.SavedToDo {
	return <-tc.ch
}

func (tc *ToDoConsumer) Stop() {
	tc.quit <- struct{}{}
}
