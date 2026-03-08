<script setup>
import { ref, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useThemeStore } from '../stores/themeStore'
import { GetSettings, SaveSettings, GetToolsInfo } from '../../bindings/xorriso-ui/services/settingsservice.js'

const { t, locale } = useI18n()
const router = useRouter()
const themeStore = useThemeStore()

// --- State ---
const saving = ref(false)
const saveError = ref('')
const toolsInfo = ref([])

const settings = reactive({
  language: locale.value,
  theme: themeStore.currentTheme,
})

// Apply language change immediately
watch(() => settings.language, (newLang) => {
  locale.value = newLang
  localStorage.setItem('xorriso-language', newLang)
})

// Apply theme change immediately
watch(() => settings.theme, (newTheme) => {
  themeStore.setTheme(newTheme === 'dark')
})

// --- Actions ---
async function loadSettings() {
  try {
    const result = await GetSettings()
    if (result) {
      settings.language = result.language || 'en'
      settings.theme = result.theme || 'dark'
    }
  } catch (error) {
    console.error('Failed to load settings:', error)
  }
}

async function saveSettings() {
  saving.value = true
  saveError.value = ''
  try {
    await SaveSettings({
      language: settings.language,
      theme: settings.theme,
    })
    goBack()
  } catch (error) {
    console.error('Failed to save settings:', error)
    saveError.value = String(error)
    setTimeout(() => { saveError.value = '' }, 5000)
  } finally {
    saving.value = false
  }
}

function goBack() {
  router.push('/')
}

async function loadToolsInfo() {
  try {
    const result = await GetToolsInfo()
    if (result) {
      toolsInfo.value = result
    }
  } catch (error) {
    console.error('Failed to get tools info:', error)
  }
}

// Load on mount
loadSettings()
loadToolsInfo()
</script>

<template>
  <div class="h-full overflow-y-auto">
    <div class="max-w-2xl mx-auto p-6 space-y-8">

      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-xl font-semibold">{{ t('settings.title') }}</h1>
          <p class="text-sm text-gray-500 mt-1">{{ t('settings.subtitle') }}</p>
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

      <!-- Appearance -->
      <section class="space-y-3">
        <h2 class="text-sm font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">{{ t('settings.appearance') }}</h2>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm text-gray-700 dark:text-gray-300 mb-1">{{ t('settings.language') }}</label>
            <select
              v-model="settings.language"
              class="w-full px-3 py-2 text-sm bg-gray-100 dark:bg-gray-800 border border-gray-400 dark:border-gray-600 rounded text-gray-900 dark:text-gray-200 focus:outline-none focus:border-blue-500"
            >
              <option value="en">English</option>
              <option value="ru">Русский</option>
            </select>
          </div>
          <div>
            <label class="block text-sm text-gray-700 dark:text-gray-300 mb-1">{{ t('settings.theme') }}</label>
            <select
              v-model="settings.theme"
              class="w-full px-3 py-2 text-sm bg-gray-100 dark:bg-gray-800 border border-gray-400 dark:border-gray-600 rounded text-gray-900 dark:text-gray-200 focus:outline-none focus:border-blue-500"
            >
              <option value="dark">{{ t('settings.dark') }}</option>
              <option value="light">{{ t('settings.light') }}</option>
            </select>
          </div>
        </div>
      </section>

      <!-- General -->
      <section class="space-y-3">
        <h2 class="text-sm font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">{{ t('settings.general') }}</h2>

        <!-- External tools status -->
        <div>
          <label class="block text-sm text-gray-700 dark:text-gray-300 mb-2">{{ t('settings.externalTools') }}</label>
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
              <span v-if="tool.found" class="text-xs text-gray-500 dark:text-gray-500">{{ tool.version }}</span>
              <span v-else class="text-xs text-red-400">{{ t('settings.toolNotFound') }}</span>
            </div>
          </div>
        </div>
      </section>

      <!-- Save -->
      <div class="flex items-center gap-3 pt-4 border-t border-gray-300 dark:border-gray-700">
        <button
          @click="saveSettings"
          :disabled="saving"
          class="px-6 py-2 text-sm font-semibold rounded bg-blue-600 hover:bg-blue-500 text-white disabled:opacity-50 transition-colors"
        >
          {{ saving ? t('settings.saving') : t('settings.saveSettings') }}
        </button>
        <span v-if="saveError" class="text-sm text-red-400">{{ saveError }}</span>
      </div>
    </div>
  </div>
</template>
