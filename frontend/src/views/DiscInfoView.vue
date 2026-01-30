<script setup>
import { useI18n } from 'vue-i18n'
import { useDeviceStore } from '../stores/deviceStore'
import { useProjectStore } from '../stores/projectStore'

const { t } = useI18n()
const deviceStore = useDeviceStore()
const projectStore = useProjectStore()

const emit = defineEmits(['close'])

function formatBytes(bytes) {
  return projectStore.formatBytes(bytes)
}
</script>

<template>
  <div class="h-full overflow-y-auto p-6">
    <div class="max-w-3xl mx-auto space-y-6">

      <!-- Header with close button -->
      <div class="flex items-center justify-between">
        <h2 class="text-sm font-semibold text-gray-900 dark:text-gray-100">{{ t('tabs.discInfo') }}</h2>
        <button
          @click="emit('close')"
          class="p-1.5 rounded hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
        >
          <svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Device Selector + Actions -->
      <div class="flex items-center gap-3 flex-wrap">
        <label class="text-sm font-medium text-gray-600 dark:text-gray-400 shrink-0">
          {{ t('device.device') }}:
        </label>
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
        <button
          @click="deviceStore.fetchMediaInfo()"
          class="px-3 py-1.5 text-xs font-medium rounded bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
        >
          {{ t('device.refreshMedia') }}
        </button>
        <button
          @click="deviceStore.ejectDisc()"
          :disabled="!deviceStore.currentDevicePath"
          class="px-3 py-1.5 text-xs font-medium rounded bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 disabled:opacity-40 transition-colors"
        >
          {{ t('device.ejectDisc') }}
        </button>
        <button
          v-if="deviceStore.currentDevice?.canCloseTray"
          @click="deviceStore.loadTray()"
          :disabled="!deviceStore.currentDevicePath"
          class="px-3 py-1.5 text-xs font-medium rounded bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 disabled:opacity-40 transition-colors"
        >
          {{ t('device.loadTray') }}
        </button>
      </div>

      <!-- No device selected -->
      <div v-if="!deviceStore.currentDevice" class="text-center text-gray-500 dark:text-gray-600 py-12">
        <svg class="w-16 h-16 mx-auto mb-3 text-gray-300 dark:text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="10" stroke-width="1" />
          <circle cx="12" cy="12" r="3" stroke-width="1" />
          <circle cx="12" cy="12" r="6" stroke-width="0.5" opacity="0.5" />
        </svg>
        <p class="text-sm">{{ t('device.noDeviceSelected') }}</p>
      </div>

      <template v-else>
        <!-- Drive Info -->
        <div class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
          <h3 class="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wider mb-3">{{ t('device.driveInfo') }}</h3>
          <div class="grid grid-cols-[auto_1fr] gap-x-4 gap-y-1.5 text-sm text-gray-700 dark:text-gray-300">
            <span class="text-gray-500">{{ t('device.device') }}:</span>
            <span>{{ deviceStore.currentDevice.vendor }} {{ deviceStore.currentDevice.model }}</span>
            <span class="text-gray-500">Path:</span>
            <span class="font-mono text-xs">{{ deviceStore.currentDevice.path }}</span>
            <template v-if="deviceStore.currentDevice.linkPath && deviceStore.currentDevice.linkPath !== deviceStore.currentDevice.path">
              <span class="text-gray-500">Link:</span>
              <span class="font-mono text-xs">{{ deviceStore.currentDevice.linkPath }}</span>
            </template>
            <template v-if="deviceStore.currentDevice.revision">
              <span class="text-gray-500">Firmware:</span>
              <span>{{ deviceStore.currentDevice.revision }}</span>
            </template>
            <template v-if="deviceStore.currentDevice.driveSpeed">
              <span class="text-gray-500">{{ t('device.driveSpeed') }}:</span>
              <span>{{ deviceStore.currentDevice.driveSpeed }}x</span>
            </template>
          </div>
        </div>

        <!-- Media Info -->
        <div v-if="deviceStore.mediaInfo" class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
          <h3 class="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wider mb-3">{{ t('device.mediaInfo') }}</h3>
          <div class="grid grid-cols-[auto_1fr] gap-x-4 gap-y-1.5 text-sm text-gray-700 dark:text-gray-300">
            <span class="text-gray-500">{{ t('device.type') }}:</span>
            <span>{{ deviceStore.mediaInfo.mediaType || '—' }}</span>
            <span class="text-gray-500">{{ t('device.status') }}:</span>
            <span>{{ deviceStore.mediaInfo.mediaStatus || '—' }}</span>
            <template v-if="deviceStore.mediaInfo.mediaProduct">
              <span class="text-gray-500">{{ t('device.manufacturer') }}:</span>
              <span>{{ deviceStore.mediaInfo.mediaProduct }}</span>
            </template>
            <span class="text-gray-500">{{ t('device.capacity') }}:</span>
            <span>{{ formatBytes(deviceStore.mediaInfo.totalCapacity) }}</span>
            <span class="text-gray-500">{{ t('device.free') }}:</span>
            <span>{{ formatBytes(deviceStore.mediaInfo.freeSpace) }}</span>
            <span class="text-gray-500">{{ t('common.used') }}:</span>
            <span>{{ formatBytes(deviceStore.mediaInfo.usedSpace) }}</span>
            <span class="text-gray-500">{{ t('device.sessions') }}:</span>
            <span>{{ deviceStore.mediaInfo.sessions }}</span>
            <span class="text-gray-500">{{ t('device.erasable') }}:</span>
            <span>{{ deviceStore.mediaInfo.erasable ? t('device.yes') : t('device.no') }}</span>
          </div>
        </div>

        <div v-else class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 border border-gray-200 dark:border-gray-700 text-center text-gray-500 dark:text-gray-600 py-8">
          <p class="text-sm">{{ t('device.noMediaInserted') }}</p>
        </div>

        <!-- Volume Info (PVD) -->
        <div v-if="deviceStore.mediaInfo && (deviceStore.mediaInfo.volumeId || deviceStore.mediaInfo.publisherId || deviceStore.mediaInfo.appId)"
          class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
          <h3 class="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wider mb-3">{{ t('device.volumeInfo') }}</h3>
          <div class="grid grid-cols-[auto_1fr] gap-x-4 gap-y-1.5 text-sm text-gray-700 dark:text-gray-300">
            <template v-if="deviceStore.mediaInfo.volumeId">
              <span class="text-gray-500">{{ t('device.volumeId') }}:</span>
              <span>{{ deviceStore.mediaInfo.volumeId }}</span>
            </template>
            <template v-if="deviceStore.mediaInfo.volumeSetId">
              <span class="text-gray-500">Volume Set:</span>
              <span>{{ deviceStore.mediaInfo.volumeSetId }}</span>
            </template>
            <template v-if="deviceStore.mediaInfo.publisherId">
              <span class="text-gray-500">{{ t('device.publisher') }}:</span>
              <span>{{ deviceStore.mediaInfo.publisherId }}</span>
            </template>
            <template v-if="deviceStore.mediaInfo.preparerId">
              <span class="text-gray-500">{{ t('device.preparer') }}:</span>
              <span>{{ deviceStore.mediaInfo.preparerId }}</span>
            </template>
            <template v-if="deviceStore.mediaInfo.appId">
              <span class="text-gray-500">{{ t('device.application') }}:</span>
              <span class="break-words">{{ deviceStore.mediaInfo.appId }}</span>
            </template>
            <template v-if="deviceStore.mediaInfo.systemId">
              <span class="text-gray-500">{{ t('device.system') }}:</span>
              <span>{{ deviceStore.mediaInfo.systemId }}</span>
            </template>
            <template v-if="deviceStore.mediaInfo.creationTime">
              <span class="text-gray-500">{{ t('device.created') }}:</span>
              <span class="font-mono text-xs">{{ deviceStore.mediaInfo.creationTime }}</span>
            </template>
            <template v-if="deviceStore.mediaInfo.modifyTime">
              <span class="text-gray-500">{{ t('device.modified') }}:</span>
              <span class="font-mono text-xs">{{ deviceStore.mediaInfo.modifyTime }}</span>
            </template>
          </div>
        </div>

        <!-- Supported Profiles -->
        <div v-if="deviceStore.currentDevice?.profiles?.length"
          class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
          <h3 class="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wider mb-3">{{ t('device.driveInfo') }} — Profiles</h3>
          <div class="flex flex-wrap gap-1.5">
            <span
              v-for="profile in deviceStore.currentDevice.profiles"
              :key="profile.name"
              class="px-2 py-0.5 rounded text-xs font-medium"
              :class="profile.current
                ? 'bg-blue-500/20 text-blue-400 border border-blue-500/30'
                : 'bg-gray-200 dark:bg-gray-700 text-gray-600 dark:text-gray-400'"
            >
              {{ profile.name }}
            </span>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
