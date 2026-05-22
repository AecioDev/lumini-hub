import type React from "react";
import { ThemeProvider } from "@/components/theme-provider";
import { AuthProvider } from "@/contexts/auth-context";
import { Toaster } from "@/components/ui/toaster";
import { Inter } from "next/font/google";
import "../styles/globals.css";
import { HydrationGate } from "@/components/common/hydration-gate";

const inter = Inter({ subsets: ["latin"] });

export const metadata = {
  title: "Sistema ERP - Gestão Empresarial",
  description: "Sistema integrado de gestão empresarial",
  generator: "Eu Mesmo",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="pt-BR" suppressHydrationWarning>
      <body className={inter.className}>
        <ThemeProvider
          attribute="class"
          defaultTheme="light"
          enableSystem
          disableTransitionOnChange
        >
          <AuthProvider>
            <HydrationGate>
              {children}
              <Toaster />
            </HydrationGate>
          </AuthProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
