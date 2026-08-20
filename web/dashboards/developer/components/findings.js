function escapeHTML(value) {
  return String(value ?? '')
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;');
}

export function renderFindings(container, findings = []) {
  if (!container) return;

  if (!Array.isArray(findings) || findings.length === 0) {
    container.innerHTML = `
      <section class="developer-card developer-findings">
        <div class="developer-card__header">
          <span>Security Findings</span>
          <span>0</span>
        </div>
        <div class="developer-state-view developer-state-view--empty">
          <strong>No findings</strong>
          <p>No security findings available.</p>
        </div>
      </section>
    `;
    return;
  }

  container.innerHTML = `
    <section class="developer-card developer-findings">
      <div class="developer-card__header">
        <span>Security Findings</span>
        <span>${findings.length}</span>
      </div>

      <div class="developer-findings__list">
        ${findings.map((item) => `
          <article class="developer-finding">
            <strong>${escapeHTML(item.title || 'Finding')}</strong>
            <span>${escapeHTML(item.severity || item.level || 'UNKNOWN')}</span>
          </article>
        `).join('')}
      </div>
    </section>
  `;
}
