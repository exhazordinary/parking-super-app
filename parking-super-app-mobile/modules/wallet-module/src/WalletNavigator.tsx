import React from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import type { WalletStackParamList } from '@parking/navigation';
import { WalletHomeScreen, TopUpScreen, TransactionHistoryScreen } from './screens';

const Stack = createNativeStackNavigator<WalletStackParamList>();

export default function WalletNavigator(): React.JSX.Element {
  return (
    <Stack.Navigator screenOptions={{ headerShown: false, animation: 'slide_from_right' }}>
      <Stack.Screen name="WalletHome" component={WalletHomeScreen} />
      <Stack.Screen name="TopUp" component={TopUpScreen} />
      <Stack.Screen name="TransactionHistory" component={TransactionHistoryScreen} />
    </Stack.Navigator>
  );
}
