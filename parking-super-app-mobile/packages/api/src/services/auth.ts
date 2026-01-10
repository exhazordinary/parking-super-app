import { apiClient } from '../client';
import type {
  ApiResponse,
  User,
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  RegisterResponse,
  OTPRequestPayload,
  OTPVerifyPayload,
  OTPVerifyResponse,
  UpdateProfileRequest,
  ChangePasswordRequest,
} from '../types';

const AUTH_BASE = '/api/v1/auth';

export const authService = {
  /**
   * Login with phone and password
   */
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await apiClient.post<ApiResponse<LoginResponse>>(
      `${AUTH_BASE}/login`,
      data
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Register a new user
   */
  register: async (data: RegisterRequest): Promise<RegisterResponse> => {
    const response = await apiClient.post<ApiResponse<RegisterResponse>>(
      `${AUTH_BASE}/register`,
      data
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Request OTP
   */
  requestOTP: async (data: OTPRequestPayload): Promise<{ sent: boolean }> => {
    const response = await apiClient.post<ApiResponse<{ sent: boolean }>>(
      `${AUTH_BASE}/otp/request`,
      data
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Verify OTP
   */
  verifyOTP: async (data: OTPVerifyPayload): Promise<OTPVerifyResponse> => {
    const response = await apiClient.post<ApiResponse<OTPVerifyResponse>>(
      `${AUTH_BASE}/otp/verify`,
      data
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Get current user profile
   */
  getProfile: async (): Promise<User> => {
    const response = await apiClient.get<ApiResponse<User>>(
      `${AUTH_BASE}/profile`
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Update user profile
   */
  updateProfile: async (data: UpdateProfileRequest): Promise<User> => {
    const response = await apiClient.put<ApiResponse<User>>(
      `${AUTH_BASE}/profile`,
      data
    );
    if (!response.data.success || !response.data.data) {
      throw response.data.error;
    }
    return response.data.data;
  },

  /**
   * Change password
   */
  changePassword: async (data: ChangePasswordRequest): Promise<{ success: boolean }> => {
    const response = await apiClient.post<ApiResponse<{ success: boolean }>>(
      `${AUTH_BASE}/password/change`,
      data
    );
    if (!response.data.success) {
      throw response.data.error;
    }
    return { success: true };
  },

  /**
   * Logout
   */
  logout: async (): Promise<void> => {
    await apiClient.post(`${AUTH_BASE}/logout`);
  },
};
