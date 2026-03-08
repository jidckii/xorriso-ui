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
  RemoveEntries,
  CalculateSize,
  GetHomeDirectory,
  ListMountPoints,
  GetImagePreview,
  OpenWithDefault,
  RevealInFileManager,
  GetFileProperties,
} from '../../bindings/xorriso-ui/services/projectservice.js'
import { useTabStore } from './tabStore'
import { formatBytes } from '../composables/useFormatBytes'

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

  // Пакетное удаление нескольких записей за один IPC-вызов
  async function removeEntries(tabId, destPaths) {
    if (!destPaths || destPaths.length === 0) return
    const tabStore = useTabStore()
    const data = tabStore.getProjectData(tabId)
    if (!data) return

    try {
      const result = await RemoveEntries({ ...data }, destPaths)
      tabStore.updateProjectData(tabId, { entries: result.entries, modified: true })
      await calculateSize(tabId)
    } catch (error) {
      console.error('Failed to remove entries:', error)
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

  async function getImagePreview(filePath, maxSize = 200) {
    try {
      return await GetImagePreview(filePath, maxSize)
    } catch {
      return ''
    }
  }

  async function openWithDefault(filePath) {
    try {
      await OpenWithDefault(filePath)
    } catch (error) {
      console.error('Failed to open file:', error)
    }
  }

  async function revealInFileManager(filePath) {
    try {
      await RevealInFileManager(filePath)
    } catch (error) {
      console.error('Failed to reveal in file manager:', error)
    }
  }

  async function getFileProperties(filePath) {
    try {
      return await GetFileProperties(filePath)
    } catch (error) {
      console.error('Failed to get file properties:', error)
      return null
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
    removeEntries,
    calculateSize,
    browseDirectory,
    getHomeDirectory,
    listMountPoints,
    getImagePreview,
    openWithDefault,
    revealInFileManager,
    getFileProperties,
    formatBytes,
  }
})
