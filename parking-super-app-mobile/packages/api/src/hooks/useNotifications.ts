import {
  useQuery,
  useMutation,
  useQueryClient,
  useInfiniteQuery,
} from '@tanstack/react-query';
import { notificationService } from '../services/notification';
import { queryKeys } from '../queryKeys';
import type {
  GetNotificationsParams,
  UpdatePreferencesRequest,
  RegisterDeviceRequest,
} from '../types';

/**
 * Hook to get notifications with infinite scroll
 */
export function useNotifications(params?: Omit<GetNotificationsParams, 'page'>) {
  return useInfiniteQuery({
    queryKey: queryKeys.notifications.list(params),
    queryFn: ({ pageParam = 1 }) =>
      notificationService.getNotifications({ ...params, page: pageParam }),
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
 * Hook to get unread notification count
 */
export function useUnreadCount() {
  return useQuery({
    queryKey: queryKeys.notifications.unreadCount(),
    queryFn: () => notificationService.getUnreadCount(),
    refetchInterval: 1000 * 60, // Refetch every minute
  });
}

/**
 * Hook to mark notifications as read
 */
export function useMarkAsRead() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (notificationIds: string[]) =>
      notificationService.markAsRead({ notificationIds }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.notifications.list() });
      queryClient.invalidateQueries({ queryKey: queryKeys.notifications.unreadCount() });
    },
  });
}

/**
 * Hook to mark all notifications as read
 */
export function useMarkAllAsRead() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => notificationService.markAllAsRead(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.notifications.list() });
      queryClient.invalidateQueries({ queryKey: queryKeys.notifications.unreadCount() });
    },
  });
}

/**
 * Hook to get notification preferences
 */
export function useNotificationPreferences() {
  return useQuery({
    queryKey: queryKeys.notifications.preferences(),
    queryFn: () => notificationService.getPreferences(),
  });
}

/**
 * Hook to update notification preferences
 */
export function useUpdateNotificationPreferences() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: UpdatePreferencesRequest) =>
      notificationService.updatePreferences(data),
    onSuccess: (data) => {
      queryClient.setQueryData(queryKeys.notifications.preferences(), data);
    },
  });
}

/**
 * Hook to register device for push notifications
 */
export function useRegisterDevice() {
  return useMutation({
    mutationFn: (data: RegisterDeviceRequest) =>
      notificationService.registerDevice(data),
  });
}

/**
 * Hook to unregister device
 */
export function useUnregisterDevice() {
  return useMutation({
    mutationFn: (deviceId: string) => notificationService.unregisterDevice(deviceId),
  });
}
