/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./components/**/*.{js,vue,ts}",
    "./layouts/**/*.vue",
    "./pages/**/*.vue",
    "./plugins/**/*.{js,ts}",
    "./app.vue",
    "./error.vue",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#FF9AA2', // Pastel Pink
          dark: '#FF858F',
          light: '#FFB7B2',
        },
        secondary: {
          DEFAULT: '#C7CEEA', // Pastel Periwinkle
          dark: '#B5BEDF',
        },
        accent: {
          DEFAULT: '#B5EAD7', // Pastel Mint
          dark: '#95D5BE',
        },
        background: {
          start: '#FFF5F5',
          end: '#F0F4FF',
        },
        text: {
          DEFAULT: '#4A4A4A',
          light: '#7A7A7A',
        },
        success: '#77DD77', // Pastel Green
        error: '#FF6961',   // Pastel Red
      },
      fontFamily: {
        sans: ['Quicksand', 'ui-sans-serif', 'system-ui', 'sans-serif'],
      },
      animation: {
        'float-in': 'floatIn 0.8s ease-out',
        'slide-up': 'slideUp 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275)',
        'fade-in': 'fadeIn 0.3s ease-in',
      },
      keyframes: {
        floatIn: {
          '0%': { opacity: '0', transform: 'translateY(20px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(10px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
      },
    },
  },
  plugins: [],
}
