import DiscLayout from './DiscLayout.vue'
import { useTabStore } from '../../stores/tabStore'

export default {
  title: 'Project/DiscLayout',
  component: DiscLayout,
  parameters: {
    layout: 'fullscreen',
  },
}

export const Empty = {
  decorators: [
    (story) => ({
      components: { story },
      setup() {
        const tabStore = useTabStore()
        tabStore.tabs = [
          {
            id: 'tab1',
            label: 'Empty Project',
            projectData: {
              name: 'Empty Project',
              entries: [],
              selectedProjectEntries: [],
              totalSize: 0,
              modified: false,
            },
          },
        ]
        tabStore.activeTabId = 'tab1'
        return {}
      },
      template: '<div style="width: 400px; height: 500px;"><story /></div>',
    }),
  ],
}

export const WithEntries = {
  decorators: [
    (story) => ({
      components: { story },
      setup() {
        const tabStore = useTabStore()
        tabStore.tabs = [
          {
            id: 'tab1',
            label: 'My Backup',
            projectData: {
              name: 'My Backup',
              entries: [
                { sourcePath: '/home/user/photo.jpg', destPath: '/photo.jpg', name: 'photo.jpg', isDir: false, size: 3145728 },
                { sourcePath: '/home/user/docs', destPath: '/docs', name: 'docs', isDir: true, size: 0 },
                { sourcePath: '/home/user/docs/report.pdf', destPath: '/docs/report.pdf', name: 'report.pdf', isDir: false, size: 1048576 },
                { sourcePath: '/home/user/docs/notes.txt', destPath: '/docs/notes.txt', name: 'notes.txt', isDir: false, size: 512 },
                { sourcePath: '/home/user/backup.tar.gz', destPath: '/backup.tar.gz', name: 'backup.tar.gz', isDir: false, size: 52428800 },
              ],
              selectedProjectEntries: [],
              totalSize: 56623616,
              modified: false,
            },
          },
        ]
        tabStore.activeTabId = 'tab1'
        return {}
      },
      template: '<div style="width: 400px; height: 500px;"><story /></div>',
    }),
  ],
}

export const WithSelection = {
  decorators: [
    (story) => ({
      components: { story },
      setup() {
        const tabStore = useTabStore()
        tabStore.tabs = [
          {
            id: 'tab1',
            label: 'Selected Items',
            projectData: {
              name: 'Selected Items',
              entries: [
                { sourcePath: '/home/user/file1.txt', destPath: '/file1.txt', name: 'file1.txt', isDir: false, size: 1024 },
                { sourcePath: '/home/user/file2.txt', destPath: '/file2.txt', name: 'file2.txt', isDir: false, size: 2048 },
                { sourcePath: '/home/user/file3.txt', destPath: '/file3.txt', name: 'file3.txt', isDir: false, size: 4096 },
              ],
              selectedProjectEntries: ['/file1.txt', '/file3.txt'],
              totalSize: 7168,
              modified: true,
            },
          },
        ]
        tabStore.activeTabId = 'tab1'
        return {}
      },
      template: '<div style="width: 400px; height: 500px;"><story /></div>',
    }),
  ],
}
