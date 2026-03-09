import BurnOverlay from './BurnOverlay.vue'
import { useTabStore } from '../../stores/tabStore'
import { useDeviceStore } from '../../stores/deviceStore'
import { useBurnStore } from '../../stores/burnStore'

const mockProjectData = {
  version: 1,
  name: 'My Backup',
  filePath: '',
  volumeId: 'MY_BACKUP',
  totalSize: 55574528,
  entries: [
    { sourcePath: '/home/user/photo.jpg', destPath: '/photo.jpg', name: 'photo.jpg', isDir: false, size: 3145728 },
    { sourcePath: '/home/user/docs', destPath: '/docs', name: 'docs', isDir: true, size: 0 },
    { sourcePath: '/home/user/backup.tar.gz', destPath: '/backup.tar.gz', name: 'backup.tar.gz', isDir: false, size: 52428800 },
  ],
  isoOptions: { isoLevel: '3', rockRidge: true, joliet: true, udf: false, hfsPlus: false, zisofs: false, md5: false, backupMode: false },
  burnOptions: { speed: 'auto', verify: true, eject: false, closeDisc: true, burnMode: 'auto', dummyMode: false, streamRecording: false, multisession: false, padding: 300 },
  browsePath: '/',
  browseEntries: [],
  selectedBrowseFiles: [],
  selectedProjectEntries: [],
  modified: false,
  browseTreeExpanded: [],
  discTreeExpanded: [],
  browseShowHidden: false,
  createdAt: null,
  updatedAt: null,
}

const mockDevices = [
  { path: '/dev/sr0', vendor: 'ASUS', model: 'BW-16D1HT' },
]

const mockMediaInfo = {
  mediaType: 'DVD+R',
  mediaStatus: 'blank',
  totalCapacity: 4700000000,
  freeSpace: 4700000000,
  erasable: false,
  sessions: 0,
}

const mockSpeeds = [
  { writeSpeed: '4x', displayName: '4x (5,540 KB/s)' },
  { writeSpeed: '8x', displayName: '8x (11,080 KB/s)' },
]

function storeDecorator(overrides = {}) {
  return (story) => ({
    components: { story },
    setup() {
      const tabStore = useTabStore()
      tabStore.tabs = [{
        id: 'tab1',
        label: 'My Backup',
        projectData: { ...mockProjectData },
      }]
      tabStore.activeTabId = 'tab1'

      const deviceStore = useDeviceStore()
      deviceStore.devices = overrides.devices ?? mockDevices
      deviceStore.currentDevicePath = overrides.currentDevicePath ?? '/dev/sr0'
      deviceStore.mediaInfo = overrides.mediaInfo ?? mockMediaInfo
      deviceStore.speeds = overrides.speeds ?? mockSpeeds

      const burnStore = useBurnStore()
      burnStore.init = () => {}

      if (overrides.viewMode) {
        localStorage.setItem('xorriso-burn-mode', overrides.viewMode)
      }

      return {}
    },
    template: '<div style="width: 600px;" class="bg-gray-900 rounded-lg border border-gray-700"><story /></div>',
  })
}

export default {
  title: 'Burn/BurnOverlay',
  component: BurnOverlay,
  parameters: { layout: 'centered' },
}

export const SimpleModeBurn = {
  args: { mode: 'burn' },
  decorators: [storeDecorator({ viewMode: 'simple' })],
}

export const ExpertModeBurn = {
  args: { mode: 'burn' },
  decorators: [storeDecorator({ viewMode: 'expert' })],
}

export const SaveMode = {
  args: { mode: 'save' },
  decorators: [storeDecorator({ viewMode: 'simple' })],
}
