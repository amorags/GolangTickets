<template>
  <div class="flex items-center justify-center gap-2 mt-10">
    <!-- Previous Button -->
    <button
      @click="$emit('update:page', page - 1)"
      :disabled="!hasPrevious"
      class="btn-ghost px-4 py-2 text-sm flex items-center gap-2"
      :class="{ 'opacity-50 cursor-not-allowed': !hasPrevious }"
    >
      <ChevronLeftIcon class="w-4 h-4" />
      Previous
    </button>

    <!-- Page Numbers -->
    <div class="flex items-center gap-1">
      <button
        v-for="p in displayedPages"
        :key="p"
        @click="p !== '...' && $emit('update:page', p)"
        :class="[
          'w-10 h-10 rounded-lg text-sm font-medium transition-all duration-300',
          p === page 
            ? 'bg-gradient-to-br from-primary to-secondary text-white shadow-glow' 
            : p === '...' 
            ? 'text-text-muted cursor-default' 
            : 'text-text-light hover:text-text hover:bg-surface-light border border-surface-border'
        ]"
      >
        {{ p }}
      </button>
    </div>

    <!-- Next Button -->
    <button
      @click="$emit('update:page', page + 1)"
      :disabled="!hasNext"
      class="btn-ghost px-4 py-2 text-sm flex items-center gap-2"
      :class="{ 'opacity-50 cursor-not-allowed': !hasNext }"
    >
      Next
      <ChevronRightIcon class="w-4 h-4" />
    </button>
  </div>

  <!-- Page Info -->
  <div class="text-center mt-4 text-sm text-text-muted">
    Page {{ page }} of {{ totalPages }} • {{ total }} total events
  </div>
</template>

<script setup lang="ts">
import { ChevronLeftIcon, ChevronRightIcon } from '@heroicons/vue/24/outline'

const props = defineProps<{
  page: number
  totalPages: number
  total: number
  limit: number
  hasNext: boolean
  hasPrevious: boolean
}>()

defineEmits<{
  'update:page': [page: number]
}>()

const displayedPages = computed(() => {
  const pages: (number | string)[] = []
  const { page, totalPages } = props
  
  if (totalPages <= 7) {
    for (let i = 1; i <= totalPages; i++) {
      pages.push(i)
    }
  } else {
    if (page <= 3) {
      pages.push(1, 2, 3, 4, '...', totalPages)
    } else if (page >= totalPages - 2) {
      pages.push(1, '...', totalPages - 3, totalPages - 2, totalPages - 1, totalPages)
    } else {
      pages.push(1, '...', page - 1, page, page + 1, '...', totalPages)
    }
  }
  
  return pages
})
</script>
