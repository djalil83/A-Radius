# Security Knowledge Compare, Patch Separation, and Rollback

## Compare Version

Developer dapat membandingkan versi knowledge tanpa menganggapnya sebagai perubahan aplikasi Production. Contoh compare:

```text
SK-2.4.7 VS SK-2.4.8

ADDED
+ API-AUTH-042
+ SESSION-019
+ DEP-031

UPDATED
~ API-AUTH-018
~ DATABASE-007

REMOVED
- API-AUTH-003
```

Impact mapping untuk compare tersebut adalah `Administrator=HIGH`, `Langganan=HIGH`, `Pelanggan=MEDIUM`, `Technician=MEDIUM`, `Sales=LOW`, dan `Customer=LOW`.

## Knowledge tidak sama dengan patch

Knowledge version dapat berubah tanpa mengubah application runtime:

```text
SECURITY KNOWLEDGE
        ↓
AI ANALYSIS
        ↓
FINDING
        ↓
RECOMMENDATION
        ↓
PATCH PROPOSAL
        ↓
DEVELOPER PREVIEW
        ↓
APPROVAL
        ↓
STAGING
        ↓
PRODUCTION
```

`SK-2.4.8` hanya menjadi input analisis dan rule evaluation. Kode Production tetap menggunakan versi lama sampai Developer menyetujui patch proposal, security test dan staging test selesai, serta Production Review memberikan izin deployment.

## Rollback

Jika versi active menyebabkan false positive atau incompatibility, rollback dimulai sebagai proposal:

```text
ACTIVE SK-2.4.8
      ↓
Problem detected
      ↓
Rollback proposal
      ↓
Developer approval
      ↓
ACTIVE SK-2.4.7
```

Rollback harus menghasilkan audit event:

| Field | Nilai contoh |
|---|---|
| Action | `KNOWLEDGE_ROLLBACK` |
| From | `SK-2.4.8` |
| To | `SK-2.4.7` |
| Reason | `Compatibility issue` |
| Requested By | `Developer` |
| Approved By | `Developer` |
| Timestamp | `17/08/2026 14:42 WITA` |

Rollback tidak menghapus data versi lama. Versi lama dipertahankan untuk provenance, comparison, dan audit.

## Database blueprint

Migration `database/migrations/0003_security_knowledge.up.sql` menyediakan schema `security` dan tabel berikut:

| Tabel | Tujuan |
|---|---|
| `security_knowledge_versions` | Semver, status lifecycle, checksum, counters, confidence, dan timestamp lifecycle |
| `security_intelligence` | Finding intelligence, source, evidence, confidence, dan category |
| `security_rules` | Rule code, definition, severity, recommendation, dan enabled flag |
| `security_application_mapping` | Hubungan knowledge dengan module, component, risk, dan affected flag |
| `security_knowledge_approvals` | Approval untuk approve, stage, activate, rollback, dan archive |
| `security_knowledge_audit` | Audit append-oriented untuk perubahan status dan tindakan knowledge |

Constraint database memastikan hanya ada satu knowledge version berstatus `ACTIVE`, approval berstatus `APPROVED` harus memiliki `approved_by`, dan `production_changed` tidak dapat bernilai true pada lifecycle awal knowledge.

## Endpoint informasional

```text
GET /developer/security/knowledge/compare
GET /developer/security/knowledge/patch-pipeline
GET /developer/security/knowledge/rollback-audit
```

Endpoint tersebut hanya menyajikan compare, pipeline, dan contoh audit. Eksekusi rollback produksi tetap harus melalui approval, scope authorization, idempotency key, health check, dan audit.
