import CapacityBar from './CapacityBar.vue'

export default {
  title: 'Project/CapacityBar',
  component: CapacityBar,
  render: (args) => ({
    components: { CapacityBar },
    setup() {
      return { args }
    },
    template: '<div style="width: 500px;"><CapacityBar v-bind="args" /></div>',
  }),
}

export const Empty = {
  args: {
    usedBytes: 0,
    mediaType: 'DVD',
  },
}

export const HalfDVD = {
  args: {
    usedBytes: 2.35e9,
    mediaType: 'DVD',
  },
}

export const WarningDVD = {
  args: {
    usedBytes: 4.2e9,
    mediaType: 'DVD',
  },
}

export const OverflowDVD = {
  args: {
    usedBytes: 5e9,
    mediaType: 'DVD',
  },
}

export const BluRay = {
  args: {
    usedBytes: 12e9,
    mediaType: 'BD',
  },
}
