<script setup>
import { computed } from 'vue'
import { Folder, FolderOpen } from 'lucide-vue-next'
import { getFileIconSvg } from '../../utils/fileIcons'

const props = defineProps({
  name: { type: String, required: true },
  isDir: { type: Boolean, default: false },
  isOpen: { type: Boolean, default: false },
  size: { type: Number, default: 16 },
})

const iconSvg = computed(() => {
  if (props.isDir) return null
  return getFileIconSvg(props.name)
})
</script>

<template>
  <component
    v-if="isDir"
    :is="isOpen ? FolderOpen : Folder"
    :size="size"
    class="text-yellow-500 shrink-0"
  />
  <span
    v-else
    class="inline-flex shrink-0 [&>svg]:w-full [&>svg]:h-full"
    :style="{ width: size + 'px', height: size + 'px' }"
    v-html="iconSvg"
  />
</template>
