import { ref } from 'vue'
import DiscLayoutTree from './DiscLayoutTree.vue'

export default {
  title: 'Project/DiscLayoutTree',
  component: DiscLayoutTree,
}

export const FlatFiles = {
  render: () => ({
    components: { DiscLayoutTree },
    setup() {
      const items = [
        { _key: '/file1.txt', name: 'file1.txt', isDir: false, size: 1024, children: [] },
        { _key: '/photo.jpg', name: 'photo.jpg', isDir: false, size: 3145728, children: [] },
        { _key: '/backup.tar.gz', name: 'backup.tar.gz', isDir: false, size: 52428800, children: [] },
      ]
      const expanded = ref([])
      const selectedKeys = new Set()
      return { items, expanded, selectedKeys }
    },
    template: '<div style="width: 400px;"><DiscLayoutTree :items="items" v-model:expanded="expanded" :selected-keys="selectedKeys" /></div>',
  }),
}

export const NestedFolders = {
  render: () => ({
    components: { DiscLayoutTree },
    setup() {
      const items = [
        {
          _key: '/Documents',
          name: 'Documents',
          isDir: true,
          size: 0,
          children: [
            { _key: '/Documents/report.pdf', name: 'report.pdf', isDir: false, size: 2097152, children: [] },
            { _key: '/Documents/notes.txt', name: 'notes.txt', isDir: false, size: 512, children: [] },
          ],
        },
        {
          _key: '/Photos',
          name: 'Photos',
          isDir: true,
          size: 0,
          children: [
            { _key: '/Photos/vacation.jpg', name: 'vacation.jpg', isDir: false, size: 4194304, children: [] },
            { _key: '/Photos/family.png', name: 'family.png', isDir: false, size: 8388608, children: [] },
          ],
        },
        { _key: '/readme.md', name: 'readme.md', isDir: false, size: 256, children: [] },
      ]
      const expanded = ref(['/Documents', '/Photos'])
      const selectedKeys = new Set()
      return { items, expanded, selectedKeys }
    },
    template: '<div style="width: 400px;"><DiscLayoutTree :items="items" v-model:expanded="expanded" :selected-keys="selectedKeys" /></div>',
  }),
}

export const WithSelection = {
  render: () => ({
    components: { DiscLayoutTree },
    setup() {
      const items = [
        { _key: '/file1.txt', name: 'file1.txt', isDir: false, size: 1024, children: [] },
        { _key: '/file2.txt', name: 'file2.txt', isDir: false, size: 2048, children: [] },
        { _key: '/file3.txt', name: 'file3.txt', isDir: false, size: 4096, children: [] },
      ]
      const expanded = ref([])
      const selectedKeys = new Set(['/file1.txt', '/file3.txt'])
      return { items, expanded, selectedKeys }
    },
    template: '<div style="width: 400px;"><DiscLayoutTree :items="items" v-model:expanded="expanded" :selected-keys="selectedKeys" /></div>',
  }),
}
