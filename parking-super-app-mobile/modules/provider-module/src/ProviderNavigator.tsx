import React from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import type { ProviderStackParamList } from '@parking/navigation';
import { ProviderListScreen, ProviderDetailScreen } from './screens';

const Stack = createNativeStackNavigator<ProviderStackParamList>();

export default function ProviderNavigator(): React.JSX.Element {
  return (
    <Stack.Navigator screenOptions={{ headerShown: false, animation: 'slide_from_right' }}>
      <Stack.Screen name="ProviderList" component={ProviderListScreen} />
      <Stack.Screen name="ProviderDetail" component={ProviderDetailScreen} />
    </Stack.Navigator>
  );
}
