import React, { useRef, useState, useCallback, useEffect } from 'react';
import {
  View,
  TextInput,
  StyleSheet,
  Pressable,
  Keyboard,
} from 'react-native';
import { Text, useTheme } from 'react-native-paper';
import type { MD3Theme } from 'react-native-paper';

interface OTPInputProps {
  value: string;
  onChangeText: (value: string) => void;
  length?: number;
  error?: boolean;
  helperText?: string;
  disabled?: boolean;
  onComplete?: (code: string) => void;
  autoFocus?: boolean;
}

/**
 * OTP input with individual digit boxes
 */
export function OTPInput({
  value,
  onChangeText,
  length = 6,
  error = false,
  helperText,
  disabled = false,
  onComplete,
  autoFocus = true,
}: OTPInputProps): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const inputRef = useRef<TextInput>(null);
  const [isFocused, setIsFocused] = useState(false);

  // Focus the hidden input when any box is pressed
  const handlePress = useCallback(() => {
    inputRef.current?.focus();
  }, []);

  // Handle text change
  const handleChange = useCallback(
    (text: string) => {
      // Only allow digits
      const digits = text.replace(/\D/g, '').slice(0, length);
      onChangeText(digits);

      // Call onComplete when all digits are entered
      if (digits.length === length && onComplete) {
        onComplete(digits);
      }
    },
    [length, onChangeText, onComplete]
  );

  // Auto-focus on mount
  useEffect(() => {
    if (autoFocus) {
      setTimeout(() => {
        inputRef.current?.focus();
      }, 100);
    }
  }, [autoFocus]);

  // Get the digits array for rendering boxes
  const digits = value.split('').slice(0, length);

  return (
    <View style={styles.container}>
      <Pressable onPress={handlePress} disabled={disabled}>
        <View style={styles.boxContainer}>
          {Array.from({ length }, (_, index) => {
            const digit = digits[index] || '';
            const isCurrentIndex = index === value.length;
            const isFilled = digit !== '';

            return (
              <View
                key={index}
                style={[
                  styles.box,
                  {
                    borderColor: error
                      ? theme.colors.error
                      : isFocused && isCurrentIndex
                      ? theme.colors.primary
                      : isFilled
                      ? theme.colors.primary
                      : theme.colors.outline,
                    backgroundColor: disabled
                      ? theme.colors.surfaceDisabled
                      : theme.colors.surface,
                  },
                ]}
              >
                <Text
                  variant="headlineMedium"
                  style={[
                    styles.digit,
                    { color: theme.colors.onSurface },
                  ]}
                >
                  {digit}
                </Text>
              </View>
            );
          })}
        </View>
      </Pressable>

      {/* Hidden input */}
      <TextInput
        ref={inputRef}
        value={value}
        onChangeText={handleChange}
        keyboardType="number-pad"
        maxLength={length}
        style={styles.hiddenInput}
        onFocus={() => setIsFocused(true)}
        onBlur={() => setIsFocused(false)}
        editable={!disabled}
        caretHidden
        autoComplete="one-time-code"
        textContentType="oneTimeCode"
      />

      {/* Helper text */}
      {helperText && (
        <Text
          variant="bodySmall"
          style={[
            styles.helperText,
            { color: error ? theme.colors.error : theme.colors.onSurfaceVariant },
          ]}
        >
          {helperText}
        </Text>
      )}
    </View>
  );
}

/**
 * Countdown timer for resend OTP
 */
interface ResendTimerProps {
  seconds: number;
  onResend: () => void;
  isLoading?: boolean;
}

export function ResendTimer({
  seconds: initialSeconds,
  onResend,
  isLoading = false,
}: ResendTimerProps): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const [seconds, setSeconds] = useState(initialSeconds);
  const [canResend, setCanResend] = useState(false);

  useEffect(() => {
    if (seconds > 0) {
      const timer = setTimeout(() => setSeconds(seconds - 1), 1000);
      return () => clearTimeout(timer);
    } else {
      setCanResend(true);
    }
  }, [seconds]);

  const handleResend = useCallback(() => {
    if (canResend && !isLoading) {
      onResend();
      setSeconds(initialSeconds);
      setCanResend(false);
    }
  }, [canResend, isLoading, onResend, initialSeconds]);

  return (
    <View style={styles.resendContainer}>
      {canResend ? (
        <Pressable onPress={handleResend} disabled={isLoading}>
          <Text
            style={[
              styles.resendText,
              { color: theme.colors.primary },
            ]}
          >
            {isLoading ? 'Sending...' : 'Resend OTP'}
          </Text>
        </Pressable>
      ) : (
        <Text
          style={[
            styles.resendText,
            { color: theme.colors.onSurfaceVariant },
          ]}
        >
          Resend in {seconds}s
        </Text>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    marginBottom: 24,
  },
  boxContainer: {
    flexDirection: 'row',
    justifyContent: 'center',
    gap: 8,
  },
  box: {
    width: 48,
    height: 56,
    borderWidth: 2,
    borderRadius: 12,
    justifyContent: 'center',
    alignItems: 'center',
  },
  digit: {
    fontWeight: '600',
  },
  hiddenInput: {
    position: 'absolute',
    opacity: 0,
    height: 0,
    width: 0,
  },
  helperText: {
    marginTop: 12,
    textAlign: 'center',
  },
  resendContainer: {
    alignItems: 'center',
    marginTop: 16,
  },
  resendText: {
    fontSize: 14,
    fontWeight: '500',
  },
});
