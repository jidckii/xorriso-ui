<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { formatBytes } from '../../composables/useFormatBytes'

const props = defineProps({
  project: { type: Object, required: true },
  devices: { type: Array, required: true },
  currentDevicePath: { type: String, default: '' },
  mediaInfo: { type: Object, default: null },
  speeds: { type: Array, default: () => [] },
  isBurning: { type: Boolean, default: false },
})

const emit = defineEmits(['select-device', 'start-burn', 'create-iso', 'blank-disc', 'go-back'])

const { t } = useI18n()

const blankMode = ref('fast')

const canBurn = computed(() =>
  props.currentDevicePath &&
  props.project?.entries?.length > 0 &&
  !props.isBurning
)
</script>

<template>
  <div class="p-6 space-y-6">

    <!-- Сводка проекта -->
    <div class="space-y-2">
      <h3 class="text-sm font-medium text-gray-600 dark:text-gray-400">{{ t('burn.project') }}</h3>
      <div class="bg-white dark:bg-gray-900 rounded px-4 py-3 text-sm space-y-1">
        <div class="flex justify-between">
          <span class="text-gray-500">{{ t('common.name') }}:</span>
          <span class="text-gray-900 dark:text-gray-100">{{ project?.name }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-500">{{ t('common.volumeId') }}:</span>
          <span class="text-gray-900 dark:text-gray-100">{{ project?.volumeId }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-500">{{ t('common.files') }}:</span>
          <span class="text-gray-900 dark:text-gray-100">{{ (project?.entries?.length || 0) }} {{ t('project.items') }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-500">{{ t('common.totalSize') }}:</span>
          <span class="text-gray-900 dark:text-gray-100">{{ formatBytes(project?.totalSize || 0) }}</span>
        </div>
      </div>
    </div>

    <!-- Выбор устройства -->
    <div class="space-y-2">
      <h3 class="text-sm font-medium text-gray-600 dark:text-gray-400">{{ t('burn.device') }}</h3>
      <select
        :value="currentDevicePath"
        @change="emit('select-device', ($event.target).value)"
        class="w-full px-3 py-2 text-sm bg-white dark:bg-gray-900 border border-gray-400 dark:border-gray-600 rounded text-gray-800 dark:text-gray-200 focus:outline-none focus:border-blue-500"
      >
        <option v-for="dev in devices" :key="dev.path" :value="dev.path">
          {{ dev.vendor }} {{ dev.model }} ({{ dev.path }})
        </option>
      </select>
      <div v-if="mediaInfo" class="text-xs text-gray-500">
        {{ mediaInfo.mediaType }} - {{ mediaInfo.mediaStatus }} -
        {{ formatBytes(mediaInfo.freeSpace) }} {{ t('device.free') }}
      </div>
    </div>

    <!-- Опции записи -->
    <div class="space-y-2">
      <h3 class="text-sm font-medium text-gray-600 dark:text-gray-400">{{ t('burn.burnOptions') }}</h3>
      <div class="grid grid-cols-2 gap-3">
        <div>
          <label class="block text-xs text-gray-500 mb-1">{{ t('burn.speed') }}</label>
          <select
            v-model="project.burnOptions.speed"
            class="w-full px-2 py-1.5 text-sm bg-white dark:bg-gray-900 border border-gray-400 dark:border-gray-600 rounded text-gray-800 dark:text-gray-200 focus:outline-none focus:border-blue-500"
          >
            <option v-for="s in speeds" :key="s.value" :value="s.value">
              {{ s.label }}
            </option>
          </select>
        </div>
        <div>
          <label class="block text-xs text-gray-500 mb-1">{{ t('burn.burnMode') }}</label>
          <select
            v-model="project.burnOptions.burnMode"
            class="w-full px-2 py-1.5 text-sm bg-white dark:bg-gray-900 border border-gray-400 dark:border-gray-600 rounded text-gray-800 dark:text-gray-200 focus:outline-none focus:border-blue-500"
          >
            <option value="auto">{{ t('burn.autoDao') }}</option>
            <option value="tao">{{ t('burn.tao') }}</option>
            <option value="sao">{{ t('burn.saoDao') }}</option>
          </select>
        </div>
      </div>

      <div class="grid grid-cols-2 gap-x-6 gap-y-2 mt-3">
        <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
          <input type="checkbox" v-model="project.burnOptions.verify" class="accent-blue-500" />
          {{ t('burn.verifyAfterBurn') }}
        </label>
        <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
          <input type="checkbox" v-model="project.burnOptions.eject" class="accent-blue-500" />
          {{ t('burn.ejectWhenDone') }}
        </label>
        <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
          <input type="checkbox" v-model="project.burnOptions.dummyMode" class="accent-yellow-500" />
          {{ t('burn.simulationMode') }}
        </label>
        <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
          <input type="checkbox" v-model="project.burnOptions.closeDisc" class="accent-blue-500" />
          {{ t('burn.closeDisc') }}
        </label>
        <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
          <input type="checkbox" v-model="project.burnOptions.streamRecording" class="accent-blue-500" />
          {{ t('burn.streamRecording') }}
        </label>
        <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
          <input type="checkbox" v-model="project.burnOptions.multisession" class="accent-blue-500" />
          {{ t('burn.multisession') }}
        </label>
        <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
          <input type="checkbox" v-model="project.burnOptions.cleanupIso" class="accent-blue-500" />
          {{ t('burn.cleanupIso') }}
        </label>
      </div>
      <p v-if="project.burnOptions.multisession" class="text-xs text-yellow-500 mt-1">
        {{ t('burn.multisessionHint') }}
      </p>
    </div>

    <!-- Очистка диска -->
    <div class="space-y-2 pt-2 border-t border-gray-300 dark:border-gray-700">
      <h3 class="text-sm font-medium text-gray-600 dark:text-gray-400">{{ t('burn.blankDisc') }}</h3>
      <div class="flex items-center gap-3">
        <select
          v-model="blankMode"
          class="px-2 py-1.5 text-sm bg-white dark:bg-gray-900 border border-gray-400 dark:border-gray-600 rounded text-gray-800 dark:text-gray-200 focus:outline-none focus:border-blue-500"
        >
          <option value="fast">{{ t('burn.fastBlank') }}</option>
          <option value="full">{{ t('burn.fullBlank') }}</option>
          <option value="deformat">{{ t('burn.deformat') }}</option>
        </select>
        <button
          @click="emit('blank-disc', blankMode)"
          :disabled="!currentDevicePath || isBurning"
          class="px-3 py-1.5 text-sm font-medium rounded bg-yellow-600 hover:bg-yellow-500 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
        >
          {{ t('burn.blank') }}
        </button>
      </div>
    </div>

    <!-- Кнопки действий -->
    <div class="flex justify-end gap-3 pt-2">
      <button
        @click="emit('go-back')"
        class="px-4 py-2 text-sm font-medium rounded bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
      >
        {{ t('burn.cancel') }}
      </button>
      <button
        @click="emit('create-iso')"
        :disabled="!project?.entries?.length || isBurning"
        class="px-5 py-2 text-sm font-semibold rounded bg-blue-600 hover:bg-blue-500 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
      >
        {{ t('burn.createIso') }}
      </button>
      <button
        @click="emit('start-burn')"
        :disabled="!canBurn"
        class="px-6 py-2 text-sm font-semibold rounded bg-orange-600 hover:bg-orange-500 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
      >
        {{ t('burn.startBurn') }}
      </button>
    </div>
  </div>
</template>
