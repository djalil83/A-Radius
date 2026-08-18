# Changelog

Semua perubahan penting pada proyek **A-RADIUS** didokumentasikan dalam berkas ini.

Format berkas ini mengikuti [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), dan proyek menggunakan [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

Belum ada perubahan yang dirilis setelah v1.0.0.

## [1.0.0] - 2026-08-18

Release awal A-RADIUS sebagai platform manajemen ISP dengan backend profil langganan, keamanan berbasis RBAC, dashboard multi-peran, lifecycle subscription, GenieACS, dan alur perubahan yang mewajibkan approval manusia.

### Added

#### Subscription Profile dan backend

- API Go berbasis Chi dan PostgreSQL untuk membaca, membuat, memperbarui, dan menghapus profil berlangganan.
- Model data untuk profile, service assignment, billing configuration, network configuration, lifecycle metadata, version, dan audit information.
- Versioning dengan optimistic locking untuk mencegah lost update.
- Respons `409 Conflict` ketika request menggunakan versi data yang sudah kedaluwarsa.
- Test service, HTTP handler, dan concurrency integration test.

#### Security dan authorization

- JWT middleware HS256 dengan validasi signature, issuer, audience, expiry, subject, dan algorithm allow-list.
- RBAC server-side berbasis database dengan relasi user, role, role-permission, dan permission.
- Permission middleware serta kebijakan authorization fail-closed.
- Audit trail untuk login, authorization decision, proposal, approval, deployment, rollback, dan command execution.

#### Continuous Security Intelligence

- Developer AI Center, Security Center, Security Knowledge, AI Research, Preview, Approval, Production, dan Audit Trail.
- Security finding, recommendation, patch preview, security test, staging, deployment, health check, dan rollback flow.
- Security Knowledge versioning dengan lifecycle `DISCOVERED`, `ANALYZING`, `VALIDATED`, `REVIEW REQUIRED`, `APPROVED`, `STAGED`, `ACTIVE`, `SUPERSEDED`, dan `ARCHIVED`.
- Pemisahan Security Knowledge, analysis cache, rule engine, patch preview, staging, dan Production.

#### Subscription Lifecycle

- Status subscription seperti `ACTIVE`, `WARNING`, dan `ISOLATED`.
- Bulk action proposal untuk ganti profile, isolir, reaktivasi, perubahan billing, penghapusan, dan perubahan network.
- Preview, risk assessment, approval request, worker execution, validation, dan audit trail untuk operasi sensitif.
- Orkestrasi Pelanggan, Service, Billing, Payment, RADIUS, MikroTik, GenieACS, dan Finance.

#### Administrator module hub

Dashboard Administrator menyediakan katalog modul berikut:

| Modul | Fokus |
|---|---|
| Voucher | Generate, distribusi, aktivasi, dan pencabutan voucher. |
| Mitra | Reseller, mitra cabang, dan komisi. |
| Billing | Invoice, jatuh tempo, denda, isolir, dan reaktivasi. |
| Payment | Rekonsiliasi pembayaran dan integrasi gateway. |
| NMS | Monitoring perangkat, link, dan alarm jaringan. |
| Teknisi | Work order, penugasan, dan validasi pekerjaan. |
| Finance | Pendapatan, biaya, komisi, dan laporan cabang. |
| Sistem | Konfigurasi cabang, integrasi, lisensi, dan health check. |

Fitur Administrator juga mencakup module catalog, AI Report, proposal preview, approval, rejection, before state, proposed state, target IDs, risk level, worker ID, dan status eksekusi.

#### GenieACS dan network

- Model server ACS, ONU/CPE, status backend-calculated, command proposal, dan command audit.
- Tabel `genieacs_servers`, `onu_devices`, dan `genieacs_command_audit`.
- Status perangkat `ONLINE`, `OFFLINE_UNDER_24H`, `OFFLINE_OVER_24H`, dan `UNKNOWN` berdasarkan timestamp inform/connection.
- Command proposal untuk `SUMMON`, `REBOOT`, `RESET`, `DELETE`, `SYNC`, dan `DHCP_OPTION_43`.
- Encoder dan decoder DHCP Option 43 tervalidasi untuk protocol HTTP/HTTPS, host, port, sub-option length, dan hexadecimal payload.

#### Dashboard dan delivery

- Dashboard terpisah untuk Developer, Administrator, Teknisi, Reseller, dan Pelanggan.
- Shared dashboard shell untuk role validation, route constants, authenticated API fetch, dan rendering shell.
- Dockerfile, Docker Compose PostgreSQL, environment example, dokumentasi API OpenAPI, migration up/down, security architecture, integration test architecture, dan Administrator module integration guide.

### Changed

- Halaman Profile Berlangganan tidak lagi bergantung pada localStorage sebagai sumber data utama dan terhubung ke API Go melalui Fetch.
- Langganan diubah menjadi pusat lifecycle pelanggan yang menghubungkan service, billing, payment, network, dan finance.
- Operasi massal dipindahkan dari eksekusi langsung menjadi proposal yang melewati preview, approval, worker, validasi, dan audit.
- Alur perubahan aplikasi mengikuti Developer → Preview → Approval → Production.
- Credential service dipisahkan dari browser; rahasia GenieACS, MikroTik, RADIUS, VPN, SNMP, dan payment gateway tidak dikirim ke frontend.

### Fixed

- Race condition dan lost update pada endpoint versioning profile melalui optimistic locking.
- Risiko authorization bypass melalui validasi JWT dan RBAC middleware yang konsisten.
- Ketidakkonsistenan audit pada proposal Administrator melalui append-only audit trail.
- Self-approval dicegah: pemohon proposal tidak dapat menyetujui proposalnya sendiri.
- Status ONU tidak lagi dipercaya dari frontend, tetapi dihitung oleh backend.
- Assertion migration CI/CD diperbaiki dengan mengganti `psql -Atx` menjadi `psql -At`, sehingga output count dapat dibandingkan sebagai nilai numerik shell.

### Security

- AI tidak memiliki akses langsung ke Production; AI hanya menghasilkan intelligence, finding, recommendation, patch preview, dan proposal.
- Perubahan yang berdampak pada pelanggan, API, firewall, credential, deployment, atau layanan network wajib melalui human-in-the-loop approval.
- Authorization fail-closed ketika principal, permission, atau engine tidak valid.
- Password dan credential rahasia tidak dikembalikan oleh API browser maupun decoder Option 43.
- Seluruh tindakan penting memiliki audit trail dengan actor, role, target, result, request ID, user-agent, dan timestamp.

### Validation

Validasi yang dijalankan pada release dan perbaikan CI/CD meliputi:

```text
go test ./...
go vet ./...
gofmt -l .
git diff --check
node --check web/dashboards/administrator/dashboard.js
node --check web/dashboards/administrator/genieacs.js
CI/CD: Go tests and PostgreSQL migration — success
CI/CD: Build and publish Docker image — success
```

### Upgrade notes

Terapkan migration database secara berurutan, termasuk migration untuk security knowledge, subscription lifecycle, GenieACS, dan Administrator modules. Pastikan `JWT_SECRET`, issuer, audience, `DATABASE_URL`, dan konfigurasi Redis tersedia pada environment server.

Executor konkret untuk MikroTik, RADIUS, GenieACS, NMS, dan payment gateway harus diinjeksi pada worker sebelum operasi eksternal diaktifkan. Browser tidak boleh mengakses credential service atau perangkat network secara langsung.

### Referensi commit utama

| Commit | Deskripsi |
|---|---|
| `4a2d122` | Backend API Subscription Profile, migration, repository, service, handler, dan test. |
| `1780b5b` | Integrasi UI Profile Berlangganan dengan API menggunakan Fetch. |
| `fef08e8` | Dockerfile, Docker Compose, environment example, dan local stack documentation. |
| `609eed0` | Security intelligence, dashboard multi-role, subscription lifecycle, GenieACS, Administrator modules, approval, audit, dan worker boundary. |
| `ed2836f` | Merge Pull Request #5 ke branch `main`. |
| `d1e3f19` | Perbaikan assertion migration pada workflow CI/CD. |

[unreleased]: https://github.com/djalil83/A-Radius/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/djalil83/A-Radius/releases/tag/v1.0.0
[Keep a Changelog]: https://keepachangelog.com/en/1.1.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html
