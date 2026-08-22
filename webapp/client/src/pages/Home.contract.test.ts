import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

const homeSource = readFileSync(new URL("./Home.tsx", import.meta.url), "utf8");

describe("Home profile CRUD read-only contract", () => {
  it("wires create, save, edit, and delete controls through the shared guard", () => {
    for (const action of ["create", "save", "edit", "delete"]) {
      expect(homeSource).toContain(`ProfileCrudActions only="${action}"`);
    }
    expect((homeSource.match(/readOnly=\{licenseReadOnly/g) ?? []).length).toBeGreaterThanOrEqual(4);
  });

  it("keeps the backend mutation guard in the submit path", () => {
    expect(homeSource).toContain('if (licenseReadOnly) return toast.error("License expired: mode read-only aktif.")');
  });
});
