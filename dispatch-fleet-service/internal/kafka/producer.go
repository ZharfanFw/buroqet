package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type TrackingPayload struct {
	AWB       string    `json:"awb"`
	Status    string    `json:"status"`
	HubID     string    `json:"hub_id"`
	Location  string    `json:"location"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
}

type Producer interface {
	PublishPackageDispatched(ctx context.Context, awb string, location string) error
	Close() error
}

type producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) Producer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  "package.dispatched",
		AllowAutoTopicCreation: true,
	}

	return &producer{writer: w}
}

func (p *producer) PublishPackageDispatched(ctx context.Context, awb string, location string) error {
	payload := TrackingPayload{
		AWB:       awb,
		Status:    "OUT_FOR_DELIVERY",
		Location:  location,
		Timestamp: time.Now(),
		Source:    "Dispatch",
	}

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(awb),
		Value: bytes,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

func (p *producer) Close() error {
	return p.writer.Close()
}
