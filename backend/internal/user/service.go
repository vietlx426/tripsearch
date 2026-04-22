package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetMe(ctx context.Context, id uuid.UUID) (*UserDTO, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("GetMe: %w", err)
	}

	dto := &UserDTO{
		ID:       u.ID,
		Email:    u.Email,
		FullName: u.FullName,
		Role:     u.Role,
	}
	if u.AvatarUrl.Valid {
		dto.AvatarURL = &u.AvatarUrl.String
	}

	return dto, nil
}
