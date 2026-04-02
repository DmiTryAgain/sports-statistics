# Гайд: как добавить новое упражнение или категорию

## Добавление упражнения с существующей категорией

### 1. `pkg/tg/dictionary.go` — константа

```go
const (
    // ...
    myNewEx Exercise = "myNew"  // camelCase, уникальная строка-ключ
)
```

### 2. `pkg/tg/dictionary.go` — слова в `exerciseByLang`

Добавить все формы слова (именительный, родительный, винительный падежи + типичные опечатки) в обе секции — `langRU` и `langEN`.

```go
langRU: {
    // myNewEx
    "моё упражнение":  myNewEx,
    "моего упражнения": myNewEx,
    "мое упражнение":  myNewEx,  // без ё
    "моё упрожнение":  myNewEx,  // опечатка
},
langEN: {
    // myNewEx
    "my exercise":  myNewEx,
    "my exercises": myNewEx,
    "my exercize":  myNewEx,  // опечатка
},
```

Упражнения из нескольких слов добавляются целиком: `"bench press": benchPressEx`.
Поиск жадный — чем длиннее совпадение, тем оно приоритетнее.

### 3. `pkg/tg/dictionary.go` — `exTextByLang`

Каноничное отображаемое название (то, что видит пользователь в статистике).

```go
exTextByLang = map[language]map[Exercise]string{
    langRU: {
        // ...
        myNewEx: "моё упражнение",
    },
    langEN: {
        // ...
        myNewEx: "my exercise",
    },
}
```

### 4. `pkg/tg/dictionary.go` — `exerciseCategoryMap`

Если категория отличается от `CategoryReps` (только кол-во) — добавить явно:

```go
var exerciseCategoryMap = map[Exercise]ExerciseCategory{
    // ...
    myNewEx: CategoryRepsWeight,  // или нужная категория
}
```

Если упражнение только с кол-вом повторений — запись не нужна, `CategoryReps` используется по умолчанию.

### 5. `pkg/tg/dictionary.go` — `exerciseOptionalParamsMap`

Если у упражнения есть необязательные параметры (например, вес для подтягиваний):

```go
var exerciseOptionalParamsMap = map[Exercise][]ParamType{
    // ...
    myNewEx: {ParamWeight},
}
```

### 6. `pkg/tg/dictionary.go` — `exerciseOrder`

Добавить в список — определяет порядок отображения в кнопочном меню:

```go
var exerciseOrder = []Exercise{
    // ...
    myNewEx,
}
```

### 7. Проверка

```bash
make fmt lint
make test
```

`dictionary_test.go` автоматически проверяет, что каждое упражнение из `exerciseOrder` зарегистрировано во всех обязательных словарях.

---

## Добавление новой категории

Нужно, когда существующие категории не покрывают комбинацию параметров нового упражнения.

### Существующие категории для справки

| Категория | Обязательные | Soft-required | Опциональные |
|-----------|-------------|---------------|--------------|
| `CategoryReps` | count | — | — |
| `CategoryRepsWeight` | weight + count | — | — |
| `CategoryDistTime` | — | distance или time | — |
| `CategoryDuration` | time | — | — |
| `CategoryDurationWeight` | weight + time | — | — |
| `CategoryRepsOrDuration` | — | count или time | — |

Опциональные параметры задаются per-exercise через `exerciseOptionalParamsMap`, не через категорию.

### 1. `pkg/tg/model.go` — константа категории

```go
const (
    // ...
    CategoryDurationWeight  ExerciseCategory = iota
    CategoryRepsOrDuration
    CategoryMyNew          // добавить в конец
)
```

### 2. `pkg/tg/model.go` — `RequiredParams()`

Параметры, без которых сохранение невозможно:

```go
func (c ExerciseCategory) RequiredParams() []ParamType {
    switch c {
    // ...
    case CategoryMyNew:
        return []ParamType{ParamWeight, ParamCount}  // или nil если только soft-required
    }
}
```

### 3. `pkg/tg/model.go` — `SoftRequiredParams()`

Если нужен хотя бы один из нескольких параметров:

```go
func (c ExerciseCategory) SoftRequiredParams() []ParamType {
    switch c {
    // ...
    case CategoryMyNew:
        return []ParamType{ParamCount, ParamDuration}
    }
}
```

Если soft-required не нужен — case не добавлять, вернётся `nil` из `default`.

### 4. `pkg/tg/model.go` — `softRequiredError()` (только если добавили soft-required)

```go
func (c ExerciseCategory) softRequiredError() error {
    switch c {
    case CategoryRepsOrDuration:
        return errCountOrTimeRequired
    case CategoryMyNew:
        return errMyNewRequired  // новая ошибка
    default:
        return errDistOrTimeRequired
    }
}
```

### 5. `pkg/tg/handler.go` — новая ошибка (если нужна)

```go
var (
    // ...
    errMyNewRequired = errors.New("my new error message")
)
```

### 6. `pkg/tg/handler.go` — `missingParamMessage()`

Сообщение пользователю при вводе команды без обязательных параметров:

```go
func (m *MessageHandler) missingParamMessage(category ExerciseCategory, lang language) string {
    switch category {
    // ...
    case CategoryMyNew:
        return messagesByLang[lang][myNewRequired]
    }
}
```

### 7. `pkg/tg/handler.go` — `nextUnfilledParam()` (только если новые soft-required)

Эта функция определяет, какой параметр предложить следующим в кнопочном диалоге.
Строка, отвечающая за «optional» флаг — показывает кнопку «Пропустить»:

```go
optional := (sp == ParamDuration && params.DistanceM != nil) ||
    (sp == ParamDuration && params.Count != nil) ||
    (sp == ParamCount && params.DurationSec != nil)
    // добавить нужные условия для новой категории
```

Логика: параметр считается опциональным (с кнопкой «Пропустить»), если другой soft-required из пары уже заполнен.

### 8. `pkg/tg/keyboard.go` — функции форматирования

В `keyboard.go` одно место, где логика зависит от категории:

**`durationInlineKeyboard` и `optionalDurationInlineKeyboard`** — выбирают набор кнопок времени (короткие 30–180с vs длинные 15–60мин). Добавить категорию в условие коротких значений, если упражнение предполагает секунды/минуты, а не часы:
```go
if category == CategoryDuration || category == CategoryRepsOrDuration {
    durations = []float64{30, 45, 60, 90, 120, 180}
}
```

Форматирование подтверждения (`formatAddConfirmation`), кнопки Quick-Add (`formatQuickAddButton`) и команды копирования (`formatQuickCopyCommand`) работают универсально: проверяют наличие каждого параметра (вес, дистанция, время) и выводят то, что есть. Категорию-специфичную логику добавлять не нужно — параметры уже провалидированы перед сохранением.

### 9. `pkg/tg/handler.go` — `showParamStepCB()` (только если новые soft-required)

Функция показывает шаг диалога. В `ParamCount` и `ParamDuration` есть блок `detail` — он добавляет к названию упражнения уже введённые параметры. Если в новой категории count и duration soft-required, этот блок уже покрывает их (добавлено при реализации `CategoryRepsOrDuration`). Если новая категория использует другие комбинации — дополнить аналогично.

### 10. `pkg/tg/dictionary.go` — ключ и тексты сообщения

```go
// в iota-блоке констант ключей:
myNewRequired

// в messagesByLang[langRU]:
myNewRequired: "Нужно указать ... Пример: упражнение ...",

// в messagesByLang[langEN]:
myNewRequired: "... is required. Example: exercise ...",
```

---

## Чеклист при добавлении упражнения

- [ ] Константа `Exercise` в `dictionary.go`
- [ ] Слова в `exerciseByLang[langRU]` (падежи + опечатки)
- [ ] Слова в `exerciseByLang[langEN]` (формы + опечатки)
- [ ] Каноничное название в `exTextByLang[langRU]` и `exTextByLang[langEN]`
- [ ] Запись в `exerciseCategoryMap` (если не `CategoryReps`)
- [ ] Запись в `exerciseOptionalParamsMap` (если есть необязательные параметры)
- [ ] Добавить в `exerciseOrder`
- [ ] `make fmt lint && make test`

## Чеклист при добавлении категории

- [ ] Константа в `model.go` (в конец iota)
- [ ] Case в `RequiredParams()` в `model.go`
- [ ] Case в `SoftRequiredParams()` в `model.go` (если нужен)
- [ ] Case в `softRequiredError()` в `model.go` (если добавили soft-required)
- [ ] Ошибка `err...Required` в `handler.go` (если нужна новая)
- [ ] Case в `missingParamMessage()` в `handler.go`
- [ ] Обновить `nextUnfilledParam()` в `handler.go` (если новые soft-required)
- [ ] Обновить `showParamStepCB()` в `handler.go` (если нужно особое отображение шага)
- [ ] Обновить `durationInlineKeyboard` и `optionalDurationInlineKeyboard` в `keyboard.go` (если категория использует короткие интервалы времени)
- [ ] Ключ сообщения в iota и тексты RU+EN в `messagesByLang` в `dictionary.go`
- [ ] Далее — стандартный чеклист добавления упражнения
