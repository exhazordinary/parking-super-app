import React, { useState, useEffect, useCallback } from 'react';
import { View, StyleSheet, ScrollView } from 'react-native';
import { Text, useTheme } from 'react-native-paper';
import { SafeAreaView } from 'react-native-safe-area-context';
import MaterialCommunityIcons from 'react-native-vector-icons/MaterialCommunityIcons';
import type { MD3Theme } from 'react-native-paper';
import { useNavigation, useRoute } from '@react-navigation/native';

import { Button, Card, LoadingSpinner, ErrorState, ConfirmModal, LoadingOverlay, spacing } from '@parking/ui';
import { useSession, useEndSession } from '@parking/api';
import type { ParkingStackScreenProps } from '@parking/navigation';

type RouteProp = ParkingStackScreenProps<'ActiveSession'>['route'];
type NavigationProp = ParkingStackScreenProps<'ActiveSession'>['navigation'];

export function ActiveSessionScreen(): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const navigation = useNavigation<NavigationProp>();
  const route = useRoute<RouteProp>();
  const { sessionId } = route.params;

  const { data: session, isLoading, error, refetch } = useSession(sessionId);
  const endSession = useEndSession();

  const [duration, setDuration] = useState({ hours: 0, minutes: 0, seconds: 0 });
  const [showEndModal, setShowEndModal] = useState(false);

  // Live timer
  useEffect(() => {
    if (!session?.startTime) return;

    const updateDuration = () => {
      const start = new Date(session.startTime).getTime();
      const diff = Date.now() - start;
      const hours = Math.floor(diff / 3600000);
      const minutes = Math.floor((diff % 3600000) / 60000);
      const seconds = Math.floor((diff % 60000) / 1000);
      setDuration({ hours, minutes, seconds });
    };

    updateDuration();
    const interval = setInterval(updateDuration, 1000);
    return () => clearInterval(interval);
  }, [session?.startTime]);

  const calculateCost = useCallback(() => {
    if (!session?.rate) return 0;
    const totalMinutes = duration.hours * 60 + duration.minutes;
    const hours = Math.ceil(totalMinutes / 60);
    return hours * session.rate.baseRate;
  }, [session?.rate, duration]);

  const handleEndSession = useCallback(async () => {
    setShowEndModal(false);
    try {
      await endSession.mutateAsync({ sessionId });
      navigation.replace('SessionDetail', { sessionId });
    } catch (err) {
      // Error handled by mutation
    }
  }, [sessionId, endSession, navigation]);

  if (isLoading) return <SafeAreaView style={styles.safeArea}><LoadingSpinner message="Loading session..." /></SafeAreaView>;
  if (error || !session) return <SafeAreaView style={styles.safeArea}><ErrorState message="Failed to load session" onRetry={refetch} /></SafeAreaView>;

  const estimatedCost = calculateCost();

  return (
    <SafeAreaView style={styles.safeArea}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        {/* Timer Card */}
        <Card style={[styles.timerCard, { backgroundColor: theme.colors.primary }]}>
          <Text variant="titleMedium" style={{ color: 'rgba(255,255,255,0.8)' }}>Parking Duration</Text>
          <View style={styles.timerDisplay}>
            <View style={styles.timerUnit}>
              <Text style={styles.timerNumber}>{String(duration.hours).padStart(2, '0')}</Text>
              <Text style={styles.timerLabel}>HRS</Text>
            </View>
            <Text style={styles.timerSeparator}>:</Text>
            <View style={styles.timerUnit}>
              <Text style={styles.timerNumber}>{String(duration.minutes).padStart(2, '0')}</Text>
              <Text style={styles.timerLabel}>MIN</Text>
            </View>
            <Text style={styles.timerSeparator}>:</Text>
            <View style={styles.timerUnit}>
              <Text style={styles.timerNumber}>{String(duration.seconds).padStart(2, '0')}</Text>
              <Text style={styles.timerLabel}>SEC</Text>
            </View>
          </View>
        </Card>

        {/* Session Details */}
        <Card style={styles.detailsCard}>
          <View style={styles.detailRow}>
            <MaterialCommunityIcons name="map-marker" size={20} color={theme.colors.primary} />
            <View style={styles.detailInfo}>
              <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>Location</Text>
              <Text variant="bodyLarge" style={{ fontWeight: '500' }}>{session.location?.name || 'Parking'}</Text>
            </View>
          </View>

          <View style={styles.detailRow}>
            <MaterialCommunityIcons name="car" size={20} color={theme.colors.primary} />
            <View style={styles.detailInfo}>
              <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>Vehicle</Text>
              <Text variant="bodyLarge" style={{ fontWeight: '500' }}>{session.vehicle?.plateNumber}</Text>
            </View>
          </View>

          <View style={styles.detailRow}>
            <MaterialCommunityIcons name="clock-start" size={20} color={theme.colors.primary} />
            <View style={styles.detailInfo}>
              <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>Started At</Text>
              <Text variant="bodyLarge" style={{ fontWeight: '500' }}>{new Date(session.startTime).toLocaleTimeString()}</Text>
            </View>
          </View>

          <View style={[styles.detailRow, styles.costRow]}>
            <Text variant="titleMedium">Estimated Cost</Text>
            <Text variant="headlineSmall" style={{ fontWeight: '700', color: theme.colors.primary }}>
              RM {estimatedCost.toFixed(2)}
            </Text>
          </View>
        </Card>

        <Button
          mode="contained"
          onPress={() => setShowEndModal(true)}
          buttonColor={theme.colors.error}
          style={styles.endButton}
        >
          End Parking Session
        </Button>
      </ScrollView>

      <ConfirmModal
        visible={showEndModal}
        onDismiss={() => setShowEndModal(false)}
        onConfirm={handleEndSession}
        title="End Session?"
        message={`You will be charged RM ${estimatedCost.toFixed(2)} from your wallet.`}
        confirmText="End & Pay"
        destructive
        loading={endSession.isPending}
      />

      <LoadingOverlay visible={endSession.isPending} message="Ending session..." />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1, backgroundColor: '#F3F4F6' },
  scrollContent: { padding: spacing.md },
  timerCard: { padding: spacing.lg, alignItems: 'center', marginBottom: spacing.lg },
  timerDisplay: { flexDirection: 'row', alignItems: 'center', marginTop: spacing.md },
  timerUnit: { alignItems: 'center' },
  timerNumber: { fontSize: 48, fontWeight: '700', color: '#FFFFFF' },
  timerLabel: { fontSize: 12, color: 'rgba(255,255,255,0.7)' },
  timerSeparator: { fontSize: 48, fontWeight: '700', color: '#FFFFFF', marginHorizontal: spacing.sm },
  detailsCard: { padding: spacing.md, marginBottom: spacing.lg },
  detailRow: { flexDirection: 'row', alignItems: 'center', paddingVertical: spacing.sm, borderBottomWidth: 1, borderBottomColor: '#F3F4F6' },
  detailInfo: { marginLeft: spacing.md },
  costRow: { borderBottomWidth: 0, justifyContent: 'space-between', marginTop: spacing.md, paddingTop: spacing.md, borderTopWidth: 1, borderTopColor: '#E5E7EB' },
  endButton: { marginTop: spacing.md },
});
