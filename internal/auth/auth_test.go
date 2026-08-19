package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/djalil83/A-Radius/internal/authz"
)

type fakeAuthService struct {
	user User

	loginToken string

	loginExpiresAt time.Time

	logoutToken string
}

func (f *fakeAuthService) Login(
	ctx context.Context,
	username string,
	password string,
	ipAddress string,
	userAgent string,
) (User, string, time.Time, error) {
	if username != "customer01" || password != "secret" {
		return User{}, "", time.Time{}, ErrInvalidCredentials
	}

	return f.user, f.loginToken, f.loginExpiresAt, nil
}

func (f *fakeAuthService) Authenticate(
	ctx context.Context,
	rawToken string,
) (User, error) {
	if rawToken != f.loginToken {
		return User{}, ErrSessionNotFound
	}

	return f.user, nil
}

func (f *fakeAuthService) Logout(
	ctx context.Context,
	rawToken string,
) error {
	f.logoutToken = rawToken
	return nil
}

func TestLogin(t *testing.T) {
	expiresAt := time.Now().Add(time.Hour)

	service := &fakeAuthService{
		user: User{
			ID:       "user-123",
			Username: "customer01",
			Status:   "active",
		},
		loginToken:     "test-session-token",
		loginExpiresAt: expiresAt,
	}

	handler := NewHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/auth/login",
		strings.NewReader(
			`{"username":"customer01","password":"secret"}`,
		),
	)

	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200, got %d",
			rec.Code,
		)
	}

	cookies := rec.Result().Cookies()

	if len(cookies) == 0 {
		t.Fatal("expected session cookie")
	}

	var found bool

	for _, cookie := range cookies {
		if cookie.Name != SessionCookieName {
			continue
		}

		found = true

		if cookie.Value != "test-session-token" {
			t.Fatalf(
				"unexpected session token: %q",
				cookie.Value,
			)
		}

		if !cookie.HttpOnly {
			t.Fatal("session cookie must be HttpOnly")
		}

		if cookie.SameSite != http.SameSiteLaxMode {
			t.Fatal("session cookie must use SameSite=Lax")
		}
	}

	if !found {
		t.Fatal("session cookie not found")
	}
}

func TestLogout(t *testing.T) {
	service := &fakeAuthService{}

	handler := NewHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/auth/logout",
		nil,
	)

	req.AddCookie(&http.Cookie{
		Name:  SessionCookieName,
		Value: "test-session-token",
	})

	rec := httptest.NewRecorder()

	handler.Logout(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200, got %d",
			rec.Code,
		)
	}

	if service.logoutToken != "test-session-token" {
		t.Fatalf(
			"expected logout token to be passed through, got %q",
			service.logoutToken,
		)
	}

	var deleted bool

	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == SessionCookieName &&
			cookie.MaxAge < 0 {
			deleted = true
		}
	}

	if !deleted {
		t.Fatal("expected session cookie deletion")
	}
}

func TestRequireSession(t *testing.T) {
	service := &fakeAuthService{
		user: User{
			ID:       "user-123",
			Username: "customer01",
			Status:   "active",
		},
		loginToken: "valid-token",
	}

	handler := NewHandler(service)

	next := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			user := UserFromContext(r.Context())

			if user == nil {
				t.Fatal("expected authenticated user")
			}

			if user.ID != "user-123" {
				t.Fatalf(
					"unexpected user id: %q",
					user.ID,
				)
			}

			principal := authz.PrincipalFromContext(
				r.Context(),
			)

			if principal == nil {
				t.Fatal("expected authorization principal")
			}

			if principal.UserID != "user-123" {
				t.Fatalf(
					"unexpected principal user id: %q",
					principal.UserID,
				)
			}

			w.WriteHeader(http.StatusOK)
		},
	)

	protected := handler.RequireSession(next)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/customer",
		nil,
	)

	req.AddCookie(&http.Cookie{
		Name:  SessionCookieName,
		Value: "valid-token",
	})

	rec := httptest.NewRecorder()

	protected.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200, got %d",
			rec.Code,
		)
	}
}

func TestRequireSessionWithoutCookie(t *testing.T) {
	service := &fakeAuthService{}

	handler := NewHandler(service)

	next := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("next handler must not be called")
		},
	)

	protected := handler.RequireSession(next)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/customer",
		nil,
	)

	rec := httptest.NewRecorder()

	protected.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected 401, got %d",
			rec.Code,
		)
	}
}

func TestRequireSessionInvalidToken(t *testing.T) {
	service := &fakeAuthService{
		user: User{
			ID:       "user-123",
			Username: "customer01",
			Status:   "active",
		},
		loginToken: "valid-token",
	}

	handler := NewHandler(service)

	next := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("next handler must not be called")
		},
	)

	protected := handler.RequireSession(next)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/customer",
		nil,
	)

	req.AddCookie(&http.Cookie{
		Name:  SessionCookieName,
		Value: "invalid-token",
	})

	rec := httptest.NewRecorder()

	protected.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected 401, got %d",
			rec.Code,
		)
	}
}

func TestServiceErrorsAreNotExposed(t *testing.T) {
	service := &fakeAuthService{}

	handler := NewHandler(service)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/auth/login",
		strings.NewReader(
			`{"username":"wrong","password":"wrong"}`,
		),
	)

	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected 401, got %d",
			rec.Code,
		)
	}

	if strings.Contains(
		rec.Body.String(),
		errors.New("database failure").Error(),
	) {
		t.Fatal("internal error must not be exposed")
	}
}
