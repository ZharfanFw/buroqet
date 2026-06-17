package domain

import (
	"time"
	"github.com/google/uuid"
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

// Order is the main entity stored in PostgreSQL (Supabase)
type Order struct {
	OrderID              uuid.UUID   `gorm:"type:uuid;primaryKey;column:order_id" json:"order_id"`
	AWBNumber            string      `gorm:"column:awb_number" json:"awb_number"`
	CustomerID           *uuid.UUID  `gorm:"type:uuid;column:customer_id" json:"customer_id"`
	ServiceID            *uuid.UUID  `gorm:"type:uuid;column:service_id" json:"service_id"`
	SenderName           string      `gorm:"column:sender_name" json:"sender_name"`
	SenderPhone          string      `gorm:"column:sender_phone" json:"sender_phone"`
	SenderAddress        string      `gorm:"column:sender_address" json:"sender_address"`
	OriginPostalCode     string      `gorm:"column:origin_postal_code" json:"origin_postal_code"`
	OriginLat            float64     `gorm:"column:origin_lat" json:"origin_lat"`
	OriginLng            float64     `gorm:"column:origin_lng" json:"origin_lng"`
	ReceiverName         string      `gorm:"column:receiver_name" json:"receiver_name"`
	ReceiverPhone        string      `gorm:"column:receiver_phone" json:"receiver_phone"`
	ReceiverAddress      string      `gorm:"column:receiver_address" json:"receiver_address"`
	DestPostalCode       string      `gorm:"column:dest_postal_code" json:"dest_postal_code"`
	DestLat              float64     `gorm:"column:dest_lat" json:"dest_lat"`
	DestLng              float64     `gorm:"column:dest_lng" json:"dest_lng"`
	ActualWeightKg       float64     `gorm:"column:actual_weight_kg" json:"actual_weight_kg"`
	LengthCm             float64     `gorm:"column:length_cm" json:"length_cm"`
	WidthCm              float64     `gorm:"column:width_cm" json:"width_cm"`
	HeightCm             float64     `gorm:"column:height_cm" json:"height_cm"`
	VolumetricWeightKg   float64     `gorm:"column:volumetric_weight_kg" json:"volumetric_weight_kg"`
	PricingMethod        string      `gorm:"column:pricing_method" json:"pricing_method"`
	BaseTariff           float64     `gorm:"column:base_tariff" json:"base_tariff"`
	InsuranceFee         float64     `gorm:"column:insurance_fee" json:"insurance_fee"`
	Discount             float64     `gorm:"column:discount" json:"discount"`
	TotalCost            float64     `gorm:"column:total_cost" json:"total_cost"`
	UseInsurance         bool        `gorm:"column:use_insurance" json:"use_insurance"`
	PaymentType          PaymentType `gorm:"column:payment_type" json:"payment_type"`
	IsCod                bool        `gorm:"column:is_cod" json:"is_cod"`
	RouteID              *uuid.UUID  `gorm:"type:uuid;column:route_id" json:"route_id"`
	PaymentRef           string      `gorm:"column:payment_ref" json:"payment_ref"`
	Status               OrderStatus `gorm:"column:status" json:"status"`
	Notes                string      `gorm:"column:notes" json:"notes"`
	CreatedAt            time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt            time.Time   `gorm:"column:updated_at" json:"updated_at"`
}

// --- Request / Response DTOs ---

// CreateOrderRequest is the payload accepted from client
type CreateOrderRequest struct {
	CustomerID    string `json:"customer_id" binding:"omitempty,uuid"`
	ServiceID     string `json:"service_id" binding:"omitempty,uuid"`

	// Sender
	SenderName    string `json:"sender_name" binding:"required"`
	SenderPhone   string `json:"sender_phone" binding:"required"`
	SenderAddress string `json:"sender_address" binding:"required"`
	OriginPostal  string `json:"origin_postal" binding:"required"`
	OriginLat     float64 `json:"origin_lat"`
	OriginLng     float64 `json:"origin_lng"`

	// Receiver
	ReceiverName    string `json:"receiver_name" binding:"required"`
	ReceiverPhone   string `json:"receiver_phone" binding:"required"`
	ReceiverAddress string `json:"receiver_address" binding:"required"`
	DestPostal      string `json:"dest_postal" binding:"required"`
	DestLat         float64 `json:"dest_lat"`
	DestLng         float64 `json:"dest_lng"`

	// Package
	WeightActual float64 `json:"weight_actual" binding:"required,gt=0"`
	Length       float64 `json:"length" binding:"required,gt=0"`
	Width        float64 `json:"width" binding:"required,gt=0"`
	Height       float64 `json:"height" binding:"required,gt=0"`

	// Service
	ServiceType ServiceType `json:"service_type" binding:"required,oneof=REGULER EXPRESS"`
	PaymentType PaymentType `json:"payment_type" binding:"required,oneof=COD NON_COD"`
	UseInsurance bool       `json:"use_insurance"`
}

// CreateOrderResponse is returned after successful order creation
type CreateOrderResponse struct {
	OrderID       string      `json:"order_id"`
	AWBNumber     string      `json:"awb_number"`
	PaymentRef    string      `json:"payment_ref"`
	Status        OrderStatus `json:"status"`
	TotalCost     float64     `json:"total_cost"`
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
	BaseFare     float64 `json:"base_fare"`
	Insurance    float64 `json:"insurance"`
	Discount     float64 `json:"discount"`
	TotalPrice   float64 `json:"total_price"`
	EstimatedSLA string  `json:"estimated_sla"`
}
