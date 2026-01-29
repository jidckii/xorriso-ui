import { defineStore } from 'pinia'
import { ref, reactive, computed } from 'vue'
// TODO: import Wails service bindings when available
// import { NewProject, AddFiles, RemoveEntry, CalculateSize, BrowseDirectory } from '../../bindings/xorriso-ui/ProjectService'

export const useProjectStore = defineStore('project', () => {
  // --- State ---
  const project = reactive({
    name: 'Untitled Project',
    volumeId: 'UNTITLED',
    entries: [],
    isoOptions: {
      rockRidge: true,
      joliet: true,
      md5: true,
      backupMode: false,
    },
    burnOptions: {
      speed: 'auto',
      dummyMode: false,
      verify: true,
      closeDisc: false,
      eject: true,
      burnMode: 'auto',
      streamRecording: false,
    },
  })

  const totalSize = ref(0)
  const modified = ref(false)
  const browseLoading = ref(false)

  // --- Getters ---
  const entryCount = computed(() => project.entries.length)

  const totalSizeFormatted = computed(() => formatBytes(totalSize.value))

  // --- Actions ---

  async function newProject(name = 'Untitled Project', volumeId = 'UNTITLED') {
    // TODO: replace with actual Wails call
    // const result = await NewProject(name, volumeId)

    project.name = name
    project.volumeId = volumeId
    project.entries = []
    project.isoOptions = {
      rockRidge: true,
      joliet: true,
      md5: true,
      backupMode: false,
    }
    project.burnOptions = {
      speed: 'auto',
      dummyMode: false,
      verify: true,
      closeDisc: false,
      eject: true,
      burnMode: 'auto',
      streamRecording: false,
    }
    totalSize.value = 0
    modified.value = false
  }

  async function addFiles(sourcePaths, destDir = '/') {
    if (!sourcePaths || sourcePaths.length === 0) return

    try {
      // TODO: replace with actual Wails call
      // const result = await AddFiles(project, sourcePaths, destDir)
      // project.entries = result.entries

      // Mock: add entries to the project
      for (const path of sourcePaths) {
        const name = path.split('/').pop()
        project.entries.push({
          name,
          sourcePath: path,
          destPath: destDir === '/' ? `/${name}` : `${destDir}/${name}`,
          isDir: false,
          size: Math.floor(Math.random() * 1024 * 1024 * 100), // Mock random size
        })
      }

      modified.value = true
      await calculateSize()
    } catch (error) {
      console.error('Failed to add files:', error)
    }
  }

  async function removeEntry(destPath) {
    try {
      // TODO: replace with actual Wails call
      // const result = await RemoveEntry(project, destPath)
      // project.entries = result.entries

      const idx = project.entries.findIndex(e => e.destPath === destPath)
      if (idx !== -1) {
        project.entries.splice(idx, 1)
        modified.value = true
        await calculateSize()
      }
    } catch (error) {
      console.error('Failed to remove entry:', error)
    }
  }

  async function calculateSize() {
    try {
      // TODO: replace with actual Wails call
      // const result = await CalculateSize(project)
      // totalSize.value = result

      // Mock: sum up entry sizes
      totalSize.value = project.entries.reduce((sum, e) => sum + (e.size || 0), 0)
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
      // TODO: replace with actual Wails call
      // const result = await BrowseDirectory(path)
      // return result

      // Mock data for UI development
      return [
        { name: '..', path: path === '/' ? '/' : path.split('/').slice(0, -1).join('/') || '/', isDir: true, size: 0 },
        { name: 'Documents', path: `${path === '/' ? '' : path}/Documents`, isDir: true, size: 0 },
        { name: 'Downloads', path: `${path === '/' ? '' : path}/Downloads`, isDir: true, size: 0 },
        { name: 'Music', path: `${path === '/' ? '' : path}/Music`, isDir: true, size: 0 },
        { name: 'Photos', path: `${path === '/' ? '' : path}/Photos`, isDir: true, size: 0 },
        { name: 'readme.txt', path: `${path === '/' ? '' : path}/readme.txt`, isDir: false, size: 1234 },
        { name: 'backup.tar.gz', path: `${path === '/' ? '' : path}/backup.tar.gz`, isDir: false, size: 52428800 },
        { name: 'image.iso', path: `${path === '/' ? '' : path}/image.iso`, isDir: false, size: 734003200 },
      ]
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
    addFiles,
    removeEntry,
    calculateSize,
    browseDirectory,
    // Helpers
    formatBytes,
  }
})
