<template>
  <div class="min-h-screen flex items-center justify-center p-4 relative overflow-hidden">
    <!-- Animated background elements -->
    <div class="absolute inset-0 overflow-hidden">
      <div class="absolute top-1/4 left-1/4 w-96 h-96 bg-primary/20 rounded-full blur-3xl animate-float"></div>
      <div class="absolute bottom-1/4 right-1/4 w-80 h-80 bg-secondary/20 rounded-full blur-3xl animate-float" style="animation-delay: -3s;"></div>
      <div class="absolute top-1/2 right-1/3 w-64 h-64 bg-accent/20 rounded-full blur-3xl animate-float" style="animation-delay: -5s;"></div>
    </div>

    <!-- Grid pattern overlay -->
    <div class="absolute inset-0 bg-cy-grid opacity-50"></div>

    <!-- Main card -->
    <div class="relative w-full max-w-md">
      <!-- Glow effect behind card -->
      <div class="absolute inset-0 bg-gradient-to-r from-primary/30 to-secondary/30 rounded-3xl blur-2xl transform scale-105"></div>
      
      <div class="relative glass-card p-8 animate-float-in">
        <!-- Logo & Header -->
        <div class="text-center mb-8">
          <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-gradient-to-br from-primary to-secondary mb-4 shadow-glow">
            <span class="text-3xl font-bold text-white">T</span>
          </div>
          <h1 class="text-3xl font-display font-bold text-gradient mb-2">
            {{ showSignup ? 'Create Account' : 'Welcome Back' }}
          </h1>
          <p class="text-text-light">
            {{ showSignup ? 'Join the experience today' : 'Sign in to continue' }}
          </p>
        </div>

        <!-- Login Form -->
        <div v-if="!showSignup" class="space-y-5 animate-fade-in">
          <form @submit.prevent="handleLogin" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-text-light mb-2">Email</label>
              <div class="relative">
                <EnvelopeIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-text-muted" />
                <input
                  v-model="loginForm.email"
                  type="email"
                  class="input-cyber pl-11"
                  placeholder="hello@example.com"
                  required
                />
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-text-light mb-2">Password</label>
              <div class="relative">
                <LockClosedIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-text-muted" />
                <input
                  v-model="loginForm.password"
                  type="password"
                  class="input-cyber pl-11"
                  placeholder="••••••••"
                  required
                />
              </div>
            </div>
            <button 
              type="submit"
              class="btn-neon w-full py-4 text-base mt-6"
            >
              Sign In
            </button>
          </form>
          
          <p class="text-center text-sm text-text-light">
            Don't have an account?
            <a href="#" @click.prevent="showSignup = true" class="font-bold text-primary hover:text-primary-light transition-colors ml-1">
              Sign up
            </a>
          </p>
        </div>

        <!-- Signup Form -->
        <div v-else class="space-y-5 animate-fade-in">
          <form @submit.prevent="handleSignup" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-text-light mb-2">Username</label>
              <div class="relative">
                <UserIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-text-muted" />
                <input
                  v-model="signupForm.username"
                  type="text"
                  class="input-cyber pl-11"
                  placeholder="johndoe"
                  required
                />
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-text-light mb-2">Email</label>
              <div class="relative">
                <EnvelopeIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-text-muted" />
                <input
                  v-model="signupForm.email"
                  type="email"
                  class="input-cyber pl-11"
                  placeholder="hello@example.com"
                  required
                />
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-text-light mb-2">Password</label>
              <div class="relative">
                <LockClosedIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-text-muted" />
                <input
                  v-model="signupForm.password"
                  type="password"
                  class="input-cyber pl-11"
                  placeholder="Min 6 characters"
                  required
                />
              </div>
            </div>
            <button 
              type="submit"
              class="btn-neon w-full py-4 text-base mt-6"
            >
              Create Account
            </button>
          </form>
          
          <p class="text-center text-sm text-text-light">
            Already have an account?
            <a href="#" @click.prevent="showSignup = false" class="font-bold text-primary hover:text-primary-light transition-colors ml-1">
              Sign in
            </a>
          </p>
        </div>

        <!-- Decorative elements -->
        <div class="absolute -top-px left-1/4 right-1/4 h-px bg-gradient-to-r from-transparent via-primary to-transparent"></div>
        <div class="absolute -bottom-px left-1/4 right-1/4 h-px bg-gradient-to-r from-transparent via-secondary to-transparent"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useToast } from "vue-toastification"
import { EnvelopeIcon, LockClosedIcon, UserIcon } from '@heroicons/vue/24/outline'

const api = useApi()
const router = useRouter()
const toast = useToast()

const showSignup = ref(false)

const loginForm = ref({
  email: '',
  password: '',
})

const signupForm = ref({
  username: '',
  email: '',
  password: '',
})

onMounted(() => {
  if (api.getToken()) {
    router.push('/events')
  }
})

const handleLogin = async () => {
  try {
    await api.login(loginForm.value.email, loginForm.value.password)
    localStorage.setItem('user_email', loginForm.value.email)
    toast.success('Welcome back!')
    router.push('/events')
  } catch (error: any) {
    toast.error(error.message || 'Login failed')
  }
}

const handleSignup = async () => {
  if (signupForm.value.password.length < 6) {
    toast.error('Password must be at least 6 characters')
    return
  }

  try {
    await api.signup(
      signupForm.value.username,
      signupForm.value.email,
      signupForm.value.password
    )
    localStorage.setItem('user_email', signupForm.value.email)
    toast.success('Account created successfully!')
    router.push('/events')
  } catch (error: any) {
    toast.error(error.message || 'Signup failed')
  }
}
</script>
