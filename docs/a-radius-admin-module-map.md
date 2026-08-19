# A-RADIUS Administrator Module Map

## Domain structure

A-RADIUS memiliki satu dashboard Administrator yang mengorkestrasi Pelanggan, Langganan, Service, Router & Server, GenieACS, Voucher, Mitra, Billing, Payment, NMS, Teknisi, Finance, dan Sistem. Setiap domain memiliki ownership serta permission sendiri; Langganan tetap menjadi pusat lifecycle pelanggan, bukan database terpisah.

## GenieACS Administrator

Entry point frontend tersedia di `web/dashboards/administrator/genieacs.html`. Dashboard menampilkan total ONU, online, offline di bawah 24 jam, offline di atas 24 jam, filter server/status, pencarian ID/SN/username, dan tabel device.

Command `SUMMON`, `REBOOT`, `DELETE`, `RESET`, `SINKRON`, dan `DHCP43` tidak dieksekusi langsung dari browser. Browser hanya membuat command proposal preview. Command produksi harus melalui JWT, RBAC, target scope, approval, worker, validation, dan audit trail. Command `DELETE` dan `RESET` memerlukan perlindungan tambahan karena dapat mengganggu perangkat pelanggan.

## Integrasi dengan lifecycle Langganan

Status Langganan menjadi sumber keputusan untuk Billing, RADIUS/MikroTik, Customer App, dan monitoring GenieACS. Perubahan status atau command device harus dapat ditelusuri menggunakan correlation ID dan audit event.

## Modul lanjutan

Menu Voucher, Mitra, NMS, Teknisi, Finance, dan Sistem dapat dikembangkan sebagai bounded modules. Integrasi lintas-domain dilakukan melalui API/domain events, bukan akses langsung antartabel yang melewati ownership.
