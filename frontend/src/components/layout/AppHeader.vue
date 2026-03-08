<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import { useThemeStore } from '../../stores/themeStore'
import { SaveSettings, GetSettings } from '../../../bindings/xorriso-ui/services/settingsservice.js'
import { Languages } from 'lucide-vue-next'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()
const themeStore = useThemeStore()

const langMenuOpen = ref(false)

const languages = [
  { code: 'en', label: 'English' },
  { code: 'ru', label: 'Русский' },
]

async function setLanguage(code) {
  locale.value = code
  langMenuOpen.value = false
  localStorage.setItem('xorriso-language', code)
  try {
    const settings = await GetSettings()
    if (settings) {
      settings.language = code
      await SaveSettings(settings)
    }
  } catch (e) {
    console.error('Failed to save language:', e)
  }
}
</script>

<template>
  <header class="bg-gray-100 dark:bg-gray-800 border-b border-gray-300 dark:border-gray-700 px-4 py-2 flex items-center gap-4">
    <!-- Логотип -->
    <div class="flex items-center gap-2 mr-4">
      <svg class="w-6 h-6 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <circle cx="12" cy="12" r="10" stroke-width="1.5" />
        <circle cx="12" cy="12" r="3" stroke-width="1.5" />
        <circle cx="12" cy="12" r="6" stroke-width="0.5" opacity="0.5" />
      </svg>
      <span class="text-lg font-bold text-gray-900 dark:text-gray-100 whitespace-nowrap">xorriso-ui</span>
    </div>

    <!-- Spacer -->
    <div class="flex-1" />

    <!-- Language selector -->
    <div class="relative">
      <button
        @click="langMenuOpen = !langMenuOpen"
        class="p-2 rounded hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors text-gray-600 dark:text-gray-400"
        :title="t('info.language')"
      >
        <Languages :size="20" />
      </button>
      <div
        v-if="langMenuOpen"
        class="absolute right-0 top-full mt-1 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded shadow-lg z-50 min-w-[120px] py-1"
        @mouseleave="langMenuOpen = false"
      >
        <button
          v-for="lang in languages"
          :key="lang.code"
          @click="setLanguage(lang.code)"
          class="w-full text-left px-3 py-1.5 text-sm hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
          :class="locale === lang.code ? 'text-blue-500 font-medium' : 'text-gray-700 dark:text-gray-300'"
        >
          {{ lang.label }}
        </button>
      </div>
    </div>

    <!-- Theme toggle -->
    <button
      @click="themeStore.toggleTheme()"
      class="p-2 rounded hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
    >
      <svg v-if="themeStore.isDark" class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
      </svg>
      <svg v-else class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
      </svg>
    </button>

    <!-- Settings button -->
    <button
      @click="route.path === '/settings' ? router.push('/') : router.push('/settings')"
      class="p-2 rounded transition-colors"
      :class="route.path === '/settings'
        ? 'bg-blue-600 text-white hover:bg-blue-500'
        : 'hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-400'"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <circle cx="12" cy="12" r="10" stroke-width="2" />
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 16v-4m0-4h.01" />
      </svg>
    </button>
  </header>
</template>
