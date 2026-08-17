# APB CLIENT APP + ACS — COMPLETE BLUEPRINT

## Customer App
- Dashboard
- Login / OTP
- Profile
- Services
- FTTH
- PPPoE
- Hotspot
- Static IP
- Billing
- Invoice
- Payment Gateway
- Connection monitoring
- Usage
- Tickets
- Troubleshooting
- Notifications
- AI Customer Assistant
- Settings

## Device Management
### Router / ONT
- Device status
- Model
- Serial Number (masked where appropriate)
- Firmware
- Hardware version
- IP
- MAC (masked where appropriate)
- Uptime
- Last Seen
- Optical information
- WAN
- LAN
- WiFi

### WiFi
- WiFi 2.4 GHz
- WiFi 5 GHz
- SSID
- WiFi status
- Connected device count
- Connected device list
- Channel
- Signal
- RX/TX
- Change WiFi password

### LAN
- LAN 1..LAN 4
- Connected/disconnected
- Speed
- Duplex
- Connected device
- Port diagnostics

### Router actions
- Restart router
- Change customer WiFi password
- Change limited router password if supported
- Device information
- Diagnostics

All router actions pass through Customer API → Authorization → ACS.
Customer APP never connects directly to the router.

## ACS
ACS Integration Engine
- Device inventory
- Provisioning
- Parameter read/write allowlist
- Diagnostics
- WiFi management
- LAN management
- Firmware management
- Reboot
- Events
- Device health

Supports the ACS/device management protocol used by the deployed equipment,
such as TR-069/TR-369 where applicable.

## AI
AI can diagnose:
- billing/service state
- PPPoE/session
- router status
- WiFi congestion
- connected-device count
- optical status
- OLT status
- NMS status

AI only recommends actions.
Any action requires explicit customer authorization where customer-facing,
and appropriate operator/developer approval for privileged changes.

## Security
- Customer-scoped access
- RBAC
- API authorization
- rate limiting
- audit trail
- no plaintext network credentials in frontend
- secrets remain server-side
- parameter allowlist for ACS
- sensitive values masked

## Architecture

APP CLIENT
  ↓ HTTPS
CUSTOMER API
  ↓
DEVICE ENGINE
  ↓
ACS INTEGRATION
  ↓
ACS
  ↓
ONT / ROUTER
  ├── WiFi
  ├── LAN
  └── WAN

APB CORE
  ├── CRM
  ├── Service
  ├── Billing
  ├── Payment
  ├── RADIUS
  ├── MikroTik
  ├── FTTH/OLT
  ├── NMS
  ├── Inventory
  ├── Helpdesk
  └── AI
