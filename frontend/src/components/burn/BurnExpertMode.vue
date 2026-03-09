<script setup>
import { ref, computed, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { Pencil, Check, Disc, Flame, Save } from 'lucide-vue-next'
import { formatBytes } from '../../composables/useFormatBytes'
import InfoTooltip from '../ui/InfoTooltip.vue'

const props = defineProps({
  mode: { type: String, default: 'burn' },
  project: { type: Object, required: true },
  devices: { type: Array, required: true },
  currentDevicePath: { type: String, default: '' },
  mediaInfo: { type: Object, default: null },
  speeds: { type: Array, default: () => [] },
  isBurning: { type: Boolean, default: false },
  mediaCapacityBytes: { type: Number, default: 0 },
})

const emit = defineEmits([
  'select-device',
  'start-burn',
  'create-iso',
  'blank-disc',
  'format-disc',
  'refresh-media',
  'eject',
  'copy-command',
  'save-project',
  'update-name',
])

const { t } = useI18n()

// --- Редактирование имени ---
const editingName = ref(false)
const nameInput = ref(null)

async function startEditName() {
  editingName.value = true
  await nextTick()
  nameInput.value?.focus()
  nameInput.value?.select()
}

function confirmName() {
  editingName.value = false
}

// Режим очистки (для секции Blank/Format)
const blankMode = ref('fast')

// --- Вычисляемые свойства состояния медиа ---

const mediaStatus = computed(() => props.mediaInfo?.mediaStatus || '')

const isErasable = computed(() => props.mediaInfo?.erasable || false)

const isFormattable = computed(() => {
  const type = (props.mediaInfo?.mediaType || '').toUpperCase()
  return (
    type.includes('BD-RE') ||
    type.includes('DVD-RAM') ||
    type.includes('DVD+RW') ||
    type.includes('DVD-RW')
  )
})

const hasMedia = computed(() => !!props.mediaInfo && !!props.mediaInfo.mediaType)

const hasEnoughSpace = computed(() => {
  if (!props.project || !props.mediaCapacityBytes) return false
  return props.project.totalSize <= props.mediaCapacityBytes
})

const capacityPercent = computed(() => {
  if (!props.mediaCapacityBytes || !props.project) return 0
  return (props.project.totalSize / props.mediaCapacityBytes) * 100
})

const capacityColor = computed(() => {
  const p = capacityPercent.value
  if (p > 100) return 'bg-red-500'
  if (p > 85) return 'bg-yellow-500'
  return 'bg-green-500'
})

// Можно ли начать запись
const canBurn = computed(() => {
  return (
    props.currentDevicePath &&
    props.project?.entries?.length > 0 &&
    !props.isBurning &&
    hasMedia.value &&
    (mediaStatus.value.includes('blank') || mediaStatus.value.includes('appendable')) &&
    hasEnoughSpace.value
  )
})
</script>

<template>
  <div class="space-y-4">

    <!-- Секция 1: Выбор устройства (только в режиме burn) -->
    <div v-if="mode === 'burn'" class="flex items-center gap-3 flex-wrap">
      <label class="text-sm font-medium text-gray-600 dark:text-gray-400 shrink-0">
        {{ t('device.device') }}:
      </label>
      <div class="relative">
        <select
          :value="currentDevicePath"
          @change="emit('select-device', $event.target.value)"
          class="appearance-none bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 text-sm rounded px-3 py-1.5 pr-8 border border-gray-400 dark:border-gray-600 hover:border-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent cursor-pointer min-w-[200px]"
        >
          <option value="" disabled>{{ t('device.selectDevice') }}</option>
          <option
            v-for="dev in devices"
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
        @click="emit('refresh-media')"
        class="px-3 py-1.5 text-xs font-medium rounded bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
      >
        {{ t('device.refreshMedia') }}
      </button>
      <button
        @click="emit('eject')"
        :disabled="!currentDevicePath"
        class="px-3 py-1.5 text-xs font-medium rounded bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 disabled:opacity-40 transition-colors"
      >
        {{ t('device.ejectDisc') }}
      </button>
    </div>

    <!-- Нет устройства (только в режиме burn) -->
    <div
      v-if="mode === 'burn' && !currentDevicePath"
      class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 border border-gray-200 dark:border-gray-700 text-center text-gray-500 dark:text-gray-600 py-8"
    >
      <p class="text-sm">{{ t('device.noDeviceSelected') }}</p>
    </div>

    <template v-if="mode === 'save' || currentDevicePath">
      <!-- Карточка носителя (только в режиме burn) -->
      <div v-if="mode === 'burn'" class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
        <h3 class="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wider mb-3">
          {{ t('device.mediaInfo') }}
        </h3>

        <!-- Нет медиа -->
        <div v-if="!hasMedia" class="text-sm text-gray-500 dark:text-gray-600 py-2 text-center">
          {{ t('device.noMediaInserted') }}
        </div>

        <div v-else class="grid grid-cols-[auto_1fr] gap-x-4 gap-y-1.5 text-sm text-gray-700 dark:text-gray-300">
          <span class="text-gray-500">{{ t('device.type') }}:</span>
          <span>{{ mediaInfo.mediaType || '—' }}</span>
          <span class="text-gray-500">{{ t('device.status') }}:</span>
          <span>{{ mediaInfo.mediaStatus || '—' }}</span>
          <span class="text-gray-500">{{ t('device.capacity') }}:</span>
          <span>{{ formatBytes(mediaInfo.totalCapacity) }}</span>
          <span class="text-gray-500">{{ t('device.erasable') }}:</span>
          <span>{{ mediaInfo.erasable ? t('device.yes') : t('device.no') }}</span>
          <span class="text-gray-500">{{ t('device.sessions') }}:</span>
          <span>{{ mediaInfo.sessions }}</span>
        </div>
      </div>

      <!-- Секция 2: Сводка проекта + индикатор ёмкости -->
      <div class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
        <h3 class="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wider mb-3">
          {{ t('burn.project') }}
        </h3>
        <div class="space-y-2 text-sm text-gray-700 dark:text-gray-300">
          <div class="flex justify-between items-center">
            <span class="text-gray-500">{{ t('common.name') }}:</span>
            <div class="flex items-center gap-1.5">
              <template v-if="editingName">
                <input
                  ref="nameInput"
                  :value="project?.name || ''"
                  @input="emit('update-name', $event.target.value)"
                  @keydown.enter="confirmName"
                  class="bg-gray-100 dark:bg-gray-700 text-right rounded px-2 py-0.5 text-sm text-gray-700 dark:text-gray-300 border border-gray-300 dark:border-gray-600 focus:border-blue-500 focus:outline-none max-w-[200px]"
                />
                <button
                  @click="confirmName"
                  class="p-0.5 rounded hover:bg-gray-200 dark:hover:bg-gray-700 text-green-500 transition-colors"
                >
                  <Check :size="14" />
                </button>
              </template>
              <template v-else>
                <span>{{ project?.name || '—' }}</span>
                <button
                  @click="startEditName"
                  class="p-0.5 rounded hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
                >
                  <Pencil :size="12" />
                </button>
              </template>
            </div>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-500">{{ t('common.files') }}:</span>
            <span>{{ project?.entries?.length || 0 }} {{ t('project.items') }}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-500">{{ t('common.totalSize') }}:</span>
            <span>{{ formatBytes(project?.totalSize || 0) }}</span>
          </div>
        </div>

        <!-- Индикатор заполнения диска (только в режиме burn) -->
        <div v-if="mode === 'burn' && mediaCapacityBytes > 0" class="mt-3 space-y-1">
          <div class="h-3 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
            <div
              class="h-full rounded-full transition-all"
              :class="capacityColor"
              :style="{ width: Math.min(100, capacityPercent) + '%' }"
            ></div>
          </div>
          <span class="text-xs text-gray-500">
            {{ capacityPercent.toFixed(1) }}% —
            {{ hasEnoughSpace ? t('burn.dataFits') : t('burn.notEnoughSpace') }}
          </span>
        </div>
      </div>

      <!-- Секция 3: Опции ISO / файловая система -->
      <div class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
        <h3 class="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wider mb-3">
          {{ t('burn.isoOptions') }}
        </h3>

        <!-- ISO Level — отдельный ряд -->
        <div class="mb-3">
          <label class="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400 mb-1">
            {{ t('burn.isoLevel') }}
            <InfoTooltip :text="t('burn.tooltips.isoLevel')" :link="t('burn.tooltipLinks.isoLevel')" />
          </label>
          <div class="relative inline-block">
            <select
              v-model="project.isoOptions.isoLevel"
              class="appearance-none w-full bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 text-sm rounded px-3 py-1.5 pr-8 border border-gray-400 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent cursor-pointer"
            >
              <option value="1">1</option>
              <option value="2">2</option>
              <option value="3">3</option>
              <option value="4">4</option>
            </select>
            <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
              <svg class="w-4 h-4 text-gray-600 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </div>
          </div>
        </div>

        <!-- Чекбоксы в сетке 2 колонки -->
        <div class="grid grid-cols-2 gap-x-6 gap-y-2">
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 cursor-pointer">
            <input type="checkbox" v-model="project.isoOptions.udf" class="accent-blue-500" />
            {{ t('burn.udf') }}
            <InfoTooltip :text="t('burn.tooltips.udf')" :link="t('burn.tooltipLinks.udf')" />
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 cursor-pointer">
            <input type="checkbox" v-model="project.isoOptions.rockRidge" class="accent-blue-500" />
            {{ t('burn.rockRidge') }}
            <InfoTooltip :text="t('burn.tooltips.rockRidge')" :link="t('burn.tooltipLinks.rockRidge')" />
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 cursor-pointer">
            <input type="checkbox" v-model="project.isoOptions.joliet" class="accent-blue-500" />
            {{ t('burn.joliet') }}
            <InfoTooltip :text="t('burn.tooltips.joliet')" :link="t('burn.tooltipLinks.joliet')" />
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 cursor-pointer">
            <input type="checkbox" v-model="project.isoOptions.hfsPlus" class="accent-blue-500" />
            {{ t('burn.hfsPlus') }}
            <InfoTooltip :text="t('burn.tooltips.hfsPlus')" :link="t('burn.tooltipLinks.hfsPlus')" />
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 cursor-pointer">
            <input type="checkbox" v-model="project.isoOptions.zisofs" class="accent-blue-500" />
            {{ t('burn.zisofs') }}
            <InfoTooltip :text="t('burn.tooltips.zisofs')" :link="t('burn.tooltipLinks.zisofs')" />
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 cursor-pointer">
            <input type="checkbox" v-model="project.isoOptions.md5" class="accent-blue-500" />
            {{ t('burn.md5') }}
            <InfoTooltip :text="t('burn.tooltips.md5')" />
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 cursor-pointer">
            <input type="checkbox" v-model="project.isoOptions.backupMode" class="accent-blue-500" />
            {{ t('burn.backupMode') }}
            <InfoTooltip :text="t('burn.tooltips.backupMode')" />
          </label>
        </div>
      </div>

      <!-- Секция 4: Параметры записи -->
      <div class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
        <h3 class="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wider mb-3">
          {{ t('burn.burnOptions') }}
        </h3>

        <!-- Скорость + Режим записи -->
        <div class="grid grid-cols-2 gap-4 mb-3">
          <div>
            <label class="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400 mb-1">
              {{ t('burn.speed') }}
              <InfoTooltip :text="t('burn.tooltips.speed')" />
            </label>
            <div class="relative">
              <select
                v-model="project.burnOptions.speed"
                class="appearance-none w-full bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 text-sm rounded px-3 py-1.5 pr-8 border border-gray-400 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent cursor-pointer"
              >
                <option value="auto">{{ t('common.auto') }}</option>
                <option v-for="s in speeds" :key="s.writeSpeed" :value="s.writeSpeed">
                  {{ s.displayName }}
                </option>
              </select>
              <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                <svg class="w-4 h-4 text-gray-600 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </div>
            </div>
          </div>
          <div>
            <label class="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400 mb-1">
              {{ t('burn.burnMode') }}
              <InfoTooltip :text="t('burn.tooltips.burnMode')" :link="t('burn.tooltipLinks.burnMode')" />
            </label>
            <div class="relative">
              <select
                v-model="project.burnOptions.burnMode"
                class="appearance-none w-full bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 text-sm rounded px-3 py-1.5 pr-8 border border-gray-400 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent cursor-pointer"
              >
                <option value="auto">{{ t('burn.autoDao') }}</option>
                <option value="tao">{{ t('burn.tao') }}</option>
                <option value="sao">{{ t('burn.saoDao') }}</option>
              </select>
              <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                <svg class="w-4 h-4 text-gray-600 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </div>
            </div>
          </div>
        </div>

        <!-- Чекбоксы в grid 2 колонки -->
        <div class="grid grid-cols-2 gap-x-6 gap-y-2 mb-3">
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 cursor-pointer">
            <input type="checkbox" v-model="project.burnOptions.verify" class="accent-blue-500" />
            {{ t('burn.verifyAfterBurn') }}
            <InfoTooltip :text="t('burn.tooltips.verify')" />
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 cursor-pointer">
            <input type="checkbox" v-model="project.burnOptions.eject" class="accent-blue-500" />
            {{ t('burn.ejectWhenDone') }}
            <InfoTooltip :text="t('burn.tooltips.eject')" />
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 cursor-pointer">
            <input type="checkbox" v-model="project.burnOptions.dummyMode" class="accent-yellow-500" />
            {{ t('burn.simulationMode') }}
            <InfoTooltip :text="t('burn.tooltips.dummyMode')" />
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 cursor-pointer">
            <input type="checkbox" v-model="project.burnOptions.closeDisc" class="accent-blue-500" />
            {{ t('burn.closeDisc') }}
            <InfoTooltip :text="t('burn.tooltips.closeDisc')" />
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 cursor-pointer">
            <input type="checkbox" v-model="project.burnOptions.streamRecording" class="accent-blue-500" />
            {{ t('burn.streamRecording') }}
            <InfoTooltip :text="t('burn.tooltips.streamRecording')" />
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 cursor-pointer">
            <input type="checkbox" v-model="project.burnOptions.multisession" class="accent-blue-500" />
            {{ t('burn.multisession') }}
            <InfoTooltip :text="t('burn.tooltips.multisession')" />
          </label>
        </div>

        <!-- Подсказка мультисессии -->
        <p v-if="project.burnOptions.multisession" class="text-xs text-yellow-500 mb-3">
          {{ t('burn.multisessionHint') }}
        </p>

        <!-- Padding — отдельно -->
        <div class="flex items-center gap-3">
          <label class="flex items-center gap-1 text-xs text-gray-500 dark:text-gray-400 shrink-0">
            {{ t('burn.padding') }}:
            <InfoTooltip :text="t('burn.tooltips.padding')" />
          </label>
          <input
            type="number"
            v-model.number="project.burnOptions.padding"
            min="0"
            class="w-24 bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 text-sm rounded px-2 py-1 border border-gray-400 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
          <span class="text-xs text-gray-500">KiB</span>
        </div>
      </div>

      <!-- Секция 5: Blank / Format (только в режиме burn) -->
      <div v-if="mode === 'burn'" class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
        <h3 class="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wider mb-3">
          {{ t('burn.blankDisc') }}
        </h3>
        <div class="flex items-center gap-3 flex-wrap">
          <div class="relative">
            <select
              v-model="blankMode"
              class="appearance-none w-full bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 text-sm rounded px-3 py-1.5 pr-8 border border-gray-400 dark:border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent cursor-pointer"
            >
              <option value="fast">{{ t('burn.fastBlank') }}</option>
              <option value="full">{{ t('burn.fullBlank') }}</option>
              <option value="deformat">{{ t('burn.deformat') }}</option>
            </select>
            <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
              <svg class="w-4 h-4 text-gray-600 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </div>
          </div>
          <button
            @click="emit('blank-disc', blankMode)"
            :disabled="!currentDevicePath || isBurning"
            class="px-4 py-1.5 text-sm font-medium rounded bg-yellow-600 hover:bg-yellow-500 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            {{ t('burn.blank') }}
          </button>
          <button
            @click="emit('format-disc', blankMode)"
            :disabled="!currentDevicePath || isBurning"
            class="px-4 py-1.5 text-sm font-medium rounded bg-yellow-600 hover:bg-yellow-500 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            {{ t('burn.formatDisc') }}
          </button>
        </div>
      </div>

      <!-- Секция 6: Кнопки действий — режим save -->
      <div v-if="mode === 'save'" class="flex justify-end gap-3">
        <button
          @click="emit('save-project')"
          :disabled="!project?.entries?.length"
          class="inline-flex items-center gap-1.5 px-6 py-2 text-sm font-semibold rounded bg-blue-600 hover:bg-blue-500 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
        >
          <Save :size="16" />
          {{ t('header.save') }}
        </button>
      </div>

      <!-- Секция 6: Кнопки действий — режим burn -->
      <div v-else class="flex justify-end gap-3 flex-wrap">
        <!-- Скопировать команду -->
        <button
          @click="emit('copy-command')"
          class="px-4 py-2 text-sm font-medium rounded bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
        >
          {{ t('burn.copyCommand') }}
        </button>

        <!-- Создать ISO -->
        <button
          @click="emit('create-iso')"
          :disabled="!project?.entries?.length || isBurning"
          class="inline-flex items-center gap-1.5 px-5 py-2 text-sm font-semibold rounded bg-blue-600 hover:bg-blue-500 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
        >
          <Disc :size="16" />
          {{ t('burn.createIso') }}
        </button>

        <!-- Прожиг -->
        <button
          @click="emit('start-burn')"
          :disabled="!canBurn"
          class="inline-flex items-center gap-1.5 px-6 py-2 text-sm font-semibold rounded bg-orange-600 hover:bg-orange-500 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
        >
          <Flame :size="16" />
          {{ t('burn.startBurn') }}
        </button>
      </div>
    </template>
  </div>
</template>
