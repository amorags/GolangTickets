<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- Loading State -->
    <div v-if="pending" class="flex justify-center py-20">
      <div class="relative">
        <div class="w-16 h-16 spinner-cyber"></div>
        <div class="absolute inset-0 w-16 h-16 rounded-full border-2 border-primary/20"></div>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error || !event" class="glass-card p-8 text-center">
      <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-error/10 flex items-center justify-center">
        <ExclamationTriangleIcon class="w-8 h-8 text-error" />
      </div>
      <p class="text-error mb-4">Event not found</p>
      <NuxtLink to="/events" class="btn-ghost inline-flex items-center gap-2">
        <ArrowLeftIcon class="w-4 h-4" />
        Back to events
      </NuxtLink>
    </div>

    <!-- Event Content -->
    <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-10 animate-float-in">
      <!-- Left Column: Image & Description -->
      <div class="space-y-6">
        <!-- Image -->
        <div class="relative rounded-2xl overflow-hidden glass-card">
          <div class="aspect-[4/3] relative">
            <img 
              v-if="event.image_url"
              :src="event.image_url" 
              :alt="event.name"
              class="w-full h-full object-cover"
            />
            <div v-else class="w-full h-full bg-dark-50 flex items-center justify-center">
              <PhotoIcon class="w-16 h-16 text-text-muted" />
            </div>
            
            <!-- Gradient overlay -->
            <div class="absolute inset-0 bg-gradient-to-t from-dark via-transparent to-transparent opacity-40"></div>
            
            <!-- WebSocket status -->
            <div class="absolute top-4 right-4">
              <WebSocketStatus :isConnected="ws.isConnected.value" />
            </div>
          </div>
        </div>

        <!-- Description Card -->
        <div class="glass-card p-6">
          <h3 class="text-lg font-bold text-text mb-4 flex items-center gap-2">
            <InformationCircleIcon class="w-5 h-5 text-primary" />
            About this event
          </h3>
          <p class="text-text-light leading-relaxed">
            {{ event.description || 'No description available for this event.' }}
          </p>
        </div>
      </div>

      <!-- Right Column: Details & Booking -->
      <div class="space-y-6">
        <!-- Event Header -->
        <div class="glass-card p-6">
          <span class="tag-cyber mb-4 inline-block">
            {{ event.event_type }}
          </span>
          <h1 class="text-3xl font-display font-bold text-text mb-3">
            {{ event.name }}
          </h1>
          <div class="flex items-center text-text-light">
            <MapPinIcon class="w-5 h-5 mr-2 text-secondary" />
            <span>{{ event.venue_name }} • {{ event.city }}</span>
          </div>
        </div>

        <!-- Date & Time -->
        <div class="grid grid-cols-2 gap-4">
          <div class="glass-card p-5 text-center">
            <CalendarIcon class="w-6 h-6 text-primary mx-auto mb-2" />
            <p class="text-xs text-text-muted uppercase tracking-wider mb-1">Date</p>
            <p class="font-bold text-text">{{ formatDate(event.date) }}</p>
          </div>
          <div class="glass-card p-5 text-center">
            <ClockIcon class="w-6 h-6 text-secondary mx-auto mb-2" />
            <p class="text-xs text-text-muted uppercase tracking-wider mb-1">Time</p>
            <p class="font-bold text-text">{{ formatTime(event.date) }}</p>
          </div>
        </div>

        <!-- Booking Card -->
        <div class="glass-card p-6 relative overflow-hidden">
          <!-- Decorative glow -->
          <div class="absolute -top-20 -right-20 w-40 h-40 bg-primary/20 rounded-full blur-3xl"></div>
          
          <div class="relative z-10">
            <!-- Price & Availability -->
            <div class="flex justify-between items-end mb-6 pb-6 border-b border-surface-border">
              <div>
                <p class="text-xs text-text-muted uppercase tracking-wider mb-1">Price per ticket</p>
                <p class="text-4xl font-display font-bold text-gradient">${{ event.price.toFixed(2) }}</p>
              </div>
              <div class="text-right">
                <p class="text-xs text-text-muted uppercase tracking-wider mb-1">Available</p>
                <div class="flex items-center gap-2 justify-end">
                  <span 
                    class="text-2xl font-bold transition-all duration-300"
                    :class="{'scale-110 text-success': justUpdated}"
                  >
                    {{ localAvailableTickets }}
                  </span>
                  <span class="text-text-muted">/ {{ event.capacity }}</span>
                </div>
              </div>
            </div>

            <!-- Login prompt -->
            <div v-if="!isLoggedIn" class="text-center py-6">
              <p class="text-text-light mb-4">Sign in to book tickets for this event</p>
              <NuxtLink to="/" class="btn-neon inline-block">
                Sign In to Book
              </NuxtLink>
            </div>

            <!-- Booking Form -->
            <form v-else @submit.prevent="handleBooking" class="space-y-5">
              <div>
                <label class="block text-sm font-medium text-text-light mb-3">Quantity</label>
                <div class="flex items-center gap-4">
                  <button 
                    type="button" 
                    class="w-12 h-12 rounded-xl bg-dark-50 border border-surface-border flex items-center justify-center hover:border-primary hover:text-primary transition-all font-bold text-xl"
                    @click="quantity = Math.max(1, quantity - 1)"
                  >−</button>
                  <div class="flex-1 text-center">
                    <span class="text-3xl font-bold text-text">{{ quantity }}</span>
                  </div>
                  <button 
                    type="button" 
                    class="w-12 h-12 rounded-xl bg-dark-50 border border-surface-border flex items-center justify-center hover:border-primary hover:text-primary transition-all font-bold text-xl"
                    @click="quantity = Math.min(localAvailableTickets, 10, quantity + 1)"
                  >+</button>
                </div>
              </div>

              <!-- Total -->
              <div class="flex justify-between items-center py-4 border-y border-surface-border">
                <span class="text-text-light font-medium">Total</span>
                <span class="text-2xl font-bold text-gradient">${{ (event.price * quantity).toFixed(2) }}</span>
              </div>

              <!-- Submit Button -->
              <button 
                type="submit" 
                class="btn-neon w-full py-4 text-lg"
                :disabled="isBooking || localAvailableTickets === 0"
              >
                <span v-if="isBooking" class="flex items-center justify-center gap-2">
                  <svg class="animate-spin w-5 h-5" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Processing...
                </span>
                <span v-else-if="localAvailableTickets === 0">Sold Out</span>
                <span v-else>Book Now</span>
              </button>
            </form>
          </div>
        </div>

        <!-- Location Card -->
        <div class="glass-card p-6">
          <h3 class="text-sm font-bold text-text-muted uppercase tracking-wider mb-4 flex items-center gap-2">
            <MapPinIcon class="w-4 h-4" />
            Location Details
          </h3>
          <div class="space-y-2 text-text-light">
            <p class="font-medium text-text">{{ event.venue_name }}</p>
            <p>{{ event.address }}</p>
            <p>{{ event.city }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useToast } from "vue-toastification"
import { 
  MapPinIcon, 
  CalendarIcon, 
  ClockIcon, 
  InformationCircleIcon,
  ExclamationTriangleIcon,
  ArrowLeftIcon,
  PhotoIcon
} from '@heroicons/vue/24/outline'
import type { AvailabilityUpdate } from '~/types'

const route = useRoute()
const router = useRouter()
const api = useApi()
const ws = useWebSocket()
const toast = useToast()

const eventId = computed(() => parseInt(route.params.id as string))
const isLoggedIn = ref(false)
const quantity = ref(1)
const isBooking = ref(false)
const localAvailableTickets = ref(0)
const justUpdated = ref(false)

onMounted(() => {
  isLoggedIn.value = !!api.getToken()

  const token = api.getToken()
  if (token) {
    ws.connect(token)
  }

  if (event.value) {
    localAvailableTickets.value = event.value.available_tickets
  }

  ws.on('availability_update', (update: AvailabilityUpdate) => {
    if (update.event_id === eventId.value) {
      localAvailableTickets.value = update.available_tickets
      justUpdated.value = true
      setTimeout(() => {
        justUpdated.value = false
      }, 500)
    }
  })
})

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

watch(event, (newEvent) => {
  if (newEvent) {
    localAvailableTickets.value = newEvent.available_tickets
  }
}, { immediate: true })

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

const formatTime = (dateString: string) => {
  return new Date(dateString).toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
  })
}

const handleBooking = async () => {
  isBooking.value = true
  try {
    await api.createBooking(eventId.value, quantity.value)
    toast.success('Booking successful! Enjoy the event.')
    setTimeout(() => {
      router.push('/profile')
    }, 1500)
  } catch (error: any) {
    toast.error(error.message || 'Booking failed')
  } finally {
    isBooking.value = false
  }
}
</script>
