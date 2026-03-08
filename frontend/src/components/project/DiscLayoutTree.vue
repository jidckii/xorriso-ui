<script setup>
import { computed } from 'vue'
import { TreeRoot, TreeItem } from 'reka-ui'
import { ChevronRight } from 'lucide-vue-next'
import FileIcon from '../ui/FileIcon.vue'
import ImagePreviewTooltip from '../ui/ImagePreviewTooltip.vue'
import { formatBytes } from '../../composables/useFormatBytes'
import { useImagePreview } from '../../composables/useImagePreview'

const props = defineProps({
  items: {
    type: Array,
    required: true,
  },
  expanded: {
    type: Array,
    default: () => [],
  },
  selectedKeys: {
    type: Set,
    required: true,
  },
})

const emit = defineEmits(['update:expanded', 'toggle-selection', 'contextmenu'])

const expandedModel = computed({
  get: () => props.expanded,
  set: (val) => emit('update:expanded', val),
})

const { previewVisible, previewX, previewY, previewSourcePath, onItemMouseEnter, onItemMouseMove, onItemMouseLeave } = useImagePreview()

function getKey(item) {
  return item._key || item.destPath
}

function getChildren(item) {
  if (!item.children || item.children.length === 0) return undefined
  return item.children
}

function isItemSelected(key) {
  return props.selectedKeys.has(key)
}
</script>

<template>
  <TreeRoot
    :items="items"
    :get-key="getKey"
    :get-children="getChildren"
    v-model:expanded="expandedModel"
    class="w-full text-sm select-none"
  >
    <template #default="{ flattenItems }">
      <TreeItem
        v-for="item in flattenItems"
        :key="item._id"
        v-bind="item.bind"
        class="flex items-center gap-1.5 py-1 cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-800/50 transition-colors outline-none"
        :class="{ 'bg-blue-500/15': isItemSelected(item.value._key) }"
        :style="{ paddingLeft: (item.level * 16 + 8) + 'px', paddingRight: '8px' }"
        @click="emit('toggle-selection', item.value)"
        @contextmenu.prevent="emit('contextmenu', item.value, $event)"
        @mouseenter="onItemMouseEnter($event, item.value)"
        @mousemove="onItemMouseMove"
        @mouseleave="onItemMouseLeave"
      >
        <template #default="{ isExpanded }">
          <!-- Стрелка раскрытия -->
          <span class="w-4 h-4 flex items-center justify-center shrink-0">
            <ChevronRight
              v-if="item.value.children?.length"
              :size="14"
              class="text-gray-500 transition-transform duration-150"
              :class="{ 'rotate-90': isExpanded }"
            />
          </span>

          <!-- Чекбокс выбора -->
          <input
            type="checkbox"
            :checked="isItemSelected(item.value._key)"
            @click.stop="emit('toggle-selection', item.value)"
            class="w-3.5 h-3.5 accent-blue-600 shrink-0 cursor-pointer"
          />

          <!-- Иконка файла/папки -->
          <FileIcon
            :name="item.value.name"
            :is-dir="item.value.isDir"
            :is-open="isExpanded"
            :size="16"
          />

          <!-- Имя -->
          <span class="truncate flex-1 text-gray-800 dark:text-gray-200">
            {{ item.value.name }}
          </span>

          <!-- Размер -->
          <span v-if="item.value.size" class="text-xs text-gray-500 shrink-0 ml-2">
            {{ formatBytes(item.value.size) }}
          </span>
        </template>
      </TreeItem>
    </template>
  </TreeRoot>

  <!-- Тултип предпросмотра изображения -->
  <ImagePreviewTooltip
    :file-path="previewSourcePath"
    :visible="previewVisible"
    :x="previewX"
    :y="previewY"
  />
</template>
