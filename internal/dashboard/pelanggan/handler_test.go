package pelanggan

import (
	"net/http"
	"net/http/httptest"
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

func TestRoutesWithPortalRequireAuthentication(t *testing.T) {
	portal := customerportal.NewHandler(nil)
	handler := NewHandler(portal)
	router := handler.Routes()

	tests := []string{
		"/",
		"/dashboard",
		"/me",
		"/services",
	}

	for _, path := range tests {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				path,
				nil,
			)

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Fatalf(
					"expected HTTP 401 for %s, got %d",
					path,
					rec.Code,
				)
			}
		})
	}
}
