<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Dialogs } from '@wailsio/runtime'
import { useDeviceStore } from '../stores/deviceStore'
import { useProjectStore } from '../stores/projectStore'
import { useBurnStore } from '../stores/burnStore'
import { useTabStore } from '../stores/tabStore'
import BurnConfigure from '../components/burn/BurnConfigure.vue'
import BurnRunning from '../components/burn/BurnRunning.vue'
import BurnResult from '../components/burn/BurnResult.vue'

const { t } = useI18n()
const router = useRouter()
const deviceStore = useDeviceStore()
const projectStore = useProjectStore()
const burnStore = useBurnStore()
const tabStore = useTabStore()

// Захватываем активный проект при входе на экран записи
const burnProject = computed(() => tabStore.activeProject)

// --- Состояние шагов ---
const step = ref('configure') // 'configure' | 'burning' | 'done'

onMounted(() => {
  burnStore.init()
  if (!deviceStore.currentDevicePath) {
    deviceStore.init()
  }
})

// --- Computed ---
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

// --- Действия ---
async function startBurn() {
  step.value = 'burning'
  await burnStore.startBurn(deviceStore.currentDevicePath)

  // Отслеживаем завершение
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

async function blankDisc(mode) {
  await burnStore.blankDisc(deviceStore.currentDevicePath, mode)
}

function goBack() {
  burnStore.resetJob()
  router.push('/')
}

function burnAgain() {
  burnStore.resetJob()
  step.value = 'configure'
}

async function createISO() {
  try {
    const outputPath = await Dialogs.SaveFile({
      title: t('burn.createIsoTitle'),
      filters: [{ displayName: 'ISO Image', pattern: '*.iso' }],
    })
    if (!outputPath) return

    step.value = 'burning'
    await burnStore.createISO(burnProject.value, outputPath)

    // Ожидаем завершения
    const checkDone = setInterval(() => {
      if (!burnStore.isBurning && burnStore.currentJob) {
        step.value = 'done'
        clearInterval(checkDone)
      }
    }, 500)
  } catch (error) {
    console.error('Failed to create ISO:', error)
  }
}
</script>

<template>
  <div class="flex items-center justify-center h-full p-6">
    <div class="w-full max-w-2xl bg-gray-100 dark:bg-gray-800 rounded-lg shadow-xl border border-gray-300 dark:border-gray-700">

      <!-- Заголовок -->
      <div class="flex items-center gap-3 px-6 py-4 border-b border-gray-300 dark:border-gray-700">
        <svg class="w-6 h-6 text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M15.362 5.214A8.252 8.252 0 0112 21 8.25 8.25 0 016.038 7.048 8.287 8.287 0 009 9.6a8.983 8.983 0 013.361-6.867 8.21 8.21 0 003 2.48z" />
        </svg>
        <h2 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ t('burn.title') }}</h2>
      </div>

      <!-- Шаг: Настройка -->
      <BurnConfigure
        v-if="step === 'configure'"
        :project="burnProject"
        :devices="deviceStore.devices"
        :current-device-path="deviceStore.currentDevicePath"
        :media-info="deviceStore.mediaInfo"
        :speeds="deviceStore.speeds"
        :is-burning="burnStore.isBurning"
        @select-device="deviceStore.selectDevice($event)"
        @start-burn="startBurn"
        @create-iso="createISO"
        @blank-disc="blankDisc"
        @go-back="goBack"
      />

      <!-- Шаг: Запись -->
      <BurnRunning
        v-else-if="step === 'burning'"
        :progress="burnStore.progress"
        :log-lines="burnStore.logLines"
        :phase-label="phaseLabel"
        @cancel="cancelBurn"
      />

      <!-- Шаг: Результат -->
      <BurnResult
        v-else-if="step === 'done'"
        :job="burnStore.currentJob"
        :log-lines="burnStore.logLines"
        @go-back="goBack"
        @burn-again="burnAgain"
      />
    </div>
  </div>
</template>
