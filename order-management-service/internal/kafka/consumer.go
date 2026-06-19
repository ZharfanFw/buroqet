package kafka

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"order-management-service/internal/domain"
	"order-management-service/internal/repository"

	kafkago "github.com/segmentio/kafka-go"
)

// TrackingUpdatedEvent is the payload consumed from "tracking.updated" topic
type TrackingUpdatedEvent struct {
	AWB    string `json:"awb"`
	Status string `json:"status"`
}

// Consumer defines the consumer interface
type Consumer interface {
	Start(ctx context.Context)
	Close() error
}

type kafkaConsumer struct {
	reader *kafkago.Reader
	repo   repository.OrderRepository
}

// NewConsumer creates a new instance of Kafka Consumer
func NewConsumer(brokers []string, groupID string, repo repository.OrderRepository) Consumer {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       "tracking.updated",
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		MaxWait:     1 * time.Second,
		StartOffset: kafkago.FirstOffset,
	})

	return &kafkaConsumer{
		reader: reader,
		repo:   repo,
	}
}

// Start runs the message consumption loop in a blocking manner.
// Should be run inside a goroutine.
func (c *kafkaConsumer) Start(ctx context.Context) {
	log.Println("[KAFKA-CONSUMER] Started listening to topic 'tracking.updated'...")
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("[KAFKA-CONSUMER] Stopping consumer loop due to context cancel")
				return
			}
			log.Printf("[KAFKA-CONSUMER] Error reading message: %v\n", err)
			continue
		}

		var event TrackingUpdatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("[KAFKA-CONSUMER] Error unmarshalling payload: %v\n", err)
			continue
		}

		if event.AWB == "" || event.Status == "" {
			log.Println("[KAFKA-CONSUMER] Invalid event payload: missing AWB or Status")
			continue
		}

		// Map status tracking to OrderStatus domain
		var targetStatus domain.OrderStatus
		statusUpper := strings.ToUpper(event.Status)
		switch statusUpper {
		case "INBOUND":
			targetStatus = domain.StatusPickedUp
		case "ON_TRANSIT":
			targetStatus = domain.StatusOnTransit
		case "AT_HUB":
			targetStatus = domain.StatusAtDestinationHub
		case "OUT_FOR_DELIVERY":
			targetStatus = domain.StatusOutForDelivery
		case "DELIVERED":
			targetStatus = domain.StatusDelivered
		case "FAILED":
			targetStatus = domain.StatusFailed
		case "RETURNED":
			targetStatus = domain.StatusReturned
		default:
			log.Printf("[KAFKA-CONSUMER] Unknown status: %s, skipping database update\n", event.Status)
			continue
		}

		log.Printf("[KAFKA-CONSUMER] Updating Order AWB=%s to Status=%s (from Tracking status=%s)\n", event.AWB, targetStatus, event.Status)
		
		// Update status in PostgreSQL DB
		if err := c.repo.UpdateStatus(ctx, event.AWB, targetStatus); err != nil {
			log.Printf("[KAFKA-CONSUMER] Error updating order status in DB: %v\n", err)
		} else {
			log.Printf("[KAFKA-CONSUMER] Successfully updated Order AWB=%s to Status=%s\n", event.AWB, targetStatus)
		}
	}
}

// Close closes the underlying Kafka reader
func (c *kafkaConsumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}
