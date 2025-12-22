<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <div class="flex flex-col md:flex-row justify-between items-center mb-8 gap-4">
      <div>
        <h2 class="text-3xl font-bold text-gray-900">Upcoming Events</h2>
        <p class="text-gray-500 mt-1">
          {{ paginationData?.total || 0 }} events found
        </p>
      </div>
      <WebSocketStatus :isConnected="ws.isConnected.value" />
    </div>

    <div v-if="pending" class="flex justify-center py-20">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
    </div>

    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-2xl p-8 text-center">
      <p class="text-red-600">Failed to load events: {{ error.message }}</p>
      <button @click="refresh()" class="mt-4 px-6 py-2 bg-red-600 text-white rounded-xl hover:bg-red-700 transition-colors">
        Retry
      </button>
    </div>

    <div v-else-if="events?.length === 0" class="bg-gray-50 border border-gray-200 rounded-2xl p-12 text-center">
      <p class="text-gray-500 text-lg">No events found matching your criteria.</p>
      <button @click="clearFilters" class="mt-4 px-6 py-2 bg-primary text-white rounded-xl hover:bg-primary-dark transition-colors">
        Clear Filters
      </button>
    </div>

    <div v-else>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
      <div
        v-for="event in events"
        :key="event.ID"
        class="group bg-white rounded-3xl overflow-hidden shadow-lg shadow-gray-100 hover:shadow-xl hover:shadow-primary/10 border border-gray-100 transition-all duration-300 hover:-translate-y-1 cursor-pointer"
        @click="router.push(`/events/${event.ID}`)"
      >
        <div class="relative h-48 overflow-hidden">
          <div class="absolute inset-0 bg-gray-200 animate-pulse" v-if="!event.image_url"></div>
          <img 
            v-if="event.image_url"
            :src="event.image_url" 
            :alt="event.name"
            class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
            loading="lazy"
          />
          <div class="absolute top-4 left-4">
            <span class="px-3 py-1 rounded-full bg-white/90 backdrop-blur text-xs font-bold uppercase tracking-wider text-gray-800 shadow-sm">
              {{ event.event_type }}
            </span>
          </div>
        </div>
        
        <div class="p-6">
          <h3 class="text-xl font-bold text-gray-900 mb-2 line-clamp-1 group-hover:text-primary transition-colors">
            {{ event.name }}
          </h3>
          
          <div class="space-y-2 mb-4">
            <div class="flex items-center text-gray-500 text-sm">
              <CalendarIcon class="w-4 h-4 mr-2" />
              {{ formatDate(event.date) }}
            </div>
            <div class="flex items-center text-gray-500 text-sm">
              <MapPinIcon class="w-4 h-4 mr-2" />
              {{ event.venue_name }}, {{ event.city }}
            </div>
          </div>

          <div class="flex items-center justify-between pt-4 border-t border-gray-100">
            <div class="flex flex-col">
              <span class="text-xs text-gray-400 font-medium uppercase">Price</span>
              <span class="text-lg font-bold text-primary">${{ event.price.toFixed(2) }}</span>
            </div>
            <div class="text-right">
              <span class="text-xs text-gray-400 font-medium uppercase">Availability</span>
              <div class="flex items-center gap-1 font-medium" :class="getAvailabilityColor(getAvailableTickets(event.ID), event.capacity)">
                {{ getAvailableTickets(event.ID) }} / {{ event.capacity }}
              </div>
            </div>
          </div>
        </div>
      </div>
      </div>

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
import { CalendarIcon, MapPinIcon } from '@heroicons/vue/24/outline'

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
  if (ratio < 0.1) return 'text-red-500'
  if (ratio < 0.3) return 'text-yellow-600'
  return 'text-green-600'
}

const handlePageChange = (page: number) => {
  router.push({ query: { ...route.query, page } })
}

const clearFilters = () => {
  router.push({ query: { page: 1, limit: 20 } })
}
</script>