<script setup>
import { ref } from 'vue'

const props = defineProps({
  entries: {
    type: Array,
    default: () => [],
    // Each entry: { name, path, isDir, size, children? }
  },
  selectable: { type: Boolean, default: false },
  depth: { type: Number, default: 0 },
})

const emit = defineEmits(['select', 'dblclick'])

const expandedPaths = ref(new Set())
const selectedPath = ref('')

function toggleExpand(entry) {
  if (!entry.isDir) return
  if (expandedPaths.value.has(entry.path)) {
    expandedPaths.value.delete(entry.path)
  } else {
    expandedPaths.value.add(entry.path)
  }
}

function isExpanded(entry) {
  return expandedPaths.value.has(entry.path)
}

function selectEntry(entry) {
  selectedPath.value = entry.path
  emit('select', entry)
}

function dblclickEntry(entry) {
  emit('dblclick', entry)
}

function formatSize(bytes) {
  if (bytes == null || bytes === 0) return ''
  if (bytes < 1024) return bytes + ' B'
  const units = ['KB', 'MB', 'GB']
  let i = -1
  let size = bytes
  do {
    size /= 1024
    i++
  } while (size >= 1024 && i < units.length - 1)
  return size.toFixed(1) + ' ' + units[i]
}
</script>

<template>
  <div class="text-sm">
    <div
      v-for="entry in entries"
      :key="entry.path"
    >
      <div
        class="flex items-center gap-1 px-2 py-0.5 rounded cursor-pointer hover:bg-gray-200/50 dark:hover:bg-gray-700/50 transition-colors"
        :class="{
          'bg-blue-900/30': selectable && selectedPath === entry.path,
        }"
        :style="{ paddingLeft: depth * 16 + 8 + 'px' }"
        @click="entry.isDir ? toggleExpand(entry) : selectEntry(entry)"
        @dblclick="dblclickEntry(entry)"
      >
        <!-- Expand arrow for dirs -->
        <svg
          v-if="entry.isDir"
          class="w-3 h-3 text-gray-500 transition-transform flex-shrink-0"
          :class="{ 'rotate-90': isExpanded(entry) }"
          fill="currentColor"
          viewBox="0 0 20 20"
        >
          <path d="M6 6l8 4-8 4V6z" />
        </svg>
        <span v-else class="w-3" />

        <!-- Icon -->
        <svg
          v-if="entry.isDir"
          class="w-4 h-4 text-yellow-500 flex-shrink-0"
          fill="currentColor"
          viewBox="0 0 20 20"
        >
          <path d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" />
        </svg>
        <svg
          v-else
          class="w-4 h-4 text-gray-600 dark:text-gray-400 flex-shrink-0"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
        </svg>

        <!-- Name -->
        <span class="text-gray-800 dark:text-gray-200 truncate flex-1" @click.stop="selectEntry(entry)">
          {{ entry.name }}
        </span>

        <!-- Size -->
        <span v-if="!entry.isDir && entry.size" class="text-gray-500 text-xs flex-shrink-0 ml-2">
          {{ formatSize(entry.size) }}
        </span>
      </div>

      <!-- Children (recursive) -->
      <div v-if="entry.isDir && isExpanded(entry) && entry.children">
        <FileTree
          :entries="entry.children"
          :selectable="selectable"
          :depth="depth + 1"
          @select="$emit('select', $event)"
          @dblclick="$emit('dblclick', $event)"
        />
      </div>
    </div>
  </div>
</template>
