<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Modal from '../ui/Modal.vue'
import FileIcon from '../ui/FileIcon.vue'
import FilePropertiesBasic from './FilePropertiesBasic.vue'
import FilePropertiesExif from './FilePropertiesExif.vue'
import FilePropertiesMedia from './FilePropertiesMedia.vue'
import { useProjectStore } from '../../stores/projectStore'

const props = defineProps({
  show: { type: Boolean, default: false },
  filePath: { type: String, default: '' },
})

const emit = defineEmits(['close'])

const { t } = useI18n()
const projectStore = useProjectStore()

const properties = ref(null)
const loading = ref(false)
const previewUrl = ref('')

watch(() => props.show, async (val) => {
  if (val && props.filePath) {
    loading.value = true
    properties.value = null
    previewUrl.value = ''

    const result = await projectStore.getFileProperties(props.filePath)
    properties.value = result
    loading.value = false

    // Загрузка превью для изображений
    if (result && result.mimeType && result.mimeType.startsWith('image/')) {
      previewUrl.value = await projectStore.getImagePreview(props.filePath, 300)
    }
  }
})

const hasExifData = computed(() => {
  const p = properties.value
  return p && (p.cameraMake || p.cameraModel || p.fNumber || p.exposureTime || p.isoSpeed || p.focalLength || p.dateTaken || p.flash)
})

const hasMediaData = computed(() => {
  const p = properties.value
  return p && (p.videoCodec || p.audioCodec || p.frameRate || p.sampleRate)
})

const hasExtraData = computed(() => hasExifData.value || hasMediaData.value)

const modalSize = computed(() => hasExtraData.value || previewUrl.value ? 'lg' : 'sm')
</script>

<template>
  <Modal :show="show" :title="t('fileProperties.title')" :size="modalSize" @close="emit('close')">
    <div v-if="loading" class="flex items-center justify-center py-8">
      <span class="text-sm text-gray-500">{{ t('fileProperties.loading') }}</span>
    </div>

    <div v-else-if="properties">
      <!-- Заголовок с иконкой и именем файла -->
      <div class="flex items-center gap-3 pb-3 border-b border-gray-300 dark:border-gray-600">
        <FileIcon :name="properties.name" :is-dir="properties.isDir" :size="32" />
        <div class="min-w-0 flex-1">
          <div class="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
            {{ properties.name }}
          </div>
          <div class="text-xs text-gray-500 truncate">{{ properties.path }}</div>
        </div>
      </div>

      <!-- Двухколоночный layout при наличии дополнительных данных -->
      <div :class="hasExtraData || previewUrl ? 'flex gap-6 pt-4' : 'pt-4'">
        <!-- Левая колонка: превью + базовые свойства -->
        <div :class="hasExtraData || previewUrl ? 'flex-1 min-w-0' : ''">
          <!-- Превью изображения -->
          <div v-if="previewUrl" class="flex justify-center pb-4">
            <img
              :src="previewUrl"
              class="rounded max-w-full max-h-[250px] object-contain"
              alt=""
            />
          </div>

          <FilePropertiesBasic :properties="properties" />
        </div>

        <!-- Правая колонка: EXIF / Media метаданные -->
        <div v-if="hasExtraData" class="flex-1 min-w-0">
          <FilePropertiesExif v-if="hasExifData" :properties="properties" />
          <FilePropertiesMedia v-if="hasMediaData" :properties="properties" :class="hasExifData ? 'pt-4' : ''" />
        </div>
      </div>
    </div>

    <template #footer>
      <button
        class="px-4 py-1.5 text-sm rounded bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
        @click="emit('close')"
      >
        {{ t('fileProperties.close') }}
      </button>
    </template>
  </Modal>
</template>
