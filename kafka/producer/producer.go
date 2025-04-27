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

func Open(producerTopic string) (*Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_BROKER")})
	if err != nil {
		return nil, err
	}
	return &Producer{producer: producer, topic: producerTopic}, nil
}

func (p *Producer) Produce(task *act.TaskAction) error {
	taskJson, err := json.Marshal(task)
	if err != nil {
		return err
	}
	deliveryChan := make(chan kafka.Event)
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(task.Id),
		Value: taskJson,
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
