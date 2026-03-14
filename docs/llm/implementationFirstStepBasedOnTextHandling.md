# План реализации: упражнения с дополнительными параметрами (Этап 1 — текст с суффиксами)

## Оглавление

1. [Изменения в БД](#1-изменения-в-бд)
2. [Изменения в слое pkg/db](#2-изменения-в-слое-pkgdb)
3. [Новые типы в слое pkg/tg](#3-новые-типы-в-слое-pkgtg)
4. [Обновление словарей](#4-обновление-словарей)
5. [Парсинг параметров с суффиксами](#5-парсинг-параметров-с-суффиксами)
6. [Изменения в обработчиках](#6-изменения-в-обработчиках)
7. [Конвертация и форматирование единиц](#7-конвертация-и-форматирование-единиц)
8. [Тесты](#8-тесты)

---

## 1. Изменения в БД

### Схема НЕ меняется

Таблица `statistics` уже содержит поле `params jsonb` — именно для этих целей оно было зарезервировано. Колонка `count float8 NOT NULL` остаётся: для упражнений без повторений (планка, бег) будем записывать `count = 1` (одна "попытка"/подход).

**Патчи в `docs/patches/` НЕ нужны.**

Если в будущем понадобится индексировать по параметрам (например, для быстрой фильтрации по весу), можно будет добавить GIN-индекс на `params`:
```sql
CREATE INDEX statistics_params_gin ON statistics USING gin ("params");
```
Но на текущем этапе это не требуется.

---

## 2. Изменения в слое `pkg/db`

### 2.1. `model_params.go` — расширить `StatisticParams`

Текущее состояние — пустая структура. Нужно добавить поля для хранения нормализованных параметров:

```go
type StatisticParams struct {
    WeightKg    *float64 `json:"weightKg,omitempty"`
    DistanceM   *float64 `json:"distanceM,omitempty"`
    DurationSec *float64 `json:"durationSec,omitempty"`
}

// IsEmpty возвращает true, если ни один параметр не задан
func (sp *StatisticParams) IsEmpty() bool {
    if sp == nil {
        return true
    }
    return sp.WeightKg == nil && sp.DistanceM == nil && sp.DurationSec == nil
}
```

Все значения хранятся в **базовых единицах**:
- `WeightKg` — всегда в килограммах (ввод `500г` -> хранение `0.5`)
- `DistanceM` — всегда в метрах (ввод `5км` -> хранение `5000`)
- `DurationSec` — всегда в секундах (ввод `25мин` -> хранение `1500`)

Указатели (`*float64`) позволяют отличить "не задано" от "задано 0".

### 2.2. `statistic_ext.go` — обновить `GroupedStatistic` и SQL-запрос

Текущая структура:
```go
type GroupedStatistic struct {
    TgUserID string  `pg:"tgUserId,use_zero"`
    Exercise string  `pg:"exercise,use_zero"`
    SumCount float64 `pg:"sumCount,use_zero"`
    Sets     int     `pg:"sets,use_zero"`
}
```

Новая структура:
```go
type GroupedStatistic struct {
    TgUserID       string   `pg:"tgUserId,use_zero"`
    Exercise       string   `pg:"exercise,use_zero"`
    SumCount       float64  `pg:"sumCount,use_zero"`
    Sets           int      `pg:"sets,use_zero"`
    WeightKg       *float64 `pg:"weightKg"`
    DistanceM      *float64 `pg:"distanceM"`
    SumDurationSec *float64 `pg:"sumDurationSec"`
}
```

Обновлённый SQL в `GroupedStatisticByFilters`:

```sql
SELECT
    t."tgUserId",
    t."exercise",
    sum(t."count") as "sumCount",
    count(*) as "sets",
    (t."params"->>'weightKg')::float8 as "weightKg",
    (t."params"->>'distanceM')::float8 as "distanceM",
    sum((t."params"->>'durationSec')::float8) as "sumDurationSec"
FROM statistics t
WHERE ...
GROUP BY
    t."tgUserId",
    t."exercise",
    t."params"->>'weightKg',
    t."params"->>'distanceM'
ORDER BY t."exercise", "weightKg" DESC NULLS LAST, "distanceM" DESC NULLS LAST, "sumCount" DESC
```

**Почему GROUP BY по всем параметрам работает корректно:**
- Для reps-упражнений (подтягивания): `weightKg` = NULL, `distanceM` = NULL -> все строки группируются вместе (NULL = NULL в GROUP BY). Результат идентичен текущему.
- Для reps+weight (жим): `distanceM` = NULL -> группируется по exercise + weightKg. Каждый вес — отдельная строка.
- Для distance+time (бег): `weightKg` = NULL -> группируется по exercise + distanceM. Каждая дистанция — отдельная строка, время суммируется.
- Для duration-only (планка): все params NULL, кроме durationSec, но мы по нему не группируем, а суммируем -> одна строка с общим временем.

**ORDER BY** сортирует: по упражнению, затем по весу (тяжёлый первым), затем по дистанции (длинная первой), затем по количеству.

### 2.3. Файлы, сгенерированные mfd-generator

`model.go`, `model_search.go`, `model_validate.go` — **не трогаем**. Модель `Statistic` уже имеет поле `Params *StatisticParams`, а `StatisticParams` мы расширяем в ручном файле `model_params.go`. Перегенерация не требуется.

---

## 3. Новые типы в слое `pkg/tg`

### 3.1. `model.go` — категории упражнений

```go
// ExerciseCategory определяет, какие параметры ожидаются у упражнения
type ExerciseCategory int

const (
    CategoryReps       ExerciseCategory = iota // только повторения (текущие упражнения)
    CategoryRepsWeight                         // повторения + вес
    CategoryDistTime                           // дистанция + время
    CategoryDuration                           // только время (длительность)
)

// RequiredParams возвращает список обязательных типов параметров для категории
func (c ExerciseCategory) RequiredParams() []ParamType {
    switch c {
    case CategoryReps:
        return []ParamType{ParamCount}
    case CategoryRepsWeight:
        return []ParamType{ParamWeight, ParamCount}
    case CategoryDistTime:
        return []ParamType{ParamDistance, ParamDuration}
    case CategoryDuration:
        return []ParamType{ParamDuration}
    default:
        return nil
    }
}

// OptionalParams возвращает необязательные параметры
func (c ExerciseCategory) OptionalParams() []ParamType {
    switch c {
    case CategoryDistTime:
        return []ParamType{ParamCount} // для бега count необязателен
    default:
        return nil
    }
}
```

### 3.2. `model.go` — типы параметров и единицы измерения

```go
// ParamType — тип параметра упражнения
type ParamType int

const (
    ParamCount    ParamType = iota // количество повторений (голое число без суффикса)
    ParamWeight                    // вес: кг, г
    ParamDistance                   // дистанция: км, м
    ParamDuration                  // длительность: ч, мин, сек
)

// UnitDef описывает единицу измерения и коэффициент нормализации к базовой единице
type UnitDef struct {
    ParamType  ParamType
    Multiplier float64 // множитель для приведения к базовой единице
}
```

### 3.3. `model.go` — результат парсинга параметров

```go
// ParsedParams — результат парсинга текстовых параметров упражнения
type ParsedParams struct {
    Count       *float64 // повторения
    WeightKg    *float64 // вес в кг (нормализован)
    DistanceM   *float64 // дистанция в метрах (нормализована)
    DurationSec *float64 // длительность в секундах (нормализована)
}

// ToDBParams конвертирует в структуру для БД
func (pp ParsedParams) ToDBParams() *db.StatisticParams {
    p := &db.StatisticParams{
        WeightKg:    pp.WeightKg,
        DistanceM:   pp.DistanceM,
        DurationSec: pp.DurationSec,
    }
    if p.IsEmpty() {
        return nil
    }
    return p
}

// CountOrDefault возвращает count или 1, если не задан
func (pp ParsedParams) CountOrDefault() float64 {
    if pp.Count != nil {
        return *pp.Count
    }
    return 1
}
```

### 3.4. `model.go` — обновить метод `Exercise`

Добавить метод для получения категории:

```go
func (e Exercise) Category() ExerciseCategory {
    cat, ok := exerciseCategoryMap[e]
    if !ok {
        return CategoryReps // по умолчанию — reps
    }
    return cat
}
```

Заменить текущий метод `mustHaveCnt`:

```go
// mustHaveCnt проверяет, обязателен ли count для упражнения
func (e Exercise) mustHaveCnt() bool {
    for _, p := range e.Category().RequiredParams() {
        if p == ParamCount {
            return true
        }
    }
    return false
}
```

### 3.5. `model.go` — обновить `GroupedStatistic` (TG-модель)

```go
type GroupedStatistic struct {
    db.GroupedStatistic
    TranslatedExercise string
    Category           ExerciseCategory
}

func NewGroupedStatistic(in db.GroupedStatistic, lang language) GroupedStatistic {
    ex := Exercise(in.Exercise)
    return GroupedStatistic{
        GroupedStatistic:   in,
        TranslatedExercise: exTextByLang[lang][ex],
        Category:           ex.Category(),
    }
}
```

---

## 4. Обновление словарей

### 4.1. `dictionary.go` — новые упражнения

Новые константы:

```go
const (
    // ... существующие ...
    benchPressEx Exercise = "benchPress"   // жим лёжа
    deadliftEx   Exercise = "deadlift"     // становая тяга
    barbellSquatEx Exercise = "barbellSquat" // присед со штангой
    joggingEx    Exercise = "jogging"      // бег (раскомментировать)
    plankEx      Exercise = "plank"        // планка
)
```

Добавить в `exerciseByLang` словари для каждого нового упражнения (RU + EN + опечатки). Пример для жима:

```go
// benchPressEx
"жим":          benchPressEx,
"жим лёжа":     benchPressEx,
"жим лежа":     benchPressEx,
"жым":          benchPressEx,
"жым лёжа":     benchPressEx,
"жым лежа":     benchPressEx,
// EN
"bench":        benchPressEx,
"bench press":  benchPressEx,
"benchpress":   benchPressEx,
// ... опечатки ...
```

### 4.2. `dictionary.go` — маппинг категорий

```go
var exerciseCategoryMap = map[Exercise]ExerciseCategory{
    // CategoryReps — все текущие упражнения (pullUpEx, pushUpEx, и т.д.)
    // Их можно не перечислять, т.к. CategoryReps — значение по умолчанию

    // CategoryRepsWeight
    benchPressEx:   CategoryRepsWeight,
    deadliftEx:     CategoryRepsWeight,
    barbellSquatEx: CategoryRepsWeight,

    // CategoryDistTime
    joggingEx: CategoryDistTime,

    // CategoryDuration
    plankEx: CategoryDuration,
}
```

### 4.3. `dictionary.go` — словарь суффиксов единиц измерения

```go
var unitSuffixByLang = map[language]map[string]UnitDef{
    langRU: {
        // Вес
        "кг":  {ParamType: ParamWeight, Multiplier: 1},        // кг -> кг
        "г":   {ParamType: ParamWeight, Multiplier: 0.001},    // г -> кг
        // Дистанция
        "км":  {ParamType: ParamDistance, Multiplier: 1000},    // км -> м
        "м":   {ParamType: ParamDistance, Multiplier: 1},       // м -> м
        // Время
        "ч":   {ParamType: ParamDuration, Multiplier: 3600},   // ч -> сек
        "мин": {ParamType: ParamDuration, Multiplier: 60},     // мин -> сек
        "сек": {ParamType: ParamDuration, Multiplier: 1},      // сек -> сек
        "с":   {ParamType: ParamDuration, Multiplier: 1},      // с -> сек
        // Повторения (явный суффикс)
        "раз": {ParamType: ParamCount, Multiplier: 1},
        "р":   {ParamType: ParamCount, Multiplier: 1},
    },
    langEN: {
        // Weight
        "kg":   {ParamType: ParamWeight, Multiplier: 1},
        "lbs":  {ParamType: ParamWeight, Multiplier: 0.453592}, // фунты -> кг
        "lb":   {ParamType: ParamWeight, Multiplier: 0.453592},
        "g":    {ParamType: ParamWeight, Multiplier: 0.001},
        // Distance
        "km":   {ParamType: ParamDistance, Multiplier: 1000},
        "m":    {ParamType: ParamDistance, Multiplier: 1},
        "mi":   {ParamType: ParamDistance, Multiplier: 1609.34}, // мили -> м
        // Duration
        "h":    {ParamType: ParamDuration, Multiplier: 3600},
        "hr":   {ParamType: ParamDuration, Multiplier: 3600},
        "min":  {ParamType: ParamDuration, Multiplier: 60},
        "sec":  {ParamType: ParamDuration, Multiplier: 1},
        "s":    {ParamType: ParamDuration, Multiplier: 1},
        // Count
        "reps": {ParamType: ParamCount, Multiplier: 1},
        "rep":  {ParamType: ParamCount, Multiplier: 1},
        "x":    {ParamType: ParamCount, Multiplier: 1},
    },
}
```

### 4.4. `dictionary.go` — обновить `exHasCnt`

Заменить `exHasCnt` на `exerciseCategoryMap`. Старая `exHasCnt` больше не нужна — метод `mustHaveCnt()` теперь работает через `Category().RequiredParams()`.

### 4.5. `dictionary.go` — новые сообщения

Добавить ключи в `messagesByLang`:

```go
const (
    // ... существующие ...
    weightRequired   // "Для этого упражнения нужно указать вес. Пример: жим 80кг 10"
    durationRequired // "Для этого упражнения нужно указать время. Пример: планка 90сек"
    distanceRequired // "Для этого упражнения нужно указать дистанцию. Пример: бег 5км 25мин"
    paramInvalid     // "Не удалось распознать параметр: %s"
    tableWeightCol   // "вес" / "weight"
    tableDistCol     // "дистанция" / "distance"
    tableTimeCol     // "время" / "time"
)
```

### 4.6. `dictionary.go` — обновить `exTextByLang`, `exHasCnt`

Добавить переводы для новых упражнений:

```go
exTextByLang = map[language]map[Exercise]string{
    langRU: {
        // ... существующие ...
        benchPressEx:   "жим лёжа",
        deadliftEx:     "становая тяга",
        barbellSquatEx: "присед со штангой",
        joggingEx:      "бег",
        plankEx:        "планка",
    },
    langEN: {
        // ... существующие ...
        benchPressEx:   "bench press",
        deadliftEx:     "deadlift",
        barbellSquatEx: "barbell squat",
        joggingEx:      "jogging",
        plankEx:        "plank",
    },
}
```

Удалить `exHasCnt`, т.к. заменяется на `exerciseCategoryMap`.

---

## 5. Парсинг параметров с суффиксами

### 5.1. Новая функция `parseValueWithUnit`

Размещение: `pkg/tg/handler.go`

Парсит одно слово или пару слов (число + суффикс отдельно).

```go
// parseValueWithUnit пытается извлечь числовое значение и единицу измерения из слова.
// Поддерживает форматы: "80кг" (слитно) и возвращает (80, UnitDef{...}, true).
// Если слово — только число без суффикса, возвращает (80, UnitDef{}, false) —
// вызывающий код решает, что делать.
func parseValueWithUnit(word string, lang language) (value float64, unit UnitDef, hasUnit bool, err error)
```

Внутри: регулярное выражение `^(\d+[.,]?\d*)\s*([a-zA-Zа-яА-ЯёЁ]+)$` разделяет число и текстовый суффикс. Суффикс ищется в `unitSuffixByLang[lang]`.

### 5.2. Новая функция `parseExerciseParams`

Размещение: `pkg/tg/handler.go`

Главная функция парсинга параметров. Принимает слова после упражнения и категорию.

```go
// parseExerciseParams парсит слова после упражнения, извлекая параметры в соответствии с категорией.
//
// Логика:
// 1. Идёт по словам, для каждого пробует parseValueWithUnit.
// 2. Если число с суффиксом — записывает в соответствующий параметр (по ParamType из UnitDef).
// 3. Если число без суффикса — проверяет следующее слово: если это суффикс, объединяет.
//    Иначе трактует как count (повторения) — обратная совместимость.
// 4. Нормализует значения (умножает на Multiplier).
// 5. Если один ParamType встречается несколько раз, суммирует
//    (например "1ч 30мин" -> DurationSec = 3600 + 1800 = 5400).
// 6. Валидирует: проверяет, что все обязательные параметры категории заполнены.
//
// Возвращает ParsedParams и ошибку/текст подсказки при нехватке параметров.
func (m *MessageHandler) parseExerciseParams(
    words []string,
    category ExerciseCategory,
    lang language,
) (ParsedParams, error)
```

### 5.3. Алгоритм парсинга — пошагово

```
Вход: words=["80кг", "10"], category=CategoryRepsWeight, lang=langRU

Шаг 1: word="80кг"
  → parseValueWithUnit("80кг", langRU)
  → value=80, unit={ParamWeight, 1.0}, hasUnit=true
  → params.WeightKg = 80 * 1.0 = 80

Шаг 2: word="10"
  → parseValueWithUnit("10", langRU)
  → value=10, unit={}, hasUnit=false
  → Следующее слово? Нет.
  → Голое число → count
  → params.Count = 10

Валидация: CategoryRepsWeight требует [ParamWeight, ParamCount] → оба есть. OK.
```

```
Вход: words=["5", "км", "1", "ч", "30", "мин"], category=CategoryDistTime, lang=langRU

Шаг 1: word="5"
  → parseValueWithUnit("5", langRU) → value=5, hasUnit=false
  → Следующее слово = "км" → это суффикс! → unit={ParamDistance, 1000}
  → params.DistanceM = 5 * 1000 = 5000
  → Пропускаем "км"

Шаг 2: word="1"
  → parseValueWithUnit("1", langRU) → value=1, hasUnit=false
  → Следующее слово = "ч" → суффикс! → unit={ParamDuration, 3600}
  → params.DurationSec = 1 * 3600 = 3600

Шаг 3: word="30"
  → parseValueWithUnit("30", langRU) → value=30, hasUnit=false
  → Следующее слово = "мин" → суффикс! → unit={ParamDuration, 60}
  → params.DurationSec += 30 * 60 = 1800 → итого 5400

Валидация: CategoryDistTime требует [ParamDistance, ParamDuration] → оба есть. OK.
```

```
Вход: words=["10"], category=CategoryReps, lang=langRU (обратная совместимость)

Шаг 1: word="10"
  → parseValueWithUnit("10", langRU) → value=10, hasUnit=false
  → Нет суффикса → count
  → params.Count = 10

Валидация: CategoryReps требует [ParamCount] → есть. OK.
```

### 5.4. Обработка особого случая: `clearRawMsg` и суффиксы

Текущий `clearRawMsg` удаляет пунктуацию (`[[:punct:]]`), но буквы суффиксов (`кг`, `км`, `мин`) не являются пунктуацией — они сохраняются.

Точки в дробных числах уже защищены плейсхолдерами (`"80.5кг"` -> точка сохраняется).

**Изменения в `clearRawMsg` не требуются.**

---

## 6. Изменения в обработчиках

### 6.1. `handleAdd` — основные изменения

Текущая логика:
```
1. extractExerciseAndItsPosition → exercise, position
2. if mustHaveCnt: parseFloat(words[position+1]) → count
3. AddStatistic(count, params=nil)
```

Новая логика:
```
1. extractExerciseAndItsPosition → exercise, position
2. category = exercise.Category()
3. remainingWords = words[position+1:]
4. parsedParams = parseExerciseParams(remainingWords, category, lang)
5. count = parsedParams.CountOrDefault()
6. dbParams = parsedParams.ToDBParams()
7. AddStatistic(count, params=dbParams)
```

Для CategoryReps упражнений (все текущие) поведение идентично: `parseExerciseParams` найдёт голое число, запишет в Count, `ToDBParams()` вернёт `nil`. Полная обратная совместимость.

### 6.2. `handleShow` / `buildTableByStat` — вывод статистики

Текущий `buildTableByStat` строит одну таблицу с колонками: упражнение, кол-во, подходы.

Новый подход: **одна единая таблица** с динамическими колонками. Колонка доп. параметра (вес, дистанция, время) появляется только если хотя бы у одного упражнения в результате есть этот параметр. У упражнений, где параметр отсутствует, выводится `-`.

Если все упражнения в результате — reps-only (текущее поведение), дополнительные колонки не появляются. Полная обратная совместимость вывода.

```go
func (m *MessageHandler) buildTableByStat(ctx context.Context, in []db.GroupedStatistic, lang language) (string, error) {
    // 1. Преобразовать в TG-модели
    tgStats := NewGroupedStatisticList(in, lang)

    // 2. Определить, какие доп. колонки нужны
    hasWeight := anyHasWeight(tgStats)
    hasDistance := anyHasDistance(tgStats)
    hasDuration := anyHasDuration(tgStats)

    // 3. Построить единую таблицу с динамическими колонками
    // Колонки всегда: упражнение, [вес], [дистанция], [время], кол-во, подходы
    // Колонки в [] появляются только если hasXxx == true
}
```

Вспомогательные функции:

```go
func anyHasWeight(stats []GroupedStatistic) bool {
    for _, s := range stats {
        if s.WeightKg != nil {
            return true
        }
    }
    return false
}

func anyHasDistance(stats []GroupedStatistic) bool { ... }
func anyHasDuration(stats []GroupedStatistic) bool { ... }
```

Форматирование значения ячейки: если параметр есть — форматированное значение, если нет — `"-"`.

### 6.3. Пример: вывод статистики со всеми категориями одновременно

Пользователь за тренировку сделал упражнения из всех категорий и просит показать всё:

```
User: покажи всё за сегодня
```

Данные в БД (8 записей):

| exercise | count | params |
|----------|-------|--------|
| pullUp | 10 | NULL |
| pullUp | 12 | NULL |
| benchPress | 10 | {"weightKg": 60} |
| benchPress | 10 | {"weightKg": 80} |
| benchPress | 8 | {"weightKg": 80} |
| jogging | 1 | {"distanceM": 5000, "durationSec": 1500} |
| plank | 1 | {"durationSec": 60} |
| plank | 1 | {"durationSec": 90} |

SQL-запрос с `GROUP BY exercise, params->>'weightKg', params->>'distanceM'` вернёт:

| exercise | weightKg | distanceM | sumCount | sets | sumDurationSec |
|----------|----------|-----------|----------|------|----------------|
| pullUp | NULL | NULL | 22 | 2 | NULL |
| benchPress | 60 | NULL | 10 | 1 | NULL |
| benchPress | 80 | NULL | 18 | 2 | NULL |
| jogging | NULL | 5000 | 1 | 1 | 1500 |
| plank | NULL | NULL | 2 | 2 | 150 |

Go-код определяет: `hasWeight=true`, `hasDistance=true`, `hasDuration=true` — показываем все доп. колонки. Итоговый ответ бота:

```
упражнение      вес     дистанция    время         кол-во    подходы
подтягивания    -       -            -             22        2
жим лёжа        60кг    -            -             10        1
жим лёжа        80кг    -            -             18        2
бег             -       5км          25мин         1         1
планка          -       -            2мин 30сек    2         2
```

Если пользователь запросит только reps-упражнения (`покажи подтягивания за сегодня`):
- `hasWeight=false`, `hasDistance=false`, `hasDuration=false`
- Доп. колонки не появляются, таблица выглядит как текущая:

```
упражнение      кол-во    подходы
подтягивания    22        2
```

**Почему SQL работает корректно для всех категорий:**
- Для `pullUp` (CategoryReps): `weightKg=NULL`, `distanceM=NULL` — все записи в одной группе, `sumCount=22`. Идентично текущей логике.
- Для `benchPress` (CategoryRepsWeight): группировка по `weightKg` — каждый вес отдельной строкой.
- Для `jogging` (CategoryDistTime): группировка по `distanceM`, `sumDurationSec` суммируется.
- Для `plank` (CategoryDuration): `weightKg=NULL`, `distanceM=NULL` — все записи в одной группе, `sumDurationSec=150`.

NULL-значения в PostgreSQL GROUP BY группируются вместе, поэтому неиспользуемые параметры не влияют на группировку упражнений других категорий.

---

## 7. Конвертация и форматирование единиц

### 7.1. Нормализация на входе (уже описана в п. 5)

```
Ввод "80кг"    -> WeightKg = 80.0
Ввод "500г"    -> WeightKg = 0.5
Ввод "5км"     -> DistanceM = 5000.0
Ввод "800м"    -> DistanceM = 800.0
Ввод "1ч 30мин"-> DurationSec = 5400.0
Ввод "90сек"   -> DurationSec = 90.0
```

### 7.2. Форматирование на выходе

Новые функции в `pkg/tg/model.go`:

```go
// formatWeight форматирует вес из кг в удобный вид
func formatWeight(kg float64) string {
    if kg < 1 {
        return fmt.Sprintf("%.0fг", kg*1000)  // 0.5 -> "500г"
    }
    if kg == float64(int(kg)) {
        return fmt.Sprintf("%.0fкг", kg)       // 80.0 -> "80кг"
    }
    return fmt.Sprintf("%.1fкг", kg)            // 80.5 -> "80.5кг"
}

// formatDistance форматирует дистанцию из метров в удобный вид
func formatDistance(meters float64) string {
    if meters >= 1000 {
        km := meters / 1000
        if km == float64(int(km)) {
            return fmt.Sprintf("%.0fкм", km)   // 5000 -> "5км"
        }
        return fmt.Sprintf("%.1fкм", km)        // 2500 -> "2.5км"
    }
    return fmt.Sprintf("%.0fм", meters)          // 800 -> "800м"
}

// formatDuration форматирует длительность из секунд в удобный вид
func formatDuration(sec float64) string {
    totalSec := int(sec)
    if totalSec < 60 {
        return fmt.Sprintf("%dсек", totalSec)                    // 45 -> "45сек"
    }
    if totalSec < 3600 {
        min := totalSec / 60
        remSec := totalSec % 60
        if remSec == 0 {
            return fmt.Sprintf("%dмин", min)                     // 1500 -> "25мин"
        }
        return fmt.Sprintf("%dмин %dсек", min, remSec)          // 1530 -> "25мин 30сек"
    }
    hours := totalSec / 3600
    remMin := (totalSec % 3600) / 60
    if remMin == 0 {
        return fmt.Sprintf("%dч", hours)                         // 3600 -> "1ч"
    }
    return fmt.Sprintf("%dч %dмин", hours, remMin)               // 5400 -> "1ч 30мин"
}
```

Эти функции принимают значения в базовых единицах (кг, м, сек) и форматируют в человекочитаемый вид. Для EN-версий — аналогичные функции с английскими суффиксами или параметризация по языку.

Полиморфный вариант с учётом языка:

```go
func formatWeight(kg float64, lang language) string { ... }
func formatDistance(meters float64, lang language) string { ... }
func formatDuration(sec float64, lang language) string { ... }
```

---

## 8. Тесты

### 8.1. Unit-тесты — `handler_test.go`

Тесты на парсинг и форматирование не требуют БД.

**TestParseValueWithUnit:**
```go
{"80кг", langRU}        -> value=80, unit=ParamWeight, hasUnit=true
{"5.5км", langRU}       -> value=5.5, unit=ParamDistance, hasUnit=true
{"30мин", langRU}       -> value=30, unit=ParamDuration, hasUnit=true
{"10", langRU}          -> value=10, unit={}, hasUnit=false
{"80kg", langEN}        -> value=80, unit=ParamWeight, hasUnit=true
{"150lbs", langEN}      -> value=150, unit=ParamWeight, hasUnit=true
{"abc", langRU}         -> err
```

**TestParseExerciseParams:**
```go
// CategoryReps — обратная совместимость
{words: ["10"], cat: CategoryReps}
  -> Count=10, Weight=nil, Distance=nil, Duration=nil

// CategoryRepsWeight — слитный суффикс
{words: ["80кг", "10"], cat: CategoryRepsWeight}
  -> Count=10, WeightKg=80, Distance=nil, Duration=nil

// CategoryRepsWeight — раздельный суффикс
{words: ["80", "кг", "10"], cat: CategoryRepsWeight}
  -> Count=10, WeightKg=80, Distance=nil, Duration=nil

// CategoryRepsWeight — с явным суффиксом повторений
{words: ["80кг", "10раз"], cat: CategoryRepsWeight}
  -> Count=10, WeightKg=80

// CategoryDistTime — составное время
{words: ["5км", "1ч", "30мин"], cat: CategoryDistTime}
  -> Distance=5000, DurationSec=5400

// CategoryDuration
{words: ["90сек"], cat: CategoryDuration}
  -> DurationSec=90

// Ошибка — не хватает обязательного параметра
{words: ["80кг"], cat: CategoryRepsWeight}
  -> error (нет count)

// Ошибка — невалидное значение
{words: ["абвкг", "10"], cat: CategoryRepsWeight}
  -> error
```

**TestClearRawMsg — проверка сохранности суффиксов:**
```go
{rawMsg: "сделал жим 80кг 10", want: "сделал жим 80кг 10"}
{rawMsg: "сделал жим 80.5кг 10", want: "сделал жим 80.5кг 10"}
{rawMsg: "пробежал 5км 25мин", want: "пробежал 5км 25мин"}
```

**TestFormatWeight / TestFormatDistance / TestFormatDuration:**
```go
formatWeight(80, langRU)       -> "80кг"
formatWeight(0.5, langRU)      -> "500г"
formatWeight(80.5, langRU)     -> "80.5кг"
formatDistance(5000, langRU)    -> "5км"
formatDistance(800, langRU)     -> "800м"
formatDistance(2500, langRU)    -> "2.5км"
formatDuration(45, langRU)     -> "45сек"
formatDuration(1500, langRU)   -> "25мин"
formatDuration(5400, langRU)   -> "1ч 30мин"
formatDuration(1530, langRU)   -> "25мин 30сек"
```

### 8.2. Интеграционные тесты — `handler_test.go` (с БД)

Для интеграционных тестов используем пакет `pkg/db/test/`. Он предоставляет:
- `test.Setup(t)` — создание подключения к БД с автоматическим cleanup
- `test.Statistic(t, dbo, &db.Statistic{...}, ops...)` — создание записи в БД с возвратом `Cleaner`
- `test.WithFakeStatistic` — заполнение обязательных полей случайными данными
- `test.Cleaner` — функция очистки, вызывается через `defer`

Паттерн написания тестов: каждый тест-кейс создаёт нужные записи через `test.Statistic`, вызывает тестируемый метод, проверяет результат. Все записи автоматически удаляются через `Cleaner`.

**TestHandleAdd — новые упражнения:**
```go
{rawMsg: "жим 80кг 10", lang: langRU}      -> exAdded
{rawMsg: "bench 80kg 10", lang: langEN}     -> exAdded
{rawMsg: "планка 90сек", lang: langRU}      -> exAdded
{rawMsg: "plank 90sec", lang: langEN}       -> exAdded
{rawMsg: "бег 5км 25мин", lang: langRU}     -> exAdded
{rawMsg: "жим 10", lang: langRU}            -> weightRequired (нет веса)
{rawMsg: "жим 80кг", lang: langRU}          -> cntRequired (нет повторений)
```

**TestHandleShow — вывод с динамическими колонками:**
```go
// Только reps-упражнения — таблица без доп. колонок
// Reps + weight — таблица с колонкой "вес"
// Все категории — таблица со всеми доп. колонками, "-" для отсутствующих
```

---

## Порядок реализации

| Шаг | Что делаем | Файлы |
|-----|-----------|-------|
| 1 | Добавить типы: `ExerciseCategory`, `ParamType`, `UnitDef`, `ParsedParams` | `pkg/tg/model.go` |
| 2 | Расширить `StatisticParams` | `pkg/db/model_params.go` |
| 3 | Добавить словари: новые упражнения, суффиксы, `exerciseCategoryMap`, сообщения | `pkg/tg/dictionary.go` |
| 4 | Реализовать `parseValueWithUnit` и `parseExerciseParams` | `pkg/tg/handler.go` |
| 5 | Написать unit-тесты на парсинг | `pkg/tg/handler_test.go` |
| 6 | Обновить `handleAdd`: использовать `parseExerciseParams` вместо прямого `ParseFloat` | `pkg/tg/handler.go` |
| 7 | Написать функции форматирования: `formatWeight`, `formatDistance`, `formatDuration` | `pkg/tg/model.go` |
| 8 | Обновить `GroupedStatistic` (DB) и SQL-запрос | `pkg/db/statistic_ext.go` |
| 9 | Обновить `buildTableByStat`: единая таблица с динамическими колонками, `-` для отсутствующих | `pkg/tg/handler.go` |
| 10 | Написать интеграционные тесты на `handleAdd` с новыми упражнениями | `pkg/tg/handler_test.go` |
| 11 | Удалить `exHasCnt` | `pkg/tg/dictionary.go` |
