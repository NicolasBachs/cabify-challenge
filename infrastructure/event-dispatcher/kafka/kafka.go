package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/app"
	eventDispatcher "gitlab-hiring.cabify.tech/cabify/interviewing/car-pooling-challenge-go/infrastructure/event-dispatcher"
)

type KafkaEventDispatcher struct {
	writer *kafka.Writer
}

func NewKafkaEventDispatcher(brokerURL string) (eventDispatcher.EventDispatcher, error) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerURL},
	})

	return &KafkaEventDispatcher{
		writer: writer,
	}, nil
}

var kafkaName = "KAFKA"

func (d *KafkaEventDispatcher) Dispatch(ctx context.Context, topic string, event interface{}) error {
	rawEventBytes, err := json.Marshal(event)

	if err != nil {
		return fmt.Errorf("Failed to marshal event to dispatch, error: %s", err.Error())
	}

	jsonEventBytes, err := json.Marshal(string(rawEventBytes))

	if err != nil {
		return fmt.Errorf("Failed to marshal event json string to dispatch: %s", err.Error())
	}

	app.Logger.Debug(kafkaName, "Dispatching event for topic: %s, event: %s", topic, string(rawEventBytes))

	msg := kafka.Message{
		Topic: topic,
		Value: jsonEventBytes,
	}

	if err := d.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	app.Logger.Debug(kafkaName, "Event dispatched to topic: %s with value: %s", msg.Topic, string(rawEventBytes))

	return nil
}

func (d *KafkaEventDispatcher) Close() error {
	if err := d.writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}
