import React, { useCallback } from 'react';
import { View, StyleSheet, ScrollView, Linking, Pressable } from 'react-native';
import { Text, useTheme, Avatar, Divider } from 'react-native-paper';
import { SafeAreaView } from 'react-native-safe-area-context';
import MaterialCommunityIcons from 'react-native-vector-icons/MaterialCommunityIcons';
import type { MD3Theme } from 'react-native-paper';
import { useNavigation, useRoute } from '@react-navigation/native';

import { Button, Card, LoadingSpinner, ErrorState, spacing } from '@parking/ui';
import { useProvider, useProviderLocations, type ParkingLocation } from '@parking/api';
import type { ProviderStackScreenProps } from '@parking/navigation';

type RouteProp = ProviderStackScreenProps<'ProviderDetail'>['route'];
type NavigationProp = ProviderStackScreenProps<'ProviderDetail'>['navigation'];

function LocationItem({ location, onPress }: { location: ParkingLocation; onPress: () => void }): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const statusColor = location.status === 'open' ? '#16A34A' : location.status === 'full' ? '#DC2626' : '#D97706';

  return (
    <Pressable onPress={onPress} style={styles.locationItem}>
      <View style={styles.locationInfo}>
        <Text variant="bodyLarge" style={{ fontWeight: '500' }}>{location.name}</Text>
        <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>{location.address}</Text>
        <View style={styles.locationStats}>
          <View style={[styles.statusDot, { backgroundColor: statusColor }]} />
          <Text variant="labelSmall">{location.status} • {location.availableSpaces} spots</Text>
        </View>
      </View>
      <Text variant="labelLarge" style={{ color: theme.colors.primary }}>RM {location.rates[0]?.baseRate.toFixed(2)}/hr</Text>
    </Pressable>
  );
}

export function ProviderDetailScreen(): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const navigation = useNavigation<NavigationProp>();
  const route = useRoute<RouteProp>();
  const { providerId } = route.params;

  const { data: provider, isLoading, error, refetch } = useProvider(providerId);
  const { data: locationsData } = useProviderLocations(providerId, { limit: 5 });
  const locations = locationsData?.pages[0]?.items || [];

  const handleContact = useCallback((type: 'phone' | 'email' | 'website') => {
    if (!provider) return;
    if (type === 'phone' && provider.phone) Linking.openURL(`tel:${provider.phone}`);
    else if (type === 'email' && provider.email) Linking.openURL(`mailto:${provider.email}`);
    else if (type === 'website' && provider.website) Linking.openURL(provider.website);
  }, [provider]);

  const handleViewLocations = useCallback(() => navigation.navigate('ProviderLocations', { providerId }), [navigation, providerId]);
  const handleLocationPress = useCallback((id: string) => {
    // Navigate to parking module location detail
    // @ts-expect-error - cross-module navigation
    navigation.navigate('Main', { screen: 'ParkingTab', params: { screen: 'LocationDetail', params: { locationId: id } } });
  }, [navigation]);

  if (isLoading) return <SafeAreaView style={styles.safeArea}><LoadingSpinner message="Loading provider..." /></SafeAreaView>;
  if (error || !provider) return <SafeAreaView style={styles.safeArea}><ErrorState message="Failed to load provider" onRetry={refetch} /></SafeAreaView>;

  return (
    <SafeAreaView style={styles.safeArea}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        {/* Header */}
        <View style={styles.header}>
          {provider.logo ? <Avatar.Image size={80} source={{ uri: provider.logo }} /> : <Avatar.Icon size={80} icon="parking" style={{ backgroundColor: theme.colors.primaryContainer }} />}
          <Text variant="headlineSmall" style={styles.providerName}>{provider.name}</Text>
          {provider.rating && (
            <View style={styles.ratingRow}>
              <MaterialCommunityIcons name="star" size={18} color="#F59E0B" />
              <Text variant="titleMedium" style={{ marginLeft: 4 }}>{provider.rating.toFixed(1)}</Text>
              <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant, marginLeft: 4 }}>({provider.reviewCount} reviews)</Text>
            </View>
          )}
        </View>

        {/* Description */}
        {provider.description && <Card style={styles.section}><Text variant="bodyMedium">{provider.description}</Text></Card>}

        {/* Contact */}
        <Card style={styles.section}>
          <Text variant="labelLarge" style={styles.sectionTitle}>Contact</Text>
          {provider.phone && (
            <Pressable onPress={() => handleContact('phone')} style={styles.contactRow}>
              <MaterialCommunityIcons name="phone" size={20} color={theme.colors.primary} />
              <Text variant="bodyMedium" style={styles.contactText}>{provider.phone}</Text>
            </Pressable>
          )}
          {provider.email && (
            <Pressable onPress={() => handleContact('email')} style={styles.contactRow}>
              <MaterialCommunityIcons name="email" size={20} color={theme.colors.primary} />
              <Text variant="bodyMedium" style={styles.contactText}>{provider.email}</Text>
            </Pressable>
          )}
          {provider.website && (
            <Pressable onPress={() => handleContact('website')} style={styles.contactRow}>
              <MaterialCommunityIcons name="web" size={20} color={theme.colors.primary} />
              <Text variant="bodyMedium" style={styles.contactText}>{provider.website}</Text>
            </Pressable>
          )}
        </Card>

        {/* Locations */}
        <View style={styles.locationsSection}>
          <View style={styles.sectionHeader}>
            <Text variant="titleMedium" style={{ fontWeight: '600' }}>Locations ({provider.totalLocations})</Text>
            <Button mode="text" onPress={handleViewLocations} compact>View All</Button>
          </View>
          <Card style={{ padding: 0 }}>
            {locations.map((loc, i) => (
              <React.Fragment key={loc.id}>
                <LocationItem location={loc} onPress={() => handleLocationPress(loc.id)} />
                {i < locations.length - 1 && <Divider />}
              </React.Fragment>
            ))}
          </Card>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1, backgroundColor: '#F3F4F6' },
  scrollContent: { padding: spacing.md },
  header: { alignItems: 'center', paddingVertical: spacing.lg },
  providerName: { fontWeight: '700', marginTop: spacing.md },
  ratingRow: { flexDirection: 'row', alignItems: 'center', marginTop: spacing.xs },
  section: { padding: spacing.md, marginBottom: spacing.md },
  sectionTitle: { marginBottom: spacing.sm },
  contactRow: { flexDirection: 'row', alignItems: 'center', paddingVertical: spacing.sm },
  contactText: { marginLeft: spacing.md },
  locationsSection: { marginBottom: spacing.lg },
  sectionHeader: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: spacing.sm },
  locationItem: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', padding: spacing.md },
  locationInfo: { flex: 1 },
  locationStats: { flexDirection: 'row', alignItems: 'center', marginTop: spacing.xs },
  statusDot: { width: 8, height: 8, borderRadius: 4, marginRight: spacing.xs },
});
