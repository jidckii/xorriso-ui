<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { ArrowUp, Home, Usb, Eye, EyeOff } from 'lucide-vue-next'
import PanelHeader from '../ui/PanelHeader.vue'
import SortButtons from '../ui/SortButtons.vue'

const { t } = useI18n()

const props = defineProps({
  mountPoints: { type: Array, default: () => [] },
  browsePath: { type: String, default: '/' },
  breadcrumbs: { type: Array, default: () => [] },
  editingPath: { type: Boolean, default: false },
  pathInput: { type: String, default: '' },
  showHidden: { type: Boolean, default: false },
  sortBy: { type: String, default: 'name' },
  sortDir: { type: String, default: 'asc' },
})

const emit = defineEmits([
  'navigate',
  'go-up',
  'start-edit',
  'confirm-path',
  'cancel-edit',
  'toggle-hidden',
  'toggle-sort',
  'update:pathInput',
])

// Внутренний ref для input-элемента редактирования пути
const pathInputRef = ref(null)

// Локальная модель для v-model на input
const localPathInput = computed({
  get: () => props.pathInput,
  set: (val) => emit('update:pathInput', val),
})

// Автофокус при переходе в режим редактирования
watch(() => props.editingPath, (editing) => {
  if (editing) {
    nextTick(() => {
      pathInputRef.value?.focus()
      pathInputRef.value?.select()
    })
  }
})

function onConfirmPath() {
  emit('confirm-path', localPathInput.value.trim())
}
</script>

<template>
  <PanelHeader>
    <template #row1>
      <button
        v-for="mp in mountPoints"
        :key="mp.path"
        @click="emit('navigate', mp.path)"
        class="flex items-center gap-1 px-2 py-0.5 text-xs rounded hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors shrink-0"
        :class="{
          'bg-blue-500/15 text-blue-600 dark:text-blue-400': browsePath?.startsWith(mp.path),
          'text-gray-600 dark:text-gray-400': !browsePath?.startsWith(mp.path),
        }"
        :title="mp.path"
      >
        <Home v-if="mp.icon === 'home'" :size="12" />
        <Usb v-else :size="12" />
        <span>{{ mp.label }}</span>
      </button>

      <span class="flex-1" />

      <SortButtons
        :sort-by="sortBy"
        :sort-dir="sortDir"
        @toggle-sort="emit('toggle-sort', $event)"
      />

      <div class="w-px h-4 bg-gray-300 dark:bg-gray-600 mx-0.5 shrink-0" />

      <button
        @click="emit('toggle-hidden')"
        class="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors shrink-0"
        :title="showHidden ? t('project.hideHidden') : t('project.showHidden')"
      >
        <Eye v-if="showHidden" :size="14" class="text-blue-500" />
        <EyeOff v-else :size="14" class="text-gray-500" />
      </button>
    </template>

    <template #row2>
      <button
        @click="emit('go-up')"
        class="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors shrink-0"
        :title="t('project.goUp')"
      >
        <ArrowUp :size="14" class="text-gray-600 dark:text-gray-400" />
      </button>

      <!-- Редактируемый путь -->
      <input
        v-if="editingPath"
        ref="pathInputRef"
        v-model="localPathInput"
        @keydown.enter="onConfirmPath"
        @keydown.escape="emit('cancel-edit')"
        @blur="emit('cancel-edit')"
        class="flex-1 ml-1 px-2 py-0.5 text-xs bg-white dark:bg-gray-900 border border-blue-500 rounded text-gray-800 dark:text-gray-200 outline-none"
      />

      <!-- Хлебные крошки (клик по пустой области — редактирование) -->
      <div
        v-else
        class="flex items-center gap-0.5 text-xs text-gray-600 dark:text-gray-400 overflow-hidden flex-1 ml-1 cursor-text"
        @click.self="emit('start-edit')"
      >
        <span
          class="cursor-pointer hover:text-gray-800 dark:hover:text-gray-200 transition-colors shrink-0"
          @click="emit('navigate', '/')"
        >/</span>
        <template v-for="(crumb, idx) in breadcrumbs" :key="crumb.path">
          <span
            class="cursor-pointer hover:text-gray-800 dark:hover:text-gray-200 transition-colors truncate"
            :class="{ 'text-gray-800 dark:text-gray-200 font-medium': idx === breadcrumbs.length - 1 }"
            @click="emit('navigate', crumb.path)"
          >{{ crumb.name }}</span>
          <span v-if="idx < breadcrumbs.length - 1" class="text-gray-400 dark:text-gray-600 shrink-0">/</span>
        </template>
        <!-- Невидимый спейсер для перехвата кликов по пустой области -->
        <span class="flex-1" @click="emit('start-edit')" />
      </div>
    </template>
  </PanelHeader>
</template>
