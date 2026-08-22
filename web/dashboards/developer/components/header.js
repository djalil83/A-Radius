export function renderDeveloperHeader({
  title = 'A-RADIUS DEVELOPER',
  section = 'Security Center',
  server = '',
  status = 'LIVE'
} = {}) {
  const header = document.createElement('header');
  header.className = 'developer-header';

  header.innerHTML = `
    <div class="developer-header__identity">
      <strong>${escapeHTML(title)}</strong>
      ${server ? `<span>Server: ${escapeHTML(server)}</span>` : ''}
    </div>

    <div class="developer-header__section">
      <span class="developer-header__icon" aria-hidden="true">🛡️</span>
      <strong>${escapeHTML(section)}</strong>
    </div>

    <div class="developer-header__status">
      <span>${escapeHTML(status)}</span>
    </div>
  `;

  return header;
}

function escapeHTML(value) {
  return String(value)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}
