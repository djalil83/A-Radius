import { requireRole, renderDashboardShell, apiFetch } from '../../shared/dashboard-shell.js';

requireRole('developer');

const featuredFinding = {
  id: 'SEC-2026-00127', severity: 'HIGH', module: 'Langganan API', endpoint: '/api/langganan/update', problem: 'Authorization check tidak konsisten.', recommendation: 'Tambahkan role-based authorization middleware.', roles: ['Administrator', 'Technician', 'Sales', 'Reseller'], before: 'authorization → partial', after: 'authorization → RBAC middleware', production: 'UNCHANGED',
};

const findings = [
  { level: 'HIGH', tone: 'high', title: 'API endpoint memiliki authorization policy tidak konsisten', module: '/api/langganan', detail: 'Policy authorization berbeda antara endpoint baca dan tulis. Verifikasi kembali permission profile:read dan profile:write.' },
  { level: 'MEDIUM', tone: 'medium', title: 'Dependency membutuhkan pembaruan keamanan', module: 'frontend', detail: 'Terdapat dependency frontend yang perlu dibandingkan dengan advisory keamanan terbaru.' },
  { level: 'LOW', tone: 'low', title: 'Session timeout belum optimal', module: 'authentication', detail: 'Durasi session aktif lebih panjang dari baseline keamanan yang direkomendasikan.' },
];

const scanCategories = [
  ['SQL Injection', 'Code, API, Database'], ['XSS', 'Code, API'], ['CSRF', 'API, Authentication'], ['Broken Access Control', 'Authorization'], ['Authentication Weakness', 'Authentication, Session'], ['Insecure Secrets', 'Credential, Configuration'], ['Dependency Vulnerability', 'Dependency'], ['API Exposure', 'API'], ['Dangerous Configuration', 'Configuration, Infrastructure'], ['Privilege Escalation', 'Authorization, Access'],
];
const securityLayers = {
  prevention: [['RBAC', 'Active'], ['MFA', 'Required'], ['Session Security', 'Monitored'], ['Rate Limiting', 'Active'], ['Input Validation', 'Enforced'], ['Secret Management', 'Vault-backed'], ['API Authorization', 'Enforced']],
  detection: [['AI Security Checker', 'Ready'], ['Anomaly Detection', 'Monitoring'], ['Login Monitoring', 'Monitoring'], ['API Monitoring', 'Monitoring'], ['File Changes', 'Monitoring'], ['Database Changes', 'Monitoring'], ['Permission Changes', 'Audited']],
  response: [['Alert', 'Proposed'], ['Quarantine / Lock', 'Approval required'], ['Rollback', 'Approval required'], ['Disable Credential', 'Approval required'], ['Incident Report', 'Available'], ['Audit Trail', 'Append-only']],
};
const securityLayerMarkup = (key, title, icon) => `<article class="layer-card content-card"><div class="section-heading"><h1>${icon} ${title}</h1><span>${key.toUpperCase()}</span></div><div class="layer-list">${securityLayers[key].map(([name, status]) => `<div class="layer-row"><span>${name}</span><strong>${status}</strong></div>`).join('')}</div></article>`;
const continuousKnowledge = [
  ['CVE / CISA KEV', 'Prioritas eksploitasi aktif'], ['Security Advisories / CSAF', 'Advisory terstruktur'], ['OWASP API Security', 'Taxonomy API risk'], ['Dependency Vulnerability', 'CVE dan package evidence'], ['Framework / Vendor Update', 'Perubahan keamanan upstream'],
];
const continuousIntelligenceMarkup = () => `<section class="continuous-intelligence content-card"><div class="section-heading"><div><h1>🧠 CONTINUOUS SECURITY INTELLIGENCE</h1><span>Security knowledge → validated analysis → developer report</span></div><span class="analysis-badge">AI ADVISORY ONLY</span></div><div class="intelligence-flow"><div class="knowledge-panel"><p class="eyebrow">INTERNET / SECURITY KNOWLEDGE</p>${continuousKnowledge.map(([name, purpose]) => `<div class="knowledge-row"><strong>${name}</strong><span>${purpose}</span></div>`).join('')}</div><div class="intelligence-arrow">→</div><div class="engine-panel"><strong>AI SECURITY<br>LEARNING ENGINE</strong><span>FILTER &amp; VALIDATION</span><span>PROVENANCE REQUIRED</span></div><div class="intelligence-arrow">→</div><div class="analysis-panel"><p class="eyebrow">A-RADIUS APPLICATION</p><div class="analysis-targets"><span>CODE</span><span>API</span><span>DB</span></div><strong>SECURITY ANALYSIS</strong><span>NEW FINDING / INSIGHT</span><span>DEVELOPER REPORT</span></div></div><div class="intelligence-policy"><span>🔒 AI tidak memiliki akses langsung ke Production.</span><strong>Patch preview, automated test, staging test, dan human approval wajib sebelum deploy.</strong></div><div class="intelligence-actions"><button type="button" class="secondary-action" data-intelligence="sources">VIEW SOURCES</button><button type="button" class="primary-action" data-intelligence="finding">LOAD FEATURED FINDING</button><span data-intelligence-feedback aria-live="polite"></span></div></section>`;

const newSecurityIntelligenceMarkup = () => `<section class="new-intelligence-report content-card"><div class="section-heading"><div><h1>🧠 NEW SECURITY INTELLIGENCE</h1><span>Security Intelligence Feed · 17/08/2026 14:02 WITA</span></div><span class="featured-severity">HIGH</span></div><div class="intelligence-report-grid"><div><span>TANGGAL</span><strong>17/08/2026 14:02 WITA</strong></div><div><span>SOURCE</span><strong>Security Intelligence Feed</strong></div><div><span>KATEGORI</span><strong>API Authorization</strong></div><div><span>SEVERITY</span><strong class="high-text">HIGH</strong></div></div><div class="new-report-body"><h3>AI FINDING</h3><p>Ditemukan pola vulnerability baru yang relevan dengan modul API.</p><h3>A-RADIUS IMPACT</h3><div class="impact-list"><span>✓ API Langganan</span><span>✓ API Administrator</span><span>✓ API Technician</span></div><div class="report-status"><span>STATUS</span><strong>⚠ Belum diterapkan</strong></div><h3>AI RECOMMENDATION</h3><p>Periksa authorization middleware dan konsistensi RBAC.</p><div class="report-meta"><span>PATCH: Preview tersedia.</span><span>TEST: Security test tersedia.</span><span>PRODUCTION: <strong>UNCHANGED</strong></span></div></div><div class="intelligence-report-actions"><button type="button" class="secondary-action" data-report-action="analysis">VIEW ANALYSIS</button><button type="button" class="primary-action" data-report-action="patch">VIEW PATCH</button><button type="button" class="secondary-action" data-report-action="test">RUN TEST</button><button type="button" class="primary-action" data-report-action="approval">REQUEST DEVELOPER APPROVAL</button></div><p class="report-feedback" data-report-feedback aria-live="polite"></p></section>`;

const knowledgeVersionMarkup = () => `<section class="knowledge-version-panel content-card"><div class="section-heading"><div><h1>🧠 SECURITY KNOWLEDGE</h1><span>A-RADIUS Developer Security Intelligence · Knowledge DB isolated</span></div><span class="analysis-badge">PRODUCTION SEPARATED</span></div><div class="active-knowledge-banner"><div><p class="eyebrow">ACTIVE VERSION</p><strong>SK-2.4.7</strong><span>Security Knowledge Engine</span></div><div><b class="active-status">● ACTIVE</b><small>Last Update: 17/08/2026 14:20 WITA</small></div></div><div class="knowledge-status-counters"><div><strong>12</strong><span>NEW</span></div><div><strong>4</strong><span>REVIEW</span></div><div><strong>1</strong><span>ACTIVE</span></div><div><strong>18</strong><span>ARCHIVED</span></div></div><div class="knowledge-table-wrap"><div class="knowledge-table-heading"><strong>KNOWLEDGE VERSIONS</strong><span>Human-controlled promotion</span></div><div class="knowledge-table"><div class="knowledge-table-row table-head"><span>VERSION</span><span>STATUS</span><span>FINDINGS</span><span>SOURCE</span><span>DATE</span><span>ACTION</span></div><div class="knowledge-table-row"><strong>SK-2.4.7</strong><span class="status-active">● ACTIVE</span><span>23</span><span>Security AI</span><span>17/08/26</span><span><button type="button" data-version-action="view">VIEW</button> <button type="button" data-version-action="compare">COMPARE</button></span></div><div class="knowledge-table-row"><strong>SK-2.4.6</strong><span>ARCHIVED</span><span>19</span><span>Security AI</span><span>14/08/26</span><span><button type="button" data-version-action="view">VIEW</button> <button type="button" data-version-action="rollback">ROLLBACK</button></span></div><div class="knowledge-table-row"><strong>SK-2.4.5</strong><span>ARCHIVED</span><span>17</span><span>Security AI</span><span>11/08/26</span><span><button type="button" data-version-action="view">VIEW</button> <button type="button" data-version-action="compare">COMPARE</button></span></div><div class="knowledge-table-row"><strong>SK-2.4.4</strong><span>ARCHIVED</span><span>16</span><span>Security AI</span><span>08/08/26</span><span><button type="button" data-version-action="view">VIEW</button> <button type="button" data-version-action="rollback">ROLLBACK</button></span></div></div></div><div class="knowledge-boundary-flow"><span>INTERNET</span><i>↓</i><span>AI KNOWLEDGE DB</span><i>↓</i><span>ANALYSIS CACHE</span><i>↓</i><span>RULE ENGINE</span><i>↓</i><strong>A-RADIUS SCAN</strong></div><div class="knowledge-policy-note"><strong>AI LEARNING ISOLATED</strong><span>AI tidak boleh langsung mengubah ACTIVE. APPROVED → STAGED → ACTIVE memerlukan kontrol Developer.</span><b>PRODUCTION: UNCHANGED</b></div><div class="knowledge-version-actions"><button type="button" class="secondary-action" data-knowledge-action="policy">VIEW POLICY</button><button type="button" class="primary-action" data-knowledge-action="featured">VIEW NEW KNOWLEDGE</button><span data-knowledge-feedback aria-live="polite"></span></div></section>`;

const knowledgeCompareMarkup = () => `<section class="knowledge-compare-panel content-card"><div class="section-heading"><div><h1>🔬 COMPARE VERSION</h1><span>Knowledge change ≠ Production application change</span></div><span class="analysis-badge">PRODUCTION: UNCHANGED</span></div><div class="compare-version-title"><strong>SK-2.4.7</strong><span>VS</span><strong>SK-2.4.8</strong></div><div class="compare-columns"><div><h3>ADDED</h3><span>+ API-AUTH-042</span><span>+ SESSION-019</span><span>+ DEP-031</span></div><div><h3>UPDATED</h3><span>~ API-AUTH-018</span><span>~ DATABASE-007</span></div><div><h3>REMOVED</h3><span>- API-AUTH-003</span></div></div><div class="affected-modules"><h3>AFFECTED MODULE</h3><div><span>Administrator <b class="high-text">HIGH</b></span><span>Langganan <b class="high-text">HIGH</b></span><span>Pelanggan <b class="medium-text">MEDIUM</b></span><span>Technician <b class="medium-text">MEDIUM</b></span><span>Sales <b class="low-text">LOW</b></span><span>Customer <b class="low-text">LOW</b></span></div></div><div class="knowledge-patch-pipeline"><span>SECURITY KNOWLEDGE</span><i>↓</i><span>AI ANALYSIS</span><i>↓</i><span>FINDING</span><i>↓</i><span>RECOMMENDATION</span><i>↓</i><strong>PATCH PROPOSAL</strong><i>↓</i><span>DEVELOPER PREVIEW</span><i>↓</i><span>APPROVAL</span><i>↓</i><span>STAGING</span><i>↓</i><span>PRODUCTION</span></div><p class="compare-feedback" data-compare-feedback aria-live="polite"></p></section>`;

const newIntelligenceDetailMarkup = () => `<section class="new-intelligence-detail content-card"><div class="section-heading"><div><h1>🧠 NEW SECURITY INTELLIGENCE</h1><span>ID: INT-2026-00821 · Candidate SK-2.4.8</span></div><span class="featured-severity">HIGH</span></div><div class="intelligence-report-grid"><div><span>CATEGORY</span><strong>API Security</strong></div><div><span>RISK</span><strong class="high-text">HIGH</strong></div><div><span>CONFIDENCE</span><strong>91%</strong></div><div><span>STATUS</span><strong class="review-status">⚠ REVIEW REQUIRED</strong></div></div><div class="new-report-body"><h3>AI ANALYSIS</h3><p>Ditemukan security advisory baru yang relevan dengan modul API.</p><div class="impact-split"><div><h3>RELEVAN DENGAN</h3><span>✓ /api/langganan</span><span>✓ /api/pelanggan</span><span>✓ /api/admin</span><span>✓ /api/technician</span></div><div><h3>TIDAK RELEVAN</h3><span>○ UI statis</span><span>○ Template voucher</span></div></div><h3>PROPOSED RULE</h3><p><strong>API-AUTH-042</strong><br>Pastikan setiap endpoint sensitif memiliki authorization policy berdasarkan role dan permission.</p><div class="report-status"><span>APPLICATION IMPACT</span><strong>HIGH</strong></div></div><div class="intelligence-report-actions"><button type="button" class="secondary-action" data-new-intel-action="evidence">VIEW EVIDENCE</button><button type="button" class="secondary-action" data-new-intel-action="test">RUN TEST</button><button type="button" class="secondary-action" data-new-intel-action="compare">COMPARE</button><button type="button" class="primary-action" data-new-intel-action="approval">REQUEST APPROVAL</button></div><p class="report-feedback" data-new-intel-feedback aria-live="polite"></p></section>`;

const riskLevels = [
  ['CRITICAL', 'Potensi kompromi sistem/data sangat serius', 'critical'], ['HIGH', 'Risiko besar terhadap akun, API atau data', 'high'], ['MEDIUM', 'Kelemahan yang perlu diperbaiki', 'medium'], ['LOW', 'Hardening / konfigurasi', 'low'], ['INFO', 'Rekomendasi keamanan', 'info'],
];

const groups = [
  { icon: '🧠', title: 'AI Center', items: ['AI Dashboard', 'AI Knowledge', 'Learning History', 'New Security Intelligence', 'Application Analysis', 'AI Recommendations'] },
  { icon: '🛡️', title: 'Security Center', items: ['Security Score', 'AI Security Checker', 'Threat Detection', 'Vulnerability Scanner', 'Dependency Scanner', 'API Security', 'Database Security', 'Security Audit'] },
  { icon: '📚', title: 'Security Knowledge', items: ['New Knowledge', 'CVE Intelligence', 'OWASP Rules', 'Dependency Intelligence', 'Framework Updates', 'Infrastructure Security', 'Knowledge Version'] },
  { icon: '🔬', title: 'AI Research', items: ['Analyze New Threat', 'Compare With A-RADIUS', 'Impact Analysis', 'Generate Recommendation', 'Generate Patch Preview'] },
  { icon: '👁️', title: 'Preview', items: ['Code Preview', 'Security Fix Preview', 'UI Preview', 'Database Migration Preview', 'Configuration Preview'] },
  { icon: '✅', title: 'Approval', items: ['Security Fix', 'Code Update', 'Dependency Update', 'Database Change', 'Production Deployment'] },
  { icon: '🧾', title: 'Audit Trail', items: ['AI Learning', 'AI Recommendation', 'Developer Approval', 'Deployment', 'Rollback'] },
];

const navMarkup = () => groups.map((group) => `<section class="nav-group"><h2>${group.icon} ${group.title}</h2><div class="nav-items">${group.items.map((item) => `<button type="button" class="nav-item" data-feature="${item}">${item}</button>`).join('')}</div></section>`).join('');
const findingMarkup = () => findings.map((finding, index) => `<article class="finding-row ${finding.tone}" data-finding="${index}"><div class="finding-level"><span class="severity-dot"></span><strong>${finding.level}</strong></div><div class="finding-copy"><strong>${finding.title}</strong><span>Module : ${finding.module}</span><p class="finding-detail" hidden>${finding.detail}</p></div><button type="button" class="finding-detail-button" data-detail="${index}">DETAIL</button></article>`).join('');
const riskMarkup = () => riskLevels.map(([level, description, tone]) => `<div class="risk-item ${tone}"><span class="risk-symbol">●</span><strong>${level}</strong><span>${description}</span></div>`).join('');
const featuredFindingMarkup = () => `<article class="featured-finding"><div class="finding-banner"><div><span class="finding-kicker">SECURITY FINDING</span><h2>${featuredFinding.id}</h2></div><span class="featured-severity">🔴 ${featuredFinding.severity}</span></div><div class="featured-grid"><div><span>ID</span><strong>${featuredFinding.id}</strong></div><div><span>MODULE</span><strong>${featuredFinding.module}</strong></div><div><span>ENDPOINT</span><strong>${featuredFinding.endpoint}</strong></div></div><div class="featured-copy"><h3>Problem</h3><p>${featuredFinding.problem}</p><h3>Recommendation</h3><p>${featuredFinding.recommendation}</p><h3>Affected Roles</h3><div class="role-tags">${featuredFinding.roles.map((role) => `<span>${role}</span>`).join('')}</div><div class="production-unchanged"><span>Production</span><strong>${featuredFinding.production}</strong></div></div><button type="button" class="primary-action create-preview-button" data-finding-action="preview">CREATE FIX PREVIEW</button></article>`;
const previewMarkup = () => `<section class="preview-card content-card" data-preview-card hidden><div class="section-heading"><div><h1>DEVELOPER PREVIEW</h1><span>Proposed fix only — production remains unchanged</span></div><span class="preview-status">DRAFT</span></div><div class="before-after"><div><span class="diff-label before-label">BEFORE</span><code>${featuredFinding.before}</code></div><div class="diff-arrow">→</div><div><span class="diff-label after-label">AFTER</span><code>${featuredFinding.after}</code></div></div><div class="preview-actions"><button type="button" class="primary-action" data-finding-action="security-test">RUN SECURITY TEST</button><button type="button" class="secondary-action" data-finding-action="request-approval">REQUEST APPROVAL</button></div><p class="preview-feedback" data-preview-feedback aria-live="polite"></p></section>`;
const categoryMarkup = () => scanCategories.map(([name, scope]) => `<div class="category-item"><strong>${name}</strong><span>${scope}</span></div>`).join('');
const workflowMarkup = () => `<div class="workflow-track"><div class="workflow-node">DEVELOPER</div><i>↓</i><div class="workflow-node">AI SECURITY CHECKER</div><i>↓</i><div class="workflow-node">FINDING</div><i>↓</i><div class="workflow-node">AI RECOMMENDATION</div><i>↓</i><div class="workflow-node">GENERATE PROPOSED FIX</div><i>↓</i><div class="workflow-node">DEVELOPER PREVIEW</div><div class="workflow-branches"><button type="button" class="danger-action" data-workflow="reject">✕ REJECT</button><button type="button" class="primary-action" data-workflow="approve">✓ APPROVE</button></div><i>↓</i><div class="workflow-node gated">SECURITY TEST</div><i>↓</i><div class="workflow-node gated">STAGING TEST</div><i>↓</i><div class="workflow-node gated">PRODUCTION APPROVAL</div><i>↓</i><div class="workflow-node gated">DEPLOY</div><i>↓</i><div class="workflow-node gated">HEALTH CHECK</div><div class="workflow-outcomes"><span class="success-outcome">SUCCESS → AUDIT</span><span class="failure-outcome">FAILURE → ROLLBACK</span></div></div>`;

renderDashboardShell({
  role: 'developer',
  title: 'Security Center',
  content: `<div class="developer-layout"><aside class="developer-sidebar"><div class="brand-lockup"><span class="brand-mark">A</span><div><strong>A-RADIUS</strong><small>DEVELOPER</small></div></div><nav aria-label="Developer navigation">${navMarkup()}</nav></aside><main class="developer-main"><header class="security-header"><div><strong>A-RADIUS DEVELOPER</strong><span>Server: D74111</span></div><div class="security-header-title"><span class="shield">🛡️</span><strong>SECURITY CENTER</strong></div><div class="security-header-meta"><span>17/08/2026 13:50 WITA</span><span class="live-dot">● LIVE</span></div></header><section class="security-intro"><div class="security-score-block"><div><p class="eyebrow">SECURITY SCORE</p><strong class="score-number">87<span> / 100</span></strong><span class="score-label">GOOD</span></div><div class="score-meter"><span></span></div></div><div class="threat-block"><p class="eyebrow">THREATS</p><strong class="threat-total">3</strong><div class="threat-line high-text">⚠ HIGH <b>1</b></div><div class="threat-line medium-text">⚠ MEDIUM <b>2</b></div></div></section><section class="scan-toolbar" aria-label="Security scan actions"><button class="scan-button primary-action" type="button" data-scan="full">🔍 FULL SCAN</button><button class="scan-button" type="button" data-scan="code">CODE</button><button class="scan-button" type="button" data-scan="api">API</button><button class="scan-button" type="button" data-scan="database">DATABASE</button><button class="scan-button" type="button" data-scan="configuration">CONFIG</button><span class="scan-status" data-scan-status>Ready</span></section><section class="layer-overview"><div class="layer-intro"><p class="eyebrow">CROSS-APPLICATION SECURITY LAYER</p><h2>Protection for all A-Radius applications</h2><p>Developer Security Center memusatkan prevention, detection, dan response untuk API, dashboard, database, authentication, serta layanan pelanggan.</p></div>${securityLayerMarkup('prevention', 'PREVENTION', '🛡️')}${securityLayerMarkup('detection', 'DETECTION', '🔎')}${securityLayerMarkup('response', 'RESPONSE', '🚨')}</section>${continuousIntelligenceMarkup()}${knowledgeVersionMarkup()}${knowledgeCompareMarkup()}${newIntelligenceDetailMarkup()}${newSecurityIntelligenceMarkup()}<section class="security-findings content-card"><div class="section-heading"><h1>🚨 SECURITY FINDINGS</h1><span>3 findings</span></div>${findingMarkup()}</section>${featuredFindingMarkup()}${previewMarkup()}<section class="security-catalog"><article class="content-card"><div class="section-heading"><h1>🔎 SCAN CATEGORIES</h1><span>10 controls</span></div><div class="category-grid">${categoryMarkup()}</div></article><article class="content-card"><div class="section-heading"><h1>🚦 RISK LEVEL</h1><span>Severity policy</span></div><div class="risk-list">${riskMarkup()}</div></article></section><section class="ai-analysis content-card"><div class="section-heading"><h1>AI ANALYSIS</h1><span class="analysis-badge">ANALYZED</span></div><div class="analysis-steps"><span class="step-done">✓ Masalah terdeteksi</span><span class="step-done">✓ Dampak dianalisis</span><span class="step-done">✓ Rekomendasi dibuat</span><span class="step-pending">✕ Production BELUM diubah</span></div><div class="analysis-actions"><button type="button" class="secondary-action" data-analysis="detail">LIHAT DETAIL</button><button type="button" class="primary-action" data-analysis="preview">BUAT PREVIEW FIX</button><button type="button" class="danger-action" data-analysis="ignore">IGNORE</button></div><p class="action-feedback" data-action-feedback aria-live="polite"></p></section><section class="workflow-card content-card"><div class="section-heading"><div><h1>🔐 AI-TO-PRODUCTION GATE</h1><span class="no-production-access">AI TIDAK MEMILIKI AKSES LANGSUNG KE PRODUCTION</span></div></div>${workflowMarkup()}<p class="workflow-feedback" data-workflow-feedback aria-live="polite"></p></section></main></div>`, 
});

document.querySelectorAll('[data-feature]').forEach((button) => button.addEventListener('click', () => {
  document.querySelectorAll('[data-feature]').forEach((item) => item.classList.remove('active'));
  button.classList.add('active');
}));

document.querySelectorAll('[data-detail]').forEach((button) => button.addEventListener('click', () => {
  const detail = document.querySelector(`[data-finding="${button.dataset.detail}"] .finding-detail`);
  const open = !detail.hidden;
  detail.hidden = open;
  button.textContent = open ? 'DETAIL' : 'TUTUP';
}));

document.querySelectorAll('[data-scan]').forEach((button) => button.addEventListener('click', async () => {
  const status = document.querySelector('[data-scan-status]');
  const type = button.dataset.scan;
  status.textContent = `${type.toUpperCase()} scan diproses...`;
  document.querySelectorAll('[data-scan]').forEach((item) => { item.disabled = true; });
  try {
    await apiFetch('/developer/security/scans', { method: 'POST', body: JSON.stringify({ type: type === 'configuration' ? 'config' : type }) });
    status.textContent = `${type.toUpperCase()} scan dijadwalkan`;
  } catch (error) {
    status.textContent = error.status === 404 ? 'Backend scan belum aktif' : 'Scan gagal';
  } finally {
    setTimeout(() => { document.querySelectorAll('[data-scan]').forEach((item) => { item.disabled = false; }); }, 1400);
  }
}));

document.querySelectorAll('[data-workflow]').forEach((button) => button.addEventListener('click', () => {
  const feedback = document.querySelector('[data-workflow-feedback]');
  feedback.textContent = button.dataset.workflow === 'approve' ? 'Preview disetujui. Security Test, Staging Test, dan Production Approval tetap wajib dilewati.' : 'Proposed fix ditolak. Production tidak disentuh dan perubahan dikembalikan ke draft.';
}));

document.querySelectorAll('[data-finding-action]').forEach((button) => button.addEventListener('click', async () => {
  const action = button.dataset.findingAction;
  const preview = document.querySelector('[data-preview-card]');
  const feedback = document.querySelector('[data-preview-feedback]');
  if (action === 'preview') { preview.hidden = false; preview.scrollIntoView({ behavior: 'smooth', block: 'nearest' }); feedback.textContent = 'Fix preview dibuat. Tidak ada perubahan ke Production.'; return; }
  feedback.textContent = action === 'security-test' ? 'Security Test dijadwalkan untuk SEC-2026-00127.' : 'Approval request dibuat. Developer approval diperlukan sebelum tahap Production.';
  try { await apiFetch(`/developer/security/findings/${featuredFinding.id}/${action}`, { method: 'POST', body: JSON.stringify({ finding_id: featuredFinding.id }) }); } catch (_) { /* UI tetap memberi feedback; backend dapat diaktifkan bertahap. */ }
}));

document.querySelectorAll('[data-layer-action]').forEach((button) => button.addEventListener('click', () => {
  const feedback = document.querySelector('[data-action-feedback]');
  feedback.textContent = 'AI hanya membuat rekomendasi/proposal. Tindakan ini memerlukan approval Administrator atau Developer yang berwenang.';
}));

document.querySelectorAll('[data-intelligence]').forEach((button) => button.addEventListener('click', async () => {
  const feedback = document.querySelector('[data-intelligence-feedback]');
  const action = button.dataset.intelligence;
  feedback.textContent = action === 'sources' ? 'Trusted sources: CISA KEV, OWASP API Security, CSAF, dependency advisories, dan vendor advisories.' : 'Featured finding SEC-2026-00127 siap dianalisis. Production tetap UNCHANGED.';
  try { await apiFetch(`/developer/security/continuous/${action === 'sources' ? 'sources' : 'featured-finding'}`); } catch (_) { /* panel tetap informasional saat endpoint belum aktif */ }
}));

document.querySelectorAll('[data-version-action]').forEach((button) => button.addEventListener('click', async () => {
  const action = button.dataset.versionAction;
  const feedback = document.querySelector('[data-knowledge-feedback]');
  feedback.textContent = action === 'rollback' ? 'Rollback hanya membuat proposal dan memerlukan approval Developer. ACTIVE tidak diubah otomatis.' : `${action.toUpperCase()} knowledge version dipilih. Production tetap UNCHANGED.`;
  if (action === 'compare') { const compare = document.querySelector('[data-compare-feedback]'); compare.textContent = 'Compare SK-2.4.7 vs SK-2.4.8 dimuat. Knowledge berubah; aplikasi Production belum berubah.'; }
  try { await apiFetch(action === 'compare' ? '/developer/security/knowledge/compare' : '/developer/security/knowledge/versions'); } catch (_) { /* advisory-only */ }
}));

document.querySelectorAll('[data-new-intel-action]').forEach((button) => button.addEventListener('click', async () => {
  const feedback = document.querySelector('[data-new-intel-feedback]');
  const messages = { evidence: 'Evidence INT-2026-00821 ditampilkan dari knowledge source tervalidasi.', test: 'Security test dijadwalkan di sandbox/staging.', compare: 'Candidate SK-2.4.8 dibandingkan dengan active SK-2.4.7.', approval: 'Request approval dibuat. Status tetap REVIEW REQUIRED sampai Developer menyetujui.' };
  feedback.textContent = messages[button.dataset.newIntelAction];
  try { await apiFetch('/developer/security/knowledge/new-intelligence'); } catch (_) { /* status UI tetap advisory-only */ }
}));

document.querySelectorAll('[data-knowledge-action]').forEach((button) => button.addEventListener('click', async () => {
  const action = button.dataset.knowledgeAction;
  const feedback = document.querySelector('[data-knowledge-feedback]');
  feedback.textContent = action === 'policy' ? 'Policy: AI Knowledge DB terisolasi, analysis cache read-only, auto promotion ke Production disabled.' : 'New knowledge v1.3 dimuat. Status NEW dan Production tetap UNCHANGED.';
  try { await apiFetch(`/developer/security/knowledge/${action}`); } catch (_) { /* panel tetap informasional saat endpoint belum aktif */ }
}));

document.querySelectorAll('[data-report-action]').forEach((button) => button.addEventListener('click', async () => {
  const action = button.dataset.reportAction;
  const feedback = document.querySelector('[data-report-feedback]');
  const messages = { analysis: 'Analisis AI dibuka: API Authorization memiliki relevansi HIGH.', patch: 'Patch preview tersedia. Production tetap UNCHANGED.', test: 'Security test dijadwalkan di sandbox/staging.', approval: 'Developer approval request dibuat. Production belum berubah.' };
  feedback.textContent = messages[action];
  if (action === 'patch' || action === 'test' || action === 'approval') {
    try { await apiFetch(action === 'patch' ? '/developer/security/continuous/patch-preview' : `/developer/security/findings/${featuredFinding.id}/${action === 'test' ? 'security-test' : 'request-approval'}`, { method: 'POST', body: JSON.stringify({ finding_id: 'SEC-2026-00127' }) }); } catch (_) { /* feedback tetap aman saat endpoint belum tersedia */ }
  }
}));

document.querySelectorAll('[data-analysis]').forEach((button) => button.addEventListener('click', () => {
  const feedback = document.querySelector('[data-action-feedback]');
  const messages = { detail: 'Detail evidence dan rekomendasi akan dibuka dari Security Report.', preview: 'Preview fix dibuat tanpa mengubah production.', ignore: 'Finding ditandai untuk di-review sebelum diabaikan.' };
  feedback.textContent = messages[button.dataset.analysis];
}));
