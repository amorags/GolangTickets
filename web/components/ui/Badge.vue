<template>
  <span
    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border"
    :class="classes"
  >
    <slot />
  </span>
</template>

<script setup lang="ts">
const props = defineProps<{
  variant?: 'success' | 'error' | 'warning' | 'info' | 'default'
}>()

const classes = computed(() => {
  switch (props.variant) {
    case 'success':
      return 'bg-green-100 text-green-800 border-green-200'
    case 'error':
      return 'bg-red-100 text-red-800 border-red-200'
    case 'warning':
      return 'bg-yellow-100 text-yellow-800 border-yellow-200'
    case 'info':
      return 'bg-blue-100 text-blue-800 border-blue-200'
    default:
      return 'bg-gray-100 text-gray-800 border-gray-200'
  }
})
</script>