# Developer Security Dashboard

## Ruang lingkup

Developer Security Dashboard adalah command center untuk tim platform A-Radius. Dashboard ini menggabungkan observability keamanan, analisis AI, pemeriksaan dependency, pengendalian akses, approval perubahan, dan operasi produksi dalam satu area yang hanya dapat diakses role `developer` dengan permission yang sesuai.

Struktur navigasi yang diberikan pengguna telah diterjemahkan ke `web/dashboards/developer/dashboard.js`. Menu UI bersifat katalog kapabilitas; setiap aksi nyata tetap harus divalidasi ulang di backend melalui permission RBAC.

## Pemetaan menu dan permission

| Area | Contoh fungsi | Permission utama | Approval |
|---|---|---|---|
| AI Engine | Code, architecture, database, configuration, dan security analyzer | `code:read`, `database:read`, `security:read` | Untuk analisis berisiko tinggi |
| Security Overview | Security score dan temuan berdasarkan level | `security:read` | Tidak untuk baca |
| AI Security Checker | Full scan, code, API, database, auth, session, dependency, dan infrastructure scan | `security:scan` | Ya, terutama scan yang membuat perubahan atau mengakses produksi |
| Threat Detection | Login mencurigakan, brute force, token abuse, dan anomaly | `threat:read` | Untuk tindakan mitigasi |
| Access Security | Developer, administrator, technician, sales, reseller, customer | `security:read` | Perubahan akses memakai `approval:decide` |
| Credential Security | Password policy, API key, token, secret, SSH key, exposure | `credential:read` dan `credential:rotate` | Rotasi/revoke wajib approval sesuai kebijakan |
| Dependency Security | Vulnerable/outdated package, CVE, license, supply chain | `dependency:read` dan `dependency:fix` | Perbaikan dependency produksi wajib approval |
| API Security | Auth, authz, rate limit, validation, CORS, webhook, exposure | `api:read` | Perubahan konfigurasi memakai approval |
| Database Security | SQL injection, privilege, credential, backup, sensitive data, query anomaly | `database:read` | Perubahan database wajib approval dan migration review |
| Security Report | Findings, evidence, risk, module, recommendation, history | `security:read` | Tidak untuk baca |
| Security Audit Trail | Login, permission, configuration, deployment, finding, approval | `audit:read` | Audit bersifat append-only |
| Development | Code, database, API, config, preview, test, security test | `code:read`, `code:write`, `preview:read` | Code/config/database production wajib approval |
| Production | Release, deployment, rollback, version, health check | `deployment:read`, `deployment:run`, `deployment:rollback` | Ya |

## Kontrak backend

`internal/dashboard/developer/security_contract.go` menyediakan `FeatureContract` yang mengikat label menu dengan permission, risk level, dan penanda `RequiresApproval`. Ini mencegah UI menjadi satu-satunya sumber aturan keamanan.

`internal/dashboard/developer/handler.go` menyediakan endpoint awal berikut:

| Endpoint | Fungsi | Status |
|---|---|---|
| `GET /` | Health/status dashboard | Aktif |
| `GET /security/overview` | Ringkasan score dan jumlah temuan | Kerangka aktif |
| `GET /security/features` | Daftar kontrak kapabilitas | Aktif |
| `POST /security/scans` | Menjadwalkan full scan | Menerima job dan mengembalikan `202 Accepted`; worker scanner belum dipasang |

Handler tersebut harus dipasang di bawah JWT middleware dan permission middleware. Contoh permission untuk `POST /security/scans` adalah `security:scan`. Jangan memberikan `deployment:run`, `credential:rotate`, atau `approval:decide` secara otomatis hanya karena pengguna memiliki akses ke dashboard developer.

## Alur full scan

```text
Developer request
    -> JWT validation
    -> Principal injection
    -> RBAC: security:scan
    -> Audit decision
    -> Scan job queue
    -> Scanner workers
    -> Evidence storage
    -> Findings + risk scoring
    -> Approval jika remediation mengubah sistem
    -> Audit trail
```

Implementasi saat ini berhenti pada pembuatan job response `queued`. Scanner nyata harus ditambahkan sebagai worker asynchronous, dengan batas waktu, idempotency key, pengamanan akses repository, evidence retention, dan isolasi dari credential produksi.

## Kontrol produksi

Perubahan di area Production tidak boleh langsung dilakukan dari browser. Backend perlu menerapkan approval gate, perubahan immutable release, health check setelah deployment, rollback yang tercatat, serta audit event untuk actor, target, versi, waktu, dan hasil. AI hanya boleh memberi analisis atau rekomendasi; AI tidak boleh memperoleh permission deployment secara implisit.
