'use client';

import { useEffect } from 'react';
import { atom, useAtom } from 'jotai';
import { updateThemeColor } from '@/utils/update-theme-color';
import {
  DEFAULT_PRESET_COLOR_NAME,
  DEFAULT_PRESET_COLORS,
  usePresets,
} from '@/config/color-presets';

// Helper para ler dados do localStorage de forma segura
function getLocalStoragePreset() {
  if (typeof window !== 'undefined') {
    const localStorageValue = localStorage.getItem('lumini-color-preset');
    try {
      return localStorageValue ? JSON.parse(localStorageValue) : DEFAULT_PRESET_COLORS;
    } catch (e) {
      return DEFAULT_PRESET_COLORS;
    }
  }
  return DEFAULT_PRESET_COLORS;
}

const colorPresetsAtom = atom(
  typeof window !== 'undefined' ? getLocalStoragePreset() : DEFAULT_PRESET_COLORS
);

const colorPresetsAtomWithPersistence = atom(
  (get) => get(colorPresetsAtom),
  (get, set, newStorage: any) => {
    set(colorPresetsAtom, newStorage);
    if (typeof window !== 'undefined') {
      localStorage.setItem('lumini-color-preset', JSON.stringify(newStorage));
    }
  }
);

export function useColorPresets() {
  const [colorPresets, setColorPresets] = useAtom(
    colorPresetsAtomWithPersistence
  );
  return {
    colorPresets: colorPresets === null ? DEFAULT_PRESET_COLORS : colorPresets,
    setColorPresets,
  };
}

// Color preset name atom
const colorPresetNameAtom = atom(
  typeof window !== 'undefined'
    ? localStorage.getItem('lumini-color-preset-name') || DEFAULT_PRESET_COLOR_NAME
    : DEFAULT_PRESET_COLOR_NAME
);

const colorPresetNameAtomWithPersistence = atom(
  (get) => get(colorPresetNameAtom),
  (get, set, newStorage: string) => {
    set(colorPresetNameAtom, newStorage);
    if (typeof window !== 'undefined') {
      localStorage.setItem('lumini-color-preset-name', newStorage);
    }
  }
);

export function useColorPresetName() {
  const [colorPresetName, setColorPresetName] = useAtom(
    colorPresetNameAtomWithPersistence
  );
  return {
    colorPresetName:
      colorPresetName === null ? DEFAULT_PRESET_COLOR_NAME : colorPresetName,
    setColorPresetName,
  };
}

// Aplica o preset de cores
export function useApplyColorPreset<T extends Record<string, any>>(
  colorPresets: T
) {
  const COLOR_PRESETS = usePresets();

  useEffect(() => {
    let colorLighter = COLOR_PRESETS[0].colors.lighter;
    let colorLight = COLOR_PRESETS[0].colors.light;
    let colorDefault = COLOR_PRESETS[0].colors.default;
    let colorDark = COLOR_PRESETS[0].colors.dark;
    let colorForeground = COLOR_PRESETS[0].colors.foreground;

    if (colorPresets) {
      colorLighter = colorPresets.lighter;
      colorLight = colorPresets.light;
      colorDefault = colorPresets.default;
      colorDark = colorPresets.dark;
      colorForeground = colorPresets.foreground;
    }

    updateThemeColor(
      colorLighter,
      colorLight,
      colorDefault,
      colorDark,
      colorForeground
    );
  }, [colorPresets, COLOR_PRESETS]);
}
