/**
 * Authentication types
 */

import { Timestamps } from './common';

export interface User extends Timestamps {
  id: string;
  phone: string;
  name: string;
  email?: string;
  avatar?: string;
  isVerified: boolean;
  preferences?: UserPreferences;
}

export interface UserPreferences {
  notifications: boolean;
  pushNotifications: boolean;
  emailNotifications: boolean;
  smsNotifications: boolean;
  language: string;
  theme: 'light' | 'dark' | 'system';
}

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  expiresIn: number;
}

export interface LoginRequest {
  phone: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  tokens: AuthTokens;
}

export interface RegisterRequest {
  phone: string;
  name: string;
  password: string;
  email?: string;
}

export interface RegisterResponse {
  user: User;
  tokens: AuthTokens;
}

export interface OTPRequestPayload {
  phone: string;
  type: 'login' | 'register' | 'reset_password' | 'verify_phone';
}

export interface OTPVerifyPayload {
  phone: string;
  code: string;
  type: 'login' | 'register' | 'reset_password' | 'verify_phone';
}

export interface OTPVerifyResponse {
  verified: boolean;
  tokens?: AuthTokens;
}

export interface RefreshTokenRequest {
  refreshToken: string;
}

export interface UpdateProfileRequest {
  name?: string;
  email?: string;
  avatar?: string;
}

export interface ChangePasswordRequest {
  currentPassword: string;
  newPassword: string;
}
