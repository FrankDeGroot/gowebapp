package producer

import (
	"encoding/json"
	"log"
	"os"
	"todo-app/act"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var (
	producer *kafka.Producer
	topic    string
)

func Connect(producerTopic string) error {
	topic = producerTopic
	var err error
	producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_BROKER")})
	return err
}

func Produce(todo *act.TodoAction) error {
	if producer == nil {
		log.Printf("Producer not connected\n")
		return nil
	}
	todoJson, err := json.Marshal(todo)
	if err != nil {
		return err
	}
	deliveryChan := make(chan kafka.Event)
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
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

func Close() {
	producer.Flush(500)
	producer.Close()
}
