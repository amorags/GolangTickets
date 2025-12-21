<template>
  <div class="container">
    <div class="card">
      <h1>Welcome</h1>

      <div v-if="!showSignup" class="form-section">
        <h2>Login</h2>
        <form @submit.prevent="handleLogin">
          <input
            v-model="loginForm.email"
            type="email"
            placeholder="Email"
            required
          />
          <input
            v-model="loginForm.password"
            type="password"
            placeholder="Password"
            required
          />
          <button type="submit">Login</button>
        </form>
        <p class="switch-form">
          Don't have an account?
          <a href="#" @click.prevent="showSignup = true">Sign up</a>
        </p>
      </div>

      <div v-else class="form-section">
        <h2>Sign Up</h2>
        <form @submit.prevent="handleSignup">
          <input
            v-model="signupForm.username"
            type="text"
            placeholder="Username"
            required
          />
          <input
            v-model="signupForm.email"
            type="email"
            placeholder="Email"
            required
          />
          <input
            v-model="signupForm.password"
            type="password"
            placeholder="Password (min 6 chars)"
            required
          />
          <button type="submit">Sign Up</button>
        </form>
        <p class="switch-form">
          Already have an account?
          <a href="#" @click.prevent="showSignup = false">Login</a>
        </p>
      </div>

      <div v-if="message" :class="['message', 'show', messageType]">
        {{ message }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const api = useApi()
const router = useRouter()

const showSignup = ref(false)
const message = ref('')
const messageType = ref<'success' | 'error'>('success')

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

const showMessage = (text: string, type: 'success' | 'error') => {
  message.value = text
  messageType.value = type
  setTimeout(() => {
    message.value = ''
  }, 3000)
}

const handleLogin = async () => {
  try {
    await api.login(loginForm.value.email, loginForm.value.password)
    showMessage('Login successful! Redirecting...', 'success')
    setTimeout(() => {
      router.push('/events')
    }, 1000)
  } catch (error: any) {
    showMessage(error.message || 'Login failed', 'error')
  }
}

const handleSignup = async () => {
  if (signupForm.value.password.length < 6) {
    showMessage('Password must be at least 6 characters', 'error')
    return
  }

  try {
    await api.signup(
      signupForm.value.username,
      signupForm.value.email,
      signupForm.value.password
    )
    showMessage('Account created! Redirecting...', 'success')
    setTimeout(() => {
      router.push('/events')
    }, 1000)
  } catch (error: any) {
    showMessage(error.message || 'Signup failed', 'error')
  }
}
</script>