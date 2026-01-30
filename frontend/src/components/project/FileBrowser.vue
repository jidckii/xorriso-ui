<script setup>
import { ref, reactive, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ArrowUp, Home, Usb, Eye, EyeOff } from 'lucide-vue-next'
import FileBrowserItem from './FileBrowserItem.vue'
import { useProjectStore } from '../../stores/projectStore'
import { useTabStore } from '../../stores/tabStore'

const { t } = useI18n()
const projectStore = useProjectStore()
const tabStore = useTabStore()

const currentProject = computed(() => tabStore.activeProject)
const tabId = computed(() => tabStore.activeTabId)

// Flat list of browse entries (current directory contents)
const entries = ref([])
const selectedPaths = ref(new Set())

// Mount points for quick navigation
const mountPoints = ref([])

// Show/hide hidden files
const showHidden = computed({
  get: () => currentProject.value?.browseShowHidden ?? false,
  set: (val) => {
    if (currentProject.value) {
      currentProject.value.browseShowHidden = val
    }
  },
})

// Filtered entries (hide dot-files unless showHidden)
const filteredEntries = computed(() => {
  if (showHidden.value) return entries.value
  return entries.value.filter(e => !e.name.startsWith('.'))
})

// Editable path (Ctrl+L)
const editingPath = ref(false)
const pathInput = ref('')
const pathInputRef = ref(null)

function startEditPath() {
  pathInput.value = currentProject.value?.browsePath || '/'
  editingPath.value = true
  nextTick(() => {
    pathInputRef.value?.focus()
    pathInputRef.value?.select()
  })
}

function confirmPath() {
  editingPath.value = false
  const path = pathInput.value.trim()
  if (path) {
    navigateTo(path)
  }
}

function cancelEditPath() {
  editingPath.value = false
}

function onKeydown(e) {
  if (e.ctrlKey && e.key === 'l') {
    e.preventDefault()
    startEditPath()
  }
}

// Breadcrumb segments
const breadcrumbs = computed(() => {
  const path = currentProject.value?.browsePath || '/'
  if (path === '/') return []
  const parts = path.split('/').filter(Boolean)
  return parts.map((part, i) => ({
    name: part,
    path: '/' + parts.slice(0, i + 1).join('/'),
  }))
})

// Load directory contents
async function loadDirectory(path) {
  const result = await projectStore.browseDirectory(path)
  entries.value = result || []
  selectedPaths.value = new Set()
  if (currentProject.value) {
    currentProject.value.selectedBrowseFiles = []
  }
}

// Navigate to path
async function navigateTo(path) {
  if (!currentProject.value) return
  currentProject.value.browsePath = path
  await loadDirectory(path)
}

function goUp() {
  const path = currentProject.value?.browsePath || '/'
  if (path === '/') return
  const parts = path.split('/').filter(Boolean)
  parts.pop()
  navigateTo(parts.length === 0 ? '/' : '/' + parts.join('/'))
}

async function goHome() {
  const home = await projectStore.getHomeDirectory()
  navigateTo(home || '/')
}

// Expand/collapse directories inline
const expandedDirs = ref(new Set())
const dirChildren = reactive({}) // sourcePath -> FileEntry[]

async function toggleExpand(entry) {
  const key = entry.sourcePath
  if (expandedDirs.value.has(key)) {
    expandedDirs.value.delete(key)
    expandedDirs.value = new Set(expandedDirs.value)
  } else {
    // Load children if not yet loaded
    if (!dirChildren[key]) {
      const children = await projectStore.browseDirectory(key)
      dirChildren[key] = children || []
    }
    expandedDirs.value.add(key)
    expandedDirs.value = new Set(expandedDirs.value)
  }
}

// Selection (multi-select with Ctrl/Cmd)
function toggleSelection(entry, event) {
  const key = entry.sourcePath
  if (event.ctrlKey || event.metaKey) {
    // Toggle individual item
    if (selectedPaths.value.has(key)) {
      selectedPaths.value.delete(key)
    } else {
      selectedPaths.value.add(key)
    }
    selectedPaths.value = new Set(selectedPaths.value)
  } else {
    // Replace selection
    selectedPaths.value = new Set([key])
  }
  syncSelection()
}

function syncSelection() {
  if (currentProject.value) {
    currentProject.value.selectedBrowseFiles = [...selectedPaths.value]
  }
}

// Double-click: navigate into directory
function onDblClick(entry) {
  if (entry.isDir) {
    navigateTo(entry.sourcePath)
  }
}

// Add selected to project
async function addSelectedToProject() {
  if (!currentProject.value) return
  const paths = [...selectedPaths.value]
  if (paths.length > 0) {
    await projectStore.addFiles(tabId.value, paths)
    selectedPaths.value = new Set()
    syncSelection()
  }
}

function formatBytes(bytes) {
  return projectStore.formatBytes(bytes)
}

// Load mount points with change detection
let mountPollTimer = null

async function loadMountPoints() {
  const result = await projectStore.listMountPoints()
  // Only update if list changed (compare serialized)
  const newJson = JSON.stringify(result)
  const oldJson = JSON.stringify(mountPoints.value)
  if (newJson !== oldJson) {
    mountPoints.value = result
  }
}

// Init
onMounted(async () => {
  if (currentProject.value) {
    const path = currentProject.value.browsePath
    if (!path || path === '/') {
      const home = await projectStore.getHomeDirectory()
      currentProject.value.browsePath = home || '/'
    }
    await loadDirectory(currentProject.value.browsePath)
  }
  await loadMountPoints()
  // Poll for mount changes every 3 seconds
  mountPollTimer = setInterval(loadMountPoints, 3000)
})

onUnmounted(() => {
  if (mountPollTimer) {
    clearInterval(mountPollTimer)
    mountPollTimer = null
  }
})

watch(tabId, async () => {
  if (currentProject.value) {
    expandedDirs.value = new Set()
    await loadDirectory(currentProject.value.browsePath || '/')
  }
})
</script>

<template>
  <div class="flex flex-col h-full" @keydown="onKeydown" @contextmenu.prevent tabindex="-1">
    <!-- Toolbar: mount points + hidden toggle -->
    <div class="flex items-center gap-1 px-2 py-1 bg-gray-100 dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700/50">
      <button
        v-for="mp in mountPoints"
        :key="mp.path"
        @click="navigateTo(mp.path)"
        class="flex items-center gap-1 px-2 py-0.5 text-xs rounded hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors shrink-0"
        :class="{
          'bg-blue-500/15 text-blue-600 dark:text-blue-400': currentProject?.browsePath?.startsWith(mp.path),
          'text-gray-600 dark:text-gray-400': !currentProject?.browsePath?.startsWith(mp.path),
        }"
        :title="mp.path"
      >
        <Home v-if="mp.icon === 'home'" :size="12" />
        <Usb v-else :size="12" />
        <span>{{ mp.label }}</span>
      </button>

      <span class="flex-1" />

      <button
        @click="showHidden = !showHidden"
        class="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors shrink-0"
        :title="showHidden ? t('project.hideHidden') : t('project.showHidden')"
      >
        <Eye v-if="showHidden" :size="14" class="text-blue-500" />
        <EyeOff v-else :size="14" class="text-gray-500" />
      </button>
    </div>

    <!-- Path bar: breadcrumb or editable input -->
    <div class="flex items-center gap-1 px-3 py-1.5 bg-gray-50 dark:bg-gray-800/50 border-b border-gray-300 dark:border-gray-700 min-h-[34px]">
      <button
        @click="goUp"
        class="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors shrink-0"
        :title="t('project.goUp')"
      >
        <ArrowUp :size="14" class="text-gray-600 dark:text-gray-400" />
      </button>

      <!-- Editable path input -->
      <input
        v-if="editingPath"
        ref="pathInputRef"
        v-model="pathInput"
        @keydown.enter="confirmPath"
        @keydown.escape="cancelEditPath"
        @blur="cancelEditPath"
        class="flex-1 ml-1 px-2 py-0.5 text-xs bg-white dark:bg-gray-900 border border-blue-500 rounded text-gray-800 dark:text-gray-200 outline-none"
      />

      <!-- Breadcrumb (click empty area to edit) -->
      <div
        v-else
        class="flex items-center gap-0.5 text-xs text-gray-600 dark:text-gray-400 overflow-hidden flex-1 ml-1 cursor-text"
        @click.self="startEditPath"
      >
        <span
          class="cursor-pointer hover:text-gray-800 dark:hover:text-gray-200 transition-colors shrink-0"
          @click="navigateTo('/')"
        >/</span>
        <template v-for="(crumb, idx) in breadcrumbs" :key="crumb.path">
          <span
            class="cursor-pointer hover:text-gray-800 dark:hover:text-gray-200 transition-colors truncate"
            :class="{ 'text-gray-800 dark:text-gray-200 font-medium': idx === breadcrumbs.length - 1 }"
            @click="navigateTo(crumb.path)"
          >{{ crumb.name }}</span>
          <span v-if="idx < breadcrumbs.length - 1" class="text-gray-400 dark:text-gray-600 shrink-0">/</span>
        </template>
        <!-- Invisible spacer to capture clicks on empty area -->
        <span class="flex-1" @click="startEditPath" />
      </div>
    </div>

    <!-- File list with inline expand -->
    <div class="flex-1 overflow-y-auto text-sm select-none">
      <div v-if="filteredEntries.length === 0" class="text-center text-gray-500 py-8">
        {{ t('project.emptyDirectory') }}
      </div>

      <FileBrowserItem
        v-for="entry in filteredEntries"
        :key="entry.sourcePath"
        :entry="entry"
        :depth="0"
        :expanded-dirs="expandedDirs"
        :dir-children="dirChildren"
        :selected-paths="selectedPaths"
        :format-bytes="formatBytes"
        :show-hidden="showHidden"
        @toggle-expand="toggleExpand"
        @toggle-selection="toggleSelection"
        @dblclick="onDblClick"
      />
    </div>

    <!-- Toolbar -->
    <div class="flex items-center gap-2 px-3 py-2 bg-gray-100 dark:bg-gray-800 border-t border-gray-300 dark:border-gray-700">
      <button
        @click="addSelectedToProject"
        :disabled="selectedPaths.size === 0"
        class="px-3 py-1 text-xs font-medium rounded bg-blue-600 hover:bg-blue-500 text-white disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
      >
        {{ t('project.addToProject') }}
      </button>
      <span class="text-xs text-gray-500">
        {{ selectedPaths.size }} {{ t('project.selected') }}
      </span>
    </div>
  </div>
</template>
