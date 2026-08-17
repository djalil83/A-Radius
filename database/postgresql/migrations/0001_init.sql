-- A-Radius / APB
-- Initial database migration

BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS apb;

COMMIT;
