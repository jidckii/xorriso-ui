<script setup>
import { computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import AppHeader from './components/layout/AppHeader.vue'
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

function onKeydown(e) {
  if (e.key === 'Escape') {
    if (tabStore.showDiscInfo) {
      tabStore.showDiscInfo = false
    }
  }
}

onMounted(async () => {
  document.addEventListener('keydown', onKeydown)
  deviceStore.init()
  if (tabStore.tabs.length === 0) {
    const tabId = tabStore.addProjectTab()
    const tab = tabStore.tabs.find(t => t.id === tabId)
    await projectStore.newProject(tabId, tab.label, tab.label)
  }
})

onUnmounted(() => {
  document.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <div class="flex flex-col h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
    <AppHeader />
    <TabBar v-if="isMainView" />
    <main class="flex-1 overflow-hidden relative">
      <template v-if="isMainView">
        <ProjectView :key="tabStore.activeTabId" />
        <!-- Disc Info Modal -->
        <Teleport to="body">
          <div v-if="tabStore.showDiscInfo" class="fixed inset-0 z-50 flex items-center justify-center">
            <div class="absolute inset-0 bg-black/50" @click="tabStore.showDiscInfo = false" />
            <div class="relative w-[900px] max-h-[85vh] bg-white dark:bg-gray-900 rounded-lg shadow-2xl border border-gray-300 dark:border-gray-700 overflow-y-auto">
              <DiscInfoView @close="tabStore.showDiscInfo = false" />
            </div>
          </div>
        </Teleport>
        <!-- Burn / Save Modal -->
        <Teleport to="body">
          <div v-if="tabStore.showBurnModal" class="fixed inset-0 z-50 flex items-center justify-center">
            <div class="absolute inset-0 bg-black/50" />
            <div class="relative w-[900px] max-h-[85vh] bg-white dark:bg-gray-900 rounded-lg shadow-2xl border border-gray-300 dark:border-gray-700 overflow-y-auto">
              <BurnOverlay :mode="tabStore.burnModalMode" @close="tabStore.closeBurnModal()" />
            </div>
          </div>
        </Teleport>
      </template>
      <router-view v-else />
    </main>
  </div>
</template>
