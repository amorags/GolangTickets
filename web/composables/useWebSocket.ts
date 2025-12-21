interface WebSocketMessage {
  type: string
  event_id?: number
  timestamp?: string
  data?: any
}

interface AvailabilityUpdate {
  event_id: number
  available_tickets: number
  capacity: number
  last_updated: string
}

type EventHandler = (data: any) => void

export const useWebSocket = () => {
  const socket = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const reconnectAttempts = ref(0)
  const maxReconnectAttempts = 5
  const reconnectDelay = ref(1000)
  const lastToken = ref('')

  // Event handlers map
  const eventHandlers = ref<Map<string, Set<EventHandler>>>(new Map())

  // Subscribe to event updates
  const subscribe = (eventId: number) => {
    if (socket.value && isConnected.value) {
      socket.value.send(JSON.stringify({
        type: 'subscribe',
        event_id: eventId
      }))
    }
  }

  // Unsubscribe from event updates
  const unsubscribe = (eventId: number) => {
    if (socket.value && isConnected.value) {
      socket.value.send(JSON.stringify({
        type: 'unsubscribe',
        event_id: eventId
      }))
    }
  }

  // Add event listener
  const on = (eventType: string, handler: EventHandler) => {
    if (!eventHandlers.value.has(eventType)) {
      eventHandlers.value.set(eventType, new Set())
    }
    eventHandlers.value.get(eventType)?.add(handler)
  }

  // Remove event listener
  const off = (eventType: string, handler: EventHandler) => {
    eventHandlers.value.get(eventType)?.delete(handler)
  }

  // Emit event to handlers
  const emit = (eventType: string, data: any) => {
    eventHandlers.value.get(eventType)?.forEach(handler => handler(data))
  }

  const connect = (token: string) => {
    if (socket.value) return

    lastToken.value = token
    const wsUrl = `ws://localhost:8080/ws?token=${encodeURIComponent(token)}`
    socket.value = new WebSocket(wsUrl)

    socket.value.onopen = () => {
      isConnected.value = true
      reconnectAttempts.value = 0
      reconnectDelay.value = 1000
      emit('connected', null)
      console.log('WebSocket connected')
    }

    socket.value.onmessage = (event) => {
      try {
        const message: WebSocketMessage = JSON.parse(event.data)

        if (message.type === 'availability_update') {
          emit('availability_update', message.data as AvailabilityUpdate)
        } else if (message.type === 'connection_ack') {
          console.log('WebSocket connection acknowledged:', message.data)
        } else if (message.type === 'error') {
          emit('error', message.data)
          console.error('WebSocket error:', message.data)
        }
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error)
      }
    }

    socket.value.onerror = (error) => {
      emit('error', error)
      console.error('WebSocket error:', error)
    }

    socket.value.onclose = () => {
      isConnected.value = false
      socket.value = null
      emit('disconnected', null)
      console.log('WebSocket disconnected')

      // Attempt reconnection
      if (reconnectAttempts.value < maxReconnectAttempts) {
        setTimeout(() => {
          reconnectAttempts.value++
          reconnectDelay.value *= 2 // Exponential backoff
          console.log(`Reconnecting... Attempt ${reconnectAttempts.value}`)
          connect(lastToken.value)
        }, reconnectDelay.value)
      } else {
        console.log('Max reconnection attempts reached')
        emit('max_reconnect_attempts_reached', null)
      }
    }
  }

  const disconnect = () => {
    if (socket.value) {
      socket.value.close()
      socket.value = null
      isConnected.value = false
    }
  }

  // Cleanup on unmount
  onUnmounted(() => {
    disconnect()
  })

  return {
    isConnected: readonly(isConnected),
    connect,
    disconnect,
    subscribe,
    unsubscribe,
    on,
    off
  }
}
