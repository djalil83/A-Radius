# Cross-Application Security Layer

## Tujuan

Developer Security Center bukan hanya halaman monitoring untuk role developer. Ia menjadi **security layer lintas-aplikasi** A-Radius yang menerima signal dan menerapkan kebijakan untuk API, dashboard Administrator, Teknisi, Reseller, Pelanggan, authentication, database, konfigurasi, serta deployment.

Security Center memusatkan visibility dan policy, tetapi enforcement tetap berada sedekat mungkin dengan resource yang dilindungi. API tetap memvalidasi JWT, RBAC, tenant scope, input, dan rate limit pada request. Security Center tidak boleh menjadi bypass terhadap guard domain.

## Tiga lapisan

| Lapisan | Tanggung jawab | Contoh kontrol |
|---|---|---|
| Prevention | Mencegah request atau perubahan tidak aman sebelum terjadi | RBAC, MFA, session security, rate limiting, input validation, secret management, API authorization |
| Detection | Mengumpulkan signal dan menemukan anomali | AI Security Checker, anomaly detection, login/API monitoring, file/database/permission change monitoring |
| Response | Menyiapkan tindakan mitigasi, approval, pemulihan, dan bukti | Alert, quarantine/lock proposal, rollback, credential disable approval, incident report, audit trail |

## AI advisory boundary

AI boleh melakukan scanning, membuat finding, menganalisis dampak, menyusun rekomendasi, dan menghasilkan proposed fix. AI **tidak** boleh memiliki credential deployment, akses langsung ke Production, atau kemampuan menjalankan tindakan disruptif.

Tindakan yang dapat berdampak pada layanan pelanggan harus menjadi proposal yang menunggu persetujuan manusia. Contohnya adalah memblokir akun, memutus API, mengubah firewall, menghapus credential, dan deployment.

```text
Signal / finding
    -> AI recommendation
    -> Proposed action
    -> Developer Preview
    -> Security Test
    -> Staging Test
    -> Administrator/Developer approval
    -> Execution by controlled service
    -> Health check
    -> Success: audit
    -> Failure: rollback + audit
```

## Ownership

Developer bertanggung jawab atas platform security, code/configuration, scanner, dependency, deployment, dan technical remediation. Administrator bertanggung jawab atas akses operasional, approval bisnis, user impact, dan keputusan yang menyentuh pelanggan. Role Technician, Sales, Reseller, dan Pelanggan dapat menjadi sumber telemetry atau affected subject, tetapi tidak memperoleh hak menjalankan tindakan disruptif dari Security Center.

| Tindakan | AI | Developer | Administrator |
|---|---|---|---|
| Membuat finding | Boleh | Boleh | Boleh |
| Membuat recommendation | Boleh | Boleh | Boleh |
| Generate proposed fix | Boleh | Review | Review |
| Block account | Proposal saja | Approve sesuai scope | Approve sesuai dampak |
| Disconnect API | Proposal saja | Approve | Approve |
| Change firewall | Proposal saja | Approve | Approve |
| Delete credential | Proposal saja | Approve | Approve |
| Deploy | Tidak boleh langsung | Approve dan execute melalui pipeline | Approve sesuai kebijakan |
| Rollback | Proposal atau trigger otomatis terbatas | Approve/execute | Approve sesuai dampak |

## Enforcement requirements

Setiap response action wajib membawa actor terverifikasi dari JWT, tenant/resource scope, change ticket atau approval request, idempotency key, reason, expiry bila relevan, dan audit event. Approval harus memeriksa permission, status finding, hasil security test, hasil staging test, dan separation of duties. Tidak boleh menganggap role dari browser atau klaim AI sebagai bukti approval.

Guard harus **fail-closed**: jika approval, scope, audit, atau health-check evidence tidak tersedia, action tidak boleh dieksekusi. Setelah deployment, health check menentukan success atau rollback; kedua hasil tersebut harus dicatat sebagai audit trail immutable.
