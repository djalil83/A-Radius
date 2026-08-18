# A-RADIUS Continuous Security Intelligence

## Tujuan

Developer Security Center berfungsi sebagai **security layer lintas-aplikasi** untuk code, API, database, dashboard, authentication, dan layanan pelanggan. Sistem mengumpulkan knowledge keamanan, memvalidasi sumbernya, membandingkan knowledge dengan kondisi A-RADIUS, lalu membuat finding dan proposed fix. AI tidak memperoleh akses langsung ke Production.

## Arsitektur pipeline

```text
Trusted security sources
  ├─ CISA KEV / CVE
  ├─ Vendor advisories / CSAF
  ├─ OWASP API Security
  ├─ Dependency advisories
  └─ Framework / infrastructure updates
          │
          ▼
  Knowledge ingestion
          │
          ▼
  Deduplication + schema validation + provenance + content hash
          │
          ▼
  AI Security Learning Engine
          │
          ▼
  A-RADIUS analysis targets: Code | API | DB
          │
          ▼
  Finding + evidence + risk score + recommendation
          │
          ▼
  Patch preview in sandbox/developer preview
          │
          ▼
  Automated security/regression test
          │
          ▼
  Human developer approval
          ├─ Reject → archive proposal
          └─ Approve → staging test → production review → controlled deploy
                                                    │
                                                    ▼
                                              health check → audit / rollback
```

## Sumber dan provenance

Sumber tidak dianggap sebagai instruksi eksekusi. Setiap knowledge item harus menyimpan source URL, publisher, retrieved timestamp, published timestamp bila tersedia, advisory/CVE ID, content hash, parser version, validation status, dan confidence. Item yang gagal validasi masuk quarantine dan tidak dipakai untuk menghasilkan patch.

CISA KEV dipakai sebagai sinyal prioritas eksploitasi; OWASP API Security menjadi taxonomy API risk; CSAF menjadi format advisory terstruktur; dan dependency advisory/NVD evidence digunakan untuk korelasi package/version/CVE. Referensi resmi disimpan dalam `docs/security-knowledge-sources.md`.

## Kontrak backend

`internal/dashboard/developer/continuous_security.go` menyediakan kontrak untuk `KnowledgeItem`, `SecurityAnalysis`, `ContinuousFinding`, `PatchPreview`, dan `ContinuousSecurityPolicy`. Policy default bersifat fail-closed:

| Aturan | Nilai |
|---|---|
| AI production access | `false` |
| Human approval | `true` |
| Required stages | developer approval, automated security test, staging test, production review, health check |
| Allowed AI environments | sandbox, staging |

Endpoint dashboard yang tersedia:

| Endpoint | Fungsi | Efek Production |
|---|---|---|
| `GET /developer/security/continuous/policy` | Membaca policy gate | Tidak ada |
| `GET /developer/security/continuous/sources` | Membaca trusted source registry | Tidak ada |
| `GET /developer/security/continuous/featured-finding` | Membaca contoh finding | Tidak ada |
| `GET/POST /developer/security/continuous/patch-preview` | Membaca atau membuat proposal preview | `UNCHANGED` |

Semua endpoint tetap harus dipasang di belakang JWT authentication dan RBAC. Endpoint proposal tidak boleh diperlakukan sebagai deployment endpoint.

## Batas AI dan approval

AI dapat memantau sumber yang dipercaya, memperbarui knowledge base, melakukan korelasi, menemukan potensi vulnerability, menganalisis perubahan, membuat recommendation, membuat report, membuat patch preview, melakukan test sandbox/staging, dan memberikan risk score. AI tidak dapat memblokir akun, memutus API, mengubah firewall, menghapus credential, mengganti permission produksi, atau melakukan deployment.

Tindakan disruptif memerlukan actor JWT terverifikasi, resource/tenant scope, approval record, change ticket, evidence test, idempotency key, dan audit event. Jika salah satu prasyarat tidak tersedia, backend harus menolak eksekusi.

## Operasi otomatis

Untuk pemeriksaan terjadwal beberapa kali per hari, worker terjadwal dapat mengambil feed, memvalidasi manifest, dan membuat job analisis. Untuk polling per menit atau selalu aktif, gunakan proses background yang tetap hidup di hosting terkelola; jangan membuat sesi AI baru untuk setiap polling deterministik. Jika analisis memerlukan penilaian AI yang kompleks, worker dapat menempatkan job ke antrean analisis dan tetap mempertahankan approval gate.

Dua pilihan operasional yang layak:

| Pendekatan | Tradeoff | Cost | Setup Complexity |
|---|---|---|---|
| Worker terjadwal dalam aplikasi dengan dashboard | Murah dan terintegrasi; cocok untuk sinkronisasi berkala dan sumber publik; bukan untuk polling sub-menit | Rendah, berbasis pemakaian | Sedang |
| Worker selalu aktif dengan antrean dan dashboard | Respons lebih cepat, cocok untuk monitoring kontinyu; membutuhkan resource dan observability lebih kuat | Lebih tinggi dibanding job berkala | Tinggi |

Pilihan ringan yang direkomendasikan untuk tahap awal adalah sinkronisasi terjadwal, dengan interval yang disesuaikan berdasarkan freshness advisory. Naikkan ke worker selalu aktif hanya ketika kebutuhan latency dan telemetry sudah terbukti.

## Validasi

```text
node --check web/dashboards/developer/dashboard.js  PASS
go test ./...                                      PASS
gofmt                                              PASS
git diff --check                                  PASS
```

## References

[1]: https://www.cisa.gov/known-exploited-vulnerabilities-catalog "CISA Known Exploited Vulnerabilities Catalog"
[2]: https://owasp.org/www-project-api-security/ "OWASP API Security Project"
[3]: https://docs.oasis-open.org/csaf/csaf/v2.0/os/csaf-v2.0-os.html "OASIS Common Security Advisory Framework 2.0"
[4]: https://owasp.org/www-project-dependency-check/ "OWASP Dependency-Check"
