import { apiClient } from '../client';
import type {
  ApiResponse,
  PaginatedResponse,
  ParkingSession,
  ParkingLocation,
  Vehicle,
  StartSessionRequest,
  StartSessionResponse,
  EndSessionRequest,
  EndSessionResponse,
  CreateVehicleRequest,
  UpdateVehicleRequest,
  SearchLocationsParams,
  GetSessionsParams,
} from '../types';

const PARKING_BASE = '/api/v1/parking';

export const parkingService = {
  // Sessions
  /**
   * Start a parking session
   */
  startSession: async (data: StartSessionRequest): Promise<StartSessionResponse> => {
    const response = await apiClient.post<ApiResponse<StartSessionResponse>>(
      `${PARKING_BASE}/sessions`,
      data
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * End a parking session
   */
  endSession: async (data: EndSessionRequest): Promise<EndSessionResponse> => {
    const response = await apiClient.post<ApiResponse<EndSessionResponse>>(
      `${PARKING_BASE}/sessions/${data.sessionId}/end`,
      { paymentMethod: data.paymentMethod }
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Get active session
   */
  getActiveSession: async (): Promise<ParkingSession | null> => {
    const response = await apiClient.get<ApiResponse<ParkingSession | null>>(
      `${PARKING_BASE}/sessions/active`
    );
    if (!response.data.success) {
      throw response.data.error;
    }
    return response.data.data ?? null;
  },

  /**
   * Get session by ID
   */
  getSession: async (id: string): Promise<ParkingSession> => {
    const response = await apiClient.get<ApiResponse<ParkingSession>>(
      `${PARKING_BASE}/sessions/${id}`
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Get session history
   */
  getSessions: async (
    params?: GetSessionsParams
  ): Promise<PaginatedResponse<ParkingSession>> => {
    const response = await apiClient.get<ApiResponse<PaginatedResponse<ParkingSession>>>(
      `${PARKING_BASE}/sessions`,
      { params }
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  // Locations
  /**
   * Search nearby parking locations
   */
  searchLocations: async (
    params: SearchLocationsParams
  ): Promise<ParkingLocation[]> => {
    const response = await apiClient.get<ApiResponse<ParkingLocation[]>>(
      `${PARKING_BASE}/locations`,
      { params }
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Get location by ID
   */
  getLocation: async (id: string): Promise<ParkingLocation> => {
    const response = await apiClient.get<ApiResponse<ParkingLocation>>(
      `${PARKING_BASE}/locations/${id}`
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  // Vehicles
  /**
   * Get user's vehicles
   */
  getVehicles: async (): Promise<Vehicle[]> => {
    const response = await apiClient.get<ApiResponse<Vehicle[]>>(
      `${PARKING_BASE}/vehicles`
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Get vehicle by ID
   */
  getVehicle: async (id: string): Promise<Vehicle> => {
    const response = await apiClient.get<ApiResponse<Vehicle>>(
      `${PARKING_BASE}/vehicles/${id}`
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Create a new vehicle
   */
  createVehicle: async (data: CreateVehicleRequest): Promise<Vehicle> => {
    const response = await apiClient.post<ApiResponse<Vehicle>>(
      `${PARKING_BASE}/vehicles`,
      data
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Update a vehicle
   */
  updateVehicle: async (id: string, data: UpdateVehicleRequest): Promise<Vehicle> => {
    const response = await apiClient.put<ApiResponse<Vehicle>>(
      `${PARKING_BASE}/vehicles/${id}`,
      data
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Delete a vehicle
   */
  deleteVehicle: async (id: string): Promise<void> => {
    await apiClient.delete(`${PARKING_BASE}/vehicles/${id}`);
  },
};
