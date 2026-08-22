export function renderStatusBadge(status, options = {}) {
  const value = String(status ?? 'UNKNOWN').trim() || 'UNKNOWN';
  const tone = String(options.tone || value).toLowerCase();

  const badge = document.createElement('span');
  badge.className = `developer-status-badge developer-status-badge--${tone}`;
  badge.textContent = value;
  badge.setAttribute('role', 'status');

  return badge;
}

export function statusBadgeHTML(status, options = {}) {
  const value = String(status ?? 'UNKNOWN').trim() || 'UNKNOWN';
  const tone = String(options.tone || value).toLowerCase();

  return `<span class="developer-status-badge developer-status-badge--${tone}" role="status">${escapeHTML(value)}</span>`;
}

function escapeHTML(value) {
  return String(value)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}
