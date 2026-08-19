import { apiFetch as sharedApiFetch } from '/shared/dashboard-shell.js';

const root = document.querySelector('[data-dashboard-root]');

if (!root) {
  throw new Error('Elemen [data-dashboard-root] tidak ditemukan');
}

document.title = 'A-Radius | Dashboard Pelanggan';

root.innerHTML = `
<style>
  :root {
    --primary: #0f172a;
    --primary-soft: #1e293b;
    --accent: #2563eb;
    --success: #16a34a;
    --success-soft: #f0fdf4;
    --danger: #dc2626;
    --border: #e2e8f0;
    --muted: #64748b;
    --surface: #ffffff;
    --background: #f8fafc;
  }

  * {
    box-sizing: border-box;
  }

  body {
    margin: 0;
    background: var(--background);
    color: var(--primary);
    font-family: Inter, ui-sans-serif, system-ui, -apple-system,
      BlinkMacSystemFont, "Segoe UI", sans-serif;
  }

  .customer-dashboard {
    min-height: 100vh;
  }

  .dashboard-header {
    background: linear-gradient(135deg, #0f172a, #1e293b);
    color: white;
    padding: 22px clamp(18px, 4vw, 48px);
  }

  .header-inner {
    max-width: 1280px;
    margin: auto;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 20px;
  }

  .brand {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .brand-icon {
    width: 42px;
    height: 42px;
    display: grid;
    place-items: center;
    border-radius: 12px;
    background: rgba(255,255,255,.12);
    font-weight: 800;
  }

  .brand-title {
    font-weight: 800;
    font-size: 18px;
  }

  .brand-subtitle {
    color: #cbd5e1;
    font-size: 12px;
    margin-top: 2px;
  }

  .status-pill {
    display: inline-flex;
    align-items: center;
    gap: 7px;
    border-radius: 999px;
    padding: 8px 12px;
    font-size: 12px;
    font-weight: 700;
    background: rgba(255,255,255,.1);
  }

  .status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #22c55e;
  }

  .container {
    max-width: 1280px;
    margin: auto;
    padding: 28px clamp(18px, 4vw, 48px) 48px;
  }

  .welcome {
    margin-bottom: 24px;
  }

  .welcome h1 {
    margin: 0;
    font-size: clamp(25px, 4vw, 34px);
    letter-spacing: -.02em;
  }

  .welcome p {
    margin: 7px 0 0;
    color: var(--muted);
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 16px;
    margin-bottom: 24px;
  }

  .card {
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 16px;
    box-shadow: 0 5px 20px rgba(15,23,42,.04);
  }

  .summary-card {
    padding: 20px;
  }

  .summary-label {
    color: var(--muted);
    font-size: 13px;
    font-weight: 600;
  }

  .summary-value {
    margin-top: 8px;
    font-size: 30px;
    font-weight: 800;
  }

  .summary-caption {
    margin-top: 4px;
    color: var(--muted);
    font-size: 12px;
  }

  .section {
    padding: 22px;
    margin-bottom: 18px;
  }

  .section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    margin-bottom: 18px;
  }

  .section-title {
    margin: 0;
    font-size: 18px;
    font-weight: 800;
  }

  .section-subtitle {
    color: var(--muted);
    font-size: 12px;
    margin-top: 4px;
  }

  .profile-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 14px;
  }

  .field {
    padding: 14px;
    border: 1px solid var(--border);
    border-radius: 12px;
    background: #fcfdff;
  }

  .field-label {
    color: var(--muted);
    font-size: 11px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: .04em;
  }

  .field-value {
    margin-top: 5px;
    font-size: 14px;
    font-weight: 650;
    word-break: break-word;
  }

  .service {
    border: 1px solid var(--border);
    border-radius: 14px;
    padding: 17px;
    margin-bottom: 12px;
  }

  .service:last-child {
    margin-bottom: 0;
  }

  .service-top {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 15px;
  }

  .service-name {
    font-weight: 800;
    font-size: 16px;
  }

  .service-code {
    color: var(--muted);
    font-size: 12px;
    margin-top: 4px;
  }

  .badge {
    display: inline-flex;
    align-items: center;
    border-radius: 999px;
    padding: 6px 10px;
    font-size: 11px;
    font-weight: 800;
  }

  .badge-active {
    color: #166534;
    background: var(--success-soft);
  }

  .badge-inactive {
    color: #92400e;
    background: #fffbeb;
  }

  .service-meta {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 10px;
    margin-top: 15px;
  }

  .meta {
    padding: 11px;
    background: var(--background);
    border-radius: 10px;
  }

  .meta-label {
    color: var(--muted);
    font-size: 11px;
  }

  .meta-value {
    margin-top: 3px;
    font-weight: 750;
    font-size: 13px;
  }

  .loading,
  .error {
    max-width: 1280px;
    margin: 50px auto;
    padding: 28px;
    text-align: center;
  }

  .error {
    color: var(--danger);
  }

  .retry {
    margin-top: 14px;
    border: 0;
    border-radius: 10px;
    padding: 10px 16px;
    background: var(--accent);
    color: white;
    font-weight: 700;
    cursor: pointer;
  }

  @media (max-width: 800px) {
    .grid,
    .profile-grid {
      grid-template-columns: 1fr;
    }

    .service-meta {
      grid-template-columns: 1fr;
    }

    .header-inner {
      align-items: flex-start;
    }
  }
</style>

<div class="customer-dashboard">
  <header class="dashboard-header">
    <div class="header-inner">
      <div class="brand">
        <div class="brand-icon">AR</div>
        <div>
          <div class="brand-title">A-Radius</div>
          <div class="brand-subtitle">Portal Pelanggan</div>
        </div>
      </div>

      <div class="status-pill">
        <span class="status-dot"></span>
        Akun Aktif
      </div>
    </div>
  </header>

  <main class="container">
    <div id="dashboard-content">
      <div class="card loading">
        Memuat dashboard pelanggan...
      </div>
    </div>
  </main>
</div>
`;

const content = document.querySelector('#dashboard-content');

function escapeHTML(value) {
  return String(value ?? '')
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}

function formatSpeed(kbps) {
  const value = Number(kbps || 0);

  if (value >= 1000000) {
    return `${(value / 1000000).toFixed(1)} Gbps`;
  }

  if (value >= 1000) {
    return `${(value / 1000).toFixed(value % 1000 === 0 ? 0 : 1)} Mbps`;
  }

  return `${value} Kbps`;
}

const apiFetch = sharedApiFetch;


function renderDashboard(data) {
  const customer = data.customer || {};
  const services = Array.isArray(data.services) ? data.services : [];
  const summary = data.summary || {};

  const fullAddress = [
    customer.address,
    customer.village,
    customer.district,
    customer.regency,
    customer.province,
    customer.postal_code,
  ]
    .filter(Boolean)
    .join(', ');

  const serviceHTML = services.length
    ? services.map(service => `
      <article class="service">
        <div class="service-top">
          <div>
            <div class="service-name">
              ${escapeHTML(service.package_name || service.service_type || 'Layanan')}
            </div>
            <div class="service-code">
              ${escapeHTML(service.service_code || '-')}
            </div>
          </div>

          <span class="badge ${
            service.status === 'active'
              ? 'badge-active'
              : 'badge-inactive'
          }">
            ${escapeHTML(service.status || 'unknown')}
          </span>
        </div>

        <div class="service-meta">
          <div class="meta">
            <div class="meta-label">Jenis</div>
            <div class="meta-value">
              ${escapeHTML(service.service_type || '-')}
            </div>
          </div>

          <div class="meta">
            <div class="meta-label">Download</div>
            <div class="meta-value">
              ${formatSpeed(service.download_speed)}
            </div>
          </div>

          <div class="meta">
            <div class="meta-label">Upload</div>
            <div class="meta-value">
              ${formatSpeed(service.upload_speed)}
            </div>
          </div>
        </div>
      </article>
    `).join('')
    : `
      <div class="field">
        <div class="field-value">Belum ada layanan.</div>
      </div>
    `;

  content.innerHTML = `
    <div class="welcome">
      <h1>Halo, ${escapeHTML(customer.name || 'Pelanggan')}</h1>
      <p>Berikut ringkasan akun dan layanan internet Anda.</p>
    </div>

    <section class="grid">
      <div class="card summary-card">
        <div class="summary-label">Total Layanan</div>
        <div class="summary-value">
          ${Number(summary.service_count || services.length)}
        </div>
        <div class="summary-caption">Layanan terdaftar</div>
      </div>

      <div class="card summary-card">
        <div class="summary-label">Layanan Aktif</div>
        <div class="summary-value">
          ${Number(summary.active_service_count || 0)}
        </div>
        <div class="summary-caption">Sedang berjalan</div>
      </div>

      <div class="card summary-card">
        <div class="summary-label">Status Pelanggan</div>
        <div class="summary-value" style="font-size:22px">
          ${escapeHTML(customer.status || '-')}
        </div>
        <div class="summary-caption">
          ${escapeHTML(customer.customer_code || '-')}
        </div>
      </div>
    </section>

    <section class="card section">
      <div class="section-header">
        <div>
          <h2 class="section-title">Profil Pelanggan</h2>
          <div class="section-subtitle">
            Informasi akun yang terhubung dengan sesi login.
          </div>
        </div>
      </div>

      <div class="profile-grid">
        <div class="field">
          <div class="field-label">Nama</div>
          <div class="field-value">
            ${escapeHTML(customer.name || '-')}
          </div>
        </div>

        <div class="field">
          <div class="field-label">Kode Pelanggan</div>
          <div class="field-value">
            ${escapeHTML(customer.customer_code || '-')}
          </div>
        </div>

        <div class="field">
          <div class="field-label">Email</div>
          <div class="field-value">
            ${escapeHTML(customer.email || '-')}
          </div>
        </div>

        <div class="field">
          <div class="field-label">Telepon</div>
          <div class="field-value">
            ${escapeHTML(customer.phone || '-')}
          </div>
        </div>

        <div class="field" style="grid-column:1/-1">
          <div class="field-label">Alamat</div>
          <div class="field-value">
            ${escapeHTML(fullAddress || '-')}
          </div>
        </div>
      </div>
    </section>

    <section class="card section">
      <div class="section-header">
        <div>
          <h2 class="section-title">Layanan Internet</h2>
          <div class="section-subtitle">
            Layanan yang dimiliki oleh akun pelanggan ini.
          </div>
        </div>
      </div>

      ${serviceHTML}
    </section>
  `;
}

function renderError(error) {
  let message = 'Dashboard tidak dapat dimuat.';

  if (error?.status === 401) {
    message = 'Sesi login tidak valid atau sudah berakhir.';
  } else if (error?.status === 403) {
    message =
      'Akun Anda tidak memiliki izin mengakses dashboard pelanggan.';
  } else if (error?.status >= 500) {
    message =
      'Server sedang mengalami masalah. Silakan coba kembali.';
  }

  content.innerHTML = `
    <div class="card error">
      <strong>${escapeHTML(message)}</strong>
      <div>
        <button class="retry" type="button" id="retry-dashboard">
          Coba Lagi
        </button>
      </div>
    </div>
  `;

  document
    .querySelector('#retry-dashboard')
    ?.addEventListener('click', loadDashboard);
}

async function loadDashboard() {
  content.innerHTML = `
    <div class="card loading">
      Memuat dashboard pelanggan...
    </div>
  `;

  try {
    const dashboard = await apiFetch('/customer/dashboard');
    renderDashboard(dashboard);
  } catch (error) {
    console.error('Customer dashboard error:', error);
    renderError(error);
  }
}

loadDashboard();
