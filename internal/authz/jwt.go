package authz

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	Secret        []byte
	Issuer        string
	Audience      string
	Leeway        time.Duration
	AllowedMethod string
}

type JWTClaims struct {
	TenantID          string `json:"tenant_id"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	jwt.RegisteredClaims
}

type JWTVerifier struct {
	key        []byte
	issuer     string
	audience   string
	leeway     time.Duration
	allowedAlg string
	parser     *jwt.Parser
}

func NewJWTVerifier(cfg JWTConfig) (*JWTVerifier, error) {
	if len(cfg.Secret) < 32 {
		return nil, errors.New("JWT secret must be at least 32 bytes")
	}
	if cfg.Issuer == "" || cfg.Audience == "" {
		return nil, errors.New("JWT issuer and audience are required")
	}
	if cfg.Leeway <= 0 {
		cfg.Leeway = 30 * time.Second
	}
	if cfg.AllowedMethod == "" {
		cfg.AllowedMethod = "HS256"
	}
	if cfg.AllowedMethod != "HS256" {
		return nil, errors.New("only HS256 is supported by this verifier")
	}
	return &JWTVerifier{
		key: cfg.Secret, issuer: cfg.Issuer, audience: cfg.Audience, leeway: cfg.Leeway,
		allowedAlg: cfg.AllowedMethod,
		parser:     jwt.NewParser(jwt.WithValidMethods([]string{cfg.AllowedMethod}), jwt.WithIssuer(cfg.Issuer), jwt.WithAudience(cfg.Audience), jwt.WithLeeway(cfg.Leeway), jwt.WithIssuedAt()),
	}, nil
}

func (v *JWTVerifier) Parse(raw string) (*Principal, error) {
	if v == nil || v.parser == nil {
		return nil, ErrNotConfigured
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, ErrUnauthenticated
	}
	claims := &JWTClaims{}
	token, err := v.parser.ParseWithClaims(raw, claims, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != v.allowedAlg {
			return nil, errors.New("unexpected JWT signing method")
		}
		return v.key, nil
	})
	if err != nil || !token.Valid || claims.Subject == "" || claims.TenantID == "" {
		return nil, ErrUnauthenticated
	}
	return &Principal{UserID: claims.Subject, Username: claims.PreferredUsername, TenantID: claims.TenantID}, nil
}

func (v *JWTVerifier) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix = "Bearer "
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, prefix) {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		principal, err := v.Parse(strings.TrimSpace(strings.TrimPrefix(header, prefix)))
		if err != nil {
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithPrincipal(r.Context(), principal)))
	})
}

func PrincipalTenant(ctx context.Context) string {
	if p := PrincipalFromContext(ctx); p != nil {
		return p.TenantID
	}
	return ""
}
