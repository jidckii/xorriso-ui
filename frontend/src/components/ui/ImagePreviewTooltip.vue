<script setup>
import { ref, watch, onUnmounted } from 'vue'
import { useProjectStore } from '../../stores/projectStore'

const props = defineProps({
  filePath: { type: String, default: '' },
  visible: { type: Boolean, default: false },
  x: { type: Number, default: 0 },
  y: { type: Number, default: 0 },
})

const projectStore = useProjectStore()
const previewUrl = ref('')
const loading = ref(false)
let loadTimer = null
let currentPath = ''

// Simple cache: path -> data URL (max 50 entries)
const cache = new Map()
const CACHE_MAX = 50

// Image extensions that support preview
const imageExts = new Set(['.jpg', '.jpeg', '.png', '.gif', '.webp', '.bmp'])

function isImageFile(path) {
  if (!path) return false
  const ext = path.substring(path.lastIndexOf('.')).toLowerCase()
  return imageExts.has(ext)
}

watch(() => props.visible, (visible) => {
  if (!visible) {
    clearTimeout(loadTimer)
    loadTimer = null
    previewUrl.value = ''
    currentPath = ''
    loading.value = false
    return
  }

  if (!isImageFile(props.filePath)) {
    previewUrl.value = ''
    return
  }

  const path = props.filePath
  currentPath = path

  // Check cache first
  if (cache.has(path)) {
    previewUrl.value = cache.get(path)
    loading.value = false
    return
  }

  // Debounce: load after 300ms hover
  loading.value = true
  loadTimer = setTimeout(async () => {
    if (currentPath !== path) return
    const url = await projectStore.getImagePreview(path, 200)
    if (currentPath === path) {
      previewUrl.value = url
      loading.value = false
      // Cache the result
      if (url) {
        if (cache.size >= CACHE_MAX) {
          const first = cache.keys().next().value
          cache.delete(first)
        }
        cache.set(path, url)
      }
    }
  }, 300)
})

onUnmounted(() => {
  clearTimeout(loadTimer)
})

defineExpose({ isImageFile })
</script>

<template>
  <Teleport to="body">
    <div
      v-if="visible && (previewUrl || loading) && isImageFile(filePath)"
      class="fixed z-50 pointer-events-none"
      :style="{ left: x + 'px', top: y + 'px', transform: 'translate(12px, -50%)' }"
    >
      <div class="bg-gray-900 border border-gray-700 rounded-lg shadow-xl p-1.5 max-w-[220px]">
        <img
          v-if="previewUrl"
          :src="previewUrl"
          class="rounded max-w-[200px] max-h-[200px] object-contain"
          alt=""
        />
        <div v-else-if="loading" class="w-[100px] h-[80px] flex items-center justify-center">
          <span class="text-xs text-gray-400">...</span>
        </div>
      </div>
    </div>
  </Teleport>
</template>
