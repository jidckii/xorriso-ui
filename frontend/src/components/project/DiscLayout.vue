<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { TreeRoot, TreeItem } from 'reka-ui'
import { ChevronRight } from 'lucide-vue-next'
import FileIcon from '../ui/FileIcon.vue'
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
const selected = ref([])

function onSelectionChange(val) {
  if (!currentProject.value) return
  const items = Array.isArray(val) ? val : [val]
  selected.value = items
  currentProject.value.selectedProjectEntries = items.map(i => i.destPath || i._key)
}

async function removeSelectedFromProject() {
  if (!currentProject.value) return
  for (const destPath of currentProject.value.selectedProjectEntries) {
    await projectStore.removeEntry(tabId.value, destPath)
  }
  currentProject.value.selectedProjectEntries = []
  selected.value = []
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
</script>

<template>
  <div class="flex flex-col h-full">
    <!-- Header -->
    <div class="flex items-center gap-2 px-3 py-1.5 bg-gray-100 dark:bg-gray-800 border-b border-gray-300 dark:border-gray-700 min-h-[34px]">
      <svg class="w-4 h-4 text-gray-600 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <circle cx="12" cy="12" r="10" stroke-width="1.5" />
        <circle cx="12" cy="12" r="3" stroke-width="1.5" />
      </svg>
      <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('project.discLayout') }}</span>
      <span class="ml-auto text-xs text-gray-500">
        {{ currentProject?.entries?.length || 0 }} {{ t('project.items') }} — {{ formatBytes(currentProject?.totalSize || 0) }}
      </span>
    </div>

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
        :model-value="selected"
        multiple
        selection-behavior="toggle"
        class="w-full text-sm"
        @update:model-value="onSelectionChange"
      >
        <template #default="{ flattenItems }">
          <TreeItem
            v-for="item in flattenItems"
            :key="item._id"
            v-bind="item.bind"
            class="flex items-center gap-1.5 py-1 cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-800/50 transition-colors outline-none data-[selected]:bg-blue-900/20"
            :style="{ paddingLeft: (item.level * 16 + 8) + 'px', paddingRight: '8px' }"
          >
            <template #default="{ isExpanded, isSelected }">
              <!-- Expand chevron -->
              <span class="w-4 h-4 flex items-center justify-center shrink-0">
                <ChevronRight
                  v-if="item.value.children?.length"
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
              <span v-if="!item.value.isDir && item.value.size" class="text-xs text-gray-500 shrink-0 ml-2">
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
        @click="removeSelectedFromProject"
        :disabled="!currentProject || currentProject.selectedProjectEntries.length === 0"
        class="px-3 py-1 text-xs font-medium rounded bg-red-600 hover:bg-red-500 text-white disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
      >
        {{ t('project.remove') }}
      </button>
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
