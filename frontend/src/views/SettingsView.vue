<script setup>
import { ref, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { useThemeStore } from '../stores/themeStore'
import { GetSettings, SaveSettings } from '../../bindings/xorriso-ui/services/settingsservice.js'

const { t, locale } = useI18n()
const router = useRouter()
const themeStore = useThemeStore()

// --- State ---
const saving = ref(false)
const saved = ref(false)
const saveError = ref('')

const settings = reactive({
  language: locale.value,
  theme: themeStore.currentTheme,
  xorrisoBinaryPath: '/usr/bin/xorriso',
  devicePollInterval: 5,
  bdxlSafeMode: false,
  defaultBurnOptions: {
    speed: 'auto',
    dummyMode: false,
    verify: true,
    closeDisc: false,
    eject: true,
    burnMode: 'auto',
    streamRecording: false,
  },
  defaultIsoOptions: {
    udf: true,
    rockRidge: false,
    joliet: false,
    md5: true,
    backupMode: false,
  },
})

// BDXL safe mode auto-toggle
watch(() => settings.bdxlSafeMode, (enabled) => {
  if (enabled) {
    settings.defaultIsoOptions.md5 = true
    settings.defaultBurnOptions.verify = true
    settings.defaultBurnOptions.streamRecording = true
  }
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
      settings.xorrisoBinaryPath = result.xorrisoPath || '/usr/bin/xorriso'
      settings.devicePollInterval = result.devicePollInterval || 5
      settings.bdxlSafeMode = result.bdxlSafeMode || false
      settings.language = result.language || 'en'
      settings.theme = result.theme || 'dark'

      if (result.defaultBurn) {
        Object.assign(settings.defaultBurnOptions, result.defaultBurn)
      }
      if (result.defaultIso) {
        Object.assign(settings.defaultIsoOptions, result.defaultIso)
      }
    }
  } catch (error) {
    console.error('Failed to load settings:', error)
  }
}

async function saveSettings() {
  saving.value = true
  saved.value = false
  saveError.value = ''
  try {
    await SaveSettings({
      xorrisoPath: settings.xorrisoBinaryPath,
      defaultBurn: { ...settings.defaultBurnOptions },
      defaultIso: { ...settings.defaultIsoOptions },
      bdxlSafeMode: settings.bdxlSafeMode,
      autoEject: settings.defaultBurnOptions.eject,
      devicePollInterval: settings.devicePollInterval,
      language: settings.language,
      theme: settings.theme,
    })
    saved.value = true
    setTimeout(() => { saved.value = false }, 3000)
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

// Load on mount
loadSettings()
</script>

<template>
  <div class="h-full overflow-y-auto">
    <div class="max-w-2xl mx-auto p-6 space-y-8">

      <div class="flex items-center gap-3">
        <button
          @click="goBack"
          class="p-1.5 rounded hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors text-gray-600 dark:text-gray-400"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <div>
          <h1 class="text-xl font-semibold">{{ t('settings.title') }}</h1>
          <p class="text-sm text-gray-500 mt-1">{{ t('settings.subtitle') }}</p>
        </div>
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

        <div>
          <label class="block text-sm text-gray-700 dark:text-gray-300 mb-1">{{ t('settings.xorrisoBinaryPath') }}</label>
          <input
            v-model="settings.xorrisoBinaryPath"
            type="text"
            placeholder="/usr/bin/xorriso"
            class="w-full px-3 py-2 text-sm bg-gray-100 dark:bg-gray-800 border border-gray-400 dark:border-gray-600 rounded text-gray-900 dark:text-gray-200 placeholder-gray-500 dark:placeholder-gray-600 focus:outline-none focus:border-blue-500"
          />
          <p class="text-xs text-gray-500 dark:text-gray-600 mt-1">{{ t('settings.xorrisoBinaryPathHelp') }}</p>
        </div>

        <div>
          <label class="block text-sm text-gray-700 dark:text-gray-300 mb-1">{{ t('settings.devicePollInterval') }}</label>
          <input
            v-model.number="settings.devicePollInterval"
            type="number"
            min="1"
            max="60"
            class="w-32 px-3 py-2 text-sm bg-gray-100 dark:bg-gray-800 border border-gray-400 dark:border-gray-600 rounded text-gray-900 dark:text-gray-200 focus:outline-none focus:border-blue-500"
          />
          <p class="text-xs text-gray-500 dark:text-gray-600 mt-1">{{ t('settings.devicePollIntervalHelp') }}</p>
        </div>
      </section>

      <!-- Default Filesystem Options -->
      <section class="space-y-3">
        <h2 class="text-sm font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">{{ t('settings.defaultFsOptions') }}</h2>

        <div class="space-y-2">
          <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-gray-300">
            <input type="checkbox" v-model="settings.defaultIsoOptions.udf" class="accent-blue-500" />
            {{ t('settings.udf') }}
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-gray-300">
            <input type="checkbox" v-model="settings.defaultIsoOptions.rockRidge" class="accent-blue-500" />
            {{ t('settings.rockRidge') }}
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-gray-300">
            <input type="checkbox" v-model="settings.defaultIsoOptions.joliet" class="accent-blue-500" />
            {{ t('settings.joliet') }}
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-gray-300">
            <input type="checkbox" v-model="settings.defaultIsoOptions.md5" class="accent-blue-500" />
            {{ t('settings.md5') }}
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-gray-300">
            <input type="checkbox" v-model="settings.defaultIsoOptions.backupMode" class="accent-blue-500" />
            {{ t('settings.backupMode') }}
          </label>
        </div>
      </section>

      <!-- Default Burn Options -->
      <section class="space-y-3">
        <h2 class="text-sm font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">{{ t('settings.defaultBurnOptions') }}</h2>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm text-gray-700 dark:text-gray-300 mb-1">{{ t('settings.defaultSpeed') }}</label>
            <select
              v-model="settings.defaultBurnOptions.speed"
              class="w-full px-2 py-1.5 text-sm bg-gray-100 dark:bg-gray-800 border border-gray-400 dark:border-gray-600 rounded text-gray-900 dark:text-gray-200 focus:outline-none focus:border-blue-500"
            >
              <option value="auto">{{ t('settings.speedAuto') }}</option>
              <option value="1">{{ t('settings.speed1x') }}</option>
              <option value="2">{{ t('settings.speed2x') }}</option>
              <option value="4">{{ t('settings.speed4x') }}</option>
              <option value="8">{{ t('settings.speed8x') }}</option>
              <option value="12">{{ t('settings.speed12x') }}</option>
              <option value="16">{{ t('settings.speed16x') }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm text-gray-700 dark:text-gray-300 mb-1">{{ t('settings.defaultBurnMode') }}</label>
            <select
              v-model="settings.defaultBurnOptions.burnMode"
              class="w-full px-2 py-1.5 text-sm bg-gray-100 dark:bg-gray-800 border border-gray-400 dark:border-gray-600 rounded text-gray-900 dark:text-gray-200 focus:outline-none focus:border-blue-500"
            >
              <option value="auto">{{ t('burn.autoDao') }}</option>
              <option value="tao">{{ t('burn.tao') }}</option>
              <option value="sao">{{ t('burn.saoDao') }}</option>
            </select>
          </div>
        </div>

        <div class="space-y-2 mt-2">
          <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-gray-300">
            <input type="checkbox" v-model="settings.defaultBurnOptions.verify" class="accent-blue-500" />
            {{ t('burn.verifyAfterBurn') }}
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-gray-300">
            <input type="checkbox" v-model="settings.defaultBurnOptions.eject" class="accent-blue-500" />
            {{ t('burn.ejectWhenDone') }}
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-gray-300">
            <input type="checkbox" v-model="settings.defaultBurnOptions.dummyMode" class="accent-yellow-500" />
            {{ t('burn.simulationMode') }}
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-gray-300">
            <input type="checkbox" v-model="settings.defaultBurnOptions.closeDisc" class="accent-blue-500" />
            {{ t('burn.closeDisc') }}
          </label>
          <label class="flex items-center gap-3 text-sm text-gray-700 dark:text-gray-300">
            <input type="checkbox" v-model="settings.defaultBurnOptions.streamRecording" class="accent-blue-500" />
            {{ t('burn.streamRecording') }}
          </label>
        </div>
      </section>

      <!-- BDXL Safe Mode -->
      <section class="space-y-3">
        <h2 class="text-sm font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">{{ t('settings.bluray') }}</h2>

        <label class="flex items-start gap-3 text-sm text-gray-700 dark:text-gray-300">
          <input type="checkbox" v-model="settings.bdxlSafeMode" class="accent-cyan-500 mt-0.5" />
          <div>
            <span class="font-medium">{{ t('settings.bdxlSafeMode') }}</span>
            <p class="text-xs text-gray-500 mt-0.5">{{ t('settings.bdxlSafeModeDescription') }}</p>
          </div>
        </label>
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
        <span v-if="saved" class="text-sm text-green-400">{{ t('settings.savedSuccessfully') }}</span>
        <span v-if="saveError" class="text-sm text-red-400">{{ saveError }}</span>
      </div>
    </div>
  </div>
</template>
