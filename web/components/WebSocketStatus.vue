<template>
  <div 
    class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-bold transition-all duration-300 shadow-sm border"
    :class="[
      isConnected 
        ? 'bg-accent/20 text-green-800 border-accent' 
        : 'bg-error/10 text-red-800 border-error/30'
    ]"
  >
    <span 
      class="w-2 h-2 rounded-full animate-pulse"
      :class="isConnected ? 'bg-green-500' : 'bg-red-500'"
    ></span>
    {{ isConnected ? 'Live Updates' : 'Offline' }}
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  isConnected: boolean
}>()
</script>