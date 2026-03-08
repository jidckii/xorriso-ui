import ImagePreviewTooltip from './ImagePreviewTooltip.vue'

export default {
  title: 'UI/ImagePreviewTooltip',
  component: ImagePreviewTooltip,
  parameters: {
    layout: 'fullscreen',
  },
  argTypes: {
    filePath: { control: 'text' },
    visible: { control: 'boolean' },
    x: { control: 'number' },
    y: { control: 'number' },
  },
}

export const Hidden = {
  args: {
    visible: false,
    filePath: 'test.jpg',
    x: 200,
    y: 200,
  },
}

export const Loading = {
  args: {
    visible: true,
    filePath: 'test.jpg',
    x: 200,
    y: 300,
  },
}

export const Visible = {
  args: {
    visible: true,
    filePath: 'test.jpg',
    x: 200,
    y: 200,
  },
}
