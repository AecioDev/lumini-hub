// src/components/common/form-field.tsx
"use client";

import { Label } from "@/components/ui/label";
import { ReactNode } from "react";

interface FormFieldProps {
  label: string;
  id: string;
  children: ReactNode;
  errorMessage?: string;
  className?: string;
}

export function FormField({
  label,
  id,
  children,
  errorMessage,
  className,
}: FormFieldProps) {
  return (
    <div className={className}>
      <Label htmlFor={id} className="text-foreground">
        {label}
      </Label>
      {children}
      {errorMessage && (
        <p className="text-sm text-destructive mt-1">{errorMessage}</p>
      )}
    </div>
  );
}
