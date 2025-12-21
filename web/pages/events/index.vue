<template>
  <div class="page-container">
    <div class="page-header">
      <h2>Upcoming Events</h2>
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
              👥 {{ event.available_tickets }}/{{ event.capacity }} available
            </div>
          </div>
          <div class="event-price">${{ event.price.toFixed(2) }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Event } from '~/types'

const api = useApi()
const router = useRouter()

const { data: eventsData, pending, error } = await useAsyncData('events', () => api.getEvents())

const events = computed(() => eventsData.value?.data || [])

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}
</script>