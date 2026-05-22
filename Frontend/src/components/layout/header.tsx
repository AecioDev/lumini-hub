"use client";

import { DashboardSwitcher } from "./dashboard-switcher";
import { UserSection } from "./user-section";

export function Header() {
  return (
    <header className="flex h-16 items-center justify-between border-b bg-card px-6">
      {/* Lado Esquerdo: Título do sistema */}
      <div className="text-lg font-semibold">Sistema ERP</div>

      {/* Lado Direito: Switcher de Dashboard + Nova UserSection */}
      <div className="flex items-center gap-2 sm:gap-4">
        <DashboardSwitcher />
        <UserSection />
      </div>
    </header>
  );
}
