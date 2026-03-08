import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

/**
 * Composable for media status display (color, dot, label).
 * @param {import('vue').Ref<string>|() => string} statusGetter - reactive status value or getter
 */
export function useMediaStatus(statusGetter) {
  const { t } = useI18n()

  const status = computed(() =>
    typeof statusGetter === 'function' ? statusGetter() : statusGetter.value
  )

  const statusColor = computed(() => {
    const map = {
      blank: 'text-green-400',
      appendable: 'text-yellow-400',
      closed: 'text-red-400',
      unknown: 'text-gray-600 dark:text-gray-400',
    }
    return map[status.value] || 'text-gray-600 dark:text-gray-400'
  })

  const statusDot = computed(() => {
    const map = {
      blank: 'bg-green-500',
      appendable: 'bg-yellow-500',
      closed: 'bg-red-500',
      unknown: 'bg-gray-500',
    }
    return map[status.value] || 'bg-gray-500'
  })

  const statusLabel = computed(() => {
    const map = {
      blank: t('device.blank'),
      appendable: t('device.appendable'),
      closed: t('device.closed'),
      unknown: t('device.unknown'),
    }
    return map[status.value] || status.value
  })

  return { statusColor, statusDot, statusLabel }
}
