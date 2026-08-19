package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserDisabled       = errors.New("user account is not active")
)

type Service struct {
	repository *Repository
	sessionTTL time.Duration
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
		sessionTTL: 24 * time.Hour,
	}
}

func (s *Service) Login(
	ctx context.Context,
	username string,
	password string,
	ipAddress string,
	userAgent string,
) (User, string, time.Time, error) {
	if s == nil || s.repository == nil {
		return User{}, "", time.Time{}, errors.New(
			"auth service is not configured",
		)
	}

	username = strings.TrimSpace(username)

	if username == "" || password == "" {
		return User{}, "", time.Time{}, ErrInvalidCredentials
	}

	user, passwordHash, err := s.repository.FindUserByUsername(
		ctx,
		username,
	)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return User{}, "", time.Time{}, ErrInvalidCredentials
		}

		return User{}, "", time.Time{}, err
	}

	if user.Status != "active" {
		return User{}, "", time.Time{}, ErrUserDisabled
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(password),
	); err != nil {
		return User{}, "", time.Time{}, ErrInvalidCredentials
	}

	rawToken, err := generateToken()
	if err != nil {
		return User{}, "", time.Time{}, err
	}

	tokenHash := hashToken(rawToken)

	expiresAt := time.Now().Add(s.sessionTTL)

	if err := s.repository.CreateSession(
		ctx,
		user.ID,
		tokenHash,
		expiresAt,
		ipAddress,
		userAgent,
	); err != nil {
		return User{}, "", time.Time{}, err
	}

	return user, rawToken, expiresAt, nil
}

func (s *Service) Authenticate(
	ctx context.Context,
	rawToken string,
) (User, error) {
	if s == nil || s.repository == nil {
		return User{}, errors.New(
			"auth service is not configured",
		)
	}

	if rawToken == "" {
		return User{}, ErrSessionNotFound
	}

	user, err := s.repository.FindUserBySessionTokenHash(
		ctx,
		hashToken(rawToken),
	)
	if err != nil {
		return User{}, err
	}

	if user.Status != "active" {
		return User{}, ErrUserDisabled
	}

	return user, nil
}

func (s *Service) Logout(
	ctx context.Context,
	rawToken string,
) error {
	if s == nil || s.repository == nil {
		return errors.New(
			"auth service is not configured",
		)
	}

	if rawToken == "" {
		return nil
	}

	return s.repository.RevokeSession(
		ctx,
		hashToken(rawToken),
	)
}

func generateToken() (string, error) {
	buf := make([]byte, 32)

	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
