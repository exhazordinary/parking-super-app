import React, { useCallback } from 'react';
import { View, StyleSheet, ScrollView, RefreshControl } from 'react-native';
import { Text, useTheme } from 'react-native-paper';
import { SafeAreaView } from 'react-native-safe-area-context';
import MaterialCommunityIcons from 'react-native-vector-icons/MaterialCommunityIcons';
import type { MD3Theme } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';

import { Button, Card, LoadingSpinner, ErrorState, spacing } from '@parking/ui';
import { useWallet, useTransactions, type Transaction } from '@parking/api';
import type { WalletStackScreenProps } from '@parking/navigation';

type NavigationProp = WalletStackScreenProps<'WalletHome'>['navigation'];

function TransactionItem({ transaction }: { transaction: Transaction }): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const isCredit = transaction.type === 'topup' || transaction.type === 'refund' || transaction.type === 'transfer_in';

  const icon = {
    topup: 'plus-circle',
    payment: 'minus-circle',
    refund: 'arrow-u-left-top',
    transfer_in: 'arrow-down',
    transfer_out: 'arrow-up',
  }[transaction.type] || 'circle';

  return (
    <View style={styles.transactionItem}>
      <View style={[styles.transactionIcon, { backgroundColor: isCredit ? '#DCFCE7' : '#FEE2E2' }]}>
        <MaterialCommunityIcons name={icon} size={20} color={isCredit ? '#16A34A' : '#DC2626'} />
      </View>
      <View style={styles.transactionInfo}>
        <Text variant="bodyMedium" style={{ fontWeight: '500' }}>{transaction.description}</Text>
        <Text variant="bodySmall" style={{ color: theme.colors.onSurfaceVariant }}>
          {new Date(transaction.createdAt).toLocaleDateString()}
        </Text>
      </View>
      <Text variant="titleSmall" style={{ color: isCredit ? '#16A34A' : '#DC2626', fontWeight: '600' }}>
        {isCredit ? '+' : '-'}RM {Math.abs(transaction.amount).toFixed(2)}
      </Text>
    </View>
  );
}

export function WalletHomeScreen(): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const navigation = useNavigation<NavigationProp>();

  const { data: wallet, isLoading: walletLoading, error: walletError, refetch: refetchWallet } = useWallet();
  const { data: transactionsData, isLoading: txLoading, refetch: refetchTx } = useTransactions({ limit: 5 });

  const transactions = transactionsData?.pages[0]?.items || [];
  const isLoading = walletLoading || txLoading;

  const handleRefresh = useCallback(() => {
    refetchWallet();
    refetchTx();
  }, [refetchWallet, refetchTx]);

  const handleTopUp = useCallback(() => {
    navigation.navigate('TopUp');
  }, [navigation]);

  const handleViewHistory = useCallback(() => {
    navigation.navigate('TransactionHistory');
  }, [navigation]);

  if (walletLoading) {
    return <SafeAreaView style={styles.safeArea}><LoadingSpinner message="Loading wallet..." /></SafeAreaView>;
  }

  if (walletError || !wallet) {
    return <SafeAreaView style={styles.safeArea}><ErrorState message="Failed to load wallet" onRetry={refetchWallet} /></SafeAreaView>;
  }

  return (
    <SafeAreaView style={styles.safeArea}>
      <ScrollView
        contentContainerStyle={styles.scrollContent}
        refreshControl={<RefreshControl refreshing={isLoading} onRefresh={handleRefresh} />}
      >
        {/* Balance Card */}
        <Card style={[styles.balanceCard, { backgroundColor: theme.colors.primary }]}>
          <Text variant="bodyMedium" style={{ color: 'rgba(255,255,255,0.8)' }}>Available Balance</Text>
          <Text variant="displaySmall" style={styles.balanceAmount}>
            RM {wallet.balance.toFixed(2)}
          </Text>
          <Button mode="contained" onPress={handleTopUp} style={styles.topUpButton} buttonColor="#FFFFFF" textColor={theme.colors.primary}>
            Top Up
          </Button>
        </Card>

        {/* Quick Actions */}
        <View style={styles.quickActions}>
          <Button mode="outlined" icon="history" onPress={handleViewHistory} style={styles.actionButton}>
            History
          </Button>
        </View>

        {/* Recent Transactions */}
        <View style={styles.section}>
          <View style={styles.sectionHeader}>
            <Text variant="titleMedium" style={{ fontWeight: '600' }}>Recent Transactions</Text>
            <Button mode="text" onPress={handleViewHistory} compact>View All</Button>
          </View>
          <Card style={styles.transactionsCard}>
            {transactions.length === 0 ? (
              <Text variant="bodyMedium" style={{ color: theme.colors.onSurfaceVariant, textAlign: 'center', padding: spacing.lg }}>
                No transactions yet
              </Text>
            ) : (
              transactions.map((tx) => <TransactionItem key={tx.id} transaction={tx} />)
            )}
          </Card>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1, backgroundColor: '#F3F4F6' },
  scrollContent: { padding: spacing.md },
  balanceCard: { padding: spacing.lg, borderRadius: 16, marginBottom: spacing.md },
  balanceAmount: { color: '#FFFFFF', fontWeight: '700', marginVertical: spacing.sm },
  topUpButton: { marginTop: spacing.md, alignSelf: 'flex-start' },
  quickActions: { flexDirection: 'row', gap: spacing.sm, marginBottom: spacing.lg },
  actionButton: { flex: 1 },
  section: { marginBottom: spacing.lg },
  sectionHeader: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: spacing.sm },
  transactionsCard: { padding: 0, overflow: 'hidden' },
  transactionItem: { flexDirection: 'row', alignItems: 'center', padding: spacing.md, borderBottomWidth: 1, borderBottomColor: '#E5E7EB' },
  transactionIcon: { width: 40, height: 40, borderRadius: 20, justifyContent: 'center', alignItems: 'center', marginRight: spacing.md },
  transactionInfo: { flex: 1 },
});
