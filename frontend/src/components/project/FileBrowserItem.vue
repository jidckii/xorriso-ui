<script setup>
import { ChevronRight } from 'lucide-vue-next'
import FileIcon from '../ui/FileIcon.vue'

const props = defineProps({
  entry: { type: Object, required: true },
  depth: { type: Number, default: 0 },
  expandedDirs: { type: Set, required: true },
  dirChildren: { type: Object, required: true },
  selectedPaths: { type: Set, required: true },
  formatBytes: { type: Function, required: true },
})

const emit = defineEmits(['toggle-expand', 'toggle-selection', 'dblclick'])

function isExpanded(entry) {
  return props.expandedDirs.has(entry.sourcePath)
}

function isSelected(entry) {
  return props.selectedPaths.has(entry.sourcePath)
}

function getChildren(entry) {
  return props.dirChildren[entry.sourcePath] || []
}
</script>

<template>
  <!-- Entry row -->
  <div
    class="flex items-center gap-1.5 py-1 cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-800/50 transition-colors"
    :class="{ 'bg-blue-500/15': isSelected(entry) }"
    :style="{ paddingLeft: (depth * 16 + 8) + 'px', paddingRight: '8px' }"
    @click="emit('toggle-selection', entry, $event)"
    @dblclick="emit('dblclick', entry)"
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

    <!-- Size -->
    <span v-if="!entry.isDir" class="text-xs text-gray-500 shrink-0 ml-2">
      {{ formatBytes(entry.size) }}
    </span>
  </div>

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
      :format-bytes="formatBytes"
      @toggle-expand="emit('toggle-expand', $event)"
      @toggle-selection="(entry, ev) => emit('toggle-selection', entry, ev)"
      @dblclick="emit('dblclick', $event)"
    />
  </template>
</template>
