package producer

import (
	"log"
	"todo-app/dto"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	kafkaBroker = "redpanda-0:9092"
	topics      = "todo"
)

var producer *kafka.Producer

func Connect() {
	var err error
	producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaBroker})
	if err != nil {
		log.Fatalf("Unable to connect to broker: %v\n", err)
	}
}

func Close() {
	producer.Close()
}

func Produce(dto.SavedToDo) {

}
