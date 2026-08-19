# A-Radius
ISP MANAGEMEN

## Menjalankan API dan PostgreSQL dengan Docker Compose

Stack lokal terdiri dari PostgreSQL, container migration one-shot, dan API Go. Salin template environment terlebih dahulu, kemudian jalankan Compose:

```bash
cp .env.example .env
docker compose up --build
```

Migration `database/migrations/0002_subscription_profiles.up.sql` dijalankan oleh service `migrate` setelah PostgreSQL berstatus sehat. API tersedia pada `http://localhost:8080`; PostgreSQL tersedia pada port `5432` dari host. Untuk menghentikan container tanpa menghapus data, gunakan `docker compose down`. Untuk menghapus database lokal dan menjalankan migration dari awal, gunakan `docker compose down -v`.

Contoh pemeriksaan API setelah container berjalan:

```bash
curl -H 'X-Tenant-ID: 00000000-0000-0000-0000-000000000001' \
     -H 'X-Actor-ID: 00000000-0000-0000-0000-000000000002' \
     'http://localhost:8080/api/v1/subscription-profiles'
```

Nilai password pada `.env.example` hanya untuk pengembangan lokal. Jangan menggunakan konfigurasi tersebut di staging atau production. Untuk deployment nyata, gunakan secret manager, TLS PostgreSQL, network policy, autentikasi production, dan migration runner yang terkontrol.
