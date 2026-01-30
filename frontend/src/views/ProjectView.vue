<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useDeviceStore } from '../stores/deviceStore'
import { useProjectStore } from '../stores/projectStore'
import { useTabStore } from '../stores/tabStore'

const { t } = useI18n()
const router = useRouter()
const deviceStore = useDeviceStore()
const projectStore = useProjectStore()
const tabStore = useTabStore()

// --- Per-tab project data ---
const currentProject = computed(() => tabStore.activeProject)
const tabId = computed(() => tabStore.activeTabId)

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

// --- Init ---
onMounted(async () => {
  if (currentProject.value && currentProject.value.browseEntries.length === 0) {
    currentProject.value.browseEntries = await projectStore.browseDirectory(currentProject.value.browsePath)
  }
})

// --- File Browser Actions ---
async function navigateTo(entry) {
  if (entry.isDir && currentProject.value) {
    currentProject.value.browsePath = entry.path
    currentProject.value.selectedBrowseFiles = []
    currentProject.value.browseEntries = await projectStore.browseDirectory(entry.path)
  }
}

function toggleBrowseSelection(entry) {
  if (!currentProject.value) return
  const sel = currentProject.value.selectedBrowseFiles
  const idx = sel.indexOf(entry.path)
  if (idx === -1) {
    sel.push(entry.path)
  } else {
    sel.splice(idx, 1)
  }
}

async function addSelectedToProject() {
  if (!currentProject.value) return
  const paths = currentProject.value.selectedBrowseFiles.filter(p => p !== currentProject.value.browsePath)
  if (paths.length > 0) {
    await projectStore.addFiles(tabId.value, paths)
    currentProject.value.selectedBrowseFiles = []
  }
}

// --- Disc Layout Actions ---
function toggleProjectSelection(entry) {
  if (!currentProject.value) return
  const sel = currentProject.value.selectedProjectEntries
  const idx = sel.indexOf(entry.destPath)
  if (idx === -1) {
    sel.push(entry.destPath)
  } else {
    sel.splice(idx, 1)
  }
}

async function removeSelectedFromProject() {
  if (!currentProject.value) return
  for (const destPath of currentProject.value.selectedProjectEntries) {
    await projectStore.removeEntry(tabId.value, destPath)
  }
  currentProject.value.selectedProjectEntries = []
}

// --- Burn Navigation ---
function goToBurn() {
  router.push('/burn')
}
</script>

<template>
  <div v-if="currentProject" class="flex flex-col h-full">
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
          <span class="ml-auto text-xs text-gray-500 truncate max-w-48">{{ currentProject.browsePath }}</span>
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
                v-for="entry in currentProject.browseEntries"
                :key="entry.path"
                class="border-b border-gray-100 dark:border-gray-800 hover:bg-gray-100/50 dark:hover:bg-gray-800/50 cursor-pointer"
                :class="{ 'bg-blue-900/30': currentProject.selectedBrowseFiles.includes(entry.path) }"
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
            :disabled="currentProject.selectedBrowseFiles.length === 0"
            class="px-3 py-1 text-xs font-medium rounded bg-blue-600 hover:bg-blue-500 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            {{ t('project.addToProject') }}
          </button>
          <span class="text-xs text-gray-500">
            {{ currentProject.selectedBrowseFiles.length }} {{ t('project.selected') }}
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
            {{ currentProject.entries.length }} {{ t('project.items') }} - {{ formatBytes(currentProject.totalSize) }}
          </span>
        </div>

        <div class="flex-1 overflow-y-auto">
          <div v-if="currentProject.entries.length === 0" class="flex items-center justify-center h-full text-gray-500 dark:text-gray-600">
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
                v-for="entry in currentProject.entries"
                :key="entry.destPath"
                class="border-b border-gray-100 dark:border-gray-800 hover:bg-gray-100/50 dark:hover:bg-gray-800/50 cursor-pointer"
                :class="{ 'bg-blue-900/30': currentProject.selectedProjectEntries.includes(entry.destPath) }"
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
            :disabled="currentProject.selectedProjectEntries.length === 0"
            class="px-3 py-1 text-xs font-medium rounded bg-red-600 hover:bg-red-500 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            {{ t('project.remove') }}
          </button>
          <span class="flex-1"></span>
          <button
            @click="goToBurn"
            :disabled="currentProject.entries.length === 0"
            class="px-4 py-1 text-xs font-medium rounded bg-orange-600 hover:bg-orange-500 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            {{ t('burn.title') }}
          </button>
        </div>
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
