<script setup>
import { ref, reactive, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ExternalLink, FolderOpen, Trash2, Info } from 'lucide-vue-next'
import PanelHeader from '../ui/PanelHeader.vue'
import DiscLayoutTree from './DiscLayoutTree.vue'
import DiscLayoutToolbar from './DiscLayoutToolbar.vue'
import ContextMenu from '../ui/ContextMenu.vue'
import FilePropertiesModal from './FilePropertiesModal.vue'
import SortButtons from '../ui/SortButtons.vue'
import { useProjectStore } from '../../stores/projectStore'
import { useTabStore } from '../../stores/tabStore'
import { formatBytes } from '../../composables/useFormatBytes'
import { useFileSort } from '../../composables/useFileSort'

const { t } = useI18n()
const router = useRouter()
const projectStore = useProjectStore()
const tabStore = useTabStore()

const currentProject = computed(() => tabStore.activeProject)
const tabId = computed(() => tabStore.activeTabId)

// Сортировка дерева диска
const { sortBy: discSortBy, sortDir: discSortDir, toggleSort: discToggleSort, compareFn: discCompareFn } = useFileSort(ref([]))

// Рекурсивная сортировка дочерних узлов дерева
function sortTreeChildren(items) {
  if (!items || items.length === 0) return items
  return [...items].sort(discCompareFn).map(item => ({
    ...item,
    children: sortTreeChildren(item.children),
  }))
}

// Построение дерева из плоского списка записей
const treeItems = computed(() => {
  const entries = currentProject.value?.entries || []
  if (entries.length === 0) return []
  const tree = buildTreeFromEntries(entries)
  return sortTreeChildren(tree)
})

// Drag-and-Drop состояние
const isDragOver = ref(false)

function buildTreeFromEntries(entries) {
  const root = { children: [] }

  for (const entry of entries) {
    const parts = entry.destPath.split('/').filter(Boolean)
    let current = root

    for (let i = 0; i < parts.length; i++) {
      const name = parts[i]
      let child = current.children.find(c => c.name === name)

      if (!child) {
        const isLeaf = i === parts.length - 1
        child = {
          _key: '/' + parts.slice(0, i + 1).join('/'),
          destPath: '/' + parts.slice(0, i + 1).join('/'),
          sourcePath: isLeaf ? entry.sourcePath : '',
          name,
          isDir: isLeaf ? entry.isDir : true,
          size: isLeaf ? entry.size : 0,
          children: [],
        }
        current.children.push(child)
      }
      current = child
    }
  }

  return root.children
}

const expanded = ref([])

// Ручное управление выделением с пробросом на дочерние элементы
const selectedKeys = ref(new Set())

function toggleItemSelection(item) {
  const selecting = !selectedKeys.value.has(item._key)

  if (selecting) {
    selectedKeys.value.add(item._key)
  } else {
    selectedKeys.value.delete(item._key)
  }

  // Пробрасываем выбор на дочерние элементы
  if (item.children && item.children.length > 0) {
    propagateSelection(item.children, selecting)
  }

  selectedKeys.value = new Set(selectedKeys.value)
  syncProjectSelection()
}

function propagateSelection(children, selecting) {
  for (const child of children) {
    if (selecting) {
      selectedKeys.value.add(child._key)
    } else {
      selectedKeys.value.delete(child._key)
    }
    if (child.children && child.children.length > 0) {
      propagateSelection(child.children, selecting)
    }
  }
}

function syncProjectSelection() {
  if (currentProject.value) {
    currentProject.value.selectedProjectEntries = [...selectedKeys.value]
  }
}

const selectedCount = computed(() => selectedKeys.value.size)

async function removeSelectedFromProject() {
  if (!currentProject.value) return
  // Собираем все выбранные пути и удаляем одним пакетным IPC-вызовом
  const paths = [...currentProject.value.selectedProjectEntries]
  await projectStore.removeEntries(tabId.value, paths)
  currentProject.value.selectedProjectEntries = []
  selectedKeys.value = new Set()
}

// Выделить все / снять выделение
function selectAllItems(items) {
  for (const item of items) {
    selectedKeys.value.add(item._key)
    if (item.children && item.children.length > 0) {
      selectAllItems(item.children)
    }
  }
}

function selectAll() {
  selectAllItems(treeItems.value)
  selectedKeys.value = new Set(selectedKeys.value)
  syncProjectSelection()
}

function deselectAll() {
  selectedKeys.value = new Set()
  syncProjectSelection()
}

const allSelected = computed(() => {
  if (treeItems.value.length === 0) return false
  return selectedKeys.value.size > 0 && countAllItems(treeItems.value) === selectedKeys.value.size
})

function countAllItems(items) {
  let count = 0
  for (const item of items) {
    count++
    if (item.children && item.children.length > 0) {
      count += countAllItems(item.children)
    }
  }
  return count
}

function goToBurn() {
  router.push('/burn')
}

const canBurn = computed(() => {
  return currentProject.value && currentProject.value.entries.length > 0
})

// Drag-and-Drop обработчики
function onDragOver(event) {
  isDragOver.value = true
  event.dataTransfer.dropEffect = 'copy'
}

function onDragLeave(event) {
  // Проверяем что курсор действительно покинул контейнер (а не перешёл на дочерний элемент)
  if (!event.currentTarget.contains(event.relatedTarget)) {
    isDragOver.value = false
  }
}

async function onDrop(event) {
  isDragOver.value = false
  const raw = event.dataTransfer.getData('application/xorriso-paths')
  if (!raw) return
  try {
    const paths = JSON.parse(raw)
    if (Array.isArray(paths) && paths.length > 0) {
      await projectStore.addFiles(tabId.value, paths)
    }
  } catch {
    // Невалидный JSON — игнорируем
  }
}

// Keyboard shortcuts
function onKeydown(e) {
  switch (e.key) {
    case 'Delete':
      e.preventDefault()
      removeSelectedFromProject()
      break
    case 'Escape':
      e.preventDefault()
      deselectAll()
      break
  }
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
  const items = []
  if (contextMenu.entry.sourcePath) {
    items.push({ label: t('contextMenu.open'), icon: ExternalLink, action: 'open' })
    items.push({ label: t('contextMenu.revealInFileManager'), icon: FolderOpen, action: 'reveal' })
  }
  items.push({ label: t('contextMenu.removeFromProject'), icon: Trash2, action: 'remove' })
  if (contextMenu.entry.sourcePath) {
    items.push({ separator: true })
    items.push({ label: t('contextMenu.properties'), icon: Info, action: 'properties' })
  }
  return items
})

function onContextMenu(item, event) {
  event.preventDefault()
  event.stopPropagation()
  contextMenu.entry = item
  contextMenu.x = event.clientX
  contextMenu.y = event.clientY
  contextMenu.show = true
}

function onContextMenuSelect(action) {
  const entry = contextMenu.entry
  if (!entry) return

  switch (action) {
    case 'open':
      if (entry.sourcePath) {
        projectStore.openWithDefault(entry.sourcePath)
      }
      break
    case 'reveal':
      if (entry.sourcePath) {
        projectStore.revealInFileManager(entry.sourcePath)
      }
      break
    case 'remove':
      if (currentProject.value) {
        projectStore.removeEntries(tabId.value, [entry.destPath])
      }
      break
    case 'properties':
      if (entry.sourcePath) {
        propertiesModal.filePath = entry.sourcePath
        propertiesModal.show = true
      }
      break
  }
}

// File properties modal
const propertiesModal = reactive({
  show: false,
  filePath: '',
})
</script>

<template>
  <div class="flex flex-col h-full outline-none" tabindex="0" @keydown="onKeydown">
    <PanelHeader>
      <template #row1>
        <svg class="w-4 h-4 text-gray-600 dark:text-gray-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="10" stroke-width="1.5" />
          <circle cx="12" cy="12" r="3" stroke-width="1.5" />
        </svg>
        <span class="text-xs font-medium text-gray-700 dark:text-gray-300">{{ t('project.discLayout') }}</span>
        <span class="flex-1" />
        <SortButtons
          :sort-by="discSortBy"
          :sort-dir="discSortDir"
          @toggle-sort="discToggleSort"
        />
      </template>
      <template #row2>
        <span class="text-xs text-gray-500">
          {{ currentProject?.entries?.length || 0 }} {{ t('project.items') }} — {{ formatBytes(currentProject?.totalSize || 0) }}
        </span>
      </template>
    </PanelHeader>

    <!-- Содержимое дерева -->
    <div
      class="flex-1 overflow-y-auto transition-colors"
      :class="{ 'ring-2 ring-blue-500 ring-inset bg-blue-500/5': isDragOver }"
      @dragover.prevent="onDragOver"
      @dragleave="onDragLeave"
      @drop.prevent="onDrop"
    >
      <!-- Пустое состояние -->
      <div
        v-if="!currentProject?.entries?.length"
        class="flex flex-col items-center justify-center h-full text-gray-500 py-12"
      >
        <svg class="w-12 h-12 mb-3 text-gray-200 dark:text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
            d="M12 9v3m0 0v3m0-3h3m-3 0H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-sm">{{ isDragOver ? t('project.dropFilesHere') : t('common.addFilesFromBrowser') }}</p>
      </div>

      <!-- Дерево файлов -->
      <DiscLayoutTree
        v-else
        :items="treeItems"
        v-model:expanded="expanded"
        :selected-keys="selectedKeys"
        @toggle-selection="toggleItemSelection"
        @contextmenu="onContextMenu"
      />
    </div>

    <!-- Context Menu -->
    <ContextMenu
      :show="contextMenu.show"
      :x="contextMenu.x"
      :y="contextMenu.y"
      :items="contextMenuItems"
      @close="contextMenu.show = false"
      @select="onContextMenuSelect"
    />

    <!-- File Properties Modal -->
    <FilePropertiesModal
      :show="propertiesModal.show"
      :file-path="propertiesModal.filePath"
      @close="propertiesModal.show = false"
    />

    <!-- Панель инструментов -->
    <DiscLayoutToolbar
      :all-selected="allSelected"
      :selected-count="selectedCount"
      :can-burn="canBurn"
      @select-all="selectAll"
      @deselect-all="deselectAll"
      @remove-selected="removeSelectedFromProject"
      @go-to-burn="goToBurn"
    />
  </div>
</template>
