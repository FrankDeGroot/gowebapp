package events

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func DeleteTopic(topic string) error {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_BROKER")})
	if err != nil {
		return err
	}
	defer adminClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err = adminClient.DeleteTopics(ctx, []string{topic}, nil)
	if err != nil {
		return err
	}
	log.Printf("Topic %s deleted successfully", topic)
	return nil
}
