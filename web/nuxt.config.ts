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
  app: {
    head: {
      title: 'TicketApp - Event Booking Platform',
      meta: [
        { name: 'description', content: 'Discover and book tickets for amazing events, concerts, tours, and more.' },
        { name: 'theme-color', content: '#0F0F1A' },
        { name: 'color-scheme', content: 'dark' },
      ],
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&family=Quicksand:wght@400;500;600;700&display=swap' },
      ],
    },
  },
})
