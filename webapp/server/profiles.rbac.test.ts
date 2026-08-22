import { describe, expect, it } from "vitest";
import { appRouter } from "./routers";
import type { TrpcContext } from "./_core/context";

function context(role: "user" | "admin"): TrpcContext {
  return {
    user: { id: 7, openId: "rbac-test", name: "RBAC Test", email: "rbac@example.com", loginMethod: "test", role, createdAt: new Date(), updatedAt: new Date(), lastSignedIn: new Date() },
    req: { protocol: "https", headers: {} } as TrpcContext["req"],
    res: {} as TrpcContext["res"],
  };
}

describe("profiles RBAC", () => {
  it("rejects create for a regular user before database access", async () => {
    const caller = appRouter.createCaller(context("user"));
    await expect(caller.profiles.create({ name: "HOME-20", service: "FTTH", category: "Rumahan", media: "Fiber Optic", color: "#1677FF", isActive: true, description: "", downloadRate: "1M", uploadRate: "2M", price: "150000" })).rejects.toMatchObject({ code: "FORBIDDEN" });
  });
});
