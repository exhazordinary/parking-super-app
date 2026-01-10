import React, { useCallback } from 'react';
import { View, StyleSheet, FlatList, Pressable } from 'react-native';
import { Text, useTheme } from 'react-native-paper';
import { SafeAreaView } from 'react-native-safe-area-context';
import MaterialCommunityIcons from 'react-native-vector-icons/MaterialCommunityIcons';
import type { MD3Theme } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';

import { LoadingSpinner, EmptyState, ErrorState, spacing } from '@parking/ui';
import { useTransactions, type Transaction } from '@parking/api';
import type { WalletStackScreenProps } from '@parking/navigation';

type NavigationProp = WalletStackScreenProps<'TransactionHistory'>['navigation'];

function TransactionRow({ transaction, onPress }: { transaction: Transaction; onPress: () => void }): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const isCredit = transaction.type === 'topup' || transaction.type === 'refund' || transaction.type === 'transfer_in';

  return (
    <Pressable onPress={onPress} style={({ pressed }) => [styles.row, pressed && { backgroundColor: theme.colors.surfaceVariant }]}>
      <View style={[styles.icon, { backgroundColor: isCredit ? '#DCFCE7' : '#FEE2E2' }]}>
        <MaterialCommunityIcons name={isCredit ? 'arrow-down' : 'arrow-up'} size={20} color={isCredit ? '#16A34A' : '#DC2626'} />
      </View>
      <View style={styles.info}>
        <Text variant="bodyMedium" style={{ fontWeight: '500' }}>{transaction.description}</Text>
        <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>
          {new Date(transaction.createdAt).toLocaleString()}
        </Text>
      </View>
      <View style={styles.amountContainer}>
        <Text variant="titleSmall" style={{ color: isCredit ? '#16A34A' : '#DC2626', fontWeight: '600' }}>
          {isCredit ? '+' : '-'}RM {Math.abs(transaction.amount).toFixed(2)}
        </Text>
        <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant, textTransform: 'capitalize' }}>
          {transaction.status}
        </Text>
      </View>
    </Pressable>
  );
}

export function TransactionHistoryScreen(): React.JSX.Element {
  const navigation = useNavigation<NavigationProp>();
  const { data, isLoading, error, refetch, fetchNextPage, hasNextPage, isFetchingNextPage } = useTransactions();

  const transactions = data?.pages.flatMap((page) => page.items) || [];

  const handleTransactionPress = useCallback((id: string) => {
    navigation.navigate('TransactionDetail', { transactionId: id });
  }, [navigation]);

  const handleLoadMore = useCallback(() => {
    if (hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

  if (isLoading) return <SafeAreaView style={styles.safeArea}><LoadingSpinner message="Loading transactions..." /></SafeAreaView>;
  if (error) return <SafeAreaView style={styles.safeArea}><ErrorState message="Failed to load transactions" onRetry={refetch} /></SafeAreaView>;

  return (
    <SafeAreaView style={styles.safeArea}>
      <FlatList
        data={transactions}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => <TransactionRow transaction={item} onPress={() => handleTransactionPress(item.id)} />}
        onEndReached={handleLoadMore}
        onEndReachedThreshold={0.5}
        ListEmptyComponent={<EmptyState icon="receipt" title="No Transactions" description="Your transaction history will appear here" />}
        ListFooterComponent={isFetchingNextPage ? <LoadingSpinner size="small" /> : null}
        contentContainerStyle={transactions.length === 0 ? { flex: 1 } : undefined}
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1, backgroundColor: '#FFFFFF' },
  row: { flexDirection: 'row', alignItems: 'center', padding: spacing.md, borderBottomWidth: 1, borderBottomColor: '#F3F4F6' },
  icon: { width: 40, height: 40, borderRadius: 20, justifyContent: 'center', alignItems: 'center', marginRight: spacing.md },
  info: { flex: 1 },
  amountContainer: { alignItems: 'flex-end' },
});
