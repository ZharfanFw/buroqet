// internal/cache/redis_cache.go
// Implementasi konkret TrackingCache menggunakan Redis.
// Menyimpan status terakhir paket untuk query cepat (O(1)).

package cache

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"tracking-service/internal/domain"

	"github.com/redis/go-redis/v9"
)

const (
	// TTL status di Redis: 7 hari
	statusTTL = 7 * 24 * time.Hour

	// Prefix key di Redis untuk namespace isolation
	keyPrefix = "tracking:status:"
)

// RedisTrackingCache adalah implementasi konkret dari domain.TrackingCache
type RedisTrackingCache struct {
	client *redis.Client
}

// NewRedisTrackingCache membuat instance cache baru
func NewRedisTrackingCache(client *redis.Client) *RedisTrackingCache {
	return &RedisTrackingCache{client: client}
}

// SetStatus menyimpan status terakhir paket ke Redis.
// Key format: "tracking:status:{awb}"
func (c *RedisTrackingCache) SetStatus(ctx context.Context, status *domain.TrackingStatus) error {
	key := keyPrefix + status.AWB

	data, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("gagal marshal status: %w", err)
	}

	if err := c.client.Set(ctx, key, data, statusTTL).Err(); err != nil {
		return fmt.Errorf("redis Set gagal untuk AWB %s: %w", status.AWB, err)
	}
	return nil
}

// GetStatus mengambil status terakhir paket dari Redis.
// Mengembalikan nil, nil jika key tidak ditemukan (cache miss).
func (c *RedisTrackingCache) GetStatus(ctx context.Context, awb string) (*domain.TrackingStatus, error) {
	key := keyPrefix + awb

	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // cache miss — bukan error
	}
	if err != nil {
		return nil, fmt.Errorf("redis Get gagal untuk AWB %s: %w", awb, err)
	}

	var status domain.TrackingStatus
	if err := json.Unmarshal(data, &status); err != nil {
		return nil, fmt.Errorf("gagal unmarshal status dari Redis: %w", err)
	}
	return &status, nil
}

// DeleteStatus menghapus status dari Redis (saat paket delivered/returned).
// Menghapus key yang tidak ada tidak dianggap error.
func (c *RedisTrackingCache) DeleteStatus(ctx context.Context, awb string) error {
	key := keyPrefix + awb
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("redis Del gagal untuk AWB %s: %w", awb, err)
	}
	return nil
}

// ============================================================
// Helper untuk koneksi Redis — digunakan di main.go
// ============================================================

// ConnectRedis membuat koneksi ke Redis dengan connection pool optimal untuk high load
func ConnectRedis(addr string, password string, db int) (*redis.Client, error) {
	var opt *redis.Options

	// Cek apakah addr berupa URL (mendukung rediss:// untuk Upstash)
	if len(addr) > 8 && (addr[:8] == "redis://" || addr[:9] == "rediss://") {
		var err error
		opt, err = redis.ParseURL(addr)
		if err != nil {
			return nil, fmt.Errorf("gagal parse redis URL: %w", err)
		}
	} else {
		opt = &redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}
	}

	// Enable TLS if domain belongs to Upstash but wasn't a URL
	if strings.Contains(addr, "upstash.io") {
		opt.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	// Connection pool untuk high load (5K-10K RPS)
	opt.PoolSize = 50
	opt.MinIdleConns = 10
	opt.DialTimeout = 5 * time.Second
	opt.ReadTimeout = 3 * time.Second
	opt.WriteTimeout = 3 * time.Second
	opt.PoolTimeout = 4 * time.Second
	opt.ConnMaxIdleTime = 5 * time.Minute

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("gagal ping Redis: %w", err)
	}

	return client, nil
}