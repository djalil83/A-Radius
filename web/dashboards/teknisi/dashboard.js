import { requireRole, renderDashboardShell } from '../../shared/dashboard-shell.js';

const session = requireRole('teknisi');
renderDashboardShell({
  role: 'teknisi',
  title: 'Teknisi Dashboard',
  content: `<h1>Teknisi Dashboard</h1><p>Penanganan tiket, status jaringan, dan pekerjaan lapangan.</p><div data-dashboard-widgets></div>`,
});

// TODO(teknisi): tambahkan widget dan pemanggilan endpoint khusus dashboard ini.
console.info('Dashboard teknisi siap dikembangkan', session);
