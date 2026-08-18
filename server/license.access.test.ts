import { describe, expect, it } from "vitest";
import { getLicenseWriteBlockForStatus } from "./db";

const now = new Date("2026-08-18T00:00:00.000Z");
const endDate = (days: number) => new Date(now.getTime() + days * 86400000);

describe("license write access", () => {
  it("allows active licenses within the grace window", () => {
    expect(getLicenseWriteBlockForStatus("ACTIVE", endDate(-2), 3, now)).toBeNull();
  });

  it("blocks an expired license after grace period", () => {
    expect(getLicenseWriteBlockForStatus("ACTIVE", endDate(-4), 3, now)).toContain("expired");
  });

  it("blocks explicit expired and suspended statuses", () => {
    expect(getLicenseWriteBlockForStatus("EXPIRED", endDate(30), 3, now)).toContain("expired");
    expect(getLicenseWriteBlockForStatus("SUSPENDED", endDate(30), 3, now)).toContain("suspended");
  });
});
