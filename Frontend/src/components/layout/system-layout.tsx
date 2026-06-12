'use client';

import React from 'react';
import SystemSidebar from './system-sidebar';
import { SystemHeader } from './system-header';

interface SystemLayoutProps {
  children: React.ReactNode;
}

export default function SystemLayout({ children }: SystemLayoutProps) {
  return (
    <main className="flex min-h-screen flex-grow bg-background text-foreground">
      {/* Sidebar - Visível apenas de xl para cima */}
      <SystemSidebar className="fixed hidden xl:block" />

      {/* Container de Conteúdo e Cabeçalho */}
      <div className="flex w-full flex-col xl:ms-[270px] xl:w-[calc(100%-270px)] 2xl:ms-72 2xl:w-[calc(100%-288px)]">
        <SystemHeader />
        
        {/* Espaço do Conteúdo da Página */}
        <div className="flex flex-grow flex-col px-4 pb-6 pt-6 md:px-5 lg:px-6 lg:pb-8 xl:pl-3 2xl:pl-6 3xl:px-8 3xl:pl-6 3xl:pt-6 4xl:px-10 4xl:pb-9 4xl:pl-9">
          {children}
        </div>
      </div>
    </main>
  );
}
