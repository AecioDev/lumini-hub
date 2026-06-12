import type React from "react";
import SystemLayout from "@/components/layout/system-layout";

export default function RouteLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <SystemLayout>{children}</SystemLayout>;
}
