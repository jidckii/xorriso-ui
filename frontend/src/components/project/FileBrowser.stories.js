import FileBrowser from './FileBrowser.vue'
import { useTabStore } from '../../stores/tabStore'

export default {
  title: 'Project/FileBrowser',
  component: FileBrowser,
  parameters: {
    layout: 'fullscreen',
  },
}

export const EmptyDirectory = {
  decorators: [
    (story) => ({
      components: { story },
      setup() {
        const tabStore = useTabStore()
        tabStore.tabs = [
          {
            id: 'tab1',
            label: 'Test Project',
            projectData: {
              name: 'Test Project',
              browsePath: '/home/user',
              browseShowHidden: false,
              browseEntries: [],
              entries: [],
              selectedBrowseFiles: [],
              selectedProjectEntries: [],
              totalSize: 0,
              modified: false,
            },
          },
        ]
        tabStore.activeTabId = 'tab1'
        return {}
      },
      template: '<div style="width: 500px; height: 500px;"><story /></div>',
    }),
  ],
}

export const WithHiddenFiles = {
  decorators: [
    (story) => ({
      components: { story },
      setup() {
        const tabStore = useTabStore()
        tabStore.tabs = [
          {
            id: 'tab1',
            label: 'Test Project',
            projectData: {
              name: 'Test Project',
              browsePath: '/home/user',
              browseShowHidden: true,
              browseEntries: [],
              entries: [],
              selectedBrowseFiles: [],
              selectedProjectEntries: [],
              totalSize: 0,
              modified: false,
            },
          },
        ]
        tabStore.activeTabId = 'tab1'
        return {}
      },
      template: '<div style="width: 500px; height: 500px;"><story /></div>',
    }),
  ],
}
