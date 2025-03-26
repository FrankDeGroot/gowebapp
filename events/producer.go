package events

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"todo-app/dto"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ToDoProducer struct {
	p  *kafka.Producer
	ch chan dto.SavedToDo
}

func NewToDoProducer() (*ToDoProducer, error) {
	tp := ToDoProducer{nil, make(chan dto.SavedToDo)}
	var err error
	tp.p, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_BROKER")})
	if err != nil {
		return nil, err
	}
	go func() {
		for e := range tp.p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	go func() {
		defer tp.p.Close()
		for todo := range tp.ch {
			todoJSON, err := json.Marshal(todo)
			if err != nil {
				log.Fatalf("Failed to encode todo to JSON: %v\n", err)
			}
			err = tp.p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &toDoTopic,
					Partition: kafka.PartitionAny,
				},
				Key:   []byte(todo.Id),
				Value: todoJSON,
			}, nil)
			if err != nil {
				log.Fatalf("Error producing message: %v\n", err)
			}
		}
		tp.p.Flush(500)
	}()
	return &tp, nil
}

func (tp *ToDoProducer) Produce(toDo dto.SavedToDo) {
	tp.ch <- toDo
}

func (tp *ToDoProducer) Close() {
	close(tp.ch)
}
