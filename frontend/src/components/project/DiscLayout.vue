<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { TreeRoot, TreeItem } from 'reka-ui'
import { ChevronRight } from 'lucide-vue-next'
import FileIcon from '../ui/FileIcon.vue'
import PanelHeader from '../ui/PanelHeader.vue'
import ImagePreviewTooltip from '../ui/ImagePreviewTooltip.vue'
import { useProjectStore } from '../../stores/projectStore'
import { useTabStore } from '../../stores/tabStore'

const { t } = useI18n()
const router = useRouter()
const projectStore = useProjectStore()
const tabStore = useTabStore()

const currentProject = computed(() => tabStore.activeProject)
const tabId = computed(() => tabStore.activeTabId)

// Build tree from flat entries
const treeItems = computed(() => {
  const entries = currentProject.value?.entries || []
  if (entries.length === 0) return []
  return buildTreeFromEntries(entries)
})

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

// Manual selection with parent-child propagation
const selectedKeys = ref(new Set())

function isItemSelected(key) {
  return selectedKeys.value.has(key)
}

function toggleItemSelection(item) {
  const selecting = !selectedKeys.value.has(item._key)

  if (selecting) {
    selectedKeys.value.add(item._key)
  } else {
    selectedKeys.value.delete(item._key)
  }

  // Propagate to children
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
  for (const destPath of currentProject.value.selectedProjectEntries) {
    await projectStore.removeEntry(tabId.value, destPath)
  }
  currentProject.value.selectedProjectEntries = []
  selectedKeys.value = new Set()
}

// Select all / deselect all
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

function formatBytes(bytes) {
  return projectStore.formatBytes(bytes)
}

function getKey(item) {
  return item._key || item.destPath
}

function getChildren(item) {
  if (!item.children || item.children.length === 0) return undefined
  return item.children
}

// Image preview tooltip
const previewVisible = ref(false)
const previewX = ref(0)
const previewY = ref(0)
const previewSourcePath = ref('')

const imageExts = new Set(['.jpg', '.jpeg', '.png', '.gif', '.webp', '.bmp'])

function isImageFile(name) {
  if (!name) return false
  const dot = name.lastIndexOf('.')
  if (dot < 0) return false
  return imageExts.has(name.substring(dot).toLowerCase())
}

function onItemMouseEnter(e, item) {
  if (item.isDir || !item.sourcePath || !isImageFile(item.name)) return
  previewSourcePath.value = item.sourcePath
  previewX.value = e.clientX
  previewY.value = e.clientY
  previewVisible.value = true
}

function onItemMouseMove(e) {
  previewX.value = e.clientX
  previewY.value = e.clientY
}

function onItemMouseLeave() {
  previewVisible.value = false
}
</script>

<template>
  <div class="flex flex-col h-full">
    <PanelHeader>
      <template #row1>
        <svg class="w-4 h-4 text-gray-600 dark:text-gray-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="10" stroke-width="1.5" />
          <circle cx="12" cy="12" r="3" stroke-width="1.5" />
        </svg>
        <span class="text-xs font-medium text-gray-700 dark:text-gray-300">{{ t('project.discLayout') }}</span>
        <span class="flex-1" />
      </template>
      <template #row2>
        <span class="text-xs text-gray-500">
          {{ currentProject?.entries?.length || 0 }} {{ t('project.items') }} — {{ formatBytes(currentProject?.totalSize || 0) }}
        </span>
      </template>
    </PanelHeader>

    <!-- Tree content -->
    <div class="flex-1 overflow-y-auto">
      <!-- Empty state -->
      <div
        v-if="!currentProject?.entries?.length"
        class="flex flex-col items-center justify-center h-full text-gray-500 py-12"
      >
        <svg class="w-12 h-12 mb-3 text-gray-200 dark:text-gray-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
            d="M12 9v3m0 0v3m0-3h3m-3 0H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-sm">{{ t('common.addFilesFromBrowser') }}</p>
      </div>

      <!-- File tree -->
      <TreeRoot
        v-else
        :items="treeItems"
        :get-key="getKey"
        :get-children="getChildren"
        v-model:expanded="expanded"
        class="w-full text-sm select-none"
      >
        <template #default="{ flattenItems }">
          <TreeItem
            v-for="item in flattenItems"
            :key="item._id"
            v-bind="item.bind"
            class="flex items-center gap-1.5 py-1 cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-800/50 transition-colors outline-none"
            :class="{ 'bg-blue-500/15': isItemSelected(item.value._key) }"
            :style="{ paddingLeft: (item.level * 16 + 8) + 'px', paddingRight: '8px' }"
            @click="toggleItemSelection(item.value)"
            @mouseenter="onItemMouseEnter($event, item.value)"
            @mousemove="onItemMouseMove"
            @mouseleave="onItemMouseLeave"
          >
            <template #default="{ isExpanded }">
              <!-- Expand chevron -->
              <span class="w-4 h-4 flex items-center justify-center shrink-0">
                <ChevronRight
                  v-if="item.value.children?.length"
                  :size="14"
                  class="text-gray-500 transition-transform duration-150"
                  :class="{ 'rotate-90': isExpanded }"
                />
              </span>

              <!-- Selection checkbox -->
              <input
                type="checkbox"
                :checked="isItemSelected(item.value._key)"
                @click.stop="toggleItemSelection(item.value)"
                class="w-3.5 h-3.5 accent-blue-600 shrink-0 cursor-pointer"
              />

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
              <span v-if="item.value.size" class="text-xs text-gray-500 shrink-0 ml-2">
                {{ formatBytes(item.value.size) }}
              </span>
            </template>
          </TreeItem>
        </template>
      </TreeRoot>

      <!-- Image preview tooltip -->
      <ImagePreviewTooltip
        :file-path="previewSourcePath"
        :visible="previewVisible"
        :x="previewX"
        :y="previewY"
      />
    </div>

    <!-- Toolbar -->
    <div class="flex items-center gap-2 px-3 py-2 bg-gray-100 dark:bg-gray-800 border-t border-gray-300 dark:border-gray-700">
      <label class="flex items-center gap-1.5 text-xs text-gray-500 cursor-pointer select-none">
        <input
          type="checkbox"
          :checked="allSelected"
          @change="allSelected ? deselectAll() : selectAll()"
          class="w-3.5 h-3.5 accent-blue-600 cursor-pointer"
        />
        {{ t('project.selectAll') }}
      </label>
      <button
        @click="removeSelectedFromProject"
        :disabled="selectedCount === 0"
        class="px-3 py-1 text-xs font-medium rounded bg-red-600 hover:bg-red-500 text-white disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
      >
        {{ t('project.remove') }}
      </button>
      <span class="text-xs text-gray-500">
        {{ selectedCount }} {{ t('project.selected') }}
      </span>
      <span class="flex-1"></span>
      <button
        @click="goToBurn"
        :disabled="!currentProject || currentProject.entries.length === 0"
        class="px-4 py-1 text-xs font-medium rounded bg-orange-600 hover:bg-orange-500 text-white disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
      >
        {{ t('burn.title') }}
      </button>
    </div>
  </div>
</template>
