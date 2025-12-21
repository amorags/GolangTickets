export interface User {
  id: number
  username: string
  email: string
  created_at: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface Event {
  ID: number
  name: string
  description: string
  event_type: string
  venue_name: string
  city: string
  address: string
  date: string
  price: number
  capacity: number
  available_tickets: number
  image_url: string
  CreatedAt: string
}

export interface Booking {
  ID: number
  event_id: number
  event_name: string
  user_id: number
  quantity: number
  total_price: number
  status: 'confirmed' | 'cancelled'
  CreatedAt: string
}

export interface ApiResponse<T> {
  message: string
  data: T
}

export interface ApiError {
  error: string
}