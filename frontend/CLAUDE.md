# CLAUDE.md — Frontend

This file provides guidance to Claude Code for frontend development in xorriso-ui.

## Стек

- Vue3 с `<script setup>` синтаксисом (Composition API only)
- Pinia 3 — state management
- Vue Router 4 — hash mode (`createWebHashHistory`)
- Tailwind CSS v4 — через `@tailwindcss/vite` plugin (НЕ PostCSS)
- vue-i18n 12 alpha — интернационализация (en/ru)
- reka-ui — headless UI (Tree, Dialog и т.д.)
- lucide-vue-next — иконки интерфейса
- material-file-icons — иконки файлов (VS Code стиль)
- Vite 5 — сборщик
- @wailsio/runtime 3.0.0-alpha.79 — коммуникация с Go backend

## Команды

```bash
yarn install          # Установка зависимостей
yarn dev              # Vite dev server (:5173)
yarn build            # Production сборка → dist/
yarn build:dev        # Debug сборка без минификации
```

## Маршруты

```
/          → ProjectView  (основной экран)
/burn      → BurnView     (диалог записи)
/settings  → SettingsView (настройки)
```

DiscInfoView — отображается как overlay поверх ProjectView (управляется через tabStore.showDiscInfo).

## Pinia Stores

### deviceStore
- Управляет списком приводов, текущим устройством, медиа-информацией, скоростями
- Слушает backend-события `device:list-updated`, `device:media-changed`
- Ключевые getters: `currentDevice`, `hasMedia`, `mediaCapacityBytes`, `mediaFreeBytes`

### projectStore
- Файлы проекта, ISO-опции, размер
- Синхронизируется с tabStore: `updateTabProjectData()` при каждом изменении
- Actions: `newProject()`, `openProject()`, `saveProject()`, `addFiles()`, `removeEntry()`
- `browseDirectory()` / `getImagePreview()` — для FileBrowser

### burnStore
- Состояние записи, прогресс, логи
- Слушает: `burn:progress`, `burn:state-changed`, `burn:complete`, `burn:error`
- Actions: `startBurn()`, `cancelBurn()`, `blankDisc()`, `formatDisc()`

### tabStore
- Multi-document interface: каждый проект — вкладка
- State: `tabs[]`, `activeTabId`, `showDiscInfo`
- Каждая вкладка хранит свой `projectData` (entries, ISO options, browse state)

### themeStore
- Переключение light/dark; синхронизирует localStorage ↔ Go settings ↔ CSS класс `dark`

## Вызов Go backend

```javascript
// Автогенерированные биндинги
import { ListDevices } from '../bindings/xorriso-ui/services/deviceservice'
const devices = await ListDevices()

// Слушать события от backend
import { Events } from '@wailsio/runtime'
Events.On('burn:progress', (data) => { ... })
```

Биндинги генерируются командой `wails3 generate bindings` и лежат в `frontend/bindings/`.

## Компоненты

### layout/
- `AppHeader.vue` — верхнее меню (File: New/Open/Save, Help, Theme toggle)
- `AppStatusBar.vue` — нижняя строка состояния (текущий привод, ошибки)
- `TabBar.vue` — вкладки проектов с "+" для нового

### device/
- `DeviceSelector.vue` — dropdown выбора привода
- `DevicePanel.vue` — панель информации о приводе
- `MediaInfo.vue` — ёмкость, профиль, сессии, кнопки Eject/Load

### project/
- `FileBrowser.vue` — левая панель: навигация по файловой системе хоста
- `FileBrowserItem.vue` — элемент в FileBrowser
- `DiscLayout.vue` — правая панель: структура файлов на будущем диске (drag-drop добавление)
- `FileTree.vue` — переиспользуемое дерево файлов (reka-ui TreeRoot/TreeItem)
- `CapacityBar.vue` — визуальный индикатор заполнения диска (green → yellow → red)
- `FilePropertiesModal.vue` — свойства файла

### burn/
- `BurnDialog.vue` — основной диалог записи
- `BurnOptions.vue` — настройки записи (скорость, verify, close disc, etc.)
- `BurnProgress.vue` — прогресс бар с фазами, ETA, FIFO
- `BurnLog.vue` — лог операций

### ui/
Переиспользуемые: `Button`, `Modal`, `ProgressBar`, `Toast`, `ContextMenu`, `FileIcon`, `ImagePreviewTooltip`, `PanelHeader`

## Принципы разработки

### Компонентный подход (обязательно)

- Весь UI строится из изолированных, переиспользуемых компонентов
- Логика не дублируется — выносится в отдельный компонент или composable
- Каждый компонент отвечает за одну задачу (Single Responsibility)
- Props для входных данных, emit для событий наружу — никаких прямых мутаций родительского состояния

### Storybook (обязательно)

- Для **каждого** компонента должен быть создан story-файл
- Stories позволяют визуально проверить компонент в изоляции до интеграции в приложение
- Файлы stories располагаются рядом с компонентами: `ComponentName.stories.js`
- Каждый story должен покрывать основные состояния компонента (default, loading, error, empty, disabled и т.д.)
- При создании нового компонента — story создаётся одновременно с ним

## CSS / Темизация

```css
/* main.css */
@import "tailwindcss";
@custom-variant dark (&:where(.dark, .dark *));
```

- Tailwind v4: без tailwind.config.js, конфигурация через CSS
- Тёмная тема: класс `dark` на `<html>`, переключается через themeStore
- Цвета по умолчанию: bg-gray-900 (тёмная), bg-gray-50 (светлая)

## i18n

- Файлы: `src/locales/en.json`, `src/locales/ru.json`
- Использование: `{{ $t('key.path') }}` в шаблонах, `t('key.path')` в setup
- Fallback locale: en
- Locale хранится в localStorage + синхронизируется с Go settings
