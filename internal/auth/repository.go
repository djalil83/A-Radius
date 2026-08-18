package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrSessionNotFound = errors.New("session not found")
)

type DB interface {
	QueryRowContext(
		ctx context.Context,
		query string,
		args ...any,
	) *sql.Row

	ExecContext(
		ctx context.Context,
		query string,
		args ...any,
	) (sql.Result, error)
}

type Repository struct {
	DB DB
}

func NewRepository(db DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) FindUserByUsername(
	ctx context.Context,
	username string,
) (User, string, error) {
	if r == nil || r.DB == nil {
		return User{}, "", errors.New("auth repository is not configured")
	}

	const query = `
SELECT
id,
username,
password_hash,
status
FROM apb.users
WHERE username = $1
LIMIT 1
`

	var user User
	var passwordHash string

	err := r.DB.QueryRowContext(
		ctx,
		query,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&passwordHash,
		&user.Status,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return User{}, "", ErrUserNotFound
	}

	if err != nil {
		return User{}, "", err
	}

	return user, passwordHash, nil
}

func (r *Repository) CreateSession(
	ctx context.Context,
	userID string,
	tokenHash string,
	expiresAt time.Time,
	ipAddress string,
	userAgent string,
) error {
	if r == nil || r.DB == nil {
		return errors.New("auth repository is not configured")
	}

	const query = `
INSERT INTO apb.auth_sessions (
user_id,
token_hash,
expires_at,
ip_address,
user_agent
)
VALUES ($1, $2, $3, NULLIF($4, '')::inet, $5)
`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		userID,
		tokenHash,
		expiresAt,
		ipAddress,
		userAgent,
	)

	return err
}

func (r *Repository) FindUserBySessionTokenHash(
	ctx context.Context,
	tokenHash string,
) (User, error) {
	if r == nil || r.DB == nil {
		return User{}, errors.New("auth repository is not configured")
	}

	const query = `
SELECT
u.id,
u.username,
u.status
FROM apb.auth_sessions s
JOIN apb.users u
ON u.id = s.user_id
WHERE s.token_hash = $1
AND s.revoked_at IS NULL
AND s.expires_at > NOW()
LIMIT 1
`

	var user User

	err := r.DB.QueryRowContext(
		ctx,
		query,
		tokenHash,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Status,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrSessionNotFound
	}

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) RevokeSession(
	ctx context.Context,
	tokenHash string,
) error {
	if r == nil || r.DB == nil {
		return errors.New("auth repository is not configured")
	}

	const query = `
UPDATE apb.auth_sessions
SET revoked_at = NOW()
WHERE token_hash = $1
AND revoked_at IS NULL
`

	_, err := r.DB.ExecContext(
		ctx,
		query,
		tokenHash,
	)

	return err
}
