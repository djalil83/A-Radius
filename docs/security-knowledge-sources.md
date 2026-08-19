# Trusted Security Knowledge Sources

Catatan sumber resmi untuk A-RADIUS Continuous Security Intelligence.

| Sumber | Kegunaan | URL |
|---|---|---|
| CISA Known Exploited Vulnerabilities Catalog | Prioritas kerentanan yang diketahui dieksploitasi; kandidat eskalasi risiko tinggi | https://www.cisa.gov/known-exploited-vulnerabilities-catalog |
| OWASP API Security Project | Knowledge API security dan OWASP API Security Top 10 2023, termasuk broken object/function authorization, broken authentication, unrestricted resource consumption, SSRF, misconfiguration, inventory, dan unsafe API consumption | https://owasp.org/www-project-api-security/ |
| OASIS CSAF 2.0 | Format terstruktur yang interoperable untuk security advisories, impact, remediation, dan status produk | https://docs.oasis-open.org/csaf/csaf/v2.0/os/csaf-v2.0-os.html |
| OWASP Dependency-Check | Referensi SCA: menghubungkan dependency evidence/CPE dengan CVE dan memakai NVD feeds; dapat dijalankan berkala | https://owasp.org/www-project-dependency-check/ |

## Implementasi provenance

Setiap knowledge item wajib menyimpan source URL, publisher, retrieved_at, published_at bila tersedia, advisory/CVE identifier, content hash, parser version, validation status, dan confidence. AI tidak boleh menjadikan satu feed mentah sebagai fakta final tanpa filter, deduplikasi, schema validation, dan mapping ke package/version/product A-Radius.

## Prioritas awal

CISA KEV dipakai sebagai sinyal prioritas eksploitasi aktif. OWASP menjadi knowledge taxonomy untuk API dan application risk. CSAF menjadi kandidat format ingestion advisory terstruktur. Dependency-Check/NVD digunakan untuk evidence dependency dan CVE correlation. Semua hasil tersebut tetap menghasilkan finding atau proposal; tidak ada sumber eksternal yang dapat men-trigger deployment otomatis.
