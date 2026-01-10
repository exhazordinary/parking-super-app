import React, { useCallback } from 'react';
import { View, StyleSheet, FlatList, Pressable } from 'react-native';
import { Text, useTheme, IconButton } from 'react-native-paper';
import { SafeAreaView } from 'react-native-safe-area-context';
import MaterialCommunityIcons from 'react-native-vector-icons/MaterialCommunityIcons';
import type { MD3Theme } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';

import { Button, LoadingSpinner, EmptyState, spacing } from '@parking/ui';
import { useNotifications, useMarkAsRead, useMarkAllAsRead, type Notification } from '@parking/api';
import type { NotificationStackScreenProps } from '@parking/navigation';

type NavigationProp = NotificationStackScreenProps<'NotificationList'>['navigation'];

const typeIcons: Record<string, string> = {
  parking_reminder: 'clock-alert',
  parking_expiring: 'timer-sand',
  parking_ended: 'parking',
  payment_success: 'check-circle',
  payment_failed: 'alert-circle',
  topup_success: 'wallet-plus',
  low_balance: 'wallet-outline',
  promo: 'tag',
  system: 'information',
};

function NotificationItem({ notification, onPress }: { notification: Notification; onPress: () => void }): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const icon = typeIcons[notification.type] || 'bell';

  return (
    <Pressable onPress={onPress} style={[styles.notificationItem, !notification.isRead && { backgroundColor: theme.colors.primaryContainer + '30' }]}>
      <View style={[styles.iconContainer, { backgroundColor: theme.colors.primaryContainer }]}>
        <MaterialCommunityIcons name={icon} size={20} color={theme.colors.primary} />
      </View>
      <View style={styles.notificationContent}>
        <Text variant="bodyMedium" style={{ fontWeight: notification.isRead ? '400' : '600' }}>{notification.title}</Text>
        <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }} numberOfLines={2}>{notification.message}</Text>
        <Text variant="labelSmall" style={{ color: theme.colors.outline, marginTop: 4 }}>
          {new Date(notification.createdAt).toLocaleString()}
        </Text>
      </View>
      {!notification.isRead && <View style={[styles.unreadDot, { backgroundColor: theme.colors.primary }]} />}
    </Pressable>
  );
}

export function NotificationListScreen(): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const navigation = useNavigation<NavigationProp>();
  const { data, isLoading, fetchNextPage, hasNextPage, isFetchingNextPage } = useNotifications();
  const markAsRead = useMarkAsRead();
  const markAllAsRead = useMarkAllAsRead();

  const notifications = data?.pages.flatMap((p) => p.items) || [];
  const hasUnread = notifications.some((n) => !n.isRead);

  const handleNotificationPress = useCallback(async (notification: Notification) => {
    if (!notification.isRead) await markAsRead.mutateAsync([notification.id]);
    // Navigate based on notification type/data
    if (notification.data?.sessionId) {
      // @ts-expect-error - cross-module navigation
      navigation.navigate('Main', { screen: 'ParkingTab', params: { screen: 'SessionDetail', params: { sessionId: notification.data.sessionId } } });
    }
  }, [markAsRead, navigation]);

  const handleMarkAllRead = useCallback(() => markAllAsRead.mutate(), [markAllAsRead]);
  const handlePreferences = useCallback(() => navigation.navigate('NotificationPreferences'), [navigation]);
  const handleLoadMore = useCallback(() => { if (hasNextPage && !isFetchingNextPage) fetchNextPage(); }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

  if (isLoading) return <SafeAreaView style={styles.safeArea}><LoadingSpinner message="Loading notifications..." /></SafeAreaView>;

  return (
    <SafeAreaView style={styles.safeArea}>
      <View style={styles.header}>
        <Text variant="headlineMedium" style={{ fontWeight: '700' }}>Notifications</Text>
        <View style={styles.headerActions}>
          {hasUnread && <Button mode="text" onPress={handleMarkAllRead} compact loading={markAllAsRead.isPending}>Mark all read</Button>}
          <IconButton icon="cog" onPress={handlePreferences} />
        </View>
      </View>
      <FlatList
        data={notifications}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => <NotificationItem notification={item} onPress={() => handleNotificationPress(item)} />}
        onEndReached={handleLoadMore}
        contentContainerStyle={notifications.length === 0 ? { flex: 1 } : undefined}
        ListEmptyComponent={<EmptyState icon="bell-off" title="No Notifications" description="You're all caught up!" />}
        ListFooterComponent={isFetchingNextPage ? <LoadingSpinner size="small" /> : null}
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1, backgroundColor: '#FFFFFF' },
  header: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', paddingHorizontal: spacing.md, paddingVertical: spacing.sm },
  headerActions: { flexDirection: 'row', alignItems: 'center' },
  notificationItem: { flexDirection: 'row', alignItems: 'flex-start', padding: spacing.md, borderBottomWidth: 1, borderBottomColor: '#F3F4F6' },
  iconContainer: { width: 40, height: 40, borderRadius: 20, justifyContent: 'center', alignItems: 'center' },
  notificationContent: { flex: 1, marginLeft: spacing.md },
  unreadDot: { width: 8, height: 8, borderRadius: 4, marginTop: 6 },
});
