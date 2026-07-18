"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useToast } from "@/components/ui/use-toast";
import { Loader2 } from "lucide-react";
import IntegrationConfigService from "@/services/integrations/integration-config-service";
import { IntegrationSettings } from "@/services/integrations/integration-config-schema";
import LegacyLookupService from "@/services/integrations/legacy-lookup-service";
import {
  LegacyCompany,
  LegacyLocation,
} from "@/services/integrations/legacy-lookup-schema";

const emptySettings: IntegrationSettings = {
  li_api_key: "",
  li_app_key: "",
  li_webhook_secret: "",
  codtipnot: "",
  codlocarm_oficial: "",
  codlocarm_reserva: "",
  codemp: "",
};

export function IntegrationSettingsForm() {
  const { toast } = useToast();
  const [settings, setSettings] = useState<IntegrationSettings>(emptySettings);
  const [locations, setLocations] = useState<LegacyLocation[]>([]);
  const [companies, setCompanies] = useState<LegacyCompany[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [legacyUnavailable, setLegacyUnavailable] = useState(false);

  useEffect(() => {
    const load = async () => {
      setIsLoading(true);
      try {
        const currentSettings = await IntegrationConfigService.getSettings();
        setSettings(currentSettings);
      } catch (error: any) {
        toast({
          variant: "destructive",
          title: "Erro ao carregar configurações",
          description: error.message,
        });
      }

      try {
        const [locationsResult, companiesResult] = await Promise.all([
          LegacyLookupService.getLocations(),
          LegacyLookupService.getCompanies(),
        ]);
        setLocations(locationsResult);
        setCompanies(companiesResult);
      } catch (error: any) {
        setLegacyUnavailable(true);
      }

      setIsLoading(false);
    };

    load();
  }, [toast]);

  const handleChange = (field: keyof IntegrationSettings, value: string) => {
    setSettings((prev) => ({ ...prev, [field]: value }));
  };

  const handleSave = async () => {
    setIsSaving(true);
    try {
      const updated = await IntegrationConfigService.updateSettings(settings);
      setSettings(updated);
      toast({
        title: "Configurações salvas",
        description: "As configurações de integração foram atualizadas com sucesso.",
      });
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Erro ao salvar configurações",
        description: error.message,
      });
    } finally {
      setIsSaving(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-32">
        <Loader2 className="h-6 w-6 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="li_api_key">Chave de API (Loja Integrada)</Label>
          <Input
            id="li_api_key"
            value={settings.li_api_key}
            onChange={(e) => handleChange("li_api_key", e.target.value)}
            placeholder="Chave de API da loja"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="li_app_key">Chave de Aplicação (Loja Integrada)</Label>
          <Input
            id="li_app_key"
            value={settings.li_app_key}
            onChange={(e) => handleChange("li_app_key", e.target.value)}
            placeholder="Chave do integrador"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="li_webhook_secret">Segredo do Webhook</Label>
          <Input
            id="li_webhook_secret"
            value={settings.li_webhook_secret}
            onChange={(e) => handleChange("li_webhook_secret", e.target.value)}
            placeholder="Segredo compartilhado enviado no header X-Webhook-Secret"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="codtipnot">Código do Tipo de Nota (codtipnot)</Label>
          <Input
            id="codtipnot"
            value={settings.codtipnot}
            onChange={(e) => handleChange("codtipnot", e.target.value)}
            placeholder="Ex: 10"
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="codlocarm_oficial">Local de Armazenamento Oficial</Label>
          <Select
            value={settings.codlocarm_oficial}
            onValueChange={(value) => handleChange("codlocarm_oficial", value)}
            disabled={legacyUnavailable}
          >
            <SelectTrigger id="codlocarm_oficial">
              <SelectValue
                placeholder={
                  legacyUnavailable
                    ? "SQL Server indisponível"
                    : "Selecione o local oficial"
                }
              />
            </SelectTrigger>
            <SelectContent>
              {locations.map((location) => (
                <SelectItem
                  key={location.cod_loc_am}
                  value={String(location.cod_loc_am)}
                >
                  {location.des_loc_am}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-2">
          <Label htmlFor="codlocarm_reserva">Local de Armazenamento de Reserva</Label>
          <Select
            value={settings.codlocarm_reserva}
            onValueChange={(value) => handleChange("codlocarm_reserva", value)}
            disabled={legacyUnavailable}
          >
            <SelectTrigger id="codlocarm_reserva">
              <SelectValue
                placeholder={
                  legacyUnavailable
                    ? "SQL Server indisponível"
                    : "Selecione o local de reserva"
                }
              />
            </SelectTrigger>
            <SelectContent>
              {locations.map((location) => (
                <SelectItem
                  key={location.cod_loc_am}
                  value={String(location.cod_loc_am)}
                >
                  {location.des_loc_am}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="space-y-2">
          <Label htmlFor="codemp">Empresa Vinculada à Loja Integrada</Label>
          <Select
            value={settings.codemp}
            onValueChange={(value) => handleChange("codemp", value)}
            disabled={legacyUnavailable}
          >
            <SelectTrigger id="codemp">
              <SelectValue
                placeholder={
                  legacyUnavailable
                    ? "SQL Server indisponível"
                    : "Selecione a empresa"
                }
              />
            </SelectTrigger>
            <SelectContent>
              {companies.map((company) => (
                <SelectItem
                  key={company.cod_cencus}
                  value={String(company.cod_cencus)}
                >
                  {company.des_cencus}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      {legacyUnavailable && (
        <p className="text-sm text-muted-foreground">
          Não foi possível conectar ao SQL Server legado para carregar locais de
          armazenamento e empresas. Os demais campos ainda podem ser salvos.
        </p>
      )}

      <div className="flex justify-end">
        <Button onClick={handleSave} disabled={isSaving}>
          {isSaving ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Salvando...
            </>
          ) : (
            "Salvar Configurações"
          )}
        </Button>
      </div>
    </div>
  );
}
