package consumer

import (
	"encoding/json"
	"os"
	"time"
	"todo-app/act"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	consumer *kafka.Consumer
}

func Open(topic string, name string) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKER"),
		"group.id":          name,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer: consumer}, nil
}

func (c *Consumer) Consume() (*act.TaskAction, error) {
	for {
		msg, err := c.consumer.ReadMessage(time.Second)
		if err == nil {
			task := act.TaskAction{}
			err := json.Unmarshal(msg.Value, &task)
			if err != nil {
				return nil, err
			}
			return &task, nil
		} else if !err.(kafka.Error).IsTimeout() {
			return nil, err
		}
	}
}

func (c *Consumer) Close() {
	c.consumer.Close()
}
