<template>
  <div class="page-container">
    <div class="profile-page">
      <h2>My Profile</h2>

      <div v-if="profilePending" class="profile-content">
        <div class="loading">Loading...</div>
      </div>

      <div v-else-if="profileError" class="profile-content">
        <div class="empty-state">Failed to load profile</div>
      </div>

      <div v-else-if="user" class="profile-content">
        <div class="profile-info">
          <div class="profile-field">
            <span class="profile-label">Username:</span>
            <span class="profile-value">{{ user.username }}</span>
          </div>
          <div class="profile-field">
            <span class="profile-label">Email:</span>
            <span class="profile-value">{{ user.email }}</span>
          </div>
          <div class="profile-field">
            <span class="profile-label">User ID:</span>
            <span class="profile-value">{{ user.id }}</span>
          </div>
          <div class="profile-field">
            <span class="profile-label">Member Since:</span>
            <span class="profile-value">{{ formatDate(user.created_at) }}</span>
          </div>
        </div>
      </div>

      <h2>My Bookings</h2>

      <div v-if="bookingsPending" class="bookings-content">
        <div class="loading">Loading bookings...</div>
      </div>

      <div v-else-if="bookingsError" class="bookings-content">
        <div class="empty-state">Failed to load bookings</div>
      </div>

      <div v-else-if="bookings.length === 0" class="bookings-content">
        <div class="empty-state">
          No bookings yet. <NuxtLink to="/events">Browse events</NuxtLink>
        </div>
      </div>

      <div v-else class="bookings-content">
        <div class="bookings-grid">
          <div v-for="booking in bookings" :key="booking.ID" class="booking-card">
            <div class="booking-info">
              <h4>{{ booking.event_name || 'Event' }}</h4>
              <div class="booking-details">
                <div>Quantity: {{ booking.quantity }} tickets</div>
                <div>Total: ${{ booking.total_price.toFixed(2) }}</div>
                <div>Booked: {{ formatBookingDate(booking.CreatedAt) }}</div>
                <div>Status: {{ booking.status }}</div>
              </div>
            </div>
            <div class="booking-actions">
              <button
                v-if="booking.status === 'confirmed'"
                class="btn-small btn-danger"
                @click="handleCancelBooking(booking.ID)"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      </div>

      <button class="logout-btn" @click="handleLogout">Logout</button>
    </div>
  </div>
</template>

<script setup lang="ts">
const api = useApi()
const router = useRouter()

definePageMeta({
  middleware: 'auth',
})

const {
  data: profileData,
  pending: profilePending,
  error: profileError,
} = await useAsyncData('profile', () => api.getProfile())

const {
  data: bookingsData,
  pending: bookingsPending,
  error: bookingsError,
  refresh: refreshBookings,
} = await useAsyncData('bookings', () => api.getMyBookings())

const user = computed(() => profileData.value?.data)
const bookings = computed(() => bookingsData.value?.data || [])

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

const formatBookingDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

const handleCancelBooking = async (bookingId: number) => {
  if (!confirm('Are you sure you want to cancel this booking?')) {
    return
  }

  try {
    await api.cancelBooking(bookingId)
    alert('Booking cancelled successfully')
    refreshBookings()
  } catch (error: any) {
    alert(error.message || 'Failed to cancel booking')
  }
}

const handleLogout = () => {
  api.removeToken()
  router.push('/')
}
</script>