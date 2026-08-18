package subscriptionprofile

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterRequiresIdentityHeaders(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/subscription-profiles", nil)
	res := httptest.NewRecorder()
	Router(NewHandler(nil), nil).ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", res.Code)
	}
	if got := res.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected JSON content type, got %q", got)
	}
}
