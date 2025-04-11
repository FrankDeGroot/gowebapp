package consumer

import (
	"encoding/json"
	"os"
	"time"
	"todo-app/act"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var consumer *kafka.Consumer

func Connect(topic string, name string) error {
	var err error
	consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKER"),
		"group.id":          name,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return err
	}
	err = consumer.Subscribe(topic, nil)
	if err != nil {
		return err
	}
	return nil
}

func Consume() (*act.TodoAction, error) {
	for {
		msg, err := consumer.ReadMessage(time.Second)
		if err == nil {
			todo := act.TodoAction{}
			err := json.Unmarshal(msg.Value, &todo)
			if err != nil {
				return nil, err
			}
			return &todo, nil
		} else if !err.(kafka.Error).IsTimeout() {
			return nil, err
		}
	}
}

func Close() {
	consumer.Close()
}
