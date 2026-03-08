import ProgressBar from './ProgressBar.vue'

export default {
  title: 'UI/ProgressBar',
  component: ProgressBar,
  argTypes: {
    value: { control: { type: 'range', min: 0, max: 100 } },
    variant: {
      control: 'select',
      options: ['default', 'success', 'warning', 'danger'],
    },
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg'],
    },
    showLabel: { control: 'boolean' },
  },
  decorators: [
    () => ({
      template: '<div style="padding: 2rem; max-width: 400px;"><story /></div>',
    }),
  ],
}

export const Empty = {
  args: {
    value: 0,
  },
}

export const HalfFull = {
  args: {
    value: 50,
  },
}

export const Complete = {
  args: {
    value: 100,
  },
}

export const Success = {
  args: {
    value: 75,
    variant: 'success',
  },
}

export const Warning = {
  args: {
    value: 60,
    variant: 'warning',
  },
}

export const Danger = {
  args: {
    value: 90,
    variant: 'danger',
  },
}

export const WithLabel = {
  args: {
    value: 65,
    showLabel: true,
  },
}

export const Small = {
  args: {
    value: 50,
    size: 'sm',
  },
}

export const Large = {
  args: {
    value: 50,
    size: 'lg',
  },
}
