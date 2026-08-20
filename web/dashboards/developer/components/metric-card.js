import { statusBadgeHTML } from './status-badge.js';

export function renderMetricCard({
  label,
  value = 0,
  description = '',
  status = '',
  tone = 'default',
  icon = ''
} = {}) {
  const card = document.createElement('article');
  card.className = `developer-metric-card developer-metric-card--${tone}`;

  card.innerHTML = `
    <div class="developer-metric-card__header">
      <span class="developer-metric-card__label">
        ${escapeHTML(icon)}${escapeHTML(label || 'Metric')}
      </span>
      ${status ? statusBadgeHTML(status, {tone}) : ''}
    </div>
    <strong class="developer-metric-card__value">
      ${escapeHTML(value)}
    </strong>
    ${description
      ? `<span class="developer-metric-card__description">${escapeHTML(description)}</span>`
      : ''}
  `;

  return card;
}

function escapeHTML(value) {
  return String(value)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}
