<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- Header -->
    <div class="flex items-center justify-between mb-10">
      <div>
        <h1 class="text-4xl font-display font-bold text-gradient mb-2">My Account</h1>
        <p class="text-text-light">Manage your profile and bookings</p>
      </div>
      <button 
        @click="handleLogout"
        class="btn-ghost border-error/30 text-error hover:bg-error/10 flex items-center gap-2"
      >
        <ArrowRightOnRectangleIcon class="w-4 h-4" />
        Sign out
      </button>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Left Column: User Profile -->
      <div class="lg:col-span-1 space-y-6">
        <!-- Profile Card -->
        <div class="glass-card p-6 text-center">
          <div class="relative inline-block mb-4">
            <div class="w-24 h-24 rounded-2xl bg-gradient-to-br from-primary via-secondary to-accent flex items-center justify-center text-white text-4xl font-bold shadow-glow">
              {{ user?.username.charAt(0).toUpperCase() }}
            </div>
            <div class="absolute -bottom-1 -right-1 w-6 h-6 rounded-full bg-success border-4 border-dark flex items-center justify-center">
              <CheckIcon class="w-3 h-3 text-white" />
            </div>
          </div>
          <h2 class="text-xl font-bold text-text mb-1">{{ user?.username }}</h2>
          <p class="text-text-light text-sm">Member since {{ formatDate(user?.created_at) }}</p>
        </div>

        <!-- Stats Card -->
        <div class="glass-card p-6">
          <h3 class="text-sm font-bold text-text-muted uppercase tracking-wider mb-4">Account Stats</h3>
          <div class="space-y-4">
            <div class="flex items-center justify-between p-3 rounded-xl bg-dark-50 border border-surface-border">
              <div class="flex items-center gap-3">
                <EnvelopeIcon class="w-5 h-5 text-primary" />
                <span class="text-sm text-text-light">Email</span>
              </div>
              <span class="text-sm text-text truncate max-w-[150px]">{{ user?.email }}</span>
            </div>
            <div class="flex items-center justify-between p-3 rounded-xl bg-dark-50 border border-surface-border">
              <div class="flex items-center gap-3">
                <IdentificationIcon class="w-5 h-5 text-secondary" />
                <span class="text-sm text-text-light">User ID</span>
              </div>
              <span class="text-sm text-text">#{{ user?.id }}</span>
            </div>
            <div class="flex items-center justify-between p-3 rounded-xl bg-dark-50 border border-surface-border">
              <div class="flex items-center gap-3">
                <TicketIcon class="w-5 h-5 text-accent" />
                <span class="text-sm text-text-light">Bookings</span>
              </div>
              <span class="text-sm font-bold text-text">{{ bookings.length }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: Bookings -->
      <div class="lg:col-span-2 space-y-6">
        <div class="flex items-center justify-between">
          <h2 class="text-2xl font-display font-bold text-text">My Bookings</h2>
          <NuxtLink to="/events" class="btn-ghost text-sm">
            Browse Events
          </NuxtLink>
        </div>

        <!-- Loading State -->
        <div v-if="bookingsPending" class="flex justify-center py-12">
          <div class="w-12 h-12 spinner-cyber"></div>
        </div>

        <!-- Error State -->
        <div v-else-if="bookingsError" class="glass-card p-8 text-center">
          <ExclamationTriangleIcon class="w-12 h-12 text-error mx-auto mb-4" />
          <p class="text-error">Failed to load bookings</p>
        </div>

        <!-- Empty State -->
        <div v-else-if="bookings.length === 0" class="glass-card p-12 text-center">
          <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-surface-light flex items-center justify-center">
            <TicketIcon class="w-10 h-10 text-text-muted" />
          </div>
          <h3 class="text-xl font-bold text-text mb-2">No bookings yet</h3>
          <p class="text-text-light mb-6">Start exploring events and book your tickets!</p>
          <NuxtLink to="/events" class="btn-neon inline-block">
            Browse Events
          </NuxtLink>
        </div>

        <!-- Bookings List -->
        <div v-else class="space-y-4">
          <div 
            v-for="(booking, index) in bookings" 
            :key="booking.ID"
            class="glass-card p-5 card-hover stagger-item"
            :style="{ animationDelay: `${index * 0.05}s` }"
          >
            <div class="flex items-start gap-5">
              <!-- Event Image -->
              <div 
                class="w-28 h-28 rounded-xl bg-dark-50 flex-shrink-0 bg-cover bg-center border border-surface-border overflow-hidden"
              >
                <img 
                  v-if="booking.event?.image_url"
                  :src="booking.event.image_url"
                  :alt="booking.event?.name"
                  class="w-full h-full object-cover"
                />
                <div v-else class="w-full h-full flex items-center justify-center">
                  <PhotoIcon class="w-8 h-8 text-text-muted" />
                </div>
              </div>

              <div class="flex-1 min-w-0">
                <div class="flex justify-between items-start mb-3">
                  <div>
                    <h3 class="text-lg font-bold text-text group-hover:text-primary transition-colors truncate">
                      {{ booking.event?.name || 'Unknown Event' }}
                    </h3>
                    <div class="flex items-center gap-2 mt-1 text-sm text-text-light">
                      <CalendarIcon class="w-4 h-4 text-primary" />
                      <span>{{ booking.event ? formatDate(booking.event.date) : 'Date unavailable' }}</span>
                      <span class="text-text-muted">•</span>
                      <MapPinIcon class="w-4 h-4 text-secondary" />
                      <span class="truncate">{{ booking.event?.venue_name || 'Venue TBD' }}</span>
                    </div>
                  </div>
                  <UiBadge :variant="booking.status === 'confirmed' ? 'success' : 'error'">
                    {{ booking.status }}
                  </UiBadge>
                </div>
                
                <div class="flex items-end justify-between mt-4 pt-4 border-t border-surface-border">
                  <div class="flex items-center gap-4 text-sm">
                    <div>
                      <span class="text-text-muted">Tickets:</span>
                      <span class="font-bold text-text ml-1">{{ booking.quantity }}</span>
                    </div>
                    <div>
                      <span class="text-text-muted">Total:</span>
                      <span class="font-bold text-gradient ml-1">${{ booking.total_price.toFixed(2) }}</span>
                    </div>
                  </div>

                  <button 
                    v-if="booking.status === 'confirmed'"
                    @click="confirmCancel(booking)"
                    class="text-sm text-error font-medium hover:text-error-light hover:bg-error/10 px-4 py-2 rounded-lg transition-colors"
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

    <!-- Cancel Confirmation Modal -->
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
import { 
  TicketIcon, 
  CalendarIcon, 
  MapPinIcon, 
  EnvelopeIcon,
  IdentificationIcon,
  ArrowRightOnRectangleIcon,
  CheckIcon,
  ExclamationTriangleIcon,
  PhotoIcon
} from '@heroicons/vue/24/outline'
import { useToast } from "vue-toastification"
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
  localStorage.removeItem('user_email')
  router.push('/')
}
</script>
