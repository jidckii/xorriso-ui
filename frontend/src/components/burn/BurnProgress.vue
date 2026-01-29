<script setup>
import { computed } from 'vue'
import ProgressBar from '../ui/ProgressBar.vue'

const props = defineProps({
  percent: { type: Number, default: 0 },
  phase: { type: String, default: 'Writing' }, // Formatting, Writing, Verifying
  speed: { type: String, default: '' },
  fifoPercent: { type: Number, default: 0 },
  eta: { type: String, default: '' },
  bytesWritten: { type: Number, default: 0 },
  bytesTotal: { type: Number, default: 0 },
})

const progressVariant = computed(() => {
  if (props.phase === 'Verifying') return 'success'
  if (props.percent >= 95) return 'warning'
  return 'default'
})

const phaseIcon = computed(() => {
  const map = {
    Formatting: 'M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16',
    Writing: 'M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7C14 5 16.09 5.777 17.656 7.343A7.975 7.975 0 0120 13a7.975 7.975 0 01-2.343 5.657z',
    Verifying: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
  }
  return map[props.phase] || map.Writing
})

function formatBytes(bytes) {
  if (bytes === 0) return '0 B'
  if (bytes >= 1e9) return (bytes / 1e9).toFixed(2) + ' GB'
  if (bytes >= 1e6) return (bytes / 1e6).toFixed(1) + ' MB'
  return (bytes / 1024).toFixed(0) + ' KB'
}
</script>

<template>
  <div class="space-y-6">
    <!-- Phase + percentage display -->
    <div class="text-center">
      <div class="flex items-center justify-center gap-2 mb-2">
        <svg class="w-5 h-5 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="phaseIcon" />
        </svg>
        <span class="text-sm font-medium text-gray-300 uppercase tracking-wider">
          {{ phase }}
        </span>
      </div>
      <div class="text-5xl font-bold text-gray-100 tabular-nums">
        {{ percent.toFixed(0) }}<span class="text-2xl text-gray-400">%</span>
      </div>
    </div>

    <!-- Main progress bar -->
    <ProgressBar :value="percent" :variant="progressVariant" size="lg" />

    <!-- Stats grid -->
    <div class="grid grid-cols-2 gap-3">
      <!-- Speed -->
      <div class="bg-gray-800/50 rounded-lg p-3 border border-gray-700">
        <div class="text-xs text-gray-500 mb-1">Speed</div>
        <div class="text-sm text-gray-200 font-medium">
          {{ speed || '--' }}
        </div>
      </div>

      <!-- ETA -->
      <div class="bg-gray-800/50 rounded-lg p-3 border border-gray-700">
        <div class="text-xs text-gray-500 mb-1">Time remaining</div>
        <div class="text-sm text-gray-200 font-medium">
          {{ eta || '--:--' }}
        </div>
      </div>

      <!-- Bytes written -->
      <div class="bg-gray-800/50 rounded-lg p-3 border border-gray-700">
        <div class="text-xs text-gray-500 mb-1">Written</div>
        <div class="text-sm text-gray-200 font-medium">
          {{ formatBytes(bytesWritten) }}
          <span v-if="bytesTotal" class="text-gray-500">
            / {{ formatBytes(bytesTotal) }}
          </span>
        </div>
      </div>

      <!-- FIFO buffer -->
      <div class="bg-gray-800/50 rounded-lg p-3 border border-gray-700">
        <div class="text-xs text-gray-500 mb-1">FIFO Buffer</div>
        <div class="flex items-center gap-2">
          <div class="flex-1 bg-gray-700 rounded-full h-2 overflow-hidden">
            <div
              class="h-full rounded-full transition-all duration-300"
              :class="{
                'bg-green-500': fifoPercent > 50,
                'bg-yellow-500': fifoPercent > 20 && fifoPercent <= 50,
                'bg-red-500': fifoPercent <= 20,
              }"
              :style="{ width: fifoPercent + '%' }"
            />
          </div>
          <span class="text-xs text-gray-400 tabular-nums w-8 text-right">
            {{ fifoPercent }}%
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
