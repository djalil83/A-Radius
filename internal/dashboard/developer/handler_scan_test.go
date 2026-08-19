package developer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestScanHTTPCreateAndGet(t *testing.T) {
	h := NewHandler()
	server := h.Routes()

	req := httptest.NewRequest(
		http.MethodPost,
		"/security/scans",
		strings.NewReader(`{"type":"full"}`),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d: %s", rec.Code, rec.Body.String())
	}

	var created ScanJob
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}

	if created.ID == "" {
		t.Fatal("expected scan ID")
	}

	if created.Status != ScanCompleted {
		t.Fatalf("expected completed, got %s", created.Status)
	}

	if len(created.Findings) != 1 {
		t.Fatalf("expected one finding, got %d", len(created.Findings))
	}

	getReq := httptest.NewRequest(
		http.MethodGet,
		"/security/scans/"+created.ID,
		nil,
	)
	getRec := httptest.NewRecorder()

	server.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", getRec.Code, getRec.Body.String())
	}

	var result ScanJob
	if err := json.NewDecoder(getRec.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	if result.ID != created.ID {
		t.Fatalf("expected ID %s, got %s", created.ID, result.ID)
	}
}

func TestScanHTTPInvalidType(t *testing.T) {
	h := NewHandler()
	server := h.Routes()

	req := httptest.NewRequest(
		http.MethodPost,
		"/security/scans",
		strings.NewReader(`{"type":"arbitrary-command"}`),
	)

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestScanHTTPNotFound(t *testing.T) {
	h := NewHandler()
	server := h.Routes()

	req := httptest.NewRequest(
		http.MethodGet,
		"/security/scans/does-not-exist",
		nil,
	)

	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
