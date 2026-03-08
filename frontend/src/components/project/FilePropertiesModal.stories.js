import FilePropertiesModal from './FilePropertiesModal.vue'

export default {
  title: 'Project/FilePropertiesModal',
  component: FilePropertiesModal,
  parameters: {
    layout: 'fullscreen',
  },
  render: (args) => ({
    components: { FilePropertiesModal },
    setup() {
      return { args }
    },
    template: '<FilePropertiesModal v-bind="args" />',
  }),
}

export const Loading = {
  args: {
    show: true,
    filePath: '/home/user/test.txt',
  },
}

export const Closed = {
  args: {
    show: false,
    filePath: '',
  },
}
