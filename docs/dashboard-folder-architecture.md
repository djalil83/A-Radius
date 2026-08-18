# Arsitektur Folder Dashboard A-Radius

## Tujuan

Struktur dashboard A-Radius dipisahkan berdasarkan peran agar kode UI, endpoint, dan aturan bisnis tiap jenis pengguna dapat dikembangkan secara mandiri. Pemisahan folder tidak berarti setiap dashboard memiliki backend atau database terpisah. Seluruh dashboard tetap menggunakan API domain A-Radius yang sama, sedangkan aksesnya dibatasi oleh autentikasi dan RBAC di sisi server.

> **Prinsip keamanan:** role yang dikirim dari browser tidak boleh dianggap sebagai sumber kebenaran. Server harus mengambil role dari session atau token yang telah diverifikasi.

## Struktur folder

```text
web/
├── shared/
│   └── dashboard-shell.js
└── dashboards/
    ├── developer/
    │   ├── index.html
    │   ├── dashboard.js
    │   └── README.md
    ├── administrator/
    │   ├── index.html
    │   ├── dashboard.js
    │   └── README.md
    ├── teknisi/
    │   ├── index.html
    │   ├── dashboard.js
    │   └── README.md
    ├── reseller/
    │   ├── index.html
    │   ├── dashboard.js
    │   └── README.md
    └── pelanggan/
        ├── index.html
        ├── dashboard.js
        └── README.md

internal/dashboard/
├── shared/contract.go
├── developer/handler.go
├── administrator/handler.go
├── teknisi/handler.go
├── reseller/handler.go
└── pelanggan/handler.go
```

## Hubungan antar-dashboard

| Dashboard | Pengguna utama | Akses utama | Terhubung dengan |
|---|---|---|---|
| Developer | Tim pengembang/platform | konfigurasi sistem, audit, observability | Administrator melalui audit dan konfigurasi platform |
| Administrator | Pengelola ISP | pengguna, approval, profil berlangganan, audit | Developer untuk platform; teknisi, reseller, dan pelanggan untuk operasional |
| Teknisi | Tim operasional jaringan | tiket, status jaringan, pekerjaan lapangan | Administrator untuk penugasan; pelanggan untuk penyelesaian tiket |
| Reseller | Mitra penjualan | pelanggan reseller, paket, aktivasi | Administrator untuk approval; pelanggan untuk layanan yang dijual |
| Pelanggan | Pengguna akhir | paket, status layanan, tagihan, tiket | Administrator/reseller melalui data layanan; teknisi melalui tiket |

Administrator berfungsi sebagai pusat operasi, tetapi tidak menjadi tempat penyalinan kode dashboard lain. Hubungan antar-dashboard terjadi melalui API, event/audit, dan aturan izin yang terpusat.

## Kontrak endpoint

Modul `internal/dashboard/shared/contract.go` menyimpan metadata role, base path, dan permission sebagai referensi route. Handler khusus role berada di foldernya sendiri. Domain umum seperti Subscription Profile tetap berada di `internal/subscriptionprofile` sehingga administrator, reseller, dan pelanggan dapat menggunakan kontrak API yang konsisten.

Frontend menggunakan `web/shared/dashboard-shell.js` untuk `requireRole`, `apiFetch`, dan shell navigasi. Dengan demikian, perubahan autentikasi atau format error dapat dilakukan di satu tempat.

## Aturan implementasi berikutnya

Setiap fitur baru harus ditempatkan pada dashboard yang memiliki tanggung jawab terhadap fitur tersebut. Komponen yang dipakai dua atau lebih dashboard dipindahkan ke `web/shared`. Query database dan aturan domain tidak boleh diletakkan di handler role; handler hanya mengorkestrasi request, validasi izin, dan pemanggilan service.

Untuk menghubungkan route ke router utama, daftarkan handler role setelah middleware autentikasi dan RBAC. Route yang bersifat lintas-role harus menggunakan permission, bukan sekadar nama dashboard. Perubahan schema tetap dilakukan melalui migration PostgreSQL yang sudah digunakan proyek.
