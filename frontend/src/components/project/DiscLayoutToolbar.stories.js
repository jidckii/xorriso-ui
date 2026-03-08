import DiscLayoutToolbar from './DiscLayoutToolbar.vue'

export default {
  title: 'Project/DiscLayoutToolbar',
  component: DiscLayoutToolbar,
  render: (args) => ({
    components: { DiscLayoutToolbar },
    setup() {
      return { args }
    },
    template: '<div style="width: 500px;"><DiscLayoutToolbar v-bind="args" /></div>',
  }),
}

export const NothingSelected = {
  args: {
    allSelected: false,
    selectedCount: 0,
    canBurn: false,
  },
}

export const WithSelection = {
  args: {
    allSelected: false,
    selectedCount: 3,
    canBurn: true,
  },
}

export const EmptyProject = {
  args: {
    allSelected: false,
    selectedCount: 0,
    canBurn: false,
  },
}
