<template>
  <div class="min-h-screen bg-gradient-to-br from-background-start to-background-end font-sans text-text">
    <NuxtRouteAnnouncer />
    
    <nav v-if="showNav" class="sticky top-0 z-40 bg-white/80 backdrop-blur-md border-b border-white/20 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="flex flex-shrink-0 items-center">
            <NuxtLink to="/events" class="flex items-center gap-2 group">
              <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-primary to-primary-light flex items-center justify-center text-white font-bold text-lg shadow-md group-hover:scale-105 transition-transform">
                T
              </div>
              <span class="font-bold text-xl bg-clip-text text-transparent bg-gradient-to-r from-primary to-primary-dark hidden sm:block">
                TicketApp
              </span>
            </NuxtLink>
          </div>

          <div class="hidden md:flex flex-1 justify-center px-8">
            <NavBarSearch />
          </div>
          
          <div class="flex items-center gap-4">
            <NuxtLink 
              to="/events" 
              class="p-2 rounded-xl text-text-light hover:text-primary hover:bg-primary/10 transition-colors relative group"
              title="Events"
            >
              <CalendarIcon class="w-6 h-6" />
              <span class="absolute -bottom-8 left-1/2 -translate-x-1/2 px-2 py-1 bg-gray-800 text-white text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap pointer-events-none">
                Browse Events
              </span>
            </NuxtLink>

            <NuxtLink 
              v-if="isLoggedIn" 
              to="/events/create" 
              class="p-2 rounded-xl text-text-light hover:text-primary hover:bg-primary/10 transition-colors relative group"
              title="Create Event"
            >
              <PlusCircleIcon class="w-6 h-6" />
               <span class="absolute -bottom-8 left-1/2 -translate-x-1/2 px-2 py-1 bg-gray-800 text-white text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap pointer-events-none">
                Create Event
              </span>
            </NuxtLink>

            <div v-if="isLoggedIn" class="relative ml-2">
              <Menu as="div" class="relative">
                <MenuButton class="flex items-center gap-2 p-1 pr-3 rounded-full border border-gray-200 bg-white hover:shadow-md transition-shadow">
                  <div class="w-8 h-8 rounded-full bg-secondary/30 flex items-center justify-center text-secondary-dark">
                    <UserIcon class="w-5 h-5" />
                  </div>
                  <span class="text-sm font-medium text-text hidden sm:block">My Account</span>
                  <ChevronDownIcon class="w-4 h-4 text-text-light" />
                </MenuButton>

                <transition
                  enter-active-class="transition duration-100 ease-out"
                  enter-from-class="transform scale-95 opacity-0"
                  enter-to-class="transform scale-100 opacity-100"
                  leave-active-class="transition duration-75 ease-in"
                  leave-from-class="transform scale-100 opacity-100"
                  leave-to-class="transform scale-95 opacity-0"
                >
                  <MenuItems class="absolute right-0 mt-2 w-48 origin-top-right divide-y divide-gray-100 rounded-xl bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                    <div class="p-1">
                      <MenuItem v-slot="{ active }">
                        <NuxtLink
                          to="/profile"
                          :class="[
                            active ? 'bg-primary/10 text-primary' : 'text-text',
                            'group flex w-full items-center rounded-lg px-2 py-2 text-sm'
                          ]"
                        >
                          <UserCircleIcon class="mr-2 h-5 w-5" aria-hidden="true" />
                          Profile
                        </NuxtLink>
                      </MenuItem>
                    </div>
                    <div class="p-1">
                      <MenuItem v-slot="{ active }">
                        <button
                          @click="handleLogout"
                          :class="[
                            active ? 'bg-error/10 text-error' : 'text-text',
                            'group flex w-full items-center rounded-lg px-2 py-2 text-sm'
                          ]"
                        >
                          <ArrowRightOnRectangleIcon class="mr-2 h-5 w-5" aria-hidden="true" />
                          Logout
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
              class="px-4 py-2 rounded-xl bg-primary text-white font-bold shadow-md shadow-primary/30 hover:shadow-lg hover:-translate-y-0.5 transition-all text-sm"
            >
              Login
            </NuxtLink>
          </div>
        </div>
      </div>
    </nav>

    <main class="w-full">
      <NuxtPage />
    </main>
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