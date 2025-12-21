<template>
  <div class="ws-status" :class="statusClass">
    <span class="ws-status-dot"></span>
    <span class="ws-status-text">{{ statusText }}</span>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  isConnected: boolean
}>()

const statusClass = computed(() => {
  return props.isConnected ? 'connected' : 'disconnected'
})

const statusText = computed(() => {
  return props.isConnected ? 'Live' : 'Offline'
})
</script>

<style scoped>
.ws-status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.ws-status.connected {
  background-color: #d1fae5;
  color: #065f46;
}

.ws-status.disconnected {
  background-color: #fee2e2;
  color: #991b1b;
}

.ws-status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: pulse 2s ease-in-out infinite;
}

.ws-status.connected .ws-status-dot {
  background-color: #10b981;
}

.ws-status.disconnected .ws-status-dot {
  background-color: #ef4444;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}
</style>
