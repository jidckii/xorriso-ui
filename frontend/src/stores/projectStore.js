import { defineStore } from 'pinia'
import { ref, reactive, computed } from 'vue'
import {
  NewProject,
  OpenProject,
  SaveProject,
  SaveProjectAs,
  BrowseDirectory,
  AddFiles,
  RemoveEntry,
  CalculateSize,
} from '../../bindings/xorriso-ui/services/projectservice.js'

export const useProjectStore = defineStore('project', () => {
  // --- State ---
  const project = reactive({
    name: 'Untitled Project',
    filePath: '',
    volumeId: 'UNTITLED',
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
  })

  const totalSize = ref(0)
  const modified = ref(false)
  const browseLoading = ref(false)

  // --- Getters ---
  const entryCount = computed(() => project.entries.length)

  const totalSizeFormatted = computed(() => formatBytes(totalSize.value))

  // --- Actions ---

  async function newProject(name = 'Untitled Project', volumeId = 'UNTITLED') {
    const result = await NewProject(name, volumeId)
    Object.assign(project, result)
    totalSize.value = 0
    modified.value = false
  }

  async function openProject(filePath) {
    const result = await OpenProject(filePath)
    Object.assign(project, result)
    modified.value = false
    await calculateSize()
  }

  async function saveProject() {
    await SaveProject({ ...project })
    modified.value = false
  }

  async function saveProjectAs(filePath) {
    await SaveProjectAs({ ...project }, filePath)
    project.filePath = filePath
    modified.value = false
  }

  async function addFiles(sourcePaths, destDir = '/') {
    if (!sourcePaths || sourcePaths.length === 0) return

    try {
      const result = await AddFiles({ ...project }, sourcePaths, destDir)
      project.entries = result.entries
      modified.value = true
      await calculateSize()
    } catch (error) {
      console.error('Failed to add files:', error)
    }
  }

  async function removeEntry(destPath) {
    try {
      const result = await RemoveEntry({ ...project }, destPath)
      project.entries = result.entries
      modified.value = true
      await calculateSize()
    } catch (error) {
      console.error('Failed to remove entry:', error)
    }
  }

  async function calculateSize() {
    try {
      const result = await CalculateSize({ ...project })
      totalSize.value = result
    } catch (error) {
      console.error('Failed to calculate size:', error)
    }
  }

  /**
   * Browse host filesystem directory. Returns array of file entries.
   * Used by the file browser panel to pick files for adding to project.
   */
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

  // --- Helpers ---

  function formatBytes(bytes) {
    if (bytes === 0) return '0 B'
    const units = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${units[i]}`
  }

  return {
    // State
    project,
    totalSize,
    modified,
    browseLoading,
    // Getters
    entryCount,
    totalSizeFormatted,
    // Actions
    newProject,
    openProject,
    saveProject,
    saveProjectAs,
    addFiles,
    removeEntry,
    calculateSize,
    browseDirectory,
    // Helpers
    formatBytes,
  }
})
