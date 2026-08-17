# ACS INTEGRATION

## Flow
Customer App
→ Customer API
→ Device Engine
→ Authorization
→ ACS Adapter
→ ACS
→ ONT/Router

## ACS functions
- Get device information
- Read approved parameters
- Update approved WiFi parameters
- Reboot
- Diagnostics
- Firmware status
- WAN/LAN/WiFi status
- Connected clients
- Event collection

## Safety
The ACS adapter uses an allowlist of parameters/actions.
No arbitrary device command from the customer frontend is accepted.

## Credentials
ACS credentials, device credentials and secrets remain in the
Credential/Secret Service. They are never returned as plaintext to
the customer frontend.
