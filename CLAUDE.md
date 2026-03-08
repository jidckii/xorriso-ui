# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# xorriso-ui

GUI для записи CD/DVD/Blu-ray/BDXL дисков. Backend — Go + xorriso subprocess, frontend — Vue3, связка через Wails3. Только Linux.

## Стек

- **Backend**: Go 1.26+, Wails3 v3.0.0-alpha.74
- **Frontend**: Vue3 (`<script setup>`), Pinia 3, Vue Router 4 (hash mode), Tailwind CSS v4, vue-i18n, reka-ui, lucide-vue-next
- **Сборка**: Task (Taskfile.yml), Yarn, wails3 CLI
- **Зависимости ОС**: xorriso 1.5.6+ (обязателен), mkisofs/cdrtools (опционален, для UDF)

## Команды разработки

```bash
wails3 dev -config ./build/config.yml    # Dev mode с hot reload
task build                                # Production сборка
go build ./...                            # Проверка компиляции Go
go test ./...                             # Запуск Go тестов
wails3 generate bindings                  # Перегенерация JS биндингов

cd frontend && yarn install               # Установка зависимостей фронтенда
cd frontend && yarn build                 # Production сборка фронтенда
cd frontend && yarn dev                   # Dev server фронтенда (без Wails)
```

## Проверка Go кода (обязательно после написания/изменения)

После любых изменений в Go коде **обязательно** выполнить:

```bash
go vet ./...                              # Статический анализ
go fix ./...                              # Автоматические исправления для новых версий Go
gofmt -w .                                # Форматирование кода
golangci-lint run ./...                   # Комплексный линтинг
```

Не коммитить код, не прошедший эти проверки.

## Тестирование Go (TDD)

- Использовать методологию TDD: **сначала писать тесты**, затем реализацию
- Стремиться к максимальному покрытию кода тестами
- Запуск тестов: `go test ./...`

## Архитектура

### Двусторонняя коммуникация Frontend ↔ Backend

```
Vue компонент → Pinia store → Wails binding (JS) → Go Service method
Go Service → app.Event.Emit() → Wails Events → Events.On() в store → reactive UI
```

- **Синхронные вызовы**: фронтенд вызывает Go-методы через автогенерированные биндинги (`frontend/bindings/`)
- **Асинхронные push-уведомления**: Go отправляет события через `application.Get().Event.Emit()`, фронтенд слушает через `Events.On()` из `@wailsio/runtime`
- Константы событий: `pkg/models/events.go`

### Структура Go

```
main.go              — Точка входа: проверка xorriso, регистрация 4 сервисов, создание окна
pkg/models/          — Структуры данных (Device, Project, BurnJob, события)
pkg/xorriso/         — Интеграция с xorriso subprocess (-pkt_output on, парсинг R:/I:/M:)
pkg/mkisofs/         — Опциональная интеграция mkisofs для UDF
services/            — 4 Wails-сервиса (device, project, burn, settings)
```

### Структура Frontend

```
frontend/src/
├── stores/          — 5 Pinia stores (device, project, burn, tab, theme)
├── views/           — 4 страницы (Project, Burn, Settings, DiscInfo)
├── components/      — layout/, device/, project/, burn/, ui/
├── locales/         — en.json, ru.json
└── i18n.js          — Интернационализация
```

Подробности: `frontend/CLAUDE.md`, `pkg/CLAUDE.md`

### xorriso subprocess

Два режима работы:
- **`executor.Run()`** — короткие операции (devices, TOC, media info): запуск → ожидание → парсинг результата
- **`executor.RunWithProgress()`** — долгие операции (burn, blank, format): streaming парсинг UPDATE строк, push через Wails Events

xorriso запускается с `-pkt_output on` для машинночитаемого вывода:
- `R:` — результат команды → `CmdResult.ResultLines`
- `I:` — информация/ошибки → `CmdResult.InfoLines`
- `M:` — маркеры

## Wails3 API (критично)

- Доступ к приложению: `application.Get()` возвращает `*App`
- Отправка событий: `app.Event.Emit(name, data...)` — поле `Event`, **НЕ** `Events`
- Регистрация сервисов: `application.NewService(instance)` — generic функция
- Lifecycle: `ServiceStartup(ctx context.Context, options application.ServiceOptions) error`
- **НЕ существует**: `application.WailsEvent`, `options.Application`, `s.app.Events`
- Frontend events: `import { Events } from '@wailsio/runtime'` → `Events.On()` / `Events.Off()`

## Использование агентов (обязательно!)

Для любых крупных задач **обязательно** делегируй работу специализированным агентам:

- **researcher** (haiku) — исследование кода, поиск по проекту, изучение библиотек и документации в интернете
- **golang-architect** (opus) — разработка, рефакторинг, ревью Go-кода, работа с сервисами и xorriso
- **frontend-architect** (sonnet) — разработка UI на Vue3 + Tailwind, Pinia stores, E2E тесты через Playwright
- **ci-engineer** (sonnet) — CI/CD пайплайны, Taskfile, Docker, пакетирование, автоматизация сборки

**Когда использовать агентов:**
- Планирование и реализация новых фич
- Рефакторинг (любого масштаба)
- Редактирование кода по плану (больше 2-3 файлов)
- Исследование архитектуры и зависимостей
- Настройка CI/CD и автоматизации
- E2E и интеграционное тестирование

**Параллелизация**: если задача затрагивает и фронтенд, и бэкенд — запускай `golang-architect` и `frontend-architect` параллельно.

## Соглашения

- Только Linux, без кроссплатформенности
- Go пакеты в `pkg/`, НЕ в `internal/`
- Фронтенд: `yarn`, НЕ `npm`
- Go модуль: `xorriso-ui`
- Тёмная тема по умолчанию (bg-gray-900, класс `dark` на `<html>`)
- Формат проектов: `.xorriso-project` (JSON)
- Интерфейс многовкладочный (tabStore)
- Интернационализация: en/ru, ключи в `frontend/src/locales/`
- UI иконки: lucide-vue-next; файловые иконки: material-file-icons
- Headless UI компоненты: reka-ui (TreeRoot/TreeItem и т.д.)
