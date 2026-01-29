import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
// TODO: import Wails service bindings when available
// import { ListDevices, GetMediaInfo, GetSpeeds, EjectDisc } from '../../bindings/xorriso-ui/DeviceService'
// TODO: import Wails events when backend is ready
// import { Events } from '@wailsio/runtime'

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

  const hasMedia = computed(() => mediaInfo.value !== null && mediaInfo.value.mediaType !== 'none')

  const mediaCapacityBytes = computed(() => mediaInfo.value?.capacityBytes ?? 0)

  const mediaFreeBytes = computed(() => mediaInfo.value?.freeBytes ?? 0)

  // --- Actions ---

  async function fetchDevices() {
    loading.value = true
    try {
      // TODO: replace with actual Wails call
      // const result = await ListDevices()
      // devices.value = result

      // Mock data for UI development
      devices.value = [
        {
          path: '/dev/sr0',
          name: 'ASUS BW-16D1HT',
          vendor: 'ASUS',
          model: 'BW-16D1HT',
          capabilities: ['cd-r', 'cd-rw', 'dvd-r', 'dvd-rw', 'dvd+r', 'dvd+rw', 'bd-r', 'bd-re'],
        },
        {
          path: '/dev/sr1',
          name: 'TSSTcorp CDDVDW SH-224DB',
          vendor: 'TSSTcorp',
          model: 'CDDVDW SH-224DB',
          capabilities: ['cd-r', 'cd-rw', 'dvd-r', 'dvd-rw', 'dvd+r', 'dvd+rw'],
        },
      ]

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
      // TODO: replace with actual Wails call
      // const result = await GetMediaInfo(currentDevicePath.value)
      // mediaInfo.value = result

      // Mock data for UI development
      mediaInfo.value = {
        mediaType: 'BD-R',
        mediaStatus: 'blank',
        mediaLabel: '',
        capacityBytes: 25025314816, // ~25 GB (BD-R SL)
        usedBytes: 0,
        freeBytes: 25025314816,
        sessions: 0,
        tracks: 0,
      }
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
      // TODO: replace with actual Wails call
      // const result = await GetSpeeds(currentDevicePath.value)
      // speeds.value = result

      // Mock data for UI development
      speeds.value = [
        { label: 'Auto', value: 'auto' },
        { label: '2x (9.0 MB/s)', value: '2' },
        { label: '4x (18.0 MB/s)', value: '4' },
        { label: '6x (27.0 MB/s)', value: '6' },
        { label: '8x (36.0 MB/s)', value: '8' },
        { label: '12x (54.0 MB/s)', value: '12' },
      ]
    } catch (error) {
      console.error('Failed to fetch speeds:', error)
      speeds.value = []
    }
  }

  async function ejectDisc() {
    if (!currentDevicePath.value) return
    try {
      // TODO: replace with actual Wails call
      // await EjectDisc(currentDevicePath.value)
      console.log('Eject disc:', currentDevicePath.value)
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
    // TODO: replace with actual Wails event subscriptions
    // Events.On('device:list-updated', () => { fetchDevices() })
    // Events.On('device:media-changed', (data) => {
    //   if (data.devicePath === currentDevicePath.value) {
    //     fetchMediaInfo()
    //   }
    // })

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
