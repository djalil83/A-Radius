package authz

import (
"context"
"database/sql"
"encoding/json"
"net"
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
ip net.IP,
userAgent string,
) error {
if l == nil || l.DB == nil {
return nil
}

metadata, err := json.Marshal(map[string]any{
"permission": permission,
"allowed":    allowed,
"http_status": status,
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
