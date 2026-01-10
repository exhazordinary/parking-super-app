/**
 * Provider types
 */

import { Timestamps, PaginationParams } from './common';
import { ParkingLocation } from './parking';

export interface Provider extends Timestamps {
  id: string;
  name: string;
  description?: string;
  logo?: string;
  website?: string;
  phone?: string;
  email?: string;
  status: ProviderStatus;
  totalLocations: number;
  rating?: number;
  reviewCount?: number;
}

export type ProviderStatus = 'active' | 'inactive' | 'suspended';

export interface GetProvidersParams extends PaginationParams {
  status?: ProviderStatus;
  search?: string;
}

export interface GetProviderLocationsParams extends PaginationParams {
  status?: string;
  city?: string;
  latitude?: number;
  longitude?: number;
  radius?: number;
}

export interface ProviderWithLocations extends Provider {
  locations: ParkingLocation[];
}
