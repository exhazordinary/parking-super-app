import { useState, useEffect, useCallback } from 'react';
import NetInfo, { NetInfoState } from '@react-native-community/netinfo';
import { offlineQueue, QueuedOperation } from './queue';

/**
 * Hook to track network connectivity status
 */
export function useNetworkStatus(): { isConnected: boolean; isInternetReachable: boolean } {
  const [status, setStatus] = useState({ isConnected: true, isInternetReachable: true });

  useEffect(() => {
    const unsubscribe = NetInfo.addEventListener((state: NetInfoState) => {
      setStatus({
        isConnected: state.isConnected ?? false,
        isInternetReachable: state.isInternetReachable ?? false,
      });
    });
    return () => unsubscribe();
  }, []);

  return status;
}

/**
 * Hook to get pending offline operations count
 */
export function usePendingOperations(): { count: number; refresh: () => void } {
  const [count, setCount] = useState(0);

  const refresh = useCallback(async () => {
    const pending = await offlineQueue.getPendingCount();
    setCount(pending);
  }, []);

  useEffect(() => {
    refresh();
    const interval = setInterval(refresh, 5000);
    return () => clearInterval(interval);
  }, [refresh]);

  return { count, refresh };
}

/**
 * Hook to enqueue operations when offline
 */
export function useOfflineOperation<T extends Record<string, unknown>>(
  type: QueuedOperation['type'],
  onlineHandler: (payload: T) => Promise<unknown>,
  options?: { maxRetries?: number }
): {
  execute: (payload: T) => Promise<{ queued: boolean; result?: unknown }>;
  isOnline: boolean;
} {
  const { isConnected, isInternetReachable } = useNetworkStatus();
  const isOnline = isConnected && isInternetReachable;

  const execute = useCallback(
    async (payload: T): Promise<{ queued: boolean; result?: unknown }> => {
      if (isOnline) {
        try {
          const result = await onlineHandler(payload);
          return { queued: false, result };
        } catch (error) {
          // If online request fails, queue it
          await offlineQueue.enqueue({
            type,
            payload,
            maxRetries: options?.maxRetries ?? 3,
          });
          return { queued: true };
        }
      } else {
        await offlineQueue.enqueue({
          type,
          payload,
          maxRetries: options?.maxRetries ?? 3,
        });
        return { queued: true };
      }
    },
    [isOnline, onlineHandler, type, options?.maxRetries]
  );

  return { execute, isOnline };
}
