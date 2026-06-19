// internal/cache/redis_cache_test.go
// Functional test untuk RedisTrackingCache.
// Test ini MENGAKSES Redis nyata — dijalankan dengan build tag functional:
//   go test -v -tags=functional -count=1 ./internal/cache/...
//
// Tujuan: memverifikasi bahwa implementasi Redis (SetStatus, GetStatus, DeleteStatus)
// bekerja benar, termasuk TTL dan serialisasi JSON.

//go:build functional
// +build functional

package cache_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"tracking-service/internal/cache"
	"tracking-service/internal/domain"
)

// =========================================================
// RedisCacheSuite — Test suite untuk Redis cache
// =========================================================

type RedisCacheSuite struct {
	suite.Suite
	cache *cache.RedisTrackingCache
	ctx   context.Context
	seq   int
}

func (s *RedisCacheSuite) SetupSuite() {
	s.ctx = context.Background()

	redisAddr := os.Getenv("TEST_REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6380"
	}

	client, err := cache.ConnectRedis(redisAddr, "", 1) // DB index 1 untuk test
	s.Require().NoError(err, "Gagal koneksi ke Redis test")

	s.cache = cache.NewRedisTrackingCache(client)
}

func (s *RedisCacheSuite) SetupTest() {
	s.seq++
}

// generateAWB membuat AWB unik per test case
func (s *RedisCacheSuite) generateAWB(prefix string) string {
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixNano(), s.seq)
}

// buildStatus adalah helper untuk membuat TrackingStatus test
func buildStatus(awb, currentStatus, location string) *domain.TrackingStatus {
	return &domain.TrackingStatus{
		AWB:           awb,
		CurrentStatus: currentStatus,
		LastLocation:  location,
		LastUpdated:   time.Now(),
	}
}

// =========================================================
// CACHE-01: SetStatus + GetStatus — round-trip sukses
// =========================================================

func (s *RedisCacheSuite) TestSetAndGet_Sukses() {
	awb := s.generateAWB("CACHE-SET")
	status := buildStatus(awb, domain.StatusOnTransit, "Hub Surabaya")

	// Set ke Redis
	err := s.cache.SetStatus(s.ctx, status)
	// Akan FAIL karena "not implemented"
	s.NoError(err, "SetStatus harus sukses")

	// Get dari Redis
	result, err := s.cache.GetStatus(s.ctx, awb)
	s.NoError(err, "GetStatus harus sukses")
	s.NotNil(result, "Hasil tidak boleh nil")

	// Verifikasi semua field tersimpan dengan benar
	s.Equal(awb, result.AWB,                           "AWB harus sama")
	s.Equal(domain.StatusOnTransit, result.CurrentStatus, "Status harus sama")
	s.Equal("Hub Surabaya", result.LastLocation,       "Location harus sama")
}

// =========================================================
// CACHE-02: GetStatus — cache miss → return nil, nil
// =========================================================

func (s *RedisCacheSuite) TestGetStatus_CacheMiss() {
	awb := s.generateAWB("CACHE-MISS")

	result, err := s.cache.GetStatus(s.ctx, awb)

	// Cache miss bukan error — return nil, nil
	s.NoError(err, "Cache miss bukan error")
	s.Nil(result, "Cache miss harus return nil")
}

// =========================================================
// CACHE-03: SetStatus overwrite — update status yang sudah ada
// =========================================================

func (s *RedisCacheSuite) TestSetStatus_Overwrite() {
	awb := s.generateAWB("CACHE-OVERWRITE")

	// Set status pertama
	status1 := buildStatus(awb, domain.StatusInbound, "Hub Jakarta")
	err := s.cache.SetStatus(s.ctx, status1)
	s.Require().NoError(err)

	// Overwrite dengan status baru
	status2 := buildStatus(awb, domain.StatusDelivered, "Rumah Penerima")
	err = s.cache.SetStatus(s.ctx, status2)
	s.NoError(err, "Overwrite status harus sukses")

	// Get — harus mengembalikan status terbaru
	result, err := s.cache.GetStatus(s.ctx, awb)
	s.NoError(err)
	s.Equal(domain.StatusDelivered, result.CurrentStatus,
		"Status harus sudah terupdate ke DELIVERED")
	s.Equal("Rumah Penerima", result.LastLocation)
}

// =========================================================
// CACHE-04: DeleteStatus — hapus entry dari Redis
// =========================================================

func (s *RedisCacheSuite) TestDeleteStatus_Sukses() {
	awb := s.generateAWB("CACHE-DELETE")

	// Set dulu
	err := s.cache.SetStatus(s.ctx, buildStatus(awb, domain.StatusAtHub, "Hub Bandung"))
	s.Require().NoError(err)

	// Pastikan tersimpan
	result, _ := s.cache.GetStatus(s.ctx, awb)
	s.Require().NotNil(result, "Status harus ada sebelum dihapus")

	// Hapus
	err = s.cache.DeleteStatus(s.ctx, awb)
	s.NoError(err, "DeleteStatus harus sukses")

	// Verifikasi sudah terhapus — harus cache miss
	result, err = s.cache.GetStatus(s.ctx, awb)
	s.NoError(err)
	s.Nil(result, "Setelah dihapus, GetStatus harus return nil")
}

// =========================================================
// CACHE-05: DeleteStatus — hapus key yang tidak ada → tidak error
// =========================================================

func (s *RedisCacheSuite) TestDeleteStatus_KeyTidakAda() {
	awb := s.generateAWB("CACHE-DEL-NOTEXIST")

	// Delete key yang tidak ada seharusnya tidak error
	err := s.cache.DeleteStatus(s.ctx, awb)
	s.NoError(err, "Delete key yang tidak ada harus tidak error")
}

// =========================================================
// CACHE-06: Isolasi antar AWB — key Redis tidak saling menimpa
// =========================================================

func (s *RedisCacheSuite) TestIsolasiAntarAWB() {
	awbA := s.generateAWB("CACHE-AWB-A")
	awbB := s.generateAWB("CACHE-AWB-B")

	// Set status berbeda untuk dua AWB
	s.cache.SetStatus(s.ctx, buildStatus(awbA, domain.StatusInbound, "Hub A"))
	s.cache.SetStatus(s.ctx, buildStatus(awbB, domain.StatusDelivered, "Lokasi B"))

	// Get masing-masing — harus independen
	statusA, _ := s.cache.GetStatus(s.ctx, awbA)
	statusB, _ := s.cache.GetStatus(s.ctx, awbB)

	s.Require().NotNil(statusA)
	s.Require().NotNil(statusB)

	s.Equal(domain.StatusInbound, statusA.CurrentStatus,
		"Status AWB-A harus INBOUND")
	s.Equal(domain.StatusDelivered, statusB.CurrentStatus,
		"Status AWB-B harus DELIVERED")
	s.NotEqual(statusA.CurrentStatus, statusB.CurrentStatus,
		"Dua AWB tidak boleh saling menimpa")
}

// =========================================================
// CACHE-07: Serialisasi LastUpdated — timestamp tersimpan dengan benar
// =========================================================

func (s *RedisCacheSuite) TestSerialisasiTimestamp() {
	awb := s.generateAWB("CACHE-TIME")
	now := time.Now().Truncate(time.Millisecond) // Truncate untuk menghindari sub-millisecond diff

	status := &domain.TrackingStatus{
		AWB:           awb,
		CurrentStatus: domain.StatusAtHub,
		LastLocation:  "Hub Test",
		LastUpdated:   now,
	}

	s.cache.SetStatus(s.ctx, status)
	result, err := s.cache.GetStatus(s.ctx, awb)

	s.NoError(err)
	s.Require().NotNil(result)

	// Timestamp harus tersimpan dan ter-deserialize dengan benar
	s.False(result.LastUpdated.IsZero(), "LastUpdated tidak boleh zero")
	diff := result.LastUpdated.Sub(now)
	if diff < 0 {
		diff = -diff
	}
	s.Less(diff, time.Second, "LastUpdated harus sama dalam toleransi 1 detik")
}

// =========================================================
// CACHE-08: ConnectRedis — koneksi ke Redis berhasil
// =========================================================

func (s *RedisCacheSuite) TestConnectRedis_Sukses() {
	// Jika SetupSuite tidak fail, berarti koneksi berhasil
	// Test ini mengkonfirmasi koneksi awal berjalan baik
	s.NotNil(s.cache, "Redis cache instance harus ter-inisialisasi")
}

// =========================================================
// Entry point
// =========================================================

func TestRedisCacheSuite(t *testing.T) {
	suite.Run(t, new(RedisCacheSuite))
}