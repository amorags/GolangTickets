<template>
  <div v-if="totalPages > 1" class="flex items-center justify-center gap-2 mt-8">
    <button
      :disabled="!hasPrevious"
      @click="goToPage(page - 1)"
      class="px-4 py-2 rounded-xl font-medium transition-all disabled:opacity-30 disabled:cursor-not-allowed bg-white border border-gray-200 hover:bg-gray-50"
    >
      Previous
    </button>

    <div class="flex gap-1">
      <button
        v-for="pageNum in visiblePages"
        :key="pageNum"
        @click="goToPage(pageNum)"
        :class="[
          'px-4 py-2 rounded-xl font-medium transition-all',
          pageNum === page
            ? 'bg-primary text-white shadow-lg shadow-primary/30'
            : 'bg-white border border-gray-200 hover:bg-gray-50'
        ]"
      >
        {{ pageNum }}
      </button>
    </div>

    <button
      :disabled="!hasNext"
      @click="goToPage(page + 1)"
      class="px-4 py-2 rounded-xl font-medium transition-all disabled:opacity-30 disabled:cursor-not-allowed bg-white border border-gray-200 hover:bg-gray-50"
    >
      Next
    </button>
  </div>

  <div class="text-center text-sm text-gray-500 mt-4">
    Showing {{ startItem }}-{{ endItem }} of {{ total }} events
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  page: number
  totalPages: number
  total: number
  limit: number
  hasNext: boolean
  hasPrevious: boolean
}>()

const emit = defineEmits<{
  'update:page': [page: number]
}>()

const goToPage = (pageNum: number) => {
  emit('update:page', pageNum)
  if (import.meta.client) {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }
}

const visiblePages = computed(() => {
  const pages = []
  let start = Math.max(1, props.page - 2)
  let end = Math.min(props.totalPages, start + 4)

  if (end - start < 4) {
    start = Math.max(1, end - 4)
  }

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }

  return pages
})

const startItem = computed(() => {
  return (props.page - 1) * props.limit + 1
})

const endItem = computed(() => {
  return Math.min(props.page * props.limit, props.total)
})
</script>
