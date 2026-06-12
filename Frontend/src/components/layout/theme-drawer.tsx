'use client';

import { useTheme } from 'next-themes';
import { PiCheckBold, PiSun, PiMoon, PiGear } from 'react-icons/pi';
import { cn } from '@/lib/utils';
import { usePresets, presetDark, presetLight } from '@/config/color-presets';
import { useColorPresetName, useColorPresets } from '@/hooks/use-theme-color';
import { useEffect, useState } from 'react';
import { updateThemeColor } from '@/utils/update-theme-color';
import { Button } from '@/components/ui/button';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';

export default function ThemeDrawer() {
  const { theme, setTheme } = useTheme();
  const COLOR_PRESETS = usePresets();
  const { colorPresets, setColorPresets } = useColorPresets();
  const { colorPresetName, setColorPresetName } = useColorPresetName();
  const [mounted, setMounted] = useState(false);

  // Evita erros de hidratação
  useEffect(() => {
    setMounted(true);
  }, []);

  // Sincroniza a cor caso mude para Black e o tema seja alterado
  useEffect(() => {
    if (!mounted) return;
    if (colorPresetName === 'black') {
      if (theme === 'light') {
        updateThemeColor(
          presetLight.lighter,
          presetLight.light,
          presetLight.default,
          presetLight.dark,
          presetLight.foreground
        );
      } else if (theme === 'dark') {
        updateThemeColor(
          presetDark.lighter,
          presetDark.light,
          presetDark.default,
          presetDark.dark,
          presetDark.foreground
        );
      }
    }
  }, [theme, colorPresetName, mounted]);

  if (!mounted) return null;

  return (
    <Sheet>
      <SheetTrigger asChild>
        <Button
          variant="outline"
          size="icon"
          className="rounded-full bg-background border-border text-foreground hover:bg-muted"
          title="Configurações de Tema"
        >
          <PiGear className="h-[18px] w-[18px] animate-spin-slow" />
        </Button>
      </SheetTrigger>
      <SheetContent className="w-[320px] sm:w-[400px] bg-background border-l border-border text-foreground">
        <SheetHeader className="pb-6 border-b border-border">
          <SheetTitle>Preferências do Sistema</SheetTitle>
          <SheetDescription>
            Personalize a aparência e as cores do seu ERP.
          </SheetDescription>
        </SheetHeader>

        <div className="py-6 space-y-8">
          {/* Seção 1: Aparência (Light/Dark) */}
          <div className="space-y-4">
            <h4 className="text-sm font-semibold tracking-wide uppercase text-muted-foreground">
              Aparência
            </h4>
            <div className="grid grid-cols-2 gap-4">
              <Button
                variant="outline"
                className={cn(
                  'h-20 flex flex-col items-center justify-center gap-2 rounded-lg border-2 bg-card text-foreground transition duration-300',
                  theme === 'light'
                    ? 'border-primary ring-2 ring-primary/20 bg-muted/40'
                    : 'border-border hover:border-foreground'
                )}
                onClick={() => setTheme('light')}
              >
                <PiSun className="h-6 w-6 text-yellow-500" />
                <span className="text-xs font-medium">Claro</span>
              </Button>
              <Button
                variant="outline"
                className={cn(
                  'h-20 flex flex-col items-center justify-center gap-2 rounded-lg border-2 bg-card text-foreground transition duration-300',
                  theme === 'dark'
                    ? 'border-primary ring-2 ring-primary/20 bg-muted/40'
                    : 'border-border hover:border-foreground'
                )}
                onClick={() => setTheme('dark')}
              >
                <PiMoon className="h-6 w-6 text-indigo-400" />
                <span className="text-xs font-medium">Escuro</span>
              </Button>
            </div>
          </div>

          {/* Seção 2: Cores */}
          <div className="space-y-4">
            <h4 className="text-sm font-semibold tracking-wide uppercase text-muted-foreground">
              Cores de Destaque
            </h4>
            <div className="grid grid-cols-2 gap-4 sm:grid-cols-3">
              {COLOR_PRESETS.map((preset) => {
                const isSelected =
                  colorPresetName?.toLowerCase() === preset.name.toLowerCase();

                return (
                  <div
                    key={preset.name}
                    className="flex flex-col items-center justify-center gap-2"
                  >
                    <button
                      title={preset.name}
                      onClick={() => {
                        setColorPresets(preset.colors);
                        setColorPresetName(preset.name.toLowerCase());
                      }}
                      className={cn(
                        'grid h-10 w-full place-content-center rounded-lg border-2 border-transparent transition duration-300 focus:outline-none shadow-sm',
                        isSelected
                          ? 'ring-2 ring-primary ring-offset-2 dark:ring-offset-gray-900'
                          : 'hover:border-foreground'
                      )}
                      style={{ backgroundColor: preset.colors.default }}
                    >
                      <PiCheckBold
                        className={cn(
                          'h-4 w-4',
                          isSelected ? 'text-white' : 'text-transparent',
                          isSelected && preset.name === 'Black'
                            ? 'text-black dark:text-white'
                            : ''
                        )}
                      />
                    </button>
                    <span
                      className={cn(
                        'text-xs font-medium',
                        isSelected ? 'font-bold' : 'text-muted-foreground'
                      )}
                    >
                      {theme === 'dark' && preset.name === 'Black'
                        ? 'Branco'
                        : preset.name === 'Black'
                        ? 'Preto'
                        : preset.name === 'Blue'
                        ? 'Azul'
                        : preset.name === 'Teal'
                        ? 'Verde-Água'
                        : preset.name === 'Violet'
                        ? 'Violeta'
                        : preset.name === 'Rose'
                        ? 'Rosa'
                        : preset.name === 'Yellow'
                        ? 'Amarelo'
                        : preset.name}
                    </span>
                  </div>
                );
              })}
            </div>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}
