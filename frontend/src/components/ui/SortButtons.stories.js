import SortButtons from './SortButtons.vue'

export default {
  title: 'UI/SortButtons',
  component: SortButtons,
}

export const SortByName = {
  args: { sortBy: 'name', sortDir: 'asc' },
}

export const SortBySize = {
  args: { sortBy: 'size', sortDir: 'desc' },
}

export const SortByDate = {
  args: { sortBy: 'date', sortDir: 'asc' },
}

export const SortByType = {
  args: { sortBy: 'type', sortDir: 'asc' },
}
