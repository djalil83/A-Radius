(function () {
  'use strict';

  const API_PREFIX = '/api/developer';

  async function request(path, options = {}) {
    const response = await fetch(API_PREFIX + path, {
      credentials: 'include',
      headers: {
        'Accept': 'application/json',
        ...(options.body ? {'Content-Type': 'application/json'} : {}),
        ...(options.headers || {})
      },
      ...options
    });

    let data = null;
    try {
      data = await response.json();
    } catch (_) {}

    if (!response.ok) {
      const error = new Error(
        data?.error || `Developer API request failed: ${response.status}`
      );
      error.status = response.status;
      error.data = data;
      throw error;
    }

    return data;
  }

  window.DeveloperAPI = Object.freeze({
    get: (path, options = {}) =>
      request(path, {...options, method: 'GET'}),

    post: (path, body, options = {}) =>
      request(path, {
        ...options,
        method: 'POST',
        body: JSON.stringify(body ?? {})
      })
  });
})();
