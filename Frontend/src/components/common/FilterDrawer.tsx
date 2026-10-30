import { useEffect, useState, type ReactNode } from "react";
import { FilterOutlined } from "@ant-design/icons";
import { Badge, Button, Drawer, Space } from "antd";
import { BREAKPOINT_TABLET } from "@/components/layout/layout-constants";
import { useMediaQuery } from "@/hooks/useMediaQuery";

interface FilterButtonProps {
  /** Quantos filtros estão ativos (vira a bolinha no botão). */
  count: number;
  onClick: () => void;
}

// Botão "Filtros" do cabeçalho da lista. No celular fica só o ícone.
export function FilterButton({ count, onClick }: FilterButtonProps) {
  const showLabel = useMediaQuery(`(min-width: ${BREAKPOINT_TABLET}px)`);
  return (
    <Badge count={count} size="small" offset={[-4, 4]}>
      <Button icon={<FilterOutlined />} onClick={onClick} aria-label="Filtros">
        {showLabel ? "Filtros" : null}
      </Button>
    </Badge>
  );
}

interface FilterDrawerProps<T extends object> {
  open: boolean;
  onClose: () => void;
  /** Filtros aplicados agora (a gaveta abre mostrando estes valores). */
  value: T;
  /** Filtros sem nada marcado: é para onde o botão "Limpar" volta. */
  empty: T;
  onApply: (filters: T) => void;
  /** Os campos de filtro da entidade; recebem o rascunho e uma função para alterar um campo. */
  children: (draft: T, setField: <K extends keyof T>(key: K, fieldValue: T[K]) => void) => ReactNode;
}

// Gaveta de filtros das listas: os principais filtros da entidade, que só valem ao clicar em "Aplicar"
// (o rascunho fica na gaveta; fechar sem aplicar descarta). A busca por texto fica fora, no cabeçalho
// da lista, porque é usada o tempo todo.
export function FilterDrawer<T extends object>({
  open,
  onClose,
  value,
  empty,
  onApply,
  children,
}: FilterDrawerProps<T>) {
  const isPhone = !useMediaQuery(`(min-width: ${BREAKPOINT_TABLET}px)`);
  const [draft, setDraft] = useState<T>(value);

  // Ao abrir, o rascunho volta a ser o que está aplicado.
  useEffect(() => {
    if (open) setDraft(value);
  }, [open, value]);

  const setField = <K extends keyof T>(key: K, fieldValue: T[K]) =>
    setDraft((prev) => ({ ...prev, [key]: fieldValue }));

  const apply = (filters: T) => {
    onApply(filters);
    onClose();
  };

  return (
    <Drawer
      title="Filtros"
      open={open}
      onClose={onClose}
      width={isPhone ? "100%" : 380}
      destroyOnClose
      footer={
        <Space style={{ width: "100%", justifyContent: "flex-end" }}>
          <Button onClick={() => apply(empty)}>Limpar</Button>
          <Button type="primary" onClick={() => apply(draft)}>
            Aplicar
          </Button>
        </Space>
      }
    >
      <div style={{ display: "flex", flexDirection: "column", gap: 16 }}>{children(draft, setField)}</div>
    </Drawer>
  );
}
