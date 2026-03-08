<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { formatBytes } from '../../composables/useFormatBytes'
import { useMediaStatus } from '../../composables/useMediaStatus'

const { t } = useI18n()

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

const { statusDot: statusColor, statusLabel } = useMediaStatus(() => props.status)
</script>

<template>
  <div class="flex items-center gap-3 text-sm">
    <!-- Media type badge -->
    <span
      v-if="mediaType"
      class="px-2 py-0.5 rounded text-xs font-medium bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 border border-gray-400 dark:border-gray-600"
    >
      {{ mediaType }}
    </span>
    <span v-else class="text-gray-500 text-xs">{{ t('device.noMedia') }}</span>

    <!-- Capacity -->
    <span v-if="capacityTotal > 0" class="text-gray-600 dark:text-gray-400">
      {{ formatBytes(capacityUsed) }} / {{ formatBytes(capacityTotal) }}
    </span>

    <!-- Status indicator -->
    <span class="flex items-center gap-1.5">
      <span :class="['w-2 h-2 rounded-full', statusColor]" />
      <span class="text-gray-600 dark:text-gray-400 text-xs">{{ statusLabel }}</span>
    </span>
  </div>
</template>
