import React, { ReactNode } from 'react';
import { StyleSheet, ViewStyle } from 'react-native';
import { Card as PaperCard, useTheme } from 'react-native-paper';
import type { MD3Theme } from 'react-native-paper';

export interface CardProps {
  /** Card content */
  children: ReactNode;
  /** Card title */
  title?: string;
  /** Card subtitle */
  subtitle?: string;
  /** Whether card is pressable */
  onPress?: () => void;
  /** Card elevation */
  elevation?: 0 | 1 | 2 | 3 | 4 | 5;
  /** Card mode */
  mode?: 'elevated' | 'outlined' | 'contained';
  /** Additional styles */
  style?: ViewStyle;
}

export function Card({
  children,
  title,
  subtitle,
  onPress,
  elevation = 1,
  mode = 'elevated',
  style,
}: CardProps): React.JSX.Element {
  const theme = useTheme<MD3Theme>();

  return (
    <PaperCard
      mode={mode}
      elevation={elevation}
      onPress={onPress}
      style={[
        styles.card,
        { backgroundColor: theme.colors.surface },
        style,
      ]}
    >
      {(title || subtitle) && (
        <PaperCard.Title
          title={title}
          subtitle={subtitle}
          titleStyle={styles.title}
          subtitleStyle={styles.subtitle}
        />
      )}
      <PaperCard.Content style={styles.content}>
        {children}
      </PaperCard.Content>
    </PaperCard>
  );
}

export const CardActions = PaperCard.Actions;

const styles = StyleSheet.create({
  card: {
    borderRadius: 12,
    marginVertical: 8,
  },
  title: {
    fontSize: 18,
    fontWeight: '600',
  },
  subtitle: {
    fontSize: 14,
  },
  content: {
    paddingTop: 0,
  },
});
