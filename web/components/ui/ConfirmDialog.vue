<template>
  <BaseModal :is-open="isOpen" :title="title" :description="description" @close="cancel">
    <template #footer>
      <button
        type="button"
        class="btn-ghost text-sm"
        @click="cancel"
      >
        Cancel
      </button>
      <button
        type="button"
        class="text-sm font-bold px-5 py-2.5 rounded-xl transition-all duration-300"
        :class="isDanger 
          ? 'bg-error/20 text-error border border-error/30 hover:bg-error/30 hover:shadow-glow-pink' 
          : 'btn-neon'"
        @click="confirm"
      >
        {{ confirmText }}
      </button>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import BaseModal from './BaseModal.vue'

defineProps<{
  isOpen: boolean
  title: string
  description: string
  confirmText?: string
  isDanger?: boolean
}>()

const emit = defineEmits(['close', 'confirm'])

const cancel = () => emit('close')
const confirm = () => emit('confirm')
</script>
