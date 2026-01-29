<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'secondary',
    validator: (v) => ['primary', 'secondary', 'danger', 'ghost'].includes(v),
  },
  size: {
    type: String,
    default: 'md',
    validator: (v) => ['sm', 'md', 'lg'].includes(v),
  },
  disabled: Boolean,
  loading: Boolean,
})

defineEmits(['click'])

const variantClasses = computed(() => {
  const map = {
    primary:
      'bg-blue-600 hover:bg-blue-700 text-white focus:ring-blue-500 disabled:bg-blue-800',
    secondary:
      'bg-gray-600 hover:bg-gray-500 text-gray-100 focus:ring-gray-500 disabled:bg-gray-700',
    danger:
      'bg-red-600 hover:bg-red-700 text-white focus:ring-red-500 disabled:bg-red-800',
    ghost:
      'bg-transparent hover:bg-gray-700 text-gray-300 focus:ring-gray-500',
  }
  return map[props.variant]
})

const sizeClasses = computed(() => {
  const map = {
    sm: 'px-2 py-1 text-xs',
    md: 'px-3 py-1.5 text-sm',
    lg: 'px-4 py-2 text-base',
  }
  return map[props.size]
})
</script>

<template>
  <button
    :class="[
      'inline-flex items-center justify-center gap-1.5 rounded font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-1 focus:ring-offset-gray-900 disabled:opacity-50 disabled:cursor-not-allowed',
      variantClasses,
      sizeClasses,
    ]"
    :disabled="disabled || loading"
    @click="$emit('click', $event)"
  >
    <svg
      v-if="loading"
      class="animate-spin h-4 w-4"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle
        class="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        stroke-width="4"
      />
      <path
        class="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
      />
    </svg>
    <slot />
  </button>
</template>
