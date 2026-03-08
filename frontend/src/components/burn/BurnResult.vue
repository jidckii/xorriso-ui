<script setup>
import { useI18n } from 'vue-i18n'

defineProps({
  job: { type: Object, default: null },
  logLines: { type: Array, default: () => [] },
})

const emit = defineEmits(['go-back', 'burn-again'])

const { t } = useI18n()
</script>

<template>
  <div class="p-6 space-y-4">
    <div class="text-center py-4">
      <!-- Успех -->
      <template v-if="job?.result?.success">
        <svg class="w-16 h-16 mx-auto text-green-400 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h3 class="text-lg font-semibold text-green-400">{{ t('burn.burnComplete') }}</h3>
      </template>

      <!-- Ошибка -->
      <template v-else>
        <svg class="w-16 h-16 mx-auto text-red-400 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h3 class="text-lg font-semibold text-red-400">{{ t('burn.burnFailed') }}</h3>
      </template>

      <p class="text-sm text-gray-600 dark:text-gray-400 mt-2">{{ job?.result?.message }}</p>

      <!-- Результаты верификации -->
      <div v-if="job?.result?.success" class="mt-4 space-y-1">
        <div v-if="job.result.verifyErrors === 0" class="text-sm font-medium text-green-400">
          {{ t('burn.verificationPassed') }}
        </div>
        <div v-else-if="job.result.verifyErrors > 0" class="text-sm font-medium text-red-400">
          {{ t('burn.verificationFailed', { count: job.result.verifyErrors }) }}
        </div>
        <div v-if="job.result.md5Match !== undefined" class="text-xs text-gray-500">
          MD5: <span :class="job.result.md5Match ? 'text-green-400' : 'text-red-400'">
            {{ job.result.md5Match ? t('burn.md5Match') : t('burn.md5Mismatch') }}
          </span>
        </div>
      </div>
    </div>

    <!-- Лог -->
    <div class="bg-white dark:bg-gray-900 rounded p-3 max-h-40 overflow-y-auto font-mono text-xs text-gray-600 dark:text-gray-400">
      <div v-for="(line, i) in logLines" :key="i">{{ line }}</div>
    </div>

    <!-- Кнопки действий -->
    <div class="flex justify-end gap-3">
      <button
        @click="emit('go-back')"
        class="px-4 py-2 text-sm font-medium rounded bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
      >
        {{ t('burn.backToProject') }}
      </button>
      <button
        @click="emit('burn-again')"
        class="px-4 py-2 text-sm font-medium rounded bg-orange-600 hover:bg-orange-500 transition-colors"
      >
        {{ t('burn.burnAnother') }}
      </button>
    </div>
  </div>
</template>
