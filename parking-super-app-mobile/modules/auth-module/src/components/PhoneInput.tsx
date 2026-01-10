import React, { useState, useCallback, forwardRef } from 'react';
import { View, StyleSheet, TextInput as RNTextInput } from 'react-native';
import { Text, useTheme } from 'react-native-paper';
import { Input } from '@parking/ui';
import type { MD3Theme } from 'react-native-paper';

interface PhoneInputProps {
  value: string;
  onChangeText: (value: string) => void;
  error?: boolean;
  helperText?: string;
  disabled?: boolean;
  onSubmitEditing?: () => void;
  autoFocus?: boolean;
}

/**
 * Malaysian phone number input with +60 prefix
 * Formats input as: +60 12 345 6789
 */
export const PhoneInput = forwardRef<RNTextInput, PhoneInputProps>(function PhoneInput(
  {
    value,
    onChangeText,
    error = false,
    helperText,
    disabled = false,
    onSubmitEditing,
    autoFocus = false,
  },
  ref
): React.JSX.Element {
  const theme = useTheme<MD3Theme>();
  const [isFocused, setIsFocused] = useState(false);

  // Format the display value (e.g., "12 345 6789")
  const formatDisplayValue = useCallback((raw: string): string => {
    // Remove all non-digits
    const digits = raw.replace(/\D/g, '');

    // Format with spaces: XX XXX XXXX
    if (digits.length <= 2) {
      return digits;
    } else if (digits.length <= 5) {
      return `${digits.slice(0, 2)} ${digits.slice(2)}`;
    } else {
      return `${digits.slice(0, 2)} ${digits.slice(2, 5)} ${digits.slice(5, 9)}`;
    }
  }, []);

  // Handle input change
  const handleChange = useCallback(
    (text: string) => {
      // Remove the +60 prefix if user tries to type it
      let cleanText = text.replace(/^\+60\s*/, '');
      // Remove all non-digits
      cleanText = cleanText.replace(/\D/g, '');
      // Limit to 9-10 digits (Malaysian mobile numbers)
      cleanText = cleanText.slice(0, 10);
      onChangeText(cleanText);
    },
    [onChangeText]
  );

  return (
    <View style={styles.container}>
      <View
        style={[
          styles.inputContainer,
          {
            borderColor: error
              ? theme.colors.error
              : isFocused
              ? theme.colors.primary
              : theme.colors.outline,
            backgroundColor: disabled
              ? theme.colors.surfaceDisabled
              : theme.colors.surface,
          },
        ]}
      >
        {/* Country code prefix */}
        <View style={styles.prefixContainer}>
          <Text
            style={[
              styles.prefix,
              { color: theme.colors.onSurfaceVariant },
            ]}
          >
            +60
          </Text>
        </View>

        {/* Phone number input */}
        <RNTextInput
          ref={ref}
          value={formatDisplayValue(value)}
          onChangeText={handleChange}
          style={[
            styles.input,
            { color: theme.colors.onSurface },
          ]}
          placeholder="12 345 6789"
          placeholderTextColor={theme.colors.onSurfaceVariant}
          keyboardType="phone-pad"
          maxLength={12} // XX XXX XXXX
          editable={!disabled}
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
          onSubmitEditing={onSubmitEditing}
          autoFocus={autoFocus}
        />
      </View>

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
});

/**
 * Validate Malaysian phone number
 */
export function validateMalaysianPhone(phone: string): boolean {
  // Remove non-digits
  const digits = phone.replace(/\D/g, '');
  // Malaysian mobile numbers: 9-10 digits starting with 1
  return /^1\d{8,9}$/.test(digits);
}

/**
 * Format phone for API (with country code)
 */
export function formatPhoneForApi(phone: string): string {
  const digits = phone.replace(/\D/g, '');
  return `+60${digits}`;
}

const styles = StyleSheet.create({
  container: {
    marginBottom: 16,
  },
  inputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    borderWidth: 1,
    borderRadius: 8,
    height: 56,
    overflow: 'hidden',
  },
  prefixContainer: {
    paddingHorizontal: 16,
    height: '100%',
    justifyContent: 'center',
    borderRightWidth: 1,
    borderRightColor: '#E5E7EB',
  },
  prefix: {
    fontSize: 16,
    fontWeight: '500',
  },
  input: {
    flex: 1,
    height: '100%',
    paddingHorizontal: 16,
    fontSize: 16,
  },
  helperText: {
    marginTop: 4,
    marginLeft: 4,
  },
});
