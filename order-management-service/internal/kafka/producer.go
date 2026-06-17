package kafka

import (
	"context"
	"encoding/json"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

//go:generate mockgen -source=producer.go -destination=../../mock/mock_kafka_producer.go -package=mock

// Producer defines the contract for publishing events to Kafka.
// Keeping it as an interface lets the service layer be unit-tested
// without a real Kafka broker.
type Producer interface {
	PublishOrderCreated(ctx context.Context, event OrderCreatedEvent) error
	Close() error
}

// OrderCreatedEvent is the payload published to the "order.created" Kafka topic.
// Dispatch Service listens to this topic to schedule a pickup courier.
type OrderCreatedEvent struct {
	AWBNumber     string    `json:"awb_number"`
	TransactionID string    `json:"transaction_id"`
	SenderName    string    `json:"sender_name"`
	SenderAddress string    `json:"sender_address"`
	OriginCity    string    `json:"origin_city"`
	ReceiverName  string    `json:"receiver_name"`
	DestCity      string    `json:"dest_city"`
	ServiceType   string    `json:"service_type"`
	TotalPrice    float64   `json:"total_price"`
	CreatedAt     time.Time `json:"created_at"`
}

// kafkaProducer is the real kafka-go backed implementation.
type kafkaProducer struct {
	writer *kafkago.Writer
}

// NewProducer creates a new Kafka producer pointing at the given broker addresses
// and writing to the "order.created" topic.
func NewProducer(brokers []string) Producer {
	writer := &kafkago.Writer{
		Addr:         kafkago.TCP(brokers...),
		Topic:        "order.created",
		Balancer:     &kafkago.LeastBytes{},
		RequiredAcks: kafkago.RequireOne,
	}
	return &kafkaProducer{writer: writer}
}

// PublishOrderCreated serialises the event to JSON and writes it to Kafka.
// The AWB number is used as the message key so all events for a single order
// land on the same partition (preserving order).
func (p *kafkaProducer) PublishOrderCreated(ctx context.Context, event OrderCreatedEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafkago.Message{
		Key:   []byte(event.AWBNumber),
		Value: payload,
	})
}

// Close flushes pending messages and releases the underlying connection.
func (p *kafkaProducer) Close() error {
	return p.writer.Close()
}