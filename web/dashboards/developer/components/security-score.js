(function () {
  'use strict';

  window.DeveloperSecurityScore = {
    render(container, score) {
      if (!container) return;

      const value = Number.isFinite(Number(score))
        ? Number(score)
        : 0;

      container.innerHTML = `
        <section class="developer-card developer-security-score"
                 aria-label="Security score">
          <div class="developer-card__header">
            <span>Security Score</span>
          </div>
          <strong class="developer-security-score__value">
            ${value}
          </strong>
        </section>
      `;
    }
  };
})();
