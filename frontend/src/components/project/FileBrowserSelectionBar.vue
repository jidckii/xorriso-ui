<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

defineProps({
  allSelected: { type: Boolean, default: false },
  selectedCount: { type: Number, default: 0 },
})

const emit = defineEmits(['select-all', 'deselect-all', 'add-selected'])
</script>

<template>
  <div class="flex items-center gap-2 px-3 py-2 bg-gray-100 dark:bg-gray-800 border-t border-gray-300 dark:border-gray-700">
    <label class="flex items-center gap-1.5 text-xs text-gray-500 cursor-pointer select-none">
      <input
        type="checkbox"
        :checked="allSelected"
        @change="allSelected ? emit('deselect-all') : emit('select-all')"
        class="w-3.5 h-3.5 accent-blue-600 cursor-pointer"
      />
      {{ t('project.selectAll') }}
    </label>
    <button
      @click="emit('add-selected')"
      :disabled="selectedCount === 0"
      class="px-3 py-1 text-xs font-medium rounded bg-blue-600 hover:bg-blue-500 text-white disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
    >
      {{ t('project.addToProject') }}
    </button>
    <span class="text-xs text-gray-500">
      {{ selectedCount }} {{ t('project.selected') }}
    </span>
  </div>
</template>
