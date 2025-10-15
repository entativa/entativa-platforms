package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
		Compression:  kafka.Snappy,
	}

	return &Producer{writer: writer}
}

func (p *Producer) PublishEvent(ctx context.Context, key string, data map[string]interface{}) error {
	value, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	message := kafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

// EventPublisherImpl implements service.EventPublisher
type EventPublisherImpl struct {
	producer *Producer
}

func NewEventPublisher(brokers []string, topic string) *EventPublisherImpl {
	return &EventPublisherImpl{
		producer: NewProducer(brokers, topic),
	}
}

func (e *EventPublisherImpl) PublishUserEvent(ctx context.Context, eventType string, data map[string]interface{}) error {
	eventData := make(map[string]interface{})
	for k, v := range data {
		eventData[k] = v
	}
	eventData["event_type"] = eventType
	eventData["timestamp"] = time.Now().Unix()

	return e.producer.PublishEvent(ctx, eventType, eventData)
}
