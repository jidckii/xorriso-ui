<script setup>
import { reactive } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

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
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('burnOptions.writeSpeed') }}</label>
      <select
        v-model.number="options.speed"
        class="w-full bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 text-sm rounded px-3 py-2 border border-gray-400 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option :value="0">{{ t('burnOptions.autoMaximum') }}</option>
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
      <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{{ t('burnOptions.writeMode') }}</label>
      <select
        v-model="options.writeMode"
        class="w-full bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 text-sm rounded px-3 py-2 border border-gray-400 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
      >
        <option value="auto">{{ t('burnOptions.auto') }}</option>
        <option value="tao">{{ t('burnOptions.tao') }}</option>
        <option value="sao">{{ t('burnOptions.sao') }}</option>
      </select>
    </div>

    <!-- Checkboxes -->
    <div class="space-y-3">
      <label class="flex items-center gap-3 cursor-pointer group">
        <input
          v-model="options.verify"
          type="checkbox"
          class="w-4 h-4 rounded border-gray-400 dark:border-gray-600 bg-gray-200 dark:bg-gray-700 text-blue-500 focus:ring-blue-500 focus:ring-offset-white dark:focus:ring-offset-gray-900"
        />
        <div>
          <span class="text-sm text-gray-800 dark:text-gray-200 group-hover:text-black dark:group-hover:text-white">{{ t('burnOptions.verifyAfterBurn') }}</span>
          <p class="text-xs text-gray-500">{{ t('burnOptions.verifyDescription') }}</p>
        </div>
      </label>

      <label class="flex items-center gap-3 cursor-pointer group">
        <input
          v-model="options.dummyMode"
          type="checkbox"
          class="w-4 h-4 rounded border-gray-400 dark:border-gray-600 bg-gray-200 dark:bg-gray-700 text-blue-500 focus:ring-blue-500 focus:ring-offset-white dark:focus:ring-offset-gray-900"
        />
        <div>
          <span class="text-sm text-gray-800 dark:text-gray-200 group-hover:text-black dark:group-hover:text-white">{{ t('burnOptions.dummyMode') }}</span>
          <p class="text-xs text-gray-500">{{ t('burnOptions.dummyDescription') }}</p>
        </div>
      </label>

      <label class="flex items-center gap-3 cursor-pointer group">
        <input
          v-model="options.finalize"
          type="checkbox"
          class="w-4 h-4 rounded border-gray-400 dark:border-gray-600 bg-gray-200 dark:bg-gray-700 text-blue-500 focus:ring-blue-500 focus:ring-offset-white dark:focus:ring-offset-gray-900"
        />
        <div>
          <span class="text-sm text-gray-800 dark:text-gray-200 group-hover:text-black dark:group-hover:text-white">{{ t('burnOptions.finalizeDisc') }}</span>
          <p class="text-xs text-gray-500">{{ t('burnOptions.finalizeDescription') }}</p>
        </div>
      </label>

      <label class="flex items-center gap-3 cursor-pointer group">
        <input
          v-model="options.ejectAfter"
          type="checkbox"
          class="w-4 h-4 rounded border-gray-400 dark:border-gray-600 bg-gray-200 dark:bg-gray-700 text-blue-500 focus:ring-blue-500 focus:ring-offset-white dark:focus:ring-offset-gray-900"
        />
        <div>
          <span class="text-sm text-gray-800 dark:text-gray-200 group-hover:text-black dark:group-hover:text-white">{{ t('burnOptions.ejectAfterBurn') }}</span>
          <p class="text-xs text-gray-500">{{ t('burnOptions.ejectDescription') }}</p>
        </div>
      </label>

      <label class="flex items-center gap-3 cursor-pointer group">
        <input
          v-model="options.streamRecording"
          type="checkbox"
          class="w-4 h-4 rounded border-gray-400 dark:border-gray-600 bg-gray-200 dark:bg-gray-700 text-blue-500 focus:ring-blue-500 focus:ring-offset-white dark:focus:ring-offset-gray-900"
        />
        <div>
          <span class="text-sm text-gray-800 dark:text-gray-200 group-hover:text-black dark:group-hover:text-white">{{ t('burnOptions.streamRecording') }}</span>
          <p class="text-xs text-gray-500">{{ t('burnOptions.streamDescription') }}</p>
        </div>
      </label>
    </div>
  </div>
</template>
