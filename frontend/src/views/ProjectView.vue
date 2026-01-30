<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useDeviceStore } from '../stores/deviceStore'
import { useProjectStore } from '../stores/projectStore'
import { useTabStore } from '../stores/tabStore'
import FileBrowser from '../components/project/FileBrowser.vue'
import DiscLayout from '../components/project/DiscLayout.vue'

const { t } = useI18n()
const deviceStore = useDeviceStore()
const projectStore = useProjectStore()
const tabStore = useTabStore()

const currentProject = computed(() => tabStore.activeProject)

// --- Capacity Bar ---
const capacityPercent = computed(() => {
  if (!currentProject.value || deviceStore.mediaCapacityBytes === 0) return 0
  return Math.min(100, (currentProject.value.totalSize / deviceStore.mediaCapacityBytes) * 100)
})

const capacityColor = computed(() => {
  if (capacityPercent.value > 100) return 'bg-red-500'
  if (capacityPercent.value > 90) return 'bg-yellow-500'
  return 'bg-blue-500'
})

function formatBytes(bytes) {
  return projectStore.formatBytes(bytes)
}
</script>

<template>
  <div v-if="currentProject" class="flex flex-col h-full">
    <!-- Top Panels: File Browser + Disc Layout -->
    <div class="flex-1 flex min-h-0">
      <!-- Left Panel: File Browser -->
      <div class="w-1/2 border-r border-gray-300 dark:border-gray-700">
        <FileBrowser />
      </div>

      <!-- Right Panel: Disc Layout -->
      <div class="w-1/2">
        <DiscLayout />
      </div>
    </div>

    <!-- Bottom: Capacity Bar with Device Selector -->
    <div class="px-3 py-2 bg-gray-100 dark:bg-gray-800 border-t border-gray-300 dark:border-gray-700">
      <div class="flex items-center gap-3">
        <!-- Inline Device Selector -->
        <div class="relative shrink-0">
          <select
            :value="deviceStore.currentDevicePath || ''"
            @change="deviceStore.selectDevice($event.target.value)"
            class="appearance-none bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 text-xs rounded px-2 py-1 pr-6 border border-gray-400 dark:border-gray-600 hover:border-gray-500 focus:outline-none focus:ring-1 focus:ring-blue-500 cursor-pointer max-w-[180px]"
          >
            <option value="" disabled>{{ t('device.selectDevice') }}</option>
            <option
              v-for="dev in deviceStore.devices"
              :key="dev.path"
              :value="dev.path"
            >
              {{ dev.vendor }} {{ dev.model }}
            </option>
          </select>
          <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-1">
            <svg class="w-3 h-3 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </div>
        </div>

        <span class="text-xs text-gray-600 dark:text-gray-400 whitespace-nowrap">{{ t('project.capacity') }}:</span>
        <div class="flex-1 h-3 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
          <div
            class="h-full rounded-full transition-all duration-300"
            :class="capacityColor"
            :style="{ width: Math.min(100, capacityPercent) + '%' }"
          ></div>
        </div>
        <span class="text-xs text-gray-600 dark:text-gray-400 whitespace-nowrap">
          {{ formatBytes(currentProject.totalSize) }} / {{ formatBytes(deviceStore.mediaCapacityBytes) }}
          <span v-if="capacityPercent > 0" class="ml-1">({{ capacityPercent.toFixed(1) }}%)</span>
        </span>
      </div>
      <div class="flex items-center gap-4 mt-1">
        <div class="flex items-center gap-1">
          <span class="w-3 h-1.5 rounded-full bg-blue-600 inline-block"></span>
          <span class="text-[10px] text-gray-500">{{ t('capacityBar.cd700') }}</span>
        </div>
        <div class="flex items-center gap-1">
          <span class="w-3 h-1.5 rounded-full bg-purple-600 inline-block"></span>
          <span class="text-[10px] text-gray-500">{{ t('capacityBar.dvd') }}</span>
        </div>
        <div class="flex items-center gap-1">
          <span class="w-3 h-1.5 rounded-full bg-cyan-600 inline-block"></span>
          <span class="text-[10px] text-gray-500">{{ t('capacityBar.bd') }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
