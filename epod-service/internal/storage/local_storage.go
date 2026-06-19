package storage

import (
	"io"
	"os"
	"path/filepath"
)

// StorageInterface dibuat agar nanti kalau mau ganti ke S3/Supabase Storage tinggal ganti isinya saja
type StorageInterface interface {
	Upload(file io.Reader, fileName string) (string, error)
}

type LocalStorage struct {
	uploadDir string
}

// NewLocalStorage menginisialisasi storage lokal dengan target folder ./uploads
func NewLocalStorage() StorageInterface {
	uploadDir := "./uploads"
	
	// Otomatis bikin folder ./uploads kalau belum ada di komputer/server
	_ = os.MkdirAll(uploadDir, os.ModePerm)

	return &LocalStorage{
		uploadDir: uploadDir,
	}
}

// Upload menerima file (io.Reader) dan memproses penyimpanannya ke harddisk
func (s *LocalStorage) Upload(file io.Reader, fileName string) (string, error) {
	// Tentukan path lengkap di server, misal: uploads/bukti-resi.jpg
	filePath := filepath.Join(s.uploadDir, fileName)

	// Buat file baru di harddisk server
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Salin (copy) data file dari request (frontend) ke file yang baru dibuat di server
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	// Kembalikan URL path relatif yang akan disimpan di database
	// Contoh: "/uploads/bukti-resi.jpg"
	// Sesuai dengan routing di main.go: router.Static("/uploads", "./uploads")
	return "/uploads/" + fileName, nil
}