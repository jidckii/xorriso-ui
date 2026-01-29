<script setup>
import { computed } from 'vue'

const props = defineProps({
  value: { type: Number, default: 0 },
  variant: {
    type: String,
    default: 'default',
    validator: (v) => ['default', 'success', 'warning', 'danger'].includes(v),
  },
  size: {
    type: String,
    default: 'md',
    validator: (v) => ['sm', 'md', 'lg'].includes(v),
  },
  showLabel: { type: Boolean, default: false },
})

const clampedValue = computed(() => Math.max(0, Math.min(100, props.value)))

const barColor = computed(() => {
  const map = {
    default: 'bg-blue-500',
    success: 'bg-green-500',
    warning: 'bg-yellow-500',
    danger: 'bg-red-500',
  }
  return map[props.variant]
})

const sizeClass = computed(() => {
  const map = {
    sm: 'h-1.5',
    md: 'h-3',
    lg: 'h-5',
  }
  return map[props.size]
})
</script>

<template>
  <div class="w-full">
    <div
      :class="['w-full bg-gray-700 rounded-full overflow-hidden', sizeClass]"
    >
      <div
        :class="[
          'h-full rounded-full transition-all duration-300 ease-in-out',
          barColor,
        ]"
        :style="{ width: clampedValue + '%' }"
      />
    </div>
    <div
      v-if="showLabel"
      class="mt-1 text-xs text-gray-400 text-right"
    >
      {{ clampedValue.toFixed(0) }}%
    </div>
  </div>
</template>
