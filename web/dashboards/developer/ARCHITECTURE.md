# Developer Dashboard Architecture

The Developer Dashboard is divided into:

- `components/` — reusable UI components
- `modules/ai/` — AI security intelligence
- `modules/security/` — security operations
- `modules/development/` — development tooling
- `modules/approval/` — approval gates
- `modules/production/` — release/deployment controls
- `services/` — API, permission and state boundaries

## Security rule

The frontend never grants permission.

Authorization remains server-side through the existing RBAC engine.

High-risk actions follow:

AI/Analysis
→ Finding
→ Recommendation
→ Preview
→ Test
→ Approval
→ Staging
→ Production

AI must never obtain implicit production permission.
