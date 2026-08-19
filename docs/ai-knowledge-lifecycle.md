# A-RADIUS Security Knowledge Lifecycle

Setiap versi Security Knowledge mempunyai lifecycle terkontrol:

```text
DISCOVERED
    ↓
ANALYZING
    ↓
VALIDATED
    ↓
REVIEW REQUIRED
    ↓
APPROVED
    ↓
STAGED
    ↓
ACTIVE
    ↓
SUPERSEDED
    ↓
ARCHIVED
```

AI dapat membuat knowledge berstatus `DISCOVERED`, menjalankan proses `ANALYZING`, dan menghasilkan evidence sampai `VALIDATED`. AI tidak dapat langsung menetapkan status `ACTIVE`. Transisi `APPROVED → STAGED → ACTIVE` memerlukan kontrol Developer yang memiliki permission dan approval yang sesuai.

## State transition policy

| Dari | Ke | Actor yang diizinkan |
|---|---|---|
| `DISCOVERED` | `ANALYZING` | AI atau worker knowledge |
| `ANALYZING` | `VALIDATED` / `REVIEW REQUIRED` | AI/worker setelah validation |
| `VALIDATED` | `REVIEW REQUIRED` | AI/worker |
| `REVIEW REQUIRED` | `APPROVED` | Developer melalui approval |
| `REVIEW REQUIRED` | `ARCHIVED` | Developer atau policy reviewer |
| `APPROVED` | `STAGED` | Developer/release controller |
| `STAGED` | `ACTIVE` | Developer/release controller setelah staging test |
| `ACTIVE` | `SUPERSEDED` | Release controller saat versi baru aktif |
| `SUPERSEDED` | `ARCHIVED` | Retention/archive worker atau Developer |

Transisi AI ke `ACTIVE` selalu ditolak oleh guard backend. Status knowledge juga tidak sama dengan status deployment aplikasi; versi `ACTIVE` hanya berarti rule/knowledge version aktif untuk analisis, bukan bahwa kode Production berubah.

## Knowledge registry

Registry menampilkan active version `SK-2.4.7`, last update `17/08/2026 14:20 WITA`, dan counters `NEW=12`, `REVIEW=4`, `ACTIVE=1`, `ARCHIVED=18`. Versi historis `SK-2.4.6`, `SK-2.4.5`, dan `SK-2.4.4` dipertahankan untuk comparison, provenance, dan audit.

## New Security Intelligence

`INT-2026-00821` adalah candidate `SK-2.4.8` pada kategori API Security dengan risk `HIGH` dan confidence `91%`. Intelligence ini relevan dengan `/api/langganan`, `/api/pelanggan`, `/api/admin`, dan `/api/technician`, tetapi tidak relevan dengan UI statis maupun template voucher. Proposed rule-nya adalah `API-AUTH-042`, dengan status `REVIEW REQUIRED`.

Aksi `VIEW EVIDENCE`, `RUN TEST`, `COMPARE`, dan `REQUEST APPROVAL` hanya menghasilkan evidence, test job, comparison, atau approval request. Aksi tersebut tidak mengubah Production.
