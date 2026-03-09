import { reactive } from 'vue'
import BurnSimpleMode from './BurnSimpleMode.vue'

const mockProject = {
  name: 'My Backup',
  volumeId: 'MY_BACKUP',
  totalSize: 55574528,
  entries: [
    { sourcePath: '/home/user/photo.jpg', destPath: '/photo.jpg', name: 'photo.jpg', isDir: false, size: 3145728 },
    { sourcePath: '/home/user/docs', destPath: '/docs', name: 'docs', isDir: true, size: 0 },
    { sourcePath: '/home/user/backup.tar.gz', destPath: '/backup.tar.gz', name: 'backup.tar.gz', isDir: false, size: 52428800 },
  ],
  isoOptions: { isoLevel: '3', rockRidge: true, joliet: true, udf: false, hfsPlus: false, zisofs: false, md5: false, backupMode: false },
  burnOptions: { speed: 'auto', verify: true, eject: false, closeDisc: true, burnMode: 'auto', dummyMode: false, streamRecording: false, multisession: false, padding: 300 },
}

const mockDevices = [
  { path: '/dev/sr0', vendor: 'ASUS', model: 'BW-16D1HT' },
  { path: '/dev/sr1', vendor: 'Pioneer', model: 'DVR-221L' },
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

function renderWithReactive(args) {
  return {
    components: { BurnSimpleMode },
    setup() {
      const project = reactive({
        ...mockProject,
        ...args.project,
        burnOptions: { ...mockProject.burnOptions, ...(args.project?.burnOptions || {}) },
        isoOptions: { ...mockProject.isoOptions, ...(args.project?.isoOptions || {}) },
      })
      return { args: { ...args, project } }
    },
    template: '<BurnSimpleMode v-bind="args" />',
  }
}

export default {
  title: 'Burn/BurnSimpleMode',
  component: BurnSimpleMode,
  render: renderWithReactive,
  decorators: [
    (story) => ({
      components: { story },
      template: '<div style="width: 500px;"><story /></div>',
    }),
  ],
}

export const BlankDVD = {
  args: {
    mode: 'burn',
    project: mockProject,
    devices: mockDevices,
    currentDevicePath: '/dev/sr0',
    mediaInfo: mockMediaInfo,
    speeds: mockSpeeds,
    isBurning: false,
    mediaCapacityBytes: 4700000000,
  },
}

export const NoDevice = {
  args: {
    mode: 'burn',
    project: mockProject,
    devices: mockDevices,
    currentDevicePath: '',
    mediaInfo: null,
    speeds: [],
    isBurning: false,
    mediaCapacityBytes: 0,
  },
}

export const NoMedia = {
  args: {
    mode: 'burn',
    project: mockProject,
    devices: mockDevices,
    currentDevicePath: '/dev/sr0',
    mediaInfo: null,
    speeds: [],
    isBurning: false,
    mediaCapacityBytes: 0,
  },
}

export const NotEnoughSpace = {
  args: {
    mode: 'burn',
    project: { ...mockProject, totalSize: 5000000000 },
    devices: mockDevices,
    currentDevicePath: '/dev/sr0',
    mediaInfo: mockMediaInfo,
    speeds: mockSpeeds,
    isBurning: false,
    mediaCapacityBytes: 4700000000,
  },
}

export const RewritableDisc = {
  args: {
    mode: 'burn',
    project: mockProject,
    devices: mockDevices,
    currentDevicePath: '/dev/sr0',
    mediaInfo: {
      mediaType: 'DVD-RW',
      mediaStatus: 'closed',
      totalCapacity: 4700000000,
      freeSpace: 0,
      erasable: true,
      sessions: 1,
    },
    speeds: mockSpeeds,
    isBurning: false,
    mediaCapacityBytes: 4700000000,
  },
}

export const SaveMode = {
  args: {
    mode: 'save',
    project: mockProject,
    devices: [],
    currentDevicePath: '',
    mediaInfo: null,
    speeds: [],
    isBurning: false,
    mediaCapacityBytes: 0,
  },
}

export const AppendableMedia = {
  args: {
    mode: 'burn',
    project: mockProject,
    devices: mockDevices,
    currentDevicePath: '/dev/sr0',
    mediaInfo: {
      mediaType: 'BD-R',
      mediaStatus: 'appendable',
      totalCapacity: 25000000000,
      freeSpace: 20000000000,
      erasable: false,
      sessions: 1,
    },
    speeds: mockSpeeds,
    isBurning: false,
    mediaCapacityBytes: 25000000000,
  },
}
