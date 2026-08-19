import { apiFetch, renderDashboardShell, requireRole } from '../../shared/dashboard-shell.js';

const session = requireRole('administrator');
const moduleFallback = [
  ['VOUCHER', 'Voucher', 'Pembuatan dan distribusi voucher.'],
  ['MITRA', 'Mitra', 'Reseller, komisi, dan relasi cabang.'],
  ['BILLING', 'Billing', 'Invoice, jatuh tempo, dan isolir.'],
  ['PAYMENT', 'Payment', 'Pembayaran dan rekonsiliasi gateway.'],
  ['NMS', 'NMS', 'Monitoring perangkat dan alarm jaringan.'],
  ['TEKNISI', 'Teknisi', 'Work order dan penugasan lapangan.'],
  ['FINANCE', 'Finance', 'Pendapatan dan laporan keuangan.'],
  ['SISTEM', 'Sistem', 'Konfigurasi, integrasi, dan health check.'],
];

const esc = (value = '') => String(value).replace(/[&<>'"]/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', "'": '&#39;', '"': '&quot;' }[c]));

renderDashboardShell({
  role: 'administrator',
  title: 'Administrator Command Center',
  content: `<div class="admin-layout">
    <aside class="admin-sidebar"><div class="brand">A-RADIUS <small>ADMINISTRATOR</small></div><nav id="module-nav"></nav><div class="production-lock">PRODUCTION<br><strong>APPROVAL REQUIRED</strong></div></aside>
    <section class="admin-main"><div class="admin-heading"><div><span class="eyebrow">CABANG / ADMINISTRATOR</span><h1>Operational Management Hub</h1><p>Semua perubahan sensitif berjalan melalui Developer → Preview → Approval → Worker.</p></div><span class="safe-badge">PRODUCTION UNCHANGED</span></div>
      <div id="module-grid" class="module-grid"></div>
      <section class="ai-panel"><div class="section-title"><div><span class="eyebrow">CONTINUOUS SECURITY INTELLIGENCE</span><h2>AI REPORT</h2></div><span class="muted">AI merekomendasikan, manusia menyetujui</span></div><div id="ai-reports" class="reports"></div></section>
      <section class="proposal-panel"><div class="section-title"><h2>PROPOSAL PREVIEW</h2><span id="proposal-state" class="muted">Belum ada proposal</span></div><form id="proposal-form"><div class="form-row"><label>Modul<select id="proposal-module" required></select></label><label>Aksi<input id="proposal-action" required placeholder="Contoh: ISOLATE_EXPIRED"></label><label>Risk<select id="proposal-risk"><option>LOW</option><option>MEDIUM</option><option>HIGH</option><option>CRITICAL</option></select></label></div><label>Target IDs<textarea id="proposal-targets" placeholder="Satu ID per baris"></textarea></label><label>Alasan<textarea id="proposal-reason" required placeholder="Jelaskan tujuan dan dampak"></textarea></label><button class="primary" type="submit">CREATE PREVIEW</button></form></section><div id="feedback" role="status"></div>
    </section></div>`,
});

document.head.insertAdjacentHTML('beforeend', '<link rel="stylesheet" href="./administrator.css">');

const nav = document.querySelector('#module-nav');
const grid = document.querySelector('#module-grid');
const reports = document.querySelector('#ai-reports');
const moduleSelect = document.querySelector('#proposal-module');
const feedback = document.querySelector('#feedback');

function renderModules(modules = moduleFallback.map(([key, label, description]) => ({ key, label, description }))) {
  nav.innerHTML = modules.map((m) => `<a href="#${esc(m.key)}">${esc(m.label)}</a>`).join('');
  grid.innerHTML = modules.map((m) => `<article class="module-card" id="${esc(m.key)}"><span class="module-key">${esc(m.key)}</span><h3>${esc(m.label)}</h3><p>${esc(m.description)}</p><button data-module="${esc(m.key)}">CREATE PREVIEW</button></article>`).join('');
  moduleSelect.innerHTML = modules.map((m) => `<option value="${esc(m.key)}">${esc(m.label)}</option>`).join('');
  grid.querySelectorAll('[data-module]').forEach((button) => button.addEventListener('click', () => { moduleSelect.value = button.dataset.module; document.querySelector('#proposal-form').scrollIntoView({ behavior: 'smooth' }); }));
}

function renderReports(items = []) {
  reports.innerHTML = items.length ? items.map((report) => `<article class="report-card severity-${esc(report.severity)}"><div class="report-top"><strong>${esc(report.title)}</strong><span>${esc(report.severity)}</span></div><p>${esc(report.finding)}</p><small>Recommendation: ${esc(report.recommendation)}</small><div class="report-actions"><button data-report-action="analysis">VIEW ANALYSIS</button><button data-report-action="preview" data-report-module="${esc(report.module)}">CREATE FIX PREVIEW</button></div><footer>Production: <strong>UNCHANGED</strong></footer></article>`).join('') : '<div class="empty">Belum ada AI Report untuk cabang ini.</div>';
  reports.querySelectorAll('[data-report-action="preview"]').forEach((button) => button.addEventListener('click', () => { moduleSelect.value = button.dataset.reportModule; document.querySelector('#proposal-action').value = 'SECURITY_REMEDIATION_PREVIEW'; document.querySelector('#proposal-reason').value = 'Preview remediation dari AI Report; wajib ditinjau Administrator/Developer.'; document.querySelector('#proposal-form').scrollIntoView({ behavior: 'smooth' }); }));
}

async function load() {
  renderModules(); renderReports();
  try { const modules = await apiFetch('/administrator/modules'); renderModules(modules.modules); } catch (error) { feedback.textContent = `Modul lokal aktif; katalog API belum tersedia (${error.code || error.message}).`; }
  try { const data = await apiFetch('/administrator/ai-reports'); renderReports(data.reports); } catch (error) { feedback.textContent = `AI Report belum dapat dimuat (${error.code || error.message}).`; }
}

document.querySelector('#proposal-form').addEventListener('submit', async (event) => {
  event.preventDefault();
  const targetIDs = document.querySelector('#proposal-targets').value.split(/\n|,/).map((x) => x.trim()).filter(Boolean);
  const payload = { branch_id: session.tenant_id || session.branch_id || '0', module: moduleSelect.value, action: document.querySelector('#proposal-action').value.trim(), target_type: 'ADMINISTRATOR_TARGET', target_ids: targetIDs, before_state: {}, proposed_state: {}, risk_level: document.querySelector('#proposal-risk').value, reason: document.querySelector('#proposal-reason').value.trim() };
  try { const result = await apiFetch('/administrator/proposals/preview', { method: 'POST', body: JSON.stringify(payload) }); document.querySelector('#proposal-state').textContent = `PENDING APPROVAL · ${result.proposal_id}`; feedback.textContent = 'Preview tersimpan. Tidak ada perubahan Production. Menunggu approval manusia.'; event.target.reset(); } catch (error) { feedback.textContent = `Preview gagal: ${error.code || error.message}`; }
});

load();
