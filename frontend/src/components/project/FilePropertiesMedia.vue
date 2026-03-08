<script setup>
import { useI18n } from 'vue-i18n'

defineProps({
  properties: { type: Object, required: true },
})

const { t } = useI18n()
</script>

<template>
  <div>
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
</template>
