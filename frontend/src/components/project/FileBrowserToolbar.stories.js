import FileBrowserToolbar from './FileBrowserToolbar.vue'

export default {
  title: 'Project/FileBrowserToolbar',
  component: FileBrowserToolbar,
  render: (args) => ({
    components: { FileBrowserToolbar },
    setup() {
      return { args }
    },
    template: '<div style="width: 500px;"><FileBrowserToolbar v-bind="args" /></div>',
  }),
}

export const RootPath = {
  args: {
    browsePath: '/',
    breadcrumbs: [],
    editingPath: false,
    pathInput: '',
    showHidden: false,
    mountPoints: [
      { path: '/home/user', label: 'Home', icon: 'home' },
      { path: '/media/dvd', label: 'DVD', icon: 'usb' },
    ],
  },
}

export const NestedPath = {
  args: {
    browsePath: '/home/user/Documents/projects',
    breadcrumbs: [
      { name: 'home', path: '/home' },
      { name: 'user', path: '/home/user' },
      { name: 'Documents', path: '/home/user/Documents' },
      { name: 'projects', path: '/home/user/Documents/projects' },
    ],
    editingPath: false,
    pathInput: '',
    showHidden: false,
    mountPoints: [
      { path: '/home/user', label: 'Home', icon: 'home' },
      { path: '/media/dvd', label: 'DVD', icon: 'usb' },
    ],
  },
}

export const Editing = {
  args: {
    browsePath: '/home/user',
    breadcrumbs: [
      { name: 'home', path: '/home' },
      { name: 'user', path: '/home/user' },
    ],
    editingPath: true,
    pathInput: '/home/user',
    showHidden: false,
    mountPoints: [
      { path: '/home/user', label: 'Home', icon: 'home' },
    ],
  },
}

export const WithMountPoints = {
  args: {
    browsePath: '/media/dvd',
    breadcrumbs: [
      { name: 'media', path: '/media' },
      { name: 'dvd', path: '/media/dvd' },
    ],
    editingPath: false,
    pathInput: '',
    showHidden: true,
    mountPoints: [
      { path: '/home/user', label: 'Home', icon: 'home' },
      { path: '/media/dvd', label: 'DVD', icon: 'usb' },
      { path: '/media/usb-stick', label: 'USB Stick', icon: 'usb' },
      { path: '/mnt/backup', label: 'Backup', icon: 'usb' },
    ],
  },
}
