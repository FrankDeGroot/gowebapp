package producer

import (
	"encoding/json"
	"os"
	"todo-app/act"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
	topic    string
}

func Connect(producerTopic string) (*Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_BROKER")})
	if err != nil {
		return nil, err
	}
	return &Producer{producer: producer, topic: producerTopic}, nil
}

func (p *Producer) Produce(todo *act.TodoAction) error {
	todoJson, err := json.Marshal(todo)
	if err != nil {
		return err
	}
	deliveryChan := make(chan kafka.Event)
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(todo.Id),
		Value: todoJson,
	}, deliveryChan)
	if err != nil {
		return err
	}
	e := <-deliveryChan
	switch ev := e.(type) {
	case *kafka.Message:
		return ev.TopicPartition.Error
	}
	return nil
}

func (p *Producer) Close() {
	p.producer.Flush(500)
	p.producer.Close()
}
