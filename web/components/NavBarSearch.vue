<template>
  <div class="relative flex items-center gap-2" ref="containerRef">
    <!-- Search Bar -->
    <div class="relative group">
      <MagnifyingGlassIcon class="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400 group-focus-within:text-primary transition-colors" />
      <input
        v-model="filters.search"
        type="text"
        placeholder="Search..."
        class="w-48 sm:w-64 pl-9 pr-4 py-2 text-sm rounded-full bg-gray-100 border-none focus:bg-white focus:ring-2 focus:ring-primary/20 transition-all shadow-inner"
        @input="debouncedSearch"
        @keydown.enter="applyFilters"
      />
    </div>

    <!-- Filter Toggle Button -->
    <button
      @click="showFilters = !showFilters"
      class="p-2 rounded-full border border-gray-200 hover:border-primary hover:text-primary hover:bg-primary/5 transition-all relative"
      :class="{ 'bg-primary/5 border-primary text-primary': showFilters || hasActiveFilters }"
      title="Filters"
    >
      <FunnelIcon class="w-5 h-5" />
      <span v-if="hasActiveFilters" class="absolute top-0 right-0 w-2.5 h-2.5 bg-primary rounded-full border-2 border-white"></span>
    </button>

    <!-- Filter Dropdown -->
    <transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="transform scale-95 opacity-0"
      enter-to-class="transform scale-100 opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="transform scale-100 opacity-100"
      leave-to-class="transform scale-95 opacity-0"
    >
      <div
        v-if="showFilters"
        class="absolute top-full right-0 mt-3 w-80 sm:w-96 bg-white rounded-2xl shadow-xl border border-gray-100 p-5 z-50 origin-top-right"
      >
        <div class="space-y-4">
          <div class="flex justify-between items-center pb-2 border-b border-gray-100">
            <h3 class="font-bold text-gray-900">Filters</h3>
            <button @click="clearFilters" class="text-xs text-primary hover:text-primary-dark font-medium">
              Reset all
            </button>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="col-span-2">
              <label class="block text-xs font-medium text-gray-500 mb-1.5">Event Type</label>
              <select
                v-model="filters.type"
                class="w-full px-3 py-2 text-sm rounded-lg bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"
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
              <label class="block text-xs font-medium text-gray-500 mb-1.5">City</label>
              <input
                v-model="filters.city"
                type="text"
                placeholder="Any city"
                class="w-full px-3 py-2 text-sm rounded-lg bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"
                @change="applyFilters"
              />
            </div>

            <div>
              <label class="block text-xs font-medium text-gray-500 mb-1.5">From</label>
              <input
                v-model="filters.date_from"
                type="date"
                class="w-full px-3 py-2 text-sm rounded-lg bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"
                @change="applyFilters"
              />
            </div>

            <div>
              <label class="block text-xs font-medium text-gray-500 mb-1.5">To</label>
              <input
                v-model="filters.date_to"
                type="date"
                class="w-full px-3 py-2 text-sm rounded-lg bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"
                @change="applyFilters"
              />
            </div>

            <div class="col-span-2">
              <label class="block text-xs font-medium text-gray-500 mb-1.5">Price Range</label>
              <div class="flex items-center gap-2">
                <div class="relative flex-1">
                  <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 text-xs">$</span>
                  <input
                    v-model.number="filters.price_min"
                    type="number"
                    min="0"
                    placeholder="Min"
                    class="w-full pl-6 pr-2 py-2 text-sm rounded-lg bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"
                    @change="applyFilters"
                  />
                </div>
                <span class="text-gray-400">-</span>
                <div class="relative flex-1">
                  <span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 text-xs">$</span>
                  <input
                    v-model.number="filters.price_max"
                    type="number"
                    min="0"
                    placeholder="Max"
                    class="w-full pl-6 pr-2 py-2 text-sm rounded-lg bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none"
                    @change="applyFilters"
                  />
                </div>
              </div>
            </div>
          </div>
          
          <div class="pt-2">
            <button 
              @click="applyFilters(); showFilters = false" 
              class="w-full py-2 bg-primary text-white rounded-lg font-medium text-sm hover:bg-primary-dark transition-colors"
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
import { MagnifyingGlassIcon, FunnelIcon } from '@heroicons/vue/24/outline'
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