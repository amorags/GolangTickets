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
        <div :class="{ 'tickets-updated': justUpdated }">
          👥 {{ localAvailableTickets }}/{{ event.capacity }} tickets available
          <span v-if="justUpdated" class="update-badge">Updated</span>
        </div>
        <div>💰 ${{ event.price.toFixed(2) }}</div>
        <WebSocketStatus :isConnected="ws.isConnected.value" />
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
              :max="Math.min(localAvailableTickets, 10)"
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
import type { AvailabilityUpdate } from '~/types'

const route = useRoute()
const router = useRouter()
const api = useApi()
const ws = useWebSocket()

const eventId = computed(() => parseInt(route.params.id as string))
const isLoggedIn = ref(false)
const quantity = ref(1)
const bookingMessage = ref('')
const bookingMessageType = ref<'success' | 'error'>('success')
const localAvailableTickets = ref(0)
const justUpdated = ref(false)

onMounted(() => {
  isLoggedIn.value = !!api.getToken()

  // Initialize WebSocket connection
  const token = api.getToken()
  if (token) {
    ws.connect(token)
  }

  // Set initial availability
  if (event.value) {
    localAvailableTickets.value = event.value.available_tickets
  }

  // Listen for availability updates
  ws.on('availability_update', (update: AvailabilityUpdate) => {
    if (update.event_id === eventId.value) {
      localAvailableTickets.value = update.available_tickets
      justUpdated.value = true
      setTimeout(() => {
        justUpdated.value = false
      }, 3000)
    }
  })
})

// Subscribe when WebSocket is connected OR when event ID changes
watch([() => ws.isConnected.value, eventId], ([connected, id]) => {
  if (connected && id) {
    ws.subscribe(id)
  }
}, { immediate: true })

const { data: eventData, pending, error } = await useAsyncData(
  `event-${eventId.value}`,
  () => api.getEvent(eventId.value)
)

const event = computed(() => eventData.value?.data)

// Update local availability when event data loads
watch(event, (newEvent) => {
  if (newEvent) {
    localAvailableTickets.value = newEvent.available_tickets
  }
}, { immediate: true })

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

<style scoped>
.tickets-updated {
  position: relative;
  animation: highlight 0.5s ease;
}

@keyframes highlight {
  0% {
    background-color: transparent;
  }
  50% {
    background-color: #fef3c7;
  }
  100% {
    background-color: transparent;
  }
}

.update-badge {
  display: inline-block;
  margin-left: 8px;
  padding: 2px 8px;
  background-color: #10b981;
  color: white;
  font-size: 10px;
  border-radius: 12px;
  font-weight: 600;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.8);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}
</style>