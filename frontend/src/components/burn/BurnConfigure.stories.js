import BurnConfigure from './BurnConfigure.vue'

export default {
  title: 'Burn/BurnConfigure',
  component: BurnConfigure,
  decorators: [
    () => ({
      template: '<div style="max-width: 700px;"><story /></div>',
    }),
  ],
}

const defaultProject = {
  name: 'My Data Disc',
  volumeId: 'DATA_2026',
  entries: [
    { name: 'documents', size: 1048576 },
    { name: 'photos', size: 2097152 },
    { name: 'backup.tar.gz', size: 524288 },
  ],
  totalSize: 3670016,
  burnOptions: {
    speed: 0,
    burnMode: 'auto',
    verify: true,
    eject: true,
    dummyMode: false,
    closeDisc: false,
    streamRecording: false,
    multisession: false,
    cleanupIso: false,
  },
}

const defaultDevices = [
  { path: '/dev/sr0', vendor: 'ASUS', model: 'BW-16D1HT' },
  { path: '/dev/sr1', vendor: 'Pioneer', model: 'BD-RW BDR-212D' },
]

const defaultSpeeds = [
  { value: 0, label: 'Auto (Maximum)' },
  { value: 1, label: '1x' },
  { value: 2, label: '2x' },
  { value: 4, label: '4x' },
  { value: 8, label: '8x' },
  { value: 16, label: '16x' },
]

const defaultMediaInfo = {
  mediaType: 'BD-R',
  mediaStatus: 'Blank',
  freeSpace: 25025314816,
}

export const WithProject = {
  args: {
    project: { ...defaultProject },
    devices: defaultDevices,
    currentDevicePath: '/dev/sr0',
    mediaInfo: defaultMediaInfo,
    speeds: defaultSpeeds,
    isBurning: false,
  },
}

export const NoDevice = {
  args: {
    project: { ...defaultProject },
    devices: [],
    currentDevicePath: '',
    mediaInfo: null,
    speeds: [],
    isBurning: false,
  },
}
