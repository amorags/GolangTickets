<template>
  <BaseModal :is-open="isOpen" :title="title" :description="description" @close="cancel">
    <template #footer>
      <button
        type="button"
        class="inline-flex justify-center rounded-xl border border-transparent bg-gray-100 px-4 py-2 text-sm font-medium text-gray-900 hover:bg-gray-200 focus:outline-none focus-visible:ring-2 focus-visible:ring-gray-500 focus-visible:ring-offset-2 transition-colors"
        @click="cancel"
      >
        Cancel
      </button>
      <button
        type="button"
        class="inline-flex justify-center rounded-xl border border-transparent px-4 py-2 text-sm font-medium text-white shadow-sm focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 transition-colors"
        :class="isDanger ? 'bg-error hover:bg-error/90 focus-visible:ring-error' : 'bg-primary hover:bg-primary-dark focus-visible:ring-primary'"
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