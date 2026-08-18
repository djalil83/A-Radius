# Subscription Profile API

Package ini menyediakan CRUD profile berlangganan berbasis `database/sql` dan PostgreSQL `pgx/v5`, dengan router `chi/v5`, autentikasi JWT, dan RBAC server-side.

## Menjalankan

Migration `database/migrations/0002_subscription_profiles.up.sql` harus diterapkan terlebih dahulu. Konfigurasikan secret JWT melalui secret manager; jangan commit secret ke repository.

```bash
export DATABASE_URL='postgres://user:password@localhost:5432/a_radius?sslmode=verify-full'
export PROFILE_API_ADDR=':8080'
export JWT_SECRET='minimal-32-byte-secret-from-secret-manager'
export JWT_ISSUER='https://id.example.com/'
export JWT_AUDIENCE='a-radius-api'
go run ./cmd/profile-api
```

Token harus memakai algoritma `HS256` dan memuat claim wajib `sub`, `tenant_id`, `iss`, `aud`, `iat`, serta `exp`. Verifier memvalidasi signature, algorithm, issuer, audience, expiry, issued-at, dan keberadaan subject/tenant. Untuk deployment multi-service, gunakan provider asymmetric/JWKS yang dikelola terpusat daripada membagikan HMAC secret.

## Endpoint

| Method | Path | Permission |
|---|---|---|
| GET | `/api/v1/subscription-profiles` | `subscription_profiles.read` |
| POST | `/api/v1/subscription-profiles` | `subscription_profiles.create` |
| GET | `/api/v1/subscription-profiles/{id}` | `subscription_profiles.read` |
| PATCH | `/api/v1/subscription-profiles/{id}` | `subscription_profiles.update` |
| DELETE | `/api/v1/subscription-profiles/{id}?version=N` | `subscription_profiles.archive` |
| GET | `/api/v1/subscription-profiles/{id}/revisions` | `subscription_profiles.read_history` |

Kirim token sebagai `Authorization: Bearer <JWT>`. Identity tenant dan actor diambil dari Principal hasil verifikasi JWT; `X-Tenant-ID` dan `X-Actor-ID` tidak lagi digunakan oleh handler.

## Versioning

Update menggunakan `WHERE id = $id AND tenant_id = $tenant_id AND version = $version`, lalu menaikkan `version` satu angka. Jika versi sudah berubah, API mengembalikan `409 VERSION_CONFLICT`; frontend harus reload profile sebelum mencoba lagi. Database trigger juga menolak kenaikan versi yang bukan tepat satu angka. DELETE adalah soft delete menjadi `ARCHIVED`.

## Keamanan

Handler membatasi body JSON sampai 1 MiB dan menolak field JSON yang tidak dikenal. Repository memakai parameterized query. Endpoint tidak menerima atau mengembalikan credential RADIUS, MikroTik, OLT/ACS, atau payment gateway. RBAC dievaluasi server-side melalui `apb.user_roles`, `apb.role_permissions`, `apb.permissions`, dan status user aktif; kegagalan konfigurasi RBAC bersifat fail-closed.

Sebelum production, tambahkan TLS termination yang tepercaya, rate limiting, trusted-proxy policy, secret rotation, refresh-token policy di identity provider, dan integration test terhadap PostgreSQL nyata. JWT signature validation tidak menggantikan object-level authorization: repository tetap wajib memfilter seluruh query berdasarkan `tenant_id` dari Principal.

## Pengujian endpoint

Skrip `scripts/profile-api-smoke.sh` menguji bearer token yang hilang, create, detail, revision history, update, optimistic-lock conflict, dan archive. Skrip ini membutuhkan `curl` dan `jq`, serta JWT yang diterbitkan oleh identity provider dengan permission yang sesuai.

```bash
export BASE_URL='http://localhost:8080'
export JWT_TOKEN='eyJ...'
./scripts/profile-api-smoke.sh
```

Skrip Go `scripts/profile-api-smoke/main.go` menerbitkan JWT HS256 lokal dari `JWT_SECRET` dan menjalankan skenario yang sama. Gunakan hanya untuk local/staging; jangan gunakan secret smoke test sebagai secret production.

```bash
export DATABASE_URL='postgres://user:password@localhost:5432/a_radius?sslmode=disable'
export JWT_SECRET='minimal-32-byte-secret-from-secret-manager'
export JWT_ISSUER='a-radius'
export JWT_AUDIENCE='a-radius-api'
go run ./scripts/profile-api-smoke
```

Untuk JWT dari identity provider nyata, jalankan versi cURL dan setidaknya sediakan permission `subscription_profiles.read`, `subscription_profiles.create`, `subscription_profiles.update`, `subscription_profiles.archive`, serta `subscription_profiles.read_history` melalui RBAC database. UUID `sub` harus merujuk user aktif pada `apb.users`, sedangkan `tenant_id` harus sesuai tenant data yang diuji.
