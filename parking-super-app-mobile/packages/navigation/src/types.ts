import type { NavigatorScreenParams } from '@react-navigation/native';
import type { NativeStackScreenProps } from '@react-navigation/native-stack';
import type { BottomTabScreenProps } from '@react-navigation/bottom-tabs';

/**
 * Root stack navigator params
 */
export type RootStackParamList = {
  Auth: NavigatorScreenParams<AuthStackParamList>;
  Main: NavigatorScreenParams<MainTabParamList>;
  Providers: NavigatorScreenParams<ProviderStackParamList>;
  Notifications: NavigatorScreenParams<NotificationStackParamList>;
};

/**
 * Main tab navigator params
 */
export type MainTabParamList = {
  ParkingTab: NavigatorScreenParams<ParkingStackParamList>;
  WalletTab: NavigatorScreenParams<WalletStackParamList>;
  VehicleTab: NavigatorScreenParams<VehicleStackParamList>;
  ProfileTab: NavigatorScreenParams<AuthStackParamList>;
};

/**
 * Auth module screens
 */
export type AuthStackParamList = {
  Login: undefined;
  Register: undefined;
  OTP: {
    phone: string;
    type: 'login' | 'register' | 'reset_password' | 'verify_phone';
  };
  ForgotPassword: undefined;
  ResetPassword: {
    token: string;
  };
  Profile: undefined;
  EditProfile: undefined;
  ChangePassword: undefined;
  Settings: undefined;
};

/**
 * Wallet module screens
 */
export type WalletStackParamList = {
  WalletHome: undefined;
  TopUp: undefined;
  TopUpAmount: {
    paymentMethod: string;
  };
  TopUpPayment: {
    amount: number;
    paymentMethod: string;
    idempotencyKey: string;
  };
  TransactionHistory: undefined;
  TransactionDetail: {
    transactionId: string;
  };
};

/**
 * Parking module screens
 */
export type ParkingStackParamList = {
  ParkingHome: undefined;
  FindParking: undefined;
  ParkingMap: {
    latitude?: number;
    longitude?: number;
  };
  LocationDetail: {
    locationId: string;
  };
  StartSession: {
    locationId: string;
    vehicleId?: string;
  };
  ActiveSession: {
    sessionId: string;
  };
  EndSession: {
    sessionId: string;
  };
  SessionHistory: undefined;
  SessionDetail: {
    sessionId: string;
  };
};

/**
 * Vehicle module screens
 */
export type VehicleStackParamList = {
  VehicleList: undefined;
  AddVehicle: undefined;
  EditVehicle: {
    vehicleId: string;
  };
  VehicleDetail: {
    vehicleId: string;
  };
};

/**
 * Provider module screens
 */
export type ProviderStackParamList = {
  ProviderList: undefined;
  ProviderDetail: {
    providerId: string;
  };
  ProviderLocations: {
    providerId: string;
  };
};

/**
 * Notification module screens
 */
export type NotificationStackParamList = {
  NotificationList: undefined;
  NotificationPreferences: undefined;
};

// Screen props types
export type RootStackScreenProps<T extends keyof RootStackParamList> =
  NativeStackScreenProps<RootStackParamList, T>;

export type MainTabScreenProps<T extends keyof MainTabParamList> =
  BottomTabScreenProps<MainTabParamList, T>;

export type AuthStackScreenProps<T extends keyof AuthStackParamList> =
  NativeStackScreenProps<AuthStackParamList, T>;

export type WalletStackScreenProps<T extends keyof WalletStackParamList> =
  NativeStackScreenProps<WalletStackParamList, T>;

export type ParkingStackScreenProps<T extends keyof ParkingStackParamList> =
  NativeStackScreenProps<ParkingStackParamList, T>;

export type VehicleStackScreenProps<T extends keyof VehicleStackParamList> =
  NativeStackScreenProps<VehicleStackParamList, T>;

export type ProviderStackScreenProps<T extends keyof ProviderStackParamList> =
  NativeStackScreenProps<ProviderStackParamList, T>;

export type NotificationStackScreenProps<T extends keyof NotificationStackParamList> =
  NativeStackScreenProps<NotificationStackParamList, T>;

// Declaration merging for useNavigation type safety
declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace ReactNavigation {
    interface RootParamList extends RootStackParamList {}
  }
}
