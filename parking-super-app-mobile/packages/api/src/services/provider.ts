import { apiClient } from '../client';
import type {
  ApiResponse,
  PaginatedResponse,
  Provider,
  ParkingLocation,
  GetProvidersParams,
  GetProviderLocationsParams,
} from '../types';

const PROVIDERS_BASE = '/api/v1/providers';

export const providerService = {
  /**
   * Get list of providers
   */
  getProviders: async (
    params?: GetProvidersParams
  ): Promise<PaginatedResponse<Provider>> => {
    const response = await apiClient.get<ApiResponse<PaginatedResponse<Provider>>>(
      PROVIDERS_BASE,
      { params }
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Get provider by ID
   */
  getProvider: async (id: string): Promise<Provider> => {
    const response = await apiClient.get<ApiResponse<Provider>>(
      `${PROVIDERS_BASE}/${id}`
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Get provider's locations
   */
  getProviderLocations: async (
    providerId: string,
    params?: GetProviderLocationsParams
  ): Promise<PaginatedResponse<ParkingLocation>> => {
    const response = await apiClient.get<ApiResponse<PaginatedResponse<ParkingLocation>>>(
      `${PROVIDERS_BASE}/${providerId}/locations`,
      { params }
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },
};
