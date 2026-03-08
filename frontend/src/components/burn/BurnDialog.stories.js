import BurnDialog from './BurnDialog.vue'

export default {
  title: 'Burn/BurnDialog',
  component: BurnDialog,
  parameters: {
    layout: 'fullscreen',
  },
  argTypes: {
    show: { control: 'boolean' },
    availableSpeeds: { control: 'object' },
  },
}

export const Configure = {
  args: {
    show: true,
    availableSpeeds: [1, 2, 4, 8, 16],
  },
}
