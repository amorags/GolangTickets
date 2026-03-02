<template>
  <div class="relative flex items-center gap-2" ref="containerRef">
    <!-- Search Bar -->
    <div class="relative group flex-1">
      <MagnifyingGlassIcon class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-text-muted group-focus-within:text-primary transition-colors" />
      <input
        v-model="filters.search"
        type="text"
        placeholder="Search events..."
        class="w-48 sm:w-64 pl-9 pr-4 py-2 text-sm rounded-xl bg-dark-50 border border-surface-border text-text placeholder-text-muted focus:border-primary focus:ring-2 focus:ring-primary/20 focus:outline-none transition-all"
        @input="debouncedSearch"
        @keydown.enter="applyFilters"
      />
    </div>

    <!-- Filter Toggle Button -->
    <button
      @click="showFilters = !showFilters"
      class="p-2 rounded-xl border border-surface-border hover:border-primary hover:text-primary hover:bg-primary/5 transition-all relative"
      :class="{ 'bg-primary/10 border-primary text-primary shadow-glow-sm': showFilters || hasActiveFilters }"
      title="Filters"
    >
      <FunnelIcon class="w-5 h-5" />
      <span 
        v-if="hasActiveFilters" 
        class="absolute -top-1 -right-1 w-3 h-3 bg-accent rounded-full border-2 border-dark animate-pulse"
      ></span>
    </button>

    <!-- Filter Dropdown -->
    <transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="transform scale-95 opacity-0 -translate-y-2"
      enter-to-class="transform scale-100 opacity-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="transform scale-100 opacity-100 translate-y-0"
      leave-to-class="transform scale-95 opacity-0 -translate-y-2"
    >
      <div
        v-if="showFilters"
        class="absolute top-full right-0 mt-3 w-80 sm:w-96 glass-card p-5 z-50"
      >
        <div class="space-y-4">
          <div class="flex justify-between items-center pb-3 border-b border-surface-border">
            <h3 class="font-bold text-text flex items-center gap-2">
              <AdjustmentsHorizontalIcon class="w-4 h-4 text-primary" />
              Filters
            </h3>
            <button @click="clearFilters" class="text-xs text-primary hover:text-primary-light font-medium transition-colors">
              Reset all
            </button>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="col-span-2">
              <label class="block text-xs font-medium text-text-muted mb-2">Event Type</label>
              <select
                v-model="filters.type"
                class="input-cyber text-sm"
                @change="applyFilters"
              >
                <option value="">All Types</option>
                <option value="concert">Concert</option>
                <option value="tour">Tour</option>
                <option value="standup">Standup</option>
                <option value="lecture">Lecture</option>
                <option value="musical">Musical</option>
                <option value="other">Other</option>
              </select>
            </div>

            <div class="col-span-2">
              <label class="block text-xs font-medium text-text-muted mb-2">City</label>
              <input
                v-model="filters.city"
                type="text"
                placeholder="Any city"
                class="input-cyber text-sm"
                @change="applyFilters"
              />
            </div>

            <div>
              <label class="block text-xs font-medium text-text-muted mb-2">From</label>
              <input
                v-model="filters.date_from"
                type="date"
                class="input-cyber text-sm"
                @change="applyFilters"
              />
            </div>

            <div>
              <label class="block text-xs font-medium text-text-muted mb-2">To</label>
              <input
                v-model="filters.date_to"
                type="date"
                class="input-cyber text-sm"
                @change="applyFilters"
              />
            </div>

            <div class="col-span-2">
              <label class="block text-xs font-medium text-text-muted mb-2">Price Range</label>
              <div class="flex items-center gap-2">
                <div class="relative flex-1">
                  <span class="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted text-xs">$</span>
                  <input
                    v-model.number="filters.price_min"
                    type="number"
                    min="0"
                    placeholder="Min"
                    class="input-cyber text-sm pl-6"
                    @change="applyFilters"
                  />
                </div>
                <span class="text-text-muted">—</span>
                <div class="relative flex-1">
                  <span class="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted text-xs">$</span>
                  <input
                    v-model.number="filters.price_max"
                    type="number"
                    min="0"
                    placeholder="Max"
                    class="input-cyber text-sm pl-6"
                    @change="applyFilters"
                  />
                </div>
              </div>
            </div>
          </div>
          
          <div class="pt-2">
            <button 
              @click="applyFilters(); showFilters = false" 
              class="btn-neon w-full py-2.5 text-sm"
            >
              Apply Filters
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { MagnifyingGlassIcon, FunnelIcon, AdjustmentsHorizontalIcon } from '@heroicons/vue/24/outline'
import { useDebounceFn, onClickOutside } from '@vueuse/core'

const router = useRouter()
const route = useRoute()

const showFilters = ref(false)
const containerRef = ref(null)

const filters = ref({
  search: '',
  type: '',
  city: '',
  date_from: '',
  date_to: '',
  price_min: undefined as number | undefined,
  price_max: undefined as number | undefined,
})

// Initialize filters from URL query
const initFilters = () => {
  filters.value.search = (route.query.search as string) || ''
  filters.value.type = (route.query.type as string) || ''
  filters.value.city = (route.query.city as string) || ''
  filters.value.date_from = (route.query.date_from as string) || ''
  filters.value.date_to = (route.query.date_to as string) || ''
  filters.value.price_min = route.query.price_min ? Number(route.query.price_min) : undefined
  filters.value.price_max = route.query.price_max ? Number(route.query.price_max) : undefined
}

onMounted(initFilters)

// Update local state when URL changes (e.g. back button)
watch(() => route.query, initFilters)

onClickOutside(containerRef, () => {
  showFilters.value = false
})

const debouncedSearch = useDebounceFn(() => {
  applyFilters()
}, 500)

const applyFilters = () => {
  const query: Record<string, any> = { ...route.query }
  
  // Update query params
  if (filters.value.search) query.search = filters.value.search; else delete query.search
  if (filters.value.type) query.type = filters.value.type; else delete query.type
  if (filters.value.city) query.city = filters.value.city; else delete query.city
  if (filters.value.date_from) query.date_from = filters.value.date_from; else delete query.date_from
  if (filters.value.date_to) query.date_to = filters.value.date_to; else delete query.date_to
  if (filters.value.price_min) query.price_min = filters.value.price_min; else delete query.price_min
  if (filters.value.price_max) query.price_max = filters.value.price_max; else delete query.price_max
  
  // Reset page on filter change
  query.page = 1

  router.push({ path: '/events', query })
}

const clearFilters = () => {
  filters.value = {
    search: '',
    type: '',
    city: '',
    date_from: '',
    date_to: '',
    price_min: undefined,
    price_max: undefined,
  }
  applyFilters()
}

const hasActiveFilters = computed(() => {
  const { search, ...rest } = filters.value
  return Object.values(rest).some(val => val !== '' && val !== undefined)
})
</script>
