package repository

import (
	"database/sql"
	"errors"
	"time"

	// Import driver PostgreSQL
	_ "github.com/lib/pq"
)

// EPOD merepresentasikan struktur data di dalam database
type EPOD struct {
	ID           string     `json:"id"`
	AWB          string     `json:"awb"`
	Receiver     string     `json:"receiver"`
	PhotoURL     string     `json:"photo_url"`
	SignatureURL string     `json:"signature_url,omitempty"`
	Lat          float64    `json:"lat"`
	Lng          float64    `json:"lng"`
	Status       string     `json:"status"`
	RejectReason string     `json:"reject_reason,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	VerifiedAt   *time.Time `json:"verified_at,omitempty"`
}

// EPODRepository adalah interface/kontrak untuk layer repository
type EPODRepository interface {
	Create(epod *EPOD) error
	FindAll(status string, limit int) ([]EPOD, error)
	FindByID(id string) (*EPOD, error)
	FindByAWB(awb string) (*EPOD, error)
	UpdateStatus(id string, status string) error
}

// postgresEPODRepository adalah implementasi dari interface menggunakan PostgreSQL
type postgresEPODRepository struct {
	db *sql.DB
}

// NewPostgresEPODRepository adalah constructor yang dipanggil di main.go
func NewPostgresEPODRepository(dsn string) EPODRepository {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic("Gagal koneksi ke database: " + err.Error())
	}

	// Test ping ke database
	if err := db.Ping(); err != nil {
		panic("Database tidak merespons: " + err.Error())
	}

	return &postgresEPODRepository{db: db}
}

// Create menyimpan data ePOD baru ke database
func (r *postgresEPODRepository) Create(epod *EPOD) error {
	query := `
		INSERT INTO epod.epods (id, awb, receiver, photo_url, signature_url, lat, lng, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(query,
		epod.ID, epod.AWB, epod.Receiver, epod.PhotoURL, epod.SignatureURL,
		epod.Lat, epod.Lng, epod.Status, epod.CreatedAt,
	)
	return err
}

// FindAll mengambil semua data ePOD, opsional bisa difilter berdasarkan status
func (r *postgresEPODRepository) FindAll(status string, limit int) ([]EPOD, error) {
	var rows *sql.Rows
	var err error

	if status != "" {
		query := `SELECT id, awb, receiver, photo_url, signature_url, lat, lng, status, created_at, verified_at 
				  FROM epod.epods WHERE status = $1 ORDER BY created_at DESC LIMIT $2`
		rows, err = r.db.Query(query, status, limit)
	} else {
		query := `SELECT id, awb, receiver, photo_url, signature_url, lat, lng, status, created_at, verified_at 
				  FROM epod.epods ORDER BY created_at DESC LIMIT $1`
		rows, err = r.db.Query(query, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var epods []EPOD
	for rows.Next() {
		var epod EPOD
		if err := rows.Scan(
			&epod.ID, &epod.AWB, &epod.Receiver, &epod.PhotoURL, &epod.SignatureURL,
			&epod.Lat, &epod.Lng, &epod.Status, &epod.CreatedAt, &epod.VerifiedAt,
		); err != nil {
			return nil, err
		}
		epods = append(epods, epod)
	}
	return epods, nil
}

// FindByID mencari ePOD berdasarkan ID
func (r *postgresEPODRepository) FindByID(id string) (*EPOD, error) {
	query := `SELECT id, awb, receiver, photo_url, signature_url, lat, lng, status, created_at, verified_at 
			  FROM epod.epods WHERE id = $1`
	
	var epod EPOD
	err := r.db.QueryRow(query, id).Scan(
		&epod.ID, &epod.AWB, &epod.Receiver, &epod.PhotoURL, &epod.SignatureURL,
		&epod.Lat, &epod.Lng, &epod.Status, &epod.CreatedAt, &epod.VerifiedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("epod tidak ditemukan")
	} else if err != nil {
		return nil, err
	}
	return &epod, nil
}

// FindByAWB mencari ePOD berdasarkan nomor resi (AWB)
func (r *postgresEPODRepository) FindByAWB(awb string) (*EPOD, error) {
	query := `SELECT id, awb, receiver, photo_url, signature_url, lat, lng, status, created_at, verified_at 
			  FROM epod.epods WHERE awb = $1`
	
	var epod EPOD
	err := r.db.QueryRow(query, awb).Scan(
		&epod.ID, &epod.AWB, &epod.Receiver, &epod.PhotoURL, &epod.SignatureURL,
		&epod.Lat, &epod.Lng, &epod.Status, &epod.CreatedAt, &epod.VerifiedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("epod dengan awb tersebut tidak ditemukan")
	} else if err != nil {
		return nil, err
	}
	return &epod, nil
}

// UpdateStatus mengubah status (misal: dari PENDING_REVIEW ke VERIFIED)
func (r *postgresEPODRepository) UpdateStatus(id string, status string) error {
	var query string
	var err error

	// Jika statusnya VERIFIED, kita catat jam verifikasinya
	if status == "VERIFIED" {
		query = `UPDATE epod.epods SET status = $1, verified_at = $2 WHERE id = $3`
		_, err = r.db.Exec(query, status, time.Now(), id)
	} else {
		query = `UPDATE epod.epods SET status = $1 WHERE id = $2`
		_, err = r.db.Exec(query, status, id)
	}

	return err
}