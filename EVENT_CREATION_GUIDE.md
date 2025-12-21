# Event Creation Implementation Guide

This guide outlines the steps to implement a functional "Create Event" page in the Nuxt web application, interacting with the Go API.

## Prerequisites

-   Ensure the Go API is running (`make dev` or `docker compose up`).
-   Ensure you are logged in to the application to access protected routes.

---

## Step 1: Update API Client

We need to add the `createEvent` method to our API composable.

**File:** `web/composables/useApi.ts`

Add the following function inside the `useApi` composable:

```typescript
// ... inside useApi function

  const createEvent = (eventData: any) => 
    fetchWithAuth<Event>('/events', {
      method: 'POST',
      body: JSON.stringify(eventData),
    })

// ... update return statement
  return {
    // ... existing exports
    createEvent,
  }
```

## Step 2: Define Types (Optional but Recommended)

Update the types definition to include the structure for creating an event.

**File:** `web/types/index.ts`

Add this interface:

```typescript
export interface CreateEventRequest {
  name: string
  description: string
  event_type: string
  venue_name: string
  city: string
  address: string
  date: string // ISO 8601 string
  price: number
  capacity: number
  image_url: string
}
```

## Step 3: Create the Event Creation Page

Create a new page component with a form to collect event details.

**File:** `web/pages/events/create.vue`

```vue
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
```

## Step 4: Add Navigation Link

Update the main navigation to include a link to the creation page.

**File:** `web/app.vue`

Update the `nav-links` section:

```vue
<div class="nav-links">
  <NuxtLink to="/events" class="nav-link">Events</NuxtLink>
  <!-- Add this line -->
  <NuxtLink v-if="isLoggedIn" to="/events/create" class="nav-link">Create Event</NuxtLink>
  
  <NuxtLink v-if="isLoggedIn" to="/profile" class="nav-link">Profile</NuxtLink>
  <!-- ... rest of the links -->
</div>
```

## Step 5: Verify Implementation

1.  Login to the application.
2.  Click the new "Create Event" link in the navigation.
3.  Fill out the form with valid data.
4.  Submit the form.
5.  You should be redirected to the events list, and your new event should appear there.

## Notes on Styling

-   The form reuses existing CSS classes from `assets/css/main.css` (`card`, `form-group`, `btn-primary`, etc.) to maintain consistency.
-   Inline styles are used sparingly for specific layout adjustments not covered by the global CSS.
