import type { LinkingOptions } from '@react-navigation/native';
import type { RootStackParamList } from './types';

/**
 * Deep linking configuration for the app
 * Supports parkingapp:// scheme
 */
export const linking: LinkingOptions<RootStackParamList> = {
  prefixes: ['parkingapp://', 'https://parkingapp.com'],

  config: {
    screens: {
      Auth: {
        screens: {
          Login: 'login',
          Register: 'register',
          OTP: 'otp/:phone/:type',
          ForgotPassword: 'forgot-password',
          ResetPassword: 'reset-password/:token',
        },
      },
      Main: {
        screens: {
          ParkingTab: {
            screens: {
              ParkingHome: 'parking',
              FindParking: 'parking/find',
              ParkingMap: 'parking/map',
              LocationDetail: 'parking/location/:locationId',
              StartSession: 'parking/start/:locationId',
              ActiveSession: 'parking/session/:sessionId',
              EndSession: 'parking/session/:sessionId/end',
              SessionHistory: 'parking/history',
              SessionDetail: 'parking/history/:sessionId',
            },
          },
          WalletTab: {
            screens: {
              WalletHome: 'wallet',
              TopUp: 'wallet/topup',
              TopUpAmount: 'wallet/topup/:paymentMethod',
              TransactionHistory: 'wallet/transactions',
              TransactionDetail: 'wallet/transactions/:transactionId',
            },
          },
          VehicleTab: {
            screens: {
              VehicleList: 'vehicles',
              AddVehicle: 'vehicles/add',
              EditVehicle: 'vehicles/:vehicleId/edit',
              VehicleDetail: 'vehicles/:vehicleId',
            },
          },
          ProfileTab: {
            screens: {
              Profile: 'profile',
              EditProfile: 'profile/edit',
              ChangePassword: 'profile/password',
              Settings: 'settings',
            },
          },
        },
      },
      Providers: {
        screens: {
          ProviderList: 'providers',
          ProviderDetail: 'providers/:providerId',
          ProviderLocations: 'providers/:providerId/locations',
        },
      },
      Notifications: {
        screens: {
          NotificationList: 'notifications',
          NotificationPreferences: 'notifications/preferences',
        },
      },
    },
  },
};
