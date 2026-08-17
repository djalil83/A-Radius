# APB COMPLETE BLUEPRINT

APB = CRM + Service + ISP Billing + Payment Gateway + Payment + RADIUS +
MikroTik Management + NMS + FTTH/OLT + PPPoE + Hotspot/Voucher +
Static IP + Mitra + Inventory + Helpdesk + Technician + Finance.

CORE FLOW
CRM → Service → Billing → Payment Gateway → Payment →
RADIUS/MikroTik → FTTH/OLT/Hotspot → NMS → Finance

CORE INFRASTRUCTURE
PostgreSQL = source of truth
Redis = cache/queue/session/locking
Redis Worker = automation
Credential/Secret Service = encrypted secrets
Audit Trail = security/operation history

SECURITY
- Encrypt MikroTik API, VPN, RADIUS, SNMP, OLT and payment secrets at rest.
- Never return plaintext secrets to frontend.
- Server-side RBAC and authorization.
- Webhook signature verification.
- Audit sensitive actions.
- HTTPS, rate limiting and session controls.

WEB MODULES
dashboard
customers
service
billing
payments
payment-gateway
network
router
radius
olt
ftth
hotspot
mitra
tickets
technicians
inventory
finance
reports
notifications
settings
admin
audit
tools

HOTSPOT ENGINE
dashboard
users
profiles
profile-voucher
stok-voucher
voucher-terjual
voucher-online
voucher-offline
template-voucher
session
outlets
settings

PROFILE VOUCHER
Nama Profile, MikroTik Group, Address List, Rate Limit, Shared,
Harga Jual, Komisi Reseller, Warna, Kuota, Durasi, Masa Aktif.

STOK VOUCHER
Total Stok, HPP, Komisi, Harga, Generate, Buat User, Buat VC,
Outlet, Setting, Import, Export, Riwayat, filter Profile/Router/Mitra.

VOUCHER TERJUAL
Jumlah, Total Penjualan, Bulanan, Expired, Lock MAC, Unlock MAC,
Reset Counter, Aktif/Nonaktif, Refund, Rekap, Grafik, Export.

VOUCHER ONLINE
Kick User, Hapus, Sinkronkan, Username, Profile, Uptime,
Upload, Download, Router, Interface, Server.

VOUCHER OFFLINE
Username, Router, Interface, Server, IP, Download, Upload,
Last Connected.

TEMPLATE VOUCHER
Template, Add/Update/Delete, Preview, Header/Row/Footer editor,
Parameter: #username#, #password#, #profile#, #router#, #harga#,
#masaaktif#, #durasi#, #kuota#.

CRM
Customer profile, identity, phone/WhatsApp, email, address/GPS,
documents, service, billing, payment, tickets and history.
Statuses: ACTIVE, ISOLATED, SUSPENDED, TERMINATED, CHURN.

SERVICE ENGINE
FTTH, PPPoE, Hotspot, Voucher, Static IP.
Lifecycle: PROVISIONING → ACTIVE → DUE → ISOLATED →
RESTORED/TERMINATED.

BILLING ENGINE
Recurring billing, invoices, billing cycle, due date, reminders,
partial payment, discounts, taxes, late fee, aging, isolation,
restore, refund and churn.

PAYMENT GATEWAY
QRIS, Virtual Account, E-Wallet, Payment Link, Bank Transfer,
Manual Payment.
Flow: Invoice → Payment Request → Gateway → Webhook →
Signature Verification → Payment → Reconciliation → PAID →
Activate/Restore.

PAYMENT ENGINE
Transactions, methods, gateway reference, invoice allocation,
verification, reconciliation, refund, settlement, failed/duplicate
payment handling.

MIKROTIK ENGINE
Router ID, name, connection type (IP Public/VPN RADIUS), IP,
API port, API username/password, VPN, RADIUS, SNMP, monitoring,
online sessions and scripts.
Functions: PPPoE, Hotspot, Voucher, Queue, Firewall, Address List,
Kick, Lock/Unlock MAC, Reset Counter, Isolation, Restore, Traffic.

RADIUS ENGINE
Servers, NAS, users, groups, attributes, authentication,
authorization, accounting, sessions, interim accounting and logs.
Supports PPPoE + Hotspot.

FTTH/OLT ENGINE
OLT → PON → ODC → Backbone Core → ODP → Splitter → ONU/ONT → Customer.
OLT: vendor, IP, credentials, PON, ONU/ONT, VLAN, service-port,
profiles, optical power, alarms, provisioning.
FTTH: ODC, ODP, Core, Port, Splitter, Fiber, Route, GPS,
Optical Budget, Attenuation.
Adapters: Huawei, ZTE, FiberHome, Generic.

NMS
Ping, SNMP, CPU, memory, interfaces, traffic, latency, packet loss,
router/switch/OLT status, alarms and events.
Alerts: router down, OLT down, high CPU/memory/bandwidth,
packet loss, ONU LOS, optical abnormality, payment webhook failure.

MITRA
Types: RESELLER, BILLER, SALES.
Fields: ID, status, name, category, identity, phone, address,
username, password, join date, balance, debt limit, unique code,
voucher stock, commission, outlet, transactions.
Balance: positive=credit, zero=neutral, negative=debt.

HELPDESK
Ticket number, customer, category, priority, SLA, description,
attachments, assignment, technician, status, resolution.
OPEN, ASSIGNED, IN_PROGRESS, PENDING, RESOLVED, CLOSED.

TECHNICIAN
My Tasks, installations, trouble tickets, maintenance, GPS,
photo evidence, OPM, ONT serial, MAC, material request,
material usage, job history, customer signature.

INVENTORY
ONU/ONT, WiFi Router, Fiber/LAN cable, connector, splitter,
patch cord, SFP, switch, tools.
Warehouse, stock in/out, transfer, technician loan, serial,
material request, usage, history.

FINANCE
Revenue, expense, cashflow, receivable, aging, commission,
refund, journal, settlement and reports.

NOTIFICATION
WhatsApp, SMS, Email, Push.
Invoice, reminder, payment, isolation, restore, ticket,
technician and voucher notifications.

RBAC
Super Admin, Admin, Operator, Kasir, Teknisi, Finance, NOC,
Reseller, Biller, Sales.
Actions: VIEW, CREATE, UPDATE, DELETE, EXPORT, EXECUTE, APPROVE.

AUDIT TRAIL
Timestamp, user, role, IP, session, action, module, target,
before/after, result, error, correlation ID.
Audit credential access, router commands, OLT provisioning,
payment/refund, isolation/restore, voucher actions, permissions,
exports.

DATABASE DOMAINS
Core: users, roles, permissions, customers, contacts, addresses,
documents, services, packages, subscriptions.
Billing: invoices, invoice_items, billing_cycles, discounts,
taxes, late_fees, receivables.
Payment: payments, methods, gateways, gateway_transactions,
webhooks, reconciliations, refunds, settlements.
Network: routers, credentials, radius_servers, NAS, radius_users,
groups, attributes, sessions, IP pools, addresses, interfaces,
events.
FTTH: OLTs, credentials, PON, ONUs, ONTs, ODC, ODP, splitters,
cores, routes, ports, optical measurements, connections.
Hotspot: profiles, batches, vouchers, sales, sessions, users,
templates, outlets.
Mitra: mitras, balances, transactions, commissions, debts.
Support: tickets, messages, attachments, assignments, SLA,
resolutions.
Technician: technicians, tasks, visits, photos, material requests,
usage.
Inventory: warehouses, items, stock, movements, serials, transfers.
Finance: accounts, journals, entries, expenses, revenues, cash.
System: notifications, logs, audit, settings, API keys, webhooks.

REDIS WORKER
billing.generate / reminder / isolation / restore
payment.webhook / verify / reconcile
voucher.generate / expire / sync
network.mikrotik.sync / radius.sync / olt.sync / session.sync
monitoring.ping / snmp / bandwidth / alert
notification.whatsapp / sms / email
backup.database / configuration

DEPLOYMENT
Reverse Proxy/TLS → Web → API → PostgreSQL + Redis + Worker +
Secret Service + RADIUS + MikroTik + OLT + Monitoring.

MAIN UI
Header: Menu, User, Server Time, Notification, Theme, WhatsApp, Account.
Sidebar: Dashboard, Customer, Service, Billing, Payment, Network,
Router, RADIUS, OLT, FTTH, Hotspot, Mitra, Ticket, Technician,
Inventory, Finance, Reports, Settings.
Content: title, statistics, toolbar, filters, search, table, modal, pagination.
Footer: 2020-2026 © APB - ISP MANAGEMENT SYSTEM.
