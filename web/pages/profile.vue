<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <div class="flex items-center justify-between mb-8">
      <h1 class="text-3xl font-bold text-gray-900">My Account</h1>
      <button 
        @click="handleLogout"
        class="text-error hover:text-error/80 font-medium text-sm flex items-center gap-2"
      >
        Sign out
      </button>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Left Column: User Profile -->
      <div class="lg:col-span-1 space-y-6">
        <div class="bg-white rounded-3xl p-6 shadow-lg border border-gray-100">
          <div class="flex items-center gap-4 mb-6">
            <div class="w-16 h-16 rounded-full bg-gradient-to-br from-secondary to-secondary-dark flex items-center justify-center text-white text-2xl font-bold">
              {{ user?.username.charAt(0).toUpperCase() }}
            </div>
            <div>
              <h2 class="text-xl font-bold text-gray-900">{{ user?.username }}</h2>
              <p class="text-gray-500 text-sm">Member since {{ formatDate(user?.created_at) }}</p>
            </div>
          </div>
          
          <div class="space-y-4">
            <div class="p-4 bg-gray-50 rounded-2xl">
              <p class="text-xs text-gray-400 font-medium uppercase mb-1">Email Address</p>
              <p class="font-medium text-gray-900">{{ user?.email }}</p>
            </div>
            <div class="p-4 bg-gray-50 rounded-2xl">
              <p class="text-xs text-gray-400 font-medium uppercase mb-1">User ID</p>
              <p class="font-medium text-gray-900">#{{ user?.id }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: Bookings -->
      <div class="lg:col-span-2 space-y-6">
        <h2 class="text-2xl font-bold text-gray-900">My Bookings</h2>

        <div v-if="bookingsPending" class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        </div>

        <div v-else-if="bookingsError" class="bg-red-50 p-6 rounded-2xl border border-red-100 text-red-600">
          Failed to load bookings.
        </div>

        <div v-else-if="bookings.length === 0" class="bg-white rounded-3xl p-12 text-center border border-gray-100 shadow-sm">
          <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <TicketIcon class="w-8 h-8 text-gray-400" />
          </div>
          <h3 class="text-lg font-bold text-gray-900 mb-2">No bookings yet</h3>
          <p class="text-gray-500 mb-6">You haven't booked any events yet.</p>
          <NuxtLink to="/events" class="px-6 py-2 bg-primary text-white rounded-xl font-bold shadow-md hover:bg-primary-dark transition-colors">
            Browse Events
          </NuxtLink>
        </div>

        <div v-else class="space-y-4">
          <div 
            v-for="booking in bookings" 
            :key="booking.ID"
            class="bg-white rounded-2xl p-5 border border-gray-100 shadow-sm hover:shadow-md transition-shadow group"
          >
            <div class="flex items-start gap-5">
              <!-- Event Image Thumbnail -->
              <div 
                class="w-24 h-24 rounded-xl bg-gray-200 flex-shrink-0 bg-cover bg-center"
                :style="booking.event?.image_url ? `background-image: url(${booking.event.image_url})` : ''"
              >
                <div v-if="!booking.event?.image_url" class="w-full h-full flex items-center justify-center text-gray-400 text-xs">
                  No Image
                </div>
              </div>

              <div class="flex-1 min-w-0">
                <div class="flex justify-between items-start mb-2">
                   <div>
                      <h3 class="text-lg font-bold text-gray-900 group-hover:text-primary transition-colors truncate">
                        {{ booking.event?.name || 'Unknown Event' }}
                      </h3>
                      <p class="text-sm text-gray-500 flex items-center gap-1 mt-1">
                         <CalendarIcon class="w-4 h-4" />
                         {{ booking.event ? formatDate(booking.event.date) : 'Date unavailable' }}
                         <span v-if="booking.event" class="mx-1">•</span>
                         <span v-if="booking.event">{{ booking.event.venue_name }}</span>
                      </p>
                   </div>
                   <UiBadge :variant="booking.status === 'confirmed' ? 'success' : 'error'">
                     {{ booking.status }}
                   </UiBadge>
                </div>
                
                <div class="flex items-end justify-between mt-4">
                  <div class="text-sm">
                    <span class="text-gray-500">Tickets: </span>
                    <span class="font-bold text-gray-900">{{ booking.quantity }}</span>
                    <span class="mx-2 text-gray-300">|</span>
                    <span class="text-gray-500">Total: </span>
                    <span class="font-bold text-primary">${{ booking.total_price.toFixed(2) }}</span>
                  </div>

                  <button 
                    v-if="booking.status === 'confirmed'"
                    @click="confirmCancel(booking)"
                    class="text-sm text-red-500 font-medium hover:text-red-700 hover:bg-red-50 px-3 py-1.5 rounded-lg transition-colors"
                  >
                    Cancel Booking
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <UiConfirmDialog
      :is-open="isCancelModalOpen"
      title="Cancel Booking"
      description="Are you sure you want to cancel this booking? This action cannot be undone and your tickets will be released."
      confirm-text="Yes, Cancel Booking"
      :is-danger="true"
      @close="isCancelModalOpen = false"
      @confirm="handleCancelBooking"
    />
  </div>
</template>

<script setup lang="ts">
import { TicketIcon, CalendarIcon } from '@heroicons/vue/24/outline'
import { useToast } from "vue-toastification";
import type { Booking } from '~/types'

const api = useApi()
const router = useRouter()
const toast = useToast()

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

// Modal state
const isCancelModalOpen = ref(false)
const bookingToCancel = ref<Booking | null>(null)

const confirmCancel = (booking: Booking) => {
  bookingToCancel.value = booking
  isCancelModalOpen.value = true
}

const handleCancelBooking = async () => {
  if (!bookingToCancel.value) return

  try {
    await api.cancelBooking(bookingToCancel.value.ID)
    toast.success('Booking cancelled successfully')
    isCancelModalOpen.value = false
    refreshBookings()
  } catch (error: any) {
    toast.error(error.message || 'Failed to cancel booking')
    isCancelModalOpen.value = false
  }
}

const formatDate = (dateString?: string) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

const handleLogout = () => {
  api.removeToken()
  router.push('/')
}
</script>