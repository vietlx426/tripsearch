package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	db "github.com/vietlx426/tripsearch/db/sqlc"
)

type repository struct {
	q *db.Queries
}

func NewRepository(q *db.Queries) Repository {
	return &repository{q: q}
}

func (r *repository) Insert(ctx context.Context, email, password, fullName, role string) (db.User, error) {
	user, err := r.q.InsertUser(ctx, db.InsertUserParams{
		Email:    email,
		Password: password,
		FullName: fullName,
		Role:     role,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return db.User{}, ErrEmailTaken
		}
		return db.User{}, fmt.Errorf("Insert: %w", err)
	}
	return user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (db.User, error) {
	user, err := r.q.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.User{}, ErrUserNotFound
		}
		return db.User{}, fmt.Errorf("FindByEmail: %w", err)
	}
	return user, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	user, err := r.q.FindUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.User{}, ErrUserNotFound
		}
		return db.User{}, fmt.Errorf("FindByID: %w", err)
	}
	return user, nil
}
