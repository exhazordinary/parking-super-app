import React, { ReactNode } from 'react';
import { StyleSheet, ViewStyle, View } from 'react-native';
import { Modal as PaperModal, Portal, Text, useTheme } from 'react-native-paper';
import type { MD3Theme } from 'react-native-paper';
import { Button } from './Button';

export interface ModalProps {
  /** Whether modal is visible */
  visible: boolean;
  /** Dismiss handler */
  onDismiss: () => void;
  /** Modal title */
  title?: string;
  /** Modal content */
  children: ReactNode;
  /** Whether modal can be dismissed by tapping outside */
  dismissable?: boolean;
  /** Additional content styles */
  contentStyle?: ViewStyle;
}

export function Modal({
  visible,
  onDismiss,
  title,
  children,
  dismissable = true,
  contentStyle,
}: ModalProps): React.JSX.Element {
  const theme = useTheme<MD3Theme>();

  return (
    <Portal>
      <PaperModal
        visible={visible}
        onDismiss={onDismiss}
        dismissable={dismissable}
        contentContainerStyle={[
          styles.container,
          { backgroundColor: theme.colors.surface },
          contentStyle,
        ]}
      >
        {title && (
          <Text variant="headlineSmall\" style={styles.title}>
            {title}
          </Text>
        )}
        <View style={styles.content}>
          {children}
        </View>
      </PaperModal>
    </Portal>
  );
}

export interface ConfirmModalProps {
  /** Whether modal is visible */
  visible: boolean;
  /** Dismiss handler */
  onDismiss: () => void;
  /** Confirm handler */
  onConfirm: () => void;
  /** Modal title */
  title: string;
  /** Modal message */
  message: string;
  /** Confirm button text */
  confirmText?: string;
  /** Cancel button text */
  cancelText?: string;
  /** Whether action is destructive */
  destructive?: boolean;
  /** Loading state */
  loading?: boolean;
}

export function ConfirmModal({
  visible,
  onDismiss,
  onConfirm,
  title,
  message,
  confirmText = 'Confirm',
  cancelText = 'Cancel',
  destructive = false,
  loading = false,
}: ConfirmModalProps): React.JSX.Element {
  const theme = useTheme<MD3Theme>();

  return (
    <Modal visible={visible} onDismiss={onDismiss} title={title}>
      <Text variant="bodyMedium" style={styles.message}>
        {message}
      </Text>
      <View style={styles.actions}>
        <Button
          mode="text"
          onPress={onDismiss}
          disabled={loading}
          style={styles.actionButton}
        >
          {cancelText}
        </Button>
        <Button
          mode="contained"
          onPress={onConfirm}
          loading={loading}
          style={[
            styles.actionButton,
            destructive && { backgroundColor: theme.colors.error },
          ]}
        >
          {confirmText}
        </Button>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  container: {
    margin: 20,
    padding: 20,
    borderRadius: 16,
  },
  title: {
    marginBottom: 16,
    fontWeight: '600',
  },
  content: {
    minHeight: 40,
  },
  message: {
    marginBottom: 24,
    lineHeight: 22,
  },
  actions: {
    flexDirection: 'row',
    justifyContent: 'flex-end',
    gap: 12,
  },
  actionButton: {
    minWidth: 80,
  },
});
