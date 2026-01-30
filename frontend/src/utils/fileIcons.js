import { getIcon } from 'material-file-icons'

const iconCache = new Map()

/**
 * Returns an SVG string for a file icon based on its filename.
 * Uses material-file-icons for VS Code-style file type icons.
 * Results are cached by file extension for performance.
 */
export function getFileIconSvg(filename) {
  const ext = filename.includes('.') ? filename.split('.').pop().toLowerCase() : ''
  const cacheKey = ext || filename.toLowerCase()

  if (iconCache.has(cacheKey)) {
    return iconCache.get(cacheKey)
  }

  const icon = getIcon(filename)
  iconCache.set(cacheKey, icon.svg)
  return icon.svg
}
