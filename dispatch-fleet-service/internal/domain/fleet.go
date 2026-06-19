package domain

import (
	"time"
)

// CourierStatus merepresentasikan status operasional kurir.
type CourierStatus string

const (
	StatusAvailable CourierStatus = "available"
	StatusAssigned  CourierStatus = "assigned"
	StatusOnDelivery CourierStatus = "on_delivery"
	StatusOffline   CourierStatus = "offline"
)

// Point merepresentasikan koordinat geospasial (longitude, latitude).
type Point struct {
	Longitude float64
	Latitude  float64
}

// Courier merepresentasikan entitas kurir dalam domain.
type Courier struct {
	ID              string
	Name            string
	CurrentLocation Point
	Status          CourierStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// AssignCourierResult merepresentasikan hasil dari proses assign kurir.
type AssignCourierResult struct {
	Courier         *Courier
	DistanceMeters  float64
}