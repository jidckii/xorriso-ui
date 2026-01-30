<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Modal from '../ui/Modal.vue'
import FileIcon from '../ui/FileIcon.vue'
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

    // Load image preview if it's an image
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

function formatBytes(bytes) {
  return projectStore.formatBytes(bytes)
}

function formatDate(isoString) {
  if (!isoString) return '—'
  return new Date(isoString).toLocaleString()
}

function fileName(path) {
  if (!path) return ''
  return path.split('/').pop() || path
}
</script>

<template>
  <Modal :show="show" :title="t('fileProperties.title')" :size="modalSize" @close="emit('close')">
    <div v-if="loading" class="flex items-center justify-center py-8">
      <span class="text-sm text-gray-500">{{ t('fileProperties.loading') }}</span>
    </div>

    <div v-else-if="properties">
      <!-- Header with icon and name -->
      <div class="flex items-center gap-3 pb-3 border-b border-gray-300 dark:border-gray-600">
        <FileIcon :name="properties.name" :is-dir="properties.isDir" :size="32" />
        <div class="min-w-0 flex-1">
          <div class="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
            {{ properties.name }}
          </div>
          <div class="text-xs text-gray-500 truncate">{{ properties.path }}</div>
        </div>
      </div>

      <!-- Two-column layout when extra data exists -->
      <div :class="hasExtraData || previewUrl ? 'flex gap-6 pt-4' : 'pt-4'">
        <!-- Left column: preview + basic properties -->
        <div :class="hasExtraData || previewUrl ? 'flex-1 min-w-0' : ''">
          <!-- Image preview -->
          <div v-if="previewUrl" class="flex justify-center pb-4">
            <img
              :src="previewUrl"
              class="rounded max-w-full max-h-[250px] object-contain"
              alt=""
            />
          </div>

          <!-- Properties table -->
          <table class="w-full text-sm">
            <tbody>
              <tr class="border-b border-gray-200 dark:border-gray-700/50">
                <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.size') }}</td>
                <td class="py-1.5 text-gray-900 dark:text-gray-100">
                  {{ formatBytes(properties.size) }}
                  <span v-if="properties.size > 1024" class="text-gray-500 ml-1">
                    ({{ properties.size.toLocaleString() }} B)
                  </span>
                </td>
              </tr>

              <tr v-if="properties.isDir" class="border-b border-gray-200 dark:border-gray-700/50">
                <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.items') }}</td>
                <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ properties.itemCount }}</td>
              </tr>

              <tr v-if="properties.mimeType" class="border-b border-gray-200 dark:border-gray-700/50">
                <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.type') }}</td>
                <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ properties.mimeType }}</td>
              </tr>

              <tr v-if="properties.imageWidth" class="border-b border-gray-200 dark:border-gray-700/50">
                <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.dimensions') }}</td>
                <td class="py-1.5 text-gray-900 dark:text-gray-100">
                  {{ properties.imageWidth }} × {{ properties.imageHeight }} px
                </td>
              </tr>

              <tr v-if="properties.duration" class="border-b border-gray-200 dark:border-gray-700/50">
                <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.duration') }}</td>
                <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ properties.duration }}</td>
              </tr>

              <tr class="border-b border-gray-200 dark:border-gray-700/50">
                <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.modified') }}</td>
                <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ formatDate(properties.modTime) }}</td>
              </tr>

              <tr class="border-b border-gray-200 dark:border-gray-700/50">
                <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.accessed') }}</td>
                <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ formatDate(properties.accessTime) }}</td>
              </tr>

              <tr class="border-b border-gray-200 dark:border-gray-700/50">
                <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.permissions') }}</td>
                <td class="py-1.5 text-gray-900 dark:text-gray-100 font-mono text-xs">
                  {{ properties.permissions }}
                </td>
              </tr>

              <tr>
                <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.owner') }}</td>
                <td class="py-1.5 text-gray-900 dark:text-gray-100">
                  {{ properties.owner }}:{{ properties.group }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Right column: EXIF / Media metadata -->
        <div v-if="hasExtraData" class="flex-1 min-w-0">
          <!-- EXIF metadata (photos) -->
          <div v-if="hasExifData">
            <div class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">
              {{ t('fileProperties.exifSection') }}
            </div>
            <table class="w-full text-sm">
              <tbody>
                <tr v-if="properties.cameraMake || properties.cameraModel" class="border-b border-gray-200 dark:border-gray-700/50">
                  <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.camera') }}</td>
                  <td class="py-1.5 text-gray-900 dark:text-gray-100">
                    {{ [properties.cameraMake, properties.cameraModel].filter(Boolean).join(' ') }}
                  </td>
                </tr>
                <tr v-if="properties.dateTaken" class="border-b border-gray-200 dark:border-gray-700/50">
                  <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.dateTaken') }}</td>
                  <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ formatDate(properties.dateTaken) }}</td>
                </tr>
                <tr v-if="properties.fNumber" class="border-b border-gray-200 dark:border-gray-700/50">
                  <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.aperture') }}</td>
                  <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ properties.fNumber }}</td>
                </tr>
                <tr v-if="properties.exposureTime" class="border-b border-gray-200 dark:border-gray-700/50">
                  <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.exposure') }}</td>
                  <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ properties.exposureTime }}</td>
                </tr>
                <tr v-if="properties.isoSpeed" class="border-b border-gray-200 dark:border-gray-700/50">
                  <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.iso') }}</td>
                  <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ properties.isoSpeed }}</td>
                </tr>
                <tr v-if="properties.focalLength" class="border-b border-gray-200 dark:border-gray-700/50">
                  <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.focalLength') }}</td>
                  <td class="py-1.5 text-gray-900 dark:text-gray-100">
                    {{ properties.focalLength }}
                    <span v-if="properties.focalLength35" class="text-gray-500 ml-1">
                      ({{ properties.focalLength35 }} {{ t('fileProperties.equiv35mm') }})
                    </span>
                  </td>
                </tr>
                <tr v-if="properties.flash" class="border-b border-gray-200 dark:border-gray-700/50">
                  <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.flash') }}</td>
                  <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ properties.flash }}</td>
                </tr>
                <tr v-if="properties.orientation">
                  <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.orientation') }}</td>
                  <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ properties.orientation }}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Video/Audio metadata -->
          <div v-if="hasMediaData" :class="hasExifData ? 'pt-4' : ''">
            <div class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">
              {{ t('fileProperties.mediaSection') }}
            </div>
            <table class="w-full text-sm">
              <tbody>
                <tr v-if="properties.videoCodec" class="border-b border-gray-200 dark:border-gray-700/50">
                  <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.videoCodec') }}</td>
                  <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ properties.videoCodec }}</td>
                </tr>
                <tr v-if="properties.audioCodec" class="border-b border-gray-200 dark:border-gray-700/50">
                  <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.audioCodec') }}</td>
                  <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ properties.audioCodec }}</td>
                </tr>
                <tr v-if="properties.frameRate" class="border-b border-gray-200 dark:border-gray-700/50">
                  <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.frameRate') }}</td>
                  <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ properties.frameRate }}</td>
                </tr>
                <tr v-if="properties.videoBitrate" class="border-b border-gray-200 dark:border-gray-700/50">
                  <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.videoBitrate') }}</td>
                  <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ properties.videoBitrate }}</td>
                </tr>
                <tr v-if="properties.sampleRate" class="border-b border-gray-200 dark:border-gray-700/50">
                  <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.sampleRate') }}</td>
                  <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ properties.sampleRate }}</td>
                </tr>
                <tr v-if="properties.audioBitrate" class="border-b border-gray-200 dark:border-gray-700/50">
                  <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.audioBitrate') }}</td>
                  <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ properties.audioBitrate }}</td>
                </tr>
                <tr v-if="properties.audioChannels">
                  <td class="py-1.5 pr-3 text-gray-500 whitespace-nowrap">{{ t('fileProperties.channels') }}</td>
                  <td class="py-1.5 text-gray-900 dark:text-gray-100">{{ properties.audioChannels }}</td>
                </tr>
              </tbody>
            </table>
          </div>
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
