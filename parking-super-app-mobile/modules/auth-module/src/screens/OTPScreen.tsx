import React, { useState, useCallback } from 'react';
import { View, StyleSheet } from 'react-native';
import { Text, useTheme, IconButton } from 'react-native-paper';
import { SafeAreaView } from 'react-native-safe-area-context';
import type { MD3Theme } from 'react-native-paper';
import { useNavigation, useRoute } from '@react-navigation/native';

import { Button, LoadingOverlay, spacing } from '@parking/ui';
import { useVerifyOTP, useRequestOTP, useRegister } from '@parking/api';
import { useAuthStore } from '@parking/auth';
import type { AuthStackScreenProps } from '@parking/navigation';

import { OTPInput, ResendTimer } from '../components';

type RouteProp = AuthStackScreenProps<'OTP'>['route'];
type NavigationProp = AuthStackScreenProps<'OTP'>['navigation'];

export function OTPScreen(): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const navigation = useNavigation<NavigationProp>();
  const route = useRoute<RouteProp>();

  const { phone, type } = route.params;
  const { login: storeLogin } = useAuthStore();

  const verifyOTP = useVerifyOTP();
  const requestOTP = useRequestOTP();
  const register = useRegister();

  const [otp, setOtp] = useState('');
  const [error, setError] = useState('');

  const handleVerify = useCallback(
    async (code: string) => {
      setError('');

      try {
        const result = await verifyOTP.mutateAsync({
          phone,
          code,
          type,
        });

        if (result.verified) {
          if (type === 'register' && result.tokens) {
            // Complete registration
            // In a real app, you'd pass the registration data from previous screen
            // For now, we assume the backend handles it via OTP verification
            await storeLogin(
              { id: '', phone, name: '', isVerified: true, createdAt: '', updatedAt: '' },
              result.tokens
            );
          } else if (type === 'login' && result.tokens) {
            await storeLogin(
              { id: '', phone, name: '', isVerified: true, createdAt: '', updatedAt: '' },
              result.tokens
            );
          }
          // Navigation handled by auth state change
        } else {
          setError('Invalid verification code');
        }
      } catch (err) {
        setError('Verification failed. Please try again.');
      }
    },
    [phone, type, verifyOTP, storeLogin]
  );

  const handleResend = useCallback(async () => {
    setError('');
    setOtp('');

    try {
      await requestOTP.mutateAsync({
        phone,
        type,
      });
    } catch {
      setError('Failed to resend code. Please try again.');
    }
  }, [phone, type, requestOTP]);

  const handleBack = useCallback(() => {
    navigation.goBack();
  }, [navigation]);

  // Format phone for display
  const displayPhone = phone.replace('+60', '+60 ').replace(/(\d{2})(\d{3})(\d{4})/, '$1 $2 $3');

  return (
    <SafeAreaView style={styles.safeArea}>
      <View style={styles.header}>
        <IconButton
          icon="arrow-left"
          size={24}
          onPress={handleBack}
          style={styles.backButton}
        />
      </View>

      <View style={styles.content}>
        <Text variant="headlineMedium" style={styles.title}>
          Verify Phone
        </Text>
        <Text
          variant="bodyLarge"
          style={[styles.subtitle, { color: theme.colors.onSurfaceVariant }]}
        >
          Enter the 6-digit code sent to
        </Text>
        <Text variant="titleMedium" style={styles.phoneText}>
          {displayPhone}
        </Text>

        <View style={styles.otpContainer}>
          <OTPInput
            value={otp}
            onChangeText={setOtp}
            onComplete={handleVerify}
            error={Boolean(error)}
            helperText={error}
            disabled={verifyOTP.isPending}
          />

          <ResendTimer
            seconds={60}
            onResend={handleResend}
            isLoading={requestOTP.isPending}
          />
        </View>

        <Button
          mode="contained"
          onPress={() => handleVerify(otp)}
          loading={verifyOTP.isPending}
          disabled={otp.length !== 6 || verifyOTP.isPending}
          fullWidth
          style={styles.verifyButton}
        >
          Verify
        </Button>
      </View>

      <LoadingOverlay
        visible={verifyOTP.isPending}
        message="Verifying..."
      />
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
    backgroundColor: '#FFFFFF',
  },
  header: {
    paddingHorizontal: spacing.sm,
  },
  backButton: {
    marginLeft: -spacing.sm,
  },
  content: {
    flex: 1,
    paddingHorizontal: spacing.lg,
    paddingTop: spacing.xl,
  },
  title: {
    fontWeight: '700',
    textAlign: 'center',
    marginBottom: spacing.sm,
  },
  subtitle: {
    textAlign: 'center',
  },
  phoneText: {
    textAlign: 'center',
    fontWeight: '600',
    marginTop: spacing.xs,
    marginBottom: spacing.xxl,
  },
  otpContainer: {
    alignItems: 'center',
    marginBottom: spacing.xl,
  },
  verifyButton: {
    marginTop: 'auto',
    marginBottom: spacing.lg,
  },
});
