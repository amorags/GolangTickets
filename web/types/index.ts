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
  user_id: number
  quantity: number
  total_price: number
  status: 'confirmed' | 'cancelled'
  CreatedAt: string
  event?: Event // Added optional event since it is preloaded
}

export interface ApiResponse<T> {
  message: string
  data: T
}

export interface ApiError {
  error: string
}

export interface CreateEventRequest {
  name: string
  description: string
  event_type: string
  venue_name: string
  city: string
  address: string
  date: string // ISO 8601 string
  price: number
  capacity: number
  image_url: string
}

export interface WebSocketMessage {
  type: string
  event_id?: number
  timestamp?: string
  data?: any
}

export interface AvailabilityUpdate {
  event_id: number
  available_tickets: number
  capacity: number
  last_updated: string
}

export interface EventFilters {
  search?: string
  type?: string
  city?: string
  date_from?: string
  date_to?: string
  price_min?: number
  price_max?: number
  status?: string
  page?: number
  limit?: number
  sort?: string
  order?: 'asc' | 'desc'
}

export interface PaginatedEventsResponse {
  events: Event[]
  total: number
  page: number
  limit: number
  total_pages: number
  has_next: boolean
  has_previous: boolean
}