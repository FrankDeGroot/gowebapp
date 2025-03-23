package events

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"todo-app/dto"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var p *kafka.Producer

func TodoProducer() (chan dto.SavedToDo, error) {
	var err error
	p, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_BROKER")})
	if err != nil {
		return nil, err
	}
	ch := make(chan dto.SavedToDo)
	go func() {
		for e := range p.Events() {
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
		defer p.Close()
		for todo := range ch {
			todoJSON, err := json.Marshal(todo)
			if err != nil {
				log.Fatalf("Failed to encode todo to JSON: %v\n", err)
			}
			err = p.Produce(&kafka.Message{
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
		p.Flush(15 * 1000)
	}()
	return ch, nil
}
