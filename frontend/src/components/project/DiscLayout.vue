<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import FileTree from './FileTree.vue'
import Button from '../ui/Button.vue'

const { t } = useI18n()

const props = defineProps({
  entries: {
    type: Array,
    default: () => [],
    // Each entry: { name, path, isDir, size, children? }
  },
})

const emit = defineEmits(['remove', 'update:entries'])

const dragOver = ref(false)

const isEmpty = computed(() => props.entries.length === 0)

const totalSize = computed(() => {
  function sumSize(items) {
    return items.reduce((acc, item) => {
      let s = item.size || 0
      if (item.children) s += sumSize(item.children)
      return acc + s
    }, 0)
  }
  return sumSize(props.entries)
})

function formatSize(bytes) {
  if (!bytes) return '0 B'
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

function removeEntry(entry) {
  emit('remove', entry)
}

function onDragOver(e) {
  e.preventDefault()
  dragOver.value = true
}

function onDragLeave() {
  dragOver.value = false
}

function onDrop(e) {
  e.preventDefault()
  dragOver.value = false
  // Handle drop - in real app, parse dragged data
}

function onSelect(entry) {
  // Could highlight the entry
}
</script>

<template>
  <div
    class="flex flex-col h-full bg-white dark:bg-gray-900 rounded-lg border border-gray-300 dark:border-gray-700"
    :class="{ 'border-blue-500 border-2': dragOver }"
    @dragover="onDragOver"
    @dragleave="onDragLeave"
    @drop="onDrop"
  >
    <!-- Header -->
    <div class="flex items-center justify-between px-3 py-2 border-b border-gray-300 dark:border-gray-700 bg-gray-100 dark:bg-gray-800 rounded-t-lg">
      <div class="flex items-center gap-2">
        <svg class="w-4 h-4 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="10" stroke-width="1.5" />
          <circle cx="12" cy="12" r="3" stroke-width="1.5" />
        </svg>
        <span class="text-sm font-medium text-gray-800 dark:text-gray-200">{{ t('project.discLayout') }}</span>
      </div>
      <span class="text-xs text-gray-600 dark:text-gray-400">
        {{ entries.length }} {{ t('project.items') }}, {{ formatSize(totalSize) }}
      </span>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto">
      <!-- Empty state -->
      <div
        v-if="isEmpty"
        class="flex flex-col items-center justify-center h-full text-gray-500 py-12"
      >
        <svg class="w-16 h-16 mb-4 opacity-30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="10" stroke-width="1" />
          <circle cx="12" cy="12" r="3" stroke-width="1" />
          <circle cx="12" cy="12" r="6" stroke-width="0.5" />
        </svg>
        <p class="text-sm mb-1">{{ t('project.noFilesAdded') }}</p>
        <p class="text-xs text-gray-500 dark:text-gray-600">{{ t('project.dragFilesHere') }}</p>
      </div>

      <!-- File tree with remove buttons -->
      <div v-else class="py-1">
        <div
          v-for="entry in entries"
          :key="entry.path"
          class="group flex items-center gap-1 px-2 py-1 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
        >
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
          <span class="text-sm text-gray-800 dark:text-gray-200 truncate flex-1">{{ entry.name }}</span>

          <!-- Size -->
          <span class="text-xs text-gray-500 flex-shrink-0">
            {{ formatSize(entry.size) }}
          </span>

          <!-- Remove button -->
          <button
            class="opacity-0 group-hover:opacity-100 text-gray-500 hover:text-red-400 transition-all flex-shrink-0 ml-1"
            :title="t('project.remove')"
            @click="removeEntry(entry)"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
