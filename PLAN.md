# План: добавление UDF через mkisofs

## Контекст

xorriso-ui использует xorriso как бэкенд для записи дисков. Ранее мы удалили UDF из проекта, т.к. xorriso не поддерживает создание UDF (осознанное архитектурное решение разработчиков libisofs). Однако пользователю нужен UDF — это обязательный стандарт для Blu-ray спецификации и основной сценарий — запись BDXL дисков.

**Решение:** использовать **mkisofs** (cdrtools 3.02, уже установлен) для создания ISO-образов с UDF, а **xorriso** оставить как движок записи на диск. Это стандартный подход — K3b делает то же самое.

**Схема работы:**
1. **mkisofs** создаёт ISO-образ с UDF Bridge (UDF + ISO 9660) → temp-файл или stdout
2. **xorriso** записывает готовый образ на диск с прогрессом и контролем

**Что уже реализовано и работает:**
- Верификация (`-check_media`, `-check_md5`)
- BurnResult заполняется полностью
- CreateISO метод (но через xorriso native mode без UDF)
- Мультисессия
- Log-строки
- Все UI компоненты

---

## Изменения

### 1. Вернуть UDF в модель ISOOptions

**Файл:** `pkg/models/project.go`

```go
type ISOOptions struct {
    UDF        bool `json:"udf"`         // ДОБАВИТЬ ОБРАТНО
    ISOLevel   int  `json:"isoLevel"`
    RockRidge  bool `json:"rockRidge"`
    Joliet     bool `json:"joliet"`
    HFSPlus    bool `json:"hfsPlus"`
    Zisofs     bool `json:"zisofs"`
    MD5        bool `json:"md5"`
    BackupMode bool `json:"backupMode"`
}
```

### 2. Добавить MkisofsExecutor

**Новый файл:** `pkg/mkisofs/executor.go`

Отдельный executor для mkisofs (не xorriso), т.к. у mkisofs другой формат вывода и аргументов.

```go
type Executor struct {
    binaryPath string
}

func NewExecutor(binaryPath string) *Executor

// BuildISO создаёт ISO-образ с заданными опциями
// Возвращает путь к temp-файлу или ошибку
func (e *Executor) BuildISO(ctx context.Context, opts BuildOpts) (string, error)
```

**BuildOpts:**
```go
type BuildOpts struct {
    OutputPath string      // путь к выходному ISO
    VolumeID   string
    UDF        bool        // -udf
    RockRidge  bool        // -r
    Joliet     bool        // -J
    HFSPlus    bool        // -hfsplus
    Zisofs     bool        // -z
    Files      []FileMapping // source → dest
}

type FileMapping struct {
    Source string
    Dest   string
}
```

**Команда mkisofs:**
```bash
mkisofs -udf -r -J -V "VOLUME_ID" \
  -graft-points \
  /dest1=/source1 \
  /dest2=/source2 \
  -o /tmp/xorriso-ui-XXXXX.iso
```

Прогресс mkisofs: парсить stderr — mkisofs выводит `X.XX% done` строки.

### 3. Проверка свободного места перед созданием ISO

Перед созданием temp ISO (для UDF-записи) или пользовательского ISO (CreateISO) нужно проверить, что на целевом диске достаточно свободного места.

**Файл:** `services/burn_service.go`

- Вычислить размер будущего ISO: сумма размеров всех файлов проекта + ~5% overhead на метаданные FS
- Можно также использовать `mkisofs -print-size` для точного размера (возвращает количество секторов × 2048)
- Получить свободное место на разделе через `syscall.Statfs()` (путь к temp-директории или к outputPath)
- Если свободного места недостаточно — вернуть понятную ошибку ДО начала создания образа
- На фронтенде показать ошибку с указанием: нужно X ГБ, доступно Y ГБ

**Файл:** `frontend/src/views/BurnView.vue`

- Показывать предупреждение/ошибку если места недостаточно

### 4. Опция удаления temp ISO после успешного прожига

**Файл:** `pkg/models/project.go`

Добавить в BurnOptions:
```go
CleanupISO bool `json:"cleanupIso"` // удалять temp ISO после успешной записи (default: true)
```

**Файл:** `services/burn_service.go`

- После успешной записи (и верификации если включена): если `opts.CleanupISO == true` → `os.Remove(tempISOPath)`
- Если запись неуспешна — НЕ удалять (чтобы можно было повторить)
- Отправить log-строку: "Temporary ISO removed" / "Temporary ISO kept at: /path"

**Файл:** `frontend/src/views/BurnView.vue`

- Чекбокс "Delete temporary ISO after successful burn" в секции configure
- По умолчанию: включён

**Дефолты:** `CleanupISO: true` в `project_service.go` и `settings_service.go`

**i18n:** ключи `burn.cleanupIso` / `burn.cleanupIsoHint`

### 5. Рефакторинг BurnService — двухэтапная запись

**Файл:** `services/burn_service.go`

Текущий `runBurn()` использует native mode xorriso (`-dev`, `-map`, `-commit`). Новая логика:

**Если UDF включён:**
1. Вызвать `mkisofsExecutor.BuildISO()` → temp ISO файл с UDF
2. Записать temp ISO на диск через xorriso: `xorriso -as cdrecord dev=/dev/sr0 speed=X tmpfile.iso`
3. Удалить temp файл
4. Верификация (как сейчас)
5. Eject

**Если UDF выключен:**
- Оставить текущую логику (native mode xorriso: `-dev`, `-map`, `-commit`)

Это позволяет:
- UDF = двухэтапный подход (mkisofs + xorriso cdrecord mode)
- Без UDF = прямой native mode (быстрее, не нужен temp файл)

**Запись ISO через xorriso cdrecord mode:**
```bash
xorriso -as cdrecord dev=/dev/sr0 speed=4 -v -eject /tmp/image.iso
```
Или через native mode:
```bash
xorriso -dev /dev/sr0 -speed 4 -map /tmp/image.iso / -commit
```

Лучше использовать cdrecord mode для записи готового ISO — это стандартный паттерн.

**Добавить в CommandBuilder** (`pkg/xorriso/commands.go`):
- `CdrecordMode()` → `-as cdrecord`
- `CdrecordDev(dev string)` → `dev=<dev>`
- `CdrecordSpeed(speed string)` → `speed=<speed>`
- `Verbose()` → `-v`

### 4. Рефакторинг CreateISO

**Файл:** `services/burn_service.go`

Текущий `CreateISO()` использует xorriso native mode (`-outdev stdio:...`). Обновить:

**Если UDF:**
- Использовать mkisofs: `mkisofs -udf -r -J -V "NAME" -graft-points /dest=/src -o /output.iso`

**Если без UDF:**
- Оставить xorriso native: `-outdev stdio:<path> ... -commit`

### 5. Обновить main.go

**Файл:** `main.go`

- Добавить `exec.LookPath("mkisofs")` — если не найден, UDF будет недоступен (не fatal, просто warning)
- Создать mkisofs.Executor и передать в BurnService

### 6. Обновить BurnService конструктор

**Файл:** `services/burn_service.go`

```go
type BurnService struct {
    executor        *xorriso.Executor
    mkisofsExecutor *mkisofs.Executor  // может быть nil если mkisofs не найден
    mu              sync.Mutex
    currentJob      *models.BurnJob
    cancelFn        context.CancelFunc
}
```

### 7. Обновить дефолты

**Файл:** `services/project_service.go`

```go
ISOOptions: models.ISOOptions{
    UDF:       true,  // UDF по умолчанию для BDXL
    ISOLevel:  3,
    RockRidge: true,
    Joliet:    true,
    MD5:       true,
},
```

**Файл:** `services/settings_service.go`

Аналогично: `UDF: true` в DefaultISO.

### 8. Фронтенд: вернуть UDF чекбокс

**Файлы:**
- `frontend/src/views/SettingsView.vue` — вернуть чекбокс UDF
- `frontend/src/stores/tabStore.js` — `udf: true` в дефолтах
- `frontend/src/locales/en.json`, `ru.json` — вернуть ключи UDF
- При UDF=true и mkisofs не найден — показать предупреждение

### 9. Парсинг прогресса mkisofs

**Файл:** `pkg/mkisofs/progress.go`

mkisofs выводит прогресс в stderr в формате:
```
 10.02% done, estimate finish Sat Mar  8 12:00:00 2026
 20.04% done, estimate finish Sat Mar  8 12:01:00 2026
```

Парсить regex: `(\d+\.?\d*)%\s+done`

---

## Ключевые файлы

| Файл | Изменение |
|------|-----------|
| `pkg/mkisofs/executor.go` | **Новый** — запуск mkisofs |
| `pkg/mkisofs/progress.go` | **Новый** — парсинг прогресса mkisofs |
| `pkg/models/project.go` | Вернуть `UDF bool` |
| `pkg/xorriso/commands.go` | Добавить cdrecord mode методы |
| `services/burn_service.go` | Двухэтапная запись при UDF |
| `services/project_service.go` | `UDF: true` по умолчанию |
| `services/settings_service.go` | `UDF: true` по умолчанию |
| `main.go` | LookPath mkisofs, передать в BurnService |
| `frontend/src/views/SettingsView.vue` | Вернуть UDF чекбокс |
| `frontend/src/stores/tabStore.js` | `udf: true` |
| `frontend/src/locales/*.json` | Ключи UDF |

## Верификация

1. `go build ./...` — компиляция
2. `go test ./...` — тесты
3. `cd frontend && yarn build` — фронтенд
4. `wails3 dev` → создать проект с UDF=true → проверить в логах что вызывается mkisofs
5. Создать ISO с UDF → `mount -o loop` → `df -T` → убедиться что UDF FS
6. Записать на диск → верификация → убедиться что читается
