'use client'

import * as React from 'react'
import {
  ThemeProvider as NextThemesProvider,
  type ThemeProviderProps,
} from 'next-themes'
import { useColorPresets, useApplyColorPreset } from '@/hooks/use-theme-color'

export function ThemeProvider({ children, ...props }: ThemeProviderProps) {
  const { colorPresets } = useColorPresets()
  useApplyColorPreset(colorPresets)

  return <NextThemesProvider {...props}>{children}</NextThemesProvider>
}
