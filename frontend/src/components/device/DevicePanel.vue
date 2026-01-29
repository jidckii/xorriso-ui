<script setup>
import { computed } from 'vue'
import Button from '../ui/Button.vue'

const props = defineProps({
  device: {
    type: Object,
    default: () => ({
      id: '',
      vendor: '',
      model: '',
      mediaType: '',
      mediaStatus: 'unknown',
      capacityUsed: 0,
      capacityTotal: 0,
      speeds: [],
    }),
  },
})

const emit = defineEmits(['eject', 'blank'])

const usagePercent = computed(() => {
  if (!props.device.capacityTotal) return 0
  return (props.device.capacityUsed / props.device.capacityTotal) * 100
})

const statusColor = computed(() => {
  const map = {
    blank: 'text-green-400',
    appendable: 'text-yellow-400',
    closed: 'text-red-400',
    unknown: 'text-gray-400',
  }
  return map[props.device.mediaStatus] || 'text-gray-400'
})

const statusDot = computed(() => {
  const map = {
    blank: 'bg-green-500',
    appendable: 'bg-yellow-500',
    closed: 'bg-red-500',
    unknown: 'bg-gray-500',
  }
  return map[props.device.mediaStatus] || 'bg-gray-500'
})

function formatBytes(bytes) {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i]
}
</script>

<template>
  <div class="bg-gray-800 rounded-lg border border-gray-700 p-4 space-y-4">
    <h3 class="text-sm font-semibold text-gray-200 uppercase tracking-wider">
      Device Info
    </h3>

    <div v-if="device.id" class="space-y-3">
      <!-- Vendor / Model -->
      <div>
        <div class="text-xs text-gray-500">Device</div>
        <div class="text-sm text-gray-200">
          {{ device.vendor }} {{ device.model }}
        </div>
        <div class="text-xs text-gray-500 mt-0.5">{{ device.id }}</div>
      </div>

      <!-- Media Type -->
      <div>
        <div class="text-xs text-gray-500">Media Type</div>
        <div class="text-sm text-gray-200">
          {{ device.mediaType || 'No media inserted' }}
        </div>
      </div>

      <!-- Media Status -->
      <div>
        <div class="text-xs text-gray-500">Status</div>
        <div class="flex items-center gap-2">
          <span :class="['w-2 h-2 rounded-full', statusDot]" />
          <span :class="['text-sm capitalize', statusColor]">
            {{ device.mediaStatus }}
          </span>
        </div>
      </div>

      <!-- Capacity bar -->
      <div v-if="device.capacityTotal > 0">
        <div class="text-xs text-gray-500 mb-1">Capacity</div>
        <div class="w-full bg-gray-700 rounded-full h-2.5 overflow-hidden">
          <div
            class="h-full rounded-full transition-all duration-300"
            :class="{
              'bg-green-500': usagePercent <= 80,
              'bg-yellow-500': usagePercent > 80 && usagePercent <= 95,
              'bg-red-500': usagePercent > 95,
            }"
            :style="{ width: Math.min(usagePercent, 100) + '%' }"
          />
        </div>
        <div class="flex justify-between text-xs text-gray-400 mt-1">
          <span>{{ formatBytes(device.capacityUsed) }} used</span>
          <span>{{ formatBytes(device.capacityTotal) }} total</span>
        </div>
      </div>

      <!-- Speeds -->
      <div v-if="device.speeds && device.speeds.length">
        <div class="text-xs text-gray-500">Available Speeds</div>
        <div class="flex flex-wrap gap-1 mt-1">
          <span
            v-for="speed in device.speeds"
            :key="speed"
            class="px-1.5 py-0.5 text-xs bg-gray-700 text-gray-300 rounded"
          >
            {{ speed }}x
          </span>
        </div>
      </div>

      <!-- Actions -->
      <div class="flex gap-2 pt-2 border-t border-gray-700">
        <Button size="sm" variant="secondary" @click="emit('eject')">
          <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 19h14" />
          </svg>
          Eject
        </Button>
        <Button size="sm" variant="danger" @click="emit('blank')">
          Blank / Format
        </Button>
      </div>
    </div>

    <div v-else class="text-sm text-gray-500 py-4 text-center">
      No device selected
    </div>
  </div>
</template>
