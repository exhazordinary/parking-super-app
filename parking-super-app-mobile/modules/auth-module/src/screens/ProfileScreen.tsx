import React, { useCallback } from 'react';
import { View, StyleSheet, ScrollView, Pressable } from 'react-native';
import { Text, useTheme, Avatar, Divider } from 'react-native-paper';
import { SafeAreaView } from 'react-native-safe-area-context';
import MaterialCommunityIcons from 'react-native-vector-icons/MaterialCommunityIcons';
import type { MD3Theme } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';

import { Card, LoadingSpinner, ErrorState, spacing } from '@parking/ui';
import { useProfile, useUnreadCount } from '@parking/api';
import { useAuth } from '@parking/auth';
import type { AuthStackScreenProps } from '@parking/navigation';

type NavigationProp = AuthStackScreenProps<'Profile'>['navigation'];

interface MenuItemProps {
  icon: string;
  label: string;
  onPress: () => void;
  badge?: number;
  showArrow?: boolean;
  color?: string;
}

function MenuItem({
  icon,
  label,
  onPress,
  badge,
  showArrow = true,
  color,
}: MenuItemProps): React.JSX.Element {
  const theme = useTheme<MD3Theme>();

  return (
    <Pressable
      onPress={onPress}
      style={({ pressed }) => [
        styles.menuItem,
        pressed && { backgroundColor: theme.colors.surfaceVariant },
      ]}
    >
      <View style={styles.menuItemLeft}>
        <MaterialCommunityIcons
          name={icon}
          size={24}
          color={color ?? theme.colors.onSurfaceVariant}
        />
        <Text
          variant="bodyLarge"
          style={[styles.menuItemLabel, color && { color }]}
        >
          {label}
        </Text>
      </View>
      <View style={styles.menuItemRight}>
        {badge !== undefined && badge > 0 && (
          <View
            style={[
              styles.badge,
              { backgroundColor: theme.colors.error },
            ]}
          >
            <Text style={styles.badgeText}>
              {badge > 99 ? '99+' : badge}
            </Text>
          </View>
        )}
        {showArrow && (
          <MaterialCommunityIcons
            name="chevron-right"
            size={24}
            color={theme.colors.outline}
          />
        )}
      </View>
    </Pressable>
  );
}

export function ProfileScreen(): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const navigation = useNavigation<NavigationProp>();
  const { logout, isLoggingOut } = useAuth();

  const { data: profile, isLoading, error, refetch } = useProfile();
  const { data: unreadData } = useUnreadCount();

  const handleEditProfile = useCallback(() => {
    navigation.navigate('EditProfile');
  }, [navigation]);

  const handleChangePassword = useCallback(() => {
    navigation.navigate('ChangePassword');
  }, [navigation]);

  const handleSettings = useCallback(() => {
    navigation.navigate('Settings');
  }, [navigation]);

  const handleNotifications = useCallback(() => {
    // Navigate to notifications (different stack)
    // @ts-expect-error - cross-stack navigation
    navigation.navigate('Notifications', { screen: 'NotificationList' });
  }, [navigation]);

  const handleLogout = useCallback(async () => {
    await logout();
  }, [logout]);

  if (isLoading) {
    return (
      <SafeAreaView style={styles.safeArea}>
        <LoadingSpinner message="Loading profile..." />
      </SafeAreaView>
    );
  }

  if (error || !profile) {
    return (
      <SafeAreaView style={styles.safeArea}>
        <ErrorState
          message="Failed to load profile"
          onRetry={refetch}
        />
      </SafeAreaView>
    );
  }

  const initials = profile.name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);

  return (
    <SafeAreaView style={styles.safeArea}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        {/* Profile Header */}
        <View style={styles.profileHeader}>
          {profile.avatar ? (
            <Avatar.Image size={80} source={{ uri: profile.avatar }} />
          ) : (
            <Avatar.Text
              size={80}
              label={initials}
              style={{ backgroundColor: theme.colors.primaryContainer }}
            />
          )}
          <Text variant="headlineSmall" style={styles.profileName}>
            {profile.name}
          </Text>
          <Text
            variant="bodyMedium"
            style={{ color: theme.colors.onSurfaceVariant }}
          >
            {profile.phone}
          </Text>
          {profile.email && (
            <Text
              variant="bodySmall"
              style={{ color: theme.colors.onSurfaceVariant }}
            >
              {profile.email}
            </Text>
          )}
        </View>

        {/* Account Section */}
        <Card style={styles.section}>
          <Text
            variant="labelLarge"
            style={[styles.sectionTitle, { color: theme.colors.onSurfaceVariant }]}
          >
            Account
          </Text>
          <MenuItem
            icon="account-edit-outline"
            label="Edit Profile"
            onPress={handleEditProfile}
          />
          <Divider />
          <MenuItem
            icon="lock-outline"
            label="Change Password"
            onPress={handleChangePassword}
          />
          <Divider />
          <MenuItem
            icon="bell-outline"
            label="Notifications"
            onPress={handleNotifications}
            badge={unreadData?.count}
          />
        </Card>

        {/* Preferences Section */}
        <Card style={styles.section}>
          <Text
            variant="labelLarge"
            style={[styles.sectionTitle, { color: theme.colors.onSurfaceVariant }]}
          >
            Preferences
          </Text>
          <MenuItem
            icon="cog-outline"
            label="Settings"
            onPress={handleSettings}
          />
        </Card>

        {/* Logout */}
        <Card style={styles.section}>
          <MenuItem
            icon="logout"
            label={isLoggingOut ? 'Signing out...' : 'Sign Out'}
            onPress={handleLogout}
            showArrow={false}
            color={theme.colors.error}
          />
        </Card>

        {/* App Version */}
        <Text
          variant="bodySmall"
          style={[styles.version, { color: theme.colors.onSurfaceVariant }]}
        >
          Version 1.0.0
        </Text>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
    backgroundColor: '#F3F4F6',
  },
  scrollContent: {
    padding: spacing.md,
  },
  profileHeader: {
    alignItems: 'center',
    paddingVertical: spacing.xl,
  },
  profileName: {
    fontWeight: '600',
    marginTop: spacing.md,
    marginBottom: spacing.xs,
  },
  section: {
    marginBottom: spacing.md,
    padding: 0,
    overflow: 'hidden',
  },
  sectionTitle: {
    paddingHorizontal: spacing.md,
    paddingTop: spacing.md,
    paddingBottom: spacing.sm,
  },
  menuItem: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.md,
  },
  menuItemLeft: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
  },
  menuItemLabel: {
    fontWeight: '500',
  },
  menuItemRight: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.sm,
  },
  badge: {
    paddingHorizontal: 6,
    paddingVertical: 2,
    borderRadius: 10,
    minWidth: 20,
    alignItems: 'center',
  },
  badgeText: {
    color: '#FFFFFF',
    fontSize: 12,
    fontWeight: '600',
  },
  version: {
    textAlign: 'center',
    paddingVertical: spacing.lg,
  },
});
