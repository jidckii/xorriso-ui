import FileBrowserItem from './FileBrowserItem.vue'

export default {
  title: 'Project/FileBrowserItem',
  component: FileBrowserItem,
  render: (args) => ({
    components: { FileBrowserItem },
    setup() {
      return { args }
    },
    template: '<div style="width: 400px;"><FileBrowserItem v-bind="args" /></div>',
  }),
}

export const File = {
  args: {
    entry: {
      name: 'document.pdf',
      sourcePath: '/home/user/document.pdf',
      isDir: false,
      size: 1048576,
    },
    depth: 0,
    expandedDirs: new Set(),
    dirChildren: {},
    selectedPaths: new Set(),
    showHidden: false,
  },
}

export const Directory = {
  args: {
    entry: {
      name: 'Documents',
      sourcePath: '/home/user/Documents',
      isDir: true,
      size: 0,
    },
    depth: 0,
    expandedDirs: new Set(),
    dirChildren: {},
    selectedPaths: new Set(),
    showHidden: false,
  },
}

export const SelectedFile = {
  args: {
    entry: {
      name: 'report.xlsx',
      sourcePath: '/home/user/report.xlsx',
      isDir: false,
      size: 524288,
    },
    depth: 0,
    expandedDirs: new Set(),
    dirChildren: {},
    selectedPaths: new Set(['/home/user/report.xlsx']),
    showHidden: false,
  },
}

export const ImageFile = {
  args: {
    entry: {
      name: 'photo.jpg',
      sourcePath: '/home/user/photo.jpg',
      isDir: false,
      size: 3145728,
    },
    depth: 0,
    expandedDirs: new Set(),
    dirChildren: {},
    selectedPaths: new Set(),
    showHidden: false,
  },
}

export const NestedDepth = {
  args: {
    entry: {
      name: 'nested-file.txt',
      sourcePath: '/home/user/docs/sub/nested-file.txt',
      isDir: false,
      size: 256,
    },
    depth: 2,
    expandedDirs: new Set(),
    dirChildren: {},
    selectedPaths: new Set(),
    showHidden: false,
  },
}
