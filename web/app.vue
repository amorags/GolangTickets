<template>
  <div class="min-h-screen bg-dark font-sans text-text noise-overlay">
    <NuxtRouteAnnouncer />
    
    <!-- Animated background blobs -->
    <div class="fixed inset-0 overflow-hidden pointer-events-none z-0">
      <div class="absolute top-1/4 left-1/4 w-96 h-96 bg-primary/10 rounded-full blur-3xl animate-float"></div>
      <div class="absolute bottom-1/4 right-1/4 w-80 h-80 bg-secondary/10 rounded-full blur-3xl animate-float" style="animation-delay: -3s;"></div>
      <div class="absolute top-1/2 right-1/3 w-64 h-64 bg-accent/10 rounded-full blur-3xl animate-float" style="animation-delay: -5s;"></div>
    </div>

    <nav v-if="showNav" class="sticky top-0 z-40 glass-card border-b border-surface-border rounded-none">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <!-- Logo -->
          <div class="flex flex-shrink-0 items-center">
            <NuxtLink to="/events" class="flex items-center gap-3 group">
              <div class="relative w-10 h-10 rounded-xl bg-gradient-to-br from-primary to-secondary flex items-center justify-center text-white font-bold text-xl shadow-glow group-hover:shadow-glow-lg transition-all duration-300">
                <span class="relative z-10">T</span>
                <div class="absolute inset-0 rounded-xl bg-gradient-to-br from-primary to-secondary opacity-0 group-hover:opacity-100 blur-xl transition-opacity duration-300"></div>
              </div>
              <span class="font-display font-bold text-xl text-gradient hidden sm:block">
                TicketApp
              </span>
            </NuxtLink>
          </div>

          <!-- Search Bar -->
          <div class="hidden md:flex flex-1 justify-center px-8">
            <NavBarSearch />
          </div>
          
          <!-- Right Actions -->
          <div class="flex items-center gap-3">
            <NuxtLink 
              to="/events" 
              class="p-2.5 rounded-xl text-text-light hover:text-primary hover:bg-primary/10 transition-all duration-300 relative group"
              title="Events"
            >
              <CalendarIcon class="w-5 h-5" />
              <span class="absolute -bottom-10 left-1/2 -translate-x-1/2 px-3 py-1.5 bg-dark-50 text-text text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap pointer-events-none border border-surface-border shadow-lg">
                Browse Events
              </span>
            </NuxtLink>

            <NuxtLink 
              v-if="isLoggedIn" 
              to="/events/create" 
              class="p-2.5 rounded-xl text-text-light hover:text-secondary hover:bg-secondary/10 transition-all duration-300 relative group"
              title="Create Event"
            >
              <PlusCircleIcon class="w-5 h-5" />
              <span class="absolute -bottom-10 left-1/2 -translate-x-1/2 px-3 py-1.5 bg-dark-50 text-text text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap pointer-events-none border border-surface-border shadow-lg">
                Create Event
              </span>
            </NuxtLink>

            <!-- User Menu -->
            <div v-if="isLoggedIn" class="relative ml-2">
              <Menu as="div" class="relative">
                <MenuButton class="flex items-center gap-2 p-1.5 pr-3 rounded-xl border border-surface-border bg-surface hover:bg-surface-light transition-all duration-300 group">
                  <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-primary to-accent flex items-center justify-center text-white shadow-glow-sm">
                    <UserIcon class="w-4 h-4" />
                  </div>
                  <span class="text-sm font-medium text-text hidden sm:block">Account</span>
                  <ChevronDownIcon class="w-4 h-4 text-text-light transition-transform duration-200 group-hover:text-primary" />
                </MenuButton>

                <transition
                  enter-active-class="transition duration-200 ease-out"
                  enter-from-class="transform scale-95 opacity-0 -translate-y-2"
                  enter-to-class="transform scale-100 opacity-100 translate-y-0"
                  leave-active-class="transition duration-150 ease-in"
                  leave-from-class="transform scale-100 opacity-100 translate-y-0"
                  leave-to-class="transform scale-95 opacity-0 -translate-y-2"
                >
                  <MenuItems class="absolute right-0 mt-2 w-56 origin-top-right glass-card p-1.5 focus:outline-none">
                    <div class="px-3 py-2 border-b border-surface-border mb-1">
                      <p class="text-sm font-medium text-text">Signed in as</p>
                      <p class="text-xs text-text-light truncate">{{ userEmail }}</p>
                    </div>
                    <div class="p-1">
                      <MenuItem v-slot="{ active }">
                        <NuxtLink
                          to="/profile"
                          :class="[
                            active ? 'bg-primary/10 text-primary' : 'text-text-light',
                            'group flex w-full items-center rounded-lg px-3 py-2.5 text-sm transition-colors'
                          ]"
                        >
                          <UserCircleIcon class="mr-3 h-5 w-5" aria-hidden="true" />
                          Profile
                        </NuxtLink>
                      </MenuItem>
                    </div>
                    <div class="p-1 border-t border-surface-border mt-1">
                      <MenuItem v-slot="{ active }">
                        <button
                          @click="handleLogout"
                          :class="[
                            active ? 'bg-error/10 text-error' : 'text-text-light',
                            'group flex w-full items-center rounded-lg px-3 py-2.5 text-sm transition-colors'
                          ]"
                        >
                          <ArrowRightOnRectangleIcon class="mr-3 h-5 w-5" aria-hidden="true" />
                          Sign out
                        </button>
                      </MenuItem>
                    </div>
                  </MenuItems>
                </transition>
              </Menu>
            </div>

            <NuxtLink 
              v-else 
              to="/" 
              class="btn-neon text-sm"
            >
              Sign In
            </NuxtLink>
          </div>
        </div>
      </div>
    </nav>

    <main class="relative z-10 w-full">
      <NuxtPage />
    </main>

    <!-- Footer -->
    <footer v-if="showNav" class="relative z-10 mt-16 border-t border-surface-border bg-dark-100/50 backdrop-blur-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div class="flex flex-col md:flex-row justify-between items-center gap-4">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-primary to-secondary flex items-center justify-center text-white font-bold shadow-glow-sm">
              T
            </div>
            <span class="font-display font-bold text-text-light">TicketApp</span>
          </div>
          <p class="text-sm text-text-muted">
            © {{ new Date().getFullYear() }} TicketApp. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { Menu, MenuButton, MenuItems, MenuItem } from '@headlessui/vue'
import { 
  CalendarIcon, 
  PlusCircleIcon, 
  UserIcon, 
  ChevronDownIcon, 
  UserCircleIcon, 
  ArrowRightOnRectangleIcon 
} from '@heroicons/vue/24/outline'

const route = useRoute()
const router = useRouter()
const api = useApi()

const isLoggedIn = ref(false)
const userEmail = ref('')
const showNav = computed(() => route.path !== '/')

onMounted(() => {
  const token = api.getToken()
  isLoggedIn.value = !!token
  if (token) {
    // Fetch user email from token or API
    const userData = localStorage.getItem('user_email')
    if (userData) {
      userEmail.value = userData
    }
  }
})

watch(() => route.path, () => {
  if (import.meta.client) {
    isLoggedIn.value = !!api.getToken()
  }
})

const handleLogout = () => {
  api.removeToken()
  localStorage.removeItem('user_email')
  isLoggedIn.value = false
  router.push('/')
}
</script>
