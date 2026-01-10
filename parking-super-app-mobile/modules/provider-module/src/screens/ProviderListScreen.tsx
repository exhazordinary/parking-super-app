import React, { useCallback } from 'react';
import { View, StyleSheet, FlatList, Pressable } from 'react-native';
import { Text, useTheme, Searchbar, Avatar } from 'react-native-paper';
import { SafeAreaView } from 'react-native-safe-area-context';
import MaterialCommunityIcons from 'react-native-vector-icons/MaterialCommunityIcons';
import type { MD3Theme } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';

import { Card, LoadingSpinner, EmptyState, spacing } from '@parking/ui';
import { useProviders, type Provider } from '@parking/api';
import type { ProviderStackScreenProps } from '@parking/navigation';

type NavigationProp = ProviderStackScreenProps<'ProviderList'>['navigation'];

function ProviderCard({ provider, onPress }: { provider: Provider; onPress: () => void }): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  return (
    <Pressable onPress={onPress}>
      <Card style={styles.providerCard}>
        <View style={styles.cardContent}>
          {provider.logo ? (
            <Avatar.Image size={56} source={{ uri: provider.logo }} />
          ) : (
            <Avatar.Icon size={56} icon="parking" style={{ backgroundColor: theme.colors.primaryContainer }} />
          )}
          <View style={styles.providerInfo}>
            <Text variant="titleMedium" style={{ fontWeight: '600' }}>{provider.name}</Text>
            <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }} numberOfLines={2}>{provider.description}</Text>
            <View style={styles.statsRow}>
              <MaterialCommunityIcons name="map-marker-multiple" size={14} color={theme.colors.primary} />
              <Text variant="labelSmall" style={{ color: theme.colors.primary, marginLeft: 4 }}>{provider.totalLocations} locations</Text>
              {provider.rating && (
                <>
                  <MaterialCommunityIcons name="star" size={14} color="#F59E0B" style={{ marginLeft: spacing.md }} />
                  <Text variant="labelSmall" style={{ marginLeft: 4 }}>{provider.rating.toFixed(1)}</Text>
                </>
              )}
            </View>
          </View>
          <MaterialCommunityIcons name="chevron-right" size={24} color={theme.colors.outline} />
        </View>
      </Card>
    </Pressable>
  );
}

export function ProviderListScreen(): React.JSX.Element {
  const navigation = useNavigation<NavigationProp>();
  const { data, isLoading, fetchNextPage, hasNextPage, isFetchingNextPage } = useProviders();
  const [searchQuery, setSearchQuery] = React.useState('');

  const providers = data?.pages.flatMap((p) => p.items) || [];
  const filtered = providers.filter((p) => p.name.toLowerCase().includes(searchQuery.toLowerCase()));

  const handleProviderPress = useCallback((id: string) => navigation.navigate('ProviderDetail', { providerId: id }), [navigation]);
  const handleLoadMore = useCallback(() => { if (hasNextPage && !isFetchingNextPage) fetchNextPage(); }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

  if (isLoading) return <SafeAreaView style={styles.safeArea}><LoadingSpinner message="Loading providers..." /></SafeAreaView>;

  return (
    <SafeAreaView style={styles.safeArea}>
      <View style={styles.header}>
        <Text variant="headlineMedium" style={{ fontWeight: '700' }}>Parking Providers</Text>
        <Searchbar placeholder="Search providers" value={searchQuery} onChangeText={setSearchQuery} style={styles.searchbar} />
      </View>
      <FlatList
        data={filtered}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => <ProviderCard provider={item} onPress={() => handleProviderPress(item.id)} />}
        onEndReached={handleLoadMore}
        contentContainerStyle={styles.list}
        ListEmptyComponent={<EmptyState icon="domain" title="No Providers" description="No parking providers found" />}
        ListFooterComponent={isFetchingNextPage ? <LoadingSpinner size="small" /> : null}
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1, backgroundColor: '#F3F4F6' },
  header: { padding: spacing.md },
  searchbar: { marginTop: spacing.sm, elevation: 0, backgroundColor: '#FFFFFF' },
  list: { padding: spacing.md },
  providerCard: { marginBottom: spacing.sm, padding: 0 },
  cardContent: { flexDirection: 'row', alignItems: 'center', padding: spacing.md },
  providerInfo: { flex: 1, marginLeft: spacing.md },
  statsRow: { flexDirection: 'row', alignItems: 'center', marginTop: spacing.xs },
});
