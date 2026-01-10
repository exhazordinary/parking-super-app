/**
 * Query key factory for React Query
 * Provides type-safe, consistent query keys across the app
 */

export const queryKeys = {
  // Auth
  auth: {
    all: ['auth'] as const,
    user: () => [...queryKeys.auth.all, 'user'] as const,
    profile: () => [...queryKeys.auth.all, 'profile'] as const,
  },

  // Wallet
  wallet: {
    all: ['wallet'] as const,
    balance: () => [...queryKeys.wallet.all, 'balance'] as const,
    transactions: (filters?: Record<string, unknown>) =>
      [...queryKeys.wallet.all, 'transactions', filters] as const,
    transaction: (id: string) => [...queryKeys.wallet.all, 'transaction', id] as const,
    paymentMethods: () => [...queryKeys.wallet.all, 'paymentMethods'] as const,
  },

  // Parking
  parking: {
    all: ['parking'] as const,
    sessions: (filters?: Record<string, unknown>) =>
      [...queryKeys.parking.all, 'sessions', filters] as const,
    session: (id: string) => [...queryKeys.parking.all, 'session', id] as const,
    activeSession: () => [...queryKeys.parking.all, 'activeSession'] as const,
    locations: (params?: Record<string, unknown>) =>
      [...queryKeys.parking.all, 'locations', params] as const,
    location: (id: string) => [...queryKeys.parking.all, 'location', id] as const,
  },

  // Vehicles
  vehicles: {
    all: ['vehicles'] as const,
    list: () => [...queryKeys.vehicles.all, 'list'] as const,
    vehicle: (id: string) => [...queryKeys.vehicles.all, 'vehicle', id] as const,
  },

  // Providers
  providers: {
    all: ['providers'] as const,
    list: (filters?: Record<string, unknown>) =>
      [...queryKeys.providers.all, 'list', filters] as const,
    provider: (id: string) => [...queryKeys.providers.all, 'provider', id] as const,
    locations: (providerId: string, filters?: Record<string, unknown>) =>
      [...queryKeys.providers.all, providerId, 'locations', filters] as const,
  },

  // Notifications
  notifications: {
    all: ['notifications'] as const,
    list: (filters?: Record<string, unknown>) =>
      [...queryKeys.notifications.all, 'list', filters] as const,
    unreadCount: () => [...queryKeys.notifications.all, 'unreadCount'] as const,
    preferences: () => [...queryKeys.notifications.all, 'preferences'] as const,
  },
} as const;
