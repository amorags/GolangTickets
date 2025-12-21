<template>
  <div>
    <NuxtRouteAnnouncer />
    <nav v-if="showNav" class="navbar">
      <div class="nav-container">
        <h1 class="nav-title">Ticket App</h1>
        <div class="nav-links">
          <NuxtLink to="/events" class="nav-link">Events</NuxtLink>
          <NuxtLink v-if="isLoggedIn" to="/events/create" class="nav-link">Create Event</NuxtLink>
          <NuxtLink v-if="isLoggedIn" to="/profile" class="nav-link">Profile</NuxtLink>
          <a v-if="isLoggedIn" href="#" class="nav-link" @click.prevent="handleLogout">Logout</a>
          <NuxtLink v-else to="/" class="nav-link">Login</NuxtLink>
        </div>
      </div>
    </nav>
    <NuxtPage />
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const api = useApi()

const isLoggedIn = ref(false)
const showNav = computed(() => route.path !== '/')

onMounted(() => {
  isLoggedIn.value = !!api.getToken()
})

watch(() => route.path, () => {
  if (import.meta.client) {
    isLoggedIn.value = !!api.getToken()
  }
})

const handleLogout = () => {
  api.removeToken()
  isLoggedIn.value = false
  router.push('/')
}
</script>