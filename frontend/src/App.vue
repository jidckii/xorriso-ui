<script setup>
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import AppHeader from './components/layout/AppHeader.vue'
import AppStatusBar from './components/layout/AppStatusBar.vue'
import TabBar from './components/layout/TabBar.vue'
import ProjectView from './views/ProjectView.vue'
import DiscInfoView from './views/DiscInfoView.vue'
import BurnOverlay from './components/burn/BurnOverlay.vue'
import { useTabStore } from './stores/tabStore'
import { useDeviceStore } from './stores/deviceStore'
import { useProjectStore } from './stores/projectStore'

const route = useRoute()
const tabStore = useTabStore()
const deviceStore = useDeviceStore()
const projectStore = useProjectStore()

const isMainView = computed(() => route.path === '/')

onMounted(async () => {
  deviceStore.init()
  if (tabStore.tabs.length === 0) {
    const tabId = tabStore.addProjectTab()
    const tab = tabStore.tabs.find(t => t.id === tabId)
    await projectStore.newProject(tabId, tab.label, tab.label)
  }
})
</script>

<template>
  <div class="flex flex-col h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
    <AppHeader />
    <TabBar v-if="isMainView" />
    <main class="flex-1 overflow-hidden relative">
      <template v-if="isMainView">
        <ProjectView v-show="!tabStore.showDiscInfo && !tabStore.showBurn" :key="tabStore.activeTabId" />
        <!-- Disc Info -->
        <DiscInfoView v-if="tabStore.showDiscInfo" @close="tabStore.showDiscInfo = false" />
        <!-- Burn Overlay -->
        <BurnOverlay v-if="tabStore.showBurn" @close="tabStore.closeBurn()" />
      </template>
      <router-view v-else />
    </main>
    <AppStatusBar />
  </div>
</template>
