# Redis

Redis is used by A-Radius / APB for:

- cache
- sessions
- distributed locks
- background queues

Redis must not be treated as the permanent source of truth.
PostgreSQL remains the primary persistent database.
