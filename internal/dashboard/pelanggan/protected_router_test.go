package pelanggan

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/djalil83/A-Radius/internal/authz"
)

func TestProtectedRouterNilHandler(t *testing.T) {
	engine := authz.NewEngine(nil)

	handler := ProtectedRouter(nil, engine, nil)

	req := httptest.NewRequest(http.MethodGet, "/dashboard/pelanggan/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusServiceUnavailable,
			rec.Code,
		)
	}
}

func TestProtectedRouterNilEngine(t *testing.T) {
	handler := ProtectedRouter(nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/dashboard/pelanggan/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusServiceUnavailable,
			rec.Code,
		)
	}
}
