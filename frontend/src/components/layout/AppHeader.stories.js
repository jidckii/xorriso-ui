import { getCurrentInstance } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import AppHeader from './AppHeader.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', component: { template: '<div />' } },
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
        if (!app.config.globalProperties.$router) {
          app.use(router)
        }
        return {}
      },
      template: '<div style="width: 100%;"><story /></div>',
    }),
  ],
}

export const Default = {}

export const OnSettingsPage = {
  decorators: [
    (story) => ({
      components: { story },
      setup() {
        router.push('/settings')
        return {}
      },
      template: '<story />',
    }),
  ],
}
