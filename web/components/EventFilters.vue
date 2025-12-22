<template>
  <div class="bg-white rounded-3xl shadow-lg border border-gray-100 p-6 mb-8">
    <div class="mb-6 flex gap-4">
      <div class="relative flex-1">
        <MagnifyingGlassIcon class="absolute left-4 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
        <input
          v-model="localFilters.search"
          type="text"
          placeholder="Search events, venues, or cities..."
          class="w-full pl-12 pr-4 py-3 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
          @input="debouncedSearch"
        />
      </div>
      <button
        @click="showFilters = !showFilters"
        class="flex items-center px-6 py-3 rounded-xl border border-gray-200 hover:border-primary hover:text-primary transition-all font-medium"
        :class="showFilters ? 'bg-primary/5 border-primary text-primary' : 'bg-gray-50 text-gray-600'"
      >
        <FunnelIcon class="w-5 h-5 mr-2" />
        Filters
      </button>
    </div>

    <div v-show="showFilters" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 transition-all duration-300 ease-in-out">
      <div>
        <label class="block text-xs font-medium text-gray-500 mb-2">Event Type</label>
        <select
          v-model="localFilters.type"
          class="w-full px-4 py-2.5 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all cursor-pointer"
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

      <div>
        <label class="block text-xs font-medium text-gray-500 mb-2">City</label>
        <input
          v-model="localFilters.city"
          type="text"
          placeholder="Filter by city..."
          class="w-full px-4 py-2.5 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
          @input="debouncedSearch"
        />
      </div>

      <div>
        <label class="block text-xs font-medium text-gray-500 mb-2">From Date</label>
        <input
          v-model="localFilters.date_from"
          type="date"
          class="w-full px-4 py-2.5 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
          @change="applyFilters"
        />
      </div>

      <div>
        <label class="block text-xs font-medium text-gray-500 mb-2">To Date</label>
        <input
          v-model="localFilters.date_to"
          type="date"
          class="w-full px-4 py-2.5 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
          @change="applyFilters"
        />
      </div>

      <div class="md:col-span-2">
        <label class="block text-xs font-medium text-gray-500 mb-2">Price Range</label>
        <div class="grid grid-cols-2 gap-3">
          <div class="relative">
            <span class="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400">$</span>
            <input
              v-model.number="localFilters.price_min"
              type="number"
              min="0"
              step="0.01"
              placeholder="Min"
              class="w-full pl-8 pr-4 py-2.5 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
              @change="applyFilters"
            />
          </div>
          <div class="relative">
            <span class="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-400">$</span>
            <input
              v-model.number="localFilters.price_max"
              type="number"
              min="0"
              step="0.01"
              placeholder="Max"
              class="w-full pl-8 pr-4 py-2.5 rounded-xl bg-gray-50 border border-gray-200 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
              @change="applyFilters"
            />
          </div>
        </div>
      </div>

      <div class="md:col-span-2 flex items-end">
        <button
          @click="clearFilters"
          class="w-full px-4 py-2.5 rounded-xl bg-gray-100 text-gray-600 font-medium hover:bg-gray-200 transition-colors"
        >
          Clear Filters
        </button>
      </div>
    </div>

    <div v-if="hasActiveFilters" class="mt-4 flex flex-wrap gap-2">
      <UiBadge v-if="localFilters.search" variant="info">
        Search: {{ localFilters.search }}
        <button @click="removeFilter('search')" class="ml-1.5 hover:text-blue-900">×</button>
      </UiBadge>
      <UiBadge v-if="localFilters.type" variant="info">
        Type: {{ localFilters.type }}
        <button @click="removeFilter('type')" class="ml-1.5 hover:text-blue-900">×</button>
      </UiBadge>
      <UiBadge v-if="localFilters.city" variant="info">
        City: {{ localFilters.city }}
        <button @click="removeFilter('city')" class="ml-1.5 hover:text-blue-900">×</button>
      </UiBadge>
      <UiBadge v-if="localFilters.date_from" variant="info">
        From: {{ localFilters.date_from }}
        <button @click="removeFilter('date_from')" class="ml-1.5 hover:text-blue-900">×</button>
      </UiBadge>
      <UiBadge v-if="localFilters.date_to" variant="info">
        To: {{ localFilters.date_to }}
        <button @click="removeFilter('date_to')" class="ml-1.5 hover:text-blue-900">×</button>
      </UiBadge>
      <UiBadge v-if="localFilters.price_min !== undefined" variant="info">
        Min: ${{ localFilters.price_min }}
        <button @click="removeFilter('price_min')" class="ml-1.5 hover:text-blue-900">×</button>
      </UiBadge>
      <UiBadge v-if="localFilters.price_max !== undefined" variant="info">
        Max: ${{ localFilters.price_max }}
        <button @click="removeFilter('price_max')" class="ml-1.5 hover:text-blue-900">×</button>
      </UiBadge>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { EventFilters } from '~/types'
import { MagnifyingGlassIcon, FunnelIcon } from '@heroicons/vue/24/outline'
import { useDebounceFn } from '@vueuse/core'

const props = defineProps<{
  modelValue: EventFilters
}>()

const emit = defineEmits<{
  'update:modelValue': [filters: EventFilters]
}>()

const showFilters = ref(false)
const localFilters = ref<EventFilters>({ ...props.modelValue })

const debouncedSearch = useDebounceFn(() => {
  applyFilters()
}, 500)

const applyFilters = () => {
  emit('update:modelValue', { ...localFilters.value })
}

const clearFilters = () => {
  localFilters.value = {}
  applyFilters()
}

const removeFilter = (key: keyof EventFilters) => {
  delete localFilters.value[key]
  applyFilters()
}

const hasActiveFilters = computed(() => {
  return Object.keys(localFilters.value).some(
    key => key !== 'page' && key !== 'limit' && localFilters.value[key as keyof EventFilters]
  )
})

watch(() => props.modelValue, (newValue) => {
  localFilters.value = { ...newValue }
}, { deep: true })
</script>
