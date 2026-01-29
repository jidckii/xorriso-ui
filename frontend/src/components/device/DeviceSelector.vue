<script setup>
defineProps({
  devices: {
    type: Array,
    default: () => [],
    // Each device: { id: string, name: string, vendor: string, model: string }
  },
  modelValue: { type: String, default: '' },
})

const emit = defineEmits(['update:modelValue'])

function onChange(e) {
  emit('update:modelValue', e.target.value)
}
</script>

<template>
  <div class="relative">
    <select
      :value="modelValue"
      class="appearance-none bg-gray-700 text-gray-200 text-sm rounded px-3 py-1.5 pr-8 border border-gray-600 hover:border-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent cursor-pointer min-w-[200px]"
      @change="onChange"
    >
      <option value="" disabled>Select device...</option>
      <option
        v-for="device in devices"
        :key="device.id"
        :value="device.id"
      >
        {{ device.vendor }} {{ device.model }} ({{ device.id }})
      </option>
    </select>
    <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
      <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
    </div>
  </div>
</template>
