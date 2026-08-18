# A-Radius Web UI

`profile-berlangganan.html` adalah UI Profile Berlangganan A-Radius yang terhubung ke API Go melalui `profile-api.js` dan Fetch API.

## Integrasi API

| Operasi UI | Endpoint |
|---|---|
| Muat daftar | `GET /api/v1/subscription-profiles` |
| Tambah | `POST /api/v1/subscription-profiles` |
| Edit | `PATCH /api/v1/subscription-profiles/{id}` dengan `version` |
| Arsipkan | `DELETE /api/v1/subscription-profiles/{id}?version=N` |
| History | `GET /api/v1/subscription-profiles/{id}/revisions` |

Mapping API dilakukan di `profile-api.js`. Field UUID dipertahankan sebagai string, `service_type` dipetakan ke label UI, serta nilai bit-per-second dikonversi ke label kecepatan.

## Konfigurasi

Sebelum `profile-api.js` dimuat, konfigurasi harus tersedia melalui JavaScript atau konfigurasi server yang aman:

```html
<script>
  window.ARADIUS_API_BASE = 'http://localhost:8080';
  window.ARADIUS_TENANT_ID = '00000000-0000-0000-0000-000000000001';
  window.ARADIUS_ACTOR_ID = '00000000-0000-0000-0000-000000000002';
</script>
<script src="profile-api.js"></script>
```

Pada `profile-berlangganan.html`, `profile-api.js` sudah dimuat oleh halaman. Dalam production, jangan menanam UUID actor statis. Ganti adapter header `X-Tenant-ID` dan `X-Actor-ID` dengan middleware autentikasi/session atau token yang telah diverifikasi.

## Version conflict

Setiap profile yang dimuat menyimpan `version`. Saat menyimpan atau mengarsipkan, versi tersebut dikirim kembali ke API. Jika API menjawab `409 VERSION_CONFLICT`, UI memuat ulang daftar profile dan meminta pengguna meninjau perubahan terbaru.

## Menjalankan secara lokal

Dari root proyek, jalankan API dan server file statis secara terpisah:

```bash
export DATABASE_URL='postgres://user:password@localhost:5432/a_radius?sslmode=require'
export PROFILE_API_ADDR=':8080'
go run ./cmd/profile-api

python3 -m http.server 8081 --directory web
```

Buka `http://localhost:8081/profile-berlangganan.html`. Pastikan API menyediakan CORS jika UI dan API berada pada origin berbeda. Jangan menaruh credential database, RADIUS, MikroTik, OLT, atau payment gateway di HTML maupun JavaScript frontend.
