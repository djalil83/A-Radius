# APB AI DEVELOPMENT CENTER — BLUEPRINT

## Tujuan
Menu khusus untuk pengembangan APB menggunakan AI. AI boleh:
- menganalisis kebutuhan;
- membaca struktur/source yang diizinkan;
- mencari pengetahuan terbaru melalui koneksi internet;
- menyusun perubahan kode;
- membuat preview/diff;
- menjalankan test/sandbox;
- memberikan rekomendasi arsitektur;
- memperbarui knowledge/RAG setelah disetujui.

AI TIDAK BOLEH mengubah source production atau mengeksekusi tindakan penting
secara otomatis. Setiap tindakan yang mengubah aplikasi harus melalui
DEVELOPER APPROVAL.

## Alur wajib

Developer Command
    ↓
AI Planner
    ↓
Impact Analysis
    ↓
Proposed Changes
    ↓
Preview / Diff
    ↓
Sandbox Test
    ↓
AI Review
    ↓
Developer Approval
    ├── Reject → kembali ke Planner
    ├── Request Changes → revisi
    └── Approve
             ↓
        Execution Gate
             ↓
        Apply to Dev/Staging
             ↓
        Test / Build
             ↓
        Developer Final Approval
             ↓
        Deploy Production

## Menu AI Development

1. AI DEVELOPMENT DASHBOARD
   - Project health
   - pending approvals
   - proposed changes
   - failed tests
   - recent AI actions
   - knowledge updates
   - repository status

2. AI COMMAND CENTER
   - input perintah pengembang
   - context selector
   - module selector
   - target environment
   - risk classification
   - execute/preview mode

3. AI ANALYSIS
   - architecture analysis
   - dependency analysis
   - database impact
   - API impact
   - security impact
   - performance impact
   - regression risk

4. AI CODE PROPOSAL
   - generated patch
   - file-by-file changes
   - before/after diff
   - explanation
   - assumptions
   - affected modules

5. PREVIEW & SANDBOX
   - isolated workspace
   - dry run
   - unit tests
   - integration tests
   - lint
   - type check
   - migration validation
   - build preview
   - screenshot/UI preview

6. APPROVAL CENTER
   - pending approvals
   - approve
   - reject
   - request revision
   - approval comment
   - reviewer
   - timestamp
   - change hash

7. EXECUTION CENTER
   - development execution
   - staging execution
   - production deployment request
   - rollback
   - execution logs
   - result
   NOTE: production execution remains approval-gated.

8. AI KNOWLEDGE CENTER
   - local project knowledge
   - documentation
   - API specifications
   - MikroTik/RouterOS knowledge
   - RADIUS knowledge
   - OLT/GPON knowledge
   - FTTH knowledge
   - payment integration docs
   - troubleshooting history
   - internet research
   - source citations
   - knowledge versioning

9. AI WEB RESEARCH
   - internet-enabled research
   - official documentation priority
   - current library/version lookup
   - CVE/security advisories
   - vendor documentation
   - comparison
   - source capture
   AI may learn from research, but must not silently modify production code.

10. AI TEST LAB
    - generated test cases
    - regression suite
    - API test
    - DB test
    - network mock
    - MikroTik mock
    - RADIUS mock
    - OLT mock
    - payment webhook mock

11. AI AUDIT
    - prompt
    - context
    - recommendation
    - proposed patch
    - approval
    - execution
    - test result
    - deployment
    - rollback
    - knowledge update

## Permission model

AI permissions are deliberately separated:

AI_READ:
- read approved project files
- inspect logs
- inspect metrics
- search documentation

AI_SUGGEST:
- propose architecture
- propose code
- propose SQL/migration
- propose network configuration
- propose remediation

AI_PREVIEW:
- create isolated preview
- run tests
- generate diff
- generate build artifact

AI_EXECUTE_DEV:
- only after developer approval

AI_EXECUTE_STAGING:
- only after developer approval

AI_DEPLOY_PRODUCTION:
- separate explicit approval
- optional two-person approval for critical changes

AI_KNOWLEDGE_UPDATE:
- may collect/index information
- production knowledge changes require approval/versioning

## Risk levels

READ_ONLY
LOW
MEDIUM
HIGH
CRITICAL

Examples:
READ_ONLY: inspect code/logs.
LOW: generate test, documentation, preview.
MEDIUM: modify development branch after approval.
HIGH: database migration, network configuration, service restart.
CRITICAL: production deployment, credential changes, firewall/routing,
destructive migration.

## Internet intelligence

APB AI may remain connected to the internet for:
- current technical documentation;
- package/version information;
- vendor manuals;
- security advisories;
- API changes;
- troubleshooting research.

Internet results are treated as external knowledge, not authority.
The AI must record:
- source;
- retrieval time;
- summary;
- affected component;
- confidence;
- proposed impact.

No external source is allowed to directly execute code or commands.

## Credential safety

AI never receives plaintext:
- MikroTik API password;
- VPN secret;
- RADIUS secret;
- SNMP secret;
- OLT password;
- payment gateway secret.

AI receives only secret references/capabilities where necessary.
Execution services resolve secrets server-side.

## Developer approval object

Each proposed action should contain:
- action_id
- command
- purpose
- risk
- affected_files
- affected_modules
- database_impact
- network_impact
- security_impact
- diff
- tests
- rollback_plan
- evidence
- AI_confidence
- requested_by
- approved_by
- approval_time
- execution_time
- execution_result

## Principle

AI is a DEVELOPMENT COPILOT, not an autonomous administrator.

The default behavior is:

SUGGEST → PREVIEW → TEST → ASK APPROVAL → EXECUTE

Never:

SUGGEST → EXECUTE silently.
