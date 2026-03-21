<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ChevronRight } from 'lucide-vue-next'
import FileIcon from '../ui/FileIcon.vue'
import ImagePreviewTooltip from '../ui/ImagePreviewTooltip.vue'
import { formatBytes } from '../../composables/useFormatBytes'
import { isImageFile } from '../../composables/useImagePreview'

const props = defineProps({
  entry: { type: Object, required: true },
  depth: { type: Number, default: 0 },
  expandedDirs: { type: Set, required: true },
  dirChildren: { type: Object, required: true },
  selectedPaths: { type: Set, required: true },
  showHidden: { type: Boolean, default: false },
  sortFn: { type: Function, default: null },
  focusedPath: { type: String, default: null },
})

const emit = defineEmits(['toggle-expand', 'toggle-selection', 'dblclick', 'contextmenu', 'dragstart'])

const { t } = useI18n()

function isExpanded(entry) {
  return props.expandedDirs.has(entry.sourcePath)
}

function isSelected(entry) {
  return props.selectedPaths.has(entry.sourcePath)
}

function getChildren(entry) {
  const all = props.dirChildren[entry.sourcePath] || []
  const filtered = props.showHidden ? all : all.filter(c => !c.name.startsWith('.'))
  if (props.sortFn) return props.sortFn(filtered)
  return filtered
}

function onDragStart(event) {
  emit('dragstart', props.entry, event)
}

function onMouseDown(event) {
  // Shift+Click на draggable не генерирует click — обрабатываем на mousedown
  if (event.shiftKey) {
    event.preventDefault()
    emit('toggle-selection', props.entry, event)
  }
}

// Image preview tooltip state
const previewVisible = ref(false)
const previewX = ref(0)
const previewY = ref(0)

function onMouseEnter(e, entry) {
  if (entry.isDir || !isImageFile(entry.name)) return
  previewX.value = e.clientX
  previewY.value = e.clientY
  previewVisible.value = true
}

function onMouseMove(e) {
  previewX.value = e.clientX
  previewY.value = e.clientY
}

function onMouseLeave() {
  previewVisible.value = false
}

function getExtension(name) {
  const dot = name.lastIndexOf('.')
  return dot > 0 ? name.slice(dot + 1).toLowerCase() : ''
}

function getTypeLabel(entry) {
  if (entry.isDir) return t('common.folder')
  return getExtension(entry.name)
}

// Форматирует дату модификации в локальный формат дд.мм.гг чч:мм
function formatDate(modTime) {
  if (!modTime) return ''
  const date = new Date(modTime)
  return (
    date.toLocaleDateString(undefined, { year: 'numeric', month: '2-digit', day: '2-digit' }) +
    ' ' +
    date.toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit', hour12: false })
  )
}
</script>

<template>
  <!-- Entry row -->
  <div
    draggable="true"
    :data-path="entry.sourcePath"
    class="flex items-center gap-1.5 py-1 cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-800/50 transition-colors"
    :class="{
      'bg-blue-500/15': isSelected(entry),
      'ring-1 ring-blue-400/60 ring-inset': focusedPath === entry.sourcePath && !isSelected(entry),
      'ring-1 ring-blue-400 ring-inset': focusedPath === entry.sourcePath && isSelected(entry),
    }"
    :style="{ paddingLeft: (depth * 16 + 8) + 'px', paddingRight: '8px' }"
    @mousedown="onMouseDown"
    @click="!$event.shiftKey && emit('toggle-selection', entry, $event)"
    @dblclick="emit('dblclick', entry)"
    @contextmenu.prevent="emit('contextmenu', entry, $event)"
    @dragstart="onDragStart"
    @mouseenter="onMouseEnter($event, entry)"
    @mousemove="onMouseMove"
    @mouseleave="onMouseLeave"
  >
    <!-- Expand chevron for directories -->
    <button
      v-if="entry.isDir"
      class="w-4 h-4 flex items-center justify-center shrink-0 hover:bg-gray-300 dark:hover:bg-gray-600 rounded"
      @click.stop="emit('toggle-expand', entry)"
    >
      <ChevronRight
        :size="14"
        class="text-gray-500 transition-transform duration-150"
        :class="{ 'rotate-90': isExpanded(entry) }"
      />
    </button>
    <span v-else class="w-4 h-4 shrink-0" />

    <input
      type="checkbox"
      :checked="isSelected(entry)"
      @click.stop="emit('toggle-selection', entry, $event)"
      class="w-3.5 h-3.5 accent-blue-600 shrink-0 cursor-pointer"
    />

    <!-- Icon -->
    <FileIcon
      :name="entry.name"
      :is-dir="entry.isDir"
      :is-open="isExpanded(entry)"
      :size="16"
    />

    <!-- Name -->
    <span class="truncate flex-1 text-gray-800 dark:text-gray-200">
      {{ entry.name }}
    </span>

    <!-- Date -->
    <span class="text-xs text-gray-500 shrink-0 w-28 text-right">
      {{ formatDate(entry.modTime) }}
    </span>

    <!-- Type -->
    <span class="text-xs text-gray-500 shrink-0 w-16 text-right">
      {{ getTypeLabel(entry) }}
    </span>

    <!-- Size -->
    <span class="text-xs text-gray-500 shrink-0 w-20 text-right">
      {{ entry.isDir ? '—' : formatBytes(entry.size) }}
    </span>
  </div>

  <!-- Image preview tooltip -->
  <ImagePreviewTooltip
    :file-path="entry.sourcePath"
    :visible="previewVisible"
    :x="previewX"
    :y="previewY"
  />

  <!-- Recursive children -->
  <template v-if="entry.isDir && isExpanded(entry)">
    <FileBrowserItem
      v-for="child in getChildren(entry)"
      :key="child.sourcePath"
      :entry="child"
      :depth="depth + 1"
      :expanded-dirs="expandedDirs"
      :dir-children="dirChildren"
      :selected-paths="selectedPaths"
      :show-hidden="showHidden"
      :sort-fn="sortFn"
      :focused-path="focusedPath"
      @toggle-expand="emit('toggle-expand', $event)"
      @toggle-selection="(entry, ev) => emit('toggle-selection', entry, ev)"
      @dblclick="emit('dblclick', $event)"
      @contextmenu="(entry, ev) => emit('contextmenu', entry, ev)"
      @dragstart="(entry, ev) => emit('dragstart', entry, ev)"
    />
  </template>
</template>
