<script setup>
import { Dialogs } from '@wailsio/runtime'
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import { useThemeStore } from '../../stores/themeStore'
import { useTabStore } from '../../stores/tabStore'
import { useProjectStore } from '../../stores/projectStore'
import { SaveSettings, GetSettings } from '../../../bindings/xorriso-ui/services/settingsservice.js'
import { Languages } from 'lucide-vue-next'
import Button from '../ui/Button.vue'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()
const themeStore = useThemeStore()
const tabStore = useTabStore()
const projectStore = useProjectStore()

import { ref } from 'vue'

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

async function newProject() {
  const name = t('tabs.untitledProject')
  const tabId = tabStore.addProjectTab(name, 'UNTITLED')
  await projectStore.newProject(tabId, name, 'UNTITLED')
}

async function openProject() {
  const filePath = await Dialogs.OpenFile({
    Title: t('header.open'),
    Filters: [{ DisplayName: 'Xorriso Project', Pattern: '*.xorriso-project' }],
    CanChooseFiles: true,
  })
  if (!filePath) return
  const tabId = tabStore.addProjectTab('Loading...', '')
  await projectStore.openProject(tabId, filePath)
}

async function saveProject() {
  const tabId = tabStore.activeTabId
  const data = tabStore.getProjectData(tabId)
  if (!data) return

  if (data.filePath) {
    await projectStore.saveProject(tabId)
  } else {
    const filePath = await Dialogs.SaveFile({
      Title: t('header.save'),
      Filename: (data.name || 'project') + '.xorriso-project',
      Filters: [{ DisplayName: 'Xorriso Project', Pattern: '*.xorriso-project' }],
    })
    if (!filePath) return
    await projectStore.saveProjectAs(tabId, filePath)
  }
}

function toggleDiscInfo() {
  tabStore.toggleDiscInfo()
}
</script>

<template>
  <header class="bg-gray-100 dark:bg-gray-800 border-b border-gray-300 dark:border-gray-700 px-4 py-2 flex items-center gap-4">
    <!-- App name -->
    <div class="flex items-center gap-2 mr-4">
      <svg class="w-6 h-6 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <circle cx="12" cy="12" r="10" stroke-width="1.5" />
        <circle cx="12" cy="12" r="3" stroke-width="1.5" />
        <circle cx="12" cy="12" r="6" stroke-width="0.5" opacity="0.5" />
      </svg>
      <span class="text-lg font-bold text-gray-900 dark:text-gray-100 whitespace-nowrap">xorriso-ui</span>
    </div>

    <!-- Toolbar buttons -->
    <div class="flex items-center gap-1">
      <Button variant="ghost" size="sm" @click="newProject">
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        {{ t('header.new') }}
      </Button>
      <Button variant="ghost" size="sm" @click="openProject">
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
        </svg>
        {{ t('header.open') }}
      </Button>
      <Button variant="ghost" size="sm" @click="saveProject">
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4" />
        </svg>
        {{ t('header.save') }}
      </Button>

      <div class="w-px h-6 bg-gray-300 dark:bg-gray-700 mx-2" />

      <!-- Disc Info toggle -->
      <Button
        :variant="tabStore.showDiscInfo ? 'primary' : 'ghost'"
        size="sm"
        @click="toggleDiscInfo"
      >
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="10" stroke-width="1.5" />
          <circle cx="12" cy="12" r="3" stroke-width="1.5" />
        </svg>
        {{ t('tabs.discInfo') }}
      </Button>
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

    <!-- Info toggle -->
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
