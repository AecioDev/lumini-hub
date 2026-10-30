import type { Key, ReactNode } from "react";
import { MenuOutlined } from "@ant-design/icons";
import {
  App,
  Button,
  Card,
  Dropdown,
  Empty,
  Pagination,
  Popconfirm,
  Space,
  Spin,
  Table,
  Tooltip,
  theme,
  type MenuProps,
  type TableProps,
} from "antd";
import { Link } from "react-router-dom";
import { BREAKPOINT_TABLET } from "@/components/layout/layout-constants";
import { useActionButtonSize } from "@/hooks/useActionButtonSize";
import { useMediaQuery } from "@/hooks/useMediaQuery";

// Uma ação sobre a linha (ver, editar, excluir). A tela a descreve uma vez só, com as permissões já
// aplicadas, e ela vira botão na coluna "Ações" da tabela (RowActionButtons) e item do menu do cartão.
export interface RowAction {
  key: string;
  label: string;
  icon: ReactNode;
  danger?: boolean;
  /** Rota para onde a ação leva (ex.: a tela de edição). */
  to?: string;
  onClick?: () => void;
  /** Ação destrutiva: pede confirmação antes (Popconfirm na tabela, caixa de confirmação no cartão). */
  confirm?: { title: string; description: string; onConfirm: () => void | Promise<void> };
}

// Botões de ícone da coluna "Ações" da tabela.
export function RowActionButtons({ actions }: { actions: RowAction[] }) {
  const size = useActionButtonSize();
  return (
    <Space size={6}>
      {actions.map((action) => {
        const button = (
          <Button
            size={size}
            danger={action.danger}
            icon={action.icon}
            aria-label={action.label}
            onClick={action.to ? undefined : action.onClick}
          />
        );
        if (action.confirm) {
          return (
            <Popconfirm
              key={action.key}
              title={action.confirm.title}
              description={action.confirm.description}
              okText="Excluir"
              okButtonProps={{ danger: true }}
              cancelText="Cancelar"
              onConfirm={action.confirm.onConfirm}
            >
              {button}
            </Popconfirm>
          );
        }
        const withLink = action.to ? <Link to={action.to}>{button}</Link> : button;
        return (
          <Tooltip key={action.key} title={action.label}>
            {withLink}
          </Tooltip>
        );
      })}
    </Space>
  );
}

// Como uma linha da lista vira um cartão no celular. Cada tela descreve o seu; o resto (tabela, dados,
// paginação, carregamento) é o mesmo.
export interface CardLayout<T> {
  /** Destaque do cartão (ex.: o nome, com avatar). */
  title: (row: T) => ReactNode;
  /** Linha mais fraca logo abaixo do título (ex.: e-mail). */
  subtitle?: (row: T) => ReactNode;
  /** Etiquetas (perfil, status...): o que na tabela fica escondido até rolar para o lado. */
  tags?: (row: T) => ReactNode;
  /** Dados extras em "rótulo: valor"; valor vazio não aparece. */
  fields?: { label: string; value: (row: T) => ReactNode }[];
  /** Ações da linha: as mesmas da coluna "Ações" da tabela, num menu no canto do cartão. */
  actions?: (row: T) => RowAction[];
}

interface ResponsiveTableProps<T> extends TableProps<T> {
  /** Como a linha vira cartão no celular. */
  card: CardLayout<T>;
  /** Dados do celular, quando diferem dos da tabela (ex.: a árvore de empresas vira lista plana). */
  cardDataSource?: T[];
  /** Total de registros mostrado no pé; sem ele vale o `pagination.total` ou a quantidade de linhas. */
  total?: number;
}

const totalLabel = (total: number) => `${total} ${total === 1 ? "registro" : "registros"}`;

// Lista de cadastro: tabela do antd a partir do tablet e cartões no celular (abaixo de BREAKPOINT_TABLET).
// Toda tela de lista usa este componente, para os módulos novos já nascerem com a versão de celular.
export function ResponsiveTable<T extends object>({
  card,
  cardDataSource,
  total,
  ...tableProps
}: ResponsiveTableProps<T>) {
  const isPhone = !useMediaQuery(`(min-width: ${BREAKPOINT_TABLET}px)`);
  const { token } = theme.useToken();
  const { modal } = App.useApp();
  const actionSize = useActionButtonSize();

  const menuFor = (actions: RowAction[]): MenuProps => ({
    items: actions.map((action) => ({
      key: action.key,
      label: action.to ? <Link to={action.to}>{action.label}</Link> : action.label,
      icon: action.icon,
      danger: action.danger,
      onClick: () => {
        if (action.confirm) {
          const { title, description, onConfirm } = action.confirm;
          modal.confirm({
            title,
            content: description,
            okText: "Excluir",
            okButtonProps: { danger: true },
            cancelText: "Cancelar",
            onOk: onConfirm,
          });
        } else {
          action.onClick?.();
        }
      },
    })),
  });

  const pagination = tableProps.pagination === false ? undefined : tableProps.pagination;
  const rows = (cardDataSource ?? tableProps.dataSource ?? []) as T[];
  const totalCount = total ?? pagination?.total ?? rows.length;
  const totalText = (
    <span style={{ color: token.colorTextSecondary, fontSize: 13 }}>{totalLabel(totalCount)}</span>
  );

  if (!isPhone) {
    // Com paginação o total vai à esquerda dela (showTotal); sem paginação, numa linha abaixo da tabela.
    return (
      <>
        <Table<T>
          {...tableProps}
          pagination={
            pagination ? { ...pagination, showTotal: () => totalText } : tableProps.pagination
          }
        />
        {!pagination && <div style={{ marginTop: 12 }}>{totalText}</div>}
      </>
    );
  }

  const rowKey = tableProps.rowKey;
  const keyOf = (row: T, index: number): Key => {
    if (typeof rowKey === "function") return rowKey(row);
    if (typeof rowKey === "string") return (row as Record<string, Key>)[rowKey];
    return index;
  };

  return (
    <Spin spinning={Boolean(tableProps.loading)}>
      {rows.length === 0 && !tableProps.loading ? (
        <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />
      ) : (
        <div style={{ display: "flex", flexDirection: "column", gap: 12 }}>
          {rows.map((row, index) => {
            const fields = (card.fields ?? []).filter((f) => {
              const value = f.value(row);
              return value !== null && value !== undefined && value !== "";
            });
            const actions = card.actions?.(row) ?? [];
            return (
              <Card
                key={keyOf(row, index)}
                size="small"
                style={{ background: token.colorFillQuaternary }}
              >
                <div style={{ display: "flex", alignItems: "flex-start", gap: 8 }}>
                  <div style={{ fontWeight: 500, flex: 1, minWidth: 0 }}>{card.title(row)}</div>
                  {actions.length > 0 && (
                    <Dropdown menu={menuFor(actions)} trigger={["click"]} placement="bottomRight">
                      <Button
                        type="text"
                        size={actionSize}
                        icon={<MenuOutlined />}
                        aria-label="Opções"
                      />
                    </Dropdown>
                  )}
                </div>
                {card.subtitle && (
                  <div style={{ fontSize: 12, color: token.colorTextSecondary, wordBreak: "break-word" }}>
                    {card.subtitle(row)}
                  </div>
                )}
                {card.tags && <div style={{ marginTop: 8 }}>{card.tags(row)}</div>}
                {fields.length > 0 && (
                  <div style={{ marginTop: 8, fontSize: 12 }}>
                    {fields.map((f) => (
                      <div key={f.label}>
                        <span style={{ color: token.colorTextSecondary }}>{f.label}: </span>
                        {f.value(row)}
                      </div>
                    ))}
                  </div>
                )}
              </Card>
            );
          })}
        </div>
      )}
      {/* Mesmo desenho da tabela: total à esquerda e paginação à direita (compacta no celular). */}
      <div
        style={{
          marginTop: 16,
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          flexWrap: "wrap",
          gap: 8,
        }}
      >
        {totalText}
        {pagination && (
          <Pagination
            size="small"
            showLessItems
            current={pagination.current}
            pageSize={pagination.pageSize}
            total={pagination.total}
            showSizeChanger={false}
            onChange={pagination.onChange}
          />
        )}
      </div>
    </Spin>
  );
}
