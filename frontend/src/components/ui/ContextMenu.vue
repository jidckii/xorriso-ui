<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  show: { type: Boolean, default: false },
  x: { type: Number, default: 0 },
  y: { type: Number, default: 0 },
  items: { type: Array, default: () => [] },
})

const emit = defineEmits(['close', 'select'])

const menuRef = ref(null)

function onClickOutside(e) {
  if (menuRef.value && !menuRef.value.contains(e.target)) {
    emit('close')
  }
}

function onSelect(item) {
  if (item.disabled) return
  emit('select', item.action)
  emit('close')
}

watch(() => props.show, (val) => {
  if (val) {
    setTimeout(() => document.addEventListener('mousedown', onClickOutside), 0)
  } else {
    document.removeEventListener('mousedown', onClickOutside)
  }
})

onUnmounted(() => {
  document.removeEventListener('mousedown', onClickOutside)
})
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-100 ease-out"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition duration-75 ease-in"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="show"
        ref="menuRef"
        class="fixed z-50 min-w-[180px] bg-gray-100 dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg shadow-xl py-1 origin-top-left"
        :style="{ left: x + 'px', top: y + 'px' }"
      >
        <template v-for="(item, idx) in items" :key="idx">
          <div
            v-if="item.separator"
            class="my-1 border-t border-gray-300 dark:border-gray-600"
          />
          <button
            v-else
            class="w-full flex items-center gap-2.5 px-3 py-1.5 text-sm text-left transition-colors"
            :class="item.disabled
              ? 'text-gray-400 dark:text-gray-600 cursor-not-allowed'
              : 'text-gray-800 dark:text-gray-200 hover:bg-gray-200 dark:hover:bg-gray-700'"
            @click="onSelect(item)"
          >
            <component v-if="item.icon" :is="item.icon" :size="15" class="shrink-0 opacity-70" />
            <span>{{ item.label }}</span>
          </button>
        </template>
      </div>
    </Transition>
  </Teleport>
</template>
