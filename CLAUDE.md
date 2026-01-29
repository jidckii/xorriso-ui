# xorriso-ui

GUI приложение для записи CD / DVD / Blu-ray / BDXL дисков с современным интерфейсом.
Использует xorriso как бэкенд для работы с дисками и ISO-образами.
Построено на Wails3 (Vue3 + Go), только Linux.

## Стек технологий

- **Backend**: Go 1.25+, Wails3 v3.0.0-alpha.24 (команда `wails3`)
- **Frontend**: Vue3 (`<script setup>`), Pinia, Vue Router (hash mode), Tailwind CSS v4
- **Инструменты**: asdf (версии Go/Node), Yarn (пакеты фронтенда), Task (сборка)
- **Зависимость**: xorriso 1.5.6+ (`/usr/bin/xorriso`)

## Разработка

```bash
# Запуск в режиме разработки
wails3 dev -config ./build/config.yml

# Сборка фронтенда
cd frontend && yarn build

# Сборка Go
go build ./...

# Генерация биндингов Wails
wails3 generate bindings

# Сборка продакшен
task build
```

## Структура проекта

```
main.go                    # Точка входа, регистрация Wails-сервисов
Taskfile.yml               # Task runner (только linux)
go.mod                     # Модуль: xorriso-ui

pkg/
├── xorriso/
│   ├── executor.go        # Запуск xorriso subprocess с -pkt_output on
│   ├── parser.go          # Парсинг R:/I:/M: каналов, устройств, скоростей
│   ├── commands.go        # Fluent builder командных строк xorriso
│   └── progress.go        # Парсинг UPDATE/pacifier строк прогресса
├── models/
│   ├── device.go          # Device, MediaInfo, SpeedDescriptor
│   ├── project.go         # Project, FileEntry, ISOOptions, BurnOptions
│   ├── burn.go            # BurnState, BurnJob, BurnProgress, BurnResult
│   └── events.go          # Константы имён Wails-событий

services/
├── device_service.go      # ListDevices, GetMediaInfo, GetSpeeds, EjectDisc, polling
├── project_service.go     # NewProject, Save/Open, AddFiles, RemoveEntry, BrowseDirectory
├── burn_service.go        # StartBurn, CancelBurn, BlankDisc, FormatDisc
└── settings_service.go    # GetSettings, SaveSettings (~/.config/xorriso-ui/settings.json)

frontend/src/
├── main.js                # Vue3 + Pinia + Router
├── App.vue                # Корневой: AppHeader + router-view + AppStatusBar
├── stores/
│   ├── deviceStore.js     # Приводы, текущий привод, медиа-инфо
│   ├── projectStore.js    # Файлы проекта, ISO-опции
│   └── burnStore.js       # Состояние записи, прогресс
├── views/
│   ├── ProjectView.vue    # Главный экран (split: FileBrowser + DiscLayout)
│   ├── BurnView.vue       # Диалог записи с прогрессом
│   └── SettingsView.vue   # Настройки
└── components/
    ├── layout/            # AppHeader, AppSidebar, AppStatusBar
    ├── device/            # DeviceSelector, DevicePanel, MediaInfo, SpeedSelector
    ├── project/           # FileBrowser, DiscLayout, FileTree, CapacityBar, DragDropZone
    ├── burn/              # BurnDialog, BurnOptions, BurnProgress, BurnLog
    └── ui/                # Button, Modal, ProgressBar, Toast, DropdownMenu
```

## Архитектура xorriso-интеграции

xorriso запускается как subprocess с `-pkt_output on` для машинночитаемого вывода:
- `R:` — результат команды (парсится в `ResultLines`)
- `I:` — информация/ошибки (парсится в `InfoLines`)
- `M:` — маркеры

Два режима работы:
- **`executor.Run()`** — короткие операции (devices, toc, media info): запуск → ожидание → парсинг
- **`executor.RunWithProgress()`** — долгие операции (burn, blank, format): парсинг UPDATE строк в реальном времени, трансляция через Wails Events

## Wails3 API (важно)

- Доступ к приложению из сервисов: `application.Get()` возвращает `*App`
- Отправка событий: `app.Event.Emit(name, data...)` (поле `Event`, НЕ `Events`)
- Регистрация сервисов: `application.NewService(instance)` — generic функция
- Lifecycle сервисов: `ServiceStartup(ctx context.Context, options application.ServiceOptions) error`
- **НЕ существует**: `application.WailsEvent`, `options.Application`, `s.app.Events`

## Соглашения

- Только Linux, без кроссплатформенности
- Пакеты в `pkg/`, НЕ в `internal/`
- Тёмная тема по умолчанию (bg-gray-900)
- Фронтенд собирается через `yarn`, НЕ `npm`
- Go модуль: `xorriso-ui`
- Формат проектов: `.xorriso-project` (JSON)
