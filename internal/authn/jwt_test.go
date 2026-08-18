package authn

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/djalil83/A-Radius/internal/authz"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "01234567890123456789012345678901"

func TestMiddlewareValidTokenInjectsPrincipal(t *testing.T) {
	cfg := Config{Secret: []byte(testSecret), Issuer: "a-radius", Audience: "profile-api"}
	token := signToken(t, cfg, jwt.SigningMethodHS256, Claims{
		Username: "admin",
		TenantID: "tenant-1",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "user-1",
			Issuer:    cfg.Issuer,
			Audience:  jwt.ClaimStrings{cfg.Audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		},
	})

	var got *authz.Principal
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = authz.PrincipalFromContext(r.Context())
		w.WriteHeader(http.StatusNoContent)
	})
	middleware, err := cfg.Middleware(next)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res := httptest.NewRecorder()
	middleware.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusNoContent)
	}
	if got == nil || got.UserID != "user-1" || got.TenantID != "tenant-1" || got.Username != "admin" {
		t.Fatalf("principal = %#v, want verified claims", got)
	}
}

func TestMiddlewareRejectsMissingToken(t *testing.T) {
	cfg := Config{Secret: []byte(testSecret), Issuer: "a-radius", Audience: "profile-api"}
	middleware, err := cfg.Middleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next must not be called")
	}))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	middleware.ServeHTTP(res, httptest.NewRequest(http.MethodGet, "/", nil))
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnauthorized)
	}
}

func TestMiddlewareRejectsWrongAudience(t *testing.T) {
	cfg := Config{Secret: []byte(testSecret), Issuer: "a-radius", Audience: "profile-api"}
	token := signToken(t, cfg, jwt.SigningMethodHS256, Claims{RegisteredClaims: jwt.RegisteredClaims{
		Subject:   "user-1",
		Issuer:    cfg.Issuer,
		Audience:  jwt.ClaimStrings{"another-api"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
	}})
	middleware, err := cfg.Middleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("next must not be called")
	}))
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	middleware.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnauthorized)
	}
}

func TestConfigRejectsWeakSecret(t *testing.T) {
	if err := (Config{Secret: []byte("short"), Issuer: "a-radius", Audience: "profile-api"}).Validate(); err == nil {
		t.Fatal("Validate() accepted a weak secret")
	}
}

func signToken(t *testing.T, cfg Config, method jwt.SigningMethod, claims Claims) string {
	t.Helper()
	token := jwt.NewWithClaims(method, claims)
	value, err := token.SignedString(cfg.Secret)
	if err != nil {
		t.Fatal(err)
	}
	return value
}
