// src/components/customers/customers-table.tsx
import { useRouter } from "next/navigation";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { Edit, Eye, MoreHorizontal, Trash2 } from "lucide-react";
import { type Customer } from "@/services/customer/customer-service";
import { routes } from "@/config/routes";

interface CustomersTableProps {
  customers: Customer[];
  onConfirmDelete: (customer: Customer) => void;
}

export function CustomersTable({
  customers,
  onConfirmDelete,
}: CustomersTableProps) {
  const router = useRouter();

  const viewCustomerDetails = (id: number) => {
    router.push(`/sales/customers/${id}`);
  };

  const editCustomer = (id: number) => {
    router.push(routes.customers.edit(id));
  };

  // Função utilitária para renderizar valores ou "-" se forem nulos/vazios
  const renderCell = (value?: string | null) => value || "-";

  return (
    <div className="overflow-auto rounded-md border">
      <table className="w-full text-sm">
        <thead className="bg-muted/50">
          <tr className="text-left font-medium">
            <th className="px-4 py-3">Nome</th>
            <th className="px-4 py-3">Documento</th>
            <th className="px-4 py-3">Telefone</th>
            <th className="px-4 py-3">Email</th>
            <th className="px-4 py-3">Cidade/UF</th>
            <th className="px-4 py-3 text-right">Ações</th>
          </tr>
        </thead>
        <tbody>
          {customers.length === 0 ? (
            <tr>
              <td
                colSpan={6}
                className="px-4 py-8 text-center text-muted-foreground"
              >
                Nenhum cliente encontrado
              </td>
            </tr>
          ) : (
            customers.map((customer) => {
              // Ajustando para usar renderCell
              const phone = renderCell(
                customer.contacts?.find((c) => c.type === "phone")?.contact
              );
              const email = renderCell(
                customer.contacts?.find((c) => c.type === "email")?.contact
              );
              const address = customer.addresses?.[0];
              const cityState = renderCell(
                address?.city && address?.city.state
                  ? `${address.city.name}/${address.city.state.uf}`
                  : null
              );

              return (
                <tr key={customer.id} className="border-b">
                  {/* Usando renderCell para campos diretos também */}
                  <td className="px-4 py-3">
                    {renderCell(customer.first_name)}
                  </td>
                  <td className="px-4 py-3">
                    {customer.document_number ? (
                      <span className="rounded-full bg-primary/10 px-2 py-1 text-xs font-medium text-primary">
                        {customer.person_type === "cpf" ? "CPF" : "CNPJ"}:{" "}
                        {customer.document_number}
                      </span>
                    ) : (
                      "-"
                    )}
                  </td>
                  <td className="px-4 py-3">{phone}</td>
                  <td className="px-4 py-3">{email}</td>
                  <td className="px-4 py-3">{cityState}</td>
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
                          onClick={() => viewCustomerDetails(customer.id)}
                        >
                          <Eye className="mr-2 h-4 w-4" />
                          Visualizar
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          onClick={() => editCustomer(customer.id)}
                        >
                          <Edit className="mr-2 h-4 w-4" />
                          Editar
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          onClick={() => onConfirmDelete(customer)}
                          className="text-destructive"
                        >
                          <Trash2 className="mr-2 h-4 w-4" />
                          Excluir
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </td>
                </tr>
              );
            })
          )}
        </tbody>
      </table>
    </div>
  );
}
