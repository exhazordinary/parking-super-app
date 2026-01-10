/**
 * Wallet types
 */

import { Timestamps, PaginationParams } from './common';

export interface Wallet extends Timestamps {
  id: string;
  userId: string;
  balance: number;
  currency: string;
  status: WalletStatus;
}

export type WalletStatus = 'active' | 'suspended' | 'closed';

export interface Transaction extends Timestamps {
  id: string;
  walletId: string;
  type: TransactionType;
  amount: number;
  currency: string;
  status: TransactionStatus;
  reference?: string;
  description: string;
  metadata?: TransactionMetadata;
}

export type TransactionType =
  | 'topup'
  | 'payment'
  | 'refund'
  | 'transfer_in'
  | 'transfer_out';

export type TransactionStatus =
  | 'pending'
  | 'processing'
  | 'completed'
  | 'failed'
  | 'cancelled';

export interface TransactionMetadata {
  parkingSessionId?: string;
  providerId?: string;
  locationId?: string;
  vehicleId?: string;
  paymentMethod?: string;
  receiptUrl?: string;
}

export interface TopUpRequest {
  amount: number;
  paymentMethod: PaymentMethod;
  idempotencyKey: string;
}

export interface TopUpResponse {
  transaction: Transaction;
  paymentUrl?: string;
  paymentReference?: string;
}

export type PaymentMethod =
  | 'fpx'
  | 'card'
  | 'ewallet_tng'
  | 'ewallet_boost'
  | 'ewallet_grabpay';

export interface PaymentMethodInfo {
  type: PaymentMethod;
  name: string;
  icon: string;
  enabled: boolean;
  minAmount: number;
  maxAmount: number;
}

export interface GetTransactionsParams extends PaginationParams {
  type?: TransactionType;
  status?: TransactionStatus;
  fromDate?: string;
  toDate?: string;
}
