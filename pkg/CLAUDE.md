# CLAUDE.md — Backend (Go)

This file provides guidance to Claude Code for Go backend development in xorriso-ui.

## Пакеты

### pkg/models/

Структуры данных, разделяемые между сервисами и xorriso-пакетом:

- **device.go**: `Device` (path, vendor, model, profiles), `MediaInfo` (type, capacity, sessions, speeds), `SpeedDescriptor`, `MediaProfile`, `Session`, `TOC`
- **project.go**: `Project` (name, volumeID, entries, ISO/burn options), `FileEntry` (рекурсивное дерево: sourcePath, destPath, isDir, size, children), `ISOOptions`, `BurnOptions`
- **burn.go**: `BurnState` (enum: pending → writing → verifying → done/error/cancelled), `BurnJob`, `BurnProgress` (phase, percent, speed, ETA, FIFO), `BurnResult`
- **events.go**: Константы имён Wails-событий (`device:list-updated`, `burn:progress`, и т.д.)

### pkg/xorriso/

Интеграция с xorriso через subprocess:

- **executor.go**: `Executor` — запуск xorriso с `-pkt_output on`
  - `Run(ctx, args...)` — синхронное выполнение для коротких операций
  - `RunWithProgress(ctx, progressFn, args...)` — streaming для долгих операций
- **parser.go**: Парсинг pkt_output формата (R:/I:/M: каналы)
  - `ParsePktLine()`, `ParsePktOutput()` → `CmdResult{ResultLines, InfoLines}`
  - `ParseDevices()`, `ParseMediaInfo()`, `ParseSpeeds()`, `ParseProfiles()`, `ParseTOCSessions()`
- **commands.go**: `CommandBuilder` — fluent API для построения аргументов xorriso
  - Цепочка: `NewCommandBuilder().Device(path).VolumeID(id).Map(src, dst).Commit()`
- **progress.go**: Парсинг UPDATE строк из stderr xorriso → `Progress{Percent, Speed, Written, ETA}`

### pkg/mkisofs/

Опциональная интеграция mkisofs для создания ISO с UDF:

- **executor.go**: `Executor` с `BuildISO(ctx, opts, progressFn)` — создаёт ISO, парсит `N% done` из stderr
- `FileMappingsFromEntries()` — конвертирует `[]FileEntry` в graft-points для mkisofs

## Сервисы (services/)

Каждый сервис реализует Wails3 lifecycle: `ServiceName()`, `ServiceStartup()`, `ServiceShutdown()`.

### DeviceService

- `ListDevices()` — сканирует `/proc/sys/dev/cdrom/info` + `/sys/block/` для обнаружения приводов
- `GetMediaInfo(path)`, `GetSpeeds(path)`, `GetDriveProfiles(path)` — через xorriso
- `EjectDisc()`, `LoadTray()` — управление лотком
- `pollDevices()` — фоновый polling каждые 2с, эмитирует `device:list-updated`, `device:media-changed`

### ProjectService

- `NewProject()`, `SaveProject()`, `SaveProjectAs()`, `OpenProject()` — CRUD проектов (JSON .xorriso-project)
- `AddFiles(paths, destDir)` — добавление файлов/папок с построением дерева FileEntry
- `RemoveEntry(path)`, `RenameEntry()`, `CreateDirectory()` — манипуляции с деревом
- `BrowseDirectory(path)` — навигация по файловой системе для FileBrowser
- `GetImagePreview(path)` — генерация JPEG-миниатюр (поддержка EXIF-ориентации)
- `GetFileProperties(path)` — метаинформация о файле

### BurnService

- `StartBurn(project)` — основной flow:
  1. Создание временного ISO через mkisofs (если UDF) или xorriso
  2. Запись ISO на диск через xorriso с streaming прогрессом
  3. Опциональная верификация
  4. Emit событий: `burn:progress`, `burn:state-changed`, `burn:complete`
- `CancelBurn()` — отмена через context cancellation
- `BlankDisc()`, `FormatDisc()` — подготовка RW-медиа

### SettingsService

- `GetSettings()` / `SaveSettings()` — `~/.config/xorriso-ui/settings.json`
- Поля: XorrisoPath, Language, Theme

## Отправка событий в frontend

```go
app := application.Get()
app.Event.Emit(models.EventBurnProgress, progressData)
```

**Важно**: поле `Event` (единственное число), НЕ `Events`.

## Тестирование

```bash
go test ./...                          # Все тесты
go test ./pkg/xorriso/ -run TestParse  # Конкретный тест
```

Существующие тесты:

- `pkg/xorriso/parser_test.go` — парсинг TOC, устройств, медиа
- `services/project_service_test.go` — операции с проектами
