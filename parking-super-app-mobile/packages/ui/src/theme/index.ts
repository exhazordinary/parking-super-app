import { MD3LightTheme, MD3DarkTheme, configureFonts } from 'react-native-paper';
import type { MD3Theme } from 'react-native-paper';

/**
 * Parking App Color Palette
 * Primary: Deep blue for trust and reliability
 * Secondary: Teal for action and accessibility
 * Error: Red for warnings and errors
 * Success: Green for confirmations
 */
const colors = {
  // Primary palette
  primary: '#1E3A8A', // Deep blue
  primaryContainer: '#DBEAFE',
  onPrimary: '#FFFFFF',
  onPrimaryContainer: '#1E3A8A',

  // Secondary palette
  secondary: '#0D9488', // Teal
  secondaryContainer: '#CCFBF1',
  onSecondary: '#FFFFFF',
  onSecondaryContainer: '#0D9488',

  // Tertiary palette
  tertiary: '#7C3AED', // Purple for accents
  tertiaryContainer: '#EDE9FE',
  onTertiary: '#FFFFFF',
  onTertiaryContainer: '#7C3AED',

  // Error palette
  error: '#DC2626',
  errorContainer: '#FEE2E2',
  onError: '#FFFFFF',
  onErrorContainer: '#DC2626',

  // Neutral palette
  background: '#FFFFFF',
  surface: '#FFFFFF',
  surfaceVariant: '#F3F4F6',
  surfaceDisabled: '#E5E7EB',
  onBackground: '#111827',
  onSurface: '#111827',
  onSurfaceVariant: '#6B7280',
  onSurfaceDisabled: '#9CA3AF',

  // Outline
  outline: '#D1D5DB',
  outlineVariant: '#E5E7EB',

  // Inverse
  inverseSurface: '#1F2937',
  inverseOnSurface: '#F9FAFB',
  inversePrimary: '#93C5FD',

  // Shadows
  shadow: '#000000',
  scrim: '#000000',

  // Custom app colors
  success: '#059669',
  warning: '#D97706',
  info: '#0284C7',
};

const darkColors = {
  ...colors,
  primary: '#60A5FA',
  primaryContainer: '#1E3A8A',
  onPrimary: '#1E3A8A',
  onPrimaryContainer: '#DBEAFE',

  secondary: '#2DD4BF',
  secondaryContainer: '#0D9488',
  onSecondary: '#0D9488',
  onSecondaryContainer: '#CCFBF1',

  background: '#111827',
  surface: '#1F2937',
  surfaceVariant: '#374151',
  onBackground: '#F9FAFB',
  onSurface: '#F9FAFB',
  onSurfaceVariant: '#D1D5DB',

  outline: '#4B5563',
  outlineVariant: '#374151',

  inverseSurface: '#F9FAFB',
  inverseOnSurface: '#1F2937',
  inversePrimary: '#1E3A8A',
};

/**
 * Typography configuration
 */
const fontConfig = {
  displayLarge: {
    fontFamily: 'System',
    fontSize: 57,
    fontWeight: '400' as const,
    letterSpacing: -0.25,
    lineHeight: 64,
  },
  displayMedium: {
    fontFamily: 'System',
    fontSize: 45,
    fontWeight: '400' as const,
    letterSpacing: 0,
    lineHeight: 52,
  },
  displaySmall: {
    fontFamily: 'System',
    fontSize: 36,
    fontWeight: '400' as const,
    letterSpacing: 0,
    lineHeight: 44,
  },
  headlineLarge: {
    fontFamily: 'System',
    fontSize: 32,
    fontWeight: '600' as const,
    letterSpacing: 0,
    lineHeight: 40,
  },
  headlineMedium: {
    fontFamily: 'System',
    fontSize: 28,
    fontWeight: '600' as const,
    letterSpacing: 0,
    lineHeight: 36,
  },
  headlineSmall: {
    fontFamily: 'System',
    fontSize: 24,
    fontWeight: '600' as const,
    letterSpacing: 0,
    lineHeight: 32,
  },
  titleLarge: {
    fontFamily: 'System',
    fontSize: 22,
    fontWeight: '500' as const,
    letterSpacing: 0,
    lineHeight: 28,
  },
  titleMedium: {
    fontFamily: 'System',
    fontSize: 16,
    fontWeight: '500' as const,
    letterSpacing: 0.15,
    lineHeight: 24,
  },
  titleSmall: {
    fontFamily: 'System',
    fontSize: 14,
    fontWeight: '500' as const,
    letterSpacing: 0.1,
    lineHeight: 20,
  },
  bodyLarge: {
    fontFamily: 'System',
    fontSize: 16,
    fontWeight: '400' as const,
    letterSpacing: 0.5,
    lineHeight: 24,
  },
  bodyMedium: {
    fontFamily: 'System',
    fontSize: 14,
    fontWeight: '400' as const,
    letterSpacing: 0.25,
    lineHeight: 20,
  },
  bodySmall: {
    fontFamily: 'System',
    fontSize: 12,
    fontWeight: '400' as const,
    letterSpacing: 0.4,
    lineHeight: 16,
  },
  labelLarge: {
    fontFamily: 'System',
    fontSize: 14,
    fontWeight: '500' as const,
    letterSpacing: 0.1,
    lineHeight: 20,
  },
  labelMedium: {
    fontFamily: 'System',
    fontSize: 12,
    fontWeight: '500' as const,
    letterSpacing: 0.5,
    lineHeight: 16,
  },
  labelSmall: {
    fontFamily: 'System',
    fontSize: 11,
    fontWeight: '500' as const,
    letterSpacing: 0.5,
    lineHeight: 16,
  },
};

/**
 * Spacing scale
 */
export const spacing = {
  xs: 4,
  sm: 8,
  md: 16,
  lg: 24,
  xl: 32,
  xxl: 48,
} as const;

/**
 * Border radius scale
 */
export const borderRadius = {
  xs: 4,
  sm: 8,
  md: 12,
  lg: 16,
  xl: 24,
  full: 9999,
} as const;

/**
 * Light theme
 */
export const theme: MD3Theme = {
  ...MD3LightTheme,
  colors: {
    ...MD3LightTheme.colors,
    ...colors,
  },
  fonts: configureFonts({ config: fontConfig }),
  roundness: 12,
};

/**
 * Dark theme
 */
export const darkTheme: MD3Theme = {
  ...MD3DarkTheme,
  colors: {
    ...MD3DarkTheme.colors,
    ...darkColors,
  },
  fonts: configureFonts({ config: fontConfig }),
  roundness: 12,
};

/**
 * Custom colors accessible from theme
 */
export const customColors = {
  success: colors.success,
  warning: colors.warning,
  info: colors.info,
};

export type AppTheme = typeof theme;
export type CustomColors = typeof customColors;
