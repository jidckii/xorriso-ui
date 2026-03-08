import MediaInfo from './MediaInfo.vue'

export default {
  title: 'Device/MediaInfo',
  component: MediaInfo,
  argTypes: {
    mediaType: { control: 'text' },
    capacityUsed: { control: 'number' },
    capacityTotal: { control: 'number' },
    status: {
      control: 'select',
      options: ['blank', 'appendable', 'closed', 'unknown'],
    },
  },
  render: (args) => ({
    components: { MediaInfo },
    setup() {
      return { args }
    },
    template: '<MediaInfo v-bind="args" />',
  }),
}

export const NoMedia = {
  args: {
    mediaType: '',
    capacityUsed: 0,
    capacityTotal: 0,
    status: 'unknown',
  },
}

export const BlankDVD = {
  args: {
    mediaType: 'DVD+R',
    capacityUsed: 0,
    capacityTotal: 4700000000,
    status: 'blank',
  },
}

export const AppendableBD = {
  args: {
    mediaType: 'BD-R',
    capacityUsed: 8500000000,
    capacityTotal: 25000000000,
    status: 'appendable',
  },
}

export const ClosedCD = {
  args: {
    mediaType: 'CD-R',
    capacityUsed: 700000000,
    capacityTotal: 700000000,
    status: 'closed',
  },
}
