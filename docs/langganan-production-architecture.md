# A-RADIUS Langganan Production

## Tujuan

Modul Langganan Production menjadi orchestration boundary untuk data pelanggan, service profile, billing, RADIUS, MikroTik, payment, dan finance. Modul ini tidak boleh berdiri sendiri karena aktivasi layanan dapat memengaruhi akses jaringan, invoice, ledger, dan layanan pelanggan.

## Lifecycle perubahan

```text
DEVELOPER
    ↓
PREVIEW
    ↓
APPROVAL
    ↓
PRODUCTION
```

Semua perubahan yang memengaruhi status layanan, isolir, router, RADIUS, billing, invoice, payment, atau ledger wajib menghasilkan preview, menjalankan integration readiness check, dan memperoleh approval dari role `administrator` atau `developer` sebelum dipromosikan.

## Domain integration

| Domain | Tanggung jawab | Kegagalan wajib |
|---|---|---|
| Customer | Validasi identitas dan customer ID | Block activation |
| Service | Profile, kategori, media, dan status layanan | Block activation |
| Billing | Siklus, prorata, tanggal invoice, isolir, PPN | Block activation |
| RADIUS | Group, secret, IP, dan authorization mapping | Quarantine change |
| MikroTik | Router, server, address list, dan rate limit | Quarantine change |
| Payment | Metode pembayaran dan recurring status | Mark payment pending |
| Finance | Invoice posting dan ledger mapping | Block invoice posting |

Setiap binding memiliki source of truth dan failure mode sehingga kegagalan satu domain tidak boleh menghasilkan perubahan parsial yang tidak terkontrol.

## Batas Production

JWT authentication, RBAC, optimistic locking, approval record, integration evidence, dan audit trail harus diperiksa sebelum operasi Production. AI tidak memiliki `production access`; AI hanya boleh membuat analysis, finding, recommendation, patch preview, dan test job.

Endpoint preview mengembalikan `202 Accepted` dan `production_changed=false`. Endpoint readiness hanya membaca status dependency. Implementasi saat ini belum mengeksekusi konektor RADIUS, MikroTik, Payment, atau Finance; konektor tersebut harus ditambahkan sebagai adapter dengan timeout, idempotency key, retry policy, dan compensating action.

## Failure handling

Jika Customer, Service, atau Billing belum siap, aktivasi ditolak. Jika RADIUS atau MikroTik gagal, perubahan dikarantina dan tidak dipromosikan. Jika Payment belum tersedia, status payment menjadi pending tanpa memalsukan pembayaran. Jika Finance belum siap, invoice tidak diposting ke ledger.

## Route Production

| Method | Endpoint | Permission | Dampak |
|---|---|---|---|
| GET | `/api/v1/subscription-production/policy` | `subscription:read` | Read-only |
| GET | `/api/v1/subscription-production/integrations` | `subscription:read` | Readiness metadata |
| GET | `/api/v1/subscription-production/readiness` | `subscription:read` | Readiness check |
| POST | `/api/v1/subscription-production/preview` | `subscription:preview` | Preview only |

Route tersebut dipasang di belakang JWT middleware dan authorization audit yang sama dengan Subscription Profile API.
