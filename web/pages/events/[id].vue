<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <div v-if="pending" class="flex justify-center py-20">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
    </div>

    <div v-else-if="error || !event" class="bg-red-50 border border-red-200 rounded-2xl p-8 text-center">
      <p class="text-red-600">Event not found</p>
      <NuxtLink to="/events" class="text-primary hover:underline mt-2 inline-block">Back to events</NuxtLink>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-12 animate-float-in">
      <!-- Left Column: Image -->
      <div class="space-y-6">
        <div class="relative rounded-3xl overflow-hidden shadow-2xl shadow-primary/10 aspect-[4/3]">
          <img 
            v-if="event.image_url"
            :src="event.image_url" 
            :alt="event.name"
            class="w-full h-full object-cover"
          />
          <div v-else class="w-full h-full bg-gray-200 flex items-center justify-center text-gray-400">
            No Image
          </div>
          
          <div class="absolute top-4 right-4">
             <WebSocketStatus :isConnected="ws.isConnected.value" />
          </div>
        </div>

        <div class="bg-white rounded-2xl p-6 border border-gray-100 shadow-lg">
          <h3 class="text-lg font-bold text-gray-900 mb-4">About this event</h3>
          <p class="text-gray-600 leading-relaxed">
            {{ event.description || 'No description available for this event.' }}
          </p>
        </div>
      </div>

      <!-- Right Column: Details & Booking -->
      <div class="space-y-8">
        <div>
          <span class="px-3 py-1 rounded-full bg-primary/10 text-primary font-bold text-sm uppercase tracking-wider">
            {{ event.event_type }}
          </span>
          <h1 class="text-4xl font-bold text-gray-900 mt-4 mb-2">{{ event.name }}</h1>
          <div class="flex items-center text-gray-500 text-lg">
            <MapPinIcon class="w-5 h-5 mr-2" />
            {{ event.venue_name }} • {{ event.city }}
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div class="bg-secondary/10 rounded-2xl p-4">
            <p class="text-sm text-gray-500 mb-1">Date</p>
            <p class="font-bold text-gray-900">{{ formatDate(event.date) }}</p>
          </div>
           <div class="bg-secondary/10 rounded-2xl p-4">
            <p class="text-sm text-gray-500 mb-1">Time</p>
            <p class="font-bold text-gray-900">{{ formatTime(event.date) }}</p>
          </div>
        </div>

        <div class="bg-white rounded-3xl p-8 border border-gray-100 shadow-xl relative overflow-hidden">
          <div class="absolute top-0 right-0 w-32 h-32 bg-gradient-to-br from-primary/20 to-transparent rounded-bl-full -mr-8 -mt-8"></div>
          
          <div class="relative z-10">
            <div class="flex justify-between items-end mb-6">
              <div>
                <p class="text-sm text-gray-500 font-medium uppercase mb-1">Price per ticket</p>
                <p class="text-4xl font-bold text-gray-900">${{ event.price.toFixed(2) }}</p>
              </div>
              <div class="text-right">
                 <p class="text-sm text-gray-500 font-medium uppercase mb-1">Availability</p>
                 <div class="flex items-center gap-2 justify-end">
                   <span 
                     class="text-2xl font-bold transition-all duration-300"
                     :class="{'scale-110 text-green-600': justUpdated}"
                   >
                     {{ localAvailableTickets }}
                   </span>
                   <span class="text-gray-400">/ {{ event.capacity }}</span>
                 </div>
              </div>
            </div>

            <div v-if="!isLoggedIn" class="text-center py-4">
              <p class="text-gray-600 mb-4">Sign in to book tickets for this event.</p>
              <NuxtLink to="/" class="inline-block px-8 py-3 bg-primary text-white font-bold rounded-xl shadow-lg shadow-primary/30 hover:-translate-y-0.5 transition-transform">
                Login to Book
              </NuxtLink>
            </div>

            <form v-else @submit.prevent="handleBooking" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Quantity</label>
                <div class="flex items-center gap-4">
                   <button 
                    type="button" 
                    class="w-10 h-10 rounded-lg bg-gray-100 flex items-center justify-center hover:bg-gray-200 font-bold text-xl"
                    @click="quantity = Math.max(1, quantity - 1)"
                  >-</button>
                  <input
                    v-model.number="quantity"
                    type="number"
                    min="1"
                    :max="Math.min(localAvailableTickets, 10)"
                    class="w-20 text-center font-bold text-xl border-none bg-transparent focus:ring-0"
                    readonly
                  />
                  <button 
                    type="button" 
                    class="w-10 h-10 rounded-lg bg-gray-100 flex items-center justify-center hover:bg-gray-200 font-bold text-xl"
                    @click="quantity = Math.min(localAvailableTickets, 10, quantity + 1)"
                  >+</button>
                </div>
              </div>

              <div class="pt-4 border-t border-gray-100 flex justify-between items-center mb-4">
                <span class="font-bold text-gray-700">Total</span>
                <span class="font-bold text-2xl text-primary">${{ (event.price * quantity).toFixed(2) }}</span>
              </div>

              <button 
                type="submit" 
                class="w-full py-4 bg-gradient-to-r from-primary to-primary-dark text-white font-bold text-lg rounded-xl shadow-lg shadow-primary/30 hover:shadow-xl hover:-translate-y-0.5 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                :disabled="isBooking || localAvailableTickets === 0"
              >
                <span v-if="isBooking">Processing...</span>
                <span v-else-if="localAvailableTickets === 0">Sold Out</span>
                <span v-else>Book Tickets</span>
              </button>
            </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useToast } from "vue-toastification";
import { MapPinIcon } from '@heroicons/vue/24/outline'
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
    weekday: 'long',
    month: 'long',
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