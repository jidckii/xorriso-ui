<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useDeviceStore } from '../stores/deviceStore'
import { useProjectStore } from '../stores/projectStore'
import { useBurnStore } from '../stores/burnStore'
import { useTabStore } from '../stores/tabStore'

const { t } = useI18n()
const router = useRouter()
const deviceStore = useDeviceStore()
const projectStore = useProjectStore()
const burnStore = useBurnStore()
const tabStore = useTabStore()

// Capture the active project at the time of entering burn view
const burnProject = computed(() => tabStore.activeProject)

// --- Burn Dialog State ---
const step = ref('configure') // 'configure' | 'burning' | 'done'
const showLog = ref(false)
const blankMode = ref('fast')

onMounted(() => {
  burnStore.init()
  if (!deviceStore.currentDevicePath) {
    deviceStore.init()
  }
})

// --- Computed ---
const canBurn = computed(() =>
  deviceStore.currentDevicePath &&
  burnProject.value?.entries?.length > 0 &&
  !burnStore.isBurning
)

const progressPercent = computed(() => burnStore.progress.percent)

const phaseLabel = computed(() => {
  const labels = {
    idle: t('phases.idle'),
    preparing: t('phases.preparing'),
    burning: t('phases.burning'),
    verifying: t('phases.verifying'),
    blanking: t('phases.blanking'),
    complete: t('phases.complete'),
    error: t('phases.error'),
    cancelled: t('phases.cancelled'),
  }
  return labels[burnStore.progress.phase] || burnStore.progress.phase
})

// --- Actions ---
async function startBurn() {
  step.value = 'burning'
  await burnStore.startBurn(deviceStore.currentDevicePath)

  // Watch for completion
  const checkDone = setInterval(() => {
    if (!burnStore.isBurning && burnStore.currentJob) {
      step.value = 'done'
      clearInterval(checkDone)
    }
  }, 500)
}

async function cancelBurn() {
  await burnStore.cancelBurn()
  step.value = 'done'
}

async function blankDisc() {
  await burnStore.blankDisc(deviceStore.currentDevicePath, blankMode.value)
}

function goBack() {
  burnStore.resetJob()
  router.push('/')
}

function burnAgain() {
  burnStore.resetJob()
  step.value = 'configure'
}

function formatBytes(bytes) {
  return projectStore.formatBytes(bytes)
}
</script>

<template>
  <div class="flex items-center justify-center h-full p-6">
    <div class="w-full max-w-2xl bg-gray-100 dark:bg-gray-800 rounded-lg shadow-xl border border-gray-300 dark:border-gray-700">

      <!-- Header -->
      <div class="flex items-center gap-3 px-6 py-4 border-b border-gray-300 dark:border-gray-700">
        <svg class="w-6 h-6 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M15.362 5.214A8.252 8.252 0 0112 21 8.25 8.25 0 016.038 7.048 8.287 8.287 0 009 9.6a8.983 8.983 0 013.361-6.867 8.21 8.21 0 003 2.48z" />
        </svg>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ t('burn.title') }}</h2>
      </div>

      <!-- Step: Configure -->
      <div v-if="step === 'configure'" class="p-6 space-y-6">

        <!-- Project Summary -->
        <div class="space-y-2">
          <h3 class="text-sm font-medium text-gray-600 dark:text-gray-400">{{ t('burn.project') }}</h3>
          <div class="bg-white dark:bg-gray-900 rounded px-4 py-3 text-sm space-y-1">
            <div class="flex justify-between">
              <span class="text-gray-500">{{ t('common.name') }}:</span>
              <span class="text-gray-900 dark:text-gray-100">{{ burnProject.value?.name }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500">{{ t('common.volumeId') }}:</span>
              <span class="text-gray-900 dark:text-gray-100">{{ burnProject.value?.volumeId }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500">{{ t('common.files') }}:</span>
              <span class="text-gray-900 dark:text-gray-100">{{ (burnProject.value?.entries?.length || 0) }} {{ t('project.items') }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500">{{ t('common.totalSize') }}:</span>
              <span class="text-gray-900 dark:text-gray-100">{{ formatBytes(burnProject.value?.totalSize || 0) }}</span>
            </div>
          </div>
        </div>

        <!-- Device -->
        <div class="space-y-2">
          <h3 class="text-sm font-medium text-gray-600 dark:text-gray-400">{{ t('burn.device') }}</h3>
          <select
            :value="deviceStore.currentDevicePath"
            @change="deviceStore.selectDevice(($event.target).value)"
            class="w-full px-3 py-2 text-sm bg-white dark:bg-gray-900 border border-gray-400 dark:border-gray-600 rounded text-gray-800 dark:text-gray-200 focus:outline-none focus:border-blue-500"
          >
            <option v-for="dev in deviceStore.devices" :key="dev.path" :value="dev.path">
              {{ dev.vendor }} {{ dev.model }} ({{ dev.path }})
            </option>
          </select>
          <div v-if="deviceStore.mediaInfo" class="text-xs text-gray-500">
            {{ deviceStore.mediaInfo.mediaType }} - {{ deviceStore.mediaInfo.mediaStatus }} -
            {{ formatBytes(deviceStore.mediaInfo.freeSpace) }} {{ t('device.free') }}
          </div>
        </div>

        <!-- Burn Options -->
        <div class="space-y-2">
          <h3 class="text-sm font-medium text-gray-600 dark:text-gray-400">{{ t('burn.burnOptions') }}</h3>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs text-gray-500 mb-1">{{ t('burn.speed') }}</label>
              <select
                v-model="burnProject.value.burnOptions.speed"
                class="w-full px-2 py-1.5 text-sm bg-white dark:bg-gray-900 border border-gray-400 dark:border-gray-600 rounded text-gray-800 dark:text-gray-200 focus:outline-none focus:border-blue-500"
              >
                <option v-for="s in deviceStore.speeds" :key="s.value" :value="s.value">
                  {{ s.label }}
                </option>
              </select>
            </div>
            <div>
              <label class="block text-xs text-gray-500 mb-1">{{ t('burn.burnMode') }}</label>
              <select
                v-model="burnProject.value.burnOptions.burnMode"
                class="w-full px-2 py-1.5 text-sm bg-white dark:bg-gray-900 border border-gray-400 dark:border-gray-600 rounded text-gray-800 dark:text-gray-200 focus:outline-none focus:border-blue-500"
              >
                <option value="auto">{{ t('burn.autoDao') }}</option>
                <option value="tao">{{ t('burn.tao') }}</option>
                <option value="sao">{{ t('burn.saoDao') }}</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-x-6 gap-y-2 mt-3">
            <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
              <input type="checkbox" v-model="burnProject.value.burnOptions.verify" class="accent-blue-500" />
              {{ t('burn.verifyAfterBurn') }}
            </label>
            <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
              <input type="checkbox" v-model="burnProject.value.burnOptions.eject" class="accent-blue-500" />
              {{ t('burn.ejectWhenDone') }}
            </label>
            <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
              <input type="checkbox" v-model="burnProject.value.burnOptions.dummyMode" class="accent-yellow-500" />
              {{ t('burn.simulationMode') }}
            </label>
            <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
              <input type="checkbox" v-model="burnProject.value.burnOptions.closeDisc" class="accent-blue-500" />
              {{ t('burn.closeDisc') }}
            </label>
            <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
              <input type="checkbox" v-model="burnProject.value.burnOptions.streamRecording" class="accent-blue-500" />
              {{ t('burn.streamRecording') }}
            </label>
          </div>
        </div>

        <!-- Blank disc section -->
        <div class="space-y-2 pt-2 border-t border-gray-300 dark:border-gray-700">
          <h3 class="text-sm font-medium text-gray-600 dark:text-gray-400">{{ t('burn.blankDisc') }}</h3>
          <div class="flex items-center gap-3">
            <select
              v-model="blankMode"
              class="px-2 py-1.5 text-sm bg-white dark:bg-gray-900 border border-gray-400 dark:border-gray-600 rounded text-gray-800 dark:text-gray-200 focus:outline-none focus:border-blue-500"
            >
              <option value="fast">{{ t('burn.fastBlank') }}</option>
              <option value="full">{{ t('burn.fullBlank') }}</option>
              <option value="deformat">{{ t('burn.deformat') }}</option>
            </select>
            <button
              @click="blankDisc"
              :disabled="!deviceStore.currentDevicePath || burnStore.isBurning"
              class="px-3 py-1.5 text-sm font-medium rounded bg-yellow-600 hover:bg-yellow-500 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >
              {{ t('burn.blank') }}
            </button>
          </div>
        </div>

        <!-- Action buttons -->
        <div class="flex justify-end gap-3 pt-2">
          <button
            @click="goBack"
            class="px-4 py-2 text-sm font-medium rounded bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
          >
            {{ t('burn.cancel') }}
          </button>
          <button
            @click="startBurn"
            :disabled="!canBurn"
            class="px-6 py-2 text-sm font-semibold rounded bg-orange-600 hover:bg-orange-500 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            {{ t('burn.startBurn') }}
          </button>
        </div>
      </div>

      <!-- Step: Burning -->
      <div v-else-if="step === 'burning'" class="p-6 space-y-4">
        <div class="text-center mb-4">
          <p class="text-sm text-gray-600 dark:text-gray-400">{{ phaseLabel }}</p>
          <p class="text-3xl font-bold text-orange-400 mt-1">{{ progressPercent }}%</p>
        </div>

        <!-- Progress bar -->
        <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
          <div
            class="h-full bg-orange-500 rounded-full transition-all duration-300"
            :style="{ width: progressPercent + '%' }"
          ></div>
        </div>

        <!-- Stats -->
        <div class="grid grid-cols-3 gap-4 text-center text-sm">
          <div>
            <p class="text-gray-500 text-xs">{{ t('burnProgress.speed') }}</p>
            <p class="text-gray-700 dark:text-gray-300">{{ burnStore.progress.speed || '-' }}</p>
          </div>
          <div>
            <p class="text-gray-500 text-xs">{{ t('burnProgress.written') }}</p>
            <p class="text-gray-700 dark:text-gray-300">
              {{ formatBytes(burnStore.progress.bytesWritten) }} / {{ formatBytes(burnStore.progress.bytesTotal) }}
            </p>
          </div>
          <div>
            <p class="text-gray-500 text-xs">{{ t('burnProgress.eta') }}</p>
            <p class="text-gray-700 dark:text-gray-300">{{ burnStore.progress.eta || '-' }}</p>
          </div>
        </div>

        <!-- FIFO -->
        <div class="flex items-center gap-2 text-xs text-gray-500">
          <span>{{ t('burnProgress.fifo') }}:</span>
          <div class="flex-1 h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
            <div
              class="h-full bg-green-500 rounded-full transition-all"
              :style="{ width: (burnStore.progress.fifoFill || 0) + '%' }"
            ></div>
          </div>
          <span>{{ (burnStore.progress.fifoFill || 0).toFixed(0) }}%</span>
        </div>

        <!-- Log toggle -->
        <div>
          <button
            @click="showLog = !showLog"
            class="text-xs text-gray-500 hover:text-gray-600 dark:hover:text-gray-400 transition-colors"
          >
            {{ showLog ? t('burnProgress.hideLog') : t('burnProgress.showLog') }}
          </button>
          <div v-if="showLog" class="mt-2 bg-white dark:bg-gray-900 rounded p-3 max-h-40 overflow-y-auto font-mono text-xs text-gray-600 dark:text-gray-400">
            <div v-for="(line, i) in burnStore.logLines" :key="i">{{ line }}</div>
            <div v-if="burnStore.logLines.length === 0" class="text-gray-500 dark:text-gray-600">{{ t('burnProgress.noLogOutput') }}</div>
          </div>
        </div>

        <!-- Cancel -->
        <div class="flex justify-end">
          <button
            @click="cancelBurn"
            class="px-4 py-2 text-sm font-medium rounded bg-red-600 hover:bg-red-500 transition-colors"
          >
            {{ t('burn.cancelBurn') }}
          </button>
        </div>
      </div>

      <!-- Step: Done -->
      <div v-else-if="step === 'done'" class="p-6 space-y-4">
        <div class="text-center py-4">
          <!-- Success -->
          <template v-if="burnStore.currentJob?.result?.success">
            <svg class="w-16 h-16 mx-auto text-green-400 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <h3 class="text-lg font-semibold text-green-400">{{ t('burn.burnComplete') }}</h3>
          </template>

          <!-- Failure -->
          <template v-else>
            <svg class="w-16 h-16 mx-auto text-red-400 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <h3 class="text-lg font-semibold text-red-400">{{ t('burn.burnFailed') }}</h3>
          </template>

          <p class="text-sm text-gray-600 dark:text-gray-400 mt-2">{{ burnStore.currentJob?.result?.message }}</p>
        </div>

        <!-- Log -->
        <div class="bg-white dark:bg-gray-900 rounded p-3 max-h-40 overflow-y-auto font-mono text-xs text-gray-600 dark:text-gray-400">
          <div v-for="(line, i) in burnStore.logLines" :key="i">{{ line }}</div>
        </div>

        <!-- Actions -->
        <div class="flex justify-end gap-3">
          <button
            @click="goBack"
            class="px-4 py-2 text-sm font-medium rounded bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
          >
            {{ t('burn.backToProject') }}
          </button>
          <button
            @click="burnAgain"
            class="px-4 py-2 text-sm font-medium rounded bg-orange-600 hover:bg-orange-500 transition-colors"
          >
            {{ t('burn.burnAnother') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
