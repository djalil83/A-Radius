import React, { type ReactNode } from "react";
import { Button } from "@/components/ui/button";

type CrudAction = "create" | "edit" | "delete" | "save";

type ProfileCrudActionsProps = {
  readOnly: boolean;
  only?: CrudAction;
  onCreate?: () => void;
  onEdit?: () => void;
  onDelete?: () => void;
  onSave?: () => void;
  className?: string;
  buttonClassName?: string;
  label?: string;
  content?: ReactNode;
  ariaLabel?: string;
  variant?: "default" | "outline" | "ghost";
  size?: "default" | "sm" | "lg" | "icon";
};

export function ProfileCrudActions({ readOnly, only, onCreate, onEdit, onDelete, onSave, className = "flex flex-wrap gap-2", buttonClassName, label, content, ariaLabel, variant = "default", size = "default" }: ProfileCrudActionsProps) {
  const show = (action: CrudAction) => !only || only === action;
  const render = (action: CrudAction, fallback: string, handler?: () => void) => show(action) && (
    <Button disabled={readOnly} onClick={handler} className={buttonClassName} aria-label={ariaLabel} variant={variant} size={size}>
      {content ?? label ?? fallback}
    </Button>
  );
  return (
    <div className={className} aria-label="Profile CRUD actions">
      {render("create", "Profile Baru", onCreate)}
      {render("edit", "Edit", onEdit)}
      {render("delete", "Hapus", onDelete)}
      {render("save", "Simpan", onSave)}
    </div>
  );
}
