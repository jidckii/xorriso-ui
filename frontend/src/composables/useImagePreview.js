import { ref } from 'vue'

const imageExts = new Set(['.jpg', '.jpeg', '.png', '.gif', '.webp', '.bmp'])

/**
 * Check if a filename is an image by extension.
 * @param {string} name
 * @returns {boolean}
 */
export function isImageFile(name) {
  if (!name) return false
  const dot = name.lastIndexOf('.')
  if (dot < 0) return false
  return imageExts.has(name.substring(dot).toLowerCase())
}

/**
 * Composable for image preview tooltip behavior.
 * Manages visibility, position, and source path for the preview tooltip.
 */
export function useImagePreview() {
  const previewVisible = ref(false)
  const previewX = ref(0)
  const previewY = ref(0)
  const previewSourcePath = ref('')

  function onItemMouseEnter(e, item) {
    if (item.isDir || !item.sourcePath || !isImageFile(item.name)) return
    previewSourcePath.value = item.sourcePath
    previewX.value = e.clientX
    previewY.value = e.clientY
    previewVisible.value = true
  }

  function onItemMouseMove(e) {
    previewX.value = e.clientX
    previewY.value = e.clientY
  }

  function onItemMouseLeave() {
    previewVisible.value = false
  }

  return {
    previewVisible,
    previewX,
    previewY,
    previewSourcePath,
    onItemMouseEnter,
    onItemMouseMove,
    onItemMouseLeave,
  }
}
