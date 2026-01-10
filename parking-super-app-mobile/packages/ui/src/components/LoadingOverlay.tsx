import React from 'react';
import { StyleSheet, View, Modal, ViewStyle } from 'react-native';
import { ActivityIndicator, Text, useTheme } from 'react-native-paper';
import type { MD3Theme } from 'react-native-paper';

export interface LoadingOverlayProps {
  /** Whether overlay is visible */
  visible: boolean;
  /** Loading message */
  message?: string;
  /** Whether overlay covers entire screen */
  fullScreen?: boolean;
  /** Background opacity */
  opacity?: number;
  /** Container style */
  style?: ViewStyle;
}

export function LoadingOverlay({
  visible,
  message,
  fullScreen = true,
  opacity = 0.6,
  style,
}: LoadingOverlayProps): React.JSX.Element | null {
  const theme = useTheme<MD3Theme>();

  if (!visible) {
    return null;
  }

  const content = (
    <View
      style={[
        styles.container,
        { backgroundColor: `rgba(0, 0, 0, ${opacity})` },
        style,
      ]}
    >
      <View
        style={[
          styles.content,
          { backgroundColor: theme.colors.surface },
        ]}
      >
        <ActivityIndicator
          size="large"
          color={theme.colors.primary}
          style={styles.indicator}
        />
        {message && (
          <Text variant="bodyMedium" style={styles.message}>
            {message}
          </Text>
        )}
      </View>
    </View>
  );

  if (fullScreen) {
    return (
      <Modal transparent visible={visible} animationType="fade">
        {content}
      </Modal>
    );
  }

  return content;
}

/**
 * Inline loading spinner
 */
export interface LoadingSpinnerProps {
  /** Spinner size */
  size?: 'small' | 'large';
  /** Custom color */
  color?: string;
  /** Optional message */
  message?: string;
  /** Container style */
  style?: ViewStyle;
}

export function LoadingSpinner({
  size = 'large',
  color,
  message,
  style,
}: LoadingSpinnerProps): React.JSX.Element {
  const theme = useTheme<MD3Theme>();

  return (
    <View style={[styles.spinnerContainer, style]}>
      <ActivityIndicator
        size={size}
        color={color ?? theme.colors.primary}
      />
      {message && (
        <Text
          variant="bodySmall"
          style={[styles.spinnerMessage, { color: theme.colors.onSurfaceVariant }]}
        >
          {message}
        </Text>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    ...StyleSheet.absoluteFillObject,
    justifyContent: 'center',
    alignItems: 'center',
  },
  content: {
    paddingHorizontal: 32,
    paddingVertical: 24,
    borderRadius: 16,
    alignItems: 'center',
    minWidth: 140,
  },
  indicator: {
    marginBottom: 8,
  },
  message: {
    textAlign: 'center',
    marginTop: 8,
  },
  spinnerContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 24,
  },
  spinnerMessage: {
    marginTop: 12,
    textAlign: 'center',
  },
});
