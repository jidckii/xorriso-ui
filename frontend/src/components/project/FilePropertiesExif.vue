<script setup>
import { useI18n } from 'vue-i18n'

defineProps({
  properties: { type: Object, required: true },
})

const { t } = useI18n()

function formatDate(isoString) {
  if (!isoString) return '—'
  return new Date(isoString).toLocaleString()
}
</script>

<template>
  <div>
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
</template>
