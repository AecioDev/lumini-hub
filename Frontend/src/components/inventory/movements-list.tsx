"use client"

import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { useToast } from "@/components/ui/use-toast"
import { Loader2 } from "lucide-react"
import { format } from "date-fns"
import InventoryService, { type InventoryMovement } from "@/services/inventory-service"
import ProductService, { type Product } from "@/services/product-service"

interface MovementsListProps {
  limit?: number
  showFilters?: boolean
  showPagination?: boolean
  productId?: number
}

export function InventoryMovementsList({
  limit = 10,
  showFilters = false,
  showPagination = true,
  productId,
}: MovementsListProps) {
  const [movements, setMovements] = useState<InventoryMovement[]>([])
  const [products, setProducts] = useState<Product[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [currentPage, setCurrentPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [filters, setFilters] = useState({
    product_id: productId?.toString() || "",
    movement_type: "",
    startDate: "",
    endDate: "",
  })

  const { toast } = useToast()

  // Carregar movimentações e produtos
  useEffect(() => {
    const loadData = async () => {
      try {
        setIsLoading(true)

        const [movementsResponse, productsResponse] = await Promise.all([
          InventoryService.getMovements(currentPage, limit, {
            product_id: filters.product_id ? Number.parseInt(filters.product_id) : undefined,
            movement_type: filters.movement_type || undefined,
            startDate: filters.startDate || undefined,
            endDate: filters.endDate || undefined,
          }),
          ProductService.getProducts(1, 100),
        ])

        setMovements(movementsResponse.data)
        setTotalPages(Math.ceil(movementsResponse.total / limit))
        setProducts(productsResponse.data)
      } catch (error) {
        console.error("Erro ao carregar movimentações:", error)
        toast({
          variant: "destructive",
          title: "Erro ao carregar dados",
          description: "Não foi possível carregar as movimentações de estoque.",
        })
      } finally {
        setIsLoading(false)
      }
    }

    loadData()
  }, [currentPage, limit, filters, toast, productId])

  const handleFilterChange = (key: string, value: string) => {
    setFilters((prev) => ({ ...prev, [key]: value }))
    setCurrentPage(1) // Resetar para a primeira página ao filtrar
  }

  const getProductName = (id: number) => {
    const product = products.find((p) => p.id === id)
    return product ? product.name : "-"
  }

  const getMovementTypeLabel = (type: string) => {
    switch (type) {
      case "entrada":
        return "Entrada"
      case "saida":
        return "Saída"
      case "ajuste":
        return "Ajuste"
      default:
        return type
    }
  }

  if (isLoading) {
    return (
      <div className="flex h-40 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {showFilters && (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-4">
          {!productId && (
            <div>
              <Select value={filters.product_id} onValueChange={(value) => handleFilterChange("product_id", value)}>
                <SelectTrigger>
                  <SelectValue placeholder="Produto" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">Todos os Produtos</SelectItem>
                  {products.map((product) => (
                    <SelectItem key={product.id} value={product.id.toString()}>
                      {product.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          )}

          <div>
            <Select value={filters.movement_type} onValueChange={(value) => handleFilterChange("movement_type", value)}>
              <SelectTrigger>
                <SelectValue placeholder="Tipo de Movimentação" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos os Tipos</SelectItem>
                <SelectItem value="entrada">Entrada</SelectItem>
                <SelectItem value="saida">Saída</SelectItem>
                <SelectItem value="ajuste">Ajuste</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div>
            <Input
              type="date"
              placeholder="Data Inicial"
              value={filters.startDate}
              onChange={(e) => handleFilterChange("startDate", e.target.value)}
            />
          </div>

          <div>
            <Input
              type="date"
              placeholder="Data Final"
              value={filters.endDate}
              onChange={(e) => handleFilterChange("endDate", e.target.value)}
            />
          </div>
        </div>
      )}

      <div className="rounded-md border">
        <table className="w-full">
          <thead>
            <tr className="border-b bg-muted/50 text-left text-sm font-medium">
              <th className="px-4 py-3">Data</th>
              {!productId && <th className="px-4 py-3">Produto</th>}
              <th className="px-4 py-3">Tipo</th>
              <th className="px-4 py-3">Quantidade</th>
              <th className="px-4 py-3">Estoque Anterior</th>
              <th className="px-4 py-3">Novo Estoque</th>
              <th className="px-4 py-3">Observações</th>
            </tr>
          </thead>
          <tbody>
            {movements.length === 0 ? (
              <tr>
                <td colSpan={productId ? 6 : 7} className="px-4 py-8 text-center text-muted-foreground">
                  Nenhuma movimentação encontrada
                </td>
              </tr>
            ) : (
              movements.map((movement) => (
                <tr key={movement.id} className="border-b">
                  <td className="px-4 py-3">{format(new Date(movement.created_at), "dd/MM/yyyy HH:mm")}</td>
                  {!productId && <td className="px-4 py-3">{getProductName(movement.product_id)}</td>}
                  <td className="px-4 py-3">
                    <span
                      className={`rounded-full px-2 py-1 text-xs font-medium ${
                        movement.movement_type === "entrada"
                          ? "bg-green-100 text-green-700"
                          : movement.movement_type === "saida"
                            ? "bg-red-100 text-red-700"
                            : "bg-amber-100 text-amber-700"
                      }`}
                    >
                      {getMovementTypeLabel(movement.movement_type)}
                    </span>
                  </td>
                  <td className="px-4 py-3 font-medium">
                    {movement.quantity > 0 ? `+${movement.quantity}` : movement.quantity}
                  </td>
                  <td className="px-4 py-3">{movement.previous_stock}</td>
                  <td className="px-4 py-3">{movement.new_stock}</td>
                  <td className="px-4 py-3 max-w-xs truncate">{movement.notes || "-"}</td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>

      {showPagination && totalPages > 1 && (
        <div className="flex items-center justify-end space-x-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
            disabled={currentPage === 1}
          >
            Anterior
          </Button>
          <span className="text-sm text-muted-foreground">
            Página {currentPage} de {totalPages}
          </span>
          <Button
            variant="outline"
            size="sm"
            onClick={() => setCurrentPage((prev) => Math.min(prev + 1, totalPages))}
            disabled={currentPage === totalPages}
          >
            Próxima
          </Button>
        </div>
      )}
    </div>
  )
}

