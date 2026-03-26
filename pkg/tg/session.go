package tg

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

const sessionTTL = 10 * time.Minute

// SessionState определяет текущий шаг диалога
type SessionState int

const (
	StateIdle SessionState = iota // Ожидание команды (нет активного диалога)

	// Добавление упражнения
	StateAwaitExercise    // Ожидание выбора упражнения
	StateAwaitWeight      // Ожидание ввода веса
	StateAwaitCount       // Ожидание ввода количества повторений
	StateAwaitDistance    // Ожидание ввода дистанции
	StateAwaitDuration    // Ожидание ввода времени
	StateAwaitCustomValue // Ожидание произвольного ввода (нажата кнопка "Другой")

	// Показ статистики
	StateShowAwaitExercise // Ожидание выбора упражнения для статистики
	StateShowAwaitPeriod   // Ожидание выбора периода
)

// ParamType возвращает тип параметра, соответствующий состоянию ожидания ввода
func (s SessionState) ParamType() (ParamType, bool) {
	switch s {
	case StateAwaitWeight:
		return ParamWeight, true
	case StateAwaitCount:
		return ParamCount, true
	case StateAwaitDistance:
		return ParamDistance, true
	case StateAwaitDuration:
		return ParamDuration, true
	default:
		return 0, false
	}
}

// CustomValueTarget указывает, для какого параметра ожидается произвольный ввод
type CustomValueTarget int

const (
	TargetWeight CustomValueTarget = iota
	TargetCount
	TargetDistance
	TargetDuration
)

func customTargetToParamType(t CustomValueTarget) ParamType {
	switch t {
	case TargetWeight:
		return ParamWeight
	case TargetCount:
		return ParamCount
	case TargetDistance:
		return ParamDistance
	case TargetDuration:
		return ParamDuration
	default:
		return ParamCount
	}
}

// UserSession хранит состояние текущего диалога пользователя
type UserSession struct {
	State SessionState
	Lang  language

	// Контекст добавления
	Exercise Exercise
	Params   ParsedParams

	// Для "Другой" — какой параметр ожидаем
	CustomTarget CustomValueTarget

	// Пропущенные необязательные параметры
	SkippedParams map[ParamType]struct{}

	// Контекст показа статистики
	ShowExercises Exercises

	// Сообщение с inline-кнопками (для редактирования in-place)
	LastBotMessageID int
	ChatID           int64

	// Таймер
	UpdatedAt time.Time
}

func (s *UserSession) String() string {
	parts := []string{
		fmt.Sprintf("state: %d", s.State),
		fmt.Sprintf("lang: %s", s.Lang),
		fmt.Sprintf("exercise: %s", s.Exercise),
		fmt.Sprintf("params: %s", s.Params.String()),
		fmt.Sprintf("custom target: %d", s.CustomTarget),
		fmt.Sprintf("show exercises: %s", s.ShowExercises.String()),
		fmt.Sprintf("last bot message ID: %d", s.LastBotMessageID),
		fmt.Sprintf("chat ID: %d", s.ChatID),
		fmt.Sprintf("updated at: %s", s.UpdatedAt),
	}

	return strings.Join(parts, ", ")
}

// isParamSkipped проверяет, был ли параметр пропущен пользователем
func (s *UserSession) isParamSkipped(p ParamType) bool {
	if s.SkippedParams == nil {
		return false
	}
	_, ok := s.SkippedParams[p]
	return ok
}

// IsExpired проверяет, не истекла ли сессия
func (s *UserSession) IsExpired(ttl time.Duration) bool {
	return time.Since(s.UpdatedAt) > ttl
}

// Reset сбрасывает сессию в начальное состояние
func (s *UserSession) Reset() {
	lang := s.Lang
	*s = UserSession{
		Lang:      lang,
		UpdatedAt: time.Now(),
	}
}

// SessionStore — потокобезопасное хранилище сессий.
// Хранит значения UserSession (не указатели), чтобы Get возвращал независимую копию.
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]UserSession
}

// NewSessionStore создаёт хранилище и запускает фоновую горутину очистки.
// Горутина завершается при отмене ctx.
func NewSessionStore(ctx context.Context) *SessionStore {
	ss := &SessionStore{
		sessions: make(map[string]UserSession),
	}
	go ss.cleanupLoop(ctx)
	return ss
}

// Get возвращает сессию пользователя. Если сессия истекла или не существует — nil.
// Возвращает указатель на копию — вызывающий код может безопасно модифицировать сессию
// без блокировки, а затем записать обратно через Set.
func (ss *SessionStore) Get(userID string) *UserSession {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	s, ok := ss.sessions[userID]
	if !ok || s.IsExpired(sessionTTL) {
		return nil
	}
	return &s
}

// Set создаёт или обновляет сессию
func (ss *SessionStore) Set(userID string, session *UserSession) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	session.UpdatedAt = time.Now()
	ss.sessions[userID] = *session
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
