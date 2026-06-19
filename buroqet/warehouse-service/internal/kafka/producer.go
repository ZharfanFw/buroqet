package kafka

import (
	"context"
	"log"

	kafkalib "github.com/segmentio/kafka-go"
)

type kafkaProducer struct {
	broker string
}

func NewKafkaProducer(broker string) *kafkaProducer {
	return &kafkaProducer{broker: broker}
}

func (k *kafkaProducer) PublishEvent(ctx context.Context, topic string, key string, value []byte) error {
	w := &kafkalib.Writer{
		Addr:     kafkalib.TCP(k.broker),
		Topic:    topic,
		Balancer: &kafkalib.LeastBytes{},
	}
	defer w.Close()

	err := w.WriteMessages(ctx,
		kafkalib.Message{
			Key:   []byte(key),
			Value: value,
		},
	)

	if err != nil {
		log.Printf("[KAFKA ERROR] Failed to publish topic=%s key=%s err=%v", topic, key, err)
		return err
	}

	log.Printf("[KAFKA] PUBLISHED topic=%s key=%s value=%s", topic, key, string(value))
	return nil
}