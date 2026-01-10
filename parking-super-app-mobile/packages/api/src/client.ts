import axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig } from 'axios';
import { v4 as uuidv4 } from 'uuid';
import type { ApiResponse, ApiError } from './types';

/**
 * Token provider interface for dependency injection
 * This allows the auth package to provide tokens without circular deps
 */
export interface TokenProvider {
  getAccessToken: () => Promise<string | null>;
  getRefreshToken: () => Promise<string | null>;
  setTokens: (accessToken: string, refreshToken: string) => Promise<void>;
  clearTokens: () => Promise<void>;
}

let tokenProvider: TokenProvider | null = null;

/**
 * Set the token provider for the API client
 */
export function setTokenProvider(provider: TokenProvider): void {
  tokenProvider = provider;
}

/**
 * API Configuration
 */
const API_CONFIG = {
  // Leave baseURL empty for user to configure
  baseURL: '', // e.g., 'https://api.parkingapp.com'
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
    Accept: 'application/json',
  },
};

/**
 * Create the API client instance
 */
const apiClient: AxiosInstance = axios.create(API_CONFIG);

/**
 * Configure the API base URL
 */
export function configureApi(baseURL: string): void {
  apiClient.defaults.baseURL = baseURL;
}

/**
 * Request interceptor
 * - Adds Authorization header if token is available
 * - Adds idempotency key for mutation requests
 */
apiClient.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    // Add auth token
    if (tokenProvider) {
      const token = await tokenProvider.getAccessToken();
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }

    // Add idempotency key for POST/PUT/PATCH requests to wallet endpoints
    if (
      ['post', 'put', 'patch'].includes(config.method?.toLowerCase() ?? '') &&
      config.url?.includes('/wallet')
    ) {
      if (!config.headers['Idempotency-Key']) {
        config.headers['Idempotency-Key'] = uuidv4();
      }
    }

    return config;
  },
  (error) => Promise.reject(error)
);

/**
 * Response interceptor
 * - Handles token refresh on 401
 * - Standardizes error format
 */
let isRefreshing = false;
let refreshSubscribers: ((token: string) => void)[] = [];

function subscribeTokenRefresh(callback: (token: string) => void): void {
  refreshSubscribers.push(callback);
}

function onTokenRefreshed(token: string): void {
  refreshSubscribers.forEach((callback) => callback(token));
  refreshSubscribers = [];
}

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<ApiResponse<unknown>>) => {
    const originalRequest = error.config;

    // Handle 401 Unauthorized
    if (error.response?.status === 401 && originalRequest && tokenProvider) {
      if (!isRefreshing) {
        isRefreshing = true;

        try {
          const refreshToken = await tokenProvider.getRefreshToken();
          if (refreshToken) {
            const response = await axios.post<ApiResponse<{ accessToken: string; refreshToken: string }>>(
              `${apiClient.defaults.baseURL}/api/v1/auth/refresh`,
              { refreshToken }
            );

            if (response.data.success && response.data.data) {
              const { accessToken, refreshToken: newRefreshToken } = response.data.data;
              await tokenProvider.setTokens(accessToken, newRefreshToken);
              onTokenRefreshed(accessToken);
              isRefreshing = false;

              // Retry original request
              originalRequest.headers.Authorization = `Bearer ${accessToken}`;
              return apiClient(originalRequest);
            }
          }
        } catch {
          // Refresh failed, clear tokens
          await tokenProvider.clearTokens();
          isRefreshing = false;
          refreshSubscribers = [];
        }
      } else {
        // Wait for token refresh
        return new Promise((resolve) => {
          subscribeTokenRefresh((token) => {
            originalRequest.headers.Authorization = `Bearer ${token}`;
            resolve(apiClient(originalRequest));
          });
        });
      }
    }

    // Transform error to standard format
    const apiError: ApiError = {
      code: error.response?.data?.error?.code ?? 'UNKNOWN_ERROR',
      message: error.response?.data?.error?.message ?? error.message ?? 'An unknown error occurred',
      details: error.response?.data?.error?.details,
    };

    return Promise.reject(apiError);
  }
);

export { apiClient };
