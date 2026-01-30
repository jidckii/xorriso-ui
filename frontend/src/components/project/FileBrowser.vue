<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Button from '../ui/Button.vue'

const { t } = useI18n()

const emit = defineEmits(['add-files'])

// Mock file system data (will be replaced by Wails backend calls)
const currentPath = ref('/home')
const entries = ref([
  { name: 'Documents', path: '/home/Documents', isDir: true, size: 0 },
  { name: 'Downloads', path: '/home/Downloads', isDir: true, size: 0 },
  { name: 'Music', path: '/home/Music', isDir: true, size: 0 },
  { name: 'Pictures', path: '/home/Pictures', isDir: true, size: 0 },
  { name: 'readme.txt', path: '/home/readme.txt', isDir: false, size: 4096 },
  { name: 'backup.tar.gz', path: '/home/backup.tar.gz', isDir: false, size: 52428800 },
])

const selectedEntries = ref(new Set())
const pathHistory = ref(['/home'])

const breadcrumbs = computed(() => {
  const parts = currentPath.value.split('/').filter(Boolean)
  return parts.map((part, i) => ({
    name: part,
    path: '/' + parts.slice(0, i + 1).join('/'),
  }))
})

function navigateTo(path) {
  currentPath.value = path
  selectedEntries.value.clear()
  // In real app: call Wails backend to list directory
  // e.g. entries.value = await ListDirectory(path)
}

function goUp() {
  const parts = currentPath.value.split('/').filter(Boolean)
  if (parts.length > 1) {
    parts.pop()
    navigateTo('/' + parts.join('/'))
  } else {
    navigateTo('/')
  }
}

function onEntryClick(entry) {
  if (entry.isDir) {
    return
  }
  if (selectedEntries.value.has(entry.path)) {
    selectedEntries.value.delete(entry.path)
  } else {
    selectedEntries.value.add(entry.path)
  }
}

function onEntryDblClick(entry) {
  if (entry.isDir) {
    navigateTo(entry.path)
  } else {
    addSelected([entry])
  }
}

function addSelected(entriesToAdd) {
  const items = entriesToAdd || entries.value.filter((e) => selectedEntries.value.has(e.path))
  if (items.length > 0) {
    emit('add-files', items)
    selectedEntries.value.clear()
  }
}

function formatSize(bytes) {
  if (!bytes) return ''
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
  <div class="flex flex-col h-full bg-white dark:bg-gray-900 rounded-lg border border-gray-300 dark:border-gray-700">
    <!-- Header -->
    <div class="flex items-center gap-2 px-3 py-2 border-b border-gray-300 dark:border-gray-700 bg-gray-100 dark:bg-gray-800 rounded-t-lg">
      <Button variant="ghost" size="sm" @click="goUp">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18" />
        </svg>
      </Button>

      <!-- Breadcrumb -->
      <div class="flex items-center gap-1 text-sm text-gray-600 dark:text-gray-400 overflow-hidden flex-1">
        <span
          class="cursor-pointer hover:text-gray-800 dark:hover:text-gray-200 transition-colors flex-shrink-0"
          @click="navigateTo('/')"
        >/</span>
        <template v-for="(crumb, idx) in breadcrumbs" :key="crumb.path">
          <span
            class="cursor-pointer hover:text-gray-800 dark:hover:text-gray-200 transition-colors truncate"
            :class="{ 'text-gray-800 dark:text-gray-200': idx === breadcrumbs.length - 1 }"
            @click="navigateTo(crumb.path)"
          >
            {{ crumb.name }}
          </span>
          <span v-if="idx < breadcrumbs.length - 1" class="text-gray-500 dark:text-gray-600 flex-shrink-0">/</span>
        </template>
      </div>

      <Button
        variant="primary"
        size="sm"
        :disabled="selectedEntries.size === 0"
        @click="addSelected()"
      >
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        {{ t('common.add') }}
      </Button>
    </div>

    <!-- File list -->
    <div class="flex-1 overflow-y-auto">
      <div
        v-for="entry in entries"
        :key="entry.path"
        class="flex items-center gap-2 px-3 py-1.5 cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
        :class="{
          'bg-blue-900/30': selectedEntries.has(entry.path),
        }"
        @click="onEntryClick(entry)"
        @dblclick="onEntryDblClick(entry)"
      >
        <!-- Icon -->
        <svg
          v-if="entry.isDir"
          class="w-5 h-5 text-yellow-500 flex-shrink-0"
          fill="currentColor"
          viewBox="0 0 20 20"
        >
          <path d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" />
        </svg>
        <svg
          v-else
          class="w-5 h-5 text-gray-600 dark:text-gray-400 flex-shrink-0"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
        </svg>

        <!-- Name -->
        <span class="text-sm text-gray-800 dark:text-gray-200 truncate flex-1">
          {{ entry.name }}
        </span>

        <!-- Size -->
        <span v-if="!entry.isDir" class="text-xs text-gray-500 flex-shrink-0">
          {{ formatSize(entry.size) }}
        </span>
      </div>

      <div v-if="entries.length === 0" class="text-center text-gray-500 text-sm py-8">
        {{ t('project.emptyDirectory') }}
      </div>
    </div>
  </div>
</template>
