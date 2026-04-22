package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/vietlx426/tripsearch/db/sqlc"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailTaken   = errors.New("email already taken")
)

type Repository interface {
	Insert(ctx context.Context, email, password, fullName, role string) (db.User, error)
	FindByEmail(ctx context.Context, email string) (db.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (db.User, error)
}

type Service interface {
	GetMe(ctx context.Context, id uuid.UUID) (*UserDTO, error)
}

type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	AvatarURL *string   `json:"avatar_url"`
}
