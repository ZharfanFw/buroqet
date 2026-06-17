package kafka

import (
	"context"
	"log"
)

type kafkaProducer struct {
	broker string
}

func NewKafkaProducer(broker string) *kafkaProducer {
	return &kafkaProducer{broker: broker}
}

func (k *kafkaProducer) PublishEvent(ctx context.Context, topic string, key string, value []byte) error {
	// Untuk keperluan tugas, kita log saja dulu
	// Di production, ini akan benar-benar kirim ke Kafka
	log.Printf("[KAFKA] topic=%s key=%s value=%s", topic, key, string(value))
	return nil
}