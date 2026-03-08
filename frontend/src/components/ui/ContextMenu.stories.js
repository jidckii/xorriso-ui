import ContextMenu from './ContextMenu.vue'
import { ExternalLink, FolderPlus, Info, Trash2 } from 'lucide-vue-next'
import { markRaw } from 'vue'

export default {
  title: 'UI/ContextMenu',
  component: ContextMenu,
  parameters: {
    layout: 'fullscreen',
  },
  argTypes: {
    show: { control: 'boolean' },
    x: { control: 'number' },
    y: { control: 'number' },
  },
  args: {
    show: true,
    x: 100,
    y: 100,
  },
}

export const Default = {
  args: {
    show: true,
    x: 100,
    y: 100,
    items: [
      { label: 'Open', icon: markRaw(ExternalLink), action: 'open' },
      { label: 'New Folder', icon: markRaw(FolderPlus), action: 'newFolder' },
      { label: 'Properties', icon: markRaw(Info), action: 'properties' },
      { label: 'Delete', icon: markRaw(Trash2), action: 'delete' },
    ],
  },
}

export const WithSeparator = {
  args: {
    show: true,
    x: 100,
    y: 100,
    items: [
      { label: 'Open', icon: markRaw(ExternalLink), action: 'open' },
      { label: 'New Folder', icon: markRaw(FolderPlus), action: 'newFolder' },
      { separator: true },
      { label: 'Properties', icon: markRaw(Info), action: 'properties' },
      { separator: true },
      { label: 'Delete', icon: markRaw(Trash2), action: 'delete' },
    ],
  },
}

export const WithDisabledItem = {
  args: {
    show: true,
    x: 100,
    y: 100,
    items: [
      { label: 'Open', icon: markRaw(ExternalLink), action: 'open' },
      { label: 'New Folder', icon: markRaw(FolderPlus), action: 'newFolder', disabled: true },
      { separator: true },
      { label: 'Delete', icon: markRaw(Trash2), action: 'delete', disabled: true },
    ],
  },
}
