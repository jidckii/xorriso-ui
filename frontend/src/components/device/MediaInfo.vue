<script setup>
import { computed } from 'vue'

const props = defineProps({
  mediaType: { type: String, default: '' },
  capacityUsed: { type: Number, default: 0 },
  capacityTotal: { type: Number, default: 0 },
  status: {
    type: String,
    default: 'blank',
    validator: (v) => ['blank', 'appendable', 'closed', 'unknown'].includes(v),
  },
})

const statusColor = computed(() => {
  const map = {
    blank: 'bg-green-500',
    appendable: 'bg-yellow-500',
    closed: 'bg-red-500',
    unknown: 'bg-gray-500',
  }
  return map[props.status]
})

const statusLabel = computed(() => {
  const map = {
    blank: 'Blank',
    appendable: 'Appendable',
    closed: 'Closed',
    unknown: 'Unknown',
  }
  return map[props.status]
})

function formatBytes(bytes) {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + units[i]
}
</script>

<template>
  <div class="flex items-center gap-3 text-sm">
    <!-- Media type badge -->
    <span
      v-if="mediaType"
      class="px-2 py-0.5 rounded text-xs font-medium bg-gray-700 text-gray-200 border border-gray-600"
    >
      {{ mediaType }}
    </span>
    <span v-else class="text-gray-500 text-xs">No media</span>

    <!-- Capacity -->
    <span v-if="capacityTotal > 0" class="text-gray-400">
      {{ formatBytes(capacityUsed) }} / {{ formatBytes(capacityTotal) }}
    </span>

    <!-- Status indicator -->
    <span class="flex items-center gap-1.5">
      <span :class="['w-2 h-2 rounded-full', statusColor]" />
      <span class="text-gray-400 text-xs">{{ statusLabel }}</span>
    </span>
  </div>
</template>
