import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

let nextTabId = 1

function createProjectData(name = 'Untitled Project', volumeId = 'UNTITLED') {
  return {
    name,
    filePath: '',
    volumeId,
    entries: [],
    isoOptions: {
      udf: true,
      rockRidge: false,
      joliet: false,
      hfsPlus: false,
      zisofs: false,
      md5: true,
      backupMode: false,
      bootImage: '',
      bootCatalog: '',
      efiBootImage: '',
      bootMode: '',
    },
    burnOptions: {
      speed: 'auto',
      dummyMode: false,
      verify: true,
      closeDisc: false,
      eject: true,
      burnMode: 'auto',
      streamRecording: false,
      padding: '',
    },
    createdAt: '',
    updatedAt: '',
    // Per-tab file browser state
    browsePath: '/',
    browseEntries: [],
    selectedBrowseFiles: [],
    // Per-tab disc layout state
    selectedProjectEntries: [],
    totalSize: 0,
    modified: false,
  }
}

export const useTabStore = defineStore('tabs', () => {
  const tabs = ref([])
  const activeTabId = ref(null)
  const showDiscInfo = ref(false)

  const activeTab = computed(() =>
    tabs.value.find(t => t.id === activeTabId.value) || null
  )

  const activeProject = computed(() =>
    activeTab.value?.projectData || null
  )

  function addProjectTab(name, volumeId) {
    const id = `tab-${nextTabId++}`
    tabs.value.push({
      id,
      label: name || 'Untitled Project',
      projectData: createProjectData(name, volumeId),
    })
    activeTabId.value = id
    return id
  }

  function removeTab(tabId) {
    const idx = tabs.value.findIndex(t => t.id === tabId)
    if (idx === -1) return

    tabs.value.splice(idx, 1)

    if (activeTabId.value === tabId) {
      const newIdx = Math.min(idx, tabs.value.length - 1)
      activeTabId.value = tabs.value[newIdx]?.id || null
    }

    // If no tabs remain, create a new one
    if (tabs.value.length === 0) {
      addProjectTab()
    }
  }

  function toggleDiscInfo() {
    showDiscInfo.value = !showDiscInfo.value
  }

  function setActiveTab(tabId) {
    activeTabId.value = tabId
  }

  function getProjectData(tabId) {
    const tab = tabs.value.find(t => t.id === tabId)
    return tab?.projectData || null
  }

  function updateProjectData(tabId, updates) {
    const tab = tabs.value.find(t => t.id === tabId)
    if (tab?.projectData) {
      Object.assign(tab.projectData, updates)
    }
  }

  function updateTabLabel(tabId, label) {
    const tab = tabs.value.find(t => t.id === tabId)
    if (tab) tab.label = label
  }

  return {
    tabs,
    activeTabId,
    showDiscInfo,
    activeTab,
    activeProject,
    addProjectTab,
    removeTab,
    toggleDiscInfo,
    setActiveTab,
    getProjectData,
    updateProjectData,
    updateTabLabel,
  }
})
