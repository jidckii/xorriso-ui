import DeviceSelector from './DeviceSelector.vue'

export default {
  title: 'Device/DeviceSelector',
  component: DeviceSelector,
  argTypes: {
    devices: { control: 'object' },
    modelValue: { control: 'text' },
  },
  render: (args) => ({
    components: { DeviceSelector },
    setup() {
      return { args }
    },
    template: '<DeviceSelector v-bind="args" />',
  }),
}

export const Empty = {
  args: {
    devices: [],
    modelValue: '',
  },
}

export const SingleDevice = {
  args: {
    devices: [
      { id: '/dev/sr0', vendor: 'ASUS', model: 'BW-16D1HT' },
    ],
    modelValue: '/dev/sr0',
  },
}

export const MultipleDevices = {
  args: {
    devices: [
      { id: '/dev/sr0', vendor: 'ASUS', model: 'BW-16D1HT' },
      { id: '/dev/sr1', vendor: 'Pioneer', model: 'BD-RW BDR-212D' },
      { id: '/dev/sr2', vendor: 'LG', model: 'WH16NS60' },
    ],
    modelValue: '/dev/sr0',
  },
}
