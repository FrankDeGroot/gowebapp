package events

import (
	"encoding/json"
	"log"
	"os"
	"time"
	"todo-app/dto"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var c *kafka.Consumer

func TodoConsumer(quit chan bool) (chan dto.SavedToDo, error) {
	var err error
	c, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKER"),
		"group.id":          toDoTopic,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	err = c.Subscribe(toDoTopic, nil)
	if err != nil {
		return nil, err
	}
	ch := make(chan dto.SavedToDo)
	go func() {
		defer c.Close()
		defer close(ch)
		for {
			select {
			case <-quit:
				return
			default:
			}
			msg, err := c.ReadMessage(time.Second)
			if err == nil {
				todo := dto.SavedToDo{}
				err := json.Unmarshal(msg.Value, &todo)
				if err != nil {
					log.Fatalf("Unable to unmarshal %v", msg.Value)
				}
				ch <- todo
			} else if !err.(kafka.Error).IsTimeout() {
				log.Fatalf("Consumer error: %v\n", err)
			}
		}
	}()
	return ch, nil
}
