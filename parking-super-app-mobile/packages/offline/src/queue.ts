import AsyncStorage from '@react-native-async-storage/async-storage';
import NetInfo from '@react-native-community/netinfo';
import { v4 as uuidv4 } from 'uuid';

const QUEUE_KEY = '@parking/offline-queue';

export interface QueuedOperation {
  id: string;
  type: 'startSession' | 'endSession' | 'topUp' | 'createVehicle' | 'updateProfile';
  payload: Record<string, unknown>;
  timestamp: number;
  retries: number;
  maxRetries: number;
}

type OperationHandler = (payload: Record<string, unknown>) => Promise<unknown>;

class OfflineQueue {
  private handlers: Map<string, OperationHandler> = new Map();
  private isProcessing = false;
  private unsubscribeNetInfo: (() => void) | null = null;

  /**
   * Initialize the offline queue and start listening for network changes
   */
  async initialize(): Promise<void> {
    this.unsubscribeNetInfo = NetInfo.addEventListener((state) => {
      if (state.isConnected && state.isInternetReachable) {
        this.processQueue();
      }
    });
    await this.processQueue();
  }

  /**
   * Clean up listeners
   */
  destroy(): void {
    if (this.unsubscribeNetInfo) {
      this.unsubscribeNetInfo();
      this.unsubscribeNetInfo = null;
    }
  }

  /**
   * Register a handler for an operation type
   */
  registerHandler(type: QueuedOperation['type'], handler: OperationHandler): void {
    this.handlers.set(type, handler);
  }

  /**
   * Add an operation to the queue
   */
  async enqueue(operation: Omit<QueuedOperation, 'id' | 'timestamp' | 'retries'>): Promise<string> {
    const queue = await this.getQueue();
    const queuedOp: QueuedOperation = {
      ...operation,
      id: uuidv4(),
      timestamp: Date.now(),
      retries: 0,
    };
    queue.push(queuedOp);
    await this.saveQueue(queue);

    // Try to process immediately if online
    const netState = await NetInfo.fetch();
    if (netState.isConnected && netState.isInternetReachable) {
      this.processQueue();
    }

    return queuedOp.id;
  }

  /**
   * Remove an operation from the queue
   */
  async dequeue(id: string): Promise<void> {
    const queue = await this.getQueue();
    const filtered = queue.filter((op) => op.id !== id);
    await this.saveQueue(filtered);
  }

  /**
   * Get all queued operations
   */
  async getQueue(): Promise<QueuedOperation[]> {
    try {
      const data = await AsyncStorage.getItem(QUEUE_KEY);
      return data ? JSON.parse(data) : [];
    } catch {
      return [];
    }
  }

  /**
   * Get pending operation count
   */
  async getPendingCount(): Promise<number> {
    const queue = await this.getQueue();
    return queue.length;
  }

  /**
   * Process all queued operations
   */
  async processQueue(): Promise<void> {
    if (this.isProcessing) return;

    const netState = await NetInfo.fetch();
    if (!netState.isConnected || !netState.isInternetReachable) return;

    this.isProcessing = true;

    try {
      const queue = await this.getQueue();

      for (const operation of queue) {
        const handler = this.handlers.get(operation.type);
        if (!handler) {
          console.warn(`No handler for operation type: ${operation.type}`);
          continue;
        }

        try {
          await handler(operation.payload);
          await this.dequeue(operation.id);
        } catch (error) {
          operation.retries += 1;
          if (operation.retries >= operation.maxRetries) {
            await this.dequeue(operation.id);
            console.error(`Operation ${operation.id} failed after ${operation.maxRetries} retries`);
          } else {
            await this.saveQueue(queue);
          }
        }
      }
    } finally {
      this.isProcessing = false;
    }
  }

  private async saveQueue(queue: QueuedOperation[]): Promise<void> {
    await AsyncStorage.setItem(QUEUE_KEY, JSON.stringify(queue));
  }
}

export const offlineQueue = new OfflineQueue();
