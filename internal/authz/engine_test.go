package authz

import (
"context"
"database/sql"
"testing"
)

type fakeQuerier struct {
allowed bool
err     error
}

func (f fakeQuerier) QueryRowContext(
ctx context.Context,
query string,
args ...any,
) *sql.Row {
// Tests that require real PostgreSQL are intentionally kept separate.
// This placeholder ensures the public Engine API remains compile-safe.
return nil
}

func TestErrors(t *testing.T) {
if ErrUnauthenticated == nil {
t.Fatal("ErrUnauthenticated must exist")
}

if ErrForbidden == nil {
t.Fatal("ErrForbidden must exist")
}
}
