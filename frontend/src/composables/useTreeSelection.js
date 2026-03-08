import { ref, computed } from 'vue'

/**
 * Composable for tree selection with parent-child propagation.
 * @param {Object} options
 * @param {(item: any) => string} options.getKey - extract unique key from item
 * @param {(item: any) => any[]|undefined} options.getChildren - get children array
 */
export function useTreeSelection({ getKey, getChildren }) {
  const selectedKeys = ref(new Set())

  function isSelected(item) {
    return selectedKeys.value.has(getKey(item))
  }

  function toggleSelection(item) {
    const key = getKey(item)
    const selecting = !selectedKeys.value.has(key)

    if (selecting) {
      selectedKeys.value.add(key)
    } else {
      selectedKeys.value.delete(key)
    }

    const children = getChildren(item)
    if (children && children.length > 0) {
      propagateSelection(children, selecting)
    }

    selectedKeys.value = new Set(selectedKeys.value)
  }

  function propagateSelection(children, selecting) {
    for (const child of children) {
      if (selecting) {
        selectedKeys.value.add(getKey(child))
      } else {
        selectedKeys.value.delete(getKey(child))
      }
      const grandchildren = getChildren(child)
      if (grandchildren && grandchildren.length > 0) {
        propagateSelection(grandchildren, selecting)
      }
    }
  }

  function selectAll(items) {
    selectAllRecursive(items)
    selectedKeys.value = new Set(selectedKeys.value)
  }

  function selectAllRecursive(items) {
    for (const item of items) {
      selectedKeys.value.add(getKey(item))
      const children = getChildren(item)
      if (children && children.length > 0) {
        selectAllRecursive(children)
      }
    }
  }

  function deselectAll() {
    selectedKeys.value = new Set()
  }

  function countAll(items) {
    let count = 0
    for (const item of items) {
      count++
      const children = getChildren(item)
      if (children && children.length > 0) {
        count += countAll(children)
      }
    }
    return count
  }

  function allSelected(items) {
    if (!items || items.length === 0) return false
    return selectedKeys.value.size > 0 && countAll(items) === selectedKeys.value.size
  }

  const selectedCount = computed(() => selectedKeys.value.size)

  return {
    selectedKeys,
    isSelected,
    toggleSelection,
    selectAll,
    deselectAll,
    allSelected,
    selectedCount,
  }
}
