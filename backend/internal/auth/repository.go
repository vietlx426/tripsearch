package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	db "github.com/vietlx426/tripsearch/db/sqlc"
)

var ErrSessionNotFound = errors.New("session not found")

type SessionRepository interface {
	Insert(ctx context.Context, userID uuid.UUID, refreshToken, ip, userAgent string, expiresAt time.Time) (db.Session, error)
	FindByToken(ctx context.Context, refreshToken string) (db.Session, error)
	Delete(ctx context.Context, refreshToken string) error
	DeleteAllForUser(ctx context.Context, userID uuid.UUID) error
}

type sessionRepository struct {
	q *db.Queries
}

func NewSessionRepository(q *db.Queries) SessionRepository {
	return &sessionRepository{q: q}
}

func (r *sessionRepository) Insert(ctx context.Context, userID uuid.UUID, refreshToken, ip, userAgent string, expiresAt time.Time) (db.Session, error) {
	session, err := r.q.InsertSession(ctx, db.InsertSessionParams{
		UserID:       userID,
		RefreshToken: refreshToken,
		IpAddress:    sql.NullString{String: ip, Valid: ip != ""},
		UserAgent:    sql.NullString{String: userAgent, Valid: userAgent != ""},
		ExpiresAt:    expiresAt,
	})
	if err != nil {
		return db.Session{}, fmt.Errorf("Insert: %w", err)
	}
	return session, nil
}

func (r *sessionRepository) FindByToken(ctx context.Context, refreshToken string) (db.Session, error) {
	session, err := r.q.FindSessionByToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Session{}, ErrSessionNotFound
		}
		return db.Session{}, fmt.Errorf("FindByToken: %w", err)
	}
	return session, nil
}

func (r *sessionRepository) Delete(ctx context.Context, refreshToken string) error {
	if err := r.q.DeleteSession(ctx, refreshToken); err != nil {
		return fmt.Errorf("Delete: %w", err)
	}
	return nil
}

func (r *sessionRepository) DeleteAllForUser(ctx context.Context, userID uuid.UUID) error {
	if err := r.q.DeleteAllUserSessions(ctx, userID); err != nil {
		return fmt.Errorf("DeleteAllForUser: %w", err)
	}
	return nil
}
