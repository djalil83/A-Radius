# Customer Portal Foundation

Customer-facing portal foundation for A-Radius/APB.

## Security

Customer identity is resolved server-side through:

`apb.customer_identities`

The authenticated user must never be allowed to submit an arbitrary
customer ID to access another customer's data.

The portal must use server-side RBAC and authorization.

Customer-facing responses must not expose:

- passwords
- MikroTik credentials
- RADIUS secrets
- OLT credentials
- API keys
- internal audit data
- other customers' data
