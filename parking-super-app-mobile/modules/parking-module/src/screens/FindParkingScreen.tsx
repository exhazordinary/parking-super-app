import React, { useState, useCallback, useEffect } from 'react';
import { View, StyleSheet, FlatList, Pressable, PermissionsAndroid, Platform } from 'react-native';
import { Text, useTheme, Searchbar } from 'react-native-paper';
import { SafeAreaView } from 'react-native-safe-area-context';
import MaterialCommunityIcons from 'react-native-vector-icons/MaterialCommunityIcons';
import type { MD3Theme } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';
import Geolocation from '@react-native-community/geolocation';

import { Card, LoadingSpinner, EmptyState, spacing } from '@parking/ui';
import { useSearchLocations, type ParkingLocation } from '@parking/api';
import type { ParkingStackScreenProps } from '@parking/navigation';

type NavigationProp = ParkingStackScreenProps<'FindParking'>['navigation'];

function LocationCard({ location, onPress }: { location: ParkingLocation; onPress: () => void }): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const statusColor = location.status === 'open' ? '#16A34A' : location.status === 'full' ? '#DC2626' : '#D97706';
  const distance = location.distance ? `${(location.distance / 1000).toFixed(1)} km` : '';

  return (
    <Pressable onPress={onPress}>
      <Card style={styles.locationCard}>
        <View style={styles.locationHeader}>
          <View style={styles.locationInfo}>
            <Text variant="titleSmall" style={{ fontWeight: '600' }}>{location.name}</Text>
            <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>{location.address}</Text>
          </View>
          {distance && (
            <View style={styles.distanceBadge}>
              <MaterialCommunityIcons name="map-marker" size={14} color={theme.colors.primary} />
              <Text variant="labelSmall" style={{ color: theme.colors.primary }}>{distance}</Text>
            </View>
          )}
        </View>
        <View style={styles.locationFooter}>
          <View style={styles.statusContainer}>
            <View style={[styles.statusDot, { backgroundColor: statusColor }]} />
            <Text variant="bodySmall" style={{ textTransform: 'capitalize' }}>{location.status}</Text>
            <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>
              {' '}• {location.availableSpaces}/{location.totalSpaces} spots
            </Text>
          </View>
          <Text variant="labelLarge" style={{ color: theme.colors.primary, fontWeight: '600' }}>
            RM {location.rates[0]?.baseRate.toFixed(2)}/hr
          </Text>
        </View>
      </Card>
    </Pressable>
  );
}

export function FindParkingScreen(): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const navigation = useNavigation<NavigationProp>();

  const [searchQuery, setSearchQuery] = useState('');
  const [userLocation, setUserLocation] = useState<{ latitude: number; longitude: number } | null>(null);

  const { data: locations, isLoading, error, refetch } = useSearchLocations({
    latitude: userLocation?.latitude || 3.139, // Default to KL
    longitude: userLocation?.longitude || 101.6869,
    radius: 5000,
  });

  useEffect(() => {
    const requestLocation = async () => {
      if (Platform.OS === 'android') {
        const granted = await PermissionsAndroid.request(PermissionsAndroid.PERMISSIONS.ACCESS_FINE_LOCATION);
        if (granted !== PermissionsAndroid.RESULTS.GRANTED) return;
      }
      Geolocation.getCurrentPosition(
        (position) => setUserLocation({ latitude: position.coords.latitude, longitude: position.coords.longitude }),
        (err) => console.warn('Location error:', err),
        { enableHighAccuracy: true, timeout: 15000 }
      );
    };
    requestLocation();
  }, []);

  const handleLocationPress = useCallback((id: string) => {
    navigation.navigate('LocationDetail', { locationId: id });
  }, [navigation]);

  const handleMapView = useCallback(() => {
    navigation.navigate('ParkingMap', userLocation || undefined);
  }, [navigation, userLocation]);

  const filteredLocations = locations?.filter((loc) =>
    loc.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    loc.address.toLowerCase().includes(searchQuery.toLowerCase())
  ) || [];

  return (
    <SafeAreaView style={styles.safeArea}>
      <View style={styles.header}>
        <Searchbar
          placeholder="Search parking locations"
          value={searchQuery}
          onChangeText={setSearchQuery}
          style={styles.searchbar}
        />
        <Pressable onPress={handleMapView} style={styles.mapButton}>
          <MaterialCommunityIcons name="map" size={24} color={theme.colors.primary} />
        </Pressable>
      </View>

      {isLoading ? (
        <LoadingSpinner message="Finding nearby parking..." />
      ) : (
        <FlatList
          data={filteredLocations}
          keyExtractor={(item) => item.id}
          renderItem={({ item }) => <LocationCard location={item} onPress={() => handleLocationPress(item.id)} />}
          contentContainerStyle={styles.list}
          ListEmptyComponent={<EmptyState icon="parking" title="No Parking Found" description="Try searching in a different area" />}
        />
      )}
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1, backgroundColor: '#FFFFFF' },
  header: { flexDirection: 'row', alignItems: 'center', padding: spacing.md, gap: spacing.sm },
  searchbar: { flex: 1, elevation: 0, backgroundColor: '#F3F4F6' },
  mapButton: { padding: spacing.sm, backgroundColor: '#F3F4F6', borderRadius: 8 },
  list: { padding: spacing.md },
  locationCard: { padding: spacing.md, marginBottom: spacing.sm },
  locationHeader: { flexDirection: 'row', justifyContent: 'space-between' },
  locationInfo: { flex: 1 },
  distanceBadge: { flexDirection: 'row', alignItems: 'center', gap: 2 },
  locationFooter: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginTop: spacing.sm },
  statusContainer: { flexDirection: 'row', alignItems: 'center' },
  statusDot: { width: 8, height: 8, borderRadius: 4, marginRight: spacing.xs },
});
