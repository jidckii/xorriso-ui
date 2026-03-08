<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps({
  allSelected: {
    type: Boolean,
    default: false,
  },
  selectedCount: {
    type: Number,
    default: 0,
  },
  canBurn: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['select-all', 'deselect-all', 'remove-selected', 'go-to-burn'])

function onSelectAllChange(checked) {
  if (checked) {
    emit('select-all')
  } else {
    emit('deselect-all')
  }
}
</script>

<template>
  <div class="flex items-center gap-2 px-3 py-2 bg-gray-100 dark:bg-gray-800 border-t border-gray-300 dark:border-gray-700">
    <label class="flex items-center gap-1.5 text-xs text-gray-500 cursor-pointer select-none">
      <input
        type="checkbox"
        :checked="allSelected"
        @change="onSelectAllChange(!allSelected)"
        class="w-3.5 h-3.5 accent-blue-600 cursor-pointer"
      />
      {{ t('project.selectAll') }}
    </label>
    <button
      @click="emit('remove-selected')"
      :disabled="selectedCount === 0"
      class="px-3 py-1 text-xs font-medium rounded bg-red-600 hover:bg-red-500 text-white disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
    >
      {{ t('project.remove') }}
    </button>
    <span class="text-xs text-gray-500">
      {{ selectedCount }} {{ t('project.selected') }}
    </span>
    <span class="flex-1"></span>
    <button
      @click="emit('go-to-burn')"
      :disabled="!canBurn"
      class="flex items-center gap-1.5 px-4 py-1 text-xs font-medium rounded bg-orange-600 hover:bg-orange-500 text-white disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7C14 5 16.09 5.777 17.656 7.343A7.975 7.975 0 0120 13a7.975 7.975 0 01-2.343 5.657z" />
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M9.879 16.121A3 3 0 1012.015 11L11 14H9c0 .768.293 1.536.879 2.121z" />
      </svg>
      {{ t('header.burn') }}
    </button>
  </div>
</template>
