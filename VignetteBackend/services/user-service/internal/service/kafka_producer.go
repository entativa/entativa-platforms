package service

import (
	"socialink/user-service/pkg/kafka"
)

// KafkaProducer wraps the Kafka producer for the service layer
type KafkaProducer struct {
	producer *kafka.Producer
}

// NewKafkaProducer creates a new Kafka producer wrapper
func NewKafkaProducer(brokers []string) *KafkaProducer {
	if len(brokers) == 0 {
		return nil // Kafka is optional
	}

	producer := kafka.NewProducer(brokers)
	return &KafkaProducer{producer: producer}
}

// PublishEvent publishes an event to a Kafka topic
func (k *KafkaProducer) PublishEvent(topic string, event interface{}) error {
	if k == nil || k.producer == nil {
		return nil // Kafka not configured
	}
	return k.producer.PublishEvent(topic, event)
}

// Close closes the Kafka producer
func (k *KafkaProducer) Close() error {
	if k == nil || k.producer == nil {
		return nil
	}
	return k.producer.Close()
}
