# Пайплайн записи диска

Описание процесса преобразования проекта (`.xorriso-project`) в команды xorriso и записи на физический диск.

## Общая схема

```
Project (.xorriso-project)
    │
    └─→ xorriso native mode
           xorriso строит ISO в памяти и сразу пишет на диск
           UDF включается через -udf on (libisofs)
```

## Цепочка вызовов

```
BurnService.StartBurn()
  → горутина runBurn()
    → buildISOCommand()     — Project → аргументы xorriso
    → executor.RunWithProgress()  — запуск subprocess с real-time парсингом
    → (опционально) верификация
    → (опционально) извлечение диска
```

## Построение команды xorriso

Функция `buildISOCommand()` использует fluent API `CommandBuilder` для трансляции полей проекта в аргументы xorriso.

### Маппинг полей Project → аргументов xorriso

| Поле проекта | Аргумент xorriso | Пример |
|---|---|---|
| — | `-pkt_output on` | Всегда включён (машинночитаемый вывод) |
| `volumeId` | `-volid` | `-volid PHOTOS_2025` |
| `isoOptions.isoLevel` | `-iso_level` | `-iso_level 3` |
| `isoOptions.rockRidge` | `-rockridge on` | `-rockridge on` |
| `isoOptions.joliet` | `-joliet on` | `-joliet on` |
| `isoOptions.udf` | `-udf on` | `-udf on` |
| `isoOptions.hfsPlus` | `-hfsplus on` | `-hfsplus on` |
| `isoOptions.zisofs` | `-zisofs level=9:...` | Сжатие zlib |
| `isoOptions.md5` | `-md5 on` | `-md5 on` |
| `isoOptions.backupMode` | `-acl on -xattr on` | Сохранение прав и атрибутов |
| `entries[].sourcePath` → `destPath` | `-map` | `-map /home/user/file.txt /file.txt` |
| `burnOptions.speed` | `-speed` | `-speed 8x` |
| `burnOptions.burnMode` | `-write_type` | `-write_type TAO` или `-write_type DAO` |
| `burnOptions.dummyMode` | `-dummy on` | Симуляция без записи |
| `burnOptions.closeDisc` | `-close on` | Финализация диска |
| `burnOptions.padding` | `-padding` | `-padding 300k` (в KiB) |
| `burnOptions.streamRecording` | `-stream_recording on` | Потоковая запись (Blu-ray) |
| `burnOptions.multisession` | `-close off` | Не закрывать сессию |
| устройство | `-dev` | `-dev /dev/sr0` |
| — | `-commit` | Запуск записи |

### Пример итоговой команды

```bash
xorriso \
  -pkt_output on \
  -dev /dev/sr0 \
  -volid "BACKUP_2025" \
  -iso_level 3 \
  -rockridge on \
  -joliet on \
  -udf on \
  -md5 on \
  -map /home/user/Documents /Documents \
  -map /home/user/Photos/vacation.jpg /Photos/vacation.jpg \
  -speed 8x \
  -write_type TAO \
  -padding 300k \
  -commit
```

### Пример: создание ISO-файла (без записи на диск)

При создании ISO вместо `-dev` и `-commit` используется `-outdev` с путём к файлу:

```bash
xorriso \
  -pkt_output on \
  -outdev /home/user/output.iso \
  -volid "BACKUP_2025" \
  -iso_level 3 \
  -rockridge on \
  -udf on \
  -map /home/user/Documents /Documents \
  -commit
```

## Верификация

Если включена опция `burnOptions.verify`, после записи выполняется проверка:

```bash
xorriso \
  -pkt_output on \
  -indev /dev/sr0 \
  -abort_on FAILURE \
  -md5 on \
  -check_md5_r FAILURE / -- \
  -check_media --
```

xorriso перечитывает записанные данные и проверяет контрольные суммы.

## Subprocess и парсинг вывода

### pkt_output — машинночитаемый формат xorriso

xorriso запускается с `-pkt_output on`, что переключает вывод в структурированный формат с префиксами каналов:

```
R:result line here
I:informational message here
M:marker line
```

| Канал | Назначение | Куда попадает |
|-------|-----------|---------------|
| `R:` | Результат команды | `CmdResult.ResultLines` |
| `I:` | Информация, предупреждения, ошибки | `CmdResult.InfoLines` |
| `M:` | Маркеры завершения команд | Внутренний контроль |

### Executor — управление subprocess

Два метода запуска:

#### `Run()` — синхронные операции

Для коротких запросов (список устройств, информация о медиа, TOC):

```
Executor.Run(ctx, args...) → CmdResult{ResultLines, InfoLines, ExitCode}
```

Запускает xorriso, ждёт завершения, парсит весь вывод.

#### `RunWithProgress()` — streaming операции

Для долгих операций (запись, очистка, форматирование, верификация):

```
Executor.RunWithProgress(ctx, progressFn, args...) → CmdResult
```

Читает stdout построчно в реальном времени. Для каждой строки:
1. Определяет канал (R:/I:/M:)
2. Извлекает прогресс через `ParsePacifierLine()`
3. Вызывает `progressFn(Progress)` — callback в BurnService
4. BurnService отправляет `Event.Emit()` → фронтенд обновляет UI

### ParsePacifierLine — извлечение прогресса

Парсит UPDATE-строки xorriso, содержащие метрики записи:

```
I:UPDATE: 45.2% done, speed= 6.1xBD, fifo 94%, buf  50%, remaining 0:03:22
```

Извлекаемые поля:

| Поле | Regex | Пример значения |
|------|-------|----------------|
| Процент | `(\d+\.?\d*)%\s+done` | `45.2` |
| Скорость | `(\d+\.?\d*x[A-Z]+\|\d+\.?\d*\s*[kMG]B/s)` | `6.1xBD` |
| FIFO | `fifo\s+(\d+)%` | `94` |
| ETA | `remaining\s+(\d+:\d+:\d+\|\d+:\d+)` | `0:03:22` |

## Фазы прогресса

Фронтенд получает события с фазой и процентом:

| Фаза | Описание |
|------|----------|
| `writing` | Запись на диск |
| `verifying` | Верификация (если включена) |

## Обработка ошибок

- **Отмена записи** — `CancelBurn()` отправляет сигнал отмены через `context.Cancel()`, subprocess завершается
- **Ошибка xorriso** — ненулевой exit code, ошибки из канала `I:`, результат передаётся на фронтенд

## Управление записью

### Извлечение диска

Если `burnOptions.eject: true`, после записи (и верификации) выполняется:

```bash
xorriso -pkt_output on -dev /dev/sr0 -eject all
```

### Статус задачи

`BurnService` хранит текущую задачу (`BurnJob`) с состояниями:

| Статус | Описание |
|--------|----------|
| `pending` | Задача создана, ожидает начала |
| `running` | Запись в процессе |
| `completed` | Запись завершена успешно |
| `failed` | Ошибка записи |
| `cancelled` | Отменена пользователем |
