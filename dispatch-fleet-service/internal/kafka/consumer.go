package kafka

import (
	"context"
	"encoding/json"
	"log"

	"dispatch-fleet/internal/domain"

	"github.com/segmentio/kafka-go"
)

type OrderCreatedPayload struct {
	AWBNumber     string `json:"awb_number"`
	TransactionID string `json:"transaction_id"`
	OriginCity    string `json:"origin_city"`
}

type Consumer struct {
	reader          *kafka.Reader
	dispatchService domain.DispatchService
}

func NewConsumer(brokers []string, dispatchSvc domain.DispatchService) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   "order.created",
		GroupID: "dispatch-fleet-group",
	})
	return &Consumer{
		reader:          r,
		dispatchService: dispatchSvc,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("Kafka Consumer started listening to order.created")
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("error reading message: %v", err)
			continue
		}

		var payload OrderCreatedPayload
		if err := json.Unmarshal(m.Value, &payload); err != nil {
			log.Printf("error unmarshaling payload: %v", err)
			continue
		}

		log.Printf("Received order.created event for AWB: %s, Origin: %s", payload.AWBNumber, payload.OriginCity)

		// Mock Geocoding berdasarkan kota (OriginCity) -> Latitude, Longitude
		lat, lon := mockGeocode(payload.OriginCity)
		pickupLoc := domain.Point{Latitude: lat, Longitude: lon}

		// Radius 10km (10000m)
		result, err := c.dispatchService.AssignCourierToPickup(context.Background(), pickupLoc, 10000)
		if err != nil {
			log.Printf("Failed to assign courier for AWB %s: %v", payload.AWBNumber, err)
		} else {
			log.Printf("Successfully assigned Courier %s for AWB %s (distance: %.2fm)", result.Courier.Name, payload.AWBNumber, result.DistanceMeters)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}

func mockGeocode(city string) (float64, float64) {
	switch city {
	case "Jakarta":
		return -6.2088, 106.8456
	case "Bandung":
		return -6.9175, 107.6191
	case "Surabaya":
		return -7.2504, 112.7688
	default:
		// Default to Jakarta
		return -6.2088, 106.8456
	}
}
