import React from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import type { ParkingStackParamList } from '@parking/navigation';
import { ParkingHomeScreen, FindParkingScreen, StartSessionScreen, ActiveSessionScreen } from './screens';

const Stack = createNativeStackNavigator<ParkingStackParamList>();

export default function ParkingNavigator(): React.JSX.Element {
  return (
    <Stack.Navigator screenOptions={{ headerShown: false, animation: 'slide_from_right' }}>
      <Stack.Screen name="ParkingHome" component={ParkingHomeScreen} />
      <Stack.Screen name="FindParking" component={FindParkingScreen} />
      <Stack.Screen name="StartSession" component={StartSessionScreen} />
      <Stack.Screen name="ActiveSession" component={ActiveSessionScreen} />
    </Stack.Navigator>
  );
}
