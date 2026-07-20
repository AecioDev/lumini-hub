import { useEffect, useState } from "react";
import { Button, Card, Col, Row, Skeleton } from "antd";
import { ModulePermissionsPanel, pickDefaultModule } from "@/components/users/forms/ModulePermissionsPanel";
import { ModuleList } from "./ModuleList";
import { permissionService } from "@/services/permissions/permission-service";
import { roleService } from "@/services/roles/role-service";
import { getApiErrorMessage } from "@/utils/api-error";
import { useFeedback } from "@/hooks/useFeedback";
import type { ApiPermissionsByModule } from "@/types/permission";

interface RolePermissionsEditorProps {
  roleId: number;
  initialPermissionIds: number[];
}

export function RolePermissionsEditor({
  roleId,
  initialPermissionIds,
}: RolePermissionsEditorProps) {
  const feedback = useFeedback();

  const [permissionsByModule, setPermissionsByModule] = useState<
    ApiPermissionsByModule[]
  >([]);
  const [loading, setLoading] = useState(true);
  const [selectedIds, setSelectedIds] = useState<Set<number>>(
    new Set(initialPermissionIds)
  );
  const [activeModule, setActiveModule] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    let cancelled = false;
    permissionService
      .byModule()
      .then((result) => {
        if (cancelled) return;
        setPermissionsByModule(result);
        setActiveModule(pickDefaultModule(result, new Set(initialPermissionIds)));
      })
      .catch(() => feedback.error("Erro ao carregar o catálogo de permissões."))
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
    // initialPermissionIds propositalmente fora das deps: só serve pra
    // calcular o módulo padrão na primeira carga, não deve disparar refetch.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [feedback]);

  const handleSave = async () => {
    setSaving(true);
    try {
      await roleService.updatePermissions(roleId, { permission_ids: [...selectedIds] });
      feedback.success("Permissões do perfil atualizadas com sucesso.");
    } catch (error) {
      feedback.error(getApiErrorMessage(error, "Erro ao salvar permissões do perfil."));
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return <Skeleton active paragraph={{ rows: 6 }} />;
  }

  return (
    <Row gutter={[16, 16]}>
      <Col xs={24} sm={8} lg={6}>
        <Card title="Módulos">
          <ModuleList
            permissionsByModule={permissionsByModule}
            selectedIds={selectedIds}
            activeModule={activeModule}
            onSelectModule={setActiveModule}
          />
        </Card>
      </Col>
      <Col xs={24} sm={16} lg={18}>
        <Card title="Permissões">
          <ModulePermissionsPanel
            permissionsByModule={permissionsByModule}
            selectedIds={selectedIds}
            onChange={setSelectedIds}
            activeModule={activeModule}
          />
          <div style={{ display: "flex", justifyContent: "flex-end", marginTop: 20 }}>
            <Button type="primary" loading={saving} onClick={() => void handleSave()}>
              Salvar Permissões
            </Button>
          </div>
        </Card>
      </Col>
    </Row>
  );
}
