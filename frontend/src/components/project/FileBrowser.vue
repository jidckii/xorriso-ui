<script setup>
import { ref, reactive, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ExternalLink, FolderOpen, Info, FolderPlus } from 'lucide-vue-next'
import FileBrowserItem from './FileBrowserItem.vue'
import FileBrowserToolbar from './FileBrowserToolbar.vue'
import FileBrowserSelectionBar from './FileBrowserSelectionBar.vue'
import ContextMenu from '../ui/ContextMenu.vue'
import FilePropertiesModal from './FilePropertiesModal.vue'
import { useProjectStore } from '../../stores/projectStore'
import { useTabStore } from '../../stores/tabStore'
import { formatBytes } from '../../composables/useFormatBytes'
import { useFileSort } from '../../composables/useFileSort'

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

// Сортировка файлов
const { sortBy, sortDir, sorted: sortedEntries, toggleSort, sortChildren } = useFileSort(filteredEntries)

// Keyboard navigation: focused item index in flat visible list
const focusedIndex = ref(-1)
const listRef = ref(null)

const flatVisibleEntries = computed(() => {
  const result = []
  function collect(entries) {
    for (const entry of entries) {
      result.push(entry)
      if (entry.isDir && expandedDirs.value.has(entry.sourcePath)) {
        const children = dirChildren[entry.sourcePath] || []
        const filtered = showHidden.value ? children : children.filter(c => !c.name.startsWith('.'))
        const sorted = sortChildren ? sortChildren(filtered) : filtered
        collect(sorted)
      }
    }
  }
  collect(sortedEntries.value)
  return result
})

const focusedEntry = computed(() => {
  const flat = flatVisibleEntries.value
  if (focusedIndex.value < 0 || focusedIndex.value >= flat.length) return null
  return flat[focusedIndex.value]
})

watch(flatVisibleEntries, () => {
  if (focusedIndex.value >= flatVisibleEntries.value.length) {
    focusedIndex.value = Math.max(0, flatVisibleEntries.value.length - 1)
  }
})

function scrollFocusedIntoView() {
  nextTick(() => {
    if (!focusedEntry.value || !listRef.value) return
    const el = listRef.value.querySelector(`[data-path="${CSS.escape(focusedEntry.value.sourcePath)}"]`)
    el?.scrollIntoView({ block: 'nearest' })
  })
}

// Редактируемый путь (Ctrl+L)
const editingPath = ref(false)
const pathInput = ref('')

function startEditPath() {
  pathInput.value = currentProject.value?.browsePath || '/'
  editingPath.value = true
}

function onConfirmPath(path) {
  editingPath.value = false
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
    return
  }

  // Keyboard navigation in file list
  const flat = flatVisibleEntries.value
  if (flat.length === 0) return

  switch (e.key) {
    case 'ArrowDown':
      e.preventDefault()
      if (focusedIndex.value < 0) {
        focusedIndex.value = 0
      } else {
        focusedIndex.value = Math.min(focusedIndex.value + 1, flat.length - 1)
      }
      scrollFocusedIntoView()
      break
    case 'ArrowUp':
      e.preventDefault()
      if (focusedIndex.value < 0) {
        focusedIndex.value = 0
      } else {
        focusedIndex.value = Math.max(focusedIndex.value - 1, 0)
      }
      scrollFocusedIntoView()
      break
    case 'ArrowRight': {
      e.preventDefault()
      const entry = focusedEntry.value
      if (entry?.isDir && !expandedDirs.value.has(entry.sourcePath)) {
        toggleExpand(entry)
      }
      break
    }
    case 'ArrowLeft': {
      e.preventDefault()
      const entry = focusedEntry.value
      if (entry?.isDir && expandedDirs.value.has(entry.sourcePath)) {
        toggleExpand(entry)
      }
      break
    }
    case ' ':
      e.preventDefault()
      if (focusedEntry.value) {
        toggleSelection(focusedEntry.value, e)
      }
      break
    case 'Enter':
      e.preventDefault()
      addSelectedToProject()
      break
    case 'Escape':
      e.preventDefault()
      deselectAll()
      focusedIndex.value = -1
      break
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

// Selection (always toggle mode for multi-select)
function toggleSelection(entry, event) {
  const key = entry.sourcePath
  const selecting = !selectedPaths.value.has(key)

  if (selecting) {
    selectedPaths.value.add(key)
  } else {
    selectedPaths.value.delete(key)
  }

  // Propagate to loaded children if directory
  if (entry.isDir) {
    propagateSelection(key, selecting)
  }

  selectedPaths.value = new Set(selectedPaths.value)
  syncSelection()
}

// Recursively select/deselect loaded children of a directory
function propagateSelection(dirPath, selecting) {
  const children = dirChildren[dirPath]
  if (!children) return
  for (const child of children) {
    if (selecting) {
      selectedPaths.value.add(child.sourcePath)
    } else {
      selectedPaths.value.delete(child.sourcePath)
    }
    if (child.isDir) {
      propagateSelection(child.sourcePath, selecting)
    }
  }
}

function selectAll() {
  selectedPaths.value = new Set(filteredEntries.value.map(e => e.sourcePath))
  syncSelection()
}

function deselectAll() {
  selectedPaths.value = new Set()
  syncSelection()
}

const allSelected = computed(() =>
  filteredEntries.value.length > 0 &&
  filteredEntries.value.every(e => selectedPaths.value.has(e.sourcePath))
)

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
  const allPaths = [...selectedPaths.value]
  if (allPaths.length === 0) return

  // Filter out paths whose parent directory is already in the selection,
  // because AddFiles recursively expands directories
  const pathSet = new Set(allPaths)
  const paths = allPaths.filter(p => {
    let parent = p.substring(0, p.lastIndexOf('/'))
    while (parent) {
      if (pathSet.has(parent)) return false
      parent = parent.substring(0, parent.lastIndexOf('/'))
    }
    return true
  })

  await projectStore.addFiles(tabId.value, paths)
  selectedPaths.value = new Set()
  syncSelection()
}

// Drag-and-Drop: формируем список путей для переноса
function onDragStart(entry, event) {
  let paths

  if (selectedPaths.value.has(entry.sourcePath)) {
    // Файл выделен — тащим всё выделенное
    const allPaths = [...selectedPaths.value]
    const pathSet = new Set(allPaths)
    paths = allPaths.filter(p => {
      let parent = p.substring(0, p.lastIndexOf('/'))
      while (parent) {
        if (pathSet.has(parent)) return false
        parent = parent.substring(0, parent.lastIndexOf('/'))
      }
      return true
    })
  } else {
    // Файл не выделен — тащим только его
    paths = [entry.sourcePath]
  }

  event.dataTransfer.setData('application/xorriso-paths', JSON.stringify(paths))
  event.dataTransfer.effectAllowed = 'copy'
}

// Context menu
const contextMenu = reactive({
  show: false,
  x: 0,
  y: 0,
  entry: null,
})

const contextMenuItems = computed(() => {
  if (!contextMenu.entry) return []
  const items = [
    { label: t('contextMenu.open'), icon: ExternalLink, action: 'open' },
    { label: t('contextMenu.revealInFileManager'), icon: FolderOpen, action: 'reveal' },
    { label: t('contextMenu.addToProject'), icon: FolderPlus, action: 'add' },
    { separator: true },
    { label: t('contextMenu.properties'), icon: Info, action: 'properties' },
  ]
  return items
})

function onContextMenu(entry, event) {
  event.preventDefault()
  event.stopPropagation()
  contextMenu.entry = entry
  contextMenu.x = event.clientX
  contextMenu.y = event.clientY
  contextMenu.show = true
}

function onContextMenuSelect(action) {
  const entry = contextMenu.entry
  if (!entry) return

  switch (action) {
    case 'open':
      projectStore.openWithDefault(entry.sourcePath)
      break
    case 'reveal':
      projectStore.revealInFileManager(entry.sourcePath)
      break
    case 'add':
      if (currentProject.value) {
        projectStore.addFiles(tabId.value, [entry.sourcePath])
      }
      break
    case 'properties':
      propertiesModal.filePath = entry.sourcePath
      propertiesModal.show = true
      break
  }
}

// File properties modal
const propertiesModal = reactive({
  show: false,
  filePath: '',
})

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
  <div class="flex flex-col h-full outline-none" @keydown="onKeydown" @contextmenu.prevent tabindex="0">
    <FileBrowserToolbar
      :mount-points="mountPoints"
      :browse-path="currentProject?.browsePath || '/'"
      :breadcrumbs="breadcrumbs"
      :editing-path="editingPath"
      v-model:path-input="pathInput"
      :show-hidden="showHidden"
      :sort-by="sortBy"
      :sort-dir="sortDir"
      @navigate="navigateTo"
      @go-up="goUp"
      @start-edit="startEditPath"
      @confirm-path="onConfirmPath"
      @cancel-edit="cancelEditPath"
      @toggle-hidden="showHidden = !showHidden"
      @toggle-sort="toggleSort"
    />

    <!-- File list with inline expand -->
    <div ref="listRef" class="flex-1 overflow-y-auto text-sm select-none">
      <div v-if="sortedEntries.length === 0" class="text-center text-gray-500 py-8">
        {{ t('project.emptyDirectory') }}
      </div>

      <FileBrowserItem
        v-for="entry in sortedEntries"
        :key="entry.sourcePath"
        :entry="entry"
        :depth="0"
        :expanded-dirs="expandedDirs"
        :dir-children="dirChildren"
        :selected-paths="selectedPaths"
        :show-hidden="showHidden"
        :sort-fn="sortChildren"
        :focused-path="focusedEntry?.sourcePath"
        @toggle-expand="toggleExpand"
        @toggle-selection="toggleSelection"
        @dblclick="onDblClick"
        @contextmenu="onContextMenu"
        @dragstart="onDragStart"
      />
    </div>

    <!-- Нижняя панель выбора -->
    <FileBrowserSelectionBar
      :all-selected="allSelected"
      :selected-count="selectedPaths.size"
      @select-all="selectAll"
      @deselect-all="deselectAll"
      @add-selected="addSelectedToProject"
    />

    <!-- Context menu -->
    <ContextMenu
      :show="contextMenu.show"
      :x="contextMenu.x"
      :y="contextMenu.y"
      :items="contextMenuItems"
      @close="contextMenu.show = false"
      @select="onContextMenuSelect"
    />

    <!-- File properties modal -->
    <FilePropertiesModal
      :show="propertiesModal.show"
      :file-path="propertiesModal.filePath"
      @close="propertiesModal.show = false"
    />
  </div>
</template>
