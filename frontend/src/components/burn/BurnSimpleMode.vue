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

// --- Вычисляемые свойства состояния медиа ---

const mediaStatus = computed(() => props.mediaInfo?.mediaStatus || '')

const isErasable = computed(() => props.mediaInfo?.erasable || false)

const isFormattable = computed(() => {
  const type = (props.mediaInfo?.mediaType || '').toUpperCase()
  return type.includes('BD-RE') || type.includes('DVD-RAM') || type.includes('DVD+RW') || type.includes('DVD-RW')
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

const isRewritable = computed(() => isErasable.value || isFormattable.value)
const isNotBlank = computed(() => !mediaStatus.value.includes('blank'))

// Локальное состояние: очистить перед записью
const eraseBeforeBurn = ref(true)
</script>

<template>
  <div class="space-y-4">

    <!-- Выбор устройства (только в режиме burn) -->
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
    <div v-if="mode === 'burn' && !currentDevicePath" class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 border border-gray-200 dark:border-gray-700 text-center text-gray-500 dark:text-gray-600 py-8">
      <p class="text-sm">{{ t('device.noDeviceSelected') }}</p>
    </div>

    <template v-if="mode === 'save' || currentDevicePath">
      <!-- Карточка носителя (только в режиме burn) -->
      <div v-if="mode === 'burn'" class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
        <h3 class="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wider mb-3">
          {{ t('device.mediaInfo') }}
        </h3>

        <!-- Нет медиа -->
        <div v-if="noMedia" class="text-sm text-gray-500 dark:text-gray-600 py-2 text-center">
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

      <!-- Сводка проекта -->
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

      <!-- Быстрые опции (только в режиме burn) -->
      <div v-if="mode === 'burn'" class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
        <h3 class="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wider mb-3">
          {{ t('burn.burnOptions') }}
        </h3>
        <div class="flex flex-col gap-2">
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 cursor-pointer">
            <input
              type="checkbox"
              v-model="project.burnOptions.verify"
              class="accent-blue-500"
            />
            {{ t('burn.verifyAfterBurn') }}
            <InfoTooltip :text="t('burn.tooltips.verify')" />
          </label>
          <label class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 cursor-pointer">
            <input
              type="checkbox"
              v-model="project.burnOptions.eject"
              class="accent-blue-500"
            />
            {{ t('burn.ejectWhenDone') }}
            <InfoTooltip :text="t('burn.tooltips.eject')" />
          </label>
          <label
            v-if="isRewritable && isNotBlank"
            class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300 cursor-pointer"
          >
            <input
              type="checkbox"
              v-model="eraseBeforeBurn"
              class="accent-yellow-500"
            />
            {{ t('burn.eraseBeforeBurn') }}
            <InfoTooltip :text="t('burn.tooltips.eraseBeforeBurn')" />
          </label>
        </div>
      </div>

      <!-- Кнопки действий: режим save -->
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

      <!-- Кнопки действий: режим burn -->
      <div v-else class="flex justify-end items-center gap-3 flex-wrap">
        <!-- Предупреждение о нехватке места -->
        <span
          v-if="!hasEnoughSpace && project?.entries?.length > 0 && hasMedia"
          class="text-xs text-yellow-500 flex items-center gap-1 mr-auto"
        >
          <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v4m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z" />
          </svg>
          {{ t('burn.notEnoughSpace') }}
        </span>

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
          @click="emit('start-burn', eraseBeforeBurn && isRewritable && isNotBlank)"
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
