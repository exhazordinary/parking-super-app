import React from 'react';
import { StyleSheet, ViewStyle } from 'react-native';
import { Button as PaperButton, useTheme } from 'react-native-paper';
import type { MD3Theme } from 'react-native-paper';

export interface ButtonProps {
  /** Button text */
  children: React.ReactNode;
  /** Button variant */
  mode?: 'text' | 'outlined' | 'contained' | 'elevated' | 'contained-tonal';
  /** Whether button is disabled */
  disabled?: boolean;
  /** Whether button shows loading state */
  loading?: boolean;
  /** Icon to show on the left */
  icon?: string;
  /** Press handler */
  onPress?: () => void;
  /** Additional styles */
  style?: ViewStyle;
  /** Whether button is compact */
  compact?: boolean;
  /** Whether button should fill container width */
  fullWidth?: boolean;
}

export function Button({
  children,
  mode = 'contained',
  disabled = false,
  loading = false,
  icon,
  onPress,
  style,
  compact = false,
  fullWidth = false,
}: ButtonProps): React.JSX.Element {
  const theme = useTheme<MD3Theme>();

  return (
    <PaperButton
      mode={mode}
      disabled={disabled}
      loading={loading}
      icon={icon}
      onPress={onPress}
      compact={compact}
      style={[
        fullWidth && styles.fullWidth,
        style,
      ]}
      contentStyle={styles.content}
      labelStyle={[
        styles.label,
        { fontFamily: theme.fonts.labelLarge.fontFamily },
      ]}
    >
      {children}
    </PaperButton>
  );
}

const styles = StyleSheet.create({
  fullWidth: {
    width: '100%',
  },
  content: {
    paddingVertical: 4,
  },
  label: {
    fontSize: 16,
    fontWeight: '600',
  },
});
