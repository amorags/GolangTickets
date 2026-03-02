<template>
  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <div class="animate-float-in">
      <!-- Header -->
      <div class="text-center mb-10">
        <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-primary to-secondary mb-4 shadow-glow">
          <SparklesIcon class="w-8 h-8 text-white" />
        </div>
        <h1 class="text-3xl font-display font-bold text-gradient mb-2">Create New Event</h1>
        <p class="text-text-light">Fill in the details to publish your event</p>
      </div>

      <!-- Form Card -->
      <div class="glass-card overflow-hidden">
        <form @submit.prevent="handleSubmit" class="p-8 space-y-8">
          <!-- Basic Info Section -->
          <div class="space-y-5">
            <h3 class="text-sm font-bold text-primary uppercase tracking-wider flex items-center gap-2">
              <InformationCircleIcon class="w-4 h-4" />
              Event Details
            </h3>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
              <div class="col-span-2">
                <label for="name" class="block text-sm font-medium text-text-light mb-2">Event Name</label>
                <input 
                  id="name" 
                  v-model="form.name" 
                  type="text" 
                  required 
                  placeholder="e.g. Summer Rock Festival" 
                  class="input-cyber" 
                />
              </div>

              <div class="col-span-2">
                <label for="description" class="block text-sm font-medium text-text-light mb-2">Description</label>
                <textarea 
                  id="description" 
                  v-model="form.description" 
                  required 
                  placeholder="Brief description of the event" 
                  rows="3"
                  class="input-cyber resize-none"
                ></textarea>
              </div>

              <div>
                <label for="type" class="block text-sm font-medium text-text-light mb-2">Event Type</label>
                <select 
                  id="type" 
                  v-model="form.event_type" 
                  required
                  class="input-cyber cursor-pointer"
                >
                  <option value="concert">Concert</option>
                  <option value="tour">Tour</option>
                  <option value="standup">Standup</option>
                  <option value="lecture">Lecture</option>
                  <option value="musical">Musical</option>
                  <option value="other">Other</option>
                </select>
              </div>

              <div>
                <label for="image" class="block text-sm font-medium text-text-light mb-2">Image URL</label>
                <input 
                  id="image" 
                  v-model="form.image_url" 
                  type="url" 
                  placeholder="https://example.com/image.jpg"
                  class="input-cyber" 
                />
              </div>
            </div>
          </div>

          <div class="h-px bg-gradient-to-r from-transparent via-surface-border to-transparent"></div>

          <!-- Location Section -->
          <div class="space-y-5">
            <h3 class="text-sm font-bold text-secondary uppercase tracking-wider flex items-center gap-2">
              <MapPinIcon class="w-4 h-4" />
              Location
            </h3>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
              <div class="col-span-2">
                <label for="venue" class="block text-sm font-medium text-text-light mb-2">Venue Name</label>
                <input 
                  id="venue" 
                  v-model="form.venue_name" 
                  type="text" 
                  required 
                  placeholder="e.g. City Arena"
                  class="input-cyber" 
                />
              </div>

              <div>
                <label for="city" class="block text-sm font-medium text-text-light mb-2">City</label>
                <input 
                  id="city" 
                  v-model="form.city" 
                  type="text" 
                  required 
                  placeholder="e.g. New York"
                  class="input-cyber" 
                />
              </div>

              <div>
                <label for="address" class="block text-sm font-medium text-text-light mb-2">Address</label>
                <input 
                  id="address" 
                  v-model="form.address" 
                  type="text" 
                  required 
                  placeholder="e.g. 123 Main St"
                  class="input-cyber" 
                />
              </div>
            </div>
          </div>

          <div class="h-px bg-gradient-to-r from-transparent via-surface-border to-transparent"></div>

          <!-- Date & Pricing Section -->
          <div class="space-y-5">
            <h3 class="text-sm font-bold text-accent uppercase tracking-wider flex items-center gap-2">
              <CurrencyDollarIcon class="w-4 h-4" />
              Date & Pricing
            </h3>
            
            <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
              <div>
                <label for="date" class="block text-sm font-medium text-text-light mb-2">Date & Time</label>
                <input 
                  id="date" 
                  v-model="form.date" 
                  type="datetime-local" 
                  required
                  class="input-cyber" 
                />
              </div>

              <div>
                <label for="capacity" class="block text-sm font-medium text-text-light mb-2">Capacity</label>
                <input 
                  id="capacity" 
                  v-model="form.capacity" 
                  type="number" 
                  min="1" 
                  required
                  class="input-cyber" 
                />
              </div>

              <div>
                <label for="price" class="block text-sm font-medium text-text-light mb-2">Price ($)</label>
                <input 
                  id="price" 
                  v-model="form.price" 
                  type="number" 
                  min="0" 
                  step="0.01" 
                  required
                  class="input-cyber" 
                />
              </div>
            </div>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="p-4 rounded-xl bg-error/10 border border-error/30 text-error text-sm animate-fade-in">
            <div class="flex items-center gap-2">
              <ExclamationCircleIcon class="w-5 h-5" />
              {{ error }}
            </div>
          </div>

          <!-- Actions -->
          <div class="flex gap-4 pt-4">
            <button 
              type="submit" 
              :disabled="loading" 
              class="btn-neon flex-1 py-4"
            >
              <span v-if="loading" class="flex items-center justify-center gap-2">
                <svg class="animate-spin w-5 h-5" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Creating...
              </span>
              <span v-else class="flex items-center justify-center gap-2">
                <SparklesIcon class="w-5 h-5" />
                Create Event
              </span>
            </button>
            
            <NuxtLink 
              to="/events" 
              class="btn-ghost px-8"
            >
              Cancel
            </NuxtLink>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useToast } from "vue-toastification"
import { 
  SparklesIcon, 
  InformationCircleIcon, 
  MapPinIcon, 
  CurrencyDollarIcon,
  ExclamationCircleIcon 
} from '@heroicons/vue/24/outline'
import type { CreateEventRequest } from '~/types'

definePageMeta({
  middleware: ['auth']
})

const router = useRouter()
const api = useApi()
const toast = useToast()

const loading = ref(false)
const error = ref('')

const form = reactive<CreateEventRequest>({
  name: '',
  description: '',
  event_type: 'concert',
  venue_name: '',
  city: '',
  address: '',
  date: '',
  price: 0,
  capacity: 100,
  image_url: ''
})

const handleSubmit = async () => {
  loading.value = true
  error.value = ''

  try {
    const payload = {
      ...form,
      date: new Date(form.date).toISOString()
    }

    await api.createEvent(payload)
    toast.success('Event created successfully!')
    router.push('/events')
  } catch (err: any) {
    error.value = err.message || 'Failed to create event'
    toast.error(error.value)
  } finally {
    loading.value = false
  }
}
</script>
