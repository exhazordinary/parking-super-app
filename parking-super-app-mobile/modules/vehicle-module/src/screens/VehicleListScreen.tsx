import React, { useCallback } from 'react';
import { View, StyleSheet, FlatList, Pressable } from 'react-native';
import { Text, useTheme, IconButton, FAB } from 'react-native-paper';
import { SafeAreaView } from 'react-native-safe-area-context';
import MaterialCommunityIcons from 'react-native-vector-icons/MaterialCommunityIcons';
import type { MD3Theme } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';

import { Card, LoadingSpinner, EmptyState, spacing } from '@parking/ui';
import { useVehicles, useDeleteVehicle, type Vehicle } from '@parking/api';
import type { VehicleStackScreenProps } from '@parking/navigation';

type NavigationProp = VehicleStackScreenProps<'VehicleList'>['navigation'];

function VehicleCard({ vehicle, onEdit, onDelete }: { vehicle: Vehicle; onEdit: () => void; onDelete: () => void }): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const typeIcons = { car: 'car', motorcycle: 'motorbike', truck: 'truck', van: 'van-utility' };

  return (
    <Card style={styles.vehicleCard}>
      <View style={styles.cardContent}>
        <View style={[styles.iconContainer, { backgroundColor: theme.colors.primaryContainer }]}>
          <MaterialCommunityIcons name={typeIcons[vehicle.type] || 'car'} size={28} color={theme.colors.primary} />
        </View>
        <View style={styles.vehicleInfo}>
          <View style={styles.plateRow}>
            <Text variant="titleMedium" style={{ fontWeight: '600' }}>{vehicle.plateNumber}</Text>
            {vehicle.isDefault && (
              <View style={[styles.defaultBadge, { backgroundColor: theme.colors.secondaryContainer }]}>
                <Text variant="labelSmall" style={{ color: theme.colors.secondary }}>Default</Text>
              </View>
            )}
          </View>
          {vehicle.make && <Text variant="bodyMedium" style={{ color: theme.colors.onSurfaceVariant }}>{vehicle.make} {vehicle.model}</Text>}
          {vehicle.color && <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>{vehicle.color}</Text>}
        </View>
        <View style={styles.actions}>
          <IconButton icon="pencil" size={20} onPress={onEdit} />
          <IconButton icon="delete" size={20} onPress={onDelete} iconColor={theme.colors.error} />
        </View>
      </View>
    </Card>
  );
}

export function VehicleListScreen(): React.JSX.Element {
  const navigation = useNavigation<NavigationProp>();
  const { data: vehicles, isLoading, refetch } = useVehicles();
  const deleteVehicle = useDeleteVehicle();

  const handleAddVehicle = useCallback(() => navigation.navigate('AddVehicle'), [navigation]);
  const handleEditVehicle = useCallback((id: string) => navigation.navigate('EditVehicle', { vehicleId: id }), [navigation]);
  const handleDeleteVehicle = useCallback(async (id: string) => {
    await deleteVehicle.mutateAsync(id);
  }, [deleteVehicle]);

  if (isLoading) return <SafeAreaView style={styles.safeArea}><LoadingSpinner message="Loading vehicles..." /></SafeAreaView>;

  return (
    <SafeAreaView style={styles.safeArea}>
      <View style={styles.header}>
        <Text variant="headlineMedium" style={{ fontWeight: '700' }}>My Vehicles</Text>
      </View>
      <FlatList
        data={vehicles}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => (
          <VehicleCard vehicle={item} onEdit={() => handleEditVehicle(item.id)} onDelete={() => handleDeleteVehicle(item.id)} />
        )}
        contentContainerStyle={styles.list}
        ListEmptyComponent={
          <EmptyState icon="car-off" title="No Vehicles" description="Add your first vehicle to start parking" actionText="Add Vehicle" onAction={handleAddVehicle} />
        }
      />
      <FAB icon="plus" style={styles.fab} onPress={handleAddVehicle} />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1, backgroundColor: '#F3F4F6' },
  header: { paddingHorizontal: spacing.md, paddingVertical: spacing.sm },
  list: { padding: spacing.md, paddingBottom: 80 },
  vehicleCard: { marginBottom: spacing.sm, padding: 0, overflow: 'hidden' },
  cardContent: { flexDirection: 'row', alignItems: 'center', padding: spacing.md },
  iconContainer: { width: 56, height: 56, borderRadius: 28, justifyContent: 'center', alignItems: 'center' },
  vehicleInfo: { flex: 1, marginLeft: spacing.md },
  plateRow: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm },
  defaultBadge: { paddingHorizontal: 8, paddingVertical: 2, borderRadius: 4 },
  actions: { flexDirection: 'row' },
  fab: { position: 'absolute', right: 16, bottom: 16 },
});
