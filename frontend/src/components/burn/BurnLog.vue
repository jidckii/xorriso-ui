<script setup>
import { ref, watch, nextTick, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  lines: {
    type: Array,
    default: () => [],
  },
})

const logContainer = ref(null)
const autoScroll = ref(true)

watch(
  () => props.lines.length,
  async () => {
    if (autoScroll.value) {
      await nextTick()
      scrollToBottom()
    }
  }
)

function scrollToBottom() {
  if (logContainer.value) {
    logContainer.value.scrollTop = logContainer.value.scrollHeight
  }
}

function onScroll() {
  if (!logContainer.value) return
  const el = logContainer.value
  const isAtBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 30
  autoScroll.value = isAtBottom
}
</script>

<template>
  <div
    ref="logContainer"
    class="bg-gray-50 dark:bg-gray-950 border border-gray-300 dark:border-gray-700 rounded font-mono text-xs text-gray-600 dark:text-gray-400 p-3 overflow-y-auto max-h-64 min-h-[100px]"
    @scroll="onScroll"
  >
    <div v-if="lines.length === 0" class="text-gray-500 dark:text-gray-600 italic">
      {{ t('burnProgress.noLogOutput') }}
    </div>
    <div
      v-for="(line, idx) in lines"
      :key="idx"
      class="leading-5 whitespace-pre-wrap break-all"
      :class="{
        'text-red-400': line.toLowerCase().includes('error') || line.toLowerCase().includes('fatal'),
        'text-yellow-400': line.toLowerCase().includes('warn'),
        'text-green-400': line.toLowerCase().includes('success') || line.toLowerCase().includes('done'),
      }"
    >{{ line }}</div>
  </div>
</template>
