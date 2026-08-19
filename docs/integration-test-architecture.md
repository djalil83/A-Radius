# Arsitektur Integration Test Subscription Profile

## 1. Tujuan

Dokumen ini mendefinisikan alur, struktur direktori, komponen, dan skenario pengujian otomatis untuk endpoint CRUD **Subscription Profile** pada A-Radius. Pengujian dirancang untuk memvalidasi perilaku API Go/Chi terhadap PostgreSQL nyata, bukan sekadar mock repository.

Ruang lingkup mencakup pembuatan profile, pembacaan daftar dan detail, pembaruan dengan optimistic locking, konflik versi, revision history, soft delete/archive, validasi request, autentikasi header sementara, serta audit trail dari trigger database.

## 2. Arsitektur komponen

| Komponen | Tanggung jawab | Teknologi |
|---|---|---|
| Test runner | Menjalankan test dengan build tag `integration` | Go `testing` |
| Test harness | Membuka koneksi, menerapkan migration, membuat fixture, dan cleanup | `database/sql`, `pgx/v5` |
| HTTP test client | Mengirim request ke router tanpa port eksternal | `httptest` |
| Router/API | Menjalankan handler CRUD dan versioning | Chi |
| Repository | Menjalankan query PostgreSQL dan optimistic locking | `database/sql` |
| PostgreSQL | Menyimpan profile, revisions, approval, dan audit | PostgreSQL 17 |
| Migration | Membuat enum, tabel, index, function, dan trigger | SQL |
| CI/Compose | Menyediakan database disposable untuk pengujian | Docker Compose |

## 3. Flowchart utama

```mermaid
flowchart TD
    A([Mulai: go test -tags=integration]) --> B{INTEGRATION_DATABASE_URL tersedia?}
    B -- Tidak --> C[Skip dengan pesan prasyarat]
    B -- Ya --> D[Parse DSN dan buka database]
    D --> E{PostgreSQL tersedia?}
    E -- Tidak --> F[Fail: database tidak dapat diakses]
    E -- Ya --> G[Ping dengan timeout]
    G --> H[Reset schema test / gunakan database disposable]
    H --> I[Terapkan migration 0002_subscription_profiles.up.sql]
    I --> J{Migration berhasil?}
    J -- Tidak --> K[Fail: hentikan test dan tampilkan error]
    J -- Ya --> L[Buat repository, handler, dan httptest router]
    L --> M[Pasang X-Tenant-ID dan X-Actor-ID]
    M --> N[POST: create profile]
    N --> O{HTTP 201 dan profile valid?}
    O -- Tidak --> P[Fail assertion]
    O -- Ya --> Q[GET: list dan detail]
    Q --> R{Data dan version=1 benar?}
    R -- Tidak --> P
    R -- Ya --> S[PATCH dengan version=1]
    S --> T{HTTP 200 dan version=2?}
    T -- Tidak --> P
    T -- Ya --> U[GET revisions]
    U --> V{Revision CREATE + UPDATE tersedia?}
    V -- Tidak --> P
    V -- Ya --> W[PATCH ulang memakai version lama]
    W --> X{HTTP 409 VERSION_CONFLICT?}
    X -- Tidak --> P
    X -- Ya --> Y[DELETE/archive memakai version=2]
    Y --> Z{HTTP 204 dan status ARCHIVED?}
    Z -- Tidak --> P
    Z -- Ya --> AA[Verifikasi audit trail dan revision]
    AA --> AB{Audit dan histori konsisten?}
    AB -- Tidak --> P
    AB -- Ya --> AC[Cleanup data/schema test]
    AC --> AD([Lulus])
    P --> AE[Cleanup melalui defer]
    AE --> AF([Gagal])
```

## 4. Sequence flow request

```mermaid
sequenceDiagram
    participant T as Test Runner
    participant H as Chi Handler
    participant R as Repository
    participant DB as PostgreSQL
    participant TR as DB Triggers

    T->>H: POST /api/v1/subscription-profiles
    H->>H: Validasi JSON dan identity header
    H->>R: Create(tenant, actor, payload)
    R->>DB: INSERT profile
    DB->>TR: Snapshot revision + audit event
    DB-->>R: Profile version=1
    R-->>H: Profile
    H-->>T: 201 Created

    T->>H: PATCH /{id} version=1
    H->>R: Update(..., version=1)
    R->>DB: UPDATE ... WHERE version=1
    DB->>TR: Snapshot revision + audit event
    DB-->>R: Profile version=2
    R-->>H: Profile
    H-->>T: 200 OK

    T->>H: PATCH /{id} version=1
    H->>R: Update(..., version=1)
    R->>DB: UPDATE affects 0 rows
    DB-->>R: Current version=2
    R-->>H: ErrConflict
    H-->>T: 409 VERSION_CONFLICT

    T->>H: DELETE /{id}?version=2
    H->>R: Archive(..., version=2)
    R->>DB: UPDATE status=ARCHIVED, version=3
    DB->>TR: Snapshot revision + audit event
    DB-->>R: Success
    R-->>H: nil
    H-->>T: 204 No Content
```

## 5. Struktur direktori yang disarankan

```text
A-Radius/
├── cmd/
│   └── profile-api/
│       └── main.go
├── internal/
│   └── subscriptionprofile/
│       ├── model.go
│       ├── repository.go
│       ├── service.go
│       ├── http.go
│       ├── service_test.go
│       ├── http_test.go
│       ├── integration_test.go       # build tag: integration
│       ├── integration_fixture.go    # DB, migration, identity, cleanup
│       └── README.md
├── database/
│   └── migrations/
│       ├── 0002_subscription_profiles.up.sql
│       └── 0002_subscription_profiles.down.sql
├── scripts/
│   └── test-integration.sh           # Compose + go test + cleanup
├── docker-compose.yml
├── Dockerfile
├── .env.example
├── go.mod
├── go.sum
└── docs/
    └── integration-test-architecture.md
```

## 6. Matriks skenario pengujian

| ID | Skenario | Request | Ekspektasi |
|---|---|---|---|
| IT-001 | Identity wajib | GET tanpa header | `401 UNAUTHENTICATED` |
| IT-002 | Create valid | POST payload lengkap | `201`, UUID valid, `version=1`, status `ACTIVE` |
| IT-003 | List tenant isolation | GET dengan tenant berbeda | `200`, profile tenant lain tidak muncul |
| IT-004 | Get detail | GET `/{id}` | `200`, field utama sama dengan create |
| IT-005 | Update valid | PATCH dengan `version=1` | `200`, `version=2`, nilai berubah |
| IT-006 | Update stale version | PATCH dengan `version=1` setelah update | `409 VERSION_CONFLICT` |
| IT-007 | Revision history | GET `/{id}/revisions` | Revision CREATE dan UPDATE tersedia |
| IT-008 | Archive valid | DELETE dengan `version=2` | `204`, tidak muncul pada list aktif |
| IT-009 | Archive stale version | DELETE dengan version lama | `409 VERSION_CONFLICT` |
| IT-010 | Not found | GET/PATCH/DELETE UUID tidak ada | `404 NOT_FOUND` |
| IT-011 | Invalid payload | POST enum/color/nilai tidak valid | `400 VALIDATION_ERROR` |
| IT-012 | Audit integrity | Query `audit_events` | CREATE, UPDATE, DELETE/STATUS_CHANGE tercatat |

## 7. Kontrak fixture

Fixture harus memakai `tenant_id` dan `actor_id` UUID deterministik yang berbeda dari data development lain. Setiap test sebaiknya membuat nama profile unik, misalnya `integration-<timestamp>`, sehingga test dapat dijalankan paralel tanpa collision. Database yang dipakai harus disposable; jangan menjalankan reset schema pada database bersama.

```go
const testTenantID = "00000000-0000-0000-0000-000000000101"
const testActorID = "00000000-0000-0000-0000-000000000102"
```

## 8. Perintah eksekusi

Jika PostgreSQL lokal tersedia:

```bash
INTEGRATION_DATABASE_URL='postgres://a_radius:a_radius_dev_only@localhost:5432/a_radius?sslmode=disable' \
  go test -tags=integration -v ./internal/subscriptionprofile
```

Dengan Docker Compose, wrapper script sebaiknya menjalankan PostgreSQL, menunggu healthcheck, menjalankan migration, menjalankan test, lalu menghapus container dengan `trap`:

```bash
./scripts/test-integration.sh
```

## 9. Kriteria kelulusan

Integration test dinyatakan lulus apabila seluruh assertion HTTP, payload, version increment, tenant isolation, revision snapshot, audit event, dan archive behavior terpenuhi. Test harus gagal apabila database tidak tersedia, migration gagal, optimistic locking tidak bekerja, atau audit trail tidak konsisten.

Pengujian di sandbox yang tidak memiliki Docker hanya dapat menjalankan static check dan unit test. Pengujian integration harus dijalankan pada CI atau host yang memiliki Docker dan PostgreSQL.

## 10. Catatan keamanan

Header `X-Tenant-ID` dan `X-Actor-ID` masih merupakan adapter identity sementara. Pada production, integration test harus menambahkan skenario autentikasi middleware aktual dan memastikan tenant/actor tidak dapat dipalsukan oleh client. Password database harus berasal dari secret manager, bukan dari file yang di-commit.

## 11. Skenario race condition dan konkurensi tinggi

Pengujian konkurensi harus membuktikan bahwa optimistic locking pada kolom `version` benar-benar mencegah **lost update**. Semua worker membaca profile yang sama pada `version=1`, kemudian menunggu barrier agar request PATCH dikirim hampir bersamaan. Hanya satu request yang boleh berhasil mengubah baris; request lain harus menerima `409 VERSION_CONFLICT` dan tidak boleh mengubah data atau menambah revision.

| ID | Skenario konkurensi | Metode | Ekspektasi |
|---|---|---|---|
| RC-001 | Concurrent update, same version | 32 goroutine PATCH dengan `version=1` | Tepat 1 respons `200`, 31 respons `409 VERSION_CONFLICT`, final version `2` |
| RC-002 | Concurrent update, high fan-out | 100–250 worker PATCH memakai version snapshot yang sama | Tepat 1 pemenang; tidak ada `500`, deadlock, atau timeout; version hanya naik satu kali |
| RC-003 | Concurrent update dengan payload berbeda | Worker mengirim nilai `billing_note` berbeda | Nilai final hanya berasal dari pemenang; tidak ada field partial/lost update |
| RC-004 | Update versus archive | PATCH dan DELETE bersamaan pada version sama | Tepat satu operasi menang; operasi kalah mendapat `409`; status/version konsisten |
| RC-005 | Sequential retry after conflict | Worker yang kalah GET ulang lalu PATCH dengan version terbaru | Retry berhasil pada version berikutnya; revision dan audit bertambah tepat satu kali |
| RC-006 | Cross-tenant concurrency | Dua tenant mengubah UUID profile yang sama | Tenant yang bukan pemilik mendapat `404`, bukan `409` atau data bocor |
| RC-007 | Concurrent reads during writes | Banyak GET bersamaan dengan PATCH | Semua response valid; tidak ada JSON rusak; pembaca hanya melihat snapshot committed |
| RC-008 | Audit/revision consistency | Query revision dan audit setelah fan-out update | Jumlah revision dan event sesuai jumlah operasi sukses, bukan jumlah request gagal |

### Invariant yang harus diverifikasi

Test harus memeriksa invariant berikut setelah seluruh worker selesai:

| Invariant | Kondisi yang diharapkan |
|---|---|
| Single winner | Untuk satu expected version, hanya satu mutasi yang berhasil |
| Monotonic version | Version tidak pernah menurun dan bertambah tepat satu per mutasi sukses |
| No lost update | Nilai final berasal dari payload pemenang yang diketahui test |
| Revision alignment | Setiap mutasi sukses memiliki satu revision yang sesuai dengan version finalnya |
| Audit alignment | Request gagal karena konflik tidak membuat audit `UPDATE` atau `STATUS_CHANGE` palsu |
| Atomic archive | Archive mengubah status, deleted timestamp, dan version dalam satu transaksi update |
| Tenant isolation | Konflik dan data antar-tenant tidak dapat digunakan untuk menebak resource tenant lain |

### Pola implementasi Go

Gunakan `sync.WaitGroup` untuk menahan worker pada barrier dan channel buffered untuk mengumpulkan hasil tanpa blocking. Jangan menggunakan sleep sebagai sinkronisasi utama karena dapat menghasilkan test flaky. Setiap worker wajib memiliki context timeout dan mengembalikan status HTTP, error code, version response, serta payload response.

```go
start := make(chan struct{})
results := make(chan result, workers)
var wg sync.WaitGroup

for i := 0; i < workers; i++ {
    wg.Add(1)
    go func(workerID int) {
        defer wg.Done()
        <-start
        results <- patchWithVersion(t, router, profileID, 1, workerID)
    }(i)
}
close(start)
wg.Wait()
close(results)

var success, conflict int
for got := range results {
    switch got.StatusCode {
    case http.StatusOK:
        success++
    case http.StatusConflict:
        conflict++
    default:
        t.Errorf("unexpected status: %d", got.StatusCode)
    }
}
if success != 1 || conflict != workers-1 {
    t.Fatalf("success=%d conflict=%d workers=%d", success, conflict, workers)
}
```

### Validasi database setelah race test

Setelah request selesai, test harus mengambil profile final melalui API, kemudian menjalankan query verifikasi terkontrol pada database test. Untuk satu create dan satu update sukses, nilai yang diharapkan adalah `version=2`, dua revision profile (`CREATE` dan `UPDATE`), serta satu event audit untuk create dan satu event audit untuk update. Jika pengujian archive menjadi pemenang, final version menjadi `3`, status `ARCHIVED`, dan `deleted_at` tidak null.

Race test harus dijalankan dengan `-race` untuk mendeteksi data race pada kode Go, tetapi flag tersebut tidak menggantikan pengujian concurrency pada PostgreSQL:

```bash
go test -race -tags=integration -run 'TestIntegration_Concurrent' -count=10 -v ./internal/subscriptionprofile
```

Gunakan `-count=10` pada CI atau nightly pipeline untuk meningkatkan peluang menangkap interleaving yang jarang terjadi. Test tidak boleh mengandalkan urutan worker; identitas worker hanya dipakai untuk menghasilkan payload unik dan menemukan payload pemenang.
