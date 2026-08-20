/**
 * Fondasi bersama seluruh dashboard A-Radius.
 * Implementasi autentikasi final dapat menggantikan pembacaan session ini
 * tanpa mengubah kontrak setiap dashboard.
 */

export const DASHBOARD_ROLES = Object.freeze({
  developer: 'developer',
  administrator: 'administrator',
  teknisi: 'teknisi',
  reseller: 'reseller',
  pelanggan: 'pelanggan',
});

export const DASHBOARD_PATHS = Object.freeze({
  developer: '/dashboard/developer/',
  administrator: '/dashboard/administrator/',
  teknisi: '/dashboard/teknisi/',
  reseller: '/dashboard/reseller/',
  pelanggan: '/dashboard/pelanggan/',
});

export function requireRole(expectedRole, session = window.ARADIUS_SESSION) {
  if (!session || session.role !== expectedRole) {
    const target = DASHBOARD_PATHS[session?.role] || '/';
    window.location.assign(target);
    throw new Error(`Akses ditolak untuk role: ${expectedRole}`);
  }
  return session;
}

export async function apiFetch(path, options = {}) {
  const response = await fetch(`/api/v1${path}`, {
    credentials: 'include',
    headers: { Accept: 'application/json', 'Content-Type': 'application/json', ...(options.headers || {}) },
    ...options,
  });
  if (!response.ok) {
    const body = await response.json().catch(() => ({}));
    const error = new Error(body?.error?.message || `API error ${response.status}`);
    error.status = response.status;
    error.code = body?.error?.code;
    throw error;
  }
  return response.status === 204 ? null : response.json();
}

export function renderDashboardShell({ role, title, content }) {
  document.title = `A-Radius | ${title}`;

  const root = document.querySelector('[data-dashboard-root]');
  if (!root) {
    throw new Error('Elemen [data-dashboard-root] tidak ditemukan');
  }

  root.replaceChildren();

  const shell = document.createElement('main');
  shell.className = 'dashboard-shell';
  shell.dataset.role = role;

  const header = document.createElement('header');

  const brand = document.createElement('strong');
  brand.textContent = 'A-Radius';

  const heading = document.createElement('span');
  heading.textContent = title;

  header.append(brand, heading);

  const section = document.createElement('section');

  if (typeof content === 'string') {
    section.innerHTML = content;
  } else if (content instanceof Node) {
    section.appendChild(content);
  } else if (content instanceof DocumentFragment) {
    section.appendChild(content);
  } else if (content != null) {
    throw new TypeError('Dashboard content harus berupa string, Node, atau DocumentFragment');
  }

  shell.append(header, section);
  root.appendChild(shell);

  return {
    root,
    shell,
    header,
    section,
  };
}
