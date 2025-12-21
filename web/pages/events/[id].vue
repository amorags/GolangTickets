<template>
  <div class="page-container">
    <div v-if="pending" class="event-details-container">
      <div class="loading">Loading event details...</div>
    </div>

    <div v-else-if="error" class="event-details-container">
      <div class="empty-state">Event not found</div>
    </div>

    <div v-else-if="event" class="event-details-container">
      <div class="event-header">
        <div
          class="event-detail-image"
          :style="
            event.image_url
              ? `background-image: url(${event.image_url}); background-size: cover;`
              : ''
          "
        ></div>
        <h1 class="event-title">{{ event.name }}</h1>
        <div class="event-type">{{ event.event_type || 'event' }}</div>
      </div>

      <div class="event-meta">
        <div>📅 {{ formatDate(event.date) }}</div>
        <div>📍 {{ event.venue_name }}</div>
        <div>🏙️ {{ event.city }}, {{ event.address }}</div>
        <div>
          👥 {{ event.available_tickets }}/{{ event.capacity }} tickets available
        </div>
        <div>💰 ${{ event.price.toFixed(2) }}</div>
      </div>

      <div class="event-description">
        <h3>About this event</h3>
        <p>{{ event.description || 'No description available' }}</p>
      </div>

      <div class="booking-section">
        <h3>Book Tickets</h3>
        <div v-if="bookingMessage" :class="['message', 'show', bookingMessageType]">
          {{ bookingMessage }}
        </div>
        <form v-if="isLoggedIn" class="booking-form" @submit.prevent="handleBooking">
          <div class="form-group">
            <label for="quantity">Number of Tickets</label>
            <input
              id="quantity"
              v-model.number="quantity"
              type="number"
              min="1"
              :max="event.available_tickets"
              required
            />
          </div>
          <button type="submit" class="btn-primary">
            Book Now - ${{ event.price.toFixed(2) }} each
          </button>
        </form>
        <p v-else>
          Please <NuxtLink to="/">login</NuxtLink> to book tickets
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const api = useApi()

const eventId = computed(() => parseInt(route.params.id as string))
const isLoggedIn = ref(false)
const quantity = ref(1)
const bookingMessage = ref('')
const bookingMessageType = ref<'success' | 'error'>('success')

onMounted(() => {
  isLoggedIn.value = !!api.getToken()
})

const { data: eventData, pending, error } = await useAsyncData(
  `event-${eventId.value}`,
  () => api.getEvent(eventId.value)
)

const event = computed(() => eventData.value?.data)

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    weekday: 'long',
    month: 'long',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const handleBooking = async () => {
  try {
    await api.createBooking(eventId.value, quantity.value)
    bookingMessage.value = 'Booking successful! Redirecting to your profile...'
    bookingMessageType.value = 'success'
    setTimeout(() => {
      router.push('/profile')
    }, 2000)
  } catch (error: any) {
    bookingMessage.value = error.message || 'Booking failed'
    bookingMessageType.value = 'error'
  }
}
</script>