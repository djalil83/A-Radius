const API_BASE = '/api/v1/subscription-production';
const feedback = document.querySelector('[data-feedback]');
const integrations = document.querySelector('[data-integrations]');

async function request(path, options = {}) {
  const response = await fetch(`${API_BASE}${path}`, { ...options, headers: { 'Content-Type': 'application/json', ...(options.headers || {}) } });
  if (!response.ok) throw new Error(`Request failed: ${response.status}`);
  return response.json();
}

function renderIntegrations(items) {
  integrations.innerHTML = items.map((item) => `<div class="integration-item"><strong>${item.domain}</strong><small>${item.readiness}</small><p>${item.required ? 'REQUIRED' : 'OPTIONAL'} · ${item.failure_mode}</p></div>`).join('');
}

async function load() {
  try {
    const data = await request('/integrations');
    renderIntegrations(data.items || []);
  } catch (_) {
    integrations.innerHTML = '<div class="integration-item"><strong>Integration status unavailable</strong><small>Backend connection required</small></div>';
  }
}

document.querySelector('[data-readiness-action]').addEventListener('click', async () => {
  try {
    const result = await request('/readiness?subscription_id=preview-subscription');
    feedback.textContent = result.ready ? 'Readiness OK. Approval tetap diperlukan sebelum Production.' : `Belum siap: ${(result.blockers || []).join(', ')}.`;
  } catch (_) { feedback.textContent = 'Readiness check belum dapat dijalankan.'; }
});

document.querySelectorAll('[data-preview-action]').forEach((button) => button.addEventListener('click', async () => {
  try {
    const result = await request('/preview', { method: 'POST', body: JSON.stringify({ subscription_id: 'preview-subscription', action: button.dataset.previewAction, expected_version: 1 }) });
    feedback.textContent = `Preview ${result.status} untuk ${result.action}. Production tetap UNCHANGED.`;
  } catch (_) { feedback.textContent = 'Preview gagal dibuat. Tidak ada perubahan Production.'; }
}));

load();
