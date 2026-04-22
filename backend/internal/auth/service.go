package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vietlx426/tripsearch/internal/user"
)

type service struct {
	users    user.Repository
	sessions SessionRepository
	secret   string
}

func NewService(users user.Repository, sessions SessionRepository, secret string) Service {
	return &service{users: users, sessions: sessions, secret: secret}
}

func (s *service) Register(ctx context.Context, req RegisterRequest, ip, userAgent string) (*AuthResponse, error) {
	hashed, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("Register: %w", err)
	}

	u, err := s.users.Insert(ctx, req.Email, hashed, req.FullName, "traveler")
	if err != nil {
		if errors.Is(err, user.ErrEmailTaken) {
			return nil, ErrEmailTaken
		}
		return nil, fmt.Errorf("Register: %w", err)
	}

	tokens, err := issueTokenPair(s.secret, u.ID, u.Role)
	if err != nil {
		return nil, fmt.Errorf("Register: %w", err)
	}

	_, err = s.sessions.Insert(ctx, u.ID, tokens.refreshToken, ip, userAgent, time.Now().Add(refreshTokenTTL))
	if err != nil {
		return nil, fmt.Errorf("Register: %w", err)
	}

	return &AuthResponse{
		AccessToken:  tokens.accessToken,
		RefreshToken: tokens.refreshToken,
		User: UserDTO{
			ID:       u.ID,
			Email:    u.Email,
			FullName: u.FullName,
			Role:     u.Role,
		},
	}, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest, ip, userAgent string) (*AuthResponse, error) {
	u, err := s.users.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("Login: %w", err)
	}

	if err := verifyPassword(u.Password, req.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	tokens, err := issueTokenPair(s.secret, u.ID, u.Role)
	if err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}

	_, err = s.sessions.Insert(ctx, u.ID, tokens.refreshToken, ip, userAgent, time.Now().Add(refreshTokenTTL))
	if err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}

	return &AuthResponse{
		AccessToken:  tokens.accessToken,
		RefreshToken: tokens.refreshToken,
		User: UserDTO{
			ID:       u.ID,
			Email:    u.Email,
			FullName: u.FullName,
			Role:     u.Role,
		},
	}, nil
}

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	if err := s.sessions.Delete(ctx, refreshToken); err != nil {
		return fmt.Errorf("Logout: %w", err)
	}
	return nil
}

func (s *service) Refresh(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	session, err := s.sessions.FindByToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, fmt.Errorf("Refresh: %w", err)
	}

	u, err := s.users.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("Refresh: %w", err)
	}

	if err := s.sessions.Delete(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("Refresh: %w", err)
	}

	tokens, err := issueTokenPair(s.secret, u.ID, u.Role)
	if err != nil {
		return nil, fmt.Errorf("Refresh: %w", err)
	}

	_, err = s.sessions.Insert(ctx, u.ID, tokens.refreshToken, session.IpAddress.String, session.UserAgent.String, time.Now().Add(refreshTokenTTL))
	if err != nil {
		return nil, fmt.Errorf("Refresh: %w", err)
	}

	return &AuthResponse{
		AccessToken:  tokens.accessToken,
		RefreshToken: tokens.refreshToken,
		User: UserDTO{
			ID:       u.ID,
			Email:    u.Email,
			FullName: u.FullName,
			Role:     u.Role,
		},
	}, nil
}
