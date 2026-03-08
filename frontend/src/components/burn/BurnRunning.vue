<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { formatBytes } from '../../composables/useFormatBytes'

defineProps({
  progress: { type: Object, required: true },
  logLines: { type: Array, default: () => [] },
  phaseLabel: { type: String, default: '' },
})

const emit = defineEmits(['cancel'])

const { t } = useI18n()

const showLog = ref(false)
</script>

<template>
  <div class="p-6 space-y-4">
    <div class="text-center mb-4">
      <p class="text-sm text-gray-600 dark:text-gray-400">{{ phaseLabel }}</p>
      <p class="text-3xl font-bold text-orange-400 mt-1">{{ progress.percent }}%</p>
    </div>

    <!-- Прогресс-бар -->
    <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
      <div
        class="h-full bg-orange-500 rounded-full transition-all duration-300"
        :style="{ width: progress.percent + '%' }"
      ></div>
    </div>

    <!-- Статистика -->
    <div class="grid grid-cols-3 gap-4 text-center text-sm">
      <div>
        <p class="text-gray-500 text-xs">{{ t('burnProgress.speed') }}</p>
        <p class="text-gray-700 dark:text-gray-300">{{ progress.speed || '-' }}</p>
      </div>
      <div>
        <p class="text-gray-500 text-xs">{{ t('burnProgress.written') }}</p>
        <p class="text-gray-700 dark:text-gray-300">
          {{ formatBytes(progress.bytesWritten) }} / {{ formatBytes(progress.bytesTotal) }}
        </p>
      </div>
      <div>
        <p class="text-gray-500 text-xs">{{ t('burnProgress.eta') }}</p>
        <p class="text-gray-700 dark:text-gray-300">{{ progress.eta || '-' }}</p>
      </div>
    </div>

    <!-- FIFO -->
    <div class="flex items-center gap-2 text-xs text-gray-500">
      <span>{{ t('burnProgress.fifo') }}:</span>
      <div class="flex-1 h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
        <div
          class="h-full bg-green-500 rounded-full transition-all"
          :style="{ width: (progress.fifoFill || 0) + '%' }"
        ></div>
      </div>
      <span>{{ (progress.fifoFill || 0).toFixed(0) }}%</span>
    </div>

    <!-- Переключение лога -->
    <div>
      <button
        @click="showLog = !showLog"
        class="text-xs text-gray-500 hover:text-gray-600 dark:hover:text-gray-400 transition-colors"
      >
        {{ showLog ? t('burnProgress.hideLog') : t('burnProgress.showLog') }}
      </button>
      <div v-if="showLog" class="mt-2 bg-white dark:bg-gray-900 rounded p-3 max-h-40 overflow-y-auto font-mono text-xs text-gray-600 dark:text-gray-400">
        <div v-for="(line, i) in logLines" :key="i">{{ line }}</div>
        <div v-if="logLines.length === 0" class="text-gray-500 dark:text-gray-600">{{ t('burnProgress.noLogOutput') }}</div>
      </div>
    </div>

    <!-- Отмена -->
    <div class="flex justify-end">
      <button
        @click="emit('cancel')"
        class="px-4 py-2 text-sm font-medium rounded bg-red-600 hover:bg-red-500 transition-colors"
      >
        {{ t('burn.cancelBurn') }}
      </button>
    </div>
  </div>
</template>
