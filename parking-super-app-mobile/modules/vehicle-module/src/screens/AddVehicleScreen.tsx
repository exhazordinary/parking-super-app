import React, { useState, useCallback } from 'react';
import { View, StyleSheet, ScrollView } from 'react-native';
import { Text, useTheme, SegmentedButtons, Switch } from 'react-native-paper';
import { SafeAreaView } from 'react-native-safe-area-context';
import type { MD3Theme } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';

import { Button, Input, LoadingOverlay, spacing } from '@parking/ui';
import { useCreateVehicle, type VehicleType } from '@parking/api';
import type { VehicleStackScreenProps } from '@parking/navigation';

type NavigationProp = VehicleStackScreenProps<'AddVehicle'>['navigation'];

// Malaysian plate format: ABC 1234 or AB 1234 C
const validateMalaysianPlate = (plate: string): boolean => {
  const cleaned = plate.replace(/\s/g, '').toUpperCase();
  return /^[A-Z]{1,3}\d{1,4}[A-Z]?$/.test(cleaned);
};

const formatPlate = (plate: string): string => {
  const cleaned = plate.replace(/[^A-Za-z0-9]/g, '').toUpperCase();
  const letters = cleaned.match(/^[A-Z]+/)?.[0] || '';
  const numbers = cleaned.slice(letters.length).match(/^\d+/)?.[0] || '';
  const suffix = cleaned.slice(letters.length + numbers.length).match(/^[A-Z]$/)?.[0] || '';
  return `${letters} ${numbers}${suffix ? ' ' + suffix : ''}`.trim();
};

export function AddVehicleScreen(): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const navigation = useNavigation<NavigationProp>();
  const createVehicle = useCreateVehicle();

  const [plateNumber, setPlateNumber] = useState('');
  const [make, setMake] = useState('');
  const [model, setModel] = useState('');
  const [color, setColor] = useState('');
  const [vehicleType, setVehicleType] = useState<VehicleType>('car');
  const [isDefault, setIsDefault] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const handlePlateChange = useCallback((text: string) => {
    setPlateNumber(formatPlate(text));
    setErrors((e) => ({ ...e, plateNumber: '' }));
  }, []);

  const handleSubmit = useCallback(async () => {
    const newErrors: Record<string, string> = {};
    if (!plateNumber) newErrors.plateNumber = 'Plate number is required';
    else if (!validateMalaysianPlate(plateNumber)) newErrors.plateNumber = 'Invalid Malaysian plate format';

    setErrors(newErrors);
    if (Object.keys(newErrors).length > 0) return;

    try {
      await createVehicle.mutateAsync({
        plateNumber: plateNumber.replace(/\s/g, ''),
        make: make || undefined,
        model: model || undefined,
        color: color || undefined,
        type: vehicleType,
        isDefault,
      });
      navigation.goBack();
    } catch (err) {
      setErrors({ general: 'Failed to add vehicle. Please try again.' });
    }
  }, [plateNumber, make, model, color, vehicleType, isDefault, createVehicle, navigation]);

  return (
    <SafeAreaView style={styles.safeArea}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <Text variant="headlineSmall" style={styles.title}>Add Vehicle</Text>

        <Input
          label="Plate Number *"
          value={plateNumber}
          onChangeText={handlePlateChange}
          placeholder="ABC 1234"
          autoCapitalize="characters"
          error={Boolean(errors.plateNumber)}
          helperText={errors.plateNumber || 'Malaysian format: ABC 1234'}
        />

        <Text variant="labelLarge" style={styles.label}>Vehicle Type</Text>
        <SegmentedButtons
          value={vehicleType}
          onValueChange={(v) => setVehicleType(v as VehicleType)}
          buttons={[
            { value: 'car', label: 'Car', icon: 'car' },
            { value: 'motorcycle', label: 'Bike', icon: 'motorbike' },
            { value: 'truck', label: 'Truck', icon: 'truck' },
          ]}
          style={styles.segmented}
        />

        <Input label="Make (Optional)" value={make} onChangeText={setMake} placeholder="e.g., Toyota" />
        <Input label="Model (Optional)" value={model} onChangeText={setModel} placeholder="e.g., Camry" />
        <Input label="Color (Optional)" value={color} onChangeText={setColor} placeholder="e.g., White" />

        <View style={styles.switchRow}>
          <Text variant="bodyLarge">Set as default vehicle</Text>
          <Switch value={isDefault} onValueChange={setIsDefault} />
        </View>

        {errors.general && <Text variant="bodySmall" style={{ color: theme.colors.error, textAlign: 'center' }}>{errors.general}</Text>}

        <Button mode="contained" onPress={handleSubmit} loading={createVehicle.isPending} disabled={createVehicle.isPending} fullWidth style={styles.submitButton}>
          Add Vehicle
        </Button>
      </ScrollView>
      <LoadingOverlay visible={createVehicle.isPending} message="Adding vehicle..." />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1, backgroundColor: '#FFFFFF' },
  scrollContent: { padding: spacing.md },
  title: { fontWeight: '600', marginBottom: spacing.lg },
  label: { marginTop: spacing.md, marginBottom: spacing.sm },
  segmented: { marginBottom: spacing.md },
  switchRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', paddingVertical: spacing.md, borderTopWidth: 1, borderTopColor: '#E5E7EB', marginTop: spacing.md },
  submitButton: { marginTop: spacing.lg },
});
