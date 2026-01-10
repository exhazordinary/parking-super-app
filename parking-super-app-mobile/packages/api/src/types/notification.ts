/**
 * Notification types
 */

import { Timestamps, PaginationParams } from './common';

export interface Notification extends Timestamps {
  id: string;
  userId: string;
  type: NotificationType;
  title: string;
  message: string;
  data?: NotificationData;
  isRead: boolean;
  readAt?: string;
}

export type NotificationType =
  | 'parking_reminder'
  | 'parking_expiring'
  | 'parking_ended'
  | 'payment_success'
  | 'payment_failed'
  | 'topup_success'
  | 'low_balance'
  | 'promo'
  | 'system';

export interface NotificationData {
  sessionId?: string;
  transactionId?: string;
  locationId?: string;
  actionUrl?: string;
}

export interface NotificationPreferences {
  pushEnabled: boolean;
  emailEnabled: boolean;
  smsEnabled: boolean;
  parkingReminders: boolean;
  promotions: boolean;
  systemUpdates: boolean;
  reminderMinutes: number; // minutes before session expires
}

export interface UpdatePreferencesRequest {
  pushEnabled?: boolean;
  emailEnabled?: boolean;
  smsEnabled?: boolean;
  parkingReminders?: boolean;
  promotions?: boolean;
  systemUpdates?: boolean;
  reminderMinutes?: number;
}

export interface GetNotificationsParams extends PaginationParams {
  type?: NotificationType;
  isRead?: boolean;
}

export interface MarkReadRequest {
  notificationIds: string[];
}

export interface RegisterDeviceRequest {
  token: string;
  platform: 'ios' | 'android';
  deviceId: string;
}
