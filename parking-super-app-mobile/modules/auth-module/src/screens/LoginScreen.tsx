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
import { useAuth } from '@parking/auth';
import type { AuthStackScreenProps } from '@parking/navigation';

import { PhoneInput, validateMalaysianPhone, formatPhoneForApi } from '../components';

type NavigationProp = AuthStackScreenProps<'Login'>['navigation'];

export function LoginScreen(): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const navigation = useNavigation<NavigationProp>();
  const { login, isLoggingIn, loginError } = useAuth();

  const [phone, setPhone] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [phoneError, setPhoneError] = useState('');
  const [passwordError, setPasswordError] = useState('');

  const passwordRef = useRef<TextInput>(null);

  const validate = useCallback((): boolean => {
    let isValid = true;

    if (!phone) {
      setPhoneError('Phone number is required');
      isValid = false;
    } else if (!validateMalaysianPhone(phone)) {
      setPhoneError('Please enter a valid Malaysian phone number');
      isValid = false;
    } else {
      setPhoneError('');
    }

    if (!password) {
      setPasswordError('Password is required');
      isValid = false;
    } else if (password.length < 6) {
      setPasswordError('Password must be at least 6 characters');
      isValid = false;
    } else {
      setPasswordError('');
    }

    return isValid;
  }, [phone, password]);

  const handleLogin = useCallback(async () => {
    if (!validate()) return;

    try {
      await login({
        phone: formatPhoneForApi(phone),
        password,
      });
      // Navigation handled by auth state change
    } catch (error) {
      // Error handled by useAuth
    }
  }, [validate, login, phone, password]);

  const handleRegister = useCallback(() => {
    navigation.navigate('Register');
  }, [navigation]);

  const handleForgotPassword = useCallback(() => {
    navigation.navigate('ForgotPassword');
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
              Welcome Back
            </Text>
            <Text
              variant="bodyLarge"
              style={[styles.subtitle, { color: theme.colors.onSurfaceVariant }]}
            >
              Sign in to continue to Parking App
            </Text>
          </View>

          <View style={styles.form}>
            <Text variant="labelLarge" style={styles.label}>
              Phone Number
            </Text>
            <PhoneInput
              value={phone}
              onChangeText={setPhone}
              error={Boolean(phoneError)}
              helperText={phoneError}
              onSubmitEditing={() => passwordRef.current?.focus()}
            />

            <Text variant="labelLarge" style={styles.label}>
              Password
            </Text>
            <Input
              ref={passwordRef}
              label=""
              value={password}
              onChangeText={setPassword}
              secureTextEntry={!showPassword}
              error={Boolean(passwordError)}
              helperText={passwordError}
              right={
                <Input.Icon
                  icon={showPassword ? 'eye-off' : 'eye'}
                  onPress={() => setShowPassword(!showPassword)}
                />
              }
              onSubmitEditing={handleLogin}
              returnKeyType="done"
            />

            {loginError && (
              <Text
                variant="bodySmall"
                style={[styles.errorText, { color: theme.colors.error }]}
              >
                {loginError.message || 'Login failed. Please try again.'}
              </Text>
            )}

            <Button
              mode="text"
              onPress={handleForgotPassword}
              style={styles.forgotButton}
            >
              Forgot Password?
            </Button>

            <Button
              mode="contained"
              onPress={handleLogin}
              loading={isLoggingIn}
              disabled={isLoggingIn}
              fullWidth
              style={styles.loginButton}
            >
              Sign In
            </Button>
          </View>

          <View style={styles.footer}>
            <Text variant="bodyMedium" style={{ color: theme.colors.onSurfaceVariant }}>
              Don't have an account?{' '}
            </Text>
            <Button mode="text" onPress={handleRegister} compact>
              Sign Up
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
    marginTop: spacing.xxl,
    marginBottom: spacing.xl,
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
  forgotButton: {
    alignSelf: 'flex-end',
    marginTop: spacing.sm,
  },
  loginButton: {
    marginTop: spacing.lg,
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
