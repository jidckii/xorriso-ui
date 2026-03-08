# Пайплайн записи диска

Описание процесса преобразования проекта (`.xorriso-project`) в команды xorriso/mkisofs и записи на физический диск.

## Общая схема

```
Project (.xorriso-project)
    │
    ├─ UDF выключен ──→ Native mode (только xorriso)
    │                      xorriso строит ISO в памяти и сразу пишет на диск
    │
    └─ UDF включён ───→ UDF mode (mkisofs + xorriso)
                           Фаза 1: mkisofs создаёт временный ISO-файл
                           Фаза 2: xorriso записывает ISO на диск
```

Выбор режима определяется полем `isoOptions.udf` в проекте.

## Native mode (UDF выключен)

Используется только xorriso. ISO-образ создаётся в памяти и сразу записывается на диск — без промежуточного файла.

### Цепочка вызовов

```
BurnService.StartBurn()
  → горутина runBurn()
    → buildISOCommand()     — Project → аргументы xorriso
    → executor.RunWithProgress()  — запуск subprocess с real-time парсингом
    → (опционально) верификация
    → (опционально) извлечение диска
```

### Построение команды xorriso

Функция `buildISOCommand()` использует fluent API `CommandBuilder` для трансляции полей проекта в аргументы xorriso.

#### Маппинг полей Project → аргументов xorriso

| Поле проекта | Аргумент xorriso | Пример |
|---|---|---|
| — | `-pkt_output on` | Всегда включён (машинночитаемый вывод) |
| `volumeId` | `-volid` | `-volid PHOTOS_2025` |
| `isoOptions.isoLevel` | `-iso_level` | `-iso_level 3` |
| `isoOptions.rockRidge` | `-rockridge on` | `-rockridge on` |
| `isoOptions.joliet` | `-joliet on` | `-joliet on` |
| `isoOptions.hfsPlus` | `-hfsplus on` | `-hfsplus on` |
| `isoOptions.zisofs` | `-zisofs level=9:...` | Сжатие zlib |
| `isoOptions.md5` | `-md5 on` | `-md5 on` |
| `isoOptions.backupMode` | `-acl on -xattr on` | Сохранение прав и атрибутов |
| `entries[].sourcePath` → `destPath` | `-map` | `-map /home/user/file.txt /file.txt` |
| `burnOptions.speed` | `-speed` | `-speed 8x` |
| `burnOptions.burnMode` | `-write_type` | `-write_type TAO` или `-write_type DAO` |
| `burnOptions.dummyMode` | `-dummy on` | Симуляция без записи |
| `burnOptions.closeDisc` | `-close on` | Финализация диска |
| `burnOptions.padding` | `-padding` | `-padding 300` (в секторах) |
| `burnOptions.streamRecording` | `-stream_recording on` | Потоковая запись (Blu-ray) |
| `burnOptions.multisession` | `-close off` | Не закрывать сессию |
| устройство | `-dev` | `-dev /dev/sr0` |
| — | `-commit` | Запуск записи |

#### Пример итоговой команды

```bash
xorriso \
  -pkt_output on \
  -dev /dev/sr0 \
  -volid "BACKUP_2025" \
  -iso_level 3 \
  -rockridge on \
  -joliet on \
  -md5 on \
  -map /home/user/Documents /Documents \
  -map /home/user/Photos/vacation.jpg /Photos/vacation.jpg \
  -speed 8x \
  -write_type auto \
  -padding 300 \
  -commit
```

#### Пример: создание ISO-файла (без записи на диск)

При создании ISO вместо `-dev` и `-commit` используется `-outdev` с путём к файлу:

```bash
xorriso \
  -pkt_output on \
  -outdev /home/user/output.iso \
  -volid "BACKUP_2025" \
  -iso_level 3 \
  -rockridge on \
  -map /home/user/Documents /Documents \
  -commit
```

### Верификация

Если включена опция `burnOptions.verify`, после записи выполняется проверка:

```bash
xorriso \
  -pkt_output on \
  -indev /dev/sr0 \
  -check_media \
  -check_md5 on
```

xorriso перечитывает записанные данные и проверяет контрольные суммы.

## UDF mode (UDF включён)

Для UDF-дисков (Blu-ray, BDXL, DVD с UDF) используется двухфазный процесс: mkisofs для создания ISO, xorriso для записи.

Причина: xorriso имеет ограниченную поддержку UDF, а mkisofs (из cdrtools) поддерживает UDF полноценно.

### Цепочка вызовов

```
BurnService.StartBurn()
  → горутина runBurnUDF()
    → Фаза 1: проверка свободного места на диске хоста
    → Фаза 1: mkisofs.BuildISO()  — создание временного ISO
    → Фаза 2: xorriso (cdrecord mode) — запись ISO на диск
    → (опционально) верификация
    → (опционально) извлечение диска
    → (опционально) удаление временного ISO
```

### Фаза 1: создание ISO через mkisofs

#### Маппинг полей Project → аргументов mkisofs

| Поле проекта | Аргумент mkisofs | Пример |
|---|---|---|
| — | `-udf` | Всегда включён в UDF mode |
| `volumeId` | `-V` | `-V PHOTOS_2025` |
| `isoOptions.isoLevel` | `-iso-level` | `-iso-level 3` |
| `isoOptions.rockRidge` | `-r` | Rock Ridge с исправленными правами |
| `isoOptions.joliet` | `-J` | Joliet расширение |
| `isoOptions.hfsPlus` | `-hfsplus` | HFS+ расширение |
| `isoOptions.zisofs` | `-z` | Сжатие zisofs |
| `entries[]` | `-graft-points` | Режим graft-points для маппинга путей |
| `entries[].destPath=sourcePath` | (позиционные) | `/Photos/img.jpg=/home/user/img.jpg` |
| — | `-o` | `-o /tmp/xorriso-ui-XXXXX.iso` |

#### Формат graft-points

Файлы из `entries` преобразуются функцией `FileMappingsFromEntries()`:

```
destPath=sourcePath
```

Папки (`isDir: true`) пропускаются — mkisofs создаёт их автоматически из путей файлов.

#### Пример команды mkisofs

```bash
mkisofs \
  -udf \
  -r \
  -J \
  -V "PHOTOS_2025" \
  -iso-level 3 \
  -graft-points \
  -o /tmp/xorriso-ui-abc123.iso \
  /Documents/report.pdf=/home/user/Documents/report.pdf \
  /Photos/vacation.jpg=/home/user/Photos/vacation.jpg
```

#### Парсинг прогресса mkisofs

mkisofs выводит прогресс в stderr в формате:

```
 10.00% done, estimate finish Sat Mar  8 15:30:00 2025
 20.01% done, estimate finish Sat Mar  8 15:30:05 2025
```

Парсер использует regex `(\d+\.\d+)%\s+done` и кастомный сплиттер `scanMkisofsLines`, обрабатывающий `\r`-delimited строки (mkisofs использует `\r` для обновления прогресса в терминале).

### Фаза 2: запись ISO через xorriso (cdrecord mode)

xorriso запускается в режиме совместимости с cdrecord:

```bash
xorriso \
  -as cdrecord \
  -v \
  dev=/dev/sr0 \
  speed=8x \
  /tmp/xorriso-ui-abc123.iso
```

#### Маппинг burnOptions → аргументов cdrecord mode

| Поле проекта | Аргумент | Пример |
|---|---|---|
| устройство | `dev=` | `dev=/dev/sr0` |
| `burnOptions.speed` | `speed=` | `speed=8x` |
| `burnOptions.dummyMode` | `-dummy` | Симуляция |
| `burnOptions.streamRecording` | `stream_recording=on` | Потоковая запись |
| `burnOptions.burnMode` | `-dao` / `-tao` | Режим записи |
| `burnOptions.padding` | `padsize=` | `padsize=300s` |
| `burnOptions.closeDisc` | `-close` | Финализация |

### Очистка

Если `burnOptions.cleanupIso: true`, временный ISO-файл удаляется после успешной записи. При ошибке файл сохраняется для диагностики.

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

Фронтенд получает события с фазой и процентом. Фазы зависят от режима:

### Native mode

| Фаза | Описание |
|------|----------|
| `writing` | Запись на диск |
| `verifying` | Верификация (если включена) |

### UDF mode

| Фаза | Описание |
|------|----------|
| `creating_iso` | Создание ISO через mkisofs |
| `writing` | Запись ISO на диск через xorriso |
| `verifying` | Верификация (если включена) |

## Обработка ошибок

- **Нет места на хосте** (UDF mode) — перед созданием ISO проверяется свободное место на разделе, где будет временный файл
- **Отмена записи** — `CancelBurn()` отправляет сигнал отмены через `context.Cancel()`, subprocess завершается
- **Ошибка xorriso** — ненулевой exit code, ошибки из канала `I:`, результат передаётся на фронтенд
- **Ошибка mkisofs** — ненулевой exit code, stderr содержит описание ошибки

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
