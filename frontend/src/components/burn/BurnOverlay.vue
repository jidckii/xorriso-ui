<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Dialogs } from '@wailsio/runtime'
import { useTabStore } from '../../stores/tabStore'
import { useDeviceStore } from '../../stores/deviceStore'
import { useBurnStore } from '../../stores/burnStore'
import { useProjectStore } from '../../stores/projectStore'
import BurnSimpleMode from './BurnSimpleMode.vue'
import BurnExpertMode from './BurnExpertMode.vue'
import BurnRunning from './BurnRunning.vue'
import BurnResult from './BurnResult.vue'

const { t } = useI18n()

const props = defineProps({
  mode: { type: String, default: 'burn' },
})

const tabStore = useTabStore()
const deviceStore = useDeviceStore()
const burnStore = useBurnStore()
const projectStore = useProjectStore()

const emit = defineEmits(['close'])

// --- Машина состояний ---
const step = ref('configure') // 'configure' | 'burning' | 'done'

// --- Режим отображения (сохраняется в localStorage) ---
const viewMode = ref(localStorage.getItem('xorriso-burn-mode') || 'simple')

function setViewMode(mode) {
  viewMode.value = mode
  localStorage.setItem('xorriso-burn-mode', mode)
}

// --- Стили для toggle-кнопок ---
const activeTabClass = 'px-3 py-1 text-xs font-medium rounded bg-white dark:bg-gray-600 text-gray-900 dark:text-gray-100 shadow-sm transition-colors'
const inactiveTabClass = 'px-3 py-1 text-xs font-medium rounded text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200 transition-colors'

// --- Фаза операции (для BurnRunning) ---
const phaseLabel = computed(() => {
  const labels = {
    idle: t('phases.idle'),
    preparing: t('phases.preparing'),
    burning: t('phases.burning'),
    verifying: t('phases.verifying'),
    blanking: t('phases.blanking'),
    creating_iso: t('phases.creating_iso'),
    complete: t('phases.complete'),
    error: t('phases.error'),
    cancelled: t('phases.cancelled'),
  }
  return labels[burnStore.progress.phase] || burnStore.progress.phase
})

// --- Жизненный цикл ---

function onKeydown(e) {
  if (e.key === 'Escape') {
    handleClose()
  }
}

onMounted(() => {
  document.addEventListener('keydown', onKeydown)
  burnStore.init()
  if (!deviceStore.speeds?.length && deviceStore.currentDevicePath) {
    deviceStore.fetchSpeeds()
  }
})

onUnmounted(() => {
  document.removeEventListener('keydown', onKeydown)
})

// --- Действия записи ---

async function startBurn() {
  const project = tabStore.activeProject
  step.value = 'burning'
  await burnStore.startBurn(project, deviceStore.currentDevicePath, project.burnOptions)
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

function burnAgain() {
  burnStore.resetJob()
  step.value = 'configure'
}

async function createISO() {
  const project = tabStore.activeProject
  const outputPath = await Dialogs.SaveFile({
    Title: t('burn.createIsoTitle'),
    Filename: (project?.name || 'image') + '.iso',
    Filters: [{ DisplayName: 'ISO Image', Pattern: '*.iso' }],
  })
  if (!outputPath) return
  step.value = 'burning'
  await burnStore.createISO(tabStore.activeProject, outputPath)
  const checkDone = setInterval(() => {
    if (!burnStore.isBurning && burnStore.currentJob) {
      step.value = 'done'
      clearInterval(checkDone)
    }
  }, 500)
}

async function handleBlank(mode) {
  await burnStore.blankDisc(deviceStore.currentDevicePath, mode)
}

async function handleFormat(mode) {
  await burnStore.formatDisc(deviceStore.currentDevicePath, mode)
}

async function handleCopyCommand() {
  const project = tabStore.activeProject
  const command = await burnStore.getBurnCommand(project, deviceStore.currentDevicePath, project.burnOptions)
  if (command) {
    await navigator.clipboard.writeText(command)
  }
}

// --- Сохранение проекта (режим 'save') ---

async function saveProject() {
  const activeTabId = tabStore.activeTabId
  const data = tabStore.getProjectData(activeTabId)
  if (!data) return

  if (data.filePath) {
    await projectStore.saveProject(activeTabId)
  } else {
    const filePath = await Dialogs.SaveFile({
      Title: t('header.save'),
      Filename: (data.name || 'project') + '.xorriso-project',
      Filters: [{ DisplayName: 'Xorriso Project', Pattern: '*.xorriso-project' }],
    })
    if (!filePath) return
    await projectStore.saveProjectAs(activeTabId, filePath)
  }
  emit('close')
}

// --- Переименование проекта ---

function updateProjectName(name) {
  const activeTabId = tabStore.activeTabId
  const data = tabStore.getProjectData(activeTabId)
  if (data) {
    data.name = name
    data.volumeId = name
    tabStore.updateTabLabel(activeTabId, name)
  }
}

// --- Закрытие overlay ---

function handleClose() {
  if (step.value === 'burning') {
    const confirmed = window.confirm(t('burn.cancelConfirm'))
    if (!confirmed) return
  }
  emit('close')
}
</script>

<template>
  <div class="p-6">
    <div class="space-y-6">

      <!-- Заголовок -->
      <div class="flex items-center justify-between gap-4">
        <h2 class="text-sm font-semibold text-gray-900 dark:text-gray-100 shrink-0">
          {{ mode === 'save' ? t('header.save') : t('burn.title') }}
        </h2>

        <!-- Toggle простой / экспертный режим (только в шаге configure) -->
        <div v-if="step === 'configure'" class="flex bg-gray-200 dark:bg-gray-700 rounded p-0.5">
          <button
            @click="setViewMode('simple')"
            :class="viewMode === 'simple' ? activeTabClass : inactiveTabClass"
          >
            {{ t('burn.simpleMode') }}
          </button>
          <button
            @click="setViewMode('expert')"
            :class="viewMode === 'expert' ? activeTabClass : inactiveTabClass"
          >
            {{ t('burn.expertMode') }}
          </button>
        </div>

        <!-- Кнопка закрытия -->
        <button
          @click="handleClose"
          class="p-1.5 rounded hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors ml-auto"
        >
          <svg class="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Шаг: Настройка -->
      <template v-if="step === 'configure'">
        <!-- Простой режим -->
        <BurnSimpleMode
          v-if="viewMode === 'simple'"
          :mode="mode"
          :project="tabStore.activeProject"
          :devices="deviceStore.devices"
          :current-device-path="deviceStore.currentDevicePath || ''"
          :media-info="deviceStore.mediaInfo"
          :speeds="deviceStore.speeds"
          :is-burning="burnStore.isBurning"
          :media-capacity-bytes="deviceStore.mediaCapacityBytes"
          @select-device="deviceStore.selectDevice($event)"
          @start-burn="startBurn"
          @blank-disc="handleBlank"
          @format-disc="handleFormat"
          @refresh-media="deviceStore.fetchMediaInfo()"
          @eject="deviceStore.ejectDisc()"
          @create-iso="createISO"
          @save-project="saveProject"
          @update-name="updateProjectName"
        />

        <!-- Экспертный режим -->
        <BurnExpertMode
          v-else
          :mode="mode"
          :project="tabStore.activeProject"
          :devices="deviceStore.devices"
          :current-device-path="deviceStore.currentDevicePath || ''"
          :media-info="deviceStore.mediaInfo"
          :speeds="deviceStore.speeds"
          :is-burning="burnStore.isBurning"
          :media-capacity-bytes="deviceStore.mediaCapacityBytes"
          @select-device="deviceStore.selectDevice($event)"
          @start-burn="startBurn"
          @create-iso="createISO"
          @blank-disc="handleBlank"
          @format-disc="handleFormat"
          @refresh-media="deviceStore.fetchMediaInfo()"
          @eject="deviceStore.ejectDisc()"
          @copy-command="handleCopyCommand"
          @save-project="saveProject"
          @update-name="updateProjectName"
        />
      </template>

      <!-- Шаг: Запись (только в режиме burn) -->
      <BurnRunning
        v-else-if="mode === 'burn' && step === 'burning'"
        :progress="burnStore.progress"
        :log-lines="burnStore.logLines"
        :phase-label="phaseLabel"
        @cancel="cancelBurn"
      />

      <!-- Шаг: Результат (только в режиме burn) -->
      <BurnResult
        v-else-if="mode === 'burn' && step === 'done'"
        :job="burnStore.currentJob"
        :log-lines="burnStore.logLines"
        @go-back="handleClose"
        @burn-again="burnAgain"
      />

    </div>
  </div>
</template>
