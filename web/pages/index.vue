<template>
  <div class="min-h-screen flex items-center justify-center p-4">
    <div class="w-full max-w-md bg-white/80 backdrop-blur-xl rounded-3xl shadow-2xl border border-white/50 p-8 animate-float-in">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-text mb-2">Welcome Back</h1>
        <p class="text-text-light">Manage your event tickets with ease</p>
      </div>

      <div v-if="!showSignup" class="space-y-6 animate-fade-in">
        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-text-light mb-1">Email</label>
            <input
              v-model="loginForm.email"
              type="email"
              class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-100 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
              placeholder="hello@example.com"
              required
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-text-light mb-1">Password</label>
            <input
              v-model="loginForm.password"
              type="password"
              class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-100 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
              placeholder="••••••••"
              required
            />
          </div>
          <button 
            type="submit"
            class="w-full py-3.5 rounded-xl bg-gradient-to-r from-primary to-primary-dark text-white font-bold shadow-lg shadow-primary/30 hover:shadow-xl hover:-translate-y-0.5 transition-all"
          >
            Sign In
          </button>
        </form>
        
        <p class="text-center text-sm text-text-light">
          Don't have an account?
          <a href="#" @click.prevent="showSignup = true" class="font-bold text-primary hover:text-primary-dark transition-colors">Sign up</a>
        </p>
      </div>

      <div v-else class="space-y-6 animate-fade-in">
        <form @submit.prevent="handleSignup" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-text-light mb-1">Username</label>
            <input
              v-model="signupForm.username"
              type="text"
              class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-100 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
              placeholder="johndoe"
              required
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-text-light mb-1">Email</label>
            <input
              v-model="signupForm.email"
              type="email"
              class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-100 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
              placeholder="hello@example.com"
              required
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-text-light mb-1">Password</label>
            <input
              v-model="signupForm.password"
              type="password"
              class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-100 focus:border-primary focus:ring-2 focus:ring-primary/20 outline-none transition-all"
              placeholder="Min 6 chars"
              required
            />
          </div>
          <button 
            type="submit"
            class="w-full py-3.5 rounded-xl bg-gradient-to-r from-primary to-primary-dark text-white font-bold shadow-lg shadow-primary/30 hover:shadow-xl hover:-translate-y-0.5 transition-all"
          >
            Create Account
          </button>
        </form>
        
        <p class="text-center text-sm text-text-light">
          Already have an account?
          <a href="#" @click.prevent="showSignup = false" class="font-bold text-primary hover:text-primary-dark transition-colors">Login</a>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useToast } from "vue-toastification";

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
    toast.success('Account created successfully!')
    router.push('/events')
  } catch (error: any) {
    toast.error(error.message || 'Signup failed')
  }
}
</script>