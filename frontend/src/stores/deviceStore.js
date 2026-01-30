import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { ListDevices, GetMediaInfo, GetSpeeds, EjectDisc } from '../../bindings/xorriso-ui/services/deviceservice.js'
import { Events } from '@wailsio/runtime'

export const useDeviceStore = defineStore('device', () => {
  // --- State ---
  const devices = ref([])
  const currentDevicePath = ref(null)
  const mediaInfo = ref(null)
  const speeds = ref([])
  const loading = ref(false)

  // --- Getters ---
  const currentDevice = computed(() =>
    devices.value.find(d => d.path === currentDevicePath.value) || null
  )

  const hasMedia = computed(() =>
    mediaInfo.value !== null &&
    mediaInfo.value.mediaType !== '' &&
    mediaInfo.value.mediaStatus !== ''
  )

  const mediaCapacityBytes = computed(() => mediaInfo.value?.totalCapacity ?? 0)

  const mediaFreeBytes = computed(() => mediaInfo.value?.freeSpace ?? 0)

  // --- Actions ---

  async function fetchDevices() {
    loading.value = true
    try {
      const result = await ListDevices()
      devices.value = result

      // Auto-select the first device if none selected
      if (!currentDevicePath.value && devices.value.length > 0) {
        currentDevicePath.value = devices.value[0].path
        await fetchMediaInfo()
      }
    } catch (error) {
      console.error('Failed to fetch devices:', error)
    } finally {
      loading.value = false
    }
  }

  async function selectDevice(path) {
    currentDevicePath.value = path
    mediaInfo.value = null
    speeds.value = []
    await fetchMediaInfo()
    await fetchSpeeds()
  }

  async function fetchMediaInfo() {
    if (!currentDevicePath.value) return
    loading.value = true
    try {
      const result = await GetMediaInfo(currentDevicePath.value)
      mediaInfo.value = result
    } catch (error) {
      console.error('Failed to fetch media info:', error)
      mediaInfo.value = null
    } finally {
      loading.value = false
    }
  }

  async function fetchSpeeds() {
    if (!currentDevicePath.value) return
    try {
      const result = await GetSpeeds(currentDevicePath.value)
      speeds.value = result
    } catch (error) {
      console.error('Failed to fetch speeds:', error)
      speeds.value = []
    }
  }

  async function ejectDisc() {
    if (!currentDevicePath.value) return
    try {
      await EjectDisc(currentDevicePath.value)
      mediaInfo.value = null
    } catch (error) {
      console.error('Failed to eject disc:', error)
    }
  }

  /**
   * Initialize event listeners for device-related Wails events.
   * Call this once from the app root or on store creation.
   */
  function init() {
    Events.On('device:list-updated', (eventData) => {
      if (eventData?.data) {
        devices.value = eventData.data
      } else {
        fetchDevices()
      }
    })

    Events.On('device:media-changed', (eventData) => {
      const data = eventData?.data
      if (data?.devicePath === currentDevicePath.value) {
        fetchMediaInfo()
      }
    })

    // Initial fetch
    fetchDevices()
  }

  return {
    // State
    devices,
    currentDevicePath,
    mediaInfo,
    speeds,
    loading,
    // Getters
    currentDevice,
    hasMedia,
    mediaCapacityBytes,
    mediaFreeBytes,
    // Actions
    fetchDevices,
    selectDevice,
    fetchMediaInfo,
    fetchSpeeds,
    ejectDisc,
    init,
  }
})
