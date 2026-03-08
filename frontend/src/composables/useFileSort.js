import { ref, computed } from 'vue'

export function useFileSort(items) {
  const sortBy = ref('name') // name | size | date | type
  const sortDir = ref('asc') // asc | desc

  function getExtension(name) {
    const dot = name.lastIndexOf('.')
    return dot > 0 ? name.slice(dot + 1).toLowerCase() : ''
  }

  function compareFn(a, b) {
    // Папки всегда первые
    if (a.isDir !== b.isDir) return a.isDir ? -1 : 1

    let result = 0
    switch (sortBy.value) {
      case 'name':
        result = a.name.localeCompare(b.name, undefined, { numeric: true, sensitivity: 'base' })
        break
      case 'size':
        result = (a.size || 0) - (b.size || 0)
        break
      case 'date':
        result = (a.modTime || 0) - (b.modTime || 0)
        break
      case 'type':
        result = getExtension(a.name).localeCompare(getExtension(b.name))
        if (result === 0) result = a.name.localeCompare(b.name, undefined, { numeric: true, sensitivity: 'base' })
        break
    }
    return sortDir.value === 'asc' ? result : -result
  }

  const sorted = computed(() => {
    if (!items.value) return []
    return [...items.value].sort(compareFn)
  })

  function toggleSort(field) {
    if (sortBy.value === field) {
      sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
    } else {
      sortBy.value = field
      sortDir.value = 'asc'
    }
  }

  function sortChildren(children) {
    if (!children || !children.length) return children
    return [...children].sort(compareFn)
  }

  return { sortBy, sortDir, sorted, toggleSort, compareFn, sortChildren }
}
