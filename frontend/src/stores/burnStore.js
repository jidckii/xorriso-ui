import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { StartBurn, CancelBurn, BlankDisc, FormatDisc, GetJobStatus } from '../../bindings/xorriso-ui/services/burnservice.js'
import { Events } from '@wailsio/runtime'

export const useBurnStore = defineStore('burn', () => {
  // --- State ---
  const currentJob = ref(null)
  const logLines = ref([])
  const isBurning = computed(() => currentJob.value !== null && ['preparing', 'burning', 'verifying', 'blanking', 'formatting'].includes(currentJob.value.state))

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

  async function startBurn(project, devicePath, opts) {
    logLines.value = []
    try {
      const jobId = await StartBurn(project, devicePath, opts)

      currentJob.value = {
        id: jobId,
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
        startedAt: new Date().toISOString(),
        finishedAt: null,
      }

      addLogLine('Starting burn process...')
      addLogLine(`Target device: ${devicePath}`)
    } catch (error) {
      console.error('Failed to start burn:', error)
      addLogLine(`ERROR: ${error.message || error}`)
      currentJob.value = null
    }
  }

  async function cancelBurn() {
    if (!currentJob.value) return
    try {
      await CancelBurn(currentJob.value.id)
      addLogLine('Burn cancellation requested.')
    } catch (error) {
      console.error('Failed to cancel burn:', error)
      addLogLine(`ERROR: Failed to cancel: ${error.message || error}`)
    }
  }

  async function blankDisc(devicePath, mode = 'fast') {
    logLines.value = []
    try {
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
        startedAt: new Date().toISOString(),
        finishedAt: null,
      }

      addLogLine(`Blanking disc on ${devicePath} (mode: ${mode})...`)
      await BlankDisc(devicePath, mode)
    } catch (error) {
      console.error('Failed to blank disc:', error)
      addLogLine(`ERROR: ${error.message || error}`)
      if (currentJob.value) {
        currentJob.value.state = 'error'
        currentJob.value.result = { success: false }
        currentJob.value.finishedAt = new Date().toISOString()
      }
    }
  }

  async function formatDisc(devicePath, mode = 'default') {
    logLines.value = []
    try {
      currentJob.value = {
        id: `format-${Date.now()}`,
        state: 'formatting',
        progress: {
          phase: 'formatting',
          percent: 0,
          speed: '',
          bytesWritten: 0,
          bytesTotal: 0,
          eta: '',
          fifoFill: 0,
        },
        result: null,
        startedAt: new Date().toISOString(),
        finishedAt: null,
      }

      addLogLine(`Formatting disc on ${devicePath} (mode: ${mode})...`)
      await FormatDisc(devicePath, mode)
    } catch (error) {
      console.error('Failed to format disc:', error)
      addLogLine(`ERROR: ${error.message || error}`)
      if (currentJob.value) {
        currentJob.value.state = 'error'
        currentJob.value.result = { success: false }
        currentJob.value.finishedAt = new Date().toISOString()
      }
    }
  }

  async function fetchJobStatus() {
    if (!currentJob.value) return null
    try {
      const job = await GetJobStatus(currentJob.value.id)
      if (job) {
        currentJob.value = job
      }
      return job
    } catch (error) {
      console.error('Failed to get job status:', error)
      return null
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
    Events.On('burn:progress', (data) => {
      if (currentJob.value) {
        currentJob.value.progress = data
      }
    })

    Events.On('burn:state-changed', (data) => {
      if (currentJob.value) {
        currentJob.value.state = data.state
      }
    })

    Events.On('burn:log-line', (data) => {
      addLogLine(data.line)
    })

    Events.On('burn:complete', (data) => {
      if (currentJob.value) {
        currentJob.value.state = 'complete'
        currentJob.value.result = data
        currentJob.value.finishedAt = new Date().toISOString()
      }
    })

    Events.On('burn:error', (data) => {
      if (currentJob.value) {
        currentJob.value.state = 'error'
        currentJob.value.result = { success: false }
        currentJob.value.finishedAt = new Date().toISOString()
      }
      addLogLine(`ERROR: ${data.error}`)
    })
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
    formatDisc,
    fetchJobStatus,
    resetJob,
    addLogLine,
    init,
  }
})
