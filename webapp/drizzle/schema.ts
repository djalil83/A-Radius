import { boolean, decimal, index, int, mysqlEnum, mysqlTable, text, timestamp, varchar } from "drizzle-orm/mysql-core";

export const users = mysqlTable("users", {
  id: int("id").autoincrement().primaryKey(),
  openId: varchar("openId", { length: 64 }).notNull().unique(),
  name: text("name"),
  email: varchar("email", { length: 320 }),
  loginMethod: varchar("loginMethod", { length: 64 }),
  role: mysqlEnum("role", ["user", "admin"]).default("user").notNull(),
  createdAt: timestamp("createdAt").defaultNow().notNull(),
  updatedAt: timestamp("updatedAt").defaultNow().onUpdateNow().notNull(),
  lastSignedIn: timestamp("lastSignedIn").defaultNow().notNull(),
});

export const subscriptionProfiles = mysqlTable("subscription_profiles", {
  id: int("id").autoincrement().primaryKey(),
  name: varchar("name", { length: 120 }).notNull(),
  service: mysqlEnum("service", ["FTTH", "PPPoE", "Hotspot / Voucher", "Static IP"]).notNull(),
  category: mysqlEnum("category", ["Rumahan", "Bisnis", "Dedicated", "Hotspot"]).notNull(),
  media: mysqlEnum("media", ["Fiber Optic", "Wireless", "LAN", "5G / LTE"]).notNull(),
  color: varchar("color", { length: 7 }).notNull().default("#1677FF"),
  isActive: boolean("isActive").notNull().default(true),
  description: text("description"),
  downloadRate: varchar("downloadRate", { length: 32 }).notNull().default("1M"),
  uploadRate: varchar("uploadRate", { length: 32 }).notNull().default("2M"),
  price: decimal("price", { precision: 12, scale: 2 }).notNull().default("0"),
  version: int("version").notNull().default(1),
  createdBy: int("createdBy"),
  updatedBy: int("updatedBy"),
  createdAt: timestamp("createdAt").defaultNow().notNull(),
  updatedAt: timestamp("updatedAt").defaultNow().onUpdateNow().notNull(),
}, table => ({ serviceIdx: index("subscription_profiles_service_idx").on(table.service), statusIdx: index("subscription_profiles_status_idx").on(table.isActive) }));

export const subscriptionProfileRevisions = mysqlTable("subscription_profile_revisions", {
  id: int("id").autoincrement().primaryKey(),
  profileId: int("profileId").notNull(),
  version: int("version").notNull(),
  action: mysqlEnum("action", ["CREATE", "UPDATE", "DELETE"]).notNull(),
  snapshot: text("snapshot").notNull(),
  actorId: int("actorId"),
  createdAt: timestamp("createdAt").defaultNow().notNull(),
}, table => ({ profileIdx: index("subscription_profile_revisions_profile_idx").on(table.profileId) }));

export type User = typeof users.$inferSelect;
export type InsertUser = typeof users.$inferInsert;
export type SubscriptionProfile = typeof subscriptionProfiles.$inferSelect;
export type InsertSubscriptionProfile = typeof subscriptionProfiles.$inferInsert;
export type SubscriptionProfileRevision = typeof subscriptionProfileRevisions.$inferSelect;

export const administratorLicenses = mysqlTable("administrator_licenses", {
  id: int("id").autoincrement().primaryKey(),
  administratorId: int("administratorId").notNull().unique(),
  branchId: varchar("branchId", { length: 64 }).notNull(),
  companyName: varchar("companyName", { length: 180 }).notNull(),
  packageName: varchar("packageName", { length: 100 }).notNull(),
  status: mysqlEnum("status", ["ACTIVE", "SUSPENDED", "EXPIRED", "TRIAL"]).notNull().default("ACTIVE"),
  startDate: timestamp("startDate").notNull(),
  endDate: timestamp("endDate").notNull(),
  activePeriodDays: int("activePeriodDays").notNull().default(365),
  price: decimal("price", { precision: 14, scale: 2 }).notNull().default("0"),
  currency: varchar("currency", { length: 3 }).notNull().default("IDR"),
  autoRenewal: boolean("autoRenewal").notNull().default(false),
  dueDate: timestamp("dueDate").notNull(),
  maxCustomers: int("maxCustomers").notNull().default(0),
  maxUsers: int("maxUsers").notNull().default(0),
  maxRouters: int("maxRouters").notNull().default(0),
  maxOlt: int("maxOlt").notNull().default(0),
  maxOdp: int("maxOdp").notNull().default(0),
  maxVouchers: int("maxVouchers").notNull().default(0),
  maxTechnicians: int("maxTechnicians").notNull().default(0),
  maxPartners: int("maxPartners").notNull().default(0),
  storageLimitGb: int("storageLimitGb").notNull().default(0),
  usedCustomers: int("usedCustomers").notNull().default(0),
  usedUsers: int("usedUsers").notNull().default(0),
  usedRouters: int("usedRouters").notNull().default(0),
  usedOlt: int("usedOlt").notNull().default(0),
  usedStorageGb: int("usedStorageGb").notNull().default(0),
  features: text("features").notNull().default("[]"),
  createdAt: timestamp("createdAt").defaultNow().notNull(),
  updatedAt: timestamp("updatedAt").defaultNow().onUpdateNow().notNull(),
});

export const administratorLicenseAudit = mysqlTable("administrator_license_audit", {
  id: int("id").autoincrement().primaryKey(),
  licenseId: int("licenseId").notNull(),
  actorId: int("actorId").notNull(),
  action: mysqlEnum("action", ["CREATE", "EDIT", "EXTEND", "SUSPEND", "ACTIVATE"]).notNull(),
  snapshot: text("snapshot").notNull(),
  createdAt: timestamp("createdAt").defaultNow().notNull(),
}, table => ({ licenseIdx: index("administrator_license_audit_license_idx").on(table.licenseId) }));

export type AdministratorLicense = typeof administratorLicenses.$inferSelect;
export type InsertAdministratorLicense = typeof administratorLicenses.$inferInsert;

export const administratorLicensePolicies = mysqlTable("administrator_license_policies", {
  id: int("id").autoincrement().primaryKey(),
  warningDays: int("warningDays").notNull().default(30),
  expiringSoonDays: int("expiringSoonDays").notNull().default(7),
  gracePeriodDays: int("gracePeriodDays").notNull().default(0),
  updatedBy: int("updatedBy"),
  updatedAt: timestamp("updatedAt").defaultNow().onUpdateNow().notNull(),
});

export type AdministratorLicensePolicy = typeof administratorLicensePolicies.$inferSelect;
export type InsertAdministratorLicensePolicy = typeof administratorLicensePolicies.$inferInsert;
