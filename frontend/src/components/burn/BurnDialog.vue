<script setup>
import { ref, computed } from 'vue'
import Modal from '../ui/Modal.vue'
import Button from '../ui/Button.vue'
import BurnOptions from './BurnOptions.vue'
import BurnProgress from './BurnProgress.vue'
import BurnLog from './BurnLog.vue'

const props = defineProps({
  show: Boolean,
  availableSpeeds: { type: Array, default: () => [1, 2, 4, 8, 16] },
})

const emit = defineEmits(['close', 'start-burn'])

// States: 'configure' | 'burning' | 'done' | 'error'
const state = ref('configure')

const burnOptionsRef = ref(null)
const showLog = ref(false)

// Progress data (will come from Wails backend events)
const progress = ref({
  percent: 0,
  phase: 'Writing',
  speed: '',
  fifoPercent: 0,
  eta: '',
  bytesWritten: 0,
  bytesTotal: 0,
})

const logLines = ref([])

const modalTitle = computed(() => {
  const map = {
    configure: 'Burn Disc',
    burning: 'Burning...',
    done: 'Burn Complete',
    error: 'Burn Failed',
  }
  return map[state.value]
})

const canClose = computed(() => state.value !== 'burning')

function startBurn() {
  if (!burnOptionsRef.value) return
  const options = burnOptionsRef.value.options
  state.value = 'burning'
  logLines.value = []
  logLines.value.push('[INFO] Starting burn process...')
  logLines.value.push(`[INFO] Speed: ${options.speed === 0 ? 'Auto' : options.speed + 'x'}`)
  logLines.value.push(`[INFO] Mode: ${options.writeMode}`)
  logLines.value.push(`[INFO] Verify: ${options.verify ? 'Yes' : 'No'}`)
  logLines.value.push(`[INFO] Dummy mode: ${options.dummyMode ? 'Yes' : 'No'}`)

  emit('start-burn', { ...options })

  // Mock progress for demo (remove when wired to backend)
  simulateProgress()
}

function simulateProgress() {
  let pct = 0
  const interval = setInterval(() => {
    pct += Math.random() * 3
    if (pct >= 100) {
      pct = 100
      clearInterval(interval)
      progress.value.phase = 'Verifying'
      logLines.value.push('[INFO] Writing complete, verifying...')
      setTimeout(() => {
        state.value = 'done'
        logLines.value.push('[SUCCESS] Burn completed successfully!')
      }, 1500)
    }
    progress.value.percent = Math.min(pct, 100)
    progress.value.speed = '8.2x (11.2 MB/s)'
    progress.value.fifoPercent = Math.floor(80 + Math.random() * 20)
    progress.value.eta = pct < 100 ? `${Math.ceil((100 - pct) / 3)}s` : '0s'
    progress.value.bytesWritten = Math.floor((pct / 100) * 4.7e9)
    progress.value.bytesTotal = 4.7e9
  }, 300)
}

function handleClose() {
  if (state.value === 'burning') return
  state.value = 'configure'
  logLines.value = []
  progress.value = {
    percent: 0,
    phase: 'Writing',
    speed: '',
    fifoPercent: 0,
    eta: '',
    bytesWritten: 0,
    bytesTotal: 0,
  }
  emit('close')
}
</script>

<template>
  <Modal :show="show" :title="modalTitle" size="md" @close="canClose && handleClose()">
    <!-- Configure state -->
    <div v-if="state === 'configure'">
      <BurnOptions ref="burnOptionsRef" :available-speeds="availableSpeeds" />
    </div>

    <!-- Burning state -->
    <div v-else-if="state === 'burning'" class="space-y-4">
      <BurnProgress
        :percent="progress.percent"
        :phase="progress.phase"
        :speed="progress.speed"
        :fifo-percent="progress.fifoPercent"
        :eta="progress.eta"
        :bytes-written="progress.bytesWritten"
        :bytes-total="progress.bytesTotal"
      />
    </div>

    <!-- Done state -->
    <div v-else-if="state === 'done'" class="text-center py-6">
      <svg class="w-16 h-16 mx-auto text-green-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <h3 class="text-lg font-semibold text-gray-100 mb-1">Burn Complete!</h3>
      <p class="text-sm text-gray-400">The disc was burned and verified successfully.</p>
    </div>

    <!-- Error state -->
    <div v-else-if="state === 'error'" class="text-center py-6">
      <svg class="w-16 h-16 mx-auto text-red-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <h3 class="text-lg font-semibold text-gray-100 mb-1">Burn Failed</h3>
      <p class="text-sm text-gray-400">An error occurred during the burn process. Check the log for details.</p>
    </div>

    <!-- Log toggle (visible in burning/done/error states) -->
    <div v-if="state !== 'configure'" class="mt-4">
      <button
        class="text-xs text-gray-500 hover:text-gray-300 transition-colors flex items-center gap-1"
        @click="showLog = !showLog"
      >
        <svg
          class="w-3 h-3 transition-transform"
          :class="{ 'rotate-90': showLog }"
          fill="currentColor"
          viewBox="0 0 20 20"
        >
          <path d="M6 6l8 4-8 4V6z" />
        </svg>
        {{ showLog ? 'Hide' : 'Show' }} log output
      </button>
      <div v-if="showLog" class="mt-2">
        <BurnLog :lines="logLines" />
      </div>
    </div>

    <!-- Footer buttons -->
    <template #footer>
      <Button
        v-if="state === 'configure'"
        variant="secondary"
        @click="handleClose"
      >
        Cancel
      </Button>
      <Button
        v-if="state === 'configure'"
        variant="primary"
        @click="startBurn"
      >
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7C14 5 16.09 5.777 17.656 7.343A7.975 7.975 0 0120 13a7.975 7.975 0 01-2.343 5.657z" />
        </svg>
        Start Burn
      </Button>
      <Button
        v-if="state === 'burning'"
        variant="danger"
        disabled
      >
        Cancel (not implemented)
      </Button>
      <Button
        v-if="state === 'done' || state === 'error'"
        variant="secondary"
        @click="handleClose"
      >
        Close
      </Button>
    </template>
  </Modal>
</template>
