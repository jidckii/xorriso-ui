import TabBar from './TabBar.vue'
import { useTabStore } from '../../stores/tabStore'

export default {
  title: 'Layout/TabBar',
  component: TabBar,
  parameters: {
    layout: 'fullscreen',
  },
}

export const SingleTab = {
  decorators: [
    (story) => ({
      components: { story },
      setup() {
        const tabStore = useTabStore()
        tabStore.tabs = [
          { id: 'tab1', label: 'My Project', projectData: { modified: false } },
        ]
        tabStore.activeTabId = 'tab1'
        return {}
      },
      template: '<div style="width: 600px;"><story /></div>',
    }),
  ],
}

export const MultipleTabs = {
  decorators: [
    (story) => ({
      components: { story },
      setup() {
        const tabStore = useTabStore()
        tabStore.tabs = [
          { id: 'tab1', label: 'Backup 2026', projectData: { modified: false } },
          { id: 'tab2', label: 'Photos Archive', projectData: { modified: false } },
          { id: 'tab3', label: 'System ISO', projectData: { modified: false } },
        ]
        tabStore.activeTabId = 'tab2'
        return {}
      },
      template: '<div style="width: 600px;"><story /></div>',
    }),
  ],
}

export const ModifiedTab = {
  decorators: [
    (story) => ({
      components: { story },
      setup() {
        const tabStore = useTabStore()
        tabStore.tabs = [
          { id: 'tab1', label: 'Backup 2026', projectData: { modified: false } },
          { id: 'tab2', label: 'Unsaved Project', projectData: { modified: true } },
          { id: 'tab3', label: 'System ISO', projectData: { modified: false } },
        ]
        tabStore.activeTabId = 'tab2'
        return {}
      },
      template: '<div style="width: 600px;"><story /></div>',
    }),
  ],
}
