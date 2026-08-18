import { describe, expect, it } from "vitest";
import { appRouter } from "./routers";
import type { TrpcContext } from "./_core/context";

function context(role: "user" | "admin"): TrpcContext {
  return {
    user: { id: 7, openId: "license-rbac-test", name: "License RBAC Test", email: "license@example.com", loginMethod: "test", role, createdAt: new Date(), updatedAt: new Date(), lastSignedIn: new Date() },
    req: { protocol: "https", headers: {} } as TrpcContext["req"],
    res: {} as TrpcContext["res"],
  };
}

describe("administrator license RBAC", () => {
  it("rejects extend for a regular user before database access", async () => {
    const caller = appRouter.createCaller(context("user"));
    await expect(caller.license.extend({ days: 365 })).rejects.toMatchObject({ code: "FORBIDDEN" });
  });

  it("rejects status changes for a regular user before database access", async () => {
    const caller = appRouter.createCaller(context("user"));
    await expect(caller.license.setStatus({ status: "SUSPENDED" })).rejects.toMatchObject({ code: "FORBIDDEN" });
  });
});
