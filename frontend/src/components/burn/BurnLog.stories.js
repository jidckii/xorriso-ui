import BurnLog from './BurnLog.vue'

export default {
  title: 'Burn/BurnLog',
  component: BurnLog,
  argTypes: {
    lines: { control: 'object' },
  },
  decorators: [
    () => ({
      template: '<div style="padding: 2rem; max-width: 600px;"><story /></div>',
    }),
  ],
}

export const Empty = {
  args: {
    lines: [],
  },
}

export const WithEntries = {
  args: {
    lines: [
      '[INFO] Starting burn...',
      '[INFO] Speed: 8x',
      '[INFO] Writing track 1...',
    ],
  },
}

export const WithErrors = {
  args: {
    lines: [
      '[INFO] Starting...',
      '[ERROR] Write error at sector 12345',
      '[FATAL] Burn failed',
    ],
  },
}
