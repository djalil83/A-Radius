export type LicensePolicy = {
  warningDays: number;
  expiringSoonDays: number;
  gracePeriodDays: number;
};

export type LicenseDisplayStatus = {
  label: "ACTIVE" | "WARNING" | "EXPIRING SOON" | "EXPIRED" | "SUSPENDED";
  tone: "emerald" | "amber" | "rose";
  days: number;
};

export function getLicenseDisplayStatus(
  license: { endDate: Date | string; status: "ACTIVE" | "SUSPENDED" | "EXPIRED" | "TRIAL" },
  policy: LicensePolicy,
  now = new Date(),
): LicenseDisplayStatus {
  const days = Math.ceil((new Date(license.endDate).getTime() - now.getTime()) / 86400000);
  if (license.status === "SUSPENDED") return { label: "SUSPENDED", tone: "rose", days };
  if (days < -policy.gracePeriodDays) return { label: "EXPIRED", tone: "rose", days };
  if (days <= policy.expiringSoonDays) return { label: "EXPIRING SOON", tone: "amber", days };
  if (days <= policy.warningDays) return { label: "WARNING", tone: "amber", days };
  return { label: "ACTIVE", tone: "emerald", days };
}

export function isLicenseReadOnly(status: LicenseDisplayStatus["label"]) {
  return status === "EXPIRED" || status === "SUSPENDED";
}
