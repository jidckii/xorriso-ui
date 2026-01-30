import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  const isDark = ref(true)

  const currentTheme = computed(() => isDark.value ? 'dark' : 'light')

  function detectSystemTheme() {
    if (typeof window === 'undefined') return 'dark'
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
  }

  function loadThemeFromSettings() {
    // TODO: load from Wails GetSettings() when bindings are ready
    const saved = localStorage.getItem('xorriso-theme')
    if (saved) {
      isDark.value = saved === 'dark'
    } else {
      isDark.value = detectSystemTheme() === 'dark'
    }
  }

  function setTheme(dark) {
    isDark.value = dark
    // TODO: save via Wails SaveSettings() when bindings are ready
    localStorage.setItem('xorriso-theme', dark ? 'dark' : 'light')
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
