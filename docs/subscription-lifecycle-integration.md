# A-RADIUS Subscription Lifecycle Integration

## Posisi Langganan

Langganan bukan database yang berdiri sendiri. Ia menjadi pusat lifecycle yang menghubungkan identitas pelanggan, service profile, billing, network provisioning, payment, finance, dan Customer App.

```text
PELANGGAN
   ↓
LANGGANAN
   ├── SERVICE → PROFILE
   ├── BILLING → INVOICE → PAYMENT → FINANCE
   └── NETWORK → RADIUS → MikroTik → GenieACS → OLT
                         ↓
                    PPPoE / HOTSPOT
                         ↓
                    CUSTOMER APP
```

## Otomasi status

| Kondisi | Status/aksi | Efek |
|---|---|---|
| Aktif dan invoice belum jatuh tempo | `ACTIVE` | Service ON |
| Invoice jatuh tempo | `WARNING` | Notifikasi dan billing follow-up |
| Melewati batas isolir | `ISOLATED` | RADIUS/MikroTik block melalui worker |
| Pembayaran diterima | `REACTIVATING` | Reconcile payment lalu re-enable service |
| Service activated | `ACTIVE` | Refresh network mapping dan Customer App |

Otomasi event tidak melewati approval untuk transisi operasional yang memang ditentukan oleh policy billing. Namun operasi massal atau tindakan yang mengubah banyak pelanggan tetap harus menggunakan proposal dan approval.

## Bulk action approval

Operasi `DELETE`, `SET_INACTIVE`, `SET_ISOLATED`, `CHANGE_ROUTER`, `CHANGE_PROFILE`, dan `CHANGE_BILLING` diperlakukan sebagai action berisiko. Alurnya adalah:

```text
AI / ADMIN REQUEST
   ↓
ANALYSIS
   ↓
PREVIEW: jumlah target, from/to, affected modules, risk
   ↓
ACTION PROPOSAL
   ↓
ADMIN APPROVAL
   ↓
AUTHORIZATION + optimistic locking
   ↓
REDIS WORKER
   ↓
EXECUTION
   ↓
VALIDATION
   ↓
AUDIT TRAIL
```

Contoh proposal perubahan profile dari `M20` ke `S150KM` untuk 69 pelanggan tidak boleh langsung menulis database atau network controller. Proposal harus menyimpan `target_filter`, jumlah target, actor, approval, dan status execution.

## Redis worker

`internal/subscriptionlifecycle/worker.go` menggunakan queue `a-radius:subscription:lifecycle`. Worker menolak proposal yang membutuhkan approval tetapi belum memiliki `approved_by`. Setiap job harus diproses dengan idempotency key pada event dan audit outcome `SUCCESS`, `FAILED`, atau `REJECTED_NO_APPROVAL`.

Adapter eksekusi untuk RADIUS, MikroTik, GenieACS, OLT, Payment, dan Finance harus diimplementasikan terpisah. Worker hanya mengorkestrasi proposal yang telah sah; worker tidak boleh menjadi jalur bypass RBAC.

## Audit example

```text
AUDIT #LGN-20260817-000821
User       : Administrator
Role       : Administrator Cabang
Action     : GANTI PROFILE
Target     : 69 Langganan
From       : M20
To         : S150KM
AI         : Recommendation
Approval   : APPROVED
Approved by: ADM-001
Execution  : SUCCESS
Worker     : redis-worker-03
Timestamp  : 17/08/2026 13:42 WITA
```

## Database

Migration `0004_subscription_lifecycle.up.sql` menambahkan tabel event lifecycle, bulk action proposal, dan audit bulk action. Event memiliki unique `idempotency_key`; proposal memiliki status execution dan approval; audit menyimpan worker, outcome, error, dan metadata.
