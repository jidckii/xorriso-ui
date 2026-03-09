import InfoTooltip from './InfoTooltip.vue'

export default {
  title: 'UI/InfoTooltip',
  component: InfoTooltip,
  parameters: {
    layout: 'centered',
  },
}

// Только текст — базовое использование
export const OnlyText = {
  render: () => ({
    components: { InfoTooltip },
    template: `
      <div class="flex items-center gap-2 text-sm text-gray-300 p-8">
        <span>ISO Level</span>
        <InfoTooltip text="Определяет максимальную длину имён файлов и размер. Уровень 1: имена 8.3, лимит 4ГБ. Уровень 3: длинные имена, без ограничений." />
      </div>
    `,
  }),
}

// Текст со ссылкой
export const WithLink = {
  render: () => ({
    components: { InfoTooltip },
    template: `
      <div class="flex items-center gap-2 text-sm text-gray-300 p-8">
        <span>UDF</span>
        <InfoTooltip
          text="Universal Disk Format — современная файловая система для DVD/BD. Поддерживает файлы >4ГБ, длинные имена и Unicode."
          link="https://en.wikipedia.org/wiki/Universal_Disk_Format"
        />
      </div>
    `,
  }),
}

// Длинный текст
export const LongText = {
  render: () => ({
    components: { InfoTooltip },
    template: `
      <div class="flex items-center gap-2 text-sm text-gray-300 p-8">
        <span>Rock Ridge</span>
        <InfoTooltip
          text="POSIX-расширение ISO 9660. Сохраняет права Unix, симлинки, длинные имена и глубокую вложенность директорий. Необходим для корректного восстановления файловой системы Linux с диска."
          link="https://en.wikipedia.org/wiki/Rock_Ridge"
        />
      </div>
    `,
  }),
}

// Внутри label с чекбоксом — проверка что клик по иконке не переключает чекбокс
export const InsideLabel = {
  render: () => ({
    components: { InfoTooltip },
    data() {
      return { checked: false }
    },
    template: `
      <div class="flex flex-col gap-4 p-8">
        <label class="flex items-center gap-2 text-sm text-gray-300 cursor-pointer">
          <input type="checkbox" v-model="checked" class="accent-blue-500" />
          UDF
          <InfoTooltip
            text="Universal Disk Format — современная файловая система для DVD/BD."
            link="https://en.wikipedia.org/wiki/Universal_Disk_Format"
          />
        </label>
        <p class="text-xs text-gray-500">Состояние чекбокса: {{ checked ? 'включён' : 'выключен' }}</p>
      </div>
    `,
  }),
}
