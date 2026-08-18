# Local Docker Compose

Konfigurasi ini menjalankan PostgreSQL dan `profile-api` secara lokal. PostgreSQL akan menginisialisasi `0001_init.sql`, migration profile `0002_subscription_profiles.up.sql`, seed RBAC smoke test, dan system seed pada database baru.

## Menjalankan

```bash
cp .env.example .env
docker compose up --build
```

Health API dapat diperiksa dengan:

```bash
curl http://localhost:8080/healthz
```

Respons yang diharapkan adalah `{"status":"ok"}`. Untuk menjalankan smoke test end-to-end menggunakan JWT HS256 lokal:

```bash
export BASE_URL=http://localhost:8080
export JWT_SECRET=local-development-secret-change-me-32-bytes-minimum
go run ./scripts/profile-api-smoke
```

Smoke test Go memakai user UUID `00000000-0000-0000-0000-000000000001` dan tenant UUID default `00000000-0000-0000-0000-000000000002`. Seed tersebut hanya untuk local/staging, bukan identity bootstrap production.

## Reset database

Skrip pada `/docker-entrypoint-initdb.d` hanya dijalankan ketika volume PostgreSQL masih kosong. Untuk mengulang bootstrap dari awal:

```bash
docker compose down -v
docker compose up --build
```

Perintah `down -v` menghapus seluruh data lokal PostgreSQL. Jangan menjalankannya terhadap environment yang berisi data penting.

## Catatan production

Ganti password database dan `JWT_SECRET` melalui secret manager. Gunakan issuer identity provider nyata, signing key asymmetric/JWKS bila token dipakai lintas service, TLS, backup, serta migration runner terpisah. Seed RBAC `03-rbac-smoke.sql` harus dikeluarkan dari deployment production.
