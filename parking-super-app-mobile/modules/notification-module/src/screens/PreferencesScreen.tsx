import React, { useCallback } from 'react';
import { View, StyleSheet, ScrollView } from 'react-native';
import { Text, useTheme, Switch, Divider } from 'react-native-paper';
import { SafeAreaView } from 'react-native-safe-area-context';
import type { MD3Theme } from 'react-native-paper';

import { Card, LoadingSpinner, ErrorState, spacing } from '@parking/ui';
import { useNotificationPreferences, useUpdateNotificationPreferences } from '@parking/api';

function PreferenceRow({ label, description, value, onToggle, disabled }: { label: string; description: string; value: boolean; onToggle: (v: boolean) => void; disabled?: boolean }): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  return (
    <View style={styles.preferenceRow}>
      <View style={styles.preferenceInfo}>
        <Text variant="bodyLarge" style={{ fontWeight: '500' }}>{label}</Text>
        <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>{description}</Text>
      </View>
      <Switch value={value} onValueChange={onToggle} disabled={disabled} />
    </View>
  );
}

export function PreferencesScreen(): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const { data: preferences, isLoading, error, refetch } = useNotificationPreferences();
  const updatePreferences = useUpdateNotificationPreferences();

  const handleToggle = useCallback((key: string, value: boolean) => {
    updatePreferences.mutate({ [key]: value });
  }, [updatePreferences]);

  if (isLoading) return <SafeAreaView style={styles.safeArea}><LoadingSpinner message="Loading preferences..." /></SafeAreaView>;
  if (error || !preferences) return <SafeAreaView style={styles.safeArea}><ErrorState message="Failed to load preferences" onRetry={refetch} /></SafeAreaView>;

  return (
    <SafeAreaView style={styles.safeArea}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <Text variant="headlineSmall" style={styles.title}>Notification Preferences</Text>

        {/* Channels */}
        <Card style={styles.section}>
          <Text variant="labelLarge" style={[styles.sectionTitle, { color: theme.colors.onSurfaceVariant }]}>Channels</Text>
          <PreferenceRow
            label="Push Notifications"
            description="Receive notifications on your device"
            value={preferences.pushEnabled}
            onToggle={(v) => handleToggle('pushEnabled', v)}
            disabled={updatePreferences.isPending}
          />
          <Divider />
          <PreferenceRow
            label="Email Notifications"
            description="Receive notifications via email"
            value={preferences.emailEnabled}
            onToggle={(v) => handleToggle('emailEnabled', v)}
            disabled={updatePreferences.isPending}
          />
          <Divider />
          <PreferenceRow
            label="SMS Notifications"
            description="Receive notifications via SMS"
            value={preferences.smsEnabled}
            onToggle={(v) => handleToggle('smsEnabled', v)}
            disabled={updatePreferences.isPending}
          />
        </Card>

        {/* Categories */}
        <Card style={styles.section}>
          <Text variant="labelLarge" style={[styles.sectionTitle, { color: theme.colors.onSurfaceVariant }]}>Categories</Text>
          <PreferenceRow
            label="Parking Reminders"
            description="Get reminded before your session expires"
            value={preferences.parkingReminders}
            onToggle={(v) => handleToggle('parkingReminders', v)}
            disabled={updatePreferences.isPending}
          />
          <Divider />
          <PreferenceRow
            label="Promotions"
            description="Receive special offers and discounts"
            value={preferences.promotions}
            onToggle={(v) => handleToggle('promotions', v)}
            disabled={updatePreferences.isPending}
          />
          <Divider />
          <PreferenceRow
            label="System Updates"
            description="Important app and service updates"
            value={preferences.systemUpdates}
            onToggle={(v) => handleToggle('systemUpdates', v)}
            disabled={updatePreferences.isPending}
          />
        </Card>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1, backgroundColor: '#F3F4F6' },
  scrollContent: { padding: spacing.md },
  title: { fontWeight: '600', marginBottom: spacing.lg },
  section: { marginBottom: spacing.md, padding: 0, overflow: 'hidden' },
  sectionTitle: { padding: spacing.md, paddingBottom: spacing.sm },
  preferenceRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', padding: spacing.md },
  preferenceInfo: { flex: 1, marginRight: spacing.md },
});
