<script setup>
import { useI18n } from 'vue-i18n'
import { useThemeStore } from '../../stores/themeStore'
import { useDeviceStore } from '../../stores/deviceStore'
import Button from '../ui/Button.vue'

const { t } = useI18n()
const themeStore = useThemeStore()
const deviceStore = useDeviceStore()

const emit = defineEmits([
  'new-project',
  'open-project',
  'save-project',
  'burn',
])
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
      <Button variant="ghost" size="sm" @click="emit('new-project')">
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        {{ t('header.new') }}
      </Button>
      <Button variant="ghost" size="sm" @click="emit('open-project')">
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
        </svg>
        {{ t('header.open') }}
      </Button>
      <Button variant="ghost" size="sm" @click="emit('save-project')">
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4" />
        </svg>
        {{ t('header.save') }}
      </Button>

      <div class="w-px h-6 bg-gray-300 dark:bg-gray-700 mx-2" />

      <Button variant="primary" size="sm" @click="emit('burn')">
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7C14 5 16.09 5.777 17.656 7.343A7.975 7.975 0 0120 13a7.975 7.975 0 01-2.343 5.657z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M9.879 16.121A3 3 0 1012.015 11L11 14H9c0 .768.293 1.536.879 2.121z" />
        </svg>
        {{ t('header.burn') }}
      </Button>
    </div>

    <!-- Spacer -->
    <div class="flex-1" />

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

    <!-- Settings -->
    <router-link
      to="/settings"
      class="p-2 rounded hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
    >
      <svg class="w-5 h-5 text-gray-600 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
      </svg>
    </router-link>

    <!-- Device selector — connected to deviceStore -->
    <div class="relative">
      <select
        :value="deviceStore.currentDevicePath || ''"
        @change="deviceStore.selectDevice($event.target.value)"
        class="appearance-none bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 text-sm rounded px-3 py-1.5 pr-8 border border-gray-400 dark:border-gray-600 hover:border-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent cursor-pointer min-w-[200px]"
      >
        <option value="" disabled>{{ t('device.selectDevice') }}</option>
        <option
          v-for="dev in deviceStore.devices"
          :key="dev.path"
          :value="dev.path"
        >
          {{ dev.vendor }} {{ dev.model }} ({{ dev.path }})
        </option>
      </select>
      <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
        <svg class="w-4 h-4 text-gray-600 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </div>
    </div>
  </header>
</template>
