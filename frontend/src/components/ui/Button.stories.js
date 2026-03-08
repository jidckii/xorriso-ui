import Button from './Button.vue'

export default {
  title: 'UI/Button',
  component: Button,
  argTypes: {
    variant: {
      control: 'select',
      options: ['primary', 'secondary', 'danger', 'ghost'],
    },
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg'],
    },
    disabled: { control: 'boolean' },
    loading: { control: 'boolean' },
  },
  args: {
    default: 'Button',
  },
  render: (args) => ({
    components: { Button },
    setup() {
      return { args }
    },
    template: '<Button v-bind="args">{{ args.default || "Button" }}</Button>',
  }),
}

export const Primary = {
  args: {
    variant: 'primary',
    default: 'Primary',
  },
}

export const Secondary = {
  args: {
    variant: 'secondary',
    default: 'Secondary',
  },
}

export const Danger = {
  args: {
    variant: 'danger',
    default: 'Danger',
  },
}

export const Ghost = {
  args: {
    variant: 'ghost',
    default: 'Ghost',
  },
}

export const Small = {
  args: {
    variant: 'primary',
    size: 'sm',
    default: 'Small',
  },
}

export const Medium = {
  args: {
    variant: 'primary',
    size: 'md',
    default: 'Medium',
  },
}

export const Large = {
  args: {
    variant: 'primary',
    size: 'lg',
    default: 'Large',
  },
}

export const Disabled = {
  args: {
    variant: 'primary',
    disabled: true,
    default: 'Disabled',
  },
}

export const Loading = {
  args: {
    variant: 'primary',
    loading: true,
    default: 'Loading...',
  },
}
