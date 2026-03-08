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
      class="px-4 py-1 text-xs font-medium rounded bg-orange-600 hover:bg-orange-500 text-white disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
    >
      {{ t('burn.title') }}
    </button>
  </div>
</template>
