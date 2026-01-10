import { apiClient } from '../client';
import type {
  ApiResponse,
  PaginatedResponse,
  Wallet,
  Transaction,
  TopUpRequest,
  TopUpResponse,
  PaymentMethodInfo,
  GetTransactionsParams,
} from '../types';

const WALLET_BASE = '/api/v1/wallet';

export const walletService = {
  /**
   * Get wallet balance and info
   */
  getWallet: async (): Promise<Wallet> => {
    const response = await apiClient.get<ApiResponse<Wallet>>(WALLET_BASE);
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Top up wallet
   */
  topUp: async (data: TopUpRequest): Promise<TopUpResponse> => {
    const response = await apiClient.post<ApiResponse<TopUpResponse>>(
      `${WALLET_BASE}/topup`,
      data,
      {
        headers: {
          'Idempotency-Key': data.idempotencyKey,
        },
      }
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Get transaction history
   */
  getTransactions: async (
    params?: GetTransactionsParams
  ): Promise<PaginatedResponse<Transaction>> => {
    const response = await apiClient.get<ApiResponse<PaginatedResponse<Transaction>>>(
      `${WALLET_BASE}/transactions`,
      { params }
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Get single transaction
   */
  getTransaction: async (id: string): Promise<Transaction> => {
    const response = await apiClient.get<ApiResponse<Transaction>>(
      `${WALLET_BASE}/transactions/${id}`
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Get available payment methods
   */
  getPaymentMethods: async (): Promise<PaymentMethodInfo[]> => {
    const response = await apiClient.get<ApiResponse<PaymentMethodInfo[]>>(
      `${WALLET_BASE}/payment-methods`
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },
};
