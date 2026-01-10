import { useQuery, useMutation, useQueryClient, useInfiniteQuery } from '@tanstack/react-query';
import { walletService } from '../services/wallet';
import { queryKeys } from '../queryKeys';
import type { TopUpRequest, GetTransactionsParams } from '../types';

/**
 * Hook to get wallet balance
 */
export function useWallet() {
  return useQuery({
    queryKey: queryKeys.wallet.balance(),
    queryFn: () => walletService.getWallet(),
  });
}

/**
 * Hook to top up wallet
 */
export function useTopUp() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: TopUpRequest) => walletService.topUp(data),
    onSuccess: () => {
      // Invalidate wallet balance and transactions
      queryClient.invalidateQueries({ queryKey: queryKeys.wallet.balance() });
      queryClient.invalidateQueries({ queryKey: queryKeys.wallet.transactions() });
    },
  });
}

/**
 * Hook to get transactions with infinite scroll
 */
export function useTransactions(params?: Omit<GetTransactionsParams, 'page'>) {
  return useInfiniteQuery({
    queryKey: queryKeys.wallet.transactions(params),
    queryFn: ({ pageParam = 1 }) =>
      walletService.getTransactions({ ...params, page: pageParam }),
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
 * Hook to get a single transaction
 */
export function useTransaction(id: string) {
  return useQuery({
    queryKey: queryKeys.wallet.transaction(id),
    queryFn: () => walletService.getTransaction(id),
    enabled: Boolean(id),
  });
}

/**
 * Hook to get payment methods
 */
export function usePaymentMethods() {
  return useQuery({
    queryKey: queryKeys.wallet.paymentMethods(),
    queryFn: () => walletService.getPaymentMethods(),
    staleTime: 1000 * 60 * 30, // 30 minutes
  });
}
