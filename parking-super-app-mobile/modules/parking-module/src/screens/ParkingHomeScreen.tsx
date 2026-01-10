import React, { useCallback } from 'react';
import { View, StyleSheet, ScrollView, RefreshControl, Pressable } from 'react-native';
import { Text, useTheme, IconButton } from 'react-native-paper';
import { SafeAreaView } from 'react-native-safe-area-context';
import MaterialCommunityIcons from 'react-native-vector-icons/MaterialCommunityIcons';
import type { MD3Theme } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';

import { Button, Card, LoadingSpinner, spacing } from '@parking/ui';
import { useActiveSession, useSessions } from '@parking/api';
import type { ParkingStackScreenProps, ParkingSession } from '@parking/navigation';

type NavigationProp = ParkingStackScreenProps<'ParkingHome'>['navigation'];

function ActiveSessionCard({ session, onPress }: { session: ParkingSession; onPress: () => void }): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const startTime = new Date(session.startTime);
  const duration = Math.floor((Date.now() - startTime.getTime()) / 60000);
  const hours = Math.floor(duration / 60);
  const mins = duration % 60;

  return (
    <Pressable onPress={onPress}>
      <Card style={[styles.activeCard, { backgroundColor: theme.colors.primaryContainer }]}>
        <View style={styles.activeHeader}>
          <MaterialCommunityIcons name="car" size={24} color={theme.colors.primary} />
          <Text variant="titleMedium" style={{ fontWeight: '600', marginLeft: spacing.sm }}>Active Session</Text>
        </View>
        <Text variant="bodyLarge" style={{ marginTop: spacing.sm }}>{session.location?.name || 'Parking Location'}</Text>
        <View style={styles.activeInfo}>
          <View>
            <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>Duration</Text>
            <Text variant="titleMedium" style={{ fontWeight: '600' }}>{hours}h {mins}m</Text>
          </View>
          <View>
            <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>Est. Cost</Text>
            <Text variant="titleMedium" style={{ fontWeight: '600' }}>RM {(session.totalCost || 0).toFixed(2)}</Text>
          </View>
          <Button mode="contained" compact onPress={onPress}>End Session</Button>
        </View>
      </Card>
    </Pressable>
  );
}

export function ParkingHomeScreen(): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const navigation = useNavigation<NavigationProp>();

  const { data: activeSession, isLoading: activeLoading, refetch: refetchActive } = useActiveSession();
  const { data: sessionsData, isLoading: historyLoading, refetch: refetchHistory } = useSessions({ limit: 3, status: 'completed' });

  const recentSessions = sessionsData?.pages[0]?.items || [];
  const isLoading = activeLoading || historyLoading;

  const handleRefresh = useCallback(() => {
    refetchActive();
    refetchHistory();
  }, [refetchActive, refetchHistory]);

  const handleFindParking = useCallback(() => {
    navigation.navigate('FindParking');
  }, [navigation]);

  const handleActiveSession = useCallback(() => {
    if (activeSession) {
      navigation.navigate('ActiveSession', { sessionId: activeSession.id });
    }
  }, [navigation, activeSession]);

  const handleSessionHistory = useCallback(() => {
    navigation.navigate('SessionHistory');
  }, [navigation]);

  return (
    <SafeAreaView style={styles.safeArea}>
      <View style={styles.header}>
        <Text variant="headlineMedium" style={{ fontWeight: '700' }}>Parking</Text>
        <IconButton icon="bell-outline" onPress={() => {}} />
      </View>

      <ScrollView
        contentContainerStyle={styles.scrollContent}
        refreshControl={<RefreshControl refreshing={isLoading} onRefresh={handleRefresh} />}
      >
        {/* Active Session */}
        {activeSession && <ActiveSessionCard session={activeSession} onPress={handleActiveSession} />}

        {/* Find Parking CTA */}
        <Card style={styles.findCard}>
          <MaterialCommunityIcons name="map-marker-radius" size={48} color={theme.colors.primary} />
          <Text variant="titleMedium" style={{ marginTop: spacing.md, fontWeight: '600' }}>Find Parking</Text>
          <Text variant="bodyMedium" style={{ color: theme.colors.onSurfaceVariant, textAlign: 'center', marginVertical: spacing.sm }}>
            Discover nearby parking locations
          </Text>
          <Button mode="contained" onPress={handleFindParking} icon="magnify">
            Search Nearby
          </Button>
        </Card>

        {/* Recent Sessions */}
        <View style={styles.section}>
          <View style={styles.sectionHeader}>
            <Text variant="titleMedium" style={{ fontWeight: '600' }}>Recent Sessions</Text>
            <Button mode="text" onPress={handleSessionHistory} compact>View All</Button>
          </View>
          {recentSessions.length === 0 ? (
            <Text variant="bodyMedium" style={{ color: theme.colors.onSurfaceVariant, textAlign: 'center', padding: spacing.lg }}>
              No recent sessions
            </Text>
          ) : (
            recentSessions.map((session) => (
              <Card key={session.id} style={styles.sessionCard} onPress={() => navigation.navigate('SessionDetail', { sessionId: session.id })}>
                <View style={styles.sessionRow}>
                  <View>
                    <Text variant="bodyMedium" style={{ fontWeight: '500' }}>{session.location?.name || 'Location'}</Text>
                    <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>
                      {new Date(session.startTime).toLocaleDateString()}
                    </Text>
                  </View>
                  <Text variant="titleSmall" style={{ fontWeight: '600' }}>RM {(session.totalCost || 0).toFixed(2)}</Text>
                </View>
              </Card>
            ))
          )}
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1, backgroundColor: '#F3F4F6' },
  header: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', paddingHorizontal: spacing.md },
  scrollContent: { padding: spacing.md },
  activeCard: { padding: spacing.md, marginBottom: spacing.md },
  activeHeader: { flexDirection: 'row', alignItems: 'center' },
  activeInfo: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginTop: spacing.md },
  findCard: { padding: spacing.lg, alignItems: 'center', marginBottom: spacing.lg },
  section: { marginBottom: spacing.lg },
  sectionHeader: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: spacing.sm },
  sessionCard: { padding: spacing.md, marginBottom: spacing.sm },
  sessionRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center' },
});
