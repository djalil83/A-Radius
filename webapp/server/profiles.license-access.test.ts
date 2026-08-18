import { beforeEach, describe, expect, it, vi } from "vitest";
import { appRouter } from "./routers";
import type { TrpcContext } from "./_core/context";

const getWriteBlock = vi.hoisted(() => vi.fn());
vi.mock("./db", async () => {
  const actual = await vi.importActual<typeof import("./db")>("./db");
  return { ...actual, getAdministratorLicenseWriteBlock: getWriteBlock };
});

function adminContext(): TrpcContext {
  return {
    user: { id: 7, openId: "license-test", name: "License Test", email: "license@example.com", loginMethod: "test", role: "admin", createdAt: new Date(), updatedAt: new Date(), lastSignedIn: new Date() },
    req: { protocol: "https", headers: {} } as TrpcContext["req"],
    res: {} as TrpcContext["res"],
  };
}

const input = { name: "HOME-20", service: "FTTH" as const, category: "Rumahan" as const, media: "Fiber Optic" as const, color: "#1677FF", isActive: true, description: "", downloadRate: "1M", uploadRate: "2M", price: "150000" };

describe("profile license access", () => {
  beforeEach(() => getWriteBlock.mockResolvedValue("License expired: perubahan profile dinonaktifkan."));

  it.each([
    ["create", (caller: ReturnType<typeof appRouter.createCaller>) => caller.profiles.create(input)],
    ["update", (caller: ReturnType<typeof appRouter.createCaller>) => caller.profiles.update({ ...input, id: 1, version: 1 })],
    ["remove", (caller: ReturnType<typeof appRouter.createCaller>) => caller.profiles.remove({ id: 1 })],
  ])("rejects %s before database access", async (_name, operation) => {
    await expect(operation(appRouter.createCaller(adminContext()))).rejects.toMatchObject({ code: "FORBIDDEN", message: "License expired: perubahan profile dinonaktifkan." });
    expect(getWriteBlock).toHaveBeenCalledWith(7);
  });

  it.each([
    ["create", (caller: ReturnType<typeof appRouter.createCaller>) => caller.profiles.create(input)],
    ["update", (caller: ReturnType<typeof appRouter.createCaller>) => caller.profiles.update({ ...input, id: 1, version: 1 })],
    ["remove", (caller: ReturnType<typeof appRouter.createCaller>) => caller.profiles.remove({ id: 1 })],
  ])("rejects %s for suspended license", async (_name, operation) => {
    getWriteBlock.mockResolvedValue("License suspended: perubahan profile dinonaktifkan.");
    await expect(operation(appRouter.createCaller(adminContext()))).rejects.toMatchObject({ code: "FORBIDDEN", message: "License suspended: perubahan profile dinonaktifkan." });
  });
});
