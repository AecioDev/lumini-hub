import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

// Função para formatar valores monetários
export function formatCurrency(value: number): string {
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency: "BRL",
  }).format(value)
}

// Função para formatar datas
export function formatDate(date: Date): string {
  return new Intl.DateTimeFormat("pt-BR").format(date)
}

// Função para verificar permissões de usuário
export function hasPermission(userRole: string, requiredRoles: string[]): boolean {
  if (!userRole || !requiredRoles.length) return false
  return requiredRoles.includes(userRole) || userRole === "ADMIN"
}

