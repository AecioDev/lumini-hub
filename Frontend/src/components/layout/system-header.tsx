'use client';

import Link from 'next/link';
import { cn } from '@/lib/utils';
import { Logo } from '@/components/ui/logo';
import { UserSection } from './user-section';
import SystemSidebar from './system-sidebar';
import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet';
import { PiList } from 'react-icons/pi';

export function SystemHeader() {
  return (
    <header
      className={cn(
        'sticky top-0 z-40 flex items-center justify-between bg-background/80 px-4 py-4 backdrop-blur-xl md:px-5 lg:px-6 xl:-ms-1.5 xl:pl-4 2xl:-ms-0 2xl:py-5 2xl:pl-6 3xl:px-8 3xl:pl-6 4xl:px-10 4xl:pl-9 h-16'
      )}
    >
      <div className="flex items-center gap-3">
        {/* Hambúrguer Menu para Celular e Tablet */}
        <Sheet>
          <SheetTrigger asChild>
            <button
              aria-label="Abrir Menu"
              className="p-2 -ml-2 rounded-md hover:bg-muted xl:hidden text-foreground"
            >
              <PiList className="h-6 w-6" />
            </button>
          </SheetTrigger>
          <SheetContent side="left" className="p-0 w-[300px] border-r border-border bg-card">
            {/* Renderiza a própria SystemSidebar dentro do Drawer móvel */}
            <SystemSidebar className="static w-full p-0 xl:p-0 2xl:w-full [&>div]:rounded-none" />
          </SheetContent>
        </Sheet>

        {/* Logo do Lumini Hub em telas móveis */}
        <Link
          href="/dashboard"
          aria-label="Lumini Hub Home"
          className="xl:hidden flex items-center"
        >
          <Logo className="h-8 w-8" />
        </Link>
      </div>

      {/* Lado Direito: Theme Settings + Seção do Usuário */}
      <div className="flex items-center gap-2 sm:gap-4">
        <UserSection />
      </div>
    </header>
  );
}
