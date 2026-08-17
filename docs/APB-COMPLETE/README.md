# APB — MASTER COMPLETE BLUEPRINT
Version: APB 1.0 Architecture
Scope: CRM + ISP Billing + Payment Gateway + Payment + RADIUS + MikroTik +
NMS + FTTH/OLT + PPPoE + Hotspot/Voucher + Static IP + Mitra +
Inventory + Helpdesk + Technician + Finance + Local AI.

======================================================================
0. EXECUTIVE ARCHITECTURE
======================================================================

                    ┌───────────────────────────────┐
                    │          APB WEB APP          │
                    │ Admin / Operator / NOC /      │
                    │ Finance / Kasir / Technician  │
                    │ Reseller / Biller / Sales     │
                    └───────────────┬───────────────┘
                                    │ HTTPS
                                    ▼
                    ┌───────────────────────────────┐
                    │       API / DOMAIN LAYER      │
                    │ Auth • RBAC • CRM • Service   │
                    │ Billing • Payment • Network   │
                    │ FTTH • Hotspot • Helpdesk     │
                    │ Inventory • Finance • AI      │
                    └───────────────┬───────────────┘
                                    │
            ┌───────────────────────┼────────────────────────┐
            │                       │                        │
            ▼                       ▼                        ▼
     ┌─────────────┐        ┌──────────────┐        ┌──────────────┐
     │ PostgreSQL  │        │ Redis Worker  │        │ Secret       │
     │ Source of   │        │ Automation    │        │ Service      │
     │ Truth       │        │ Queue/Cache   │        │ Credentials  │
     └─────────────┘        └──────────────┘        └──────────────┘
            │                       │
            └───────────────┬───────┘
                            ▼
                     ┌───────────────┐
                     │ LOCAL AI      │
                     │ LLM + RAG +   │
                     │ Diagnostics   │
                     └───────┬───────┘
                             │
     ┌───────────────┬───────┼──────────┬──────────────┐
     ▼               ▼       ▼          ▼              ▼
 MikroTik          RADIUS    OLT        NMS       Payment Gateway
     │               │       │          │              │
     └───────────────┴───┬───┴──────────┘              │
                         ▼                              ▼
                  FTTH / PPPoE /                  QRIS / VA /
                  Hotspot / IP                    E-Wallet /
                                                   Payment Link

MASTER BUSINESS FLOW
CRM → Service → Billing → Payment Gateway → Payment
→ RADIUS/MikroTik → FTTH/OLT/Hotspot → NMS → Finance

MASTER SUPPORT FLOW
Customer → Ticket → AI Diagnosis → Evidence → Technician
→ Material → Resolution → Service Update → Audit

======================================================================
1. APPLICATION SHELL / UI
======================================================================

Browser/Application Header:
- Home
- navigation
- server time
- notification
- dark/light
- WhatsApp/message
- account/profile

Sidebar:
1 Dashboard
2 CRM / Customers
3 Service
4 Billing
5 Payment
6 Payment Gateway
7 Network
8 Router / MikroTik
9 RADIUS
10 OLT
11 FTTH
12 Hotspot
13 Mitra
14 Helpdesk
15 Technician
16 Inventory
17 Finance
18 Reports
19 AI Center
20 Notifications
21 Audit Trail
22 Settings

Footer:
2020-2026 © APB - ISP MANAGEMENT SYSTEM
Version APB 1.x

======================================================================
2. DASHBOARD & ANALYTICS
======================================================================

Cards:
- Customers Active
- Isolated
- New
- Churn
- Revenue
- Receivable
- Online Sessions
- Total Bandwidth
- Router Online/Down
- OLT Online/Down
- ONU Online/LOS
- Open Tickets
- Technician Tasks
- Voucher Stock

Charts:
- Daily/monthly/yearly revenue
- Customer growth
- Churn
- Traffic In/Out
- Ticket SLA
- Network incidents
- Payment success/failure
- FTTH optical alarms

Quick alerts:
- New ticket
- Payment received
- Router down
- OLT alarm
- ONU LOS
- Payment webhook failure
- High bandwidth
- High CPU

======================================================================
3. CRM ENGINE
======================================================================

Customer:
- Customer ID
- Name
- Identity
- Phone/WhatsApp
- Email
- Address
- GPS
- Documents
- Contact history
- Notes
- Status

Service history:
- Packages
- Installations
- Billing
- Payments
- Sessions
- Tickets
- Technician visits
- Network incidents

Statuses:
ACTIVE / ISOLATED / SUSPENDED / TERMINATED / CHURN

======================================================================
4. SERVICE ENGINE
======================================================================

Supported services:
- FTTH
- PPPoE
- Hotspot
- Voucher
- Static IP

Subscription:
- service ID
- customer
- package
- router/NAS
- RADIUS profile
- IP
- VLAN
- OLT/ONU
- billing cycle
- price
- status

Lifecycle:
PROVISIONING → ACTIVE → DUE → ISOLATED
→ RESTORED or TERMINATED

======================================================================
5. BILLING ENGINE
======================================================================

- Recurring billing
- Invoice
- Invoice items
- Billing cycle
- Due date
- Discount
- Tax
- Late fee
- Reminder
- Partial payment
- Aging
- Auto isolation
- Auto restore
- Refund
- Churn

Automation:
Invoice generated → notification → due date → grace period
→ isolation → payment → restore.

======================================================================
6. PAYMENT GATEWAY ENGINE
======================================================================

Provider-neutral adapter:
- QRIS
- Virtual Account
- E-Wallet
- Payment Link
- Bank Transfer
- Manual Payment

Flow:
Invoice
→ Payment Request
→ Gateway
→ Webhook
→ Signature Verification
→ Payment Transaction
→ Reconciliation
→ Invoice PAID
→ Service Restore/Activation

Required:
- webhook idempotency
- signature verification
- transaction reference
- reconciliation
- settlement
- refund
- failed payment handling

======================================================================
7. PAYMENT ENGINE
======================================================================

- payment methods
- payment transactions
- gateway references
- allocations
- reconciliation
- settlement
- refund
- duplicate prevention
- audit

======================================================================
8. MIKROTIK ENGINE
======================================================================

Router object:
- ID
- Name *
- Connection Type *
  IP Public / VPN RADIUS
- IP Address
- API Port
- API Username
- API Password
- VPN endpoint/user/secret
- RADIUS endpoint/secret
- SNMP configuration
- monitoring status
- online session count
- script metadata

Functions:
- PPPoE provisioning
- Hotspot provisioning
- Voucher generation/sync
- Queue/rate limit
- Address list
- DHCP
- ARP
- firewall counters
- session management
- Kick User
- Lock MAC
- Unlock MAC
- Reset Counter
- Isolation
- Restore
- Traffic
- Resource monitoring

Credential rule:
Credentials encrypted in DB and never returned plaintext to frontend.

======================================================================
9. RADIUS ENGINE
======================================================================

- RADIUS servers
- NAS
- users
- groups
- attributes
- authentication
- authorization
- accounting
- sessions
- Acct-Start/Stop
- Interim updates
- logs

Supports:
PPPoE + Hotspot.

======================================================================
10. HOTSPOT / VOUCHER ENGINE
======================================================================

Sub-engines:
Voucher Engine
User Engine
Session Engine
Package/Profile Engine
MikroTik Integration
RADIUS Integration
Billing Integration
Payment Integration
Outlet Engine

PROFILE VOUCHER:
- Name
- MikroTik Group
- Address List
- Rate Limit
- Shared
- Price
- Reseller Commission
- Color
- Quota
- Duration
- Active Period

STOK VOUCHER:
- Total Stock
- HPP
- Commission
- Total Price
- Generate User
- Generate VC
- Outlet
- Import/Export
- History
- Router/Profile/Mitra filters

VOUCHER TERJUAL:
- quantity
- sales
- monthly sales
- expired
- Lock MAC
- Unlock MAC
- Reset Counter
- Active/Inactive
- Refund
- Recap
- Graph
- Export

VOUCHER ONLINE:
- Kick
- Delete
- Sync
- Username
- Profile
- Uptime
- Upload
- Download
- Router
- Interface
- Server

VOUCHER OFFLINE:
- Username
- Router
- Interface
- Server
- IP
- Download
- Upload
- Last Connected

TEMPLATE VOUCHER:
- Template
- Header
- Row
- Footer
- Preview
- Update
- Delete
- Parameters

======================================================================
11. STATIC IP / IPAM
======================================================================

- IP pools
- CIDR
- allocation
- reservations
- static customer IP
- DHCP leases
- ARP/MAC mapping
- VLAN mapping
- conflict detection
- utilization
- history

======================================================================
12. FTTH / OLT ENGINE
======================================================================

Topology:
OLT → PON → ODC → Backbone Core → ODP → Splitter → ONU/ONT → Customer

OLT:
- vendor
- IP
- credential reference
- PON
- ONU/ONT
- VLAN
- service profile
- line profile
- alarms
- optical power
- provisioning

FTTH inventory:
- ODC
- ODP
- core
- port
- splitter
- fiber route
- customer drop
- GPS
- optical budget
- attenuation

OLT adapters:
Huawei / ZTE / FiberHome / Generic

======================================================================
13. NMS / NETWORK MONITORING
======================================================================

Monitor:
- Ping
- SNMP
- CPU
- memory
- interfaces
- traffic
- packet loss
- latency
- router
- switch
- OLT
- PON
- ONU
- optical power

Events:
- down
- flapping
- high utilization
- high CPU
- high memory
- packet loss
- LOS
- optical abnormality

Correlation:
NMS + RADIUS + MikroTik + OLT + CRM + Ticket.

======================================================================
14. MITRA ENGINE
======================================================================

Types:
RESELLER / BILLER / SALES

Fields:
- ID
- status
- name
- category
- identity
- phone
- address
- username
- password
- join date
- balance
- debt limit
- unique code
- voucher stock
- commission
- outlet

Balance:
positive = credit
zero = neutral
negative = debt

Actions:
- activate
- deactivate
- delete
- edit
- stock
- transaction
- commission
- debt
- top-up

======================================================================
15. HELPDESK / TICKETING
======================================================================

Ticket:
- number
- customer
- service
- category
- priority
- SLA
- description
- attachment
- assignment
- technician
- status
- resolution

Statuses:
OPEN / ASSIGNED / IN_PROGRESS / PENDING / RESOLVED / CLOSED

Categories:
- Internet down
- Slow
- PPPoE
- Hotspot
- WiFi
- ONT
- Fiber
- OLT
- Payment
- Installation
- Other

======================================================================
16. TECHNICIAN / FIELD SERVICE
======================================================================

My Tasks:
- Installation
- Trouble
- Maintenance
- Survey

Execution:
- GPS
- photo evidence
- OPM result
- optical RX/TX
- ONT serial
- MAC
- speed test
- ping
- customer signature

Material:
- request
- approval
- pickup
- usage
- return

AI Assistant:
- guided diagnosis
- checklist
- SOP
- command recommendation
- material recommendation
- report generation

======================================================================
17. INVENTORY / WAREHOUSE
======================================================================

Items:
- ONU/ONT
- WiFi router
- Fiber
- LAN
- Connector
- Fast Connector
- Splitter
- Patch cord
- SFP
- Switch
- Tools

Functions:
- stock in/out
- transfer
- serial tracking
- technician loan
- material request
- usage
- history
- minimum stock alert

======================================================================
18. FINANCE
======================================================================

- revenue
- expense
- cashflow
- receivable
- aging
- commission
- refund
- journal
- accounts
- settlement
- financial reports

======================================================================
19. NOTIFICATION
======================================================================

Channels:
WhatsApp / SMS / Email / Push

Events:
- invoice
- reminder
- payment
- isolation
- restore
- ticket
- technician
- voucher
- network incident

======================================================================
20. RBAC / USER MANAGEMENT
======================================================================

Roles:
Super Admin
Admin
Operator
Kasir
Finance
NOC
Teknisi
Reseller
Biller
Sales

Permissions:
VIEW / CREATE / UPDATE / DELETE / EXPORT / EXECUTE / APPROVE

Scope:
- global
- branch
- router
- warehouse
- mitra
- assigned customer

======================================================================
21. LOCAL AI ENGINE
======================================================================

AI CENTER:
- AI Dashboard
- AI Chat
- Troubleshoot Center
- Customer Device Diagnosis
- Technician Assistant
- Network Diagnosis
- MikroTik Diagnosis
- RADIUS/PPPoE Diagnosis
- OLT/GPON Diagnosis
- FTTH Optical Diagnosis
- NMS Incident Analysis
- Knowledge Base / RAG
- Recommendations
- Automation
- AI Reports
- AI Settings

LOCAL MODEL:
Provider abstraction:
- Ollama
- llama.cpp
- other local inference runtime

RAG:
Documents → chunks → embeddings → vector retrieval
→ local LLM → grounded response.

Knowledge:
- APB SOP
- MikroTik SOP/manual
- OLT manuals
- RADIUS docs
- FTTH SOP
- troubleshooting history
- technician resolutions

======================================================================
22. AI TROUBLESHOOTING ORCHESTRATOR
======================================================================

Input:
Customer + Service + Complaint + Location + Device + Router +
OLT + ONU + PPPoE/Hotspot username + Incident time.

Evidence:
CRM
Billing
RADIUS
MikroTik
OLT
FTTH
NMS
Ticket
Inventory
Historical incidents

Output:
- summary
- probable root cause
- evidence
- confidence
- severity
- next tests
- recommendation
- customer explanation
- technician SOP
- escalation

CUSTOMER OFFLINE:
Billing → Service → RADIUS → MikroTik session
→ ONU → PON/LOS → Optical → Router/LAN.

SLOW INTERNET:
Package → Rate limit → Traffic → CPU → Interface errors
→ Packet loss → Latency → Optical → Congestion.

PPPoE:
User → Service → RADIUS → NAS → Profile → IP Pool → Session.

HOTSPOT:
Voucher → User → Profile → Router → Hotspot → Pool → Session.

FTTH:
OLT → PON → ONU → LOS → Optical → Splitter → ODP → Drop.

======================================================================
23. AI DEVICE DIAGNOSIS
======================================================================

Devices:
ONT/ONU / Router WiFi / CPE / LAN

Inputs:
- LED status
- serial
- MAC
- optical RX/TX
- RSSI
- speed test
- ping
- screenshot/photo
- OPM

Output:
- likely fault
- affected component
- next test
- repair
- material
- priority

======================================================================
24. AI NOC / INCIDENT CORRELATION
======================================================================

Correlate:
- router outage
- OLT outage
- PON alarms
- ONU LOS
- RADIUS rejects
- bandwidth saturation
- customer tickets

Output:
Incident ID
Timeline
Affected devices
Affected customers
Blast radius
Root cause probability
Priority
Suggested remediation
Escalation.

======================================================================
25. AI ACTION SAFETY GATE
======================================================================

READ_ONLY:
status, ping, logs, sessions, traffic, optical data.

LOW_RISK:
sync, monitoring refresh, report.

APPROVAL_REQUIRED:
kick user, queue/profile change, isolation, restore,
ONT provisioning.

HIGH_RISK:
firewall, routing, VLAN, credential changes, destructive
configuration, reboot.

AI cannot bypass:
RBAC + allowlist + approval + audit.

======================================================================
26. CREDENTIAL / SECRET SERVICE
======================================================================

Secrets:
- MikroTik API
- VPN
- RADIUS
- SNMP
- OLT
- Payment Gateway
- WhatsApp/SMS/Email

Rules:
- encrypted at rest
- key rotation
- server-side decryption only
- never log plaintext
- never send plaintext frontend
- reference/masked value in AI prompts
- access audited

======================================================================
27. AUDIT TRAIL
======================================================================

Fields:
timestamp
user
role
IP
session/device
action
module
target
before
after
result
error
correlation ID

Audit:
credential access
network commands
OLT provisioning
payments/refunds
isolation/restore
voucher operations
permission changes
exports
AI action approvals/execution

======================================================================
28. POSTGRESQL DOMAIN MODEL
======================================================================

Core:
users, roles, permissions, user_roles
customers, contacts, addresses, documents
services, service_packages, subscriptions, service_history

Billing:
invoices, invoice_items, billing_cycles, discounts, taxes,
late_fees, receivables

Payment:
payments, payment_methods, payment_gateways,
gateway_transactions, payment_webhooks,
payment_reconciliations, refunds, settlements

Network:
routers, router_credentials, radius_servers, radius_nas,
radius_users, radius_groups, radius_attributes, radius_sessions,
ip_pools, ip_addresses, network_interfaces, network_events

FTTH:
olts, olt_credentials, pon_ports, onus, onts, odcs, odps,
splitters, fiber_cores, fiber_routes, ftth_ports,
optical_measurements, ftth_connections

Hotspot:
hotspot_profiles, voucher_batches, vouchers, voucher_sales,
hotspot_users, hotspot_sessions, voucher_templates, outlets

Mitra:
mitras, mitra_balances, mitra_transactions,
mitra_commissions, mitra_debts

Support:
tickets, ticket_messages, ticket_attachments,
ticket_assignments, ticket_sla, ticket_resolutions

Technician:
technicians, technician_tasks, field_visits,
field_photos, material_requests, material_usages

Inventory:
warehouses, inventory_items, inventory_stock,
inventory_movements, inventory_serials, inventory_transfers

Finance:
accounts, journals, journal_entries, expenses, revenues, cash_accounts

AI:
ai_conversations, ai_messages, ai_sessions, ai_diagnostics,
ai_diagnostic_steps, ai_evidence, ai_recommendations,
ai_actions, ai_action_approvals, ai_incidents,
ai_embeddings, ai_documents, ai_document_chunks,
ai_feedback, ai_model_configs, ai_prompt_versions,
ai_tool_runs, ai_guardrail_events

System:
notifications, notification_logs, audit_logs,
system_settings, api_keys, webhook_logs

======================================================================
29. REDIS WORKER / AUTOMATION
======================================================================

billing.generate
billing.reminder
billing.isolation
billing.restore

payment.webhook
payment.verify
payment.reconcile

voucher.generate
voucher.expire
voucher.sync

network.mikrotik.sync
network.radius.sync
network.olt.sync
network.session.sync

monitoring.ping
monitoring.snmp
monitoring.bandwidth
monitoring.alert

notification.whatsapp
notification.sms
notification.email

ai.diagnosis
ai.embedding
ai.index
ai.network-analysis
ai.ftth-analysis
ai.incident-correlation
ai.report
ai.recommendation

backup.database
backup.configuration

======================================================================
30. COMPLETE PROJECT STRUCTURE
======================================================================

APB/
├── apps/
│   ├── web/
│   │   ├── app/
│   │   │   ├── dashboard/
│   │   │   ├── customers/
│   │   │   ├── service/
│   │   │   ├── billing/
│   │   │   ├── payments/
│   │   │   ├── payment-gateway/
│   │   │   ├── network/
│   │   │   ├── router/
│   │   │   ├── radius/
│   │   │   ├── olt/
│   │   │   ├── ftth/
│   │   │   ├── hotspot/
│   │   │   │   ├── dashboard/
│   │   │   │   ├── users/
│   │   │   │   ├── profiles/
│   │   │   │   ├── profile-voucher/
│   │   │   │   ├── stok-voucher/
│   │   │   │   ├── voucher-terjual/
│   │   │   │   ├── voucher-online/
│   │   │   │   ├── voucher-offline/
│   │   │   │   ├── template-voucher/
│   │   │   │   ├── session/
│   │   │   │   └── outlets/
│   │   │   ├── mitra/
│   │   │   ├── tickets/
│   │   │   ├── technicians/
│   │   │   ├── inventory/
│   │   │   ├── finance/
│   │   │   ├── reports/
│   │   │   ├── ai/
│   │   │   │   ├── dashboard/
│   │   │   │   ├── troubleshoot/
│   │   │   │   ├── customer-device/
│   │   │   │   ├── technician-assistant/
│   │   │   │   ├── network-diagnosis/
│   │   │   │   ├── ftth-diagnosis/
│   │   │   │   ├── mikrotik-diagnosis/
│   │   │   │   ├── olt-diagnosis/
│   │   │   │   ├── radius-diagnosis/
│   │   │   │   ├── knowledge-base/
│   │   │   │   ├── automation/
│   │   │   │   └── reports/
│   │   │   ├── notifications/
│   │   │   ├── audit/
│   │   │   └── settings/
│   │   └── components/
│   ├── api/
│   │   ├── auth/
│   │   ├── customers/
│   │   ├── service/
│   │   ├── billing/
│   │   ├── payment/
│   │   ├── payment-gateway/
│   │   ├── mikrotik/
│   │   ├── radius/
│   │   ├── olt/
│   │   ├── ftth/
│   │   ├── network/
│   │   ├── nms/
│   │   ├── hotspot/
│   │   ├── voucher/
│   │   ├── mitra/
│   │   ├── ticket/
│   │   ├── technician/
│   │   ├── inventory/
│   │   ├── finance/
│   │   ├── reports/
│   │   ├── notification/
│   │   ├── webhook/
│   │   ├── audit/
│   │   └── ai/
│   │       ├── chat/
│   │       ├── troubleshoot/
│   │       ├── customer-device/
│   │       ├── technician/
│   │       ├── network/
│   │       ├── knowledge/
│   │       └── actions/
│   └── worker/
│       ├── billing/
│       ├── payment/
│       ├── voucher/
│       ├── network/
│       ├── monitoring/
│       ├── notification/
│       ├── backup/
│       └── ai/
│           ├── diagnostics/
│           ├── embeddings/
│           ├── recommendations/
│           └── reports/
├── packages/
│   ├── database/
│   ├── auth/
│   ├── ui/
│   ├── mikrotik-sdk/
│   ├── radius/
│   ├── olt-adapters/
│   │   ├── huawei/
│   │   ├── zte/
│   │   ├── fiberhome/
│   │   └── generic/
│   ├── payment-adapters/
│   │   ├── qris/
│   │   ├── virtual-account/
│   │   ├── ewallet/
│   │   └── payment-link/
│   ├── notification-adapters/
│   ├── secret-service/
│   ├── local-ai/
│   │   ├── llm/
│   │   ├── rag/
│   │   ├── embeddings/
│   │   ├── tools/
│   │   ├── guardrails/
│   │   ├── prompts/
│   │   ├── providers/
│   │   ├── diagnostic-rules/
│   │   ├── knowledge/
│   │   └── telemetry/
│   └── shared/
├── database/
│   ├── postgresql/
│   │   ├── schema/
│   │   ├── migrations/
│   │   ├── seeds/
│   │   ├── indexes/
│   │   └── ai/
│   └── redis/
│       ├── queues/
│       ├── cache/
│       ├── sessions/
│       ├── locks/
│       └── ai/
├── docker/
├── migrations/
├── scripts/
├── docs/
│   ├── architecture/
│   ├── api/
│   ├── database/
│   ├── mikrotik/
│   ├── radius/
│   ├── olt/
│   ├── ftth/
│   ├── hotspot/
│   ├── billing/
│   ├── payment/
│   ├── security/
│   ├── deployment/
│   └── ai/
└── monitoring/

======================================================================
31. END-TO-END AUTOMATION
======================================================================

NEW CUSTOMER
CRM → Service → Provision → Invoice → Payment
→ RADIUS/MikroTik/OLT → NMS → Active.

OVERDUE
Invoice → Reminder → Grace Period → Isolation
→ RADIUS/MikroTik → Notification.

PAYMENT
Gateway → Webhook → Verify → Payment → Reconcile
→ Invoice Paid → Restore → Notification.

NEW FTTH INSTALLATION
CRM → Ticket → Technician → Material Request
→ OLT Provisioning → ONT → Optical Test → Activation
→ Billing → NMS.

TROUBLE TICKET
Customer → Helpdesk → AI Evidence Collection
→ Diagnosis → Technician → Material
→ Repair → Test → Resolution → Report → Audit.

NETWORK INCIDENT
NMS → Alert → AI Correlation → Impacted Customers
→ Ticket/Incident → NOC/Technician → Remediation
→ Verification → Closure.

======================================================================
32. PRODUCTION PRINCIPLES
======================================================================

- PostgreSQL is the transactional source of truth.
- Redis handles asynchronous work and fast state.
- Domain services own business rules.
- Integrations use adapters.
- Secrets are isolated from frontend.
- AI is local-first and RAG-grounded.
- AI cannot bypass RBAC/approval.
- Every sensitive operation is audited.
- Webhooks are verified and idempotent.
- Network actions are allowlisted.
- The system is designed for multi-router, multi-RADIUS,
  multi-OLT, multi-gateway and multi-provider environments.
