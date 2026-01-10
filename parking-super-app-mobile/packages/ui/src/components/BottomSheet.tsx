import React, { ReactNode, useCallback, useEffect, useRef } from 'react';
import {
  StyleSheet,
  View,
  Dimensions,
  Pressable,
  Animated,
  PanResponder,
  ViewStyle,
} from 'react-native';
import { Portal, Text, useTheme, IconButton } from 'react-native-paper';
import type { MD3Theme } from 'react-native-paper';

const SCREEN_HEIGHT = Dimensions.get('window').height;

export interface BottomSheetProps {
  /** Whether bottom sheet is visible */
  visible: boolean;
  /** Close handler */
  onClose: () => void;
  /** Sheet content */
  children: ReactNode;
  /** Sheet title */
  title?: string;
  /** Sheet height as percentage of screen (0-1) */
  height?: number;
  /** Whether sheet can be dismissed by dragging */
  dragEnabled?: boolean;
  /** Whether to show close button */
  showCloseButton?: boolean;
  /** Additional content styles */
  contentStyle?: ViewStyle;
}

export function BottomSheet({
  visible,
  onClose,
  children,
  title,
  height = 0.5,
  dragEnabled = true,
  showCloseButton = true,
  contentStyle,
}: BottomSheetProps): React.JSX.Element | null {
  const theme = useTheme<MD3Theme>();
  const translateY = useRef(new Animated.Value(SCREEN_HEIGHT)).current;
  const opacity = useRef(new Animated.Value(0)).current;

  const sheetHeight = SCREEN_HEIGHT * height;

  const openSheet = useCallback(() => {
    Animated.parallel([
      Animated.spring(translateY, {
        toValue: 0,
        useNativeDriver: true,
        tension: 65,
        friction: 11,
      }),
      Animated.timing(opacity, {
        toValue: 0.5,
        duration: 200,
        useNativeDriver: true,
      }),
    ]).start();
  }, [translateY, opacity]);

  const closeSheet = useCallback(() => {
    Animated.parallel([
      Animated.timing(translateY, {
        toValue: SCREEN_HEIGHT,
        duration: 200,
        useNativeDriver: true,
      }),
      Animated.timing(opacity, {
        toValue: 0,
        duration: 200,
        useNativeDriver: true,
      }),
    ]).start(() => onClose());
  }, [translateY, opacity, onClose]);

  useEffect(() => {
    if (visible) {
      openSheet();
    }
  }, [visible, openSheet]);

  const panResponder = useRef(
    PanResponder.create({
      onStartShouldSetPanResponder: () => dragEnabled,
      onMoveShouldSetPanResponder: (_, gestureState) =>
        dragEnabled && gestureState.dy > 0,
      onPanResponderMove: (_, gestureState) => {
        if (gestureState.dy > 0) {
          translateY.setValue(gestureState.dy);
        }
      },
      onPanResponderRelease: (_, gestureState) => {
        if (gestureState.dy > sheetHeight * 0.3 || gestureState.vy > 0.5) {
          closeSheet();
        } else {
          Animated.spring(translateY, {
            toValue: 0,
            useNativeDriver: true,
          }).start();
        }
      },
    })
  ).current;

  if (!visible) {
    return null;
  }

  return (
    <Portal>
      <View style={styles.container}>
        <Animated.View
          style={[styles.backdrop, { opacity }]}
        >
          <Pressable style={StyleSheet.absoluteFill} onPress={closeSheet} />
        </Animated.View>

        <Animated.View
          style={[
            styles.sheet,
            {
              height: sheetHeight,
              backgroundColor: theme.colors.surface,
              transform: [{ translateY }],
            },
          ]}
          {...panResponder.panHandlers}
        >
          <View style={styles.handle} />

          {(title || showCloseButton) && (
            <View style={styles.header}>
              <Text variant="titleLarge" style={styles.title}>
                {title}
              </Text>
              {showCloseButton && (
                <IconButton
                  icon="close"
                  size={24}
                  onPress={closeSheet}
                  style={styles.closeButton}
                />
              )}
            </View>
          )}

          <View style={[styles.content, contentStyle]}>
            {children}
          </View>
        </Animated.View>
      </View>
    </Portal>
  );
}

const styles = StyleSheet.create({
  container: {
    ...StyleSheet.absoluteFillObject,
    justifyContent: 'flex-end',
  },
  backdrop: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: '#000000',
  },
  sheet: {
    borderTopLeftRadius: 24,
    borderTopRightRadius: 24,
    paddingHorizontal: 16,
    paddingBottom: 32,
  },
  handle: {
    width: 40,
    height: 4,
    backgroundColor: '#D1D5DB',
    borderRadius: 2,
    alignSelf: 'center',
    marginTop: 12,
    marginBottom: 8,
  },
  header: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingVertical: 8,
  },
  title: {
    fontWeight: '600',
    flex: 1,
  },
  closeButton: {
    margin: 0,
  },
  content: {
    flex: 1,
  },
});
