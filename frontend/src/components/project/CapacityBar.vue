<script setup>
import { computed } from 'vue'

const props = defineProps({
  usedBytes: { type: Number, default: 0 },
  mediaType: { type: String, default: 'DVD' },
})

// Standard disc sizes in bytes
const discSizes = {
  'CD': 700 * 1024 * 1024,           // 700 MB
  'DVD': 4.7 * 1000 * 1000 * 1000,   // 4.7 GB (decimal)
  'DVD-DL': 8.5 * 1000 * 1000 * 1000,// 8.5 GB
  'BD': 25 * 1000 * 1000 * 1000,     // 25 GB
  'BD-DL': 50 * 1000 * 1000 * 1000,  // 50 GB
  'BDXL': 100 * 1000 * 1000 * 1000,  // 100 GB
  'BDXL-QL': 128 * 1000 * 1000 * 1000, // 128 GB
}

// Tick marks for the bar — always show all standard sizes up to the max visible
const allTicks = [
  { label: 'CD 700MB', bytes: 700 * 1024 * 1024 },
  { label: 'DVD 4.7GB', bytes: 4.7e9 },
  { label: 'DVD-DL 8.5GB', bytes: 8.5e9 },
  { label: 'BD 25GB', bytes: 25e9 },
  { label: 'BD-DL 50GB', bytes: 50e9 },
  { label: 'BDXL 100GB', bytes: 100e9 },
  { label: 'BDXL 128GB', bytes: 128e9 },
]

const currentCapacity = computed(() => {
  // Try to match media type to a known disc size
  for (const key of Object.keys(discSizes)) {
    if (props.mediaType.toUpperCase().includes(key.toUpperCase())) {
      return discSizes[key]
    }
  }
  // Default to DVD
  return discSizes['DVD']
})

const usagePercent = computed(() => {
  if (currentCapacity.value === 0) return 0
  return (props.usedBytes / currentCapacity.value) * 100
})

const barColor = computed(() => {
  if (usagePercent.value > 95) return 'bg-red-500'
  if (usagePercent.value > 80) return 'bg-yellow-500'
  return 'bg-green-500'
})

// Max bytes shown on the bar — the largest disc or current capacity, whichever is larger
const maxBarBytes = computed(() => {
  return Math.max(currentCapacity.value, props.usedBytes * 1.05)
})

const visibleTicks = computed(() => {
  return allTicks.filter((t) => t.bytes <= maxBarBytes.value * 1.1)
})

function tickPosition(tickBytes) {
  return (tickBytes / maxBarBytes.value) * 100
}

function formatBytes(bytes) {
  if (bytes === 0) return '0 B'
  if (bytes >= 1e9) return (bytes / 1e9).toFixed(1) + ' GB'
  if (bytes >= 1e6) return (bytes / 1e6).toFixed(0) + ' MB'
  return (bytes / 1024).toFixed(0) + ' KB'
}
</script>

<template>
  <div class="w-full space-y-1">
    <!-- Used / Total text -->
    <div class="flex justify-between text-xs text-gray-400">
      <span>{{ formatBytes(usedBytes) }} used</span>
      <span>{{ formatBytes(currentCapacity) }} ({{ mediaType }})</span>
    </div>

    <!-- Bar with ticks -->
    <div class="relative">
      <!-- Background bar -->
      <div class="w-full bg-gray-700 rounded-full h-4 overflow-hidden relative">
        <!-- Used portion -->
        <div
          :class="['h-full rounded-full transition-all duration-500', barColor]"
          :style="{ width: Math.min(usagePercent, 100) + '%' }"
        />
        <!-- Overflow indicator -->
        <div
          v-if="usagePercent > 100"
          class="absolute inset-0 bg-red-500/20 animate-pulse rounded-full"
        />
      </div>

      <!-- Tick marks -->
      <div class="relative h-4 mt-0.5">
        <div
          v-for="tick in visibleTicks"
          :key="tick.label"
          class="absolute flex flex-col items-center"
          :style="{ left: tickPosition(tick.bytes) + '%', transform: 'translateX(-50%)' }"
        >
          <div class="w-px h-2 bg-gray-600" />
          <span class="text-[10px] text-gray-500 whitespace-nowrap mt-0.5">
            {{ tick.label }}
          </span>
        </div>
      </div>
    </div>

    <!-- Overflow warning -->
    <div
      v-if="usagePercent > 100"
      class="text-xs text-red-400 font-medium"
    >
      Overflow: {{ formatBytes(usedBytes - currentCapacity) }} over capacity!
    </div>
  </div>
</template>
