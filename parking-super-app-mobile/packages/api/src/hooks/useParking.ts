import {
  useQuery,
  useMutation,
  useQueryClient,
  useInfiniteQuery,
} from '@tanstack/react-query';
import { parkingService } from '../services/parking';
import { queryKeys } from '../queryKeys';
import type {
  StartSessionRequest,
  EndSessionRequest,
  CreateVehicleRequest,
  UpdateVehicleRequest,
  SearchLocationsParams,
  GetSessionsParams,
} from '../types';

// Session hooks

/**
 * Hook to get active parking session
 */
export function useActiveSession() {
  return useQuery({
    queryKey: queryKeys.parking.activeSession(),
    queryFn: () => parkingService.getActiveSession(),
    refetchInterval: 1000 * 60, // Refetch every minute for live updates
  });
}

/**
 * Hook to get session by ID
 */
export function useSession(id: string) {
  return useQuery({
    queryKey: queryKeys.parking.session(id),
    queryFn: () => parkingService.getSession(id),
    enabled: Boolean(id),
  });
}

/**
 * Hook to get session history with infinite scroll
 */
export function useSessions(params?: Omit<GetSessionsParams, 'page'>) {
  return useInfiniteQuery({
    queryKey: queryKeys.parking.sessions(params),
    queryFn: ({ pageParam = 1 }) =>
      parkingService.getSessions({ ...params, page: pageParam }),
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
 * Hook to start a parking session
 */
export function useStartSession() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: StartSessionRequest) => parkingService.startSession(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.parking.activeSession() });
      queryClient.invalidateQueries({ queryKey: queryKeys.parking.sessions() });
    },
  });
}

/**
 * Hook to end a parking session
 */
export function useEndSession() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: EndSessionRequest) => parkingService.endSession(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.parking.activeSession() });
      queryClient.invalidateQueries({ queryKey: queryKeys.parking.sessions() });
      queryClient.invalidateQueries({ queryKey: queryKeys.wallet.balance() });
      queryClient.invalidateQueries({ queryKey: queryKeys.wallet.transactions() });
    },
  });
}

// Location hooks

/**
 * Hook to search nearby locations
 */
export function useSearchLocations(params: SearchLocationsParams) {
  return useQuery({
    queryKey: queryKeys.parking.locations(params),
    queryFn: () => parkingService.searchLocations(params),
    enabled: Boolean(params.latitude && params.longitude),
  });
}

/**
 * Hook to get location by ID
 */
export function useLocation(id: string) {
  return useQuery({
    queryKey: queryKeys.parking.location(id),
    queryFn: () => parkingService.getLocation(id),
    enabled: Boolean(id),
  });
}

// Vehicle hooks

/**
 * Hook to get user's vehicles
 */
export function useVehicles() {
  return useQuery({
    queryKey: queryKeys.vehicles.list(),
    queryFn: () => parkingService.getVehicles(),
  });
}

/**
 * Hook to get vehicle by ID
 */
export function useVehicle(id: string) {
  return useQuery({
    queryKey: queryKeys.vehicles.vehicle(id),
    queryFn: () => parkingService.getVehicle(id),
    enabled: Boolean(id),
  });
}

/**
 * Hook to create a vehicle
 */
export function useCreateVehicle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateVehicleRequest) => parkingService.createVehicle(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.vehicles.list() });
    },
  });
}

/**
 * Hook to update a vehicle
 */
export function useUpdateVehicle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateVehicleRequest }) =>
      parkingService.updateVehicle(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.vehicles.list() });
      queryClient.invalidateQueries({ queryKey: queryKeys.vehicles.vehicle(id) });
    },
  });
}

/**
 * Hook to delete a vehicle
 */
export function useDeleteVehicle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => parkingService.deleteVehicle(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.vehicles.list() });
    },
  });
}
