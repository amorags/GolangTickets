<template>
  <TransitionRoot appear :show="isOpen" as="template">
    <Dialog as="div" @close="close" class="relative z-50">
      <!-- Backdrop -->
      <TransitionChild
        as="template"
        enter="duration-300 ease-out"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="duration-200 ease-in"
        leave-from="opacity-100"
        leave-to="opacity-0"
      >
        <div class="fixed inset-0 bg-dark/80 backdrop-blur-md" />
      </TransitionChild>

      <div class="fixed inset-0 overflow-y-auto">
        <div class="flex min-h-full items-center justify-center p-4 text-center">
          <TransitionChild
            as="template"
            enter="duration-300 ease-out"
            enter-from="opacity-0 scale-95 -translate-y-4"
            enter-to="opacity-100 scale-100 translate-y-0"
            leave="duration-200 ease-in"
            leave-from="opacity-100 scale-100 translate-y-0"
            leave-to="opacity-0 scale-95 -translate-y-4"
          >
            <DialogPanel
              class="w-full max-w-md transform overflow-hidden rounded-2xl glass-card p-6 text-left align-middle transition-all"
            >
              <!-- Decorative top line -->
              <div class="absolute top-0 left-1/4 right-1/4 h-px bg-gradient-to-r from-transparent via-primary to-transparent"></div>
              
              <DialogTitle
                as="h3"
                class="text-lg font-display font-bold text-text"
              >
                {{ title }}
              </DialogTitle>
              <div class="mt-3">
                <p class="text-sm text-text-light leading-relaxed">
                  {{ description }}
                </p>
              </div>

              <div class="mt-6 flex justify-end gap-3">
                <slot name="footer">
                  <button
                    type="button"
                    class="btn-ghost text-sm"
                    @click="close"
                  >
                    Close
                  </button>
                </slot>
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>

<script setup lang="ts">
import {
  TransitionRoot,
  TransitionChild,
  Dialog,
  DialogPanel,
  DialogTitle,
} from '@headlessui/vue'

const props = defineProps<{
  isOpen: boolean
  title: string
  description?: string
}>()

const emit = defineEmits(['close'])

const close = () => {
  emit('close')
}
</script>
