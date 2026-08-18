# Temuan Integrasi MikroTik

Sumber resmi: [RouterOS REST API](https://help.mikrotik.com/docs/spaces/ROS/pages/47579162/REST+API) dan [RouterOS API](https://help.mikrotik.com/docs/spaces/ROS/pages/47579160/API).

Dokumentasi resmi menyatakan REST API RouterOS tersedia mulai RouterOS v7.1beta4 dan menggunakan wrapper JSON pada endpoint console API, sehingga dapat melakukan operasi baca, buat, ubah, dan hapus. Layanan `www-ssl` direkomendasikan untuk akses HTTPS; layanan `www` tanpa TLS tidak boleh dipakai untuk kredensial produksi karena risiko penyadapan pasif. Integrasi A-Radius akan memakai HTTPS, timeout pendek, allowlist endpoint, dan secret yang tidak pernah dikirim ke browser.

Kontrak awal sinkronisasi akan memetakan profile A-Radius ke resource RouterOS yang relevan secara eksplisit, dimulai dari `/ppp/profile` dan dapat diperluas ke `/ip/hotspot/user/profile` setelah validasi router. Sinkronisasi dibuat idempotent dengan external identifier, mencatat status, waktu terakhir, jumlah item, error teredaksi, dan actor di database. Karena akses router merupakan operasi infrastruktur, tombol sinkronisasi manual tetap tersedia dan sinkronisasi otomatis harus dapat diaktifkan/nonaktifkan per router.

Catatan operasional: pengguna wajib mengaktifkan `www-ssl` pada RouterOS, membuat user API dengan hak minimum, dan memastikan sertifikat/TLS serta firewall hanya mengizinkan alamat server A-Radius. Implementasi tidak akan menjalankan perubahan router tanpa konfigurasi endpoint dan secret yang diberikan pengguna.
