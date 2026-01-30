<script setup>
import AppHeader from './components/layout/AppHeader.vue'
import AppStatusBar from './components/layout/AppStatusBar.vue'
</script>

<template>
  <div class="flex flex-col h-screen bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100">
    <AppHeader />
    <main class="flex-1 overflow-hidden">
      <router-view />
    </main>
    <AppStatusBar />
  </div>
</template>
