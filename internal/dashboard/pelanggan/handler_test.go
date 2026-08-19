package pelanggan

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/djalil83/A-Radius/internal/customerportal"
)

func TestRoutesWithoutPortalReturn503(t *testing.T) {
	handler := NewHandler(nil)
	router := handler.Routes()

	req := httptest.NewRequest(
		http.MethodGet,
		"/dashboard",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf(
			"expected HTTP 503, got %d",
			rec.Code,
		)
	}
}

func TestRoutesServeDashboardUI(t *testing.T) {
	portal := customerportal.NewHandler(nil)
	handler := NewHandler(portal)
	router := handler.Routes()

	req := httptest.NewRequest(
		http.MethodGet,
		"/",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected HTTP 200 for dashboard UI, got %d",
			rec.Code,
		)
	}

	if contentType := rec.Header().Get("Content-Type"); !strings.HasPrefix(
		contentType,
		"text/html",
	) {
		t.Fatalf(
			"expected HTML content type, got %q",
			contentType,
		)
	}

	if !strings.Contains(rec.Body.String(), "A-Radius") {
		t.Fatalf(
			"expected dashboard HTML to contain A-Radius",
		)
	}
}
