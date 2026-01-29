<script setup>
import { reactive } from 'vue'

const props = defineProps({
  availableSpeeds: {
    type: Array,
    default: () => [1, 2, 4, 8, 16, 24, 48],
  },
})

const options = reactive({
  speed: 0, // 0 = auto
  writeMode: 'auto',
  verify: true,
  dummyMode: false,
  finalize: false,
  ejectAfter: true,
  streamRecording: false,
})

defineExpose({ options })
</script>

<template>
  <div class="space-y-4">
    <!-- Speed -->
    <div>
      <label class="block text-sm font-medium text-gray-300 mb-1">Write Speed</label>
      <select
        v-model.number="options.speed"
        class="w-full bg-gray-700 text-gray-200 text-sm rounded px-3 py-2 border border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option :value="0">Auto (Maximum)</option>
        <option
          v-for="speed in availableSpeeds"
          :key="speed"
          :value="speed"
        >
          {{ speed }}x
        </option>
      </select>
    </div>

    <!-- Write Mode -->
    <div>
      <label class="block text-sm font-medium text-gray-300 mb-1">Write Mode</label>
      <select
        v-model="options.writeMode"
        class="w-full bg-gray-700 text-gray-200 text-sm rounded px-3 py-2 border border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="auto">Auto</option>
        <option value="tao">TAO (Track At Once)</option>
        <option value="sao">SAO (Session At Once)</option>
      </select>
    </div>

    <!-- Checkboxes -->
    <div class="space-y-3">
      <label class="flex items-center gap-3 cursor-pointer group">
        <input
          v-model="options.verify"
          type="checkbox"
          class="w-4 h-4 rounded border-gray-600 bg-gray-700 text-blue-500 focus:ring-blue-500 focus:ring-offset-gray-900"
        />
        <div>
          <span class="text-sm text-gray-200 group-hover:text-white">Verify after burn</span>
          <p class="text-xs text-gray-500">Read back and compare data after writing</p>
        </div>
      </label>

      <label class="flex items-center gap-3 cursor-pointer group">
        <input
          v-model="options.dummyMode"
          type="checkbox"
          class="w-4 h-4 rounded border-gray-600 bg-gray-700 text-blue-500 focus:ring-blue-500 focus:ring-offset-gray-900"
        />
        <div>
          <span class="text-sm text-gray-200 group-hover:text-white">Dummy mode (test)</span>
          <p class="text-xs text-gray-500">Simulate burn without writing to disc</p>
        </div>
      </label>

      <label class="flex items-center gap-3 cursor-pointer group">
        <input
          v-model="options.finalize"
          type="checkbox"
          class="w-4 h-4 rounded border-gray-600 bg-gray-700 text-blue-500 focus:ring-blue-500 focus:ring-offset-gray-900"
        />
        <div>
          <span class="text-sm text-gray-200 group-hover:text-white">Finalize disc</span>
          <p class="text-xs text-gray-500">Close the disc so no more data can be added</p>
        </div>
      </label>

      <label class="flex items-center gap-3 cursor-pointer group">
        <input
          v-model="options.ejectAfter"
          type="checkbox"
          class="w-4 h-4 rounded border-gray-600 bg-gray-700 text-blue-500 focus:ring-blue-500 focus:ring-offset-gray-900"
        />
        <div>
          <span class="text-sm text-gray-200 group-hover:text-white">Eject after burn</span>
          <p class="text-xs text-gray-500">Open disc tray when finished</p>
        </div>
      </label>

      <label class="flex items-center gap-3 cursor-pointer group">
        <input
          v-model="options.streamRecording"
          type="checkbox"
          class="w-4 h-4 rounded border-gray-600 bg-gray-700 text-blue-500 focus:ring-blue-500 focus:ring-offset-gray-900"
        />
        <div>
          <span class="text-sm text-gray-200 group-hover:text-white">Stream recording</span>
          <p class="text-xs text-gray-500">For Blu-ray discs: faster writing, no defect management</p>
        </div>
      </label>
    </div>
  </div>
</template>
