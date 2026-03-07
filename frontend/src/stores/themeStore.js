import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { GetSettings, SaveSettings } from '../../bindings/xorriso-ui/services/settingsservice.js'

export const useThemeStore = defineStore('theme', () => {
  const isDark = ref(true)

  const currentTheme = computed(() => isDark.value ? 'dark' : 'light')

  function detectSystemTheme() {
    if (typeof window === 'undefined') return 'dark'
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
  }

  async function loadThemeFromSettings() {
    // Сначала применяем из localStorage для мгновенного отображения
    const cached = localStorage.getItem('xorriso-theme')
    if (cached) {
      isDark.value = cached === 'dark'
    }

    // Затем загружаем из Go settings
    try {
      const settings = await GetSettings()
      if (settings?.theme) {
        isDark.value = settings.theme === 'dark'
        localStorage.setItem('xorriso-theme', settings.theme)
      } else if (!cached) {
        isDark.value = detectSystemTheme() === 'dark'
      }
    } catch (error) {
      console.error('Failed to load theme from settings:', error)
      if (!cached) {
        isDark.value = detectSystemTheme() === 'dark'
      }
    }
  }

  async function setTheme(dark) {
    isDark.value = dark
    const theme = dark ? 'dark' : 'light'
    localStorage.setItem('xorriso-theme', theme)

    // Сохраняем тему в Go settings
    try {
      const settings = await GetSettings()
      if (settings) {
        settings.theme = theme
        await SaveSettings(settings)
      }
    } catch (error) {
      console.error('Failed to save theme to settings:', error)
    }
  }

  function toggleTheme() {
    setTheme(!isDark.value)
  }

  watch(isDark, (dark) => {
    if (typeof document !== 'undefined') {
      document.documentElement.classList.toggle('dark', dark)
    }
  }, { immediate: true })

  return {
    isDark,
    currentTheme,
    loadThemeFromSettings,
    setTheme,
    toggleTheme,
  }
})
