'use client';

import Link from 'next/link';
import { cn } from '@/lib/utils';
import SimpleBar from '@/components/ui/simplebar';
import { SystemSidebarMenu } from './system-sidebar-menu';
import { Logo } from '@/components/ui/logo';

export default function SystemSidebar({ className }: { className?: string }) {
  return (
    <aside
      className={cn(
        'fixed top-0 start-0 z-50 h-screen w-[284px] xl:p-5 2xl:w-[308px] bg-transparent',
        className
      )}
    >
      <div className="h-full bg-card border border-primary/20 xl:rounded-2xl flex flex-col shadow-sm dark:shadow-none transition-colors duration-200">
        {/* Cabeçalho da Sidebar / Logo */}
        <div className="sticky top-0 z-40 flex items-center gap-3 px-6 pb-5 pt-5 2xl:px-8 2xl:pt-6 border-b border-border">
          <Link href={'/dashboard'} aria-label="Lumini Hub Home" className="flex items-center gap-3">
            <Logo className="h-10 w-10 flex-shrink-0" />
            <div className="flex flex-col">
              <span className="text-lg font-bold text-foreground leading-none">Lumini Hub</span>
              <span className="text-xs text-muted-foreground mt-1">Gestão ERP</span>
            </div>
          </Link>
        </div>

        {/* Menu com Rolagem Customizada */}
        <SimpleBar className="flex-1 min-h-0 mt-4">
          <SystemSidebarMenu />
        </SimpleBar>
      </div>
    </aside>
  );
}
