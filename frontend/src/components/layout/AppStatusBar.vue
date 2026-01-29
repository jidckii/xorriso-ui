<script setup>
defineProps({
  deviceInfo: { type: String, default: 'No device' },
  mediaType: { type: String, default: '' },
  freeSpace: { type: String, default: '' },
  operationStatus: { type: String, default: 'Ready' },
})
</script>

<template>
  <footer class="bg-gray-800 border-t border-gray-700 px-4 h-8 flex items-center text-sm text-gray-400">
    <!-- Left: device info -->
    <div class="flex-1 truncate">
      {{ deviceInfo }}
    </div>

    <!-- Center: media info -->
    <div class="flex-1 text-center truncate">
      <span v-if="mediaType" class="text-gray-300">{{ mediaType }}</span>
      <span v-if="mediaType && freeSpace" class="mx-2 text-gray-600">|</span>
      <span v-if="freeSpace">{{ freeSpace }} free</span>
    </div>

    <!-- Right: status -->
    <div class="flex-1 text-right truncate">
      <span
        :class="{
          'text-green-400': operationStatus === 'Ready',
          'text-blue-400': operationStatus === 'Burning',
          'text-yellow-400': operationStatus === 'Verifying',
          'text-red-400': operationStatus === 'Error',
          'text-gray-400': !['Ready', 'Burning', 'Verifying', 'Error'].includes(operationStatus),
        }"
      >
        {{ operationStatus }}
      </span>
    </div>
  </footer>
</template>
