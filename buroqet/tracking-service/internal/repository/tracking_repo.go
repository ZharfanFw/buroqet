// internal/repository/tracking_repo.go
// Implementasi konkret TrackingRepository menggunakan MongoDB.
// Layer ini TIDAK di-test dengan unit test (karena butuh DB nyata),
// melainkan di-test dengan functional test.

package repository

import (
	"context"
	"fmt"
	"time"

	"tracking-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionTrackingEvents = "tracking_events"
)

// MongoTrackingRepository adalah implementasi konkret dari domain.TrackingRepository
type MongoTrackingRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewMongoTrackingRepository membuat instance repository baru
func NewMongoTrackingRepository(db *mongo.Database) *MongoTrackingRepository {
	collection := db.Collection(collectionTrackingEvents)
	return &MongoTrackingRepository{
		db:         db,
		collection: collection,
	}
}

// EnsureIndexes membuat index MongoDB yang diperlukan untuk performa optimal.
// Dipanggil sekali saat startup aplikasi.
func (r *MongoTrackingRepository) EnsureIndexes(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "awb", Value: 1},
			{Key: "timestamp", Value: -1},
		},
		Options: options.Index().SetName("idx_awb_timestamp"),
	}
	_, err := r.collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("gagal membuat index MongoDB: %w", err)
	}
	return nil
}

// InsertEvent menyimpan satu event tracking ke MongoDB.
// Operasi ini append-only — event tidak pernah di-update atau di-delete.
func (r *MongoTrackingRepository) InsertEvent(ctx context.Context, event *domain.TrackingEvent) error {
	_, err := r.collection.InsertOne(ctx, event)
	if err != nil {
		return fmt.Errorf("InsertEvent gagal: %w", err)
	}
	return nil
}

// GetEventsByAWB mengambil semua events untuk satu AWB,
// diurutkan berdasarkan timestamp ascending (kronologis).
func (r *MongoTrackingRepository) GetEventsByAWB(ctx context.Context, awb string) ([]domain.TrackingEvent, error) {
	filter := bson.M{"awb": awb}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("GetEventsByAWB gagal: %w", err)
	}
	defer cursor.Close(ctx)

	var events []domain.TrackingEvent
	if err := cursor.All(ctx, &events); err != nil {
		return nil, fmt.Errorf("gagal decode events: %w", err)
	}

	// Kembalikan slice kosong (bukan nil) jika tidak ada event
	if events == nil {
		events = []domain.TrackingEvent{}
	}
	return events, nil
}

// GetLatestEventByAWB mengambil event terbaru (paling akhir) untuk satu AWB.
// Digunakan sebagai fallback ketika Redis cache miss.
func (r *MongoTrackingRepository) GetLatestEventByAWB(ctx context.Context, awb string) (*domain.TrackingEvent, error) {
	filter := bson.M{"awb": awb}
	opts := options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}})

	var event domain.TrackingEvent
	err := r.collection.FindOne(ctx, filter, opts).Decode(&event)
	if err == mongo.ErrNoDocuments {
		return nil, nil // AWB tidak ada — bukan error
	}
	if err != nil {
		return nil, fmt.Errorf("GetLatestEventByAWB gagal: %w", err)
	}
	return &event, nil
}

// ============================================================
// Helper untuk koneksi MongoDB — digunakan di main.go
// ============================================================

// ConnectMongoDB membuat koneksi ke MongoDB dengan connection pool optimal
func ConnectMongoDB(uri string, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(5 * time.Second).
		SetMaxPoolSize(100).
		SetMinPoolSize(10)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("gagal koneksi ke MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("gagal ping MongoDB: %w", err)
	}

	return client.Database(dbName), nil
}