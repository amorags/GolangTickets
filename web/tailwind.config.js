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
        // Dark cyberpunk theme
        primary: {
          DEFAULT: '#8B5CF6', // Vivid purple
          dark: '#7C3AED',
          light: '#A78BFA',
          glow: 'rgba(139, 92, 246, 0.5)',
        },
        secondary: {
          DEFAULT: '#06B6D4', // Cyan
          dark: '#0891B2',
          light: '#22D3EE',
          glow: 'rgba(6, 182, 212, 0.5)',
        },
        accent: {
          DEFAULT: '#F472B6', // Pink
          dark: '#EC4899',
          light: '#F9A8D4',
          glow: 'rgba(244, 114, 182, 0.5)',
        },
        success: {
          DEFAULT: '#10B981', // Emerald
          glow: 'rgba(16, 185, 129, 0.5)',
        },
        warning: {
          DEFAULT: '#F59E0B', // Amber
          glow: 'rgba(245, 158, 11, 0.5)',
        },
        error: {
          DEFAULT: '#EF4444', // Red
          glow: 'rgba(239, 68, 68, 0.5)',
        },
        // Dark backgrounds
        dark: {
          DEFAULT: '#0F0F1A',
          50: '#1A1A2E',
          100: '#16162A',
          200: '#12122A',
          300: '#0F0F23',
          400: '#0A0A1A',
          500: '#070714',
        },
        // Surface colors for cards/modals
        surface: {
          DEFAULT: 'rgba(26, 26, 46, 0.8)',
          light: 'rgba(42, 42, 68, 0.8)',
          border: 'rgba(139, 92, 246, 0.2)',
          hover: 'rgba(139, 92, 246, 0.1)',
        },
        // Text colors
        text: {
          DEFAULT: '#F8FAFC',
          light: '#94A3B8',
          muted: '#64748B',
        },
      },
      fontFamily: {
        sans: ['Inter', 'Quicksand', 'ui-sans-serif', 'system-ui', 'sans-serif'],
        display: ['Quicksand', 'Inter', 'sans-serif'],
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'gradient-conic': 'conic-gradient(from 180deg at 50% 50%, var(--tw-gradient-stops))',
        'mesh-gradient': 'linear-gradient(135deg, #0F0F1A 0%, #1A1A2E 25%, #0F0F1A 50%, #16162A 75%, #0F0F1A 100%)',
        'cyber-grid': `
          linear-gradient(rgba(139, 92, 246, 0.03) 1px, transparent 1px),
          linear-gradient(90deg, rgba(139, 92, 246, 0.03) 1px, transparent 1px)
        `,
        'glow-blob': 'radial-gradient(circle at 50% 50%, rgba(139, 92, 246, 0.15) 0%, transparent 70%)',
      },
      animation: {
        'float-in': 'floatIn 0.8s ease-out',
        'slide-up': 'slideUp 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275)',
        'fade-in': 'fadeIn 0.3s ease-in',
        'glow-pulse': 'glowPulse 2s ease-in-out infinite',
        'shimmer': 'shimmer 2s linear infinite',
        'float': 'float 6s ease-in-out infinite',
        'pulse-slow': 'pulse 4s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'gradient-x': 'gradientX 15s ease infinite',
        'gradient-y': 'gradientY 15s ease infinite',
        'gradient-xy': 'gradientXY 15s ease infinite',
        'border-glow': 'borderGlow 3s linear infinite',
        'scan-line': 'scanLine 8s linear infinite',
      },
      keyframes: {
        floatIn: {
          '0%': { opacity: '0', transform: 'translateY(30px) scale(0.95)' },
          '100%': { opacity: '1', transform: 'translateY(0) scale(1)' },
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(20px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        glowPulse: {
          '0%, 100%': { boxShadow: '0 0 20px rgba(139, 92, 246, 0.3)' },
          '50%': { boxShadow: '0 0 40px rgba(139, 92, 246, 0.6)' },
        },
        shimmer: {
          '0%': { backgroundPosition: '-200% 0' },
          '100%': { backgroundPosition: '200% 0' },
        },
        float: {
          '0%, 100%': { transform: 'translateY(0)' },
          '50%': { transform: 'translateY(-20px)' },
        },
        gradientX: {
          '0%, 100%': { backgroundPosition: '0% 50%' },
          '50%': { backgroundPosition: '100% 50%' },
        },
        gradientY: {
          '0%, 100%': { backgroundPosition: '50% 0%' },
          '50%': { backgroundPosition: '50% 100%' },
        },
        gradientXY: {
          '0%, 100%': { backgroundPosition: '0% 0%' },
          '25%': { backgroundPosition: '100% 0%' },
          '50%': { backgroundPosition: '100% 100%' },
          '75%': { backgroundPosition: '0% 100%' },
        },
        borderGlow: {
          '0%': { borderColor: 'rgba(139, 92, 246, 0.5)' },
          '50%': { borderColor: 'rgba(6, 182, 212, 0.5)' },
          '100%': { borderColor: 'rgba(139, 92, 246, 0.5)' },
        },
        scanLine: {
          '0%': { transform: 'translateY(-100%)' },
          '100%': { transform: 'translateY(100%)' },
        },
      },
      boxShadow: {
        'glow-sm': '0 0 15px -3px rgba(139, 92, 246, 0.3)',
        'glow': '0 0 30px -5px rgba(139, 92, 246, 0.4)',
        'glow-lg': '0 0 50px -10px rgba(139, 92, 246, 0.5)',
        'glow-cyan': '0 0 30px -5px rgba(6, 182, 212, 0.4)',
        'glow-pink': '0 0 30px -5px rgba(244, 114, 182, 0.4)',
        'inner-glow': 'inset 0 0 30px rgba(139, 92, 246, 0.1)',
      },
      backdropBlur: {
        xs: '2px',
      },
    },
  },
  plugins: [],
}
