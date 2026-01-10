import { apiClient } from '../client';
import type {
  ApiResponse,
  PaginatedResponse,
  Notification,
  NotificationPreferences,
  GetNotificationsParams,
  UpdatePreferencesRequest,
  MarkReadRequest,
  RegisterDeviceRequest,
} from '../types';

const NOTIFICATIONS_BASE = '/api/v1/notifications';

export const notificationService = {
  /**
   * Get notifications
   */
  getNotifications: async (
    params?: GetNotificationsParams
  ): Promise<PaginatedResponse<Notification>> => {
    const response = await apiClient.get<ApiResponse<PaginatedResponse<Notification>>>(
      NOTIFICATIONS_BASE,
      { params }
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Get unread count
   */
  getUnreadCount: async (): Promise<{ count: number }> => {
    const response = await apiClient.get<ApiResponse<{ count: number }>>(
      `${NOTIFICATIONS_BASE}/unread-count`
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Mark notifications as read
   */
  markAsRead: async (data: MarkReadRequest): Promise<void> => {
    await apiClient.post(`${NOTIFICATIONS_BASE}/mark-read`, data);
  },

  /**
   * Mark all as read
   */
  markAllAsRead: async (): Promise<void> => {
    await apiClient.post(`${NOTIFICATIONS_BASE}/mark-all-read`);
  },

  /**
   * Get notification preferences
   */
  getPreferences: async (): Promise<NotificationPreferences> => {
    const response = await apiClient.get<ApiResponse<NotificationPreferences>>(
      `${NOTIFICATIONS_BASE}/preferences`
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Update notification preferences
   */
  updatePreferences: async (
    data: UpdatePreferencesRequest
  ): Promise<NotificationPreferences> => {
    const response = await apiClient.put<ApiResponse<NotificationPreferences>>(
      `${NOTIFICATIONS_BASE}/preferences`,
      data
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Register device for push notifications
   */
  registerDevice: async (data: RegisterDeviceRequest): Promise<void> => {
    await apiClient.post(`${NOTIFICATIONS_BASE}/devices`, data);
  },

  /**
   * Unregister device
   */
  unregisterDevice: async (deviceId: string): Promise<void> => {
    await apiClient.delete(`${NOTIFICATIONS_BASE}/devices/${deviceId}`);
  },
};
