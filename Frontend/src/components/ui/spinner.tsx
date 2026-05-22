import { Loader2 } from "lucide-react";
import { cn } from "@/lib/utils"; // só se estiver usando tailwind merge util

interface SpinnerProps {
  className?: string;
  fullScreen?: boolean;
}

export function Spinner({ className, fullScreen = false }: SpinnerProps) {
  const base = (
    <Loader2
      className={cn("h-8 w-8 animate-spin text-muted-foreground", className)}
    />
  );

  if (fullScreen) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        {base}
      </div>
    );
  }

  return base;
}
