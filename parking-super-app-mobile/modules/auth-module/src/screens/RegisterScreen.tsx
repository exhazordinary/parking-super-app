import React, { useState, useCallback, useRef } from 'react';
import {
  View,
  StyleSheet,
  ScrollView,
  KeyboardAvoidingView,
  Platform,
  TextInput,
} from 'react-native';
import { Text, useTheme } from 'react-native-paper';
import { SafeAreaView } from 'react-native-safe-area-context';
import type { MD3Theme } from 'react-native-paper';
import { useNavigation } from '@react-navigation/native';

import { Button, Input, spacing } from '@parking/ui';
import { useRequestOTP } from '@parking/api';
import type { AuthStackScreenProps } from '@parking/navigation';

import { PhoneInput, validateMalaysianPhone, formatPhoneForApi } from '../components';

type NavigationProp = AuthStackScreenProps<'Register'>['navigation'];

export function RegisterScreen(): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const navigation = useNavigation<NavigationProp>();
  const requestOTP = useRequestOTP();

  const [name, setName] = useState('');
  const [phone, setPhone] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);

  const [errors, setErrors] = useState<Record<string, string>>({});

  const phoneRef = useRef<TextInput>(null);
  const emailRef = useRef<TextInput>(null);
  const passwordRef = useRef<TextInput>(null);
  const confirmPasswordRef = useRef<TextInput>(null);

  const validate = useCallback((): boolean => {
    const newErrors: Record<string, string> = {};

    if (!name.trim()) {
      newErrors.name = 'Name is required';
    }

    if (!phone) {
      newErrors.phone = 'Phone number is required';
    } else if (!validateMalaysianPhone(phone)) {
      newErrors.phone = 'Please enter a valid Malaysian phone number';
    }

    if (email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      newErrors.email = 'Please enter a valid email address';
    }

    if (!password) {
      newErrors.password = 'Password is required';
    } else if (password.length < 8) {
      newErrors.password = 'Password must be at least 8 characters';
    }

    if (password !== confirmPassword) {
      newErrors.confirmPassword = 'Passwords do not match';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  }, [name, phone, email, password, confirmPassword]);

  const handleRegister = useCallback(async () => {
    if (!validate()) return;

    try {
      // Request OTP for verification
      await requestOTP.mutateAsync({
        phone: formatPhoneForApi(phone),
        type: 'register',
      });

      // Navigate to OTP screen with registration data
      navigation.navigate('OTP', {
        phone: formatPhoneForApi(phone),
        type: 'register',
      });
    } catch (error) {
      setErrors({
        general: 'Failed to send verification code. Please try again.',
      });
    }
  }, [validate, requestOTP, phone, navigation]);

  const handleLogin = useCallback(() => {
    navigation.navigate('Login');
  }, [navigation]);

  return (
    <SafeAreaView style={styles.safeArea}>
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        style={styles.keyboardView}
      >
        <ScrollView
          contentContainerStyle={styles.scrollContent}
          keyboardShouldPersistTaps="handled"
        >
          <View style={styles.header}>
            <Text variant="headlineLarge" style={styles.title}>
              Create Account
            </Text>
            <Text
              variant="bodyLarge"
              style={[styles.subtitle, { color: theme.colors.onSurfaceVariant }]}
            >
              Sign up to get started with Parking App
            </Text>
          </View>

          <View style={styles.form}>
            <Input
              label="Full Name"
              value={name}
              onChangeText={setName}
              error={Boolean(errors.name)}
              helperText={errors.name}
              autoCapitalize="words"
              onSubmitEditing={() => phoneRef.current?.focus()}
              returnKeyType="next"
            />

            <Text variant="labelLarge" style={styles.label}>
              Phone Number
            </Text>
            <PhoneInput
              ref={phoneRef}
              value={phone}
              onChangeText={setPhone}
              error={Boolean(errors.phone)}
              helperText={errors.phone}
              onSubmitEditing={() => emailRef.current?.focus()}
            />

            <Input
              ref={emailRef}
              label="Email (Optional)"
              value={email}
              onChangeText={setEmail}
              keyboardType="email-address"
              autoCapitalize="none"
              error={Boolean(errors.email)}
              helperText={errors.email}
              onSubmitEditing={() => passwordRef.current?.focus()}
              returnKeyType="next"
            />

            <Input
              ref={passwordRef}
              label="Password"
              value={password}
              onChangeText={setPassword}
              secureTextEntry={!showPassword}
              error={Boolean(errors.password)}
              helperText={errors.password || 'At least 8 characters'}
              right={
                <Input.Icon
                  icon={showPassword ? 'eye-off' : 'eye'}
                  onPress={() => setShowPassword(!showPassword)}
                />
              }
              onSubmitEditing={() => confirmPasswordRef.current?.focus()}
              returnKeyType="next"
            />

            <Input
              ref={confirmPasswordRef}
              label="Confirm Password"
              value={confirmPassword}
              onChangeText={setConfirmPassword}
              secureTextEntry={!showPassword}
              error={Boolean(errors.confirmPassword)}
              helperText={errors.confirmPassword}
              onSubmitEditing={handleRegister}
              returnKeyType="done"
            />

            {errors.general && (
              <Text
                variant="bodySmall"
                style={[styles.errorText, { color: theme.colors.error }]}
              >
                {errors.general}
              </Text>
            )}

            <Button
              mode="contained"
              onPress={handleRegister}
              loading={requestOTP.isPending}
              disabled={requestOTP.isPending}
              fullWidth
              style={styles.registerButton}
            >
              Continue
            </Button>

            <Text
              variant="bodySmall"
              style={[styles.termsText, { color: theme.colors.onSurfaceVariant }]}
            >
              By signing up, you agree to our Terms of Service and Privacy Policy
            </Text>
          </View>

          <View style={styles.footer}>
            <Text variant="bodyMedium" style={{ color: theme.colors.onSurfaceVariant }}>
              Already have an account?{' '}
            </Text>
            <Button mode="text" onPress={handleLogin} compact>
              Sign In
            </Button>
          </View>
        </ScrollView>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
    backgroundColor: '#FFFFFF',
  },
  keyboardView: {
    flex: 1,
  },
  scrollContent: {
    flexGrow: 1,
    padding: spacing.lg,
  },
  header: {
    marginTop: spacing.xl,
    marginBottom: spacing.lg,
  },
  title: {
    fontWeight: '700',
    marginBottom: spacing.sm,
  },
  subtitle: {
    lineHeight: 24,
  },
  form: {
    flex: 1,
  },
  label: {
    marginBottom: spacing.sm,
  },
  registerButton: {
    marginTop: spacing.lg,
  },
  termsText: {
    textAlign: 'center',
    marginTop: spacing.lg,
    lineHeight: 20,
  },
  errorText: {
    marginTop: spacing.sm,
    textAlign: 'center',
  },
  footer: {
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    paddingVertical: spacing.lg,
  },
});
