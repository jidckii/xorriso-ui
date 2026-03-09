<script setup>
import {
  TooltipProvider,
  TooltipRoot,
  TooltipTrigger,
  TooltipContent,
  TooltipPortal,
  TooltipArrow,
} from 'reka-ui'
import { CircleHelp } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { Browser } from '@wailsio/runtime'

const props = defineProps({
  text: { type: String, required: true },
  link: { type: String, default: '' },
})

const { t } = useI18n()

function openLink() {
  if (props.link) {
    Browser.OpenURL(props.link)
  }
}
</script>

<template>
  <TooltipProvider :delay-duration="300">
    <TooltipRoot>
      <TooltipTrigger as-child>
        <!-- @click.prevent.stop — чтобы клик по иконке внутри <label> не переключал чекбокс -->
        <button
          type="button"
          class="inline-flex items-center justify-center shrink-0 text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300 transition-colors focus:outline-none"
          @click.prevent.stop
          tabindex="-1"
        >
          <CircleHelp :size="14" />
        </button>
      </TooltipTrigger>
      <TooltipPortal>
        <TooltipContent
          side="top"
          :side-offset="6"
          class="z-[9999] bg-gray-800 dark:bg-gray-700 text-gray-100 text-xs rounded-lg px-3 py-2 max-w-[280px] shadow-lg leading-relaxed select-none"
        >
          <p>{{ text }}</p>
          <button
            v-if="link"
            type="button"
            class="block mt-1.5 text-blue-400 hover:text-blue-300 underline transition-colors cursor-pointer text-xs"
            @click.stop="openLink"
          >
            {{ t('burn.tooltips.learnMore') }} &#x2197;
          </button>
          <TooltipArrow class="fill-gray-800 dark:fill-gray-700" />
        </TooltipContent>
      </TooltipPortal>
    </TooltipRoot>
  </TooltipProvider>
</template>
