"use client"

import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { AlertTriangle } from "lucide-react"
import { Button } from "@/components/ui/button"
import { useRouter } from "next/navigation"

interface LowStockAlertProps {
  products: any[]
}

export function LowStockAlert({ products }: LowStockAlertProps) {
  const router = useRouter()

  if (products.length === 0) {
    return null
  }

  return (
    <Alert variant="destructive" className="mb-4">
      <AlertTriangle className="h-4 w-4" />
      <AlertTitle>Alerta de Estoque Baixo</AlertTitle>
      <AlertDescription>
        <p className="mb-2">
          {products.length} {products.length === 1 ? "produto está" : "produtos estão"} com estoque abaixo do mínimo.
        </p>
        <Button variant="outline" size="sm" onClick={() => router.push("/dashboard/inventory?tab=low-stock")}>
          Ver Produtos
        </Button>
      </AlertDescription>
    </Alert>
  )
}

