package events

import (
	"encoding/json"
	"os"
	"todo-app/dto"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ToDoProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewToDoProducer(topic string) (*ToDoProducer, error) {
	tp := ToDoProducer{nil, topic}
	var err error
	tp.producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_BROKER")})
	if err != nil {
		return nil, err
	}
	return &tp, nil
}

func (tp *ToDoProducer) Produce(toDo dto.SavedToDo) error {
	toDoJson, err := json.Marshal(toDo)
	if err != nil {
		return err
	}
	deliveryChan := make(chan kafka.Event)
	err = tp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &tp.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(toDo.Id),
		Value: toDoJson,
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

func (tp *ToDoProducer) Close() {
	tp.producer.Flush(500)
	tp.producer.Close()
}
