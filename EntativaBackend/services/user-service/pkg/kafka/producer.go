package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writers map[string]*kafka.Writer
	brokers []string
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string) *Producer {
	return &Producer{
		writers: make(map[string]*kafka.Writer),
		brokers: brokers,
	}
}

// getWriter gets or creates a Kafka writer for a topic
func (p *Producer) getWriter(topic string) *kafka.Writer {
	if writer, exists := p.writers[topic]; exists {
		return writer
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(p.brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
		Async:        true, // Non-blocking writes
	}

	p.writers[topic] = writer
	return writer
}

// PublishEvent publishes an event to Kafka
func (p *Producer) PublishEvent(topic string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	writer := p.getWriter(topic)

	message := kafka.Message{
		Key:   []byte(fmt.Sprintf("%d", time.Now().UnixNano())),
		Value: data,
		Time:  time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = writer.WriteMessages(ctx, message)
	if err != nil {
		log.Printf("Failed to publish event to Kafka: %v", err)
		return err
	}

	return nil
}

// Close closes all Kafka writers
func (p *Producer) Close() error {
	for _, writer := range p.writers {
		if err := writer.Close(); err != nil {
			log.Printf("Error closing Kafka writer: %v", err)
		}
	}
	return nil
}
