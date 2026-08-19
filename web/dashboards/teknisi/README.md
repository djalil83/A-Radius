# Teknisi Dashboard

Folder ini hanya berisi UI dan controller untuk role **teknisi**. Tujuannya adalah Penanganan tiket, status jaringan, dan pekerjaan lapangan.

## File awal

| File | Tanggung jawab |
|---|---|
| `index.html` | Entry point dashboard |
| `dashboard.js` | Validasi role, shell, dan titik integrasi widget |

Gunakan `../../shared/dashboard-shell.js` untuk autentikasi, navigasi, dan pemanggilan API. Jangan menyalin logic dari dashboard lain; logic lintas-role harus diletakkan di `web/shared`.
