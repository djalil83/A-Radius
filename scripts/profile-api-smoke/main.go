package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	TenantID          string `json:"tenant_id"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	jwt.RegisteredClaims
}

type profile struct {
	ID      string `json:"id"`
	Version int64  `json:"version"`
}

func main() {
	baseURL := strings.TrimRight(env("BASE_URL", "http://localhost:8080"), "/")
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		fatal("JWT_SECRET must be at least 32 bytes")
	}
	issuer, audience := env("JWT_ISSUER", "a-radius"), env("JWT_AUDIENCE", "a-radius-api")
	userID, tenantID := env("TEST_USER_ID", "00000000-0000-0000-0000-000000000001"), env("TEST_TENANT_ID", "00000000-0000-0000-0000-000000000002")
	now := time.Now()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		TenantID: tenantID, PreferredUsername: "api-smoke-test",
		RegisteredClaims: jwt.RegisteredClaims{Subject: userID, Issuer: issuer, Audience: jwt.ClaimStrings{audience}, IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute))},
	}).SignedString([]byte(secret))
	if err != nil {
		fatal("sign JWT: %v", err)
	}

	status, _, err := call(http.MethodGet, baseURL+"/api/v1/subscription-profiles", token, nil)
	if err != nil {
		fatal("unauthenticated request: %v", err)
	}
	if status != http.StatusUnauthorized {
		fatal("expected unauthenticated request to return 401, got %d", status)
	}

	name := "Smoke Plan " + strconv.FormatInt(time.Now().UnixNano(), 10)
	create := map[string]any{"name": name, "service_type": "FTTH", "category": "HOME", "media": "FIBER", "color": "#2563EB", "description": "Temporary API smoke-test profile", "mikrotik_group": "default", "radius_group": "ftth", "rate_limit": "20M/20M", "upload_bps": 20000000, "download_bps": 20000000, "shared_users": 1, "monthly_price": 150000, "active_days": 30, "commission_amount": 0, "commission_type": "FIXED", "billing_cycle": "MONTHLY", "auto_isolate": true}
	status, body, err := call(http.MethodPost, baseURL+"/api/v1/subscription-profiles", token, create)
	if err != nil || status != http.StatusCreated {
		fatal("create failed: status=%d err=%v body=%s", status, err, body)
	}
	var created profile
	decode(body, &created)
	if created.ID == "" || created.Version < 1 {
		fatal("create response missing id/version: %s", body)
	}

	for _, path := range []string{"/api/v1/subscription-profiles/" + created.ID, "/api/v1/subscription-profiles/" + created.ID + "/revisions"} {
		status, body, err = call(http.MethodGet, baseURL+path, token, nil)
		if err != nil || status != http.StatusOK {
			fatal("GET %s failed: status=%d err=%v body=%s", path, status, err, body)
		}
	}

	update := map[string]any{"version": created.Version, "name": name + " Updated", "service_type": "FTTH", "color": "#1D4ED8", "monthly_price": 160000, "active_days": 30, "commission_amount": 0, "commission_type": "FIXED", "billing_cycle": "MONTHLY", "auto_isolate": true}
	status, body, err = call(http.MethodPatch, baseURL+"/api/v1/subscription-profiles/"+created.ID, token, update)
	if err != nil || status != http.StatusOK {
		fatal("update failed: status=%d err=%v body=%s", status, err, body)
	}
	var updated profile
	decode(body, &updated)
	if updated.Version != created.Version+1 {
		fatal("expected version %d, got %d", created.Version+1, updated.Version)
	}

	stale := map[string]any{"version": created.Version, "name": "Stale Update", "service_type": "FTTH", "color": "#1D4ED8", "monthly_price": 160000, "active_days": 30, "commission_amount": 0, "commission_type": "FIXED", "billing_cycle": "MONTHLY", "auto_isolate": true}
	status, body, err = call(http.MethodPatch, baseURL+"/api/v1/subscription-profiles/"+created.ID, token, stale)
	if err != nil || status != http.StatusConflict {
		fatal("stale update should return 409: status=%d err=%v body=%s", status, err, body)
	}

	status, body, err = call(http.MethodDelete, baseURL+"/api/v1/subscription-profiles/"+created.ID+"?version="+strconv.FormatInt(updated.Version, 10), token, nil)
	if err != nil || status != http.StatusNoContent {
		fatal("archive failed: status=%d err=%v body=%s", status, err, body)
	}
	fmt.Printf("profile API smoke test passed for %s\n", created.ID)
}

func call(method, url, token string, payload any) (int, string, error) {
	var reader io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return 0, "", err
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if err != nil {
		return 0, "", err
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	return res.StatusCode, string(raw), err
}

func decode(body string, dst any) {
	if err := json.Unmarshal([]byte(body), dst); err != nil {
		fatal("decode response: %v; body=%s", err, body)
	}
}
func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func fatal(format string, args ...any) { fmt.Fprintf(os.Stderr, format+"\n", args...); os.Exit(1) }
