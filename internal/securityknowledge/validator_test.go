package securityknowledge

import (
	"strings"
	"testing"
)

func TestValidateHash(t *testing.T) {
	valid := strings.Repeat("a", 64)

	if !ValidateHash(valid) {
		t.Fatal("expected valid SHA-256 hash")
	}

	if ValidateHash("invalid") {
		t.Fatal("expected invalid hash to fail")
	}

	if ValidateHash(strings.Repeat("g", 64)) {
		t.Fatal("expected non-hex hash to fail")
	}
}

func TestValidateManifest(t *testing.T) {
	data := []byte(`{
"knowledge_key": "a-radius.security",
"version": "1.0.0",
"status": "active",
"scope": [
"authentication",
"authorization",
"rbac",
"http-401-403",
"audit",
"secure-development"
],
"rules": {
"authorization_default": "deny",
"unauthenticated_http": 401,
"authenticated_without_permission_http": 403,
"ai_must_not_bypass_rbac": true,
"authorization_must_be_server_side": true,
"authorization_decisions_must_be_audited": true
}
}`)

	manifest, err := ValidateManifest(data)
	if err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}

	if manifest.KnowledgeKey != "a-radius.security" {
		t.Fatalf("unexpected knowledge key: %s", manifest.KnowledgeKey)
	}

	if manifest.Version != "1.0.0" {
		t.Fatalf("unexpected version: %s", manifest.Version)
	}
}

func TestValidateManifestRejectsInvalidVersion(t *testing.T) {
	data := []byte(`{
"knowledge_key": "a-radius.security",
"version": "1.0",
"status": "active",
"scope": ["authorization"],
"rules": {
"authorization_default": "deny"
}
}`)

	if _, err := ValidateManifest(data); err != ErrInvalidVersion {
		t.Fatalf("expected ErrInvalidVersion, got %v", err)
	}
}

func TestStatusValid(t *testing.T) {
	tests := []struct {
		status Status
		valid  bool
	}{
		{StatusDraft, true},
		{StatusActive, true},
		{StatusDeprecated, true},
		{StatusRevoked, true},
		{Status("unknown"), false},
	}

	for _, tt := range tests {
		if got := tt.status.Valid(); got != tt.valid {
			t.Fatalf(
				"status %q: expected %v, got %v",
				tt.status,
				tt.valid,
				got,
			)
		}
	}
}
