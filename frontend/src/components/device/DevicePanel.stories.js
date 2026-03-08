import DevicePanel from './DevicePanel.vue'

export default {
  title: 'Device/DevicePanel',
  component: DevicePanel,
  argTypes: {
    device: { control: 'object' },
  },
  render: (args) => ({
    components: { DevicePanel },
    setup() {
      return { args }
    },
    template: '<div style="width: 300px;"><DevicePanel v-bind="args" /></div>',
  }),
}

export const NoDevice = {
  args: {
    device: {
      id: '',
      vendor: '',
      model: '',
      mediaType: '',
      mediaStatus: 'unknown',
      capacityUsed: 0,
      capacityTotal: 0,
      speeds: [],
    },
  },
}

export const BlankMedia = {
  args: {
    device: {
      id: '/dev/sr0',
      vendor: 'ASUS',
      model: 'BW-16D1HT',
      mediaType: 'DVD+R',
      mediaStatus: 'blank',
      capacityUsed: 0,
      capacityTotal: 4700000000,
      speeds: [4, 8, 16],
    },
  },
}

export const AppendableMedia = {
  args: {
    device: {
      id: '/dev/sr0',
      vendor: 'ASUS',
      model: 'BW-16D1HT',
      mediaType: 'BD-R',
      mediaStatus: 'appendable',
      capacityUsed: 12500000000,
      capacityTotal: 25000000000,
      speeds: [2, 4, 6, 8],
    },
  },
}

export const ClosedMedia = {
  args: {
    device: {
      id: '/dev/sr0',
      vendor: 'ASUS',
      model: 'BW-16D1HT',
      mediaType: 'CD-R',
      mediaStatus: 'closed',
      capacityUsed: 700000000,
      capacityTotal: 700000000,
      speeds: [8, 16, 24, 32, 48],
    },
  },
}
