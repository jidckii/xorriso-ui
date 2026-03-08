import Modal from './Modal.vue'
import Button from './Button.vue'

export default {
  title: 'UI/Modal',
  component: Modal,
  parameters: {
    layout: 'fullscreen',
  },
  argTypes: {
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg', 'xl'],
    },
    show: { control: 'boolean' },
    title: { control: 'text' },
  },
  args: {
    show: true,
    title: 'Modal Title',
  },
}

export const Small = {
  args: {
    size: 'sm',
    title: 'Small Modal',
  },
  render: (args) => ({
    components: { Modal },
    setup() {
      return { args }
    },
    template: `
      <Modal v-bind="args">
        <p class="text-gray-700 dark:text-gray-300">Small modal content.</p>
      </Modal>
    `,
  }),
}

export const Medium = {
  args: {
    size: 'md',
    title: 'Medium Modal',
  },
  render: (args) => ({
    components: { Modal },
    setup() {
      return { args }
    },
    template: `
      <Modal v-bind="args">
        <p class="text-gray-700 dark:text-gray-300">Medium modal content with more text to demonstrate the default size.</p>
      </Modal>
    `,
  }),
}

export const Large = {
  args: {
    size: 'lg',
    title: 'Large Modal',
  },
  render: (args) => ({
    components: { Modal },
    setup() {
      return { args }
    },
    template: `
      <Modal v-bind="args">
        <p class="text-gray-700 dark:text-gray-300">Large modal content. This modal uses max-w-3xl for wider layouts.</p>
      </Modal>
    `,
  }),
}

export const ExtraLarge = {
  args: {
    size: 'xl',
    title: 'Extra Large Modal',
  },
  render: (args) => ({
    components: { Modal },
    setup() {
      return { args }
    },
    template: `
      <Modal v-bind="args">
        <p class="text-gray-700 dark:text-gray-300">Extra large modal content. This modal uses max-w-5xl for the widest layout option.</p>
      </Modal>
    `,
  }),
}

export const WithFooter = {
  args: {
    size: 'md',
    title: 'Modal with Footer',
  },
  render: (args) => ({
    components: { Modal, Button },
    setup() {
      return { args }
    },
    template: `
      <Modal v-bind="args">
        <p class="text-gray-700 dark:text-gray-300">This modal has a footer with action buttons.</p>
        <template #footer>
          <Button variant="ghost">Cancel</Button>
          <Button variant="primary">Save</Button>
        </template>
      </Modal>
    `,
  }),
}

export const NoFooter = {
  args: {
    size: 'md',
    title: 'Modal without Footer',
  },
  render: (args) => ({
    components: { Modal },
    setup() {
      return { args }
    },
    template: `
      <Modal v-bind="args">
        <p class="text-gray-700 dark:text-gray-300">This modal has only body content, no footer slot is provided.</p>
      </Modal>
    `,
  }),
}
