package authz

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testJWTSecret = "01234567890123456789012345678901"

func testVerifier(t *testing.T) *JWTVerifier {
	t.Helper()
	v, err := NewJWTVerifier(JWTConfig{Secret: []byte(testJWTSecret), Issuer: "a-radius", Audience: "a-radius-api"})
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func signedTestToken(t *testing.T, claims JWTClaims, method jwt.SigningMethod, key any) string {
	t.Helper()
	token := jwt.NewWithClaims(method, claims)
	raw, err := token.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func TestJWTVerifierValidToken(t *testing.T) {
	v := testVerifier(t)
	now := time.Now()
	raw := signedTestToken(t, JWTClaims{TenantID: "tenant-1", PreferredUsername: "alice", RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1", Issuer: "a-radius", Audience: jwt.ClaimStrings{"a-radius-api"}, IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour))}}, jwt.SigningMethodHS256, []byte(testJWTSecret))
	p, err := v.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if p.UserID != "user-1" || p.TenantID != "tenant-1" || p.Username != "alice" {
		t.Fatalf("unexpected principal: %+v", p)
	}
}

func TestJWTVerifierRejectsInvalidClaimsAndAlgorithm(t *testing.T) {
	v := testVerifier(t)
	now := time.Now()
	base := JWTClaims{TenantID: "tenant-1", RegisteredClaims: jwt.RegisteredClaims{Subject: "user-1", Issuer: "a-radius", Audience: jwt.ClaimStrings{"a-radius-api"}, IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(-time.Minute))}}
	if _, err := v.Parse(signedTestToken(t, base, jwt.SigningMethodHS256, []byte(testJWTSecret))); err == nil {
		t.Fatal("expired token must be rejected")
	}
	base.ExpiresAt = jwt.NewNumericDate(now.Add(time.Hour))
	base.Issuer = "other-issuer"
	if _, err := v.Parse(signedTestToken(t, base, jwt.SigningMethodHS256, []byte(testJWTSecret))); err == nil {
		t.Fatal("wrong issuer must be rejected")
	}
	base.Issuer = "a-radius"
	if _, err := v.Parse(signedTestToken(t, base, jwt.SigningMethodHS512, []byte(testJWTSecret))); err == nil {
		t.Fatal("wrong algorithm must be rejected")
	}
	base.TenantID = ""
	if _, err := v.Parse(signedTestToken(t, base, jwt.SigningMethodHS256, []byte(testJWTSecret))); err == nil {
		t.Fatal("missing tenant must be rejected")
	}
}

func TestJWTMiddlewareRejectsMissingBearer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()
	testVerifier(t).Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { t.Fatal("next must not run") })).ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized || !strings.Contains(res.Body.String(), "authentication required") {
		t.Fatalf("unexpected response: %d %s", res.Code, res.Body.String())
	}
}
