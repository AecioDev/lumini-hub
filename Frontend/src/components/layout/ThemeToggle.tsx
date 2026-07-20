import { MoonOutlined, SunOutlined } from "@ant-design/icons";
import { Switch, Tooltip } from "antd";
import { useThemeMode } from "@/contexts/ThemeContext";

export function ThemeToggle() {
  const { mode, toggleTheme } = useThemeMode();

  return (
    <Tooltip title={mode === "dark" ? "Tema claro" : "Tema escuro"}>
      <Switch
        checked={mode === "dark"}
        onChange={toggleTheme}
        checkedChildren={<MoonOutlined />}
        unCheckedChildren={<SunOutlined />}
      />
    </Tooltip>
  );
}
