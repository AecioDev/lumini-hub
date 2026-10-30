import { MoonOutlined, SunOutlined } from "@ant-design/icons";
import { Switch, Tooltip } from "antd";
import { useThemeMode } from "@/contexts/ThemeContext";

export function ThemeToggle() {
  const { mode, toggleTheme } = useThemeMode();

  return (
    <Tooltip title={mode === "dark" ? "Tema claro" : "Tema escuro"}>
      <Switch
        checked={mode === "dark"}
        aria-label={mode === "dark" ? "Mudar para o tema claro" : "Mudar para o tema escuro"}
        onChange={toggleTheme}
        checkedChildren={<MoonOutlined />}
        unCheckedChildren={<SunOutlined />}
      />
    </Tooltip>
  );
}
