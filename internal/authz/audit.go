package authz

import (
	"context"
	"database/sql"
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

type DBAuditLogger struct {
	DB *sql.DB
}

func (l *DBAuditLogger) AuthorizationDecision(
	ctx context.Context,
	principal *Principal,
	permission string,
	allowed bool,
	status int,
	r *http.Request,
) error {
	if l == nil || l.DB == nil {
		return nil
	}

	var ip net.IP
	var userAgent string
	var method string
	var path string

	if r != nil {
		ip = clientIP(r)
		userAgent = r.UserAgent()
		method = r.Method
		path = r.URL.Path
	}

	metadata, err := json.Marshal(map[string]any{
		"permission":  permission,
		"allowed":     allowed,
		"http_status": status,
		"method":      method,
		"path":        path,
	})
	if err != nil {
		return err
	}

	var actorID any
	if principal != nil && principal.UserID != "" {
		actorID = principal.UserID
	}

	_, err = l.DB.ExecContext(
		ctx,
		`INSERT INTO apb.audit_logs
(actor_id, action, resource_type, ip_address, user_agent, metadata)
VALUES
($1, 'authorization.decision', 'permission', $2, $3, $4::jsonb)`,
		actorID,
		ip,
		userAgent,
		string(metadata),
	)

	return err
}

func clientIP(r *http.Request) net.IP {
	if r == nil {
		return nil
	}

	// Do not blindly trust X-Forwarded-For.
	// Trusted proxy handling must be configured explicitly at the edge.
	host := r.RemoteAddr

	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	if ip := net.ParseIP(strings.TrimSpace(host)); ip != nil {
		return ip
	}

	return nil
}
