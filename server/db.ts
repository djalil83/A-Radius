import { and, desc, eq, like, or, sql } from "drizzle-orm";
import { drizzle } from "drizzle-orm/mysql2";
import { InsertUser, subscriptionProfileRevisions, subscriptionProfiles, users } from "../drizzle/schema";
import { ENV } from "./_core/env";

let _db: ReturnType<typeof drizzle> | null = null;
export async function getDb() {
  if (!_db && process.env.DATABASE_URL) {
    try { _db = drizzle(process.env.DATABASE_URL); } catch (error) { console.warn("[Database] Failed to connect:", error); }
  }
  return _db;
}

export async function upsertUser(user: InsertUser): Promise<void> {
  if (!user.openId) throw new Error("User openId is required for upsert");
  const db = await getDb(); if (!db) return;
  const values: InsertUser = { openId: user.openId, name: user.name, email: user.email, loginMethod: user.loginMethod, role: user.role ?? (user.openId === ENV.ownerOpenId ? "admin" : "user"), lastSignedIn: user.lastSignedIn ?? new Date() };
  await db.insert(users).values(values).onDuplicateKeyUpdate({ set: { name: values.name, email: values.email, loginMethod: values.loginMethod, role: values.role, lastSignedIn: values.lastSignedIn } });
}
export async function getUserByOpenId(openId: string) { const db = await getDb(); if (!db) return undefined; const rows = await db.select().from(users).where(eq(users.openId, openId)).limit(1); return rows[0]; }

export async function listProfiles(filters?: { service?: string; isActive?: boolean; search?: string }) {
  const db = await getDb(); if (!db) return [];
  const conditions = [];
  if (filters?.service) conditions.push(eq(subscriptionProfiles.service, filters.service as any));
  if (filters?.isActive !== undefined) conditions.push(eq(subscriptionProfiles.isActive, filters.isActive));
  if (filters?.search) conditions.push(or(like(subscriptionProfiles.name, `%${filters.search}%`), like(subscriptionProfiles.description, `%${filters.search}%`)));
  return db.select().from(subscriptionProfiles).where(conditions.length ? and(...conditions) : undefined).orderBy(desc(subscriptionProfiles.updatedAt));
}
export async function getProfile(id: number) { const db = await getDb(); if (!db) return undefined; const rows = await db.select().from(subscriptionProfiles).where(eq(subscriptionProfiles.id, id)).limit(1); return rows[0]; }
export async function getProfileRevisions(profileId: number) { const db = await getDb(); if (!db) return []; return db.select().from(subscriptionProfileRevisions).where(eq(subscriptionProfileRevisions.profileId, profileId)).orderBy(desc(subscriptionProfileRevisions.createdAt)); }
export async function getProfileStats() { const db = await getDb(); if (!db) return { active: 0, services: 0, exampleRate: "—", engine: "A-RADIUS" }; const rows = await db.select({ active: sql<number>`sum(case when ${subscriptionProfiles.isActive} = true then 1 else 0 end)`, services: sql<number>`count(distinct ${subscriptionProfiles.service})`, exampleRate: sql<string>`coalesce(max(concat(${subscriptionProfiles.downloadRate}, '/', ${subscriptionProfiles.uploadRate})), '—')` }).from(subscriptionProfiles); return { active: Number(rows[0]?.active ?? 0), services: Number(rows[0]?.services ?? 0), exampleRate: rows[0]?.exampleRate ?? "—", engine: "A-RADIUS" }; }
