import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { authService } from '../services/auth';
import { queryKeys } from '../queryKeys';
import type {
  LoginRequest,
  RegisterRequest,
  OTPRequestPayload,
  OTPVerifyPayload,
  UpdateProfileRequest,
  ChangePasswordRequest,
} from '../types';

/**
 * Hook to get current user profile
 */
export function useProfile() {
  return useQuery({
    queryKey: queryKeys.auth.profile(),
    queryFn: () => authService.getProfile(),
  });
}

/**
 * Hook to login
 */
export function useLogin() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: LoginRequest) => authService.login(data),
    onSuccess: (data) => {
      queryClient.setQueryData(queryKeys.auth.profile(), data.user);
    },
  });
}

/**
 * Hook to register
 */
export function useRegister() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: RegisterRequest) => authService.register(data),
    onSuccess: (data) => {
      queryClient.setQueryData(queryKeys.auth.profile(), data.user);
    },
  });
}

/**
 * Hook to request OTP
 */
export function useRequestOTP() {
  return useMutation({
    mutationFn: (data: OTPRequestPayload) => authService.requestOTP(data),
  });
}

/**
 * Hook to verify OTP
 */
export function useVerifyOTP() {
  return useMutation({
    mutationFn: (data: OTPVerifyPayload) => authService.verifyOTP(data),
  });
}

/**
 * Hook to update profile
 */
export function useUpdateProfile() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: UpdateProfileRequest) => authService.updateProfile(data),
    onSuccess: (data) => {
      queryClient.setQueryData(queryKeys.auth.profile(), data);
    },
  });
}

/**
 * Hook to change password
 */
export function useChangePassword() {
  return useMutation({
    mutationFn: (data: ChangePasswordRequest) => authService.changePassword(data),
  });
}

/**
 * Hook to logout
 */
export function useLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => authService.logout(),
    onSuccess: () => {
      queryClient.clear();
    },
  });
}
