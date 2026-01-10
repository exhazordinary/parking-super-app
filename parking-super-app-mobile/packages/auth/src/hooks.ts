import { useCallback } from 'react';
import { useAuthStore } from './store';
import {
  useLogin as useLoginMutation,
  useRegister as useRegisterMutation,
  useLogout as useLogoutMutation,
  useProfile,
} from '@parking/api';
import type { LoginRequest, RegisterRequest, User, AuthTokens } from '@parking/api';

/**
 * Hook to check if user is authenticated
 */
export function useIsAuthenticated(): boolean {
  return useAuthStore((state) => state.isAuthenticated);
}

/**
 * Hook to get current user
 */
export function useUser(): User | null {
  return useAuthStore((state) => state.user);
}

/**
 * Hook to handle login flow
 */
export function useAuth() {
  const { login, logout, isAuthenticated, user, setUser } = useAuthStore();

  const loginMutation = useLoginMutation();
  const registerMutation = useRegisterMutation();
  const logoutMutation = useLogoutMutation();

  // Fetch profile on mount if authenticated
  const profileQuery = useProfile();

  // Sync profile data to store
  if (profileQuery.data && !user) {
    setUser(profileQuery.data);
  }

  const handleLogin = useCallback(
    async (data: LoginRequest): Promise<{ user: User; tokens: AuthTokens }> => {
      const result = await loginMutation.mutateAsync(data);
      await login(result.user, result.tokens);
      return result;
    },
    [loginMutation, login]
  );

  const handleRegister = useCallback(
    async (data: RegisterRequest): Promise<{ user: User; tokens: AuthTokens }> => {
      const result = await registerMutation.mutateAsync(data);
      await login(result.user, result.tokens);
      return result;
    },
    [registerMutation, login]
  );

  const handleLogout = useCallback(async (): Promise<void> => {
    try {
      await logoutMutation.mutateAsync();
    } catch {
      // Logout from server failed, but still clear local state
    }
    await logout();
  }, [logoutMutation, logout]);

  return {
    user,
    isAuthenticated,
    isLoggingIn: loginMutation.isPending,
    isRegistering: registerMutation.isPending,
    isLoggingOut: logoutMutation.isPending,
    loginError: loginMutation.error,
    registerError: registerMutation.error,
    login: handleLogin,
    register: handleRegister,
    logout: handleLogout,
  };
}

/**
 * Hook to refresh user profile
 */
export function useRefreshProfile() {
  const { setUser } = useAuthStore();
  const profileQuery = useProfile();

  const refresh = useCallback(async () => {
    const result = await profileQuery.refetch();
    if (result.data) {
      setUser(result.data);
    }
    return result.data;
  }, [profileQuery, setUser]);

  return {
    refresh,
    isRefreshing: profileQuery.isRefetching,
    error: profileQuery.error,
  };
}
