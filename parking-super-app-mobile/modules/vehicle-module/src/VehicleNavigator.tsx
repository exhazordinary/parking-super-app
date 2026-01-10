import React from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import type { VehicleStackParamList } from '@parking/navigation';
import { VehicleListScreen, AddVehicleScreen } from './screens';

const Stack = createNativeStackNavigator<VehicleStackParamList>();

export default function VehicleNavigator(): React.JSX.Element {
  return (
    <Stack.Navigator screenOptions={{ headerShown: false, animation: 'slide_from_right' }}>
      <Stack.Screen name="VehicleList" component={VehicleListScreen} />
      <Stack.Screen name="AddVehicle" component={AddVehicleScreen} />
    </Stack.Navigator>
  );
}
