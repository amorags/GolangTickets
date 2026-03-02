<template>
  <span
    class="inline-flex items-center px-3 py-1 rounded-lg text-xs font-bold uppercase tracking-wider border transition-all duration-300"
    :class="classes"
  >
    <span 
      v-if="variant === 'success'" 
      class="w-1.5 h-1.5 rounded-full bg-success mr-2 animate-pulse"
    ></span>
    <span 
      v-if="variant === 'error'" 
      class="w-1.5 h-1.5 rounded-full bg-error mr-2"
    ></span>
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
      return 'bg-success/15 text-success border-success/30 hover:bg-success/20'
    case 'error':
      return 'bg-error/15 text-error border-error/30 hover:bg-error/20'
    case 'warning':
      return 'bg-warning/15 text-warning border-warning/30 hover:bg-warning/20'
    case 'info':
      return 'bg-secondary/15 text-secondary border-secondary/30 hover:bg-secondary/20'
    default:
      return 'bg-surface-light text-text-light border-surface-border hover:bg-surface'
  }
})
</script>
