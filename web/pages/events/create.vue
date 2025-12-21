<template>
  <div class="page-container">
    <div class="card event-details-container">
      <div class="page-header">
        <h2 style="color: #333; margin-bottom: 20px;">Create New Event</h2>
      </div>

      <form @submit.prevent="handleSubmit" class="form-section">
        <div class="booking-form" style="flex-direction: column; align-items: stretch;">
          
          <!-- Basic Info -->
          <div class="form-group">
            <label for="name">Event Name</label>
            <input id="name" v-model="form.name" type="text" required placeholder="e.g. Summer Rock Festival">
          </div>

          <div class="form-group">
            <label for="description">Description</label>
            <input id="description" v-model="form.description" type="text" required placeholder="Brief description of the event">
          </div>

          <div class="form-group">
            <label for="type">Event Type</label>
            <select id="type" v-model="form.event_type" class="form-input" required>
              <option value="concert">Concert</option>
              <option value="tour">Tour</option>
              <option value="standup">Standup</option>
              <option value="lecture">Lecture</option>
              <option value="musical">Musical</option>
              <option value="other">Other</option>
            </select>
          </div>

          <!-- Location -->
          <div class="form-group">
            <label for="venue">Venue Name</label>
            <input id="venue" v-model="form.venue_name" type="text" required placeholder="e.g. City Arena">
          </div>

          <div class="booking-form" style="margin-top: 0;">
            <div class="form-group">
              <label for="city">City</label>
              <input id="city" v-model="form.city" type="text" required placeholder="e.g. New York">
            </div>
            <div class="form-group">
              <label for="address">Address</label>
              <input id="address" v-model="form.address" type="text" required placeholder="e.g. 123 Main St">
            </div>
          </div>

          <!-- Details -->
          <div class="booking-form" style="margin-top: 0;">
            <div class="form-group">
              <label for="date">Date & Time</label>
              <input id="date" v-model="form.date" type="datetime-local" required>
            </div>
            <div class="form-group">
              <label for="capacity">Capacity</label>
              <input id="capacity" v-model="form.capacity" type="number" min="1" required>
            </div>
            <div class="form-group">
              <label for="price">Price ($)</label>
              <input id="price" v-model="form.price" type="number" min="0" step="0.01" required>
            </div>
          </div>

          <div class="form-group">
            <label for="image">Image URL</label>
            <input id="image" v-model="form.image_url" type="url" placeholder="https://example.com/image.jpg">
          </div>

          <!-- Actions -->
          <div class="message error" :class="{ show: error }">{{ error }}</div>
          
          <div style="display: flex; gap: 10px; margin-top: 20px;">
            <button type="submit" :disabled="loading" class="btn-primary" style="flex: 1;">
              {{ loading ? 'Creating...' : 'Create Event' }}
            </button>
            <NuxtLink to="/events" class="button btn-secondary" style="text-decoration: none; text-align: center; color: white; border-radius: 8px;">
              Cancel
            </NuxtLink>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { CreateEventRequest } from '~/types'

definePageMeta({
  middleware: ['auth']
})

const router = useRouter()
const api = useApi()

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
    // Convert date to ISO string for backend
    const payload = {
      ...form,
      date: new Date(form.date).toISOString()
    }

    await api.createEvent(payload)
    router.push('/events')
  } catch (err: any) {
    error.value = err.message || 'Failed to create event'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.form-input {
  width: 100%;
  padding: 12px 15px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 14px;
  background-color: white;
}
</style>
