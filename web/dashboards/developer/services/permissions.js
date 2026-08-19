(function () {
  'use strict';

  const permissions = new Set();

  function load(list) {
    permissions.clear();
    for (const permission of list || []) {
      permissions.add(permission);
    }
  }

  function has(permission) {
    return permissions.has(permission);
  }

  function requirePermission(permission) {
    if (!has(permission)) {
      const error = new Error(`Missing permission: ${permission}`);
      error.code = 'FORBIDDEN';
      throw error;
    }
    return true;
  }

  window.DeveloperPermissions = Object.freeze({
    load,
    has,
    require: requirePermission
  });
})();
