import { useQuery, useInfiniteQuery } from '@tanstack/react-query';
import { providerService } from '../services/provider';
import { queryKeys } from '../queryKeys';
import type { GetProvidersParams, GetProviderLocationsParams } from '../types';

/**
 * Hook to get providers with infinite scroll
 */
export function useProviders(params?: Omit<GetProvidersParams, 'page'>) {
  return useInfiniteQuery({
    queryKey: queryKeys.providers.list(params),
    queryFn: ({ pageParam = 1 }) =>
      providerService.getProviders({ ...params, page: pageParam }),
    initialPageParam: 1,
    getNextPageParam: (lastPage) => {
      if (lastPage.page < lastPage.totalPages) {
        return lastPage.page + 1;
      }
      return undefined;
    },
  });
}

/**
 * Hook to get provider by ID
 */
export function useProvider(id: string) {
  return useQuery({
    queryKey: queryKeys.providers.provider(id),
    queryFn: () => providerService.getProvider(id),
    enabled: Boolean(id),
  });
}

/**
 * Hook to get provider locations with infinite scroll
 */
export function useProviderLocations(
  providerId: string,
  params?: Omit<GetProviderLocationsParams, 'page'>
) {
  return useInfiniteQuery({
    queryKey: queryKeys.providers.locations(providerId, params),
    queryFn: ({ pageParam = 1 }) =>
      providerService.getProviderLocations(providerId, { ...params, page: pageParam }),
    initialPageParam: 1,
    getNextPageParam: (lastPage) => {
      if (lastPage.page < lastPage.totalPages) {
        return lastPage.page + 1;
      }
      return undefined;
    },
    enabled: Boolean(providerId),
  });
}
