<template>
  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <div class="bg-white rounded-3xl shadow-xl border border-gray-100 overflow-hidden">
      <div class="px-8 py-6 border-b border-gray-100 bg-gray-50/50">
        <h2 class="text-2xl font-bold text-gray-900">Create New Event</h2>
        <p class="text-gray-500 text-sm mt-1">Fill in the details to publish your event</p>
      </div>

      <form @submit.prevent="handleSubmit" class="p-8 space-y-6">
        <!-- Basic Info -->
        <div class="space-y-4">
          <h3 class="text-sm font-bold text-gray-400 uppercase tracking-wider">Event Details</h3>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="col-span-2">
              <label for="name" class="block text-sm font-medium text-gray-700 mb-1">Event Name</label>
              <input id="name" v-model="form.name" type="text" required placeholder="e.g. Summer Rock Festival" 
                class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all" />
            </div>

            <div class="col-span-2">
              <label for="description" class="block text-sm font-medium text-gray-700 mb-1">Description</label>
              <textarea id="description" v-model="form.description" required placeholder="Brief description of the event" rows="3"
                class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"></textarea>
            </div>

            <div>
              <label for="type" class="block text-sm font-medium text-gray-700 mb-1">Event Type</label>
              <select id="type" v-model="form.event_type" required
                class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all cursor-pointer appearance-none">
                <option value="concert">Concert</option>
                <option value="tour">Tour</option>
                <option value="standup">Standup</option>
                <option value="lecture">Lecture</option>
                <option value="musical">Musical</option>
                <option value="other">Other</option>
              </select>
            </div>

            <div>
               <label for="image" class="block text-sm font-medium text-gray-700 mb-1">Image URL</label>
               <input id="image" v-model="form.image_url" type="url" placeholder="https://example.com/image.jpg"
                class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all" />
            </div>
          </div>
        </div>

        <hr class="border-gray-100" />

        <!-- Location -->
        <div class="space-y-4">
          <h3 class="text-sm font-bold text-gray-400 uppercase tracking-wider">Location</h3>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="col-span-2">
              <label for="venue" class="block text-sm font-medium text-gray-700 mb-1">Venue Name</label>
              <input id="venue" v-model="form.venue_name" type="text" required placeholder="e.g. City Arena"
                class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all" />
            </div>

            <div>
              <label for="city" class="block text-sm font-medium text-gray-700 mb-1">City</label>
              <input id="city" v-model="form.city" type="text" required placeholder="e.g. New York"
                class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all" />
            </div>

            <div>
              <label for="address" class="block text-sm font-medium text-gray-700 mb-1">Address</label>
              <input id="address" v-model="form.address" type="text" required placeholder="e.g. 123 Main St"
                class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all" />
            </div>
          </div>
        </div>

        <hr class="border-gray-100" />

        <!-- Details -->
        <div class="space-y-4">
          <h3 class="text-sm font-bold text-gray-400 uppercase tracking-wider">Date & Pricing</h3>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div>
              <label for="date" class="block text-sm font-medium text-gray-700 mb-1">Date & Time</label>
              <input id="date" v-model="form.date" type="datetime-local" required
                class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all" />
            </div>

            <div>
              <label for="capacity" class="block text-sm font-medium text-gray-700 mb-1">Capacity</label>
              <input id="capacity" v-model="form.capacity" type="number" min="1" required
                class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all" />
            </div>

            <div>
              <label for="price" class="block text-sm font-medium text-gray-700 mb-1">Price ($)</label>
              <input id="price" v-model="form.price" type="number" min="0" step="0.01" required
                class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all" />
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="pt-6 flex gap-4">
          <button type="submit" :disabled="loading" 
            class="flex-1 py-3.5 bg-gradient-to-r from-primary to-primary-dark text-white font-bold rounded-xl shadow-lg shadow-primary/30 hover:shadow-xl hover:-translate-y-0.5 transition-all disabled:opacity-50 disabled:cursor-not-allowed">
            {{ loading ? 'Creating Event...' : 'Create Event' }}
          </button>
          
          <NuxtLink to="/events" 
            class="px-8 py-3.5 bg-gray-100 text-gray-600 font-bold rounded-xl hover:bg-gray-200 transition-colors">
            Cancel
          </NuxtLink>
        </div>

        <div v-if="error" class="bg-red-50 text-red-600 p-4 rounded-xl text-center text-sm font-medium animate-fade-in">
          {{ error }}
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useToast } from "vue-toastification";
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