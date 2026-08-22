import { describe, expect, it } from "vitest";
import { getLicenseDisplayStatus, isLicenseReadOnly } from "./licenseStatus";

const policy = { warningDays: 30, expiringSoonDays: 7, gracePeriodDays: 3 };
const now = new Date("2026-08-18T00:00:00.000Z");
const license = (days: number, status: "ACTIVE" | "SUSPENDED" | "EXPIRED" | "TRIAL" = "ACTIVE") => ({
  endDate: new Date(now.getTime() + days * 86400000),
  status,
});

describe("getLicenseDisplayStatus", () => {
  it.each([
    [60, "ACTIVE"],
    [30, "WARNING"],
    [7, "EXPIRING SOON"],
    [-4, "EXPIRED"],
  ] as const)("maps %s remaining days to %s", (days, expected) => {
    expect(getLicenseDisplayStatus(license(days), policy, now).label).toBe(expected);
  });

  it("keeps a license in the grace window as expiring soon", () => {
    expect(getLicenseDisplayStatus(license(-2), policy, now).label).toBe("EXPIRING SOON");
  });

  it("prioritizes suspended status", () => {
    expect(getLicenseDisplayStatus(license(60, "SUSPENDED"), policy, now).label).toBe("SUSPENDED");
  });

  it("marks expired and suspended states as read-only for UI controls", () => {
    expect(isLicenseReadOnly(getLicenseDisplayStatus(license(-4), policy, now).label)).toBe(true);
    expect(isLicenseReadOnly(getLicenseDisplayStatus(license(60, "SUSPENDED"), policy, now).label)).toBe(true);
    expect(isLicenseReadOnly("ACTIVE")).toBe(false);
  });
});
