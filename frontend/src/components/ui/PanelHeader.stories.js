import PanelHeader from './PanelHeader.vue'
import Button from './Button.vue'
import { h } from 'vue'

export default {
  title: 'UI/PanelHeader',
  component: PanelHeader,
  decorators: [
    () => ({
      template: '<div style="max-width: 600px;"><story /></div>',
    }),
  ],
}

export const Default = {
  render: () => ({
    components: { PanelHeader },
    template: `
      <PanelHeader>
        <template #row1>
          <span class="text-sm text-gray-700 dark:text-gray-300">Toolbar row</span>
        </template>
        <template #row2>
          <span class="text-sm text-gray-600 dark:text-gray-400">Info row</span>
        </template>
      </PanelHeader>
    `,
  }),
}

export const WithButtons = {
  render: () => ({
    components: { PanelHeader, Button },
    template: `
      <PanelHeader>
        <template #row1>
          <Button size="sm" variant="ghost">New</Button>
          <Button size="sm" variant="ghost">Open</Button>
          <Button size="sm" variant="ghost">Delete</Button>
        </template>
        <template #row2>
          <span class="text-xs text-gray-500 dark:text-gray-400">/ home / user / projects</span>
        </template>
      </PanelHeader>
    `,
  }),
}
