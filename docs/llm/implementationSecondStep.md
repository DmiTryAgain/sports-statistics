# План реализации: Этап 2 — Гибридный интерфейс с кнопками и пошаговым диалогом

## Оглавление

1. [Обзор и цели](#1-обзор-и-цели)
2. [Изменения в БД](#2-изменения-в-бд)
3. [Изменения в слое pkg/db](#3-изменения-в-слое-pkgdb)
4. [Архитектура сессий и управление состоянием](#4-архитектура-сессий-и-управление-состоянием)
5. [ReplyKeyboard — постоянная клавиатура с командами](#5-replykeyboard--постоянная-клавиатура-с-командами)
6. [Обработка CallbackQuery](#6-обработка-callbackquery)
7. [Быстрое добавление частых упражнений (Quick-Add)](#7-быстрое-добавление-частых-упражнений-quick-add)
8. [Пошаговый диалог добавления упражнения](#8-пошаговый-диалог-добавления-упражнения)
9. [Inline-кнопки для показа статистики](#9-inline-кнопки-для-показа-статистики)
10. [Формат callback_data](#10-формат-callback_data)
11. [Изменения в словарях](#11-изменения-в-словарях)
12. [Изменения в handler.go и ListenAndHandle](#12-изменения-в-handlergo-и-listenandhandle)
13. [Новые файлы](#13-новые-файлы)
14. [Тесты](#14-тесты)
15. [Порядок реализации](#15-порядок-реализации)

---

## 1. Обзор и цели

### Что реализуем

Гибридный интерфейс, который:
- **Сохраняет** полную обратную совместимость с текстовым вводом из Этапа 1
- Добавляет **ReplyKeyboard** (постоянную клавиатуру) с кнопками команд — пользователь видит кнопки ещё ДО того, как начал печатать
- Добавляет **Quick-Add** — кнопки с наиболее частыми упражнениями пользователя для быстрой записи в один клик
- Добавляет **пошаговый диалог** — если пользователь не указал обязательные параметры, бот задаёт вопросы с inline-кнопками
- Добавляет **inline-кнопки** для выбора упражнений и параметров при просмотре статистики

### Ключевые принципы

1. **Текст всегда работает** — кнопки лишь дополнение, никогда не единственный путь
2. **Минимум шагов** — если пользователь опытный и пишет всё одним сообщением, кнопки не появляются
3. **Персонализация** — кнопки Quick-Add основаны на реальной истории пользователя
4. **Чистый чат** — inline-кнопки редактируются in-place, не засоряя переписку
5. **Таймауты** — незавершённые диалоги автоматически истекают

### Пользовательские сценарии

```
# Сценарий A: Опытный пользователь — текст (как сейчас)
User: сделал жим 80кг 10
Bot:  Добавлено ✅

# Сценарий B: Пользователь нажимает кнопку команды
User: *нажимает кнопку [Добавить]*
Bot:  Выбери упражнение:
      [Подтягивания] [Отжимания] [Жим лёжа] ...

      Или напиши текстом, например: жим 80кг 10

      Твои частые:
      [Жим 80кг x10] [Подтягивания x12] [Планка 60сек]

# Сценарий C: Quick-Add — одно нажатие
User: *нажимает [Жим 80кг x10]*
Bot:  Добавлено ✅ Жим лёжа: 80кг × 10

# Сценарий D: Пошаговый диалог — выбрал упражнение с весом
User: *нажимает [Жим лёжа]*
Bot:  Жим лёжа — укажи вес:
      [60кг] [70кг] [80кг] [90кг] [100кг] [Другой]
User: *нажимает [80кг]*
Bot:  Жим лёжа 80кг — сколько повторений?
      [5] [8] [10] [12] [15] [Другое]
User: *нажимает [10]*
Bot:  Добавлено ✅ Жим лёжа: 80кг × 10

# Сценарий E: Пошаговый диалог — ручной ввод на шаге
User: *нажимает [Жим лёжа]*
Bot:  Жим лёжа — укажи вес:
      [60кг] [70кг] [80кг] [90кг] [100кг] [Другой]
User: *нажимает [Другой]*
Bot:  Введи вес (например: 85кг или 85)
User: 85
Bot:  Жим лёжа 85кг — сколько повторений?
      [5] [8] [10] [12] [15] [Другое]
User: 10
Bot:  Добавлено ✅ Жим лёжа: 85кг × 10

# Сценарий F: Текст без параметров → переход в диалог
User: сделал жим
Bot:  Жим лёжа — укажи вес:
      [60кг] [70кг] [80кг] [90кг] [100кг] [Другой]

      Или напиши текстом: 80кг 10

# Сценарий G: Статистика через кнопки
User: *нажимает кнопку [Статистика]*
Bot:  Выбери упражнение:
      [Подтягивания] [Отжимания] [Жим лёжа] ...
      [Всё]
User: *нажимает [Всё]*
Bot:  Выбери период:
      [Сегодня] [Вчера] [Неделя] [Месяц] [Всё время]
User: *нажимает [Сегодня]*
Bot:  <таблица со статистикой>
```

---

## 2. Изменения в БД

### 2.1. Таблица `statistics` — без изменений

Таблица `statistics` уже содержит все необходимые поля для хранения упражнений с параметрами (поле `params jsonb` было расширено на Этапе 1). Новые DB-запросы (частые комбинации, уникальные упражнения, уникальные веса) работают с существующей схемой — они делают SELECT из существующих данных с агрегацией.

**Патчи в `docs/patches/` НЕ нужны.** Команды `make mfd-xml`, `make mfd-model`, `make mfd-repo`, `make mfd-dbtest` запускать НЕ нужно — модель `Statistic` и аннотация `docs/model/statistic.xml` не меняются.

### 2.2. Хранение пользовательских настроек — опционально (НЕ обязательно для этого этапа)

В будущем может понадобиться таблица для хранения пользовательских предпочтений (язык, часовой пояс). На текущем этапе язык определяется из текста сообщения и хранится в in-memory сессии — этого достаточно.

Если в будущем потребуется таблица `tg_users`:

```sql
-- docs/patches/YYYY-MM-DD-create_tg_users_table.sql
CREATE TABLE "tg_users"
(
    "tgUserId"   varchar(255)             NOT NULL,
    "lang"       varchar(2),
    "createdAt"  Timestamp with time zone NOT NULL Default now(),
    "updatedAt"  Timestamp with time zone NOT NULL Default now(),
    PRIMARY KEY ("tgUserId")
);
```

Тогда потребуется:
1. Добавить патч в `docs/patches/`
2. Выполнить патч на БД
3. Создать аннотацию в `docs/model/tg_user.xml`
4. Выполнить `make mfd-xml`, `make mfd-model`, `make mfd-repo`, `make mfd-dbtest`
5. Сгенерированные файлы появятся в `pkg/db/`: `tg_user.go` (репозиторий), а в `model.go`, `model_search.go`, `model_validate.go` добавятся соответствующие структуры

**На текущем этапе эту таблицу НЕ создаём** — in-memory сессий достаточно.

### 2.3. Индексы

Новые запросы используют фильтры по `tgUserId`, `exercise`, `statusId`, `createdAt` — все эти колонки уже имеют индексы (см. `docs/patches/2024-03-09-create_statistic_table.sql`). Дополнительные индексы не требуются.

Если в будущем Quick-Add запросы окажутся медленными на большом объёме данных, можно будет добавить составной индекс:

```sql
CREATE INDEX statistics_user_exercise_created
    ON statistics ("tgUserId", "exercise", "createdAt" DESC);
```

Но на текущем этапе это не требуется.

---

## 3. Изменения в слое `pkg/db`

### 3.1. Сгенерированные файлы — без изменений

Файлы, сгенерированные `mfd-generator`, **не трогаем**:
- `pkg/db/model.go` — структура `Statistic`, `Columns`, `Tables`
- `pkg/db/model_search.go` — `StatisticSearch` с фильтрами
- `pkg/db/model_validate.go` — валидация
- `pkg/db/statistic.go` — CRUD-методы репозитория

Модель `Statistic` уже имеет поле `Params *StatisticParams`, расширенное на Этапе 1 в ручном файле `pkg/db/model_params.go`. Перегенерация не требуется.

### 3.2. `pkg/db/model_params.go` — без изменений

Структура `StatisticParams` уже содержит все нужные поля (`WeightKg`, `DistanceM`, `DurationSec`) с методом `IsEmpty()`. Этого достаточно для всех сценариев Этапа 2.

### 3.3. `pkg/db/statistic_ext.go` — новые методы

В этот файл добавляются **3 новых метода** к `StatisticRepo`. Существующий код (`GroupedStatistic`, `GroupedStatisticSearch`, `GroupedStatisticByFilters`) **не меняется**.

#### 3.3.1. `FrequentExercisesByUser` — топ частых комбинаций (для Quick-Add)

Возвращает наиболее часто записываемые комбинации `exercise + count + params` за последние 30 дней. Используется для формирования кнопок Quick-Add.

```go
// FrequentExercise описывает часто используемую комбинацию упражнения с параметрами
type FrequentExercise struct {
    Exercise  string           `pg:"exercise,use_zero"`
    Count     float64          `pg:"count,use_zero"`
    Params    *StatisticParams `pg:"params"`
    Frequency int              `pg:"frequency,use_zero"`
}

func (sr StatisticRepo) FrequentExercisesByUser(ctx context.Context, tgUserID string, limit int) ([]FrequentExercise, error) {
    var result []FrequentExercise
    _, err := sr.db.QueryContext(ctx, &result, `
        SELECT t."exercise", t."count", t."params",
               COUNT(*) as "frequency"
        FROM statistics t
        WHERE t."tgUserId" = ?
          AND t."statusId" = ?
          AND t."createdAt" >= NOW() - INTERVAL '30 days'
        GROUP BY t."exercise", t."count", t."params"
        ORDER BY "frequency" DESC, MAX(t."createdAt") DESC
        LIMIT ?
    `, tgUserID, StatusEnabled, limit)
    if err != nil {
        return nil, fmt.Errorf("frequent exercises for user=%s, err=%w", tgUserID, err)
    }
    return result, nil
}
```

**SQL-логика:**
- `GROUP BY exercise, count, params` — группируем по полному набору (упражнение + кол-во + все параметры в jsonb). PostgreSQL сравнивает jsonb по значению, поэтому `{"weightKg": 80}` и `{"weightKg": 80}` — одна группа.
- `COUNT(*)` — частота этой комбинации
- `MAX(createdAt) DESC` — при одинаковой частоте, более свежие комбинации первые
- Фильтр по `createdAt >= NOW() - 30 days` — только актуальные привычки

**Пример результата:**

| exercise | count | params | frequency |
|----------|-------|--------|-----------|
| benchPress | 10 | {"weightKg": 80} | 8 |
| pullUp | 12 | NULL | 6 |
| plank | 1 | {"durationSec": 60} | 5 |

#### 3.3.2. `UniqueExercisesByUser` — список упражнений пользователя (для кнопок статистики)

Возвращает все уникальные упражнения, которые пользователь когда-либо записывал. Используется для формирования списка кнопок при просмотре статистики.

```go
func (sr StatisticRepo) UniqueExercisesByUser(ctx context.Context, tgUserID string) ([]string, error) {
    var result []struct {
        Exercise string `pg:"exercise"`
    }
    _, err := sr.db.QueryContext(ctx, &result, `
        SELECT DISTINCT t."exercise"
        FROM statistics t
        WHERE t."tgUserId" = ?
          AND t."statusId" = ?
        ORDER BY t."exercise"
    `, tgUserID, StatusEnabled)
    if err != nil {
        return nil, fmt.Errorf("unique exercises for user=%s, err=%w", tgUserID, err)
    }

    exercises := make([]string, len(result))
    for i := range result {
        exercises[i] = result[i].Exercise
    }
    return exercises, nil
}
```

**Примечание:** Используем анонимную структуру для маппинга, т.к. go-pg требует структуру с тегами.

#### 3.3.3. `UniqueWeightsByExercise` — уникальные веса для подсказок

Возвращает уникальные веса, которые пользователь использовал для конкретного упражнения. Используется для формирования кнопок выбора веса в пошаговом диалоге.

```go
func (sr StatisticRepo) UniqueWeightsByExercise(ctx context.Context, tgUserID, exercise string, limit int) ([]float64, error) {
    var result []struct {
        Weight float64 `pg:"weight"`
    }
    _, err := sr.db.QueryContext(ctx, &result, `
        SELECT DISTINCT (t."params"->>'weightKg')::float8 as "weight"
        FROM statistics t
        WHERE t."tgUserId" = ?
          AND t."exercise" = ?
          AND t."statusId" = ?
          AND t."params"->>'weightKg' IS NOT NULL
        ORDER BY "weight" DESC
        LIMIT ?
    `, tgUserID, exercise, StatusEnabled, limit)
    if err != nil {
        return nil, fmt.Errorf("unique weights for user=%s exercise=%s, err=%w", tgUserID, exercise, err)
    }

    weights := make([]float64, len(result))
    for i := range result {
        weights[i] = result[i].Weight
    }
    return weights, nil
}
```

**Почему сортировка DESC:** Тяжёлые веса обычно актуальнее для пользователя (основные рабочие подходы), поэтому показываем их первыми.

### 3.4. Константа `StatusEnabled`

В SQL-запросах используется `StatusEnabled`. Эта константа уже определена в `pkg/db/filter.go`:

```go
const (
    StatusEnabled = 1
    StatusDeleted = 2
)
```

### 3.5. Сводка изменений в `pkg/db`

| Файл | Изменения |
|------|-----------|
| `model.go` | Без изменений (сгенерирован) |
| `model_search.go` | Без изменений (сгенерирован) |
| `model_validate.go` | Без изменений (сгенерирован) |
| `model_params.go` | Без изменений |
| `statistic.go` | Без изменений (сгенерирован) |
| `statistic_ext.go` | **+** `FrequentExercise` struct, **+** `FrequentExercisesByUser()`, **+** `UniqueExercisesByUser()`, **+** `UniqueWeightsByExercise()` |
| `db.go` | Без изменений |

---

## 4. Архитектура сессий и управление состоянием

### Зачем нужно состояние

Текущий бот **stateless** — каждое сообщение обрабатывается изолированно. Для пошагового диалога нужно помнить, на каком шаге находится пользователь и что он уже выбрал.

### Структура `UserSession`

```go
// pkg/tg/session.go

// SessionState определяет текущий шаг диалога
type SessionState int

const (
    StateIdle             SessionState = iota // Ожидание команды (нет активного диалога)

    // Добавление упражнения
    StateAwaitExercise                        // Ожидание выбора упражнения
    StateAwaitWeight                          // Ожидание ввода веса
    StateAwaitCount                           // Ожидание ввода количества повторений
    StateAwaitDistance                         // Ожидание ввода дистанции
    StateAwaitDuration                        // Ожидание ввода времени
    StateAwaitCustomValue                     // Ожидание произвольного ввода (нажата кнопка "Другой")

    // Показ статистики
    StateShowAwaitExercise                    // Ожидание выбора упражнения для статистики
    StateShowAwaitPeriod                      // Ожидание выбора периода
)

// CustomValueTarget указывает, для какого параметра ожидается произвольный ввод
type CustomValueTarget int

const (
    TargetWeight   CustomValueTarget = iota
    TargetCount
    TargetDistance
    TargetDuration
)

// UserSession хранит состояние текущего диалога пользователя
type UserSession struct {
    State     SessionState
    Lang      language

    // Контекст добавления
    Exercise  Exercise         // Выбранное упражнение
    Params    ParsedParams     // Уже собранные параметры

    // Для "Другой" — какой параметр ожидаем
    CustomTarget CustomValueTarget

    // Контекст показа статистики
    ShowExercises Exercises    // Упражнения для показа

    // Сообщение с inline-кнопками (для редактирования in-place)
    LastBotMessageID int
    ChatID           int64

    // Таймер
    UpdatedAt time.Time
}

// IsExpired проверяет, не истекла ли сессия
func (s *UserSession) IsExpired(ttl time.Duration) bool {
    return time.Since(s.UpdatedAt) > ttl
}

// Reset сбрасывает сессию в начальное состояние
func (s *UserSession) Reset() {
    *s = UserSession{
        Lang:      s.Lang,
        UpdatedAt: time.Now(),
    }
}
```

### Хранилище сессий `SessionStore`

In-memory хранилище с TTL. На текущем этапе не нужен Redis или БД — бот однопроцессный, потеря сессий при перезапуске допустима (пользователь просто начнёт диалог заново).

```go
// pkg/tg/session.go

const sessionTTL = 10 * time.Minute

// SessionStore — потокобезопасное хранилище сессий
type SessionStore struct {
    mu       sync.RWMutex
    sessions map[string]*UserSession // ключ: tgUserID
}

// NewSessionStore создаёт хранилище и запускает фоновую горутину очистки.
// Горутина завершается при отмене ctx.
func NewSessionStore(ctx context.Context) *SessionStore {
    ss := &SessionStore{
        sessions: make(map[string]*UserSession),
    }
    go ss.cleanupLoop(ctx)
    return ss
}

// Get возвращает сессию пользователя. Если сессия истекла или не существует — nil.
func (ss *SessionStore) Get(userID string) *UserSession {
    ss.mu.RLock()
    defer ss.mu.RUnlock()

    s, ok := ss.sessions[userID]
    if !ok || s.IsExpired(sessionTTL) {
        return nil
    }
    return s
}

// Set создаёт или обновляет сессию
func (ss *SessionStore) Set(userID string, session *UserSession) {
    ss.mu.Lock()
    defer ss.mu.Unlock()

    session.UpdatedAt = time.Now()
    ss.sessions[userID] = session
}

// Delete удаляет сессию
func (ss *SessionStore) Delete(userID string) {
    ss.mu.Lock()
    defer ss.mu.Unlock()

    delete(ss.sessions, userID)
}

// cleanupLoop периодически удаляет истекшие сессии
func (ss *SessionStore) cleanupLoop(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            ss.mu.Lock()
            for k, v := range ss.sessions {
                if v.IsExpired(sessionTTL) {
                    delete(ss.sessions, k)
                }
            }
            ss.mu.Unlock()
        }
    }
}
```

Фоновая очистка инкапсулирована внутри `SessionStore` — `ListenAndHandle` не знает о ней. `context.Context` передаётся в конструктор, чтобы горутина корректно завершалась при остановке приложения.

---

## 5. ReplyKeyboard — постоянная клавиатура с командами

### Что это

`ReplyKeyboardMarkup` — постоянная клавиатура, отображаемая внизу экрана вместо обычной. Пользователь видит кнопки ДО того, как начинает печатать. Нажатие на кнопку отправляет текст кнопки как обычное сообщение.

### Когда отправлять

- При команде `/start` (первый запуск бота)
- При команде `помощь` / `help`
- Можно отправить один раз — Telegram запоминает клавиатуру для чата

### Раскладка кнопок

Для личного чата:
```
[📝 Добавить]  [📊 Статистика]  [❓ Помощь]
```

Для группового чата ReplyKeyboard не отправляется — в группах используется @mention текстом.

### Реализация

```go
// pkg/tg/keyboard.go

func replyKeyboard(lang language) tgbotapi.ReplyKeyboardMarkup {
    return tgbotapi.NewReplyKeyboard(
        tgbotapi.NewKeyboardButtonRow(
            tgbotapi.NewKeyboardButton(replyButtonTextByLang[lang][addBtn]),
            tgbotapi.NewKeyboardButton(replyButtonTextByLang[lang][showBtn]),
            tgbotapi.NewKeyboardButton(replyButtonTextByLang[lang][helpBtn]),
        ),
    )
}
```

### Обработка нажатий

Нажатие на ReplyKeyboard-кнопку отправляет текст. Текст кнопки маппится на команды:

```go
var replyButtonCmd = map[language]map[string]cmd{
    langRU: {
        "📝 Добавить":    addCmd,
        "📊 Статистика":  showCmd,
        "❓ Помощь":      helpCmd,
    },
    langEN: {
        "📝 Add":         addCmd,
        "📊 Statistics":  showCmd,
        "❓ Help":        helpCmd,
    },
}
```

При нажатии `[📝 Добавить]` бот:
1. Определяет команду `addCmd`
2. Текст после команды пуст → переходит к выбору упражнения через inline-кнопки
3. Создаёт сессию `StateAwaitExercise`
4. Отправляет сообщение с inline-кнопками упражнений + Quick-Add

### Команда `/start`

Добавляем обработку `/start`:

```go
func (m *MessageHandler) handleStart(upd tgbotapi.Update, lang language) {
    msg := tgbotapi.NewMessage(upd.Message.Chat.ID, messagesByLang[lang][welcomeMsg])
    msg.ReplyMarkup = replyKeyboard(lang)
    msg.ParseMode = m.cfg.ReplyFormat
    m.tgBot.Send(msg)
}
```

Сообщение приветствия описывает возможности бота и объясняет, что кнопки внизу — для быстрого доступа.

---

## 6. Обработка CallbackQuery

### Что такое CallbackQuery

Когда пользователь нажимает inline-кнопку, Telegram отправляет не Message, а `CallbackQuery`. В нём:
- `CallbackQuery.Data` — строка до 64 байт (наши закодированные данные)
- `CallbackQuery.Message` — сообщение, к которому были прикреплены кнопки
- `CallbackQuery.From` — кто нажал

### Обновление `ListenAndHandle`

Текущий код обрабатывает только `upd.Message`. Нужно добавить ветку для `upd.CallbackQuery`:

```go
for upd := range updates {
    limit <- struct{}{}

    go func() {
        defer func() { <-limit }()

        // Обработка inline-кнопок
        if upd.CallbackQuery != nil {
            m.handleCallback(ctx, upd)
            return
        }

        // Обработка обычных сообщений (существующий код)
        if upd.Message == nil {
            return
        }

        // ... существующая логика ...
    }()
}
```

### Метод `handleCallback`

```go
func (m *MessageHandler) handleCallback(ctx context.Context, upd tgbotapi.Update) {
    cb := upd.CallbackQuery
    userID := strconv.FormatInt(cb.From.ID, 10)

    // Парсим callback_data
    action, err := parseCallbackData(cb.Data)
    if err != nil {
        m.Error(ctx, "parse callback data", "data", cb.Data, "err", err)
        m.answerCallback(cb.ID, "")
        return
    }

    // Получаем или создаём сессию
    session := m.sessions.Get(userID)
    if session == nil {
        session = &UserSession{UpdatedAt: time.Now()}
    }

    // Обрабатываем действие
    switch action.Type {
    case cbSelectExercise:
        m.handleCBSelectExercise(ctx, cb, session, userID, action)
    case cbSelectWeight:
        m.handleCBSelectWeight(ctx, cb, session, userID, action)
    case cbSelectCount:
        m.handleCBSelectCount(ctx, cb, session, userID, action)
    case cbSelectDistance:
        m.handleCBSelectDistance(ctx, cb, session, userID, action)
    case cbSelectDuration:
        m.handleCBSelectDuration(ctx, cb, session, userID, action)
    case cbCustomInput:
        m.handleCBCustomInput(ctx, cb, session, userID, action)
    case cbQuickAdd:
        m.handleCBQuickAdd(ctx, cb, session, userID, action)
    case cbShowExercise:
        m.handleCBShowExercise(ctx, cb, session, userID, action)
    case cbShowPeriod:
        m.handleCBShowPeriod(ctx, cb, session, userID, action)
    case cbCancel:
        m.handleCBCancel(ctx, cb, session, userID)
    case cbExercisePage:
        m.handleCBExercisePage(ctx, cb, session, userID, action)
    default:
        m.answerCallback(cb.ID, "")
    }
}
```

### Ответ на callback и редактирование сообщения

```go
// answerCallback отвечает на callback query (убирает "часики" на кнопке)
func (m *MessageHandler) answerCallback(callbackID, text string) {
    callback := tgbotapi.NewCallback(callbackID, text)
    m.tgBot.Request(callback)
}

// editMessageWithKeyboard редактирует сообщение бота (in-place), обновляя текст и кнопки
func (m *MessageHandler) editMessageWithKeyboard(chatID int64, messageID int, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
    edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, text, keyboard)
    edit.ParseMode = m.cfg.ReplyFormat
    m.tgBot.Send(edit)
}

// editMessageRemoveKeyboard редактирует сообщение, убирая кнопки
func (m *MessageHandler) editMessageRemoveKeyboard(chatID int64, messageID int, text string) {
    edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
    edit.ParseMode = m.cfg.ReplyFormat
    m.tgBot.Send(edit)
}
```

---

## 7. Быстрое добавление частых упражнений (Quick-Add)

### Концепция

При нажатии `[📝 Добавить]` бот показывает не только список упражнений, но и кнопки Quick-Add — конкретные комбинации «упражнение + параметры», которые пользователь делал чаще всего.

Пример: если пользователь часто записывает «жим 80кг 10», эта комбинация появится как кнопка `[Жим 80кг ×10]`. Одно нажатие — запись добавлена.

### Источник данных

Новый DB-запрос — топ-N часто используемых комбинаций (упражнение + params + count) за последний месяц:

```sql
SELECT exercise, count, params,
       COUNT(*) as frequency
FROM statistics
WHERE "tgUserId" = ?
  AND "statusId" = 1
  AND "createdAt" >= NOW() - INTERVAL '30 days'
GROUP BY exercise, count, params
ORDER BY frequency DESC, MAX("createdAt") DESC
LIMIT 5
```

### Формат кнопок Quick-Add

Для каждой комбинации формируем текст кнопки:

| Категория | Формат кнопки |
|-----------|---------------|
| Reps | `Подтягивания ×12` |
| Reps+Weight | `Жим 80кг ×10` |
| Dist+Time | `Бег 5км 25мин` |
| Duration | `Планка 90сек` |

### Экран выбора после нажатия `[📝 Добавить]`

```
Выбери упражнение или напиши текстом.

Твои частые:
[Жим 80кг ×10]  [Подтягивания ×12]
[Планка 60сек]  [Бег 5км 25мин]

Все упражнения:
[Подтягивания] [Отжимания] [Брусья] ...
[>>]
[Отмена]
```

Quick-Add кнопки показываются **только если** у пользователя есть история. Для нового пользователя — только список упражнений.

### Обработка нажатия Quick-Add

При нажатии Quick-Add кнопки — сразу записываем в БД (все параметры уже известны), отправляем подтверждение. Никаких дополнительных шагов.

```go
func (m *MessageHandler) handleCBQuickAdd(ctx context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string, action CallbackAction) {
    // action содержит exercise, count, params (закодированы в callback_data)
    _, err := m.statRepo.AddStatistic(ctx, &db.Statistic{
        TgUserID: userID,
        Exercise: action.Exercise.String(),
        Count:    action.Count,
        Params:   action.Params,
        StatusID: 1,
    })
    if err != nil {
        m.Error(ctx, "quick add failed", "err", err)
        m.answerCallback(cb.ID, messagesByLang[session.Lang][errMsg])
        return
    }

    // Формируем текст подтверждения
    confirmText := formatAddConfirmation(action.Exercise, action.Count, action.Params, session.Lang)
    m.editMessageRemoveKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, confirmText)
    m.answerCallback(cb.ID, messagesByLang[session.Lang][exAdded])
    m.sessions.Delete(userID)
}
```

---

## 8. Пошаговый диалог добавления упражнения

### Машина состояний

```
                          текст с полными параметрами
                         ┌──────────────────────────────────────────────┐
                         │                                              ▼
[Команда/кнопка] → StateAwaitExercise → (выбрал) → по категории:
                                                    │
                   ┌────────────────────────────────┤
                   │                                │
             CategoryReps                     CategoryRepsWeight
                   │                                │
          StateAwaitCount               StateAwaitWeight
                   │                                │
                   ▼                          (выбрал вес)
              [Сохранить]                           │
                                            StateAwaitCount
                                                    │
                                               [Сохранить]
                   │
             CategoryDistTime                 CategoryDuration
                   │                                │
          StateAwaitDistance              StateAwaitDuration
                   │                                │
           (выбрал дистанцию)                  [Сохранить]
                   │
          StateAwaitDuration
                   │
              [Сохранить]
```

На любом шаге:
- Кнопка `[Отмена]` → сброс сессии, удаление кнопок
- Текстовое сообщение → попытка распарсить оставшиеся параметры текстом

### Экран: Выбор упражнения (`StateAwaitExercise`)

Упражнения группируются по категориям и разбиваются на страницы (по 8-10 кнопок на странице + навигация).

```
Выбери упражнение:

[Подтягивания]  [Отжимания]      [Брусья]
[Пресс]         [Приседания]     [Выпады]
[Бёрпи]         [Скакалка]       [Гиперэкстензия]

[>> Ещё]
[Отмена]
```

Страница 2:
```
[Жим лёжа]         [Становая тяга]    [Присед со штангой]
[Тяга верх. блока]  [Жим ногами]       [Скамья Скотта]
[Жим стоя]          [Тяга в наклоне]   [Бицепс гантели]

[<< Назад]  [>> Ещё]
[Отмена]
```

### Экран: Выбор веса (`StateAwaitWeight`)

Кнопки формируются на основе истории пользователя. Если есть записи для этого упражнения — берём уникальные веса. Если нет — стандартные значения.

```
Жим лёжа — укажи вес:

[40кг] [60кг] [70кг] [80кг] [100кг]
[Другой]
[Отмена]

Или напиши текстом (например: 85кг или 85)
```

Стандартные значения для нового пользователя:

| Упражнение | Стандартные веса |
|-----------|-----------------|
| По умолчанию | 20кг, 40кг, 60кг, 80кг, 100кг |

Из истории — берём до 5 уникальных весов, отсортированных по частоте.

### Экран: Выбор повторений (`StateAwaitCount`)

```
Жим лёжа 80кг — сколько повторений?

[5] [8] [10] [12] [15] [20]
[Другое]
[Отмена]

Или напиши число
```

### Экран: Выбор дистанции (`StateAwaitDistance`)

```
Бег — укажи дистанцию:

[1км] [2км] [3км] [5км] [10км]
[Другая]
[Отмена]
```

### Экран: Выбор времени (`StateAwaitDuration`)

Для бега (после дистанции):
```
Бег 5км — укажи время:

[15мин] [20мин] [25мин] [30мин] [45мин] [1ч]
[Другое]
[Отмена]
```

Для планки:
```
Планка — укажи время:

[30сек] [45сек] [60сек] [90сек] [2мин] [3мин]
[Другое]
[Отмена]
```

### Кнопка «Другой» / «Другое» (`StateAwaitCustomValue`)

При нажатии «Другой»:
1. Обновляем сообщение: «Введи вес (например: 85кг или 85)»
2. Устанавливаем `session.State = StateAwaitCustomValue`, `session.CustomTarget = TargetWeight`
3. Следующее текстовое сообщение пользователя интерпретируется как значение этого параметра

### Обработка текста во время активного диалога

Если у пользователя есть активная сессия с `State != StateIdle`, входящее текстовое сообщение обрабатывается в контексте сессии:

```go
func (m *MessageHandler) handleSessionText(ctx context.Context, upd tgbotapi.Update, session *UserSession, userID, text string) (string, error) {
    switch session.State {
    case StateAwaitExercise:
        // Пытаемся распарсить полное сообщение (упражнение + параметры)
        return m.handleAdd(ctx, text, userID, session.Lang)

    case StateAwaitCustomValue:
        // Парсим значение для нужного параметра
        return m.handleCustomValueInput(ctx, session, userID, text)

    case StateAwaitWeight, StateAwaitCount, StateAwaitDistance, StateAwaitDuration:
        // Пытаемся распарсить оставшиеся параметры
        return m.handleRemainingParamsText(ctx, session, userID, text)

    default:
        // Нет активного диалога — обычная обработка
        return "", nil // сигнал вызывающему коду обработать как обычно
    }
}
```

### Логика переходов между шагами

После получения каждого параметра проверяем, какие ещё не заполнены:

```go
func (m *MessageHandler) advanceAddDialog(ctx context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string) {
    category := session.Exercise.Category()
    required := category.RequiredParams()

    // Проверяем, все ли обязательные параметры собраны
    for _, rp := range required {
        switch rp {
        case ParamWeight:
            if session.Params.WeightKg == nil {
                m.showWeightSelection(ctx, cb, session, userID)
                return
            }
        case ParamCount:
            if session.Params.Count == nil {
                m.showCountSelection(ctx, cb, session, userID)
                return
            }
        case ParamDistance:
            if session.Params.DistanceM == nil {
                m.showDistanceSelection(ctx, cb, session, userID)
                return
            }
        case ParamDuration:
            if session.Params.DurationSec == nil {
                m.showDurationSelection(ctx, cb, session, userID)
                return
            }
        }
    }

    // Все параметры собраны — сохраняем
    m.saveFromSession(ctx, cb, session, userID)
}
```

### Сохранение результата

```go
func (m *MessageHandler) saveFromSession(ctx context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string) {
    cnt := session.Params.CountOrDefault()

    _, err := m.statRepo.AddStatistic(ctx, &db.Statistic{
        TgUserID: userID,
        Exercise: session.Exercise.String(),
        Count:    cnt,
        Params:   session.Params.ToDBParams(),
        StatusID: 1,
    })
    if err != nil {
        m.Error(ctx, "save from session failed", "err", err)
        m.answerCallback(cb.ID, messagesByLang[session.Lang][errMsg])
        return
    }

    // Формируем подтверждение
    confirmText := formatAddConfirmation(session.Exercise, cnt, session.Params.ToDBParams(), session.Lang)
    m.editMessageRemoveKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, confirmText)
    m.answerCallback(cb.ID, messagesByLang[session.Lang][exAdded])
    m.sessions.Delete(userID)
}
```

---

## 9. Inline-кнопки для показа статистики

### Сценарий

При нажатии `[📊 Статистика]`:

1. Показываем inline-кнопки с упражнениями + `[Всё]`
2. После выбора упражнения — показываем кнопки периодов
3. После выбора периода — отправляем таблицу

### Экран: Выбор упражнения для статистики

```
Выбери упражнение:

[Подтягивания]  [Отжимания]  [Жим лёжа]
[Всё]
[Отмена]
```

Здесь показываем только те упражнения, которые есть в истории пользователя (делаем запрос уникальных упражнений пользователя). Если упражнений у пользователя нет — пишем "Ничего не найдено".

### Экран: Выбор периода

```
Подтягивания — за какой период?

[Сегодня]  [Вчера]  [Неделя]
[Месяц]    [Год]    [Всё время]
[Отмена]
```

### Обработка

```go
func (m *MessageHandler) handleCBShowPeriod(ctx context.Context, cb *tgbotapi.CallbackQuery, session *UserSession, userID string, action CallbackAction) {
    // Формируем поиск
    p, _ := m.periodByText(string(action.Period), time.Now(), session.Lang)

    var periodsFilter periods
    if !p.IsZero() {
        periodsFilter = periods{p}
    }

    s := db.GroupedStatisticSearch{
        StatisticSearch: db.StatisticSearch{
            TgUserID:  &userID,
            Exercises: session.ShowExercises.StringSlice(),
        },
        Periods: periodsFilter.ToDB(),
    }

    stats, err := m.statRepo.GroupedStatisticByFilters(ctx, s)
    if err != nil {
        m.Error(ctx, "fetch statistic", "err", err)
        m.answerCallback(cb.ID, messagesByLang[session.Lang][errMsg])
        return
    }

    if len(stats) == 0 {
        m.editMessageRemoveKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, messagesByLang[session.Lang][nothingFound])
        m.answerCallback(cb.ID, "")
        m.sessions.Delete(userID)
        return
    }

    table, _ := m.buildTableByStat(ctx, stats, session.Lang)
    m.editMessageRemoveKeyboard(cb.Message.Chat.ID, cb.Message.MessageID, table)
    m.answerCallback(cb.ID, "")
    m.sessions.Delete(userID)
}
```

---

## 10. Формат callback_data

### Ограничения

Telegram ограничивает `callback_data` до **64 байт**. Нужен компактный формат.

### Формат

Используем простой текстовый формат с разделителем `|`:

```
TYPE|PARAM1|PARAM2|...
```

### Типы callback

| Тип | Формат | Пример | Описание |
|-----|--------|--------|----------|
| `ex` | `ex\|EXERCISE_KEY` | `ex\|benchPress` | Выбор упражнения для добавления |
| `w` | `w\|VALUE` | `w\|80` | Выбор веса (в кг) |
| `c` | `c\|VALUE` | `c\|10` | Выбор кол-ва повторений |
| `d` | `d\|VALUE` | `d\|5000` | Выбор дистанции (в метрах) |
| `t` | `t\|VALUE` | `t\|1500` | Выбор времени (в секундах) |
| `cu` | `cu\|TARGET` | `cu\|w` | Кнопка "Другой" (ожидание ручного ввода) |
| `qa` | `qa\|EX\|CNT\|W\|D\|T` | `qa\|benchPress\|10\|80\|\|` | Quick-Add (все параметры) |
| `se` | `se\|EXERCISE_KEY` | `se\|pullUp` | Выбор упражнения для статистики |
| `sa` | `sa` | `sa` | "Всё" для статистики |
| `sp` | `sp\|PERIOD` | `sp\|today` | Выбор периода для статистики |
| `pg` | `pg\|PAGE\|CTX` | `pg\|2\|add` | Переход на страницу (контекст: add/show) |
| `x` | `x` | `x` | Отмена |

### Проверка длины

Самый длинный callback — Quick-Add: `qa|romanianDeadlift|10|100||` = 30 байт. Укладываемся.

### Реализация парсинга

```go
// pkg/tg/callback.go

type CallbackType int

const (
    cbSelectExercise CallbackType = iota
    cbSelectWeight
    cbSelectCount
    cbSelectDistance
    cbSelectDuration
    cbCustomInput
    cbQuickAdd
    cbShowExercise
    cbShowAll
    cbShowPeriod
    cbExercisePage
    cbCancel
)

type CallbackAction struct {
    Type     CallbackType
    Exercise Exercise
    Value    float64
    Count    float64
    Params   *db.StatisticParams
    Period   textPeriod
    Page     int
    Context  string // "add" или "show"
    CustomTarget CustomValueTarget
}

func parseCallbackData(data string) (CallbackAction, error) {
    parts := strings.Split(data, "|")
    if len(parts) == 0 {
        return CallbackAction{}, fmt.Errorf("empty callback data")
    }

    switch parts[0] {
    case "ex":
        if len(parts) < 2 {
            return CallbackAction{}, fmt.Errorf("missing exercise in callback")
        }
        return CallbackAction{Type: cbSelectExercise, Exercise: Exercise(parts[1])}, nil

    case "w":
        v, err := strconv.ParseFloat(parts[1], 64)
        if err != nil {
            return CallbackAction{}, err
        }
        return CallbackAction{Type: cbSelectWeight, Value: v}, nil

    case "c":
        v, err := strconv.ParseFloat(parts[1], 64)
        if err != nil {
            return CallbackAction{}, err
        }
        return CallbackAction{Type: cbSelectCount, Value: v}, nil

    case "d":
        v, err := strconv.ParseFloat(parts[1], 64)
        if err != nil {
            return CallbackAction{}, err
        }
        return CallbackAction{Type: cbSelectDistance, Value: v}, nil

    case "t":
        v, err := strconv.ParseFloat(parts[1], 64)
        if err != nil {
            return CallbackAction{}, err
        }
        return CallbackAction{Type: cbSelectDuration, Value: v}, nil

    case "cu":
        target := TargetWeight
        switch parts[1] {
        case "w":
            target = TargetWeight
        case "c":
            target = TargetCount
        case "d":
            target = TargetDistance
        case "t":
            target = TargetDuration
        }
        return CallbackAction{Type: cbCustomInput, CustomTarget: target}, nil

    case "qa":
        return parseQuickAddCallback(parts)

    case "se":
        return CallbackAction{Type: cbShowExercise, Exercise: Exercise(parts[1])}, nil

    case "sa":
        return CallbackAction{Type: cbShowAll}, nil

    case "sp":
        return CallbackAction{Type: cbShowPeriod, Period: textPeriod(parts[1])}, nil

    case "pg":
        page, _ := strconv.Atoi(parts[1])
        ctx := ""
        if len(parts) > 2 {
            ctx = parts[2]
        }
        return CallbackAction{Type: cbExercisePage, Page: page, Context: ctx}, nil

    case "x":
        return CallbackAction{Type: cbCancel}, nil

    default:
        return CallbackAction{}, fmt.Errorf("unknown callback type: %s", parts[0])
    }
}

func encodeCallbackData(parts ...string) string {
    return strings.Join(parts, "|")
}
```

---

*Подробное описание новых DB-запросов см. в [разделе 3.3](#33-pkgdbstatistic_extgo--новые-методы).*

---

## 11. Изменения в словарях

### Новые текстовые константы

```go
// dictionary.go

const (
    // ... существующие ...

    // ReplyKeyboard кнопки
    addBtn  = iota + 100 // нумерация после существующих
    showBtn
    helpBtn

    // Сообщения для пошагового диалога
    welcomeMsg
    chooseExercise       // "Выбери упражнение:"
    chooseExerciseOrText // "Выбери упражнение или напиши текстом."
    yourFrequent         // "Твои частые:"
    chooseWeight         // "%s — укажи вес:"
    chooseCount          // "%s — сколько повторений?"
    chooseDistance        // "%s — укажи дистанцию:"
    chooseDuration       // "%s — укажи время:"
    enterCustomWeight    // "Введи вес (например: 85кг или 85)"
    enterCustomCount     // "Введи количество повторений"
    enterCustomDistance   // "Введи дистанцию (например: 5км или 5000м)"
    enterCustomDuration  // "Введи время (например: 25мин или 90сек)"
    customInputBtn       // "Другой" / "Other"
    cancelBtn            // "Отмена" / "Cancel"
    moreBtn              // "Ещё >>" / "More >>"
    backBtn              // "<< Назад" / "<< Back"
    allExBtn             // "Всё" / "All"
    choosePeriod         // "%s — за какой период?"
    addedConfirmation    // "Добавлено ✅ %s: %s"
    orWriteText          // "Или напиши текстом"
)
```

### Тексты кнопок ReplyKeyboard

```go
var replyButtonTextByLang = map[language]map[int]string{
    langRU: {
        addBtn:  "📝 Добавить",
        showBtn: "📊 Статистика",
        helpBtn: "❓ Помощь",
    },
    langEN: {
        addBtn:  "📝 Add",
        showBtn: "📊 Statistics",
        helpBtn: "❓ Help",
    },
}
```

### Пагинированный список упражнений

Все упражнения разбиваются на страницы для inline-кнопок. Порядок: сначала Reps (частые), потом RepsWeight, потом остальные.

```go
var exerciseOrder = []Exercise{
    // Reps — уличные / базовые
    pullUpEx, pushUpEx, dipsEx, absEx, squatEx, lungeEx,
    burpeeEx, skippingRopeEx, muscleUpEx, hyperextensionEx, legRaiseEx,
    // Reps+Weight — зал
    benchPressEx, deadliftEx, barbellSquatEx, shoulderPressEx,
    bentOverRowEx, latPulldownEx, seatedRowEx, legPressEx,
    dumbbellCurlEx, preacherCurlEx, tricepPushdownEx,
    legExtensionEx, legCurlEx, chestFlyEx,
    romanianDeadliftEx, hipThrustEx, lateralRaiseEx, shrugEx,
    // Duration / Distance
    plankEx, joggingEx,
}

const exercisesPerPage = 9 // 3 ряда по 3 кнопки
```

---

## 12. Изменения в handler.go и ListenAndHandle

### 12.1. Добавить `SessionStore` в `MessageHandler`

```go
type MessageHandler struct {
    embedlog.Logger

    dbc      db.DB
    statRepo db.StatisticRepo
    tgBot    *tgbotapi.BotAPI
    cfg      Bot
    sessions *SessionStore // НОВОЕ
}

func New(ctx context.Context, logger embedlog.Logger, db db.DB, statRepo db.StatisticRepo, tgBot *tgbotapi.BotAPI, cfg Bot) *MessageHandler {
    h := &MessageHandler{
        Logger:   logger,
        dbc:      db,
        cfg:      cfg,
        tgBot:    tgBot,
        statRepo: statRepo,
        sessions: NewSessionStore(ctx), // НОВОЕ: ctx для остановки горутины очистки
    }
    return h
}
```

**Примечание:** В сигнатуру `New()` добавляется `ctx context.Context`. Его нужно будет прокинуть из `main()`. Горутина очистки сессий запускается автоматически внутри `NewSessionStore(ctx)` и завершается при отмене контекста.

### 12.2. Обновить `ListenAndHandle`

```go
func (m *MessageHandler) ListenAndHandle(ctx context.Context) {
    updateConfig := tgbotapi.NewUpdate(0)
    updateConfig.Timeout = int(m.cfg.Timeout.Duration.Seconds())
    updates := m.tgBot.GetUpdatesChan(updateConfig)

    // Фоновая очистка сессий запускается автоматически в NewSessionStore(ctx)

    limit := make(chan struct{}, 100)
    for upd := range updates {
        limit <- struct{}{}

        go func() {
            defer func() { <-limit }()

            // НОВОЕ: обработка CallbackQuery
            if upd.CallbackQuery != nil {
                m.handleCallback(ctx, upd)
                return
            }

            if upd.Message == nil {
                return
            }

            // НОВОЕ: обработка /start
            if upd.Message.IsCommand() && upd.Message.Command() == "start" {
                m.handleStart(upd)
                return
            }

            lowerText := strings.ToLower(upd.Message.Text)
            hasMention := m.hasBotMention(lowerText)
            if !hasMention && upd.FromChat().IsGroup() {
                return
            }

            userID := strconv.FormatInt(upd.Message.From.ID, 10)
            msgText := m.clearRawMsg(lowerText)

            lang, err := m.detectLang(msgText)
            if err != nil {
                m.Print(ctx, err.Error(), "msg", msgText, "userID", userID)
                m.sendMsg(upd, "Can't detect a language😶 Please, use the only one keyboard layout chars")
                return
            }

            if msgText == "" {
                m.Print(ctx, "received empty message", "rawMsg", upd.Message.Text, "userID", userID)
                m.sendMsg(upd, messagesByLang[lang][emptyMessage])
                return
            }

            // НОВОЕ: проверяем, не нажал ли пользователь ReplyKeyboard кнопку
            if c, ok := replyButtonCmd[lang][upd.Message.Text]; ok {
                m.handleReplyButton(ctx, upd, userID, lang, c)
                return
            }

            // НОВОЕ: проверяем, есть ли активная сессия
            session := m.sessions.Get(userID)
            if session != nil && session.State != StateIdle {
                text, err := m.handleSessionText(ctx, upd, session, userID, msgText)
                if err != nil {
                    m.sendMsg(upd, messagesByLang[lang][errMsg])
                    return
                }
                if text != "" {
                    m.sendMsg(upd, text)
                    return
                }
                // text == "" означает, что сессия не смогла обработать — fallback на обычную обработку
            }

            text, err := m.handle(ctx, msgText, userID, lang)
            if err != nil {
                text = messagesByLang[lang][errMsg]
                m.Error(ctx, "an error occurred", "rawMsg", upd.Message.Text, "userID", userID, "err", err.Error())
            }

            m.sendMsg(upd, text)
        }()
    }
}
```

### 12.3. Обработка ReplyKeyboard кнопок

```go
func (m *MessageHandler) handleReplyButton(ctx context.Context, upd tgbotapi.Update, userID string, lang language, c cmd) {
    switch c {
    case addCmd:
        m.showAddExerciseScreen(ctx, upd, userID, lang)
    case showCmd:
        m.showStatExerciseScreen(ctx, upd, userID, lang)
    case helpCmd:
        text := fmt.Sprintf(messagesByLang[lang][commonHelpMsg], m.cfg.Name)
        m.sendMsg(upd, text)
    }
}
```

### 12.4. Отправка сообщений с inline-кнопками

```go
// sendMsgWithKeyboard отправляет сообщение с inline-кнопками
func (m *MessageHandler) sendMsgWithKeyboard(upd tgbotapi.Update, text string, keyboard tgbotapi.InlineKeyboardMarkup) int {
    msg := tgbotapi.NewMessage(upd.Message.Chat.ID, text)
    msg.ReplyMarkup = keyboard
    msg.ParseMode = m.cfg.ReplyFormat
    sent, err := m.tgBot.Send(msg)
    if err != nil {
        m.Errorf("failed to send message with keyboard: %v", err)
        return 0
    }
    return sent.MessageID
}
```

### 12.5. Связь текстового ввода с пошаговым диалогом

Важный момент: если пользователь пишет `сделал жим` (без параметров), вместо текстовой ошибки теперь нужно предложить inline-кнопки.

Обновляем `handleAdd`:

```go
func (m *MessageHandler) handleAdd(ctx context.Context, rawMsg, tgUserID string, lang language) (string, error) {
    // ... существующий код до проверки параметров ...

    if len(remainingWords) < 1 {
        // НОВОЕ: вместо текстовой ошибки — начинаем диалог с кнопками
        // Но только если это личный чат. В группе — текстовая ошибка как раньше.
        // (определяем из контекста или добавляем параметр isPrivate)
        // Для простоты пока показываем и текстовую подсказку
        lg.Print(ctx, "exercise requires params", "exercise", ex)
        return m.missingParamMessage(category, lang), nil
    }

    parsedParams, err := parseExerciseParams(remainingWords, category, lang)
    if err != nil {
        lg.Print(ctx, "param parsing error", "exercise", ex, "err", err)
        return m.missingParamMessage(category, lang), nil
    }

    // ... сохранение ...
}
```

**Важно**: в `handleReplyButton` для addCmd мы НЕ вызываем `handleAdd`, а сразу показываем экран выбора упражнения. `handleAdd` остаётся для текстового ввода. Если пользователь пишет `сделал жим` текстом — получает текстовую ошибку (как сейчас). Если нажимает кнопку → получает inline-кнопки.

Для более глубокой интеграции (Сценарий F) — в будущем можно добавить вариант, когда `handleAdd` при недостатке параметров в личном чате вместо текстовой ошибки открывает inline-диалог. Это потребует передачи `upd` в `handleAdd`, что сейчас не делается (только текст). Оставляем как опциональное улучшение.

---

## 13. Новые файлы

### 13.1. `pkg/tg/session.go`

Содержит:
- `SessionState` — состояния
- `CustomValueTarget` — цель произвольного ввода
- `UserSession` — структура сессии
- `SessionStore` — потокобезопасное хранилище
- Методы: `Get`, `Set`, `Delete`, `Cleanup`

### 13.2. `pkg/tg/callback.go`

Содержит:
- `CallbackType` — типы callback
- `CallbackAction` — результат парсинга
- `parseCallbackData()` — парсинг
- `encodeCallbackData()` — кодирование
- `parseQuickAddCallback()` — парсинг Quick-Add

### 13.3. `pkg/tg/keyboard.go`

Содержит:
- `replyKeyboard()` — создание ReplyKeyboard
- `exerciseInlineKeyboard()` — страница inline-кнопок упражнений
- `weightInlineKeyboard()` — кнопки выбора веса
- `countInlineKeyboard()` — кнопки выбора повторений
- `distanceInlineKeyboard()` — кнопки выбора дистанции
- `durationInlineKeyboard()` — кнопки выбора времени
- `periodInlineKeyboard()` — кнопки выбора периода
- `quickAddInlineKeyboard()` — кнопки Quick-Add
- `cancelButton()` — кнопка отмены
- `formatAddConfirmation()` — форматирование подтверждения

### Существующие файлы (обновления)

- `pkg/tg/handler.go` — новые методы обработки callback, обновлённый `ListenAndHandle`, `sendMsgWithKeyboard`
- `pkg/tg/dictionary.go` — новые текстовые константы, `replyButtonTextByLang`, `replyButtonCmd`, `exerciseOrder`
- `pkg/db/statistic_ext.go` — новые методы: `FrequentExercisesByUser`, `UniqueExercisesByUser`, `UniqueWeightsByExercise`

---

## 14. Тесты

### 14.1. Unit-тесты

**TestParseCallbackData:**
```go
{"ex|pullUp"}       → cbSelectExercise, Exercise=pullUpEx
{"w|80"}            → cbSelectWeight, Value=80
{"c|10"}            → cbSelectCount, Value=10
{"qa|benchPress|10|80||"} → cbQuickAdd, Exercise=benchPressEx, Count=10, Weight=80
{"sp|today"}        → cbShowPeriod, Period=todayPeriod
{"x"}               → cbCancel
{""}                → error
{"invalid|data"}    → error
```

**TestEncodeCallbackData:**
```go
encodeCallbackData("ex", "pullUp")         → "ex|pullUp"
encodeCallbackData("qa", "bench", "10", "80", "", "") → "qa|bench|10|80||"
```

**TestCallbackDataLength:**
Проверяем, что все возможные callback_data не превышают 64 байт:
```go
// Проходим по всем упражнениям, формируем максимальные callback_data
for _, ex := range exerciseOrder {
    data := encodeCallbackData("qa", string(ex), "999.99", "999.99", "99999", "99999")
    assert.LessOrEqual(t, len(data), 64)
}
```

**TestSessionStore:**
```go
// Get несуществующей сессии → nil
// Set + Get → корректное значение
// Expired → nil
// Delete → nil
// Cleanup → удаляет истёкшие
```

**TestSessionExpiry:**
```go
// Создаём сессию с UpdatedAt в прошлом
// IsExpired → true
```

### 14.2. Интеграционные тесты

**TestFrequentExercisesByUser:**
```go
// Создаём 5 записей жим 80кг x10, 3 записи подтягивания x12
// FrequentExercisesByUser → жим первый, подтягивания вторые
```

**TestUniqueExercisesByUser:**
```go
// Создаём записи по 3 разным упражнениям
// UniqueExercisesByUser → 3 уникальных упражнения
```

**TestUniqueWeightsByExercise:**
```go
// Создаём жим с весами 60, 80, 80, 100
// UniqueWeightsByExercise → [100, 80, 60]
```

### 14.3. Тесты формирования сообщений и клавиатур

Мокировать методы `Send`/`Request` из `tgbotapi.BotAPI` **не нужно** — это код внешней библиотеки, покрытый её собственными тестами. Вместо этого тестируем **нашу логику формирования** данных перед отправкой.

**TestExerciseInlineKeyboard:**
```go
// Проверяем, что клавиатура содержит правильные кнопки, callback_data, пагинацию
kb := exerciseInlineKeyboard(0, langRU, "add")
// Первая страница содержит exercisesPerPage кнопок
// Есть кнопка ">>"
// callback_data корректны: "ex|pullUp", "ex|pushUp", ...
```

**TestWeightInlineKeyboard:**
```go
// С историей: кнопки из уникальных весов пользователя
kb := weightInlineKeyboard([]float64{100, 80, 60}, langRU)
// Кнопки: [60кг] [80кг] [100кг] [Другой] [Отмена]
// callback_data: "w|60", "w|80", "w|100", "cu|w", "x"

// Без истории: стандартные значения
kb := weightInlineKeyboard(nil, langRU)
// Кнопки: [20кг] [40кг] [60кг] [80кг] [100кг] [Другой] [Отмена]
```

**TestQuickAddInlineKeyboard:**
```go
// Формируем кнопки из FrequentExercise
freq := []db.FrequentExercise{
    {Exercise: "benchPress", Count: 10, Params: &db.StatisticParams{WeightKg: ptr[float64](80)}},
    {Exercise: "pullUp", Count: 12, Params: nil},
}
kb := quickAddInlineKeyboard(freq, langRU)
// Текст кнопок: "Жим лёжа 80кг ×10", "Подтягивания ×12"
// callback_data: "qa|benchPress|10|80||", "qa|pullUp|12|||"
```

**TestFormatAddConfirmation:**
```go
// CategoryReps
formatAddConfirmation(pullUpEx, 12, nil, langRU) → "Добавлено ✅ Подтягивания: ×12"
// CategoryRepsWeight
formatAddConfirmation(benchPressEx, 10, &db.StatisticParams{WeightKg: ptr[float64](80)}, langRU)
    → "Добавлено ✅ Жим лёжа: 80кг × 10"
```

**TestPeriodInlineKeyboard:**
```go
kb := periodInlineKeyboard(langRU)
// 6 кнопок периодов + Отмена
// callback_data: "sp|today", "sp|yesterday", ..., "x"
```

---

## 15. Порядок реализации

**Схема БД не меняется.** Патчи в `docs/patches/` не нужны. Команды `make mfd-xml`, `make mfd-model`, `make mfd-repo`, `make mfd-dbtest` запускать не нужно. Сгенерированные файлы (`pkg/db/model.go`, `pkg/db/model_search.go`, `pkg/db/model_validate.go`, `pkg/db/statistic.go`) не затрагиваются.

| Шаг | Что делаем | Файлы | Зависимости |
|-----|-----------|-------|-------------|
| 1 | `SessionStore` и `UserSession` | `pkg/tg/session.go` | - |
| 2 | Тесты на `SessionStore` | `pkg/tg/session_test.go` | Шаг 1 |
| 3 | Формат callback_data: `parseCallbackData`, `encodeCallbackData` | `pkg/tg/callback.go` | - |
| 4 | Тесты на callback_data | `pkg/tg/callback_test.go` | Шаг 3 |
| 5 | Новые DB-запросы (ручной SQL): `FrequentExercise` struct, `FrequentExercisesByUser`, `UniqueExercisesByUser`, `UniqueWeightsByExercise` | `pkg/db/statistic_ext.go` | - |
| 6 | Интеграционные тесты на новые DB-запросы | `pkg/db/statistic_ext_test.go` или в существующих тестах | Шаг 5 |
| 7 | Словари: новые константы, `replyButtonTextByLang`, `replyButtonCmd`, `exerciseOrder` | `pkg/tg/dictionary.go` | - |
| 8 | Построение клавиатур: `replyKeyboard`, все inline-клавиатуры, `formatAddConfirmation` | `pkg/tg/keyboard.go` | Шаги 3, 7 |
| 9 | Добавить `SessionStore` в `MessageHandler`, обновить `New()` | `pkg/tg/handler.go` | Шаг 1 |
| 10 | Обработка `/start` и ReplyKeyboard | `pkg/tg/handler.go` | Шаги 8, 9 |
| 11 | Обработка `CallbackQuery` — роутинг по типам | `pkg/tg/handler.go` | Шаги 3, 9 |
| 12 | Пошаговый диалог добавления: обработчики `handleCBSelect*`, `advanceAddDialog`, `saveFromSession` | `pkg/tg/handler.go` | Шаги 8, 11 |
| 13 | Quick-Add: `handleCBQuickAdd`, построение кнопок из `FrequentExercisesByUser` | `pkg/tg/handler.go`, `pkg/tg/keyboard.go` | Шаги 5, 12 |
| 14 | Inline-кнопки для статистики: `handleCBShowExercise`, `handleCBShowPeriod` | `pkg/tg/handler.go` | Шаги 5, 11 |
| 15 | Обработка текста в контексте сессии: `handleSessionText`, `handleCustomValueInput` | `pkg/tg/handler.go` | Шаг 12 |
| 16 | Интеграционные тесты на полные сценарии (опционально — требуется мок бота) | `pkg/tg/handler_test.go` | Все шаги |
| 17 | `make fmt lint`, финальная проверка | - | Все шаги |

---

## Примечания

### Что НЕ входит в этот этап

- **Redis/внешнее хранилище сессий** — in-memory достаточно, бот однопроцессный
- **Клавиатура в группах** — ReplyKeyboard работает только в личных чатах; в группах остаётся текстовый ввод через @mention
- **Inline-режим** (`@bot_name` в строке ввода) — отдельная фича, не связана с этим этапом
- **Редактирование/удаление записей** — отдельный этап
- **Мультиязычные кнопки** — язык определяется по первому сообщению пользователя и запоминается в сессии

### Обратная совместимость

Весь текстовый ввод работает как раньше. Кнопки — это дополнение. Если пользователь никогда не нажмёт кнопку, его опыт не изменится.

### Ограничения Telegram API

- Inline-кнопки: максимум 8 в ряду, суммарно до 100 кнопок на сообщение
- `callback_data`: максимум 64 байта
- ReplyKeyboard: при отправке нового ReplyKeyboard старый заменяется
- `editMessageText`: можно редактировать только сообщения бота, не старше 48 часов
