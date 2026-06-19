package domain

import "errors"

var (
	// ErrNoCourierAvailable dikembalikan ketika tidak ada kurir available dalam radius.
	ErrNoCourierAvailable = errors.New("no available courier found within the given radius")

	// ErrCourierNotFound dikembalikan ketika kurir dengan ID tertentu tidak ditemukan.
	ErrCourierNotFound = errors.New("courier not found")

	// ErrInvalidLocation dikembalikan ketika koordinat tidak valid.
	ErrInvalidLocation = errors.New("invalid location coordinates")
)