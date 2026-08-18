from pathlib import Path

root = Path('/home/ubuntu/A-Radius-git')
roles = {
    'developer': ('Developer', 'Pengelolaan sistem, audit, feature flag, dan observability.'),
    'administrator': ('Administrator', 'Pengelolaan profil berlangganan, pengguna, approval, dan audit.'),
    'teknisi': ('Teknisi', 'Penanganan tiket, status jaringan, dan pekerjaan lapangan.'),
    'reseller': ('Reseller', 'Pengelolaan pelanggan reseller, paket, dan aktivasi layanan.'),
    'pelanggan': ('Pelanggan', 'Melihat paket, status layanan, tagihan, dan membuat tiket bantuan.'),
}

for role, (label, purpose) in roles.items():
    frontend = root / 'web' / 'dashboards' / role
    backend = root / 'internal' / 'dashboard' / role
    frontend.mkdir(parents=True, exist_ok=True)
    backend.mkdir(parents=True, exist_ok=True)
    (frontend / 'index.html').write_text(f'''<!doctype html>\n<html lang="id">\n<head>\n  <meta charset="utf-8">\n  <meta name="viewport" content="width=device-width, initial-scale=1">\n  <title>A-Radius | {label} Dashboard</title>\n</head>\n<body>\n  <div data-dashboard-root></div>\n  <script type="module" src="./dashboard.js"></script>\n</body>\n</html>\n''')
    (frontend / 'dashboard.js').write_text(f'''import {{ requireRole, renderDashboardShell }} from '../../shared/dashboard-shell.js';\n\nconst session = requireRole('{role}');\nrenderDashboardShell({{\n  role: '{role}',\n  title: '{label} Dashboard',\n  content: `<h1>{label} Dashboard</h1><p>{purpose}</p><div data-dashboard-widgets></div>`,\n}});\n\n// TODO({role}): tambahkan widget dan pemanggilan endpoint khusus dashboard ini.\nconsole.info('Dashboard {role} siap dikembangkan', session);\n''')
    (frontend / 'README.md').write_text(f'''# {label} Dashboard\n\nFolder ini hanya berisi UI dan controller untuk role **{role}**. Tujuannya adalah {purpose}\n\n## File awal\n\n| File | Tanggung jawab |\n|---|---|\n| `index.html` | Entry point dashboard |\n| `dashboard.js` | Validasi role, shell, dan titik integrasi widget |\n\nGunakan `../../shared/dashboard-shell.js` untuk autentikasi, navigasi, dan pemanggilan API. Jangan menyalin logic dari dashboard lain; logic lintas-role harus diletakkan di `web/shared`.\n''')
    (backend / 'handler.go').write_text(f'''package {role}\n\nimport (\n    "net/http"\n)\n\n// Handler menyediakan endpoint khusus untuk {label} Dashboard.\ntype Handler struct {{}}\n\n// Routes mendaftarkan route dengan middleware role=\"{role}\" pada router induk.\nfunc (h *Handler) Routes() http.Handler {{\n    mux := http.NewServeMux()\n    mux.HandleFunc("/", h.index)\n    return mux\n}}\n\nfunc (h *Handler) index(w http.ResponseWriter, _ *http.Request) {{\n    w.Header().Set("Content-Type", "application/json")\n    w.WriteHeader(http.StatusNotImplemented)\n    _, _ = w.Write([]byte(`{{"error":{{"code":"NOT_IMPLEMENTED","message":"dashboard endpoint belum diimplementasikan"}}}}`))\n}}\n''')
    (backend / 'README.md').write_text(f'''# {label} Dashboard Backend\n\nModul backend untuk role `{role}`. Handler di sini harus didaftarkan melalui router aplikasi dan dilindungi middleware RBAC.\n\nEndpoint yang spesifik untuk role `{role}` diletakkan di modul ini. Endpoint lintas-role seperti Subscription Profile tetap menggunakan modul domain terkait di `internal/subscriptionprofile`.\n''')

(root / 'web' / 'dashboards' / 'README.md').write_text('''# Dashboard A-Radius\n\nSetiap role memiliki folder frontend terpisah. Dashboard administrator menjadi pusat pengelolaan operasional, sedangkan developer mengelola aspek platform. Teknisi, reseller, dan pelanggan mengakses fungsi yang sesuai dengan role masing-masing melalui API yang sama dan kebijakan RBAC.\n''')
(root / 'internal' / 'dashboard' / 'README.md').write_text('''# Dashboard Backend Modules\n\nSetiap folder role berisi handler backend yang terisolasi. Router utama wajib memasang middleware autentikasi dan role sebelum mendaftarkan handler. Modul tidak boleh mempercayai role dari browser; role harus berasal dari session/token yang diverifikasi server.\n''')
