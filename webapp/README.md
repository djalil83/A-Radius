# A-Radius WebApp

This directory contains the A-Radius Administrator License and subscription-profile web application. It is kept alongside the existing Go backend so the GitHub repository remains backward compatible.

The application uses the managed WebDev full-stack runtime, MySQL via Drizzle/tRPC, Manus OAuth, and RBAC. Development commands are run from this directory with `pnpm install`, `pnpm test`, `pnpm check`, and `pnpm build`.

The latest implementation includes per-administrator license records, configurable status policy, read-only enforcement for expired or suspended licenses, audit/RBAC procedures, and tested CRUD controls.
