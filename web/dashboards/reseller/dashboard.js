import { requireRole, renderDashboardShell } from '../../shared/dashboard-shell.js';

const session = requireRole('reseller');
renderDashboardShell({
  role: 'reseller',
  title: 'Reseller Dashboard',
  content: `<h1>Reseller Dashboard</h1><p>Pengelolaan pelanggan reseller, paket, dan aktivasi layanan.</p><div data-dashboard-widgets></div>`,
});

// TODO(reseller): tambahkan widget dan pemanggilan endpoint khusus dashboard ini.
console.info('Dashboard reseller siap dikembangkan', session);
