package authz

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestAuthorizationDecisionAuditIntegration(t *testing.T) {
	dsn := os.Getenv("A_RADIUS_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("A_RADIUS_TEST_DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()

	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("database connection failed: %v", err)
	}

	var userID string
	err = db.QueryRowContext(
		ctx,
		`SELECT id::text FROM apb.users ORDER BY created_at ASC LIMIT 1`,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("failed to find test user: %v", err)
	}

	principal := &Principal{
		UserID: userID,
	}

	logger := &DBAuditLogger{DB: db}

	testPath := fmt.Sprintf(
		"/api/test/users/%d",
		time.Now().UnixNano(),
	)

	req, err := http.NewRequest(
		http.MethodGet,
		testPath,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	req.RemoteAddr = "127.0.0.1:45678"
	req.Header.Set(
		"User-Agent",
		"A-Radius-RBAC-Integration-Test",
	)

	// ALLOW
	if err := logger.AuthorizationDecision(
		ctx,
		principal,
		"users.read",
		true,
		http.StatusOK,
		req,
	); err != nil {
		t.Fatalf("ALLOW audit failed: %v", err)
	}

	// DENY
	if err := logger.AuthorizationDecision(
		ctx,
		principal,
		"users.delete",
		false,
		http.StatusForbidden,
		req,
	); err != nil {
		t.Fatalf("DENY audit failed: %v", err)
	}

	rows, err := db.QueryContext(ctx, `
SELECT
metadata->>'permission',
(metadata->>'allowed')::boolean,
(metadata->>'http_status')::integer
FROM apb.audit_logs
WHERE action = 'authorization.decision'
  AND actor_id = $1
  AND metadata->>'path' = $2
ORDER BY created_at DESC
LIMIT 2
`, userID, testPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	type decision struct {
		permission string
		allowed    bool
		status     int
	}

	var decisions []decision

	for rows.Next() {
		var d decision

		if err := rows.Scan(
			&d.permission,
			&d.allowed,
			&d.status,
		); err != nil {
			t.Fatal(err)
		}

		decisions = append(decisions, d)
	}

	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}

	if len(decisions) != 2 {
		t.Fatalf(
			"expected 2 authorization audit records, got %d",
			len(decisions),
		)
	}

	// Newest = DENY
	if decisions[0].permission != "users.delete" ||
		decisions[0].allowed ||
		decisions[0].status != http.StatusForbidden {
		t.Fatalf(
			"invalid DENY record: %+v",
			decisions[0],
		)
	}

	// Previous = ALLOW
	if decisions[1].permission != "users.read" ||
		!decisions[1].allowed ||
		decisions[1].status != http.StatusOK {
		t.Fatalf(
			"invalid ALLOW record: %+v",
			decisions[1],
		)
	}

	t.Log("ALLOW audit: users.read -> allowed=true, HTTP 200")
	t.Log("DENY audit: users.delete -> allowed=false, HTTP 403")
}
