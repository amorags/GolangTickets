<template>
  <div class="page-container">
    <div class="page-header">
      <h2>Upcoming Events</h2>
      <WebSocketStatus :isConnected="ws.isConnected.value" />
    </div>

    <div v-if="pending" class="loading">Loading events...</div>

    <div v-else-if="error" class="empty-state">
      Failed to load events: {{ error.message }}
    </div>

    <div v-else-if="events?.length === 0" class="empty-state">
      No events available
    </div>

    <div v-else class="events-grid">
      <div
        v-for="event in events"
        :key="event.ID"
        class="event-card"
        @click="router.push(`/events/${event.ID}`)"
      >
        <div
          class="event-image"
          :style="
            event.image_url
              ? `background-image: url(${event.image_url}); background-size: cover;`
              : ''
          "
        ></div>
        <div class="event-content">
          <div class="event-type">{{ event.event_type || 'event' }}</div>
          <div class="event-name">{{ event.name }}</div>
          <div class="event-info">
            <div class="event-info-item">
              📅 {{ formatDate(event.date) }}
            </div>
            <div class="event-info-item">
              📍 {{ event.venue_name }}, {{ event.city }}
            </div>
            <div class="event-info-item">
              👥 {{ getAvailableTickets(event.ID) }}/{{ event.capacity }} available
            </div>
          </div>
          <div class="event-price">${{ event.price.toFixed(2) }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Event, AvailabilityUpdate } from '~/types'

const api = useApi()
const router = useRouter()
const ws = useWebSocket()

const { data: eventsData, pending, error } = await useAsyncData('events', () => api.getEvents())

const events = computed(() => eventsData.value?.data || [])

// Map to store live availability updates
const eventAvailability = ref<Map<number, number>>(new Map())

// Initialize WebSocket connection
onMounted(() => {
  const token = api.getToken()
  if (token) {
    ws.connect(token)
  }

  // Listen for availability updates
  ws.on('availability_update', (update: AvailabilityUpdate) => {
    eventAvailability.value.set(update.event_id, update.available_tickets)
  })
})

// Subscribe to all events when they load
watch(events, (newEvents) => {
  if (newEvents && newEvents.length > 0 && ws.isConnected.value) {
    // Subscribe to all visible events
    newEvents.forEach(event => {
      ws.subscribe(event.ID)
    })
  }
}, { immediate: true })

// Also subscribe when WebSocket connects (in case events loaded before connection)
watch(() => ws.isConnected.value, (connected) => {
  if (connected && events.value && events.value.length > 0) {
    events.value.forEach(event => {
      ws.subscribe(event.ID)
    })
  }
})

// Get available tickets with live updates
const getAvailableTickets = (eventId: number) => {
  return eventAvailability.value.get(eventId) ?? events.value.find(e => e.ID === eventId)?.available_tickets ?? 0
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}
</script>