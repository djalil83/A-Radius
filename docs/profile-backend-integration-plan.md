# Rencana Integrasi Backend Profile Berlangganan

**Proyek:** A-Radius  
**Sasaran:** Menggantikan penyimpanan `localStorage` pada `web/profile-berlangganan.html` dengan API Go dan PostgreSQL tanpa mengaktifkan perubahan network atau billing secara langsung.

## 1. Keputusan arsitektur

Gunakan modul backend Go yang sudah memakai `database/sql` dan driver `pgx/v5`. Tambahkan HTTP API di bawah prefix `/api/v1/subscription-profiles`, lapisan repository untuk PostgreSQL, service untuk validasi dan aturan bisnis, serta handler HTTP untuk serialisasi JSON dan status code. Frontend tetap berupa halaman statis pada tahap pertama, tetapi seluruh operasi baca/tulis berpindah ke API.

Jangan membiarkan frontend menyimpan credential database, RADIUS, MikroTik, OLT/ACS, atau payment gateway. Nilai yang bersifat konfigurasi jaringan hanya boleh disimpan dan dikembalikan sesuai izin pengguna; secret operasional harus berada pada secret service atau vault.

## 2. Kontrak data

Objek profile yang dikirim API perlu mencakup field yang sudah ada di form, namun menggunakan nama JSON yang konsisten dan tipe eksplisit.

| Field | Tipe API | Validasi utama |
|---|---|---|
| `id` | UUID atau string UUID | Diisi server saat create |
| `name` | string | Wajib, trim, unik dalam tenant |
| `service_type` | enum | `FTTH`, `PPPOE`, `HOTSPOT_VOUCHER`, `STATIC_IP` |
| `category` | enum/string | `RUMAHAN`, `BISNIS`, `DEDICATED`, `HOTSPOT` |
| `media` | enum/string | `FIBER_OPTIC`, `WIRELESS`, `LAN`, `LTE_5G` |
| `color` | string | Format `#RRGGBB` |
| `description` | string nullable | Batas panjang, misalnya 2.000 karakter |
| `status` | enum | `ACTIVE` atau `INACTIVE` |
| `mikrotik_group` | string nullable | Tidak boleh memuat credential |
| `radius_group` | string nullable | Tidak boleh memuat credential |
| `rate_limit` | string nullable | Validasi format sesuai policy MikroTik |
| `upload_speed` | string nullable | Normalisasi satuan atau simpan sebagai integer bps |
| `download_speed` | string nullable | Normalisasi satuan atau simpan sebagai integer bps |
| `shared_users` | integer | Minimal 1 |
| `vlan_id` | integer nullable | Rentang 1–4094 |
| `olt_profile` | string nullable | Identifier profile, bukan secret |
| `ip_pool` | string nullable | Identifier pool/address-list |
| `monthly_price` | integer | Minimal 0, dalam satuan rupiah |
| `active_days` | integer | Minimal 0 |
| `commission_amount` | integer | Minimal 0 |
| `commission_type` | enum | `RUPIAH` atau `PERCENT` |
| `billing_cycle` | enum | `DAILY`, `WEEKLY`, `MONTHLY`, `CUSTOM` |
| `auto_isolate` | boolean | Default `true` |
| `billing_note` | string nullable | Batas panjang |
| `version` | integer | Untuk optimistic locking |
| `created_at`, `updated_at` | timestamp | Diisi server |

## 3. Skema PostgreSQL

Buat migration versioned, misalnya `database/migrations/0002_subscription_profiles.sql`. Gunakan UUID, timestamp UTC, check constraint, unique index untuk profile aktif/nonaktif sesuai kebutuhan tenant, dan kolom `version` untuk mencegah lost update.

```sql
create table subscription_profiles (
  id uuid primary key default gen_random_uuid(),
  tenant_id uuid not null,
  name text not null,
  service_type text not null check (service_type in ('FTTH','PPPOE','HOTSPOT_VOUCHER','STATIC_IP')),
  category text,
  media text,
  color char(7) not null default '#1677ff' check (color ~ '^#[0-9A-Fa-f]{6}$'),
  description text,
  status text not null default 'ACTIVE' check (status in ('ACTIVE','INACTIVE')),
  mikrotik_group text,
  radius_group text,
  rate_limit text,
  upload_bps bigint,
  download_bps bigint,
  shared_users integer not null default 1 check (shared_users >= 1),
  vlan_id integer check (vlan_id between 1 and 4094),
  olt_profile text,
  ip_pool text,
  monthly_price bigint not null default 0 check (monthly_price >= 0),
  active_days integer not null default 30 check (active_days >= 0),
  commission_amount bigint not null default 0 check (commission_amount >= 0),
  commission_type text not null default 'RUPIAH' check (commission_type in ('RUPIAH','PERCENT')),
  billing_cycle text not null default 'MONTHLY' check (billing_cycle in ('DAILY','WEEKLY','MONTHLY','CUSTOM')),
  auto_isolate boolean not null default true,
  billing_note text,
  version integer not null default 1 check (version >= 1),
  created_by uuid,
  updated_by uuid,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create unique index subscription_profiles_tenant_name_uq
  on subscription_profiles (tenant_id, lower(name));
create index subscription_profiles_tenant_status_idx
  on subscription_profiles (tenant_id, status);
```

Tambahkan `audit_events` atau gunakan audit subsystem yang sudah ada untuk mencatat create, update, delete, status change, dan export. Untuk operasi yang mengubah parameter network atau billing, simpan perubahan sebagai approval request sebelum diterapkan ke sistem eksternal.

## 4. Endpoint API

| Method | Endpoint | Tujuan |
|---|---|---|
| `GET` | `/api/v1/subscription-profiles` | Daftar dengan `q`, `service_type`, `status`, `limit`, `cursor` |
| `GET` | `/api/v1/subscription-profiles/{id}` | Detail profile |
| `POST` | `/api/v1/subscription-profiles` | Membuat profile baru |
| `PATCH` | `/api/v1/subscription-profiles/{id}` | Mengubah profile dengan `If-Match` atau `version` |
| `DELETE` | `/api/v1/subscription-profiles/{id}` | Soft delete atau menonaktifkan profile |
| `POST` | `/api/v1/subscription-profiles/{id}/approval` | Mengajukan perubahan sensitif untuk approval |
| `GET` | `/api/v1/subscription-profiles/export` | Export server-side sesuai izin pengguna |

Gunakan respons error JSON seragam, misalnya `{ "error": { "code": "VALIDATION_ERROR", "message": "...", "fields": {...} } }`. `POST` mengembalikan `201`, `PATCH` mengembalikan `200`, delete sukses mengembalikan `204`, validasi `400`, autentikasi `401`, otorisasi `403`, konflik versi `409`, dan resource tidak ditemukan `404`.

## 5. Keamanan dan aturan bisnis

Semua endpoint harus memerlukan autentikasi dan otorisasi berbasis tenant/role. Terapkan pemeriksaan object-level authorization agar ID dari tenant lain tidak dapat diakses. Validasi dilakukan ulang di server; validasi HTML hanya untuk pengalaman pengguna.

Gunakan parameterized query melalui `database/sql`, timeout context, limit maksimum pagination, dan logging terstruktur tanpa menulis credential atau data sensitif. Terapkan optimistic locking: `UPDATE ... WHERE id = $1 AND version = $2`, lalu kembalikan `409 Conflict` jika tidak ada baris yang berubah.

Perubahan `rate_limit`, `mikrotik_group`, `radius_group`, `vlan_id`, `olt_profile`, `ip_pool`, `monthly_price`, `billing_cycle`, atau `auto_isolate` harus menghasilkan audit event. Pada fase pertama, API hanya mengubah profile database; sinkronisasi ke RADIUS/MikroTik/OLT/billing dilakukan oleh worker setelah approval eksplisit.

## 6. Migrasi frontend

Ganti inisialisasi `localStorage` dengan `loadProfiles()` yang memanggil `GET`. Ganti `store()` dengan `createProfile()`, `updateProfile()`, dan `deleteProfile()`. Tampilkan loading state, empty state, error state, retry, dan pesan konflik versi. Simpan `version` dari respons server dan kirimkan saat update.

Selama rollout, sediakan tombol **Impor data lokal** yang membaca `a_radius_profile_v1`, mengirim setiap profile ke endpoint create setelah validasi, lalu menandai hasil per item. Jangan melakukan migrasi otomatis tanpa konfirmasi karena data browser dapat berisi profile lama atau duplikat. Setelah migrasi selesai, localStorage dapat dipertahankan sebagai backup sementara, bukan sebagai sumber kebenaran.

## 7. Tahapan implementasi

| Tahap | Hasil | Kriteria selesai |
|---|---|---|
| 1. Fondasi | Migration, model, repository, health check | Migration repeatable dan repository memiliki test |
| 2. Read path | `GET` list/detail dan koneksi frontend | UI menampilkan data PostgreSQL |
| 3. Write path | Create/update/delete dengan validasi | Test API dan optimistic locking lulus |
| 4. Security | Authz tenant/role, audit event, rate limit | Negative authorization tests lulus |
| 5. Approval | Approval queue untuk perubahan sensitif | Tidak ada side effect eksternal tanpa approval |
| 6. Cutover | Import localStorage dan feature flag | Data lokal terverifikasi, API menjadi sumber utama |
| 7. Operasional | CI migration check, backup, monitoring | Rollback dan observability terdokumentasi |

## 8. Pengujian

Tambahkan unit test untuk validator, mapping DTO, service, dan parser rate limit. Tambahkan integration test PostgreSQL menggunakan database terisolasi untuk migration, CRUD, unique constraint, pagination, dan optimistic locking. Tambahkan HTTP test untuk `401`, `403`, cross-tenant access, malformed JSON, oversized input, dan duplicate name.

Untuk frontend, uji loading, error, retry, filter, create/update/delete, import localStorage, export server-side, dan konflik versi. Pipeline harus menjalankan `go test ./...`, `go vet ./...`, migration check, serta `govulncheck ./...`. Lakukan backup database dan canary rollout sebelum cutover produksi.

## 9. Keputusan yang masih perlu dikonfirmasi

Sebelum implementasi, tentukan apakah A-Radius akan memiliki multi-tenant secara resmi, mekanisme autentikasi yang dipakai, apakah profile dihapus secara soft delete, format canonical kecepatan, dan apakah approval dilakukan oleh role admin/operator atau melalui workflow terpisah. Tanpa keputusan ini, implementasi sebaiknya berhenti pada migration, repository, dan endpoint read-only.
