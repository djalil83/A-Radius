import { requireRole, renderDashboardShell } from '../../shared/dashboard-shell.js';

const session = requireRole('pelanggan');
renderDashboardShell({
  role: 'pelanggan',
  title: 'Pelanggan Dashboard',
  content: `<h1>Pelanggan Dashboard</h1><p>Melihat paket, status layanan, tagihan, dan membuat tiket bantuan.</p><div data-dashboard-widgets></div>`,
});

// TODO(pelanggan): tambahkan widget dan pemanggilan endpoint khusus dashboard ini.
console.info('Dashboard pelanggan siap dikembangkan', session);
