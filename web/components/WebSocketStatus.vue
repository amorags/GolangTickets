<template>
  <div 
    class="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg text-xs font-medium transition-all duration-300"
    :class="isConnected ? 'bg-success/10 text-success border border-success/20' : 'bg-error/10 text-error border border-error/20'"
  >
    <span 
      class="w-2 h-2 rounded-full transition-all duration-300"
      :class="isConnected ? 'bg-success animate-pulse shadow-glow-sm' : 'bg-error'"
    ></span>
    <span>{{ isConnected ? 'Live' : 'Offline' }}</span>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  isConnected: boolean
}>()
</script>
