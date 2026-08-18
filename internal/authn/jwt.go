package authn

import (
	"crypto/subtle"
	"errors"
	"net/http"
	"strings"

	"github.com/djalil83/A-Radius/internal/authz"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrMissingToken = errors.New("authorization bearer token is required")
	ErrInvalidToken = errors.New("token is invalid")
)

// Config menggunakan HMAC untuk contoh yang sederhana dan self-contained.
// Untuk produksi multi-service, gunakan RS256/EdDSA dengan JWKS dan key rotation.
type Config struct {
	Secret   []byte
	Issuer   string
	Audience string
}

type Claims struct {
	Username string `json:"username,omitempty"`
	TenantID string `json:"tenant_id,omitempty"`
	jwt.RegisteredClaims
}

func (c Config) Validate() error {
	if len(c.Secret) < 32 {
		return errors.New("JWT secret must be at least 32 bytes")
	}
	if strings.TrimSpace(c.Issuer) == "" {
		return errors.New("JWT issuer is required")
	}
	if strings.TrimSpace(c.Audience) == "" {
		return errors.New("JWT audience is required")
	}
	return nil
}

func (c Config) Middleware(next http.Handler) (http.Handler, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := bearerToken(r.Header.Get("Authorization"))
		if err != nil {
			writeAuthError(w, http.StatusUnauthorized, err.Error())
			return
		}

		claims := new(Claims)
		parser := jwt.NewParser(
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
			jwt.WithIssuer(c.Issuer),
			jwt.WithAudience(c.Audience),
			jwt.WithExpirationRequired(),
		)
		token, err := parser.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			// Jangan menerima algoritma dari header token tanpa allow-list.
			if token.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidToken
			}
			return c.Secret, nil
		})
		if err != nil || token == nil || !token.Valid || strings.TrimSpace(claims.Subject) == "" {
			writeAuthError(w, http.StatusUnauthorized, ErrInvalidToken.Error())
			return
		}

		// UserID berasal dari sub yang telah diverifikasi oleh signature dan claims.
		principal := &authz.Principal{
			UserID:   claims.Subject,
			Username: claims.Username,
			TenantID: claims.TenantID,
		}
		ctx := authz.WithPrincipal(r.Context(), principal)
		next.ServeHTTP(w, r.WithContext(ctx))
	}), nil
}

func bearerToken(value string) (string, error) {
	parts := strings.Fields(value)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", ErrMissingToken
	}
	return parts[1], nil
}

func writeAuthError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":{"code":"UNAUTHENTICATED","message":"` + message + `"}}`))
}

// constantTimeEqual tersedia jika validasi token pada masa depan membutuhkan
// perbandingan identifier rahasia tambahan.
func constantTimeEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
