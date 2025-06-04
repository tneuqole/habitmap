package session

import (
	"context"
	"time"

	"github.com/alexedwards/scs/v2"
)

const (
	UserIDKey = "UserID"
	FlashKey  = "Flash"
)

type Manager struct {
	*scs.SessionManager
}

func New() *Manager {
	// TODO: other config and use env var for secure
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour //nolint:mnd
	sessionManager.Cookie.Secure = false

	return &Manager{sessionManager}
}

func get[T any](m *Manager, ctx context.Context, key string) *T {
	val, ok := m.Get(ctx, key).(T)
	if !ok {
		return nil
	}
	return &val
}

func (m *Manager) SetUserID(ctx context.Context, userID int64) {
	m.Put(ctx, UserIDKey, userID)
}

func (m *Manager) RemoveUserID(ctx context.Context) {
	m.Remove(ctx, UserIDKey)
}

func (m *Manager) GetUserID(ctx context.Context) *int64 {
	return get[int64](m, ctx, UserIDKey)
}

func (m *Manager) SetFlash(ctx context.Context, msg string) {
	m.Put(ctx, FlashKey, msg)
}

func (m *Manager) GetFlash(ctx context.Context) *string {
	return get[string](m, ctx, FlashKey)
}

type SessionData struct {
	Flash           *string
	UserID          *int64
	IsAuthenticated bool
}

func (m *Manager) Data(ctx context.Context) SessionData {
	data := SessionData{
		Flash:  m.GetFlash(ctx),
		UserID: m.GetUserID(ctx),
	}

	if data.UserID != nil {
		data.IsAuthenticated = true
	}

	return data
}
