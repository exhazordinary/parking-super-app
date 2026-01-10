import React, { useState, useCallback } from 'react';
import { View, StyleSheet, ScrollView, Pressable } from 'react-native';
import { Text, useTheme, RadioButton } from 'react-native-paper';
import { SafeAreaView } from 'react-native-safe-area-context';
import MaterialCommunityIcons from 'react-native-vector-icons/MaterialCommunityIcons';
import type { MD3Theme } from 'react-native-paper';
import { useNavigation, useRoute } from '@react-navigation/native';

import { Button, Card, LoadingOverlay, LoadingSpinner, ErrorState, spacing } from '@parking/ui';
import { useLocation, useVehicles, useStartSession, type Vehicle } from '@parking/api';
import type { ParkingStackScreenProps } from '@parking/navigation';

type RouteProp = ParkingStackScreenProps<'StartSession'>['route'];
type NavigationProp = ParkingStackScreenProps<'StartSession'>['navigation'];

export function StartSessionScreen(): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const navigation = useNavigation<NavigationProp>();
  const route = useRoute<RouteProp>();
  const { locationId, vehicleId: preselectedVehicleId } = route.params;

  const { data: location, isLoading: locationLoading, error: locationError } = useLocation(locationId);
  const { data: vehicles, isLoading: vehiclesLoading } = useVehicles();
  const startSession = useStartSession();

  const [selectedVehicle, setSelectedVehicle] = useState<string>(preselectedVehicleId || '');
  const [error, setError] = useState('');

  const handleStartSession = useCallback(async () => {
    if (!selectedVehicle) {
      setError('Please select a vehicle');
      return;
    }

    try {
      const result = await startSession.mutateAsync({ vehicleId: selectedVehicle, locationId });
      navigation.replace('ActiveSession', { sessionId: result.session.id });
    } catch (err) {
      setError('Failed to start session. Please try again.');
    }
  }, [selectedVehicle, locationId, startSession, navigation]);

  if (locationLoading || vehiclesLoading) {
    return <SafeAreaView style={styles.safeArea}><LoadingSpinner message="Loading..." /></SafeAreaView>;
  }

  if (locationError || !location) {
    return <SafeAreaView style={styles.safeArea}><ErrorState message="Failed to load location" /></SafeAreaView>;
  }

  return (
    <SafeAreaView style={styles.safeArea}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        {/* Location Info */}
        <Card style={styles.locationCard}>
          <Text variant="titleLarge" style={{ fontWeight: '600' }}>{location.name}</Text>
          <Text variant="bodyMedium" style={{ color: theme.colors.onSurfaceVariant, marginTop: spacing.xs }}>
            {location.address}
          </Text>
          <View style={styles.rateRow}>
            <MaterialCommunityIcons name="currency-usd" size={20} color={theme.colors.primary} />
            <Text variant="titleMedium" style={{ marginLeft: spacing.xs }}>
              RM {location.rates[0]?.baseRate.toFixed(2)}/hour
            </Text>
          </View>
        </Card>

        {/* Vehicle Selection */}
        <Text variant="titleMedium" style={styles.sectionTitle}>Select Vehicle</Text>
        <Card style={styles.vehicleCard}>
          {vehicles && vehicles.length > 0 ? (
            <RadioButton.Group onValueChange={setSelectedVehicle} value={selectedVehicle}>
              {vehicles.map((vehicle: Vehicle) => (
                <Pressable
                  key={vehicle.id}
                  onPress={() => setSelectedVehicle(vehicle.id)}
                  style={[styles.vehicleItem, selectedVehicle === vehicle.id && { backgroundColor: theme.colors.primaryContainer }]}
                >
                  <MaterialCommunityIcons name="car" size={24} color={theme.colors.primary} />
                  <View style={styles.vehicleInfo}>
                    <Text variant="bodyLarge" style={{ fontWeight: '500' }}>{vehicle.plateNumber}</Text>
                    {vehicle.make && <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>{vehicle.make} {vehicle.model}</Text>}
                  </View>
                  <RadioButton value={vehicle.id} />
                </Pressable>
              ))}
            </RadioButton.Group>
          ) : (
            <View style={styles.noVehicle}>
              <MaterialCommunityIcons name="car-off" size={48} color={theme.colors.outline} />
              <Text variant="bodyMedium" style={{ color: theme.colors.onSurfaceVariant, marginTop: spacing.md }}>
                No vehicles registered
              </Text>
              <Button mode="outlined" style={{ marginTop: spacing.md }} onPress={() => navigation.navigate('VehicleTab' as any)}>
                Add Vehicle
              </Button>
            </View>
          )}
        </Card>

        {error && <Text variant="bodySmall" style={{ color: theme.colors.error, textAlign: 'center', marginTop: spacing.sm }}>{error}</Text>}

        <Button
          mode="contained"
          onPress={handleStartSession}
          loading={startSession.isPending}
          disabled={!selectedVehicle || startSession.isPending}
          fullWidth
          style={styles.startButton}
        >
          Start Parking
        </Button>
      </ScrollView>

      <LoadingOverlay visible={startSession.isPending} message="Starting session..." />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1, backgroundColor: '#F3F4F6' },
  scrollContent: { padding: spacing.md },
  locationCard: { padding: spacing.md, marginBottom: spacing.lg },
  rateRow: { flexDirection: 'row', alignItems: 'center', marginTop: spacing.md },
  sectionTitle: { fontWeight: '600', marginBottom: spacing.sm },
  vehicleCard: { padding: 0, overflow: 'hidden', marginBottom: spacing.lg },
  vehicleItem: { flexDirection: 'row', alignItems: 'center', padding: spacing.md, borderBottomWidth: 1, borderBottomColor: '#E5E7EB' },
  vehicleInfo: { flex: 1, marginLeft: spacing.md },
  noVehicle: { alignItems: 'center', padding: spacing.xl },
  startButton: { marginTop: spacing.md },
});
