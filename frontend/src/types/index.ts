// Common types shared across services

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
}

export interface ApiError {
  message: string;
  code?: string;
}

// --- Auth ---
export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token?: string;
  user: User;
}

export interface User {
  id: string;
  name: string;
  email: string;
  role: 'admin' | 'courier' | 'warehouse' | 'customer';
}

// --- Tracking ---
export interface TrackingEvent {
  id: string;
  awb: string;
  status: string;
  location: string;
  hub_id: string;
  description: string;
  latitude: number;
  longitude: number;
  timestamp: string;
  created_at: string;
  source: string;
}

export interface TrackingStatus {
  awb: string;
  current_status: string;
  last_location: string;
  last_updated: string;
}

export interface TrackingHistory {
  awb: string;
  events: TrackingEvent[];
  total: number;
}

export interface TrackingInfo {
  awb: string;
  status: string;
  origin: string;
  destination: string;
  estimatedDelivery: string;
  events: TrackingEvent[];
}

// --- Order ---
export type OrderStatus = 'pending' | 'processing' | 'shipped' | 'delivered' | 'cancelled';

export interface Order {
  id: string;
  awb: string;
  senderId: string;
  recipientName: string;
  recipientAddress: string;
  recipientPhone: string;
  weight: number;
  price: number;
  status: OrderStatus;
  createdAt: string;
  updatedAt: string;
}

// --- Dispatch & Fleet ---
export interface Vehicle {
  id: string;
  plateNumber: string;
  type: string;
  status: 'available' | 'on-duty' | 'maintenance';
  driverId?: string;
}

export interface DispatchAssignment {
  id: string;
  orderId: string;
  vehicleId: string;
  driverId: string;
  assignedAt: string;
  status: string;
}

// --- Pricing ---
export interface PriceCalculationRequest {
  origin: string;
  destination: string;
  weight: number;
  serviceType: 'regular' | 'express' | 'same-day';
}

export interface PriceCalculationResponse {
  price: number;
  estimatedDays: number;
  serviceType: string;
}

// --- ePOD ---
export interface EPod {
  id: string;
  orderId: string;
  imageUrl: string;
  signature?: string;
  receiverName: string;
  deliveredAt: string;
}

// --- Warehouse ---
export interface InventoryItem {
  id: string;
  packageId: string;
  location: string;
  status: 'inbound' | 'stored' | 'outbound';
  arrivedAt: string;
}
