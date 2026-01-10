import React, { forwardRef } from 'react';
import { StyleSheet, ViewStyle, TextInput as RNTextInput } from 'react-native';
import { TextInput, HelperText, useTheme } from 'react-native-paper';
import type { MD3Theme } from 'react-native-paper';

export interface InputProps {
  /** Input label */
  label: string;
  /** Input value */
  value: string;
  /** Change handler */
  onChangeText: (text: string) => void;
  /** Placeholder text */
  placeholder?: string;
  /** Whether input is disabled */
  disabled?: boolean;
  /** Whether input has error */
  error?: boolean;
  /** Helper/error text */
  helperText?: string;
  /** Input mode */
  mode?: 'flat' | 'outlined';
  /** Whether to hide text (password) */
  secureTextEntry?: boolean;
  /** Keyboard type */
  keyboardType?: 'default' | 'email-address' | 'numeric' | 'phone-pad' | 'decimal-pad';
  /** Auto capitalize */
  autoCapitalize?: 'none' | 'sentences' | 'words' | 'characters';
  /** Left icon */
  left?: React.ReactNode;
  /** Right icon */
  right?: React.ReactNode;
  /** Whether input is multiline */
  multiline?: boolean;
  /** Number of lines for multiline */
  numberOfLines?: number;
  /** Additional styles */
  style?: ViewStyle;
  /** Container styles */
  containerStyle?: ViewStyle;
  /** On submit */
  onSubmitEditing?: () => void;
  /** Return key type */
  returnKeyType?: 'done' | 'go' | 'next' | 'search' | 'send';
  /** Auto focus */
  autoFocus?: boolean;
  /** Max length */
  maxLength?: number;
  /** Editable */
  editable?: boolean;
}

export const Input = forwardRef<RNTextInput, InputProps>(function Input(
  {
    label,
    value,
    onChangeText,
    placeholder,
    disabled = false,
    error = false,
    helperText,
    mode = 'outlined',
    secureTextEntry = false,
    keyboardType = 'default',
    autoCapitalize = 'none',
    left,
    right,
    multiline = false,
    numberOfLines = 1,
    style,
    containerStyle,
    onSubmitEditing,
    returnKeyType,
    autoFocus = false,
    maxLength,
    editable = true,
  },
  ref
): React.JSX.Element {
  const theme = useTheme<MD3Theme>();

  return (
    <>
      <TextInput
        ref={ref}
        label={label}
        value={value}
        onChangeText={onChangeText}
        placeholder={placeholder}
        disabled={disabled}
        error={error}
        mode={mode}
        secureTextEntry={secureTextEntry}
        keyboardType={keyboardType}
        autoCapitalize={autoCapitalize}
        left={left}
        right={right}
        multiline={multiline}
        numberOfLines={numberOfLines}
        style={[styles.input, style]}
        contentStyle={styles.content}
        outlineStyle={styles.outline}
        onSubmitEditing={onSubmitEditing}
        returnKeyType={returnKeyType}
        autoFocus={autoFocus}
        maxLength={maxLength}
        editable={editable}
      />
      {helperText && (
        <HelperText
          type={error ? 'error' : 'info'}
          visible={Boolean(helperText)}
          style={styles.helper}
        >
          {helperText}
        </HelperText>
      )}
    </>
  );
});

// Icon helpers
Input.Icon = TextInput.Icon;
Input.Affix = TextInput.Affix;

const styles = StyleSheet.create({
  input: {
    marginBottom: 4,
  },
  content: {
    paddingVertical: 0,
  },
  outline: {
    borderRadius: 8,
  },
  helper: {
    marginTop: -4,
    marginBottom: 8,
  },
});
