// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: ['@nuxtjs/tailwindcss'],
  css: ['~/assets/css/tailwind.css'],
  ssr: false,
  build: {
    transpile: ['@headlessui/vue', '@heroicons/vue', 'vue-toastification'],
  },
})