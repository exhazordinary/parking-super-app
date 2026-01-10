import { create } from 'zustand';
import { secureStorage } from './storage';
import { setTokenProvider } from '@parking/api';
import type { User, AuthTokens } from '@parking/api';

interface AuthState {
  // State
  user: User | null;
  isAuthenticated: boolean;
  isInitialized: boolean;
  isLoading: boolean;

  // Actions
  initialize: () => Promise<void>;
  setUser: (user: User) => void;
  setTokens: (tokens: AuthTokens) => Promise<void>;
  login: (user: User, tokens: AuthTokens) => Promise<void>;
  logout: () => Promise<void>;
  updateUser: (updates: Partial<User>) => void;
}

export const useAuthStore = create<AuthState>((set, get) => {
  // Register token provider with API client
  setTokenProvider({
    getAccessToken: () => secureStorage.getAccessToken(),
    getRefreshToken: () => secureStorage.getRefreshToken(),
    setTokens: (accessToken, refreshToken) =>
      secureStorage.setTokens(accessToken, refreshToken),
    clearTokens: () => secureStorage.clearTokens(),
  });

  return {
    // Initial state
    user: null,
    isAuthenticated: false,
    isInitialized: false,
    isLoading: false,

    /**
     * Initialize auth state from secure storage
     * Called on app start
     */
    initialize: async () => {
      if (get().isInitialized) return;

      set({ isLoading: true });

      try {
        const hasTokens = await secureStorage.hasTokens();

        if (hasTokens) {
          // We have tokens, user is authenticated
          // The actual user data will be fetched via useProfile hook
          set({
            isAuthenticated: true,
            isInitialized: true,
            isLoading: false,
          });
        } else {
          set({
            isAuthenticated: false,
            isInitialized: true,
            isLoading: false,
          });
        }
      } catch (error) {
        console.error('Failed to initialize auth:', error);
        set({
          isAuthenticated: false,
          isInitialized: true,
          isLoading: false,
        });
      }
    },

    /**
     * Set user data
     */
    setUser: (user) => {
      set({ user });
    },

    /**
     * Store tokens securely
     */
    setTokens: async (tokens) => {
      await secureStorage.setTokens(tokens.accessToken, tokens.refreshToken);
    },

    /**
     * Login - store user and tokens
     */
    login: async (user, tokens) => {
      await secureStorage.setTokens(tokens.accessToken, tokens.refreshToken);
      set({
        user,
        isAuthenticated: true,
      });
    },

    /**
     * Logout - clear user and tokens
     */
    logout: async () => {
      await secureStorage.clearTokens();
      set({
        user: null,
        isAuthenticated: false,
      });
    },

    /**
     * Update user data
     */
    updateUser: (updates) => {
      const currentUser = get().user;
      if (currentUser) {
        set({
          user: { ...currentUser, ...updates },
        });
      }
    },
  };
});
