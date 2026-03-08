import FileIcon from './FileIcon.vue'

export default {
  title: 'UI/FileIcon',
  component: FileIcon,
  argTypes: {
    name: { control: 'text' },
    isDir: { control: 'boolean' },
    isOpen: { control: 'boolean' },
    size: { control: { type: 'range', min: 8, max: 48 } },
  },
  decorators: [
    () => ({
      template: '<div style="padding: 2rem; display: flex; align-items: center; gap: 0.5rem;"><story /></div>',
    }),
  ],
}

export const JavaScriptFile = {
  args: {
    name: 'app.js',
  },
}

export const ImageFile = {
  args: {
    name: 'photo.png',
  },
}

export const ZipFile = {
  args: {
    name: 'archive.zip',
  },
}

export const ClosedFolder = {
  args: {
    name: 'src',
    isDir: true,
  },
}

export const OpenFolder = {
  args: {
    name: 'src',
    isDir: true,
    isOpen: true,
  },
}

export const LargeSize = {
  args: {
    name: 'document.pdf',
    size: 32,
  },
}

export const SmallSize = {
  args: {
    name: 'readme.md',
    size: 12,
  },
}
