import AppStatusBar from './AppStatusBar.vue'

export default {
  title: 'Layout/AppStatusBar',
  component: AppStatusBar,
  parameters: {
    layout: 'fullscreen',
  },
  argTypes: {
    deviceInfo: { control: 'text' },
    mediaType: { control: 'text' },
    freeSpace: { control: 'text' },
    operationStatus: { control: 'text' },
  },
  render: (args) => ({
    components: { AppStatusBar },
    setup() {
      return { args }
    },
    template: '<div style="width: 900px; position: fixed; bottom: 0; left: 0; right: 0;"><AppStatusBar v-bind="args" /></div>',
  }),
}

export const NoDevice = {
  args: {
    deviceInfo: '',
    mediaType: '',
    freeSpace: '',
    operationStatus: 'Ready',
  },
}

export const WithDevice = {
  args: {
    deviceInfo: 'ASUS BW-16D1HT (/dev/sr0)',
    mediaType: 'DVD+R',
    freeSpace: '4.4 GiB',
    operationStatus: 'Ready',
  },
}

export const Burning = {
  args: {
    deviceInfo: 'ASUS BW-16D1HT (/dev/sr0)',
    mediaType: 'BD-R',
    freeSpace: '12.5 GiB',
    operationStatus: 'Burning',
  },
}

export const Error = {
  args: {
    deviceInfo: 'ASUS BW-16D1HT (/dev/sr0)',
    mediaType: 'DVD+R',
    freeSpace: '0 B',
    operationStatus: 'Error',
  },
}
