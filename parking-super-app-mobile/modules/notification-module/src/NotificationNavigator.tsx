import React from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import type { NotificationStackParamList } from '@parking/navigation';
import { NotificationListScreen, PreferencesScreen } from './screens';

const Stack = createNativeStackNavigator<NotificationStackParamList>();

export default function NotificationNavigator(): React.JSX.Element {
  return (
    <Stack.Navigator screenOptions={{ headerShown: false, animation: 'slide_from_right' }}>
      <Stack.Screen name="NotificationList" component={NotificationListScreen} />
      <Stack.Screen name="NotificationPreferences" component={PreferencesScreen} />
    </Stack.Navigator>
  );
}
