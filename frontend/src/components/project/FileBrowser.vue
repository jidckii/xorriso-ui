<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { TreeRoot, TreeItem, TreeVirtualizer } from 'reka-ui'
import { ChevronRight, ArrowUp, Home } from 'lucide-vue-next'
import FileIcon from '../ui/FileIcon.vue'
import { useProjectStore } from '../../stores/projectStore'
import { useTabStore } from '../../stores/tabStore'

const { t } = useI18n()
const projectStore = useProjectStore()
const tabStore = useTabStore()

const currentProject = computed(() => tabStore.activeProject)
const tabId = computed(() => tabStore.activeTabId)

// Tree data: enriched browse entries with lazy-loaded children
const treeItems = ref([])
const expanded = ref([])
const selected = ref([])

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

// Load root directory
async function loadDirectory(path) {
  const entries = await projectStore.browseDirectory(path)
  treeItems.value = entries.map(e => ({
    ...e,
    _key: e.sourcePath,
    children: e.isDir ? undefined : undefined,
    _loaded: false,
  }))
  expanded.value = []
  selected.value = []
}

// Navigate to a path (update browsePath, reload tree)
async function navigateTo(path) {
  if (!currentProject.value) return
  currentProject.value.browsePath = path
  currentProject.value.selectedBrowseFiles = []
  await loadDirectory(path)
}

function goUp() {
  const path = currentProject.value?.browsePath || '/'
  if (path === '/') return
  const parts = path.split('/').filter(Boolean)
  parts.pop()
  navigateTo(parts.length === 0 ? '/' : '/' + parts.join('/'))
}

function goHome() {
  navigateTo('/')
}

// Lazy-load children when expanding a directory
async function onExpandedChange(newExpanded) {
  const addedKeys = newExpanded.filter(k => !expanded.value.includes(k))
  expanded.value = newExpanded

  for (const key of addedKeys) {
    const node = findNode(treeItems.value, key)
    if (node && node.isDir && !node._loaded) {
      const entries = await projectStore.browseDirectory(node.sourcePath)
      node.children = entries.map(e => ({
        ...e,
        _key: e.sourcePath,
        children: e.isDir ? undefined : undefined,
        _loaded: false,
      }))
      node._loaded = true
    }
  }
}

function findNode(items, key) {
  for (const item of items) {
    if (item._key === key) return item
    if (item.children) {
      const found = findNode(item.children, key)
      if (found) return found
    }
  }
  return null
}

// Selection
function onSelectionChange(val) {
  if (!currentProject.value) return
  const items = Array.isArray(val) ? val : [val]
  selected.value = items
  currentProject.value.selectedBrowseFiles = items.map(i => i.sourcePath || i._key)
}

// Double-click directory → navigate into it
function onItemDblClick(item) {
  if (item.isDir) {
    navigateTo(item.sourcePath)
  }
}

// Add selected files to project
async function addSelectedToProject() {
  if (!currentProject.value) return
  const paths = currentProject.value.selectedBrowseFiles
  if (paths.length > 0) {
    await projectStore.addFiles(tabId.value, paths)
    currentProject.value.selectedBrowseFiles = []
    selected.value = []
  }
}

function formatBytes(bytes) {
  return projectStore.formatBytes(bytes)
}

function getKey(item) {
  return item._key || item.sourcePath
}

function getChildren(item) {
  if (!item.isDir) return undefined
  return item.children
}

// Init: load initial browse entries
onMounted(async () => {
  if (currentProject.value) {
    await loadDirectory(currentProject.value.browsePath || '/')
  }
})

// Watch for tab changes
watch(tabId, async () => {
  if (currentProject.value) {
    await loadDirectory(currentProject.value.browsePath || '/')
  }
})
</script>

<template>
  <div class="flex flex-col h-full">
    <!-- Header: breadcrumb navigation -->
    <div class="flex items-center gap-1 px-3 py-2 bg-gray-100 dark:bg-gray-800 border-b border-gray-300 dark:border-gray-700">
      <button
        @click="goHome"
        class="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors shrink-0"
        :title="'Home'"
      >
        <Home :size="14" class="text-gray-600 dark:text-gray-400" />
      </button>
      <button
        @click="goUp"
        class="p-1 rounded hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors shrink-0"
        :title="'Up'"
      >
        <ArrowUp :size="14" class="text-gray-600 dark:text-gray-400" />
      </button>

      <div class="flex items-center gap-0.5 text-xs text-gray-600 dark:text-gray-400 overflow-hidden flex-1 ml-1">
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
      </div>
    </div>

    <!-- Tree content -->
    <div class="flex-1 overflow-y-auto">
      <div v-if="treeItems.length === 0" class="text-center text-gray-500 text-sm py-8">
        {{ t('project.emptyDirectory') }}
      </div>

      <TreeRoot
        v-else
        :items="treeItems"
        :get-key="getKey"
        :get-children="getChildren"
        :expanded="expanded"
        :model-value="selected"
        multiple
        selection-behavior="toggle"
        class="w-full text-sm"
        @update:expanded="onExpandedChange"
        @update:model-value="onSelectionChange"
      >
        <template #default="{ flattenItems }">
          <TreeItem
            v-for="item in flattenItems"
            :key="item._id"
            v-bind="item.bind"
            class="flex items-center gap-1.5 py-1 cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-800/50 transition-colors outline-none data-[selected]:bg-blue-900/20"
            :style="{ paddingLeft: (item.level * 16 + 8) + 'px', paddingRight: '8px' }"
            @dblclick="onItemDblClick(item.value)"
          >
            <template #default="{ isExpanded, isSelected }">
              <!-- Expand chevron -->
              <span class="w-4 h-4 flex items-center justify-center shrink-0">
                <ChevronRight
                  v-if="item.value.isDir"
                  :size="14"
                  class="text-gray-500 transition-transform duration-150"
                  :class="{ 'rotate-90': isExpanded }"
                />
              </span>

              <!-- File/folder icon -->
              <FileIcon
                :name="item.value.name"
                :is-dir="item.value.isDir"
                :is-open="isExpanded"
                :size="16"
              />

              <!-- Name -->
              <span class="truncate flex-1 text-gray-800 dark:text-gray-200">
                {{ item.value.name }}
              </span>

              <!-- Size -->
              <span v-if="!item.value.isDir" class="text-xs text-gray-500 shrink-0 ml-2">
                {{ formatBytes(item.value.size) }}
              </span>
            </template>
          </TreeItem>
        </template>
      </TreeRoot>
    </div>

    <!-- Toolbar -->
    <div class="flex items-center gap-2 px-3 py-2 bg-gray-100 dark:bg-gray-800 border-t border-gray-300 dark:border-gray-700">
      <button
        @click="addSelectedToProject"
        :disabled="!currentProject || currentProject.selectedBrowseFiles.length === 0"
        class="px-3 py-1 text-xs font-medium rounded bg-blue-600 hover:bg-blue-500 text-white disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
      >
        {{ t('project.addToProject') }}
      </button>
      <span class="text-xs text-gray-500">
        {{ currentProject?.selectedBrowseFiles?.length || 0 }} {{ t('project.selected') }}
      </span>
    </div>
  </div>
</template>
