<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useDeviceStore } from '../stores/deviceStore'
import { useProjectStore } from '../stores/projectStore'

const { t } = useI18n()
const router = useRouter()
const deviceStore = useDeviceStore()
const projectStore = useProjectStore()

// --- File Browser State ---
const browsePath = ref('/')
const browseEntries = ref([])
const selectedBrowseFiles = ref([])

// --- Disc Layout State ---
const selectedProjectEntries = ref([])

// --- Device Panel ---
const devicePanelOpen = ref(true)

// --- Init ---
onMounted(async () => {
  deviceStore.init()
  browseEntries.value = await projectStore.browseDirectory(browsePath.value)
})

// --- File Browser Actions ---
async function navigateTo(entry) {
  if (entry.isDir) {
    browsePath.value = entry.path
    selectedBrowseFiles.value = []
    browseEntries.value = await projectStore.browseDirectory(entry.path)
  }
}

function toggleBrowseSelection(entry) {
  const idx = selectedBrowseFiles.value.indexOf(entry.path)
  if (idx === -1) {
    selectedBrowseFiles.value.push(entry.path)
  } else {
    selectedBrowseFiles.value.splice(idx, 1)
  }
}

async function addSelectedToProject() {
  const paths = selectedBrowseFiles.value.filter(p => p !== browsePath.value)
  if (paths.length > 0) {
    await projectStore.addFiles(paths)
    selectedBrowseFiles.value = []
  }
}

// --- Disc Layout Actions ---
function toggleProjectSelection(entry) {
  const idx = selectedProjectEntries.value.indexOf(entry.destPath)
  if (idx === -1) {
    selectedProjectEntries.value.push(entry.destPath)
  } else {
    selectedProjectEntries.value.splice(idx, 1)
  }
}

async function removeSelectedFromProject() {
  for (const destPath of selectedProjectEntries.value) {
    await projectStore.removeEntry(destPath)
  }
  selectedProjectEntries.value = []
}

// --- Capacity Bar ---
const capacityPercent = computed(() => {
  if (deviceStore.mediaCapacityBytes === 0) return 0
  return Math.min(100, (projectStore.totalSize / deviceStore.mediaCapacityBytes) * 100)
})

const capacityColor = computed(() => {
  if (capacityPercent.value > 100) return 'bg-red-500'
  if (capacityPercent.value > 90) return 'bg-yellow-500'
  return 'bg-blue-500'
})

function formatBytes(bytes) {
  return projectStore.formatBytes(bytes)
}

// --- Burn Navigation ---
function goToBurn() {
  router.push('/burn')
}
</script>

<template>
  <div class="flex h-full">
    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col min-w-0">
      <!-- Top Panels: File Browser + Disc Layout -->
      <div class="flex-1 flex min-h-0">

        <!-- Left Panel: File Browser -->
        <div class="w-1/2 flex flex-col border-r border-gray-300 dark:border-gray-700">
          <div class="flex items-center gap-2 px-3 py-2 bg-gray-100 dark:bg-gray-800 border-b border-gray-300 dark:border-gray-700">
            <svg class="w-4 h-4 text-gray-600 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
            </svg>
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('project.fileBrowser') }}</span>
            <span class="ml-auto text-xs text-gray-500 truncate max-w-48">{{ browsePath }}</span>
          </div>

          <div class="flex-1 overflow-y-auto">
            <table class="w-full text-sm">
              <thead class="sticky top-0 bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400 text-xs uppercase">
                <tr>
                  <th class="px-3 py-1.5 text-left">{{ t('common.name') }}</th>
                  <th class="px-3 py-1.5 text-right w-24">{{ t('common.size') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="entry in browseEntries"
                  :key="entry.path"
                  class="border-b border-gray-100 dark:border-gray-800 hover:bg-gray-100/50 dark:hover:bg-gray-800/50 cursor-pointer"
                  :class="{ 'bg-blue-900/30': selectedBrowseFiles.includes(entry.path) }"
                  @click="toggleBrowseSelection(entry)"
                  @dblclick="navigateTo(entry)"
                >
                  <td class="px-3 py-1.5 flex items-center gap-2">
                    <svg v-if="entry.isDir" class="w-4 h-4 text-yellow-400 shrink-0" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" />
                    </svg>
                    <svg v-else class="w-4 h-4 text-gray-500 shrink-0" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z" clip-rule="evenodd" />
                    </svg>
                    <span class="truncate">{{ entry.name }}</span>
                  </td>
                  <td class="px-3 py-1.5 text-right text-gray-500 whitespace-nowrap">
                    {{ entry.isDir ? '' : formatBytes(entry.size) }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- File Browser Toolbar -->
          <div class="flex items-center gap-2 px-3 py-2 bg-gray-100 dark:bg-gray-800 border-t border-gray-300 dark:border-gray-700">
            <button
              @click="addSelectedToProject"
              :disabled="selectedBrowseFiles.length === 0"
              class="px-3 py-1 text-xs font-medium rounded bg-blue-600 hover:bg-blue-500 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >
              {{ t('project.addToProject') }}
            </button>
            <span class="text-xs text-gray-500">
              {{ selectedBrowseFiles.length }} {{ t('project.selected') }}
            </span>
          </div>
        </div>

        <!-- Right Panel: Disc Layout -->
        <div class="w-1/2 flex flex-col">
          <div class="flex items-center gap-2 px-3 py-2 bg-gray-100 dark:bg-gray-800 border-b border-gray-300 dark:border-gray-700">
            <svg class="w-4 h-4 text-gray-600 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8v2m0 8v2m9-6a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('project.discLayout') }}</span>
            <span class="ml-auto text-xs text-gray-500">
              {{ projectStore.entryCount }} {{ t('project.items') }} - {{ projectStore.totalSizeFormatted }}
            </span>
          </div>

          <div class="flex-1 overflow-y-auto">
            <div v-if="projectStore.project.entries.length === 0" class="flex items-center justify-center h-full text-gray-500 dark:text-gray-600">
              <div class="text-center">
                <svg class="w-12 h-12 mx-auto mb-2 text-gray-200 dark:text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                    d="M12 9v3m0 0v3m0-3h3m-3 0H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <p class="text-sm">{{ t('common.addFilesFromBrowser') }}</p>
              </div>
            </div>

            <table v-else class="w-full text-sm">
              <thead class="sticky top-0 bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-400 text-xs uppercase">
                <tr>
                  <th class="px-3 py-1.5 text-left">{{ t('common.destinationPath') }}</th>
                  <th class="px-3 py-1.5 text-right w-24">{{ t('common.size') }}</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="entry in projectStore.project.entries"
                  :key="entry.destPath"
                  class="border-b border-gray-100 dark:border-gray-800 hover:bg-gray-100/50 dark:hover:bg-gray-800/50 cursor-pointer"
                  :class="{ 'bg-blue-900/30': selectedProjectEntries.includes(entry.destPath) }"
                  @click="toggleProjectSelection(entry)"
                >
                  <td class="px-3 py-1.5 flex items-center gap-2">
                    <svg v-if="entry.isDir" class="w-4 h-4 text-yellow-400 shrink-0" fill="currentColor" viewBox="0 0 20 20">
                      <path d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" />
                    </svg>
                    <svg v-else class="w-4 h-4 text-gray-500 shrink-0" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z" clip-rule="evenodd" />
                    </svg>
                    <span class="truncate">{{ entry.destPath }}</span>
                  </td>
                  <td class="px-3 py-1.5 text-right text-gray-500 whitespace-nowrap">
                    {{ entry.isDir ? '' : formatBytes(entry.size) }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Disc Layout Toolbar -->
          <div class="flex items-center gap-2 px-3 py-2 bg-gray-100 dark:bg-gray-800 border-t border-gray-300 dark:border-gray-700">
            <button
              @click="removeSelectedFromProject"
              :disabled="selectedProjectEntries.length === 0"
              class="px-3 py-1 text-xs font-medium rounded bg-red-600 hover:bg-red-500 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >
              {{ t('project.remove') }}
            </button>
            <span class="flex-1"></span>
            <button
              @click="goToBurn"
              :disabled="projectStore.project.entries.length === 0"
              class="px-4 py-1 text-xs font-medium rounded bg-orange-600 hover:bg-orange-500 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >
              {{ t('burn.title') }}
            </button>
          </div>
        </div>
      </div>

      <!-- Bottom: Capacity Bar -->
      <div class="px-3 py-2 bg-gray-100 dark:bg-gray-800 border-t border-gray-300 dark:border-gray-700">
        <div class="flex items-center gap-3">
          <span class="text-xs text-gray-600 dark:text-gray-400 whitespace-nowrap">{{ t('project.capacity') }}:</span>
          <div class="flex-1 h-3 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
            <div
              class="h-full rounded-full transition-all duration-300"
              :class="capacityColor"
              :style="{ width: Math.min(100, capacityPercent) + '%' }"
            ></div>
          </div>
          <span class="text-xs text-gray-600 dark:text-gray-400 whitespace-nowrap">
            {{ formatBytes(projectStore.totalSize) }} / {{ formatBytes(deviceStore.mediaCapacityBytes) }}
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

    <!-- Right Sidebar: Device Panel (collapsible) -->
    <div
      class="border-l border-gray-300 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 transition-all duration-200"
      :class="devicePanelOpen ? 'w-64' : 'w-10'"
    >
      <!-- Toggle Button -->
      <button
        @click="devicePanelOpen = !devicePanelOpen"
        class="w-full flex items-center justify-center py-2 hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors border-b border-gray-300 dark:border-gray-700"
      >
        <svg
          class="w-4 h-4 text-gray-600 dark:text-gray-400 transition-transform"
          :class="{ 'rotate-180': !devicePanelOpen }"
          fill="none" stroke="currentColor" viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>

      <div v-show="devicePanelOpen" class="p-3 space-y-4 overflow-y-auto h-full">
        <!-- Device Selector -->
        <div>
          <label class="block text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">{{ t('device.device') }}</label>
          <select
            :value="deviceStore.currentDevicePath"
            @change="deviceStore.selectDevice(($event.target).value)"
            class="w-full px-2 py-1.5 text-sm bg-gray-200 dark:bg-gray-700 border border-gray-400 dark:border-gray-600 rounded text-gray-800 dark:text-gray-200 focus:outline-none focus:border-blue-500"
          >
            <option v-if="deviceStore.devices.length === 0" value="" disabled>{{ t('device.noDevicesFound') }}</option>
            <option v-for="dev in deviceStore.devices" :key="dev.path" :value="dev.path">
              {{ dev.name }}
            </option>
          </select>
        </div>

        <!-- Media Info -->
        <div v-if="deviceStore.mediaInfo">
          <h3 class="text-xs font-medium text-gray-600 dark:text-gray-400 mb-2">{{ t('device.mediaInfo') }}</h3>
          <div class="space-y-1 text-xs text-gray-700 dark:text-gray-300">
            <div class="flex justify-between">
              <span class="text-gray-500">{{ t('device.type') }}:</span>
              <span>{{ deviceStore.mediaInfo.mediaType }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500">{{ t('device.status') }}:</span>
              <span>{{ deviceStore.mediaInfo.mediaStatus }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500">{{ t('device.capacity') }}:</span>
              <span>{{ formatBytes(deviceStore.mediaInfo.capacityBytes) }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500">{{ t('device.free') }}:</span>
              <span>{{ formatBytes(deviceStore.mediaInfo.freeBytes) }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500">{{ t('device.sessions') }}:</span>
              <span>{{ deviceStore.mediaInfo.sessions }}</span>
            </div>
          </div>
        </div>

        <div v-else class="text-xs text-gray-500 dark:text-gray-600 text-center py-4">
          {{ t('device.noMediaDetected') }}
        </div>

        <!-- Actions -->
        <div class="space-y-2">
          <button
            @click="deviceStore.fetchMediaInfo()"
            class="w-full px-3 py-1.5 text-xs font-medium rounded bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
          >
            {{ t('device.refreshMedia') }}
          </button>
          <button
            @click="deviceStore.ejectDisc()"
            :disabled="!deviceStore.currentDevicePath"
            class="w-full px-3 py-1.5 text-xs font-medium rounded bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 disabled:opacity-40 transition-colors"
          >
            {{ t('device.ejectDisc') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
