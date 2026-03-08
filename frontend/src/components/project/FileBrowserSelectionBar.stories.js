import FileBrowserSelectionBar from './FileBrowserSelectionBar.vue'

export default {
  title: 'Project/FileBrowserSelectionBar',
  component: FileBrowserSelectionBar,
  render: (args) => ({
    components: { FileBrowserSelectionBar },
    setup() {
      return { args }
    },
    template: '<div style="width: 400px;"><FileBrowserSelectionBar v-bind="args" /></div>',
  }),
}

export const NothingSelected = {
  args: {
    allSelected: false,
    selectedCount: 0,
  },
}

export const PartiallySelected = {
  args: {
    allSelected: false,
    selectedCount: 3,
  },
}

export const AllSelected = {
  args: {
    allSelected: true,
    selectedCount: 10,
  },
}
