package model

import (
	"time"
)

// OrderStatus represents the lifecycle status of an order
type OrderStatus string

const (
	StatusOrderCreated   OrderStatus = "ORDER_CREATED"
	StatusPaymentPending OrderStatus = "PAYMENT_PENDING"
	StatusPaymentPaid    OrderStatus = "PAYMENT_PAID"
	StatusCancelled      OrderStatus = "CANCELLED"
)

// ServiceType represents the shipping service type
type ServiceType string

const (
	ServiceRegular ServiceType = "REGULER"
	ServiceExpress ServiceType = "EXPRESS"
)

// PaymentType represents how the order is paid
type PaymentType string

const (
	PaymentCOD    PaymentType = "COD"
	PaymentNonCOD PaymentType = "NON_COD"
)

// Order is the main entity stored in PostgreSQL
type Order struct {
	ID            uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	AWBNumber     string      `gorm:"uniqueIndex;not null" json:"awb_number"`
	TransactionID string      `gorm:"uniqueIndex;not null" json:"transaction_id"`
	Status        OrderStatus `gorm:"not null" json:"status"`

	// Sender info
	SenderName    string `gorm:"not null" json:"sender_name"`
	SenderPhone   string `gorm:"not null" json:"sender_phone"`
	SenderAddress string `gorm:"not null" json:"sender_address"`
	OriginCity    string `gorm:"not null" json:"origin_city"`
	OriginPostal  string `gorm:"not null" json:"origin_postal"`

	// Receiver info
	ReceiverName    string `gorm:"not null" json:"receiver_name"`
	ReceiverPhone   string `gorm:"not null" json:"receiver_phone"`
	ReceiverAddress string `gorm:"not null" json:"receiver_address"`
	DestCity        string `gorm:"not null" json:"dest_city"`
	DestPostal      string `gorm:"not null" json:"dest_postal"`

	// Package dimensions
	WeightActual    float64 `gorm:"not null" json:"weight_actual"`
	WeightVolumetri float64 `gorm:"not null" json:"weight_volumetri"`
	Length          float64 `gorm:"not null" json:"length"`
	Width           float64 `gorm:"not null" json:"width"`
	Height          float64 `gorm:"not null" json:"height"`

	// Service & Payment
	ServiceType ServiceType `gorm:"not null" json:"service_type"`
	PaymentType PaymentType `gorm:"not null" json:"payment_type"`
	TotalPrice  float64     `gorm:"not null" json:"total_price"`
	PaymentURL  string      `json:"payment_url,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// --- Request / Response DTOs ---

// CreateOrderRequest is the payload accepted from client
type CreateOrderRequest struct {
	// Sender
	SenderName    string `json:"sender_name" binding:"required"`
	SenderPhone   string `json:"sender_phone" binding:"required"`
	SenderAddress string `json:"sender_address" binding:"required"`
	OriginCity    string `json:"origin_city" binding:"required"`
	OriginPostal  string `json:"origin_postal" binding:"required"`

	// Receiver
	ReceiverName    string `json:"receiver_name" binding:"required"`
	ReceiverPhone   string `json:"receiver_phone" binding:"required"`
	ReceiverAddress string `json:"receiver_address" binding:"required"`
	DestCity        string `json:"dest_city" binding:"required"`
	DestPostal      string `json:"dest_postal" binding:"required"`

	// Package
	WeightActual float64 `json:"weight_actual" binding:"required,gt=0"`
	Length       float64 `json:"length" binding:"required,gt=0"`
	Width        float64 `json:"width" binding:"required,gt=0"`
	Height       float64 `json:"height" binding:"required,gt=0"`

	// Service
	ServiceType ServiceType `json:"service_type" binding:"required,oneof=REGULER EXPRESS"`
	PaymentType PaymentType `json:"payment_type" binding:"required,oneof=COD NON_COD"`
}

// CreateOrderResponse is returned after successful order creation
type CreateOrderResponse struct {
	AWBNumber     string      `json:"awb_number"`
	TransactionID string      `json:"transaction_id"`
	Status        OrderStatus `json:"status"`
	TotalPrice    float64     `json:"total_price"`
	PaymentURL    string      `json:"payment_url,omitempty"`
}

// PricingRequest is sent to Pricing & Routing Service
type PricingRequest struct {
	OriginPostal string      `json:"origin_postal"`
	DestPostal   string      `json:"dest_postal"`
	Weight       float64     `json:"weight"`
	Length       float64     `json:"length"`
	Width        float64     `json:"width"`
	Height       float64     `json:"height"`
	ServiceType  ServiceType `json:"service_type"`
}

// PricingResponse is received from Pricing & Routing Service
type PricingResponse struct {
	BaseFare    float64 `json:"base_fare"`
	Insurance   float64 `json:"insurance"`
	Discount    float64 `json:"discount"`
	TotalPrice  float64 `json:"total_price"`
	EstimatedSLA string `json:"estimated_sla"`
}