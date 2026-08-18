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
