//go:build integration

package subscriptionprofile

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	integrationTenantID = "00000000-0000-0000-0000-000000000201"
	integrationActorID  = "00000000-0000-0000-0000-000000000202"
)

type concurrencyResult struct {
	status int
	body   []byte
}

func TestIntegration_ConcurrentUpdatesHaveSingleWinner(t *testing.T) {
	db := prepareIntegrationDB(t)
	cleanupIntegrationTenant(t, db)
	defer cleanupIntegrationTenant(t, db)

	router := Router(NewHandler(NewRepository(db)))
	profileID, version := createIntegrationProfile(t, router)
	const workers = 32

	start := make(chan struct{})
	results := make(chan concurrencyResult, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			<-start
			payload := integrationUpdatePayload(version, fmt.Sprintf("worker-%02d", workerID))
			results <- callJSON(router, http.MethodPatch, "/api/v1/subscription-profiles/"+profileID, payload, true)
		}(i)
	}
	close(start)
	wg.Wait()
	close(results)

	var success, conflicts int
	var winner Profile
	for result := range results {
		switch result.status {
		case http.StatusOK:
			success++
			if err := json.Unmarshal(result.body, &winner); err != nil {
				t.Fatalf("decode winner response: %v", err)
			}
		case http.StatusConflict:
			conflicts++
		default:
			t.Fatalf("unexpected concurrent update status %d: %s", result.status, result.body)
		}
	}
	if success != 1 || conflicts != workers-1 {
		t.Fatalf("expected one winner and %d conflicts, got success=%d conflicts=%d", workers-1, success, conflicts)
	}
	if winner.Version != version+1 {
		t.Fatalf("expected winner version %d, got %d", version+1, winner.Version)
	}

	final := getIntegrationProfile(t, router, profileID)
	if final.Version != version+1 || final.Name != winner.Name {
		t.Fatalf("lost update detected: final=%+v winner=%+v", final, winner)
	}
	assertRevisionCount(t, db, profileID, 2)
	assertAuditCount(t, db, profileID, 2)
}

func TestIntegration_ConcurrentUpdateAndArchiveHaveSingleWinner(t *testing.T) {
	db := prepareIntegrationDB(t)
	cleanupIntegrationTenant(t, db)
	defer cleanupIntegrationTenant(t, db)

	router := Router(NewHandler(NewRepository(db)))
	profileID, version := createIntegrationProfile(t, router)
	start := make(chan struct{})
	results := make(chan concurrencyResult, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		<-start
		results <- callJSON(router, http.MethodPatch, "/api/v1/subscription-profiles/"+profileID, integrationUpdatePayload(version, "update-race"), true)
	}()
	go func() {
		defer wg.Done()
		<-start
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/subscription-profiles/"+profileID+"?version=1", nil)
		setIntegrationHeaders(req)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		results <- concurrencyResult{status: resp.Code, body: resp.Body.Bytes()}
	}()
	close(start)
	wg.Wait()
	close(results)

	var success, conflicts int
	for result := range results {
		switch result.status {
		case http.StatusOK, http.StatusNoContent:
			success++
		case http.StatusConflict:
			conflicts++
		default:
			t.Fatalf("unexpected update/archive status %d: %s", result.status, result.body)
		}
	}
	if success != 1 || conflicts != 1 {
		t.Fatalf("expected one winner and one conflict, got success=%d conflicts=%d", success, conflicts)
	}

	var versionAfter int64
	var status string
	err := db.QueryRowContext(context.Background(), `SELECT version,status FROM subscription_profiles WHERE tenant_id=$1 AND id=$2`, integrationTenantID, profileID).Scan(&versionAfter, &status)
	if err != nil {
		t.Fatalf("read final profile: %v", err)
	}
	if versionAfter != 2 {
		t.Fatalf("expected one successful mutation to produce version 2, got %d", versionAfter)
	}
	if status != "ACTIVE" && status != "ARCHIVED" {
		t.Fatalf("unexpected final status %q", status)
	}
	assertRevisionCount(t, db, profileID, 2)
}

var integrationSchemaOnce sync.Once
var integrationSchemaErr error

func prepareIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("INTEGRATION_DATABASE_URL")
	if dsn == "" {
		t.Skip("INTEGRATION_DATABASE_URL is not set; start disposable PostgreSQL first")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open integration database: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		t.Fatalf("ping integration database: %v", err)
	}
	integrationSchemaOnce.Do(func() {
		migrationPath := filepath.Join("..", "..", "database", "migrations", "0002_subscription_profiles.up.sql")
		migration, err := os.ReadFile(migrationPath)
		if err != nil {
			integrationSchemaErr = fmt.Errorf("read migration %s: %w", migrationPath, err)
			return
		}
		_, integrationSchemaErr = db.ExecContext(context.Background(), string(migration))
	})
	if integrationSchemaErr != nil {
		db.Close()
		t.Fatalf("apply integration migration: %v", integrationSchemaErr)
	}
	return db
}

func cleanupIntegrationTenant(t *testing.T, db *sql.DB) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, query := range []string{
		"DELETE FROM audit_events WHERE tenant_id=$1",
		"DELETE FROM subscription_profile_revisions WHERE tenant_id=$1",
		"DELETE FROM approval_requests WHERE tenant_id=$1",
		"DELETE FROM subscription_profiles WHERE tenant_id=$1",
	} {
		if _, err := db.ExecContext(ctx, query, integrationTenantID); err != nil {
			// The first cleanup runs before migration on a fresh database.
			if !isUndefinedTable(err) {
				t.Fatalf("cleanup tenant: %v", err)
			}
			return
		}
	}
}

func isUndefinedTable(err error) bool {
	return err != nil && (bytes.Contains([]byte(err.Error()), []byte("does not exist")) || bytes.Contains([]byte(err.Error()), []byte("undefined_table")))
}

func createIntegrationProfile(t *testing.T, router http.Handler) (string, int64) {
	t.Helper()
	payload := map[string]any{
		"name":              "integration-race-profile",
		"service_type":      "FTTH",
		"color":             "#2563EB",
		"shared_users":      1,
		"monthly_price":     150000,
		"active_days":       30,
		"commission_amount": 0,
		"commission_type":   "RUPIAH",
		"billing_cycle":     "MONTHLY",
		"auto_isolate":      false,
	}
	result := callJSON(router, http.MethodPost, "/api/v1/subscription-profiles", payload, false)
	if result.status != http.StatusCreated {
		t.Fatalf("create profile: status=%d body=%s", result.status, result.body)
	}
	var profile Profile
	if err := json.Unmarshal(result.body, &profile); err != nil {
		t.Fatalf("decode created profile: %v", err)
	}
	return profile.ID, profile.Version
}

func getIntegrationProfile(t *testing.T, router http.Handler, id string) Profile {
	t.Helper()
	result := callJSON(router, http.MethodGet, "/api/v1/subscription-profiles/"+id, nil, false)
	if result.status != http.StatusOK {
		t.Fatalf("get profile: status=%d body=%s", result.status, result.body)
	}
	var profile Profile
	if err := json.Unmarshal(result.body, &profile); err != nil {
		t.Fatalf("decode profile: %v", err)
	}
	return profile
}

func integrationUpdatePayload(version int64, suffix string) map[string]any {
	return map[string]any{
		"name":              "integration-" + suffix,
		"service_type":      "FTTH",
		"color":             "#2563EB",
		"shared_users":      1,
		"monthly_price":     150000,
		"active_days":       30,
		"commission_amount": 0,
		"commission_type":   "RUPIAH",
		"billing_cycle":     "MONTHLY",
		"auto_isolate":      false,
		"version":           version,
	}
}

func callJSON(router http.Handler, method, path string, payload any, includeVersion bool) concurrencyResult {
	var body bytes.Buffer
	if payload != nil {
		_ = json.NewEncoder(&body).Encode(payload)
	}
	req := httptest.NewRequest(method, path, &body)
	if includeVersion {
		req.Header.Set("Content-Type", "application/json")
	}
	setIntegrationHeaders(req)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return concurrencyResult{status: resp.Code, body: append([]byte(nil), resp.Body.Bytes()...)}
}

func setIntegrationHeaders(req *http.Request) {
	req.Header.Set("X-Tenant-ID", integrationTenantID)
	req.Header.Set("X-Actor-ID", integrationActorID)
}

func assertRevisionCount(t *testing.T, db *sql.DB, profileID string, expected int) {
	t.Helper()
	var got int
	if err := db.QueryRow(`SELECT count(*) FROM subscription_profile_revisions WHERE tenant_id=$1 AND profile_id=$2`, integrationTenantID, profileID).Scan(&got); err != nil {
		t.Fatalf("count revisions: %v", err)
	}
	if got != expected {
		t.Fatalf("expected %d revisions, got %d", expected, got)
	}
}

func assertAuditCount(t *testing.T, db *sql.DB, profileID string, expected int) {
	t.Helper()
	var got int
	if err := db.QueryRow(`SELECT count(*) FROM audit_events WHERE tenant_id=$1 AND entity_id=$2`, integrationTenantID, profileID).Scan(&got); err != nil {
		t.Fatalf("count audit events: %v", err)
	}
	if got != expected {
		t.Fatalf("expected %d audit events, got %d", expected, got)
	}
}
