<script setup>
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import AppHeader from './components/layout/AppHeader.vue'
import AppStatusBar from './components/layout/AppStatusBar.vue'
import TabBar from './components/layout/TabBar.vue'
import ProjectView from './views/ProjectView.vue'
import DiscInfoView from './views/DiscInfoView.vue'
import { useTabStore } from './stores/tabStore'
import { useDeviceStore } from './stores/deviceStore'

const route = useRoute()
const tabStore = useTabStore()
const deviceStore = useDeviceStore()

const isMainView = computed(() => route.path === '/')

onMounted(() => {
  deviceStore.init()
  if (tabStore.tabs.length === 0) {
    tabStore.addProjectTab('Untitled Project', 'UNTITLED')
  }
})
</script>

<template>
  <div class="flex flex-col h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
    <AppHeader />
    <TabBar v-if="isMainView" />
    <main class="flex-1 overflow-hidden relative">
      <template v-if="isMainView">
        <ProjectView :key="tabStore.activeTabId" />
        <!-- Disc Info overlay -->
        <div
          v-if="tabStore.showDiscInfo"
          class="absolute inset-0 z-10 bg-white dark:bg-gray-900"
        >
          <DiscInfoView @close="tabStore.showDiscInfo = false" />
        </div>
      </template>
      <router-view v-else />
    </main>
    <AppStatusBar />
  </div>
</template>
