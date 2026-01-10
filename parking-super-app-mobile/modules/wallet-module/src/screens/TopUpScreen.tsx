import React, { useState, useCallback } from 'react';
import { View, StyleSheet, ScrollView, Pressable } from 'react-native';
import { Text, useTheme, RadioButton } from 'react-native-paper';
import { SafeAreaView } from 'react-native-safe-area-context';
import MaterialCommunityIcons from 'react-native-vector-icons/MaterialCommunityIcons';
import type { MD3Theme } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';
import { v4 as uuidv4 } from 'uuid';

import { Button, Card, Input, LoadingOverlay, spacing } from '@parking/ui';
import { useTopUp, usePaymentMethods } from '@parking/api';
import type { WalletStackScreenProps, PaymentMethod } from '@parking/navigation';

type NavigationProp = WalletStackScreenProps<'TopUp'>['navigation'];

const QUICK_AMOUNTS = [10, 20, 50, 100];

export function TopUpScreen(): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const navigation = useNavigation<NavigationProp>();
  const topUp = useTopUp();
  const { data: paymentMethods } = usePaymentMethods();

  const [amount, setAmount] = useState('');
  const [selectedMethod, setSelectedMethod] = useState<PaymentMethod>('fpx');
  const [error, setError] = useState('');

  const handleQuickAmount = useCallback((value: number) => {
    setAmount(value.toString());
    setError('');
  }, []);

  const handleAmountChange = useCallback((text: string) => {
    const cleaned = text.replace(/[^0-9.]/g, '');
    setAmount(cleaned);
    setError('');
  }, []);

  const handleTopUp = useCallback(async () => {
    const numAmount = parseFloat(amount);
    if (isNaN(numAmount) || numAmount < 10) {
      setError('Minimum top-up amount is RM 10');
      return;
    }
    if (numAmount > 1000) {
      setError('Maximum top-up amount is RM 1,000');
      return;
    }

    try {
      await topUp.mutateAsync({
        amount: numAmount,
        paymentMethod: selectedMethod,
        idempotencyKey: uuidv4(),
      });
      navigation.goBack();
    } catch (err) {
      setError('Top-up failed. Please try again.');
    }
  }, [amount, selectedMethod, topUp, navigation]);

  const methods = paymentMethods || [
    { type: 'fpx' as PaymentMethod, name: 'Online Banking (FPX)', icon: 'bank', enabled: true, minAmount: 10, maxAmount: 1000 },
    { type: 'card' as PaymentMethod, name: 'Credit/Debit Card', icon: 'credit-card', enabled: true, minAmount: 10, maxAmount: 1000 },
    { type: 'ewallet_tng' as PaymentMethod, name: 'Touch n Go eWallet', icon: 'wallet', enabled: true, minAmount: 10, maxAmount: 500 },
  ];

  return (
    <SafeAreaView style={styles.safeArea}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <Text variant="headlineSmall" style={styles.title}>Top Up Wallet</Text>

        {/* Amount Input */}
        <Card style={styles.section}>
          <Text variant="labelLarge" style={styles.sectionTitle}>Enter Amount</Text>
          <View style={styles.amountInput}>
            <Text variant="headlineMedium" style={{ color: theme.colors.onSurfaceVariant }}>RM</Text>
            <Input
              label=""
              value={amount}
              onChangeText={handleAmountChange}
              keyboardType="decimal-pad"
              style={styles.amountField}
              error={Boolean(error)}
            />
          </View>
          {error && <Text variant="bodySmall" style={{ color: theme.colors.error, marginTop: spacing.xs }}>{error}</Text>}

          <View style={styles.quickAmounts}>
            {QUICK_AMOUNTS.map((val) => (
              <Pressable
                key={val}
                onPress={() => handleQuickAmount(val)}
                style={[styles.quickButton, amount === val.toString() && { backgroundColor: theme.colors.primaryContainer }]}
              >
                <Text style={{ color: amount === val.toString() ? theme.colors.primary : theme.colors.onSurface }}>
                  RM {val}
                </Text>
              </Pressable>
            ))}
          </View>
        </Card>

        {/* Payment Method */}
        <Card style={styles.section}>
          <Text variant="labelLarge" style={styles.sectionTitle}>Payment Method</Text>
          <RadioButton.Group onValueChange={(v) => setSelectedMethod(v as PaymentMethod)} value={selectedMethod}>
            {methods.filter(m => m.enabled).map((method) => (
              <Pressable
                key={method.type}
                onPress={() => setSelectedMethod(method.type)}
                style={[styles.methodItem, selectedMethod === method.type && { backgroundColor: theme.colors.primaryContainer }]}
              >
                <MaterialCommunityIcons name={method.icon} size={24} color={theme.colors.primary} />
                <Text variant="bodyLarge" style={styles.methodName}>{method.name}</Text>
                <RadioButton value={method.type} />
              </Pressable>
            ))}
          </RadioButton.Group>
        </Card>

        <Button mode="contained" onPress={handleTopUp} loading={topUp.isPending} disabled={!amount || topUp.isPending} fullWidth style={styles.confirmButton}>
          Confirm Top Up
        </Button>
      </ScrollView>

      <LoadingOverlay visible={topUp.isPending} message="Processing..." />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: { flex: 1, backgroundColor: '#F3F4F6' },
  scrollContent: { padding: spacing.md },
  title: { fontWeight: '600', marginBottom: spacing.lg },
  section: { marginBottom: spacing.lg, padding: spacing.md },
  sectionTitle: { marginBottom: spacing.md },
  amountInput: { flexDirection: 'row', alignItems: 'center', gap: spacing.sm },
  amountField: { flex: 1, fontSize: 24 },
  quickAmounts: { flexDirection: 'row', gap: spacing.sm, marginTop: spacing.md },
  quickButton: { flex: 1, paddingVertical: spacing.sm, borderRadius: 8, borderWidth: 1, borderColor: '#E5E7EB', alignItems: 'center' },
  methodItem: { flexDirection: 'row', alignItems: 'center', padding: spacing.md, borderRadius: 8, marginBottom: spacing.xs },
  methodName: { flex: 1, marginLeft: spacing.md },
  confirmButton: { marginTop: spacing.md },
});
