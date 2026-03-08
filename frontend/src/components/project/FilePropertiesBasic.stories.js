import FilePropertiesBasic from './FilePropertiesBasic.vue'

export default {
  title: 'Project/FilePropertiesBasic',
  component: FilePropertiesBasic,
  render: (args) => ({
    components: { FilePropertiesBasic },
    setup() {
      return { args }
    },
    template: '<div style="width: 400px;" class="p-4"><FilePropertiesBasic v-bind="args" /></div>',
  }),
}

export const File = {
  args: {
    properties: {
      name: 'test.txt',
      size: 12345,
      isDir: false,
      mimeType: 'text/plain',
      modTime: '2025-01-15T10:30:00Z',
      accessTime: '2025-01-15T12:00:00Z',
      permissions: 'rw-r--r--',
      owner: 'user',
      group: 'users',
    },
  },
}

export const Directory = {
  args: {
    properties: {
      name: 'Documents',
      size: 4096,
      isDir: true,
      itemCount: 42,
      modTime: '2025-03-20T14:45:00Z',
      accessTime: '2025-03-21T09:00:00Z',
      permissions: 'rwxr-xr-x',
      owner: 'user',
      group: 'users',
    },
  },
}
