import * as Keychain from 'react-native-keychain';

const AUTH_SERVICE = 'parking-app-auth';
const ACCESS_TOKEN_KEY = 'access_token';
const REFRESH_TOKEN_KEY = 'refresh_token';

/**
 * Secure storage for authentication tokens using react-native-keychain
 */
export const secureStorage = {
  /**
   * Store access token securely
   */
  setAccessToken: async (token: string): Promise<void> => {
    await Keychain.setGenericPassword(ACCESS_TOKEN_KEY, token, {
      service: `${AUTH_SERVICE}-access`,
      accessible: Keychain.ACCESSIBLE.WHEN_UNLOCKED,
    });
  },

  /**
   * Retrieve access token
   */
  getAccessToken: async (): Promise<string | null> => {
    try {
      const credentials = await Keychain.getGenericPassword({
        service: `${AUTH_SERVICE}-access`,
      });
      if (credentials) {
        return credentials.password;
      }
      return null;
    } catch {
      return null;
    }
  },

  /**
   * Store refresh token securely
   */
  setRefreshToken: async (token: string): Promise<void> => {
    await Keychain.setGenericPassword(REFRESH_TOKEN_KEY, token, {
      service: `${AUTH_SERVICE}-refresh`,
      accessible: Keychain.ACCESSIBLE.WHEN_UNLOCKED,
    });
  },

  /**
   * Retrieve refresh token
   */
  getRefreshToken: async (): Promise<string | null> => {
    try {
      const credentials = await Keychain.getGenericPassword({
        service: `${AUTH_SERVICE}-refresh`,
      });
      if (credentials) {
        return credentials.password;
      }
      return null;
    } catch {
      return null;
    }
  },

  /**
   * Store both tokens
   */
  setTokens: async (accessToken: string, refreshToken: string): Promise<void> => {
    await Promise.all([
      secureStorage.setAccessToken(accessToken),
      secureStorage.setRefreshToken(refreshToken),
    ]);
  },

  /**
   * Clear all stored tokens
   */
  clearTokens: async (): Promise<void> => {
    await Promise.all([
      Keychain.resetGenericPassword({ service: `${AUTH_SERVICE}-access` }),
      Keychain.resetGenericPassword({ service: `${AUTH_SERVICE}-refresh` }),
    ]);
  },

  /**
   * Check if tokens exist
   */
  hasTokens: async (): Promise<boolean> => {
    const accessToken = await secureStorage.getAccessToken();
    return accessToken !== null;
  },
};
