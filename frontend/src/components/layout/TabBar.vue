<script setup>
import { ref, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useTabStore } from '../../stores/tabStore'
import { useProjectStore } from '../../stores/projectStore'

const { t } = useI18n()
const tabStore = useTabStore()
const projectStore = useProjectStore()

const editingTabId = ref(null)
const editValue = ref('')
const editInput = ref(null)

function closeTab(tabId, event) {
  event.stopPropagation()
  tabStore.removeTab(tabId)
}

async function newTab() {
  const tabId = tabStore.addProjectTab()
  const tab = tabStore.tabs.find(t => t.id === tabId)
  await projectStore.newProject(tabId, tab.label, tab.label)
}

function startRename(tab) {
  editingTabId.value = tab.id
  editValue.value = tab.label
  nextTick(() => {
    editInput.value?.select()
  })
}

function commitRename(tabId) {
  const name = editValue.value.trim()
  if (name && editingTabId.value === tabId) {
    tabStore.updateTabLabel(tabId, name)
    tabStore.updateProjectData(tabId, { name, modified: true })
  }
  editingTabId.value = null
}

function cancelRename() {
  editingTabId.value = null
}
</script>

<template>
  <div class="flex items-end bg-gray-100 dark:bg-gray-800 border-b border-gray-300 dark:border-gray-700 px-1 overflow-x-auto">
    <div
      v-for="tab in tabStore.tabs"
      :key="tab.id"
      @click="tabStore.setActiveTab(tab.id)"
      class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium cursor-pointer border-b-2 transition-colors whitespace-nowrap shrink-0 select-none"
      :class="tab.id === tabStore.activeTabId
        ? 'border-blue-500 text-gray-900 dark:text-gray-100 bg-white dark:bg-gray-900'
        : 'border-transparent text-gray-500 hover:text-gray-700 dark:hover:text-gray-300 hover:bg-gray-200/50 dark:hover:bg-gray-700/50'"
    >
      <!-- File icon -->
      <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
          d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>

      <input
        v-if="editingTabId === tab.id"
        ref="editInput"
        v-model="editValue"
        @blur="commitRename(tab.id)"
        @keydown.enter="commitRename(tab.id)"
        @keydown.escape="cancelRename"
        @click.stop
        class="w-24 bg-transparent border border-blue-500 rounded px-1 py-0 text-xs outline-none text-gray-900 dark:text-gray-100"
      />
      <span v-else @dblclick.stop="startRename(tab)" :title="t('tabs.renameTab')">{{ tab.label }}</span>

      <!-- Modified indicator -->
      <span
        v-if="tab.projectData?.modified"
        class="w-1.5 h-1.5 rounded-full bg-orange-400 shrink-0"
      />

      <!-- Close button -->
      <button
        @click="closeTab(tab.id, $event)"
        class="ml-0.5 p-0.5 rounded hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
        :title="t('tabs.closeTab')"
      >
        <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Add new project tab -->
    <button
      @click="newTab"
      class="p-1.5 ml-1 mb-0.5 rounded hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors shrink-0"
      :title="t('tabs.newProject')"
    >
      <svg class="w-3.5 h-3.5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
    </button>
  </div>
</template>
