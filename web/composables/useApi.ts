import type { AuthResponse, User, Event, Booking, ApiResponse, EventFilters, PaginatedEventsResponse } from '~/types'

const API_URL = 'http://localhost:8080'

export const useApi = () => {
  const getToken = () => {
    if (import.meta.client) {
      return localStorage.getItem('token')
    }
    return null
  }

  const setToken = (token: string) => {
    if (import.meta.client) {
      localStorage.setItem('token', token)
    }
  }

  const removeToken = () => {
    if (import.meta.client) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }
  }

  const fetchWithAuth = async <T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> => {
    const token = getToken()
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    }

    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }

    const response = await fetch(`${API_URL}${endpoint}`, {
      ...options,
      headers,
    })

    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.error || 'Request failed')
    }

    return response.json()
  }

  // Auth API
  const signup = async (username: string, email: string, password: string) => {
    const response = await fetchWithAuth<AuthResponse>('/auth/signup', {
      method: 'POST',
      body: JSON.stringify({ username, email, password }),
    })
    if (response.data.token) {
      setToken(response.data.token)
      if (import.meta.client) {
        localStorage.setItem('user', JSON.stringify(response.data.user))
      }
    }
    return response
  }

  const login = async (email: string, password: string) => {
    const response = await fetchWithAuth<AuthResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
    if (response.data.token) {
      setToken(response.data.token)
      if (import.meta.client) {
        localStorage.setItem('user', JSON.stringify(response.data.user))
      }
    }
    return response
  }

  const getProfile = () => fetchWithAuth<User>('/profile')

  // Events API
  const getEvents = (filters?: EventFilters) => {
    let endpoint = '/events'

    if (filters) {
      const params = new URLSearchParams()

      if (filters.search) params.append('search', filters.search)
      if (filters.type) params.append('type', filters.type)
      if (filters.city) params.append('city', filters.city)
      if (filters.date_from) params.append('date_from', filters.date_from)
      if (filters.date_to) params.append('date_to', filters.date_to)
      if (filters.price_min !== undefined) params.append('price_min', filters.price_min.toString())
      if (filters.price_max !== undefined) params.append('price_max', filters.price_max.toString())
      if (filters.status) params.append('status', filters.status)
      if (filters.page) params.append('page', filters.page.toString())
      if (filters.limit) params.append('limit', filters.limit.toString())
      if (filters.sort) params.append('sort', filters.sort)
      if (filters.order) params.append('order', filters.order)

      const queryString = params.toString()
      if (queryString) {
        endpoint += `?${queryString}`
      }
    }

    return fetchWithAuth<PaginatedEventsResponse>(endpoint)
  }

  const getEvent = (id: number) => fetchWithAuth<Event>(`/events/${id}`)

  const createEvent = (eventData: any) => 
    fetchWithAuth<Event>('/events', {
      method: 'POST',
      body: JSON.stringify(eventData),
    })

  // Bookings API
  const createBooking = (event_id: number, quantity: number) =>
    fetchWithAuth<Booking>('/bookings', {
      method: 'POST',
      body: JSON.stringify({ event_id, quantity }),
    })

  const getMyBookings = () => fetchWithAuth<Booking[]>('/bookings')

  const cancelBooking = (id: number) =>
    fetchWithAuth<void>(`/bookings/${id}`, {
      method: 'DELETE',
    })

  return {
    signup,
    login,
    getProfile,
    getEvents,
    getEvent,
    createEvent,
    createBooking,
    getMyBookings,
    cancelBooking,
    getToken,
    removeToken,
  }
}