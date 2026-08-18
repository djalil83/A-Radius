import { TRPCError } from "@trpc/server";
import { and, eq } from "drizzle-orm";
import { z } from "zod";
import { COOKIE_NAME } from "@shared/const";
import { getSessionCookieOptions } from "./_core/cookies";
import { systemRouter } from "./_core/systemRouter";
import { protectedProcedure, publicProcedure, router } from "./_core/trpc";
import { getProfile, getProfileRevisions, getProfileStats, getUserByOpenId, getDb, listProfiles } from "./db";
import { subscriptionProfileRevisions, subscriptionProfiles } from "../drizzle/schema";

const profileInput = z.object({
  name: z.string().min(2).max(120),
  service: z.enum(["FTTH", "PPPoE", "Hotspot / Voucher", "Static IP"]),
  category: z.enum(["Rumahan", "Bisnis", "Dedicated", "Hotspot"]),
  media: z.enum(["Fiber Optic", "Wireless", "LAN", "5G / LTE"]),
  color: z.string().regex(/^#[0-9A-Fa-f]{6}$/),
  isActive: z.boolean(), description: z.string().max(2000).optional(),
  downloadRate: z.string().max(32).default("1M"), uploadRate: z.string().max(32).default("2M"), price: z.string().max(20).default("0"),
});
const adminProcedure = protectedProcedure.use(async ({ ctx, next }) => { if (ctx.user.role !== "admin") throw new TRPCError({ code: "FORBIDDEN", message: "Akses admin diperlukan untuk perubahan data." }); return next(); });

export const appRouter = router({
  system: systemRouter,
  auth: router({ me: publicProcedure.query(opts => opts.ctx.user), logout: publicProcedure.mutation(({ ctx }) => { ctx.res.clearCookie(COOKIE_NAME, { ...getSessionCookieOptions(ctx.req), maxAge: -1 }); return { success: true } as const; }) }),
  profiles: router({
    stats: protectedProcedure.query(() => getProfileStats()),
    list: protectedProcedure.input(z.object({ service: z.string().optional(), isActive: z.boolean().optional(), search: z.string().optional() }).optional()).query(({ input }) => listProfiles(input)),
    get: protectedProcedure.input(z.object({ id: z.number().int().positive() })).query(({ input }) => getProfile(input.id)),
    revisions: protectedProcedure.input(z.object({ id: z.number().int().positive() })).query(({ input }) => getProfileRevisions(input.id)),
    create: adminProcedure.input(profileInput).mutation(async ({ input, ctx }) => { const db = await getDb(); if (!db) throw new TRPCError({ code: "INTERNAL_SERVER_ERROR", message: "Database tidak tersedia." }); const result = await db.insert(subscriptionProfiles).values({ ...input, createdBy: ctx.user.id, updatedBy: ctx.user.id }); const id = Number(result[0].insertId); await db.insert(subscriptionProfileRevisions).values({ profileId: id, version: 1, action: "CREATE", snapshot: JSON.stringify(input), actorId: ctx.user.id }); return getProfile(id); }),
    update: adminProcedure.input(z.object({ id: z.number().int().positive(), version: z.number().int().positive() }).merge(profileInput)).mutation(async ({ input, ctx }) => { const db = await getDb(); if (!db) throw new TRPCError({ code: "INTERNAL_SERVER_ERROR" }); const { id, version, ...data } = input; const result = await db.update(subscriptionProfiles).set({ ...data, version: version + 1, updatedBy: ctx.user.id }).where(and(eq(subscriptionProfiles.id, id), eq(subscriptionProfiles.version, version))); if (Number(result[0]?.affectedRows ?? 0) !== 1) throw new TRPCError({ code: "CONFLICT", message: "Profile berubah. Muat ulang sebelum menyimpan." }); await db.insert(subscriptionProfileRevisions).values({ profileId: id, version: version + 1, action: "UPDATE", snapshot: JSON.stringify(data), actorId: ctx.user.id }); return getProfile(id); }),
    remove: adminProcedure.input(z.object({ id: z.number().int().positive() })).mutation(async ({ input, ctx }) => { const db = await getDb(); if (!db) throw new TRPCError({ code: "INTERNAL_SERVER_ERROR" }); const current = await getProfile(input.id); if (!current) throw new TRPCError({ code: "NOT_FOUND" }); await db.insert(subscriptionProfileRevisions).values({ profileId: input.id, version: current.version + 1, action: "DELETE", snapshot: JSON.stringify(current), actorId: ctx.user.id }); await db.delete(subscriptionProfiles).where(eq(subscriptionProfiles.id, input.id)); return { success: true } as const; }),
    exportCsv: protectedProcedure.input(z.object({ service: z.string().optional(), isActive: z.boolean().optional(), search: z.string().optional() }).optional()).query(async ({ input }) => { const rows = await listProfiles(input); const escape = (v: unknown) => `"${String(v ?? "").replaceAll('"', '""')}"`; return ["Nama,Layanan,Kategori,Media,Rate,Harga,Status", ...rows.map(p => [p.name, p.service, p.category, p.media, `${p.downloadRate}/${p.uploadRate}`, p.price, p.isActive ? "Aktif" : "Nonaktif"].map(escape).join(","))].join("\n"); }),
  }),
});
export type AppRouter = typeof appRouter;
