import { getCurrentInstance } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import AppHeader from './AppHeader.vue'
import { useTabStore } from '../../stores/tabStore'

// Простой mock-роутер для Storybook
const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', component: { template: '<div />' } },
    { path: '/burn', component: { template: '<div />' } },
    { path: '/settings', component: { template: '<div />' } },
  ],
})

export default {
  title: 'Layout/AppHeader',
  component: AppHeader,
  parameters: {
    layout: 'fullscreen',
  },
  decorators: [
    (story) => ({
      components: { story },
      setup() {
        const app = getCurrentInstance().appContext.app
        // Регистрируем роутер, если ещё не установлен
        if (!app.config.globalProperties.$router) {
          app.use(router)
        }
        return {}
      },
      template: '<div style="width: 100%;"><story /></div>',
    }),
  ],
}

export const Default = {
  decorators: [
    (story) => ({
      components: { story },
      setup() {
        const tabStore = useTabStore()
        tabStore.tabs = [
          { id: 'tab1', label: 'My Project', projectData: { modified: false } },
        ]
        tabStore.activeTabId = 'tab1'
        tabStore.showDiscInfo = false
        return {}
      },
      template: '<story />',
    }),
  ],
}

export const DiscInfoActive = {
  decorators: [
    (story) => ({
      components: { story },
      setup() {
        const tabStore = useTabStore()
        tabStore.tabs = [
          { id: 'tab1', label: 'My Project', projectData: { modified: false } },
        ]
        tabStore.activeTabId = 'tab1'
        tabStore.showDiscInfo = true
        return {}
      },
      template: '<story />',
    }),
  ],
}
