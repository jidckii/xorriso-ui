import FilePropertiesExif from './FilePropertiesExif.vue'

export default {
  title: 'Project/FilePropertiesExif',
  component: FilePropertiesExif,
  render: (args) => ({
    components: { FilePropertiesExif },
    setup() {
      return { args }
    },
    template: '<div style="width: 400px;" class="p-4"><FilePropertiesExif v-bind="args" /></div>',
  }),
}

export const PhotoWithExif = {
  args: {
    properties: {
      cameraMake: 'Canon',
      cameraModel: 'EOS R5',
      fNumber: 'f/2.8',
      exposureTime: '1/250s',
      isoSpeed: 400,
      focalLength: '85mm',
      focalLength35: '85mm',
      dateTaken: '2025-06-15T14:30:00Z',
      flash: 'Off',
      orientation: 'Horizontal',
    },
  },
}
