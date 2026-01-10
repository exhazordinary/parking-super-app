import React, { Suspense, lazy } from 'react';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { ActivityIndicator, View } from 'react-native';
import MaterialCommunityIcons from 'react-native-vector-icons/MaterialCommunityIcons';
import { useTheme } from 'react-native-paper';

import { useAuthStore } from '@parking/auth';
import type { RootStackParamList, MainTabParamList } from '@parking/navigation';

// Lazy load federated modules
const AuthModule = lazy(() => import('auth_module/AuthNavigator'));
const WalletModule = lazy(() => import('wallet_module/WalletNavigator'));
const ParkingModule = lazy(() => import('parking_module/ParkingNavigator'));
const VehicleModule = lazy(() => import('vehicle_module/VehicleNavigator'));
const ProviderModule = lazy(() => import('provider_module/ProviderNavigator'));
const NotificationModule = lazy(() => import('notification_module/NotificationNavigator'));

const Stack = createNativeStackNavigator<RootStackParamList>();
const Tab = createBottomTabNavigator<MainTabParamList>();

function LoadingFallback(): React.JSX.Element {
  return (
    <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
      <ActivityIndicator size="large" />
    </View>
  );
}

function MainTabNavigator(): React.JSX.Element {
  const theme = useTheme();

  return (
    <Tab.Navigator
      screenOptions={{
        headerShown: false,
        tabBarActiveTintColor: theme.colors.primary,
        tabBarInactiveTintColor: theme.colors.outline,
        tabBarStyle: {
          backgroundColor: theme.colors.surface,
          borderTopColor: theme.colors.outlineVariant,
        },
      }}
    >
      <Tab.Screen
        name="ParkingTab"
        options={{
          tabBarLabel: 'Parking',
          tabBarIcon: ({ color, size }) => (
            <MaterialCommunityIcons name="parking" size={size} color={color} />
          ),
        }}
      >
        {() => (
          <Suspense fallback={<LoadingFallback />}>
            <ParkingModule />
          </Suspense>
        )}
      </Tab.Screen>

      <Tab.Screen
        name="WalletTab"
        options={{
          tabBarLabel: 'Wallet',
          tabBarIcon: ({ color, size }) => (
            <MaterialCommunityIcons name="wallet" size={size} color={color} />
          ),
        }}
      >
        {() => (
          <Suspense fallback={<LoadingFallback />}>
            <WalletModule />
          </Suspense>
        )}
      </Tab.Screen>

      <Tab.Screen
        name="VehicleTab"
        options={{
          tabBarLabel: 'Vehicles',
          tabBarIcon: ({ color, size }) => (
            <MaterialCommunityIcons name="car" size={size} color={color} />
          ),
        }}
      >
        {() => (
          <Suspense fallback={<LoadingFallback />}>
            <VehicleModule />
          </Suspense>
        )}
      </Tab.Screen>

      <Tab.Screen
        name="ProfileTab"
        options={{
          tabBarLabel: 'Profile',
          tabBarIcon: ({ color, size }) => (
            <MaterialCommunityIcons name="account" size={size} color={color} />
          ),
        }}
      >
        {() => (
          <Suspense fallback={<LoadingFallback />}>
            <AuthModule initialScreen="Profile" />
          </Suspense>
        )}
      </Tab.Screen>
    </Tab.Navigator>
  );
}

export function RootNavigator(): React.JSX.Element {
  const { isAuthenticated, isInitialized } = useAuthStore();

  if (!isInitialized) {
    return <LoadingFallback />;
  }

  return (
    <Stack.Navigator screenOptions={{ headerShown: false }}>
      {isAuthenticated ? (
        <>
          <Stack.Screen name="Main" component={MainTabNavigator} />
          <Stack.Screen name="Providers">
            {() => (
              <Suspense fallback={<LoadingFallback />}>
                <ProviderModule />
              </Suspense>
            )}
          </Stack.Screen>
          <Stack.Screen name="Notifications">
            {() => (
              <Suspense fallback={<LoadingFallback />}>
                <NotificationModule />
              </Suspense>
            )}
          </Stack.Screen>
        </>
      ) : (
        <Stack.Screen name="Auth">
          {() => (
            <Suspense fallback={<LoadingFallback />}>
              <AuthModule />
            </Suspense>
          )}
        </Stack.Screen>
      )}
    </Stack.Navigator>
  );
}
