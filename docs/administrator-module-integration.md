# A-RADIUS Administrator Module Integration

Dokumen ini mendeskripsikan integrasi modul **Voucher, Mitra, Billing, Payment, NMS, Teknisi, Finance, dan Sistem** ke dalam dashboard Administrator.

## Batas keamanan

Semua tindakan yang dapat mengubah layanan, akses, perangkat, billing, atau konfigurasi berjalan melalui alur berikut:

> Administrator → Preview → Approval oleh manusia → Redis Worker → Validasi → Audit Trail

AI hanya menghasilkan laporan, analisis, dan proposal. Field `production_changed` dipaksa `false` pada endpoint proposal dan approval. Karena itu, approval API tidak mengeksekusi perubahan produksi; eksekusi dipisahkan ke worker yang hanya menerima proposal dengan status `APPROVED`.

## Modul dan permission

| Modul | Permission | Approval default | Contoh tindakan |
|---|---|---:|---|
| Voucher | `administrator.voucher.manage` | Ya | generate, revoke, bulk activate |
| Mitra | `administrator.mitra.manage` | Ya | ubah komisi, aktif/nonaktif mitra |
| Billing | `administrator.billing.manage` | Ya | ubah status invoice, isolir, reaktivasi |
| Payment | `administrator.payment.manage` | Ya | rekonsiliasi, retry webhook |
| NMS | `administrator.nms.manage` | Ya | maintenance, konfigurasi perangkat |
| Teknisi | `administrator.teknisi.manage` | Ya | assign, close, re-open work order |
| Finance | `administrator.finance.view` | Tidak untuk baca | laporan dan rekonsiliasi |
| Sistem | `administrator.system.manage` | Ya | integrasi, konfigurasi cabang, health check |

## Endpoint

`GET /api/v1/administrator/modules` mengembalikan katalog modul berdasarkan contract backend. `GET /api/v1/administrator/ai-reports` mengembalikan laporan AI untuk tenant/cabang dari principal atau query `branch_id`. `POST /api/v1/administrator/proposals/preview` membuat proposal berstatus `PENDING_APPROVAL`. `POST /api/v1/administrator/proposals/{id}/approve` dan `/reject` hanya mengubah keputusan approval; keduanya tidak menjalankan Production.

Semua endpoint berada di balik JWT dan permission middleware. Proposal juga memiliki aturan anti-self-approval sehingga pemohon tidak dapat menyetujui proposalnya sendiri. Setiap preview dan keputusan menulis event ke `apb.admin_module_audit` dengan `request_id`, actor, target, result, dan user-agent.

## Database

Migration `0006_administrator_modules.up.sql` menyediakan tabel proposal, audit module, dan AI report. Migration `0006_administrator_modules.down.sql` menghapusnya dalam urutan dependency yang aman. Data proposal menyimpan `before_state` dan `proposed_state` untuk kebutuhan diff, rollback, dan investigasi.

## GenieACS dan Option 43

Migration `0005_genieacs_devices_and_servers.up.sql` menyediakan `genieacs_servers`, `onu_devices`, dan `genieacs_command_audit`. Password tidak dikirim dalam DTO server; hanya `credential_ref` yang boleh diteruskan ke credential service. Status ONU dihitung backend dari `last_inform_at`/`last_connected_at` menjadi `ONLINE`, `OFFLINE_UNDER_24H`, `OFFLINE_OVER_24H`, atau `UNKNOWN`.

Encoder dan decoder Option 43 berada di `internal/genieacs/option43.go`. Decoder tidak mengembalikan username/password, dan tindakan terhadap ACS tetap harus melewati command proposal.

## Worker boundary

`internal/dashboard/administrator/worker.go` menggunakan queue `a-radius:administrator:approved-actions`. Worker menolak proposal yang bukan `APPROVED` atau memiliki `production_changed=true`. Implementasi executor konkret untuk gateway, MikroTik, RADIUS, NMS, dan GenieACS perlu disuntikkan melalui interface `Executor`, sehingga credential access tetap berada di service khusus dan tidak masuk browser.
