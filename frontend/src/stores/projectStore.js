import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  NewProject,
  OpenProject,
  SaveProject,
  SaveProjectAs,
  BrowseDirectory,
  AddFiles,
  RemoveEntry,
  CalculateSize,
  GetHomeDirectory,
  ListMountPoints,
} from '../../bindings/xorriso-ui/services/projectservice.js'
import { useTabStore } from './tabStore'

export const useProjectStore = defineStore('project', () => {
  const browseLoading = ref(false)

  async function newProject(tabId, name = 'Untitled Project', volumeId = 'UNTITLED') {
    const tabStore = useTabStore()
    const result = await NewProject(name, volumeId)
    tabStore.updateProjectData(tabId, { ...result, totalSize: 0, modified: false })
    tabStore.updateTabLabel(tabId, name)
  }

  async function openProject(tabId, filePath) {
    const tabStore = useTabStore()
    const result = await OpenProject(filePath)
    tabStore.updateProjectData(tabId, { ...result, modified: false })
    tabStore.updateTabLabel(tabId, result.name || filePath)
    await calculateSize(tabId)
  }

  async function saveProject(tabId) {
    const tabStore = useTabStore()
    const data = tabStore.getProjectData(tabId)
    if (!data) return
    await SaveProject({ ...data })
    tabStore.updateProjectData(tabId, { modified: false })
  }

  async function saveProjectAs(tabId, filePath) {
    const tabStore = useTabStore()
    const data = tabStore.getProjectData(tabId)
    if (!data) return
    await SaveProjectAs({ ...data }, filePath)
    tabStore.updateProjectData(tabId, { filePath, modified: false })
  }

  async function addFiles(tabId, sourcePaths, destDir = '/') {
    if (!sourcePaths || sourcePaths.length === 0) return
    const tabStore = useTabStore()
    const data = tabStore.getProjectData(tabId)
    if (!data) return

    try {
      const result = await AddFiles({ ...data }, sourcePaths, destDir)
      tabStore.updateProjectData(tabId, { entries: result.entries, modified: true })
      await calculateSize(tabId)
    } catch (error) {
      console.error('Failed to add files:', error)
    }
  }

  async function removeEntry(tabId, destPath) {
    const tabStore = useTabStore()
    const data = tabStore.getProjectData(tabId)
    if (!data) return

    try {
      const result = await RemoveEntry({ ...data }, destPath)
      tabStore.updateProjectData(tabId, { entries: result.entries, modified: true })
      await calculateSize(tabId)
    } catch (error) {
      console.error('Failed to remove entry:', error)
    }
  }

  async function calculateSize(tabId) {
    const tabStore = useTabStore()
    const data = tabStore.getProjectData(tabId)
    if (!data) return

    try {
      const result = await CalculateSize({ ...data })
      tabStore.updateProjectData(tabId, { totalSize: result })
    } catch (error) {
      console.error('Failed to calculate size:', error)
    }
  }

  async function browseDirectory(path = '/') {
    browseLoading.value = true
    try {
      const result = await BrowseDirectory(path)
      return result
    } catch (error) {
      console.error('Failed to browse directory:', error)
      return []
    } finally {
      browseLoading.value = false
    }
  }

  function formatBytes(bytes) {
    if (bytes === 0) return '0 B'
    const units = ['B', 'KiB', 'MiB', 'GiB', 'TiB']
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${units[i]}`
  }

  async function getHomeDirectory() {
    try {
      return await GetHomeDirectory()
    } catch {
      return '/'
    }
  }

  async function listMountPoints() {
    try {
      return await ListMountPoints()
    } catch {
      return []
    }
  }

  return {
    browseLoading,
    newProject,
    openProject,
    saveProject,
    saveProjectAs,
    addFiles,
    removeEntry,
    calculateSize,
    browseDirectory,
    getHomeDirectory,
    listMountPoints,
    formatBytes,
  }
})
