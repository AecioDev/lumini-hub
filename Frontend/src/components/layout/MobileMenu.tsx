import { Drawer } from "antd";
import { SiderContent } from "@/components/layout/SiderContent";
import { useThemeMode } from "@/contexts/ThemeContext";
import { CHROME_BG } from "@/theme/antd-theme";
import type { ApiUserMenuItem } from "@/types/menu";

interface MobileMenuProps {
  open: boolean;
  onClose: () => void;
  menuItems: ApiUserMenuItem[];
}

// Gaveta com o menu completo, para telas em que o menu lateral fixo não cabe (celular, tablet e janelas
// estreitas). O fundo é o mesmo do menu lateral; ela fecha ao tocar fora, e o AppLayout a fecha ao navegar.
export function MobileMenu({ open, onClose, menuItems }: MobileMenuProps) {
  const { mode } = useThemeMode();

  return (
    <Drawer
      open={open}
      onClose={onClose}
      placement="left"
      width={280}
      closable={false}
      styles={{
        header: { display: "none" },
        body: { padding: 0, background: CHROME_BG[mode] },
      }}
    >
      <SiderContent collapsed={false} menuItems={menuItems} showFavorites />
    </Drawer>
  );
}
