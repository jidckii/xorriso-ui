<script setup>
import { useI18n } from 'vue-i18n'
import { ArrowUpAZ, ArrowDownZA, ArrowUp01, ArrowDown10, Clock, FileType } from 'lucide-vue-next'

const { t } = useI18n()

const props = defineProps({
  sortBy: { type: String, default: 'name' },
  sortDir: { type: String, default: 'asc' },
})

const emit = defineEmits(['toggle-sort'])
</script>

<template>
  <div class="flex items-center gap-0.5">
    <!-- Имя -->
    <button
      @click="emit('toggle-sort', 'name')"
      class="p-1 rounded transition-colors"
      :class="sortBy === 'name' ? 'text-blue-400 hover:bg-blue-500/10' : 'text-gray-500 hover:bg-gray-200 dark:hover:bg-gray-700'"
      :title="t('sort.name')"
    >
      <ArrowDownZA v-if="sortBy === 'name' && sortDir === 'desc'" :size="14" />
      <ArrowUpAZ v-else :size="14" />
    </button>

    <!-- Размер -->
    <button
      @click="emit('toggle-sort', 'size')"
      class="p-1 rounded transition-colors"
      :class="sortBy === 'size' ? 'text-blue-400 hover:bg-blue-500/10' : 'text-gray-500 hover:bg-gray-200 dark:hover:bg-gray-700'"
      :title="t('sort.size')"
    >
      <ArrowDown10 v-if="sortBy === 'size' && sortDir === 'desc'" :size="14" />
      <ArrowUp01 v-else :size="14" />
    </button>

    <!-- Дата -->
    <button
      @click="emit('toggle-sort', 'date')"
      class="p-1 rounded transition-colors"
      :class="sortBy === 'date' ? 'text-blue-400 hover:bg-blue-500/10' : 'text-gray-500 hover:bg-gray-200 dark:hover:bg-gray-700'"
      :title="t('sort.date')"
    >
      <Clock :size="14" />
    </button>

    <!-- Тип -->
    <button
      @click="emit('toggle-sort', 'type')"
      class="p-1 rounded transition-colors"
      :class="sortBy === 'type' ? 'text-blue-400 hover:bg-blue-500/10' : 'text-gray-500 hover:bg-gray-200 dark:hover:bg-gray-700'"
      :title="t('sort.type')"
    >
      <FileType :size="14" />
    </button>
  </div>
</template>
