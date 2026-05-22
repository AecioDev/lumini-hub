"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useToast } from "@/components/ui/use-toast";
import { Edit, Eye, Loader2, MoreHorizontal, Trash2 } from "lucide-react";
import { formatCurrency } from "@/lib/utils";
import ProductService, { type Product } from "@/services/product-service";
import CategoryService, { type Category } from "@/services/category-service";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface ProductsListProps {
  limit?: number;
  showFilters?: boolean;
  showPagination?: boolean;
}

export function ProductsList({
  limit = 10,
  showFilters = false,
  showPagination = true,
}: ProductsListProps) {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [filters, setFilters] = useState({
    name: "",
    category_id: "",
    is_active: "",
  });
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  const { toast } = useToast();
  const router = useRouter();

  // Carregar produtos e categorias
  useEffect(() => {
    const loadData = async () => {
      try {
        setIsLoading(true);

        const [productsResponse, categoriesResponse] = await Promise.all([
          ProductService.getProducts(currentPage, limit, {
            name: filters.name || undefined,
            category_id: filters.category_id
              ? Number.parseInt(filters.category_id)
              : undefined,
            is_active:
              filters.is_active === "true"
                ? true
                : filters.is_active === "false"
                ? false
                : undefined,
          }),
          CategoryService.getCategories(),
        ]);

        setProducts(productsResponse.data);
        setTotalPages(Math.ceil(productsResponse.total / limit));
        setCategories(categoriesResponse);
      } catch (error) {
        console.error("Erro ao carregar produtos:", error);
        toast({
          variant: "destructive",
          title: "Erro ao carregar dados",
          description: "Não foi possível carregar a lista de produtos.",
        });
      } finally {
        setIsLoading(false);
      }
    };

    loadData();
  }, [currentPage, limit, filters, toast]);

  const handleFilterChange = (key: string, value: string) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
    setCurrentPage(1); // Resetar para a primeira página ao filtrar
  };

  const handleDeleteProduct = async () => {
    if (!selectedProduct) return;

    try {
      setIsDeleting(true);
      await ProductService.deleteProduct(selectedProduct.id);

      setProducts(products.filter((p) => p.id !== selectedProduct.id));
      setIsDeleteDialogOpen(false);
      setSelectedProduct(null);

      toast({
        title: "Produto excluído",
        description: "O produto foi excluído com sucesso.",
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Erro ao excluir produto",
        description:
          error.response?.data?.message ||
          "Ocorreu um erro ao excluir o produto.",
      });
    } finally {
      setIsDeleting(false);
    }
  };

  const getCategoryName = (id: number) => {
    const category = categories.find((c) => c.id === id);
    return category ? category.name : "-";
  };

  if (isLoading) {
    return (
      <div className="flex h-40 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {showFilters && (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
          <div>
            <Input
              placeholder="Buscar por nome"
              value={filters.name}
              onChange={(e) => handleFilterChange("name", e.target.value)}
            />
          </div>

          <div>
            <Select
              value={filters.category_id}
              onValueChange={(value) =>
                handleFilterChange("category_id", value)
              }
            >
              <SelectTrigger>
                <SelectValue placeholder="Categoria" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todas as Categorias</SelectItem>
                {categories.map((category) => (
                  <SelectItem key={category.id} value={category.id.toString()}>
                    {category.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div>
            <Select
              value={filters.is_active}
              onValueChange={(value) => handleFilterChange("is_active", value)}
            >
              <SelectTrigger>
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos os Status</SelectItem>
                <SelectItem value="true">Ativos</SelectItem>
                <SelectItem value="false">Inativos</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      )}

      <div className="rounded-md border">
        <table className="w-full">
          <thead>
            <tr className="border-b bg-muted/50 text-left text-sm font-medium">
              <th className="px-4 py-3">Código</th>
              <th className="px-4 py-3">Nome</th>
              <th className="px-4 py-3">Categoria</th>
              <th className="px-4 py-3">Preço</th>
              <th className="px-4 py-3">Estoque</th>
              <th className="px-4 py-3">Status</th>
              <th className="px-4 py-3 text-right">Ações</th>
            </tr>
          </thead>
          <tbody>
            {products.length === 0 ? (
              <tr>
                <td
                  colSpan={7}
                  className="px-4 py-8 text-center text-muted-foreground"
                >
                  Nenhum produto encontrado
                </td>
              </tr>
            ) : (
              products.map((product) => (
                <tr key={product.id} className="border-b">
                  <td className="px-4 py-3">{product.sku}</td>
                  <td className="px-4 py-3">{product.name}</td>
                  <td className="px-4 py-3">
                    {getCategoryName(product.category_id)}
                  </td>
                  <td className="px-4 py-3 font-medium">
                    {formatCurrency(product.selling_price)}
                  </td>
                  <td className="px-4 py-3">
                    <span
                      className={`rounded-full px-2 py-1 text-xs font-medium ${
                        product.current_stock <= product.min_stock
                          ? "bg-red-100 text-red-700"
                          : "bg-green-100 text-green-700"
                      }`}
                    >
                      {product.current_stock}
                    </span>
                  </td>
                  <td className="px-4 py-3">
                    <span
                      className={`rounded-full px-2 py-1 text-xs font-medium ${
                        product.is_active
                          ? "bg-green-100 text-green-700"
                          : "bg-red-100 text-red-700"
                      }`}
                    >
                      {product.is_active ? "Ativo" : "Inativo"}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-right">
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon">
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuLabel>Ações</DropdownMenuLabel>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem
                          onClick={() =>
                            router.push(
                              `/dashboard/inventory/products/${product.id}`
                            )
                          }
                        >
                          <Eye className="mr-2 h-4 w-4" />
                          Visualizar
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          onClick={() =>
                            router.push(
                              `/dashboard/inventory/products/${product.id}/edit`
                            )
                          }
                        >
                          <Edit className="mr-2 h-4 w-4" />
                          Editar
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          onClick={() => {
                            setSelectedProduct(product);
                            setIsDeleteDialogOpen(true);
                          }}
                          className="text-destructive"
                        >
                          <Trash2 className="mr-2 h-4 w-4" />
                          Excluir
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </td>
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
            onClick={() =>
              setCurrentPage((prev) => Math.min(prev + 1, totalPages))
            }
            disabled={currentPage === totalPages}
          >
            Próxima
          </Button>
        </div>
      )}

      {/* Dialog de confirmação de exclusão */}
      <Dialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Confirmar Exclusão</DialogTitle>
            <DialogDescription>
              Tem certeza que deseja excluir este produto? Esta ação não pode
              ser desfeita.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setIsDeleteDialogOpen(false)}
              disabled={isDeleting}
            >
              Cancelar
            </Button>
            <Button
              variant="destructive"
              onClick={handleDeleteProduct}
              disabled={isDeleting}
            >
              {isDeleting ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Excluindo...
                </>
              ) : (
                "Excluir"
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
