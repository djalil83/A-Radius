package customerportal

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/djalil83/A-Radius/internal/authz"
)

type fakeDB struct{}

func (f *fakeDB) QueryRowContext(
	ctx context.Context,
	query string,
	args ...any,
) *sql.Row {
	return nil
}

func (f *fakeDB) QueryContext(
	ctx context.Context,
	query string,
	args ...any,
) (*sql.Rows, error) {
	return nil, nil
}

func TestAuthenticatedUserIDRequiresPrincipal(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/customer/me",
		nil,
	)

	if got, ok := authenticatedUserID(req); ok || got != "" {
		t.Fatalf(
			"expected no authenticated user, got %q, ok=%v",
			got,
			ok,
		)
	}
}

func TestAuthenticatedUserIDFromPrincipal(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/customer/me",
		nil,
	)

	principal := &authz.Principal{
		UserID:   "user-test-001",
		Username: "customer-test",
	}

	req = req.WithContext(
		authz.WithPrincipal(
			req.Context(),
			principal,
		),
	)

	got, ok := authenticatedUserID(req)

	if !ok {
		t.Fatal("expected authenticated user")
	}

	if got != "user-test-001" {
		t.Fatalf(
			"expected user-test-001, got %q",
			got,
		)
	}
}

func TestAuthenticatedUserIDIgnoresClientHeader(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/customer/me",
		nil,
	)

	req.Header.Set("X-Actor-ID", "attacker-controlled-id")

	if got, ok := authenticatedUserID(req); ok || got != "" {
		t.Fatalf(
			"client-controlled X-Actor-ID must not authenticate user: got %q, ok=%v",
			got,
			ok,
		)
	}
}

func TestMeRequiresAuthentication(t *testing.T) {
	handler := NewHandler(nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/customer/me",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.Me(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected HTTP 401, got %d",
			rec.Code,
		)
	}
}

func TestServicesRequiresAuthentication(t *testing.T) {
	handler := NewHandler(nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/customer/services",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.Services(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected HTTP 401, got %d",
			rec.Code,
		)
	}
}

func TestDashboardRequiresAuthentication(t *testing.T) {
	handler := NewHandler(nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/customer/dashboard",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.Dashboard(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected HTTP 401, got %d",
			rec.Code,
		)
	}
}
