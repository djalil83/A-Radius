(function () {
  'use strict';

  const sections = Object.freeze({
    overview: 'Security Overview',
    ai: 'AI Security Center',
    security: 'Security Center',
    development: 'Development',
    approval: 'Approval Center',
    production: 'Production',
    audit: 'Audit Trail'
  });

  function navigate(section) {
    if (!Object.prototype.hasOwnProperty.call(sections, section)) {
      throw new Error(`Unknown developer dashboard section: ${section}`);
    }

    DeveloperState.patch({section});
    window.dispatchEvent(
      new CustomEvent('developer:navigate', {
        detail: {section}
      })
    );
  }

  function getSections() {
    return {...sections};
  }

  window.DeveloperDashboard = Object.freeze({
    navigate,
    getSections
  });
})();
