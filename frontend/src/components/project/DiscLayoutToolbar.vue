<script setup>
import { useI18n } from 'vue-i18n'
import { Save } from 'lucide-vue-next'

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
})

const emit = defineEmits(['select-all', 'deselect-all', 'remove-selected', 'open-project', 'save-project'])

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
    <!-- Кнопка открыть проект -->
    <button
      @click="emit('open-project')"
      class="flex items-center gap-1.5 px-3 py-1 text-xs font-medium rounded bg-gray-600 hover:bg-gray-500 text-white transition-colors"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
      </svg>
      {{ t('header.open') }}
    </button>
    <!-- Кнопка сохранить проект -->
    <button
      @click="emit('save-project')"
      class="flex items-center gap-1.5 px-3 py-1 text-xs font-medium rounded bg-blue-600 hover:bg-blue-500 text-white transition-colors"
    >
      <Save :size="16" />
      {{ t('header.save') }}
    </button>
  </div>
</template>
