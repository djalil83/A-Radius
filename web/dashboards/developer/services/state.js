(function () {
  'use strict';

  const state = {
    section: 'overview',
    securityScore: null,
    findings: [],
    threats: [],
    scans: [],
    approvals: [],
    loading: false,
    error: null
  };

  const listeners = new Set();

  function get() {
    return structuredClone(state);
  }

  function patch(values) {
    Object.assign(state, values);
    listeners.forEach(fn => fn(get()));
  }

  function subscribe(fn) {
    listeners.add(fn);
    return () => listeners.delete(fn);
  }

  window.DeveloperState = Object.freeze({
    get,
    patch,
    subscribe
  });
})();
