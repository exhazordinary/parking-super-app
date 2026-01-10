import React from 'react';
import { StyleSheet, View, ViewStyle } from 'react-native';
import { Text, useTheme } from 'react-native-paper';
import MaterialCommunityIcons from 'react-native-vector-icons/MaterialCommunityIcons';
import type { MD3Theme } from 'react-native-paper';
import { Button } from './Button';

export interface EmptyStateProps {
  /** Icon name from MaterialCommunityIcons */
  icon: string;
  /** Title text */
  title: string;
  /** Description text */
  description?: string;
  /** Action button text */
  actionText?: string;
  /** Action handler */
  onAction?: () => void;
  /** Icon size */
  iconSize?: number;
  /** Container style */
  style?: ViewStyle;
}

export function EmptyState({
  icon,
  title,
  description,
  actionText,
  onAction,
  iconSize = 64,
  style,
}: EmptyStateProps): React.JSX.Element {
  const theme = useTheme<MD3Theme>();

  return (
    <View style={[styles.container, style]}>
      <MaterialCommunityIcons
        name={icon}
        size={iconSize}
        color={theme.colors.outline}
        style={styles.icon}
      />
      <Text
        variant="titleMedium"
        style={[styles.title, { color: theme.colors.onSurface }]}
      >
        {title}
      </Text>
      {description && (
        <Text
          variant="bodyMedium"
          style={[styles.description, { color: theme.colors.onSurfaceVariant }]}
        >
          {description}
        </Text>
      )}
      {actionText && onAction && (
        <Button
          mode="contained"
          onPress={onAction}
          style={styles.action}
        >
          {actionText}
        </Button>
      )}
    </View>
  );
}

/**
 * Error state component
 */
export interface ErrorStateProps {
  /** Error title */
  title?: string;
  /** Error message */
  message: string;
  /** Retry handler */
  onRetry?: () => void;
  /** Retry button text */
  retryText?: string;
  /** Container style */
  style?: ViewStyle;
}

export function ErrorState({
  title = 'Something went wrong',
  message,
  onRetry,
  retryText = 'Try Again',
  style,
}: ErrorStateProps): React.JSX.Element {
  const theme = useTheme<MD3Theme>();

  return (
    <View style={[styles.container, style]}>
      <MaterialCommunityIcons
        name="alert-circle-outline"
        size={64}
        color={theme.colors.error}
        style={styles.icon}
      />
      <Text
        variant="titleMedium"
        style={[styles.title, { color: theme.colors.onSurface }]}
      >
        {title}
      </Text>
      <Text
        variant="bodyMedium"
        style={[styles.description, { color: theme.colors.onSurfaceVariant }]}
      >
        {message}
      </Text>
      {onRetry && (
        <Button
          mode="outlined"
          onPress={onRetry}
          icon="refresh"
          style={styles.action}
        >
          {retryText}
        </Button>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 32,
  },
  icon: {
    marginBottom: 16,
  },
  title: {
    fontWeight: '600',
    textAlign: 'center',
    marginBottom: 8,
  },
  description: {
    textAlign: 'center',
    lineHeight: 22,
    marginBottom: 24,
  },
  action: {
    marginTop: 8,
  },
});
