# APB LOCAL AI ENGINE

## Tujuan
Local AI menjadi mesin bantuan internal APB untuk:
1. Troubleshooting pelanggan.
2. Diagnosis perangkat pelanggan: ONT/ONU, router WiFi, CPE, LAN.
3. Diagnosis MikroTik.
4. Diagnosis RADIUS/PPPoE/Hotspot.
5. Diagnosis OLT/GPON/FTTH.
6. Bantuan teknisi lapangan.
7. Analisis NMS/monitoring.
8. Pencarian knowledge base/RAG.
9. Rekomendasi tindakan.
10. Pembuatan laporan teknisi otomatis.

## Prinsip
AI tidak langsung mengeksekusi perubahan jaringan berisiko.
AI menghasilkan diagnosis, bukti, confidence, dan rekomendasi.
Action yang mengubah jaringan membutuhkan permission/RBAC dan,
untuk aksi berisiko tinggi, approval manusia.

## MENU WEB
AI CENTER
├── AI Dashboard
├── AI Chat
├── Troubleshoot Center
├── Customer Device Diagnosis
├── Technician Assistant
├── Network Diagnosis
├── MikroTik Diagnosis
├── RADIUS/PPPoE Diagnosis
├── OLT/GPON Diagnosis
├── FTTH Optical Diagnosis
├── NMS Incident Analysis
├── Knowledge Base
├── AI Recommendations
├── Automation
├── AI Reports
└── AI Settings

## TROUBLESHOOT CENTER

Input:
- Customer ID
- Service ID
- Complaint
- Location
- Device/ONT
- Router
- OLT
- ONU
- PPPoE username
- Hotspot username
- Static IP
- Time of incident

AI automatically gathers permitted evidence:
- Customer service status
- Invoice/isolation status
- RADIUS authentication/accounting
- MikroTik session
- Ping/latency/packet loss
- Interface traffic
- CPU/memory
- OLT PON/ONU status
- Optical RX/TX power
- ONU LOS/alarm
- DHCP/ARP/MAC information
- Recent network alarms
- Ticket history
- Previous incidents

Output:
- Problem summary
- Probable root cause
- Evidence
- Confidence score
- Severity
- Recommended tests
- Recommended action
- Customer-safe explanation
- Technician instructions
- Escalation path

## DIAGNOSTIC TREE

Customer offline
→ Check billing/isolation
→ Check RADIUS auth
→ Check MikroTik session
→ Check ONU online
→ Check PON/LOS
→ Check optical power
→ Check customer router
→ Check LAN/WiFi
→ Determine root cause.

Slow internet
→ Check package/rate-limit
→ Check current traffic
→ Check router CPU
→ Check interface errors
→ Check packet loss/latency
→ Check ONU optical quality
→ Check congestion
→ Compare historical baseline.

PPPoE failure
→ User exists?
→ Service active?
→ RADIUS auth?
→ NAS reachable?
→ Secret/attribute?
→ MikroTik PPPoE profile?
→ IP pool?
→ Session/accounting?

Hotspot failure
→ Voucher valid?
→ User enabled?
→ Profile valid?
→ Router reachable?
→ Hotspot server?
→ Address pool?
→ Session?
→ MAC lock?
→ Counter/expiry?

FTTH failure
→ OLT reachable?
→ PON port?
→ ONU online?
→ LOS?
→ RX/TX power?
→ Fiber route?
→ ODP/splitter?
→ Drop cable?
→ ONT power?

## CUSTOMER DEVICE AI

Supported:
- ONT/ONU
- WiFi Router
- CPE
- LAN
- DHCP
- PPPoE
- WiFi signal

Technician can input:
- LED state
- ONT serial
- MAC
- RX/TX optical power
- WiFi RSSI
- Speed test
- Ping
- Screenshot/photo
- OPM result

AI returns:
- likely fault
- probable component
- next test
- repair recommendation
- required material
- estimated priority

## TECHNICIAN ASSISTANT

Features:
- "Apa yang harus saya cek?"
- Guided troubleshooting
- Step-by-step checklist
- OPM interpretation
- Optical budget check
- ONT provisioning checklist
- MikroTik command recommendation
- OLT provisioning checklist
- Material recommendation
- Job completion report
- Customer explanation
- Photo/evidence checklist

Technician mode must show only tools permitted by role.

## NETWORK AI

Analyze:
- Router down
- OLT down
- Interface saturation
- Packet loss
- High latency
- CPU/memory anomaly
- Repeated flapping
- Traffic anomaly
- Authentication failures

Correlate:
NMS + Router + RADIUS + OLT + Ticket + Customer impact.

Output incident:
- Incident ID
- affected devices
- affected customers
- timeline
- probable root cause
- blast radius
- priority
- suggested remediation
- escalation

## FTTH AI

Analyze:
- optical power
- attenuation
- distance
- splitter ratio
- ODP/ODC route
- ONU alarms
- PON utilization
- repeated LOS
- fiber history

AI can identify likely:
- dirty connector
- excessive attenuation
- bad patch cord
- splitter issue
- drop cable issue
- ONU issue
- OLT/PON issue

## MIKROTIK AI

Read-only diagnostic tools by default:
- resource
- interfaces
- PPP active
- hotspot active
- queues
- IP ARP
- DHCP leases
- logs
- firewall counters

AI can recommend commands but must respect:
- RBAC
- command allowlist
- approval
- audit trail

Destructive or high-risk commands are never auto-run.

## RADIUS AI

Analyze:
- Access-Accept/Reject
- timeout
- NAS mismatch
- wrong secret
- expired user
- service disabled
- missing attributes
- IP pool
- accounting gaps

## OLT AI

Analyze:
- PON status
- ONU online/offline
- LOS
- optical power
- distance
- alarms
- serial
- VLAN/service profile

## LOCAL RAG / KNOWLEDGE BASE

Knowledge sources:
- APB documentation
- SOP
- MikroTik manuals
- OLT vendor manuals
- RADIUS documentation
- FTTH SOP
- Technician troubleshooting guides
- Internal incident resolutions

Pipeline:
Documents → chunking → embeddings → vector index →
retrieval → local LLM → grounded response.

Do not use external data by default.
The local AI should remain useful even without Internet access.

## LOCAL MODEL ADAPTER

Provider interface:
LocalLLMProvider
├── Ollama
├── llama.cpp
└── Other local inference runtime

EmbeddingProvider
├── Local embedding model
└── Optional external provider

The application must not hard-code one model.

## AI ACTION GATE

READ_ONLY
- ping
- status
- logs
- session lookup
- traffic
- optical data

LOW_RISK
- sync status
- refresh monitoring
- generate report

APPROVAL_REQUIRED
- disconnect user
- change queue
- change profile
- isolate service
- restore service
- provision ONT
- change router configuration

HIGH_RISK
- firewall changes
- routing changes
- VLAN changes
- credential changes
- destructive configuration
- reboot device

High-risk actions require explicit human confirmation and audit.

## DATABASE

ai_conversations
ai_messages
ai_sessions
ai_diagnostics
ai_diagnostic_steps
ai_evidence
ai_recommendations
ai_actions
ai_action_approvals
ai_incidents
ai_embeddings
ai_documents
ai_document_chunks
ai_feedback
ai_model_configs
ai_prompt_versions
ai_tool_runs
ai_guardrail_events

## REDIS QUEUES

ai.diagnosis
ai.embedding
ai.index
ai.network-analysis
ai.ftth-analysis
ai.incident-correlation
ai.report
ai.recommendation

## API

POST /api/ai/chat
POST /api/ai/troubleshoot
POST /api/ai/diagnose/customer-device
POST /api/ai/diagnose/mikrotik
POST /api/ai/diagnose/radius
POST /api/ai/diagnose/olt
POST /api/ai/diagnose/ftth
POST /api/ai/incidents/analyze
GET  /api/ai/knowledge/search
POST /api/ai/recommendations
POST /api/ai/actions/preview
POST /api/ai/actions/approve
POST /api/ai/actions/execute
POST /api/ai/reports

All action endpoints enforce RBAC, allowlists, approval and audit.

## PRIVACY

AI receives the minimum data required.
Secrets/passwords are never included in prompts.
Sensitive credentials are replaced by references/masked values.
Customer PII access follows role permissions.
AI logs must not contain plaintext credentials.

## EXAMPLE INCIDENT

Customer: Internet mati.

AI:
1. Invoice = PAID.
2. Service = ACTIVE.
3. RADIUS = Access-Accept.
4. MikroTik session = absent.
5. ONU = ONLINE.
6. Optical RX = normal.
7. Router LAN = no DHCP lease.

Conclusion:
Likely customer router/LAN issue.
Confidence: 87%.

Technician:
- Power-cycle customer router.
- Check WAN/PPPoE.
- Check LAN cable.
- Re-test DHCP/PPPoE.
- Upload photo.
- Close ticket if restored.

## AI DASHBOARD

Cards:
- Open AI incidents
- Auto-diagnosed incidents
- Pending approvals
- High-risk alerts
- AI confidence average
- Technician assistance count
- Top recurring faults

Charts:
- faults by category
- root cause distribution
- resolution time
- repeat incidents
- device failure trend
