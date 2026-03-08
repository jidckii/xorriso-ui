import { ref } from 'vue'
import BurnOptions from './BurnOptions.vue'

export default {
  title: 'Burn/BurnOptions',
  component: BurnOptions,
  argTypes: {
    availableSpeeds: { control: 'object' },
  },
  decorators: [
    () => ({
      template: '<div style="padding: 2rem; max-width: 500px;"><story /></div>',
    }),
  ],
}

export const Default = {
  args: {
    availableSpeeds: [1, 2, 4, 8, 16],
  },
}

export const AllEnabled = {
  render: () => ({
    components: { BurnOptions },
    setup() {
      const optionsRef = ref(null)

      const onMounted = () => {
        if (optionsRef.value) {
          const opts = optionsRef.value.options
          opts.verify = true
          opts.dummyMode = true
          opts.finalize = true
          opts.ejectAfter = true
          opts.streamRecording = true
        }
      }

      return { optionsRef, onMounted }
    },
    template: `
      <BurnOptions
        ref="optionsRef"
        :available-speeds="[1, 2, 4, 8, 16]"
        @vue:mounted="onMounted"
      />
    `,
  }),
}
