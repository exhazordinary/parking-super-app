/**
 * Parking types
 */

import { Timestamps, PaginationParams } from './common';

export interface ParkingSession extends Timestamps {
  id: string;
  userId: string;
  vehicleId: string;
  locationId: string;
  providerId: string;
  startTime: string;
  endTime?: string;
  status: ParkingSessionStatus;
  rate: ParkingRate;
  duration?: number; // in minutes
  totalCost?: number;
  currency: string;
  vehicle?: Vehicle;
  location?: ParkingLocation;
}

export type ParkingSessionStatus =
  | 'active'
  | 'completed'
  | 'cancelled'
  | 'expired';

export interface ParkingRate {
  type: RateType;
  baseRate: number;
  currency: string;
  firstHourRate?: number;
  subsequentHourRate?: number;
  dailyMaxRate?: number;
  graceMinutes: number;
}

export type RateType = 'hourly' | 'daily' | 'flat' | 'tiered';

export interface Vehicle extends Timestamps {
  id: string;
  userId: string;
  plateNumber: string;
  make?: string;
  model?: string;
  color?: string;
  type: VehicleType;
  isDefault: boolean;
}

export type VehicleType = 'car' | 'motorcycle' | 'truck' | 'van';

export interface ParkingLocation extends Timestamps {
  id: string;
  providerId: string;
  name: string;
  address: string;
  city: string;
  state: string;
  postalCode: string;
  latitude: number;
  longitude: number;
  totalSpaces: number;
  availableSpaces: number;
  operatingHours: OperatingHours;
  rates: ParkingRate[];
  amenities: string[];
  images?: string[];
  status: LocationStatus;
  distance?: number; // in meters, calculated by API
}

export type LocationStatus = 'open' | 'closed' | 'full' | 'maintenance';

export interface OperatingHours {
  monday: DayHours;
  tuesday: DayHours;
  wednesday: DayHours;
  thursday: DayHours;
  friday: DayHours;
  saturday: DayHours;
  sunday: DayHours;
}

export interface DayHours {
  open: string; // HH:mm format
  close: string;
  isOpen: boolean;
}

export interface StartSessionRequest {
  vehicleId: string;
  locationId: string;
  plateNumber?: string; // for quick start without saved vehicle
}

export interface StartSessionResponse {
  session: ParkingSession;
}

export interface EndSessionRequest {
  sessionId: string;
  paymentMethod?: string;
}

export interface EndSessionResponse {
  session: ParkingSession;
  transaction: {
    id: string;
    amount: number;
    status: string;
  };
}

export interface CreateVehicleRequest {
  plateNumber: string;
  make?: string;
  model?: string;
  color?: string;
  type: VehicleType;
  isDefault?: boolean;
}

export interface UpdateVehicleRequest {
  make?: string;
  model?: string;
  color?: string;
  isDefault?: boolean;
}

export interface SearchLocationsParams {
  latitude: number;
  longitude: number;
  radius?: number; // in meters, default 5000
  vehicleType?: VehicleType;
}

export interface GetSessionsParams extends PaginationParams {
  status?: ParkingSessionStatus;
  fromDate?: string;
  toDate?: string;
  vehicleId?: string;
}
