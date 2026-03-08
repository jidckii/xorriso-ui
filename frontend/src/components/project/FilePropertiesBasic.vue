<script setup>
import { useI18n } from 'vue-i18n'
import { formatBytes } from '../../composables/useFormatBytes'

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
</template>
