import { ReactNode } from "react";
import { cn } from "@/lib/utils";

interface FormFooterProps {
  children: ReactNode;
  className?: string;
}

export function FormFooter({ children, className }: FormFooterProps) {
  return (
    <div
      className={cn(
        "-mx-6 -mb-6 mt-6 flex justify-end gap-2 border-t bg-muted/20 px-6 py-4",
        className
      )}
    >
      {children}
    </div>
  );
}
