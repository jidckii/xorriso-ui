import BurnProgress from './BurnProgress.vue'

export default {
  title: 'Burn/BurnProgress',
  component: BurnProgress,
  argTypes: {
    percent: { control: { type: 'range', min: 0, max: 100 } },
    phase: { control: 'select', options: ['Writing', 'Formatting', 'Verifying'] },
    speed: { control: 'text' },
    fifoPercent: { control: { type: 'range', min: 0, max: 100 } },
    eta: { control: 'text' },
    bytesWritten: { control: 'number' },
    bytesTotal: { control: 'number' },
  },
  decorators: [
    () => ({
      template: '<div style="padding: 2rem; max-width: 500px;"><story /></div>',
    }),
  ],
}

export const Start = {
  args: {
    percent: 0,
    phase: 'Writing',
  },
}

export const HalfWay = {
  args: {
    percent: 50,
    phase: 'Writing',
    speed: '8.2x',
    fifoPercent: 85,
    eta: '2:30',
    bytesWritten: 2.35e9,
    bytesTotal: 4.7e9,
  },
}

export const NearComplete = {
  args: {
    percent: 95,
    phase: 'Writing',
    speed: '8.2x',
    fifoPercent: 95,
    eta: '0:15',
  },
}

export const Verifying = {
  args: {
    percent: 60,
    phase: 'Verifying',
    speed: '16x',
    fifoPercent: 0,
    eta: '1:00',
  },
}
