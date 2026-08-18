# A-Radius Web UI

`profile-berlangganan.html` adalah adaptasi UI **APB Profile Berlangganan** untuk proyek A-Radius.

## Penyesuaian

UI menggunakan branding A-Radius, namespace `localStorage` `a_radius_profile_v1`, dan nama berkas ekspor `a-radius-profile.json`. Fitur lokal yang tersedia mencakup tab General/Network/Billing, filter profile, pencarian, tambah/edit/hapus profile, konfirmasi perubahan, serta ekspor JSON.

## Batasan saat ini

UI ini masih berjalan sebagai halaman statis dan belum terhubung ke API atau PostgreSQL. Data form disimpan pada browser melalui `localStorage`. Integrasi produksi perlu menambahkan endpoint profile berlangganan, validasi server-side, autentikasi/otorisasi, approval queue, dan audit trail sebelum operasi RADIUS, MikroTik, OLT/ACS, atau billing dijalankan.

## Menjalankan secara lokal

Dari root proyek, jalankan server file statis apa pun, kemudian buka `web/profile-berlangganan.html`. Contoh sederhana:

```bash
python3 -m http.server 8080 --directory web
```

Jangan menaruh credential database, RADIUS, MikroTik, OLT, atau payment gateway di HTML maupun JavaScript frontend.
