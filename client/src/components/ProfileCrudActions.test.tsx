import React from "react";
import { renderToStaticMarkup } from "react-dom/server";
import { describe, expect, it } from "vitest";
import { ProfileCrudActions } from "./ProfileCrudActions";

const labels = ["Profile Baru", "Edit", "Hapus", "Simpan"];

describe("ProfileCrudActions", () => {
  it("disables every CRUD control in read-only mode", () => {
    const html = renderToStaticMarkup(<ProfileCrudActions readOnly />);
    for (const label of labels) {
      expect(html).toContain(`>${label}</button>`);
    }
    expect((html.match(/disabled=""/g) ?? []).length).toBe(4);
  });

  it("keeps controls enabled for an active license", () => {
    const html = renderToStaticMarkup(<ProfileCrudActions readOnly={false} />);
    expect((html.match(/disabled=""/g) ?? []).length).toBe(0);
  });
});
