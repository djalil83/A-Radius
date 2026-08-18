# Subscription Profile API

Package ini menyediakan CRUD profile berlangganan berbasis `database/sql` dan PostgreSQL `pgx/v5`, dengan router `chi/v5`.

## Menjalankan

Migration `database/migrations/0002_subscription_profiles.up.sql` harus diterapkan terlebih dahulu. Setelah itu:

```bash
export DATABASE_URL='postgres://user:password@localhost:5432/a_radius?sslmode=disable'
export PROFILE_API_ADDR=':8080'
go run ./cmd/profile-api
```

## Endpoint

| Method | Path | Catatan |
|---|---|---|
| GET | `/api/v1/subscription-profiles` | Filter `q`, `service_type`, `status`, `limit`, `offset` |
| POST | `/api/v1/subscription-profiles` | Membuat profile baru |
| GET | `/api/v1/subscription-profiles/{id}` | Mengambil detail profile |
| PATCH | `/api/v1/subscription-profiles/{id}` | Wajib mengirim `version` dari respons terakhir |
| DELETE | `/api/v1/subscription-profiles/{id}?version=N` | Soft delete menjadi `ARCHIVED` |
| GET | `/api/v1/subscription-profiles/{id}/revisions` | Membaca snapshot version history |

Request harus menyertakan `X-Tenant-ID` dan `X-Actor-ID`. Header ini hanya adapter sementara untuk pengembangan; sebelum production, ganti dengan middleware autentikasi dan otorisasi A-Radius yang menghasilkan identity dari session/token terverifikasi.

## Versioning

Update menggunakan `WHERE id = $id AND tenant_id = $tenant_id AND version = $version`, lalu menaikkan `version` satu angka. Jika versi sudah berubah, API mengembalikan `409 VERSION_CONFLICT`; frontend harus reload profile sebelum mencoba lagi. Database trigger juga menolak kenaikan versi yang bukan tepat satu angka.

## Keamanan

Handler membatasi body JSON sampai 1 MiB dan menolak field JSON yang tidak dikenal. Repository memakai parameterized query. Endpoint tidak menerima atau mengembalikan credential RADIUS, MikroTik, OLT/ACS, atau payment gateway. Authz object-level dan rate limiting harus ditambahkan pada middleware production.
