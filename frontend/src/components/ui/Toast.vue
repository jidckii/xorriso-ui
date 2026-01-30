<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'

const props = defineProps({
  message: { type: String, required: true },
  type: {
    type: String,
    default: 'info',
    validator: (v) => ['info', 'success', 'warning', 'error'].includes(v),
  },
  duration: { type: Number, default: 4000 },
})

const emit = defineEmits(['close'])

const visible = ref(true)
let timer = null

onMounted(() => {
  if (props.duration > 0) {
    timer = setTimeout(() => {
      visible.value = false
      emit('close')
    }, props.duration)
  }
})

onUnmounted(() => {
  if (timer) clearTimeout(timer)
})

const typeStyles = computed(() => {
  const map = {
    info: 'bg-blue-50 dark:bg-blue-900/80 border-blue-600 text-blue-800 dark:text-blue-200',
    success: 'bg-green-50 dark:bg-green-900/80 border-green-600 text-green-800 dark:text-green-200',
    warning: 'bg-yellow-50 dark:bg-yellow-900/80 border-yellow-600 text-yellow-800 dark:text-yellow-200',
    error: 'bg-red-50 dark:bg-red-900/80 border-red-600 text-red-800 dark:text-red-200',
  }
  return map[props.type]
})

const iconPath = computed(() => {
  const map = {
    info: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
    success: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
    warning: 'M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
    error: 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z',
  }
  return map[props.type]
})

function close() {
  visible.value = false
  emit('close')
}
</script>

<template>
  <Transition
    enter-active-class="transition duration-200 ease-out"
    enter-from-class="opacity-0 translate-y-2"
    enter-to-class="opacity-100 translate-y-0"
    leave-active-class="transition duration-150 ease-in"
    leave-from-class="opacity-100 translate-y-0"
    leave-to-class="opacity-0 translate-y-2"
  >
    <div
      v-if="visible"
      :class="[
        'flex items-center gap-3 px-4 py-3 rounded-lg border shadow-lg min-w-[300px] max-w-md',
        typeStyles,
      ]"
    >
      <svg
        class="w-5 h-5 flex-shrink-0"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          :d="iconPath"
        />
      </svg>
      <span class="flex-1 text-sm">{{ message }}</span>
      <button
        class="flex-shrink-0 opacity-70 hover:opacity-100 transition-opacity"
        @click="close"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          />
        </svg>
      </button>
    </div>
  </Transition>
</template>
