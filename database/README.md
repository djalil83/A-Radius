# APB Database

Database layer for A-Radius / APB.

## Components

- PostgreSQL — persistent relational data
- Redis — cache, sessions, locks and queues

## Architecture

PostgreSQL
├── schema
├── migrations
├── seeds
└── indexes

Redis
├── cache
├── sessions
├── locks
└── queues
