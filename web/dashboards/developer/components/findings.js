(function () {
  'use strict';

  window.DeveloperFindings = {
    render(container, findings = []) {
      if (!container) return;

      if (!findings.length) {
        container.innerHTML =
          '<section class="developer-card">No security findings.</section>';
        return;
      }

      container.innerHTML = `
        <section class="developer-card">
          <div class="developer-card__header">
            <span>Security Findings</span>
            <span>${findings.length}</span>
          </div>
          <div class="developer-findings">
            ${findings.map(item => `
              <article class="developer-finding">
                <strong>${escapeHTML(item.title || 'Finding')}</strong>
                <span>${escapeHTML(item.severity || 'unknown')}</span>
              </article>
            `).join('')}
          </div>
        </section>
      `;
    }
  };

  function escapeHTML(value) {
    return String(value)
      .replaceAll('&', '&amp;')
      .replaceAll('<', '&lt;')
      .replaceAll('>', '&gt;')
      .replaceAll('"', '&quot;')
      .replaceAll("'", '&#039;');
  }
})();
