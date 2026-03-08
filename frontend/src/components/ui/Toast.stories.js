import Toast from './Toast.vue'

export default {
  title: 'UI/Toast',
  component: Toast,
  argTypes: {
    type: {
      control: 'select',
      options: ['info', 'success', 'warning', 'error'],
    },
    message: { control: 'text' },
    duration: { control: 'number' },
  },
  args: {
    duration: 0,
  },
  decorators: [
    () => ({
      template: '<div style="padding: 2rem;"><story /></div>',
    }),
  ],
}

export const Info = {
  args: {
    type: 'info',
    message: 'This is an informational message.',
    duration: 0,
  },
}

export const Success = {
  args: {
    type: 'success',
    message: 'Operation completed successfully!',
    duration: 0,
  },
}

export const Warning = {
  args: {
    type: 'warning',
    message: 'Please check your settings before proceeding.',
    duration: 0,
  },
}

export const Error = {
  args: {
    type: 'error',
    message: 'An error occurred while burning the disc.',
    duration: 0,
  },
}

export const AutoDismiss = {
  args: {
    type: 'info',
    message: 'This toast will NOT auto-dismiss (duration=0 for Storybook).',
    duration: 0,
  },
}
