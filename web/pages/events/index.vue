<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- Header -->
    <div class="flex flex-col md:flex-row justify-between items-start md:items-center mb-10 gap-4">
      <div>
        <h1 class="text-4xl font-display font-bold text-gradient mb-2">
          Upcoming Events
        </h1>
        <p class="text-text-light flex items-center gap-2">
          <span class="w-2 h-2 rounded-full bg-success animate-pulse"></span>
          {{ paginationData?.total || 0 }} events available
        </p>
      </div>
      <WebSocketStatus :isConnected="ws.isConnected.value" />
    </div>

    <!-- Loading State -->
    <div v-if="pending" class="flex justify-center py-20">
      <div class="relative">
        <div class="w-16 h-16 spinner-cyber"></div>
        <div class="absolute inset-0 w-16 h-16 rounded-full border-2 border-primary/20"></div>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="glass-card p-8 text-center">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-error/10 flex items-center justify-center">
        <ExclamationTriangleIcon class="w-8 h-8 text-error" />
      </div>
      <p class="text-error mb-4">Failed to load events: {{ error.message }}</p>
      <button @click="refresh()" class="btn-ghost border-error/30 text-error hover:bg-error/10">
        Try Again
      </button>
    </div>

    <!-- Empty State -->
    <div v-else-if="events?.length === 0" class="glass-card p-12 text-center">
      <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-surface-light flex items-center justify-center">
        <CalendarDaysIcon class="w-10 h-10 text-text-muted" />
      </div>
      <h3 class="text-xl font-bold text-text mb-2">No events found</h3>
      <p class="text-text-light mb-6">Try adjusting your filters or check back later</p>
      <button @click="clearFilters" class="btn-neon">
        Clear Filters
      </button>
    </div>

    <!-- Events Grid -->
    <div v-else>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="(event, index) in events"
          :key="event.ID"
          class="stagger-item"
          :style="{ animationDelay: `${index * 0.05}s` }"
        >
          <NuxtLink 
            :to="`/events/${event.ID}`"
            class="block glass-card overflow-hidden card-hover group"
          >
            <!-- Image -->
            <div class="relative h-52 overflow-hidden">
              <div class="absolute inset-0 bg-dark-50 animate-pulse" v-if="!event.image_url"></div>
              <img 
                v-if="event.image_url"
                :src="event.image_url" 
                :alt="event.name"
                class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110"
                loading="lazy"
              />
              <!-- Gradient overlay -->
              <div class="absolute inset-0 bg-gradient-to-t from-dark via-transparent to-transparent opacity-60"></div>
              
              <!-- Event type badge -->
              <div class="absolute top-4 left-4">
                <span class="tag-cyber">
                  {{ event.event_type }}
                </span>
              </div>

              <!-- Price badge -->
              <div class="absolute top-4 right-4">
                <div class="px-3 py-1.5 rounded-lg bg-dark-100/80 backdrop-blur-sm border border-surface-border">
                  <span class="text-lg font-bold text-gradient">${{ event.price.toFixed(2) }}</span>
                </div>
              </div>
            </div>
            
            <!-- Content -->
            <div class="p-5">
              <h3 class="text-lg font-bold text-text mb-3 line-clamp-1 group-hover:text-primary transition-colors">
                {{ event.name }}
              </h3>
              
              <div class="space-y-2 mb-4">
                <div class="flex items-center text-text-light text-sm">
                  <CalendarIcon class="w-4 h-4 mr-2 text-primary" />
                  {{ formatDate(event.date) }}
                </div>
                <div class="flex items-center text-text-light text-sm">
                  <MapPinIcon class="w-4 h-4 mr-2 text-secondary" />
                  <span class="truncate">{{ event.venue_name }}, {{ event.city }}</span>
                </div>
              </div>

              <!-- Availability bar -->
              <div class="pt-4 border-t border-surface-border">
                <div class="flex items-center justify-between mb-2">
                  <span class="text-xs text-text-muted uppercase tracking-wider">Availability</span>
                  <span 
                    class="text-sm font-bold"
                    :class="getAvailabilityColor(getAvailableTickets(event.ID), event.capacity)"
                  >
                    {{ getAvailableTickets(event.ID) }} / {{ event.capacity }}
                  </span>
                </div>
                <div class="h-1.5 bg-dark-100 rounded-full overflow-hidden">
                  <div 
                    class="h-full rounded-full transition-all duration-500"
                    :class="getAvailabilityBarColor(getAvailableTickets(event.ID), event.capacity)"
                    :style="{ width: `${(getAvailableTickets(event.ID) / event.capacity) * 100}%` }"
                  ></div>
                </div>
              </div>
            </div>
          </NuxtLink>
        </div>
      </div>

      <!-- Pagination -->
      <Pagination
        v-if="paginationData"
        :page="paginationData.page"
        :total-pages="paginationData.total_pages"
        :total="paginationData.total"
        :limit="paginationData.limit"
        :has-next="paginationData.has_next"
        :has-previous="paginationData.has_previous"
        @update:page="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Event, AvailabilityUpdate, EventFilters, PaginatedEventsResponse } from '~/types'
import { CalendarIcon, MapPinIcon, ExclamationTriangleIcon, CalendarDaysIcon } from '@heroicons/vue/24/outline'

const api = useApi()
const router = useRouter()
const route = useRoute()
const ws = useWebSocket()

const queryParams = computed<EventFilters>(() => ({
  page: Number(route.query.page) || 1,
  limit: Number(route.query.limit) || 20,
  search: route.query.search as string || undefined,
  type: route.query.type as string || undefined,
  city: route.query.city as string || undefined,
  date_from: route.query.date_from as string || undefined,
  date_to: route.query.date_to as string || undefined,
  price_min: route.query.price_min ? Number(route.query.price_min) : undefined,
  price_max: route.query.price_max ? Number(route.query.price_max) : undefined,
}))

const { data: eventsData, pending, error, refresh } = await useAsyncData(
  'events',
  () => api.getEvents(queryParams.value),
  {
    watch: [queryParams],
    deep: true
  }
)

const paginationData = computed(() => eventsData.value?.data as PaginatedEventsResponse | null)
const events = computed(() => paginationData.value?.events || [])

const eventAvailability = ref<Map<number, number>>(new Map())

onMounted(() => {
  const token = api.getToken()
  if (token) {
    ws.connect(token)
  }

  ws.on('availability_update', (update: AvailabilityUpdate) => {
    eventAvailability.value.set(update.event_id, update.available_tickets)
  })
})

watch(events, (newEvents) => {
  if (newEvents && newEvents.length > 0 && ws.isConnected.value) {
    newEvents.forEach(event => {
      ws.subscribe(event.ID)
    })
  }
}, { immediate: true })

watch(() => ws.isConnected.value, (connected) => {
  if (connected && events.value && events.value.length > 0) {
    events.value.forEach(event => {
      ws.subscribe(event.ID)
    })
  }
})

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

const getAvailabilityColor = (available: number, capacity: number) => {
  const ratio = available / capacity
  if (ratio < 0.1) return 'text-error'
  if (ratio < 0.3) return 'text-warning'
  return 'text-success'
}

const getAvailabilityBarColor = (available: number, capacity: number) => {
  const ratio = available / capacity
  if (ratio < 0.1) return 'bg-gradient-to-r from-error to-error/60'
  if (ratio < 0.3) return 'bg-gradient-to-r from-warning to-warning/60'
  return 'bg-gradient-to-r from-success to-success/60'
}

const handlePageChange = (page: number) => {
  router.push({ query: { ...route.query, page } })
}

const clearFilters = () => {
  router.push({ query: { page: 1, limit: 20 } })
}
</script>
