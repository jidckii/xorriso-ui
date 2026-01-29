import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
// TODO: import Wails service bindings when available
// import { StartBurn, CancelBurn, BlankDisc } from '../../bindings/xorriso-ui/BurnService'
// TODO: import Wails events when backend is ready
// import { Events } from '@wailsio/runtime'

export const useBurnStore = defineStore('burn', () => {
  // --- State ---
  const currentJob = ref(null)
  const logLines = ref([])
  const isBurning = computed(() => currentJob.value !== null && ['preparing', 'burning', 'verifying', 'blanking'].includes(currentJob.value.state))

  // Burn progress details
  const progress = computed(() => {
    if (!currentJob.value) {
      return {
        phase: 'idle',
        percent: 0,
        speed: '',
        bytesWritten: 0,
        bytesTotal: 0,
        eta: '',
        fifoFill: 0,
      }
    }
    return currentJob.value.progress || {
      phase: currentJob.value.state,
      percent: 0,
      speed: '',
      bytesWritten: 0,
      bytesTotal: 0,
      eta: '',
      fifoFill: 0,
    }
  })

  // --- Actions ---

  async function startBurn(devicePath) {
    logLines.value = []
    try {
      // TODO: replace with actual Wails call
      // const job = await StartBurn(project, devicePath, opts)
      // currentJob.value = job

      // Mock: create a job and simulate progress
      currentJob.value = {
        id: `burn-${Date.now()}`,
        state: 'preparing',
        progress: {
          phase: 'preparing',
          percent: 0,
          speed: '',
          bytesWritten: 0,
          bytesTotal: 0,
          eta: '',
          fifoFill: 0,
        },
        result: null,
      }

      addLogLine('Starting burn process...')
      addLogLine(`Target device: ${devicePath}`)

      // Mock: simulate burn progress for UI development
      simulateBurnProgress()
    } catch (error) {
      console.error('Failed to start burn:', error)
      currentJob.value = null
    }
  }

  async function cancelBurn() {
    if (!currentJob.value) return
    try {
      // TODO: replace with actual Wails call
      // await CancelBurn(currentJob.value.id)

      addLogLine('Burn cancelled by user.')
      currentJob.value.state = 'cancelled'
      currentJob.value.result = { success: false, message: 'Cancelled by user' }
    } catch (error) {
      console.error('Failed to cancel burn:', error)
    }
  }

  async function blankDisc(devicePath, mode = 'fast') {
    logLines.value = []
    try {
      // TODO: replace with actual Wails call
      // await BlankDisc(devicePath, mode)

      currentJob.value = {
        id: `blank-${Date.now()}`,
        state: 'blanking',
        progress: {
          phase: 'blanking',
          percent: 0,
          speed: '',
          bytesWritten: 0,
          bytesTotal: 0,
          eta: '',
          fifoFill: 0,
        },
        result: null,
      }

      addLogLine(`Blanking disc on ${devicePath} (mode: ${mode})...`)

      // Mock simulation
      setTimeout(() => {
        if (currentJob.value && currentJob.value.state === 'blanking') {
          currentJob.value.state = 'complete'
          currentJob.value.progress.percent = 100
          currentJob.value.result = { success: true, message: 'Disc blanked successfully' }
          addLogLine('Blanking complete.')
        }
      }, 3000)
    } catch (error) {
      console.error('Failed to blank disc:', error)
    }
  }

  function resetJob() {
    currentJob.value = null
    logLines.value = []
  }

  function addLogLine(line) {
    const timestamp = new Date().toLocaleTimeString()
    logLines.value.push(`[${timestamp}] ${line}`)
  }

  /**
   * Initialize event listeners for burn-related Wails events.
   * Call this once from the app root or on store creation.
   */
  function init() {
    // TODO: replace with actual Wails event subscriptions
    // Events.On('burn:progress', (data) => {
    //   if (currentJob.value) {
    //     currentJob.value.progress = data
    //   }
    // })
    // Events.On('burn:state-changed', (data) => {
    //   if (currentJob.value) {
    //     currentJob.value.state = data.state
    //   }
    // })
    // Events.On('burn:log-line', (data) => {
    //   addLogLine(data.line)
    // })
    // Events.On('burn:complete', (data) => {
    //   if (currentJob.value) {
    //     currentJob.value.state = 'complete'
    //     currentJob.value.result = data
    //   }
    // })
    // Events.On('burn:error', (data) => {
    //   if (currentJob.value) {
    //     currentJob.value.state = 'error'
    //     currentJob.value.result = { success: false, message: data.error }
    //   }
    //   addLogLine(`ERROR: ${data.error}`)
    // })
  }

  // --- Mock helpers (remove when real backend is connected) ---

  function simulateBurnProgress() {
    let percent = 0
    const total = 734003200 // ~700 MB
    const interval = setInterval(() => {
      if (!currentJob.value || currentJob.value.state === 'cancelled') {
        clearInterval(interval)
        return
      }

      percent += 2
      if (percent <= 100) {
        const written = Math.floor((percent / 100) * total)
        currentJob.value.state = percent < 5 ? 'preparing' : percent < 95 ? 'burning' : 'verifying'
        currentJob.value.progress = {
          phase: currentJob.value.state,
          percent,
          speed: '8x (36.0 MB/s)',
          bytesWritten: written,
          bytesTotal: total,
          eta: `${Math.ceil((100 - percent) * 0.5)}s`,
          fifoFill: Math.min(100, 80 + Math.random() * 20),
        }

        if (percent % 10 === 0) {
          addLogLine(`Progress: ${percent}% — ${(written / 1024 / 1024).toFixed(1)} MB written`)
        }
      }

      if (percent >= 100) {
        clearInterval(interval)
        currentJob.value.state = 'complete'
        currentJob.value.result = { success: true, message: 'Burn completed successfully' }
        addLogLine('Burn completed successfully.')
      }
    }, 500)
  }

  return {
    // State
    currentJob,
    logLines,
    // Getters
    isBurning,
    progress,
    // Actions
    startBurn,
    cancelBurn,
    blankDisc,
    resetJob,
    addLogLine,
    init,
  }
})
