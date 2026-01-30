import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHashHistory } from 'vue-router'
import App from './App.vue'
import i18n from './i18n'
import './assets/css/main.css'

import ProjectView from './views/ProjectView.vue'
import BurnView from './views/BurnView.vue'
import SettingsView from './views/SettingsView.vue'

const routes = [
  { path: '/', name: 'project', component: ProjectView },
  { path: '/burn', name: 'burn', component: BurnView },
  { path: '/settings', name: 'settings', component: SettingsView },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.use(router)
app.use(i18n)

import { useThemeStore } from './stores/themeStore'
const themeStore = useThemeStore()
themeStore.loadThemeFromSettings()

app.mount('#app')
