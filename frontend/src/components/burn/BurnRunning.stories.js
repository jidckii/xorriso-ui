import BurnRunning from './BurnRunning.vue'

export default {
  title: 'Burn/BurnRunning',
  component: BurnRunning,
  decorators: [
    () => ({
      template: '<div style="max-width: 600px;"><story /></div>',
    }),
  ],
}

export const Start = {
  args: {
    progress: {
      percent: 5,
      speed: '',
      bytesWritten: 0,
      bytesTotal: 4700000000,
      eta: '',
      fifoFill: 0,
    },
    logLines: [],
    phaseLabel: 'Burning',
  },
}

export const Middle = {
  args: {
    progress: {
      percent: 50,
      speed: '8x',
      bytesWritten: 2350000000,
      bytesTotal: 4700000000,
      eta: '2:30',
      fifoFill: 85,
    },
    logLines: [
      '[INFO] Starting burn process...',
      '[INFO] Speed: 8x',
      '[INFO] Writing track 1...',
    ],
    phaseLabel: 'Burning',
  },
}

export const WithLog = {
  render: () => ({
    components: { BurnRunning },
    setup() {
      const progress = {
        percent: 72,
        speed: '8x',
        bytesWritten: 3384000000,
        bytesTotal: 4700000000,
        eta: '1:15',
        fifoFill: 90,
      }

      const logLines = [
        '[INFO] Starting burn process...',
        '[INFO] Speed: 8x',
        '[INFO] Mode: auto',
        '[INFO] Verify: Yes',
        '[INFO] Writing track 1...',
        '[INFO] 25% complete',
        '[INFO] 50% complete',
        '[INFO] 72% complete',
      ]

      return { progress, logLines }
    },
    template: `
      <BurnRunning
        :progress="progress"
        :log-lines="logLines"
        phase-label="Burning"
      />
    `,
  }),
}
