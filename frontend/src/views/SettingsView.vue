<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { GetToolsInfo, GetAppVersion } from '../../bindings/xorriso-ui/services/settingsservice.js'

const { t } = useI18n()
const router = useRouter()

const toolsInfo = ref([])
const appVersion = ref('')

const hotkeys = [
  { keys: ['Arrow Up', 'Arrow Down'], action: 'info.hotkeyNavigate' },
  { keys: ['Arrow Right'], action: 'info.hotkeyExpand' },
  { keys: ['Arrow Left'], action: 'info.hotkeyCollapse' },
  { keys: ['Space'], action: 'info.hotkeyToggleSelect' },
  { keys: ['Enter'], action: 'info.hotkeyAddToProject' },
  { keys: ['Delete'], action: 'info.hotkeyRemove' },
  { keys: ['Escape'], action: 'info.hotkeyDeselectAll' },
  { keys: ['Ctrl + L'], action: 'info.hotkeyEditPath' },
  { keys: ['Tab'], action: 'info.hotkeyFocusPanel' },
]

function goBack() {
  router.push('/')
}

async function loadToolsInfo() {
  try {
    const result = await GetToolsInfo()
    if (result) toolsInfo.value = result
  } catch (e) {
    console.error('Failed to get tools info:', e)
  }
}

async function loadAppVersion() {
  try {
    const result = await GetAppVersion()
    if (result) appVersion.value = result
  } catch (e) {
    appVersion.value = 'unknown'
  }
}

loadToolsInfo()
loadAppVersion()
</script>

<template>
  <div class="h-full overflow-y-auto">
    <div class="max-w-2xl mx-auto p-6 space-y-8">

      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-xl font-semibold">{{ t('info.title') }}</h1>
          <p class="text-sm text-gray-500 mt-1">{{ t('info.subtitle') }}</p>
        </div>
        <button
          @click="goBack"
          class="p-1.5 rounded hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
        >
          <svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- App version -->
      <section class="space-y-3">
        <h2 class="text-sm font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">{{ t('info.appVersion') }}</h2>
        <div class="flex items-center gap-3">
          <svg class="w-5 h-5 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <circle cx="12" cy="12" r="10" stroke-width="1.5" />
            <circle cx="12" cy="12" r="3" stroke-width="1.5" />
            <circle cx="12" cy="12" r="6" stroke-width="0.5" opacity="0.5" />
          </svg>
          <span class="text-sm font-medium text-gray-700 dark:text-gray-300">xorriso-ui</span>
          <span class="text-sm text-gray-500">{{ appVersion || '...' }}</span>
        </div>
      </section>

      <!-- External tools -->
      <section class="space-y-3">
        <h2 class="text-sm font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">{{ t('settings.externalTools') }}</h2>
        <div class="space-y-2">
          <div v-for="tool in toolsInfo" :key="tool.name" class="flex items-center gap-3">
            <span
              class="inline-flex items-center justify-center w-5 h-5 rounded-full text-xs font-bold"
              :class="tool.found ? 'bg-green-500/20 text-green-400' : 'bg-red-500/20 text-red-400'"
            >
              <svg v-if="tool.found" class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7" />
              </svg>
              <svg v-else class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </span>
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ tool.name }}</span>
            <span v-if="tool.found" class="text-xs text-gray-500">{{ tool.version }}</span>
            <span v-else class="text-xs text-red-400">{{ t('settings.toolNotFound') }}</span>
          </div>
        </div>
      </section>

      <!-- Hotkeys -->
      <section class="space-y-3 border-t border-gray-300 dark:border-gray-700 pt-6">
        <h2 class="text-sm font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">{{ t('info.hotkeys') }}</h2>
        <div class="space-y-2">
          <div
            v-for="(hk, i) in hotkeys"
            :key="i"
            class="flex items-center gap-3"
          >
            <div class="flex gap-1 shrink-0" style="min-width: 180px;">
              <kbd
                v-for="key in hk.keys"
                :key="key"
                class="px-2 py-0.5 text-xs font-mono rounded bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 border border-gray-300 dark:border-gray-600"
              >{{ key }}</kbd>
            </div>
            <span class="text-sm text-gray-600 dark:text-gray-400">{{ t(hk.action) }}</span>
          </div>
        </div>
      </section>

    </div>
  </div>
</template>
