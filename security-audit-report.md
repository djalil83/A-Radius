# Laporan Audit Dependensi dan Keamanan A-Radius

**Tanggal pemeriksaan:** 18 Agustus 2026  
**Repositori:** [djalil83/A-Radius](https://github.com/djalil83/A-Radius)  
**Branch:** `main`  
**Pemeriksa:** Manus AI

## Ringkasan eksekutif

Proyek berhasil melewati `go mod verify`, seluruh pengujian `go test ./...`, dan `go vet ./...`. Namun, pemindaian menggunakan `govulncheck` menemukan **23 kerentanan yang dijangkau oleh call graph** pada Go standard library dan satu dependensi pihak ketiga, serta **12 temuan tambahan** pada paket yang diimpor tetapi tidak memiliki jejak panggilan sampai ke penggunaan rentan.

Temuan yang paling jelas perlu ditindaklanjuti adalah `golang.org/x/text` versi `v0.29.0`, yang digunakan melalui jalur normalisasi Unicode dan memiliki perbaikan mulai `v0.39.0`. Versi terbaru yang terdeteksi pada pkg.go.dev adalah `v0.41.0`, sehingga dependensi ini sebaiknya diperbarui setelah pengujian regresi.

Masalah yang lebih luas berasal dari toolchain Go yang berjalan sebagai `go1.25` tanpa patch version. Hasil audit menunjukkan kerentanan standard library yang diperbaiki bertahap hingga `go1.25.13`; rilis resmi Go juga mencatat bahwa `go1.25.13`, dirilis 13 Agustus 2026, memuat perbaikan keamanan pada `crypto/tls`, `encoding/asn1`, `encoding/xml`, `net/http`, dan `net/url`. Untuk lingkungan baru, Go `1.26.6` adalah cabang mayor terbaru yang tercantum pada riwayat rilis resmi.

## Inventarisasi dependensi

| Modul | Versi saat ini | Status pemeriksaan | Rekomendasi |
|---|---:|---|---|
| `github.com/jackc/pgx/v5` | `v5.10.0` | Tidak ada pembaruan yang dilaporkan oleh `go list -m -u all`; tidak muncul sebagai temuan rentan langsung | Pertahankan, tetapi pantau rilis upstream |
| `github.com/jackc/pgpassfile` | `v1.0.0` | Dependensi tidak langsung; tidak ada pembaruan yang dilaporkan | Pertahankan |
| `github.com/jackc/pgservicefile` | `v0.0.0-20240606120523-5a60cdf6a761` | Dependensi tidak langsung; tidak ada pembaruan yang dilaporkan | Pertahankan |
| `github.com/jackc/puddle/v2` | `v2.2.2` | Dependensi tidak langsung; tidak ada pembaruan yang dilaporkan | Pertahankan |
| `golang.org/x/sync` | `v0.17.0` | Pembaruan tersedia hingga `v0.22.0`; tidak dilaporkan sebagai kerentanan terjangkau dalam pemindaian ini | Perbarui setelah uji kompatibilitas |
| `golang.org/x/text` | `v0.29.0` | **Rentan: GO-2026-5970**, fixed mulai `v0.39.0`; jejak panggilan ditemukan melalui `norm.Form.Properties`, `Span`, dan `Transform` | **Prioritas tinggi: perbarui minimal ke `v0.39.0`, disarankan `v0.41.0`** |

`go.mod` juga menyatakan `go 1.25.0`. Karena Go menggunakan patch release untuk menyalurkan perbaikan keamanan standard library, deklarasi versi mayor-minor saja tidak cukup untuk menjamin binary dibangun dengan patch keamanan terbaru; pipeline build perlu mematok toolchain yang sudah dipatch.

## Hasil validasi lokal

| Pemeriksaan | Hasil |
|---|---|
| `go mod verify` | Lulus; seluruh modul terverifikasi |
| `go test ./...` | Lulus; paket `internal/authz` dan `internal/securityknowledge` berhasil diuji |
| `go vet ./...` | Lulus tanpa temuan |
| `go list -m -u all` | Menemukan pembaruan `x/text` hingga `v0.41.0` dan `x/sync` hingga `v0.22.0`; `pgx/v5` tidak menunjukkan pembaruan |
| `govulncheck ./...` | Gagal dengan exit code 3 karena temuan kerentanan |

## Temuan keamanan utama

### 1. `golang.org/x/text` — GO-2026-5970 — prioritas tinggi

`golang.org/x/text@v0.29.0` terdeteksi pada jalur panggilan proyek dan diperbaiki mulai `v0.39.0`. Deskripsi resmi menyatakan bahwa input UTF-8 tidak valid dapat menyebabkan loop tak berujung pada API normalisasi Unicode. Jejak yang dilaporkan mengarah dari `cmd/securityknowledge-loader/main.go:21:21` melalui `sql.Open`, kemudian ke `norm.Form.Properties`, `norm.Form.Span`, dan `norm.Form.Transform`.

Risiko praktis bergantung pada apakah loader dapat memproses data atau manifest yang dikendalikan pihak luar. Jika input tersebut tidak dipercaya, dampak yang paling relevan adalah denial of service melalui penggunaan CPU yang tidak berhenti. Perbaikan yang disarankan adalah memperbarui ke `golang.org/x/text@v0.41.0`, lalu menjalankan kembali test dan audit.

### 2. Go standard library — 23 temuan terjangkau — prioritas tinggi

Sebanyak 23 temuan berada pada `go1.25` standard library, termasuk area `crypto/tls`, `crypto/x509`, `encoding/asn1`, `net`, `net/url`, dan `os`. Temuan lama seperti GO-2025-4007 dan GO-2025-4175 memiliki versi perbaikan masing-masing mulai `go1.25.3` dan `go1.25.5`, tetapi hasil audit yang menggunakan database kerentanan bertanggal 2026 juga menemukan temuan baru yang membutuhkan patch hingga `go1.25.13`.

Rilis resmi Go `1.25.13` pada 13 Agustus 2026 mencakup perbaikan keamanan pada `crypto/tls`, `encoding/asn1`, `encoding/xml`, `net/http`, dan `net/url`. Oleh karena itu, jangan berhenti pada `go1.25.5`; gunakan minimal `go1.25.13` untuk branch Go 1.25, atau evaluasi migrasi ke `go1.26.6` sebagai rilis mayor terbaru.

### 3. Temuan informasional pada paket yang diimpor — prioritas sedang/rendah

Pemindaian juga melaporkan 12 kerentanan pada paket yang diimpor tetapi tidak memiliki call stack yang menunjukkan kode proyek menggunakan simbol rentan. Temuan jenis ini tidak otomatis berarti eksploitasi dapat terjadi pada aplikasi, tetapi tetap perlu dipantau karena jalur panggilan dapat berubah ketika kode atau dependensi diperbarui.

## Rencana perbaikan yang disarankan

Pertama, gunakan toolchain Go minimal `1.25.13` pada CI dan lingkungan produksi, atau `1.26.6` setelah kompatibilitas dikonfirmasi. Kedua, perbarui `golang.org/x/text` ke `v0.41.0` dan pertimbangkan memperbarui `golang.org/x/sync` ke `v0.22.0` karena pembaruan tersedia. Ketiga, jalankan `go test ./...`, `go vet ./...`, `go mod verify`, dan `govulncheck ./...` kembali setelah perubahan. Keempat, tambahkan pemindaian `govulncheck` ke CI agar pembaruan database kerentanan tidak bergantung pada pemeriksaan manual.

Perubahan dependensi **belum diterapkan atau di-push** ke GitHub dalam pemeriksaan ini. Itu sengaja dilakukan agar tidak mengubah perilaku proyek tanpa persetujuan pengguna.

## Kesimpulan

Proyek dapat dibangun dan diuji dengan sukses, tetapi **belum bebas kerentanan menurut database Go per 18 Agustus 2026**. Tindakan paling penting adalah memutakhirkan toolchain ke patch Go terbaru dan memperbarui `golang.org/x/text` setidaknya ke `v0.39.0`, idealnya `v0.41.0`. Setelah itu, audit ulang harus menunjukkan pengurangan temuan terjangkau; bila masih ada temuan standard library, pastikan binary benar-benar dibangun menggunakan toolchain patch yang dipilih.

## Referensi

[1]: https://go.dev/doc/tutorial/govulncheck "Tutorial resmi Go: Find and fix vulnerable dependencies with govulncheck"
[2]: https://osv.dev/vulnerability/GO-2026-5970 "OSV: GO-2026-5970"
[3]: https://pkg.go.dev/vuln/GO-2025-4007 "Go Vulnerability Database: GO-2025-4007"
[4]: https://pkg.go.dev/vuln/GO-2025-4175 "Go Vulnerability Database: GO-2025-4175"
[5]: https://go.dev/doc/devel/release "Go Release History"
[6]: https://pkg.go.dev/golang.org/x/text "pkg.go.dev: golang.org/x/text"
[7]: https://github.com/djalil83/A-Radius "Repositori A-Radius"

---

**Catatan metodologi:** Pemeriksaan dilakukan terhadap checkout lokal branch `main` menggunakan `go mod verify`, `go test ./...`, `go vet ./...`, `go list -m -u all`, dan `govulncheck -show verbose ./...`. Hasil `govulncheck` bersifat call-graph-aware; temuan “terjangkau” dan temuan informasional harus diprioritaskan berbeda sesuai jalur input aplikasi.
