# A-RADIUS AI Knowledge–Production Separation

## Tujuan

AI Knowledge adalah sumber pengetahuan dan analisis keamanan, bukan bagian dari application runtime dan bukan jalur perubahan otomatis ke Production. Knowledge dapat diperbarui secara berkala tanpa mengubah kode, konfigurasi, database, credential, atau layanan pelanggan.

## Boundary arsitektur

```text
INTERNET
   ↓
SECURITY KNOWLEDGE
   ↓
AI KNOWLEDGE DB
   ├── ANALYSIS CACHE
   └── RULE ENGINE
          ↓
     A-RADIUS SCAN
          ↓
     FINDING / REPORT
          ↓
  DEVELOPER APPROVAL
          ↓
       STAGING
          ↓
   PRODUCTION REVIEW
          ↓
      PRODUCTION
```

AI Knowledge DB dan analysis cache harus menggunakan environment, credential, network policy, dan storage permission yang terpisah dari Production. Rule engine hanya menghasilkan rule evaluation, finding, dan evidence. Rule engine tidak boleh memiliki capability deployment.

## Knowledge versioning

Setiap knowledge item disimpan sebagai versi immutable, misalnya `v1.0`, `v1.1`, `v1.2`, dan `v1.3`. Versi baru dapat berstatus `NEW`, kemudian berubah melalui proses review menjadi `REVIEWED`, `APPROVED`, `REJECTED`, atau `ARCHIVED`. Status knowledge tidak sama dengan status release aplikasi.

Metadata minimum yang wajib disimpan adalah sumber, waktu ditemukan, tingkat kepercayaan, relevansi terhadap A-RADIUS, modul terdampak, hasil analisis, recommendation, content hash, parser version, dan status lifecycle.

| Status | Makna | Dapat mengubah Production? |
|---|---|---|
| `NEW` | Knowledge baru diterima dan belum ditinjau | Tidak |
| `REVIEWED` | Evidence dan relevansi telah diperiksa | Tidak |
| `APPROVED` | Knowledge disetujui untuk dipakai sebagai input analysis/rule | Tidak |
| `REJECTED` | Knowledge tidak cukup valid atau tidak relevan | Tidak |
| `ARCHIVED` | Versi lama disimpan untuk provenance dan audit | Tidak |

## Promotion policy

Policy default menetapkan `auto_production_promotion=false`, `knowledge_db_isolated=true`, `analysis_cache_read_only=true`, dan `required_developer_approval=true`. Perubahan aplikasi harus melewati finding, developer approval, staging, production review, deployment terkontrol, dan health check.

AI boleh mempelajari knowledge baru, membandingkannya dengan Code/API/DB A-RADIUS, membuat finding, recommendation, patch preview, dan security test. AI tidak boleh mengeksekusi deployment, mengubah firewall, menghapus credential, memblokir akun, memutus API, atau menerapkan migrasi Production.

## Endpoint

```text
GET /developer/security/knowledge/policy
GET /developer/security/knowledge/versions
GET /developer/security/knowledge/featured
```

Endpoint tersebut bersifat informasional dan tidak menjalankan promosi. Endpoint patch, approval, dan deployment tetap harus dilindungi JWT, RBAC, approval record, scope, evidence test, idempotency key, dan audit trail.

## Audit dan observability

Setiap ingestion, pembuatan versi, perubahan status, analisis, recommendation, patch preview, approval, staging test, production review, deployment, health check, dan rollback harus menghasilkan audit event dengan actor, timestamp, source, version, correlation ID, dan hasil tindakan.
