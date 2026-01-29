<script setup>
import { ref } from 'vue'
import Button from '../ui/Button.vue'
import DeviceSelector from '../device/DeviceSelector.vue'

defineProps({
  devices: { type: Array, default: () => [] },
  selectedDevice: { type: String, default: '' },
})

const emit = defineEmits([
  'new-project',
  'open-project',
  'save-project',
  'burn',
  'update:selectedDevice',
])
</script>

<template>
  <header class="bg-gray-800 border-b border-gray-700 px-4 py-2 flex items-center gap-4">
    <!-- App name -->
    <div class="flex items-center gap-2 mr-4">
      <svg class="w-6 h-6 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <circle cx="12" cy="12" r="10" stroke-width="1.5" />
        <circle cx="12" cy="12" r="3" stroke-width="1.5" />
        <circle cx="12" cy="12" r="6" stroke-width="0.5" opacity="0.5" />
      </svg>
      <span class="text-lg font-bold text-gray-100 whitespace-nowrap">xorriso-ui</span>
    </div>

    <!-- Toolbar buttons -->
    <div class="flex items-center gap-1">
      <Button variant="ghost" size="sm" @click="emit('new-project')">
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        New
      </Button>
      <Button variant="ghost" size="sm" @click="emit('open-project')">
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
        </svg>
        Open
      </Button>
      <Button variant="ghost" size="sm" @click="emit('save-project')">
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4" />
        </svg>
        Save
      </Button>

      <div class="w-px h-6 bg-gray-700 mx-2" />

      <Button variant="primary" size="sm" @click="emit('burn')">
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7C14 5 16.09 5.777 17.656 7.343A7.975 7.975 0 0120 13a7.975 7.975 0 01-2.343 5.657z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M9.879 16.121A3 3 0 1012.015 11L11 14H9c0 .768.293 1.536.879 2.121z" />
        </svg>
        Burn
      </Button>
    </div>

    <!-- Spacer -->
    <div class="flex-1" />

    <!-- Device selector -->
    <DeviceSelector
      :devices="devices"
      :model-value="selectedDevice"
      @update:model-value="emit('update:selectedDevice', $event)"
    />
  </header>
</template>
