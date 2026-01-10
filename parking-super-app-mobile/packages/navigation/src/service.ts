import {
  createNavigationContainerRef,
  CommonActions,
  StackActions,
} from '@react-navigation/native';
import type { RootStackParamList } from './types';

/**
 * Navigation container ref for imperative navigation
 */
export const navigationRef = createNavigationContainerRef<RootStackParamList>();

/**
 * Check if navigation is ready
 */
export function isNavigationReady(): boolean {
  return navigationRef.isReady();
}

/**
 * Navigate to a screen
 */
export function navigate<T extends keyof RootStackParamList>(
  name: T,
  params?: RootStackParamList[T]
): void {
  if (navigationRef.isReady()) {
    // @ts-expect-error - navigation types are complex
    navigationRef.navigate(name, params);
  }
}

/**
 * Go back
 */
export function goBack(): void {
  if (navigationRef.isReady() && navigationRef.canGoBack()) {
    navigationRef.goBack();
  }
}

/**
 * Reset navigation state
 */
export function resetRoot(routeName: keyof RootStackParamList): void {
  if (navigationRef.isReady()) {
    navigationRef.dispatch(
      CommonActions.reset({
        index: 0,
        routes: [{ name: routeName }],
      })
    );
  }
}

/**
 * Push a new screen onto the stack
 */
export function push<T extends keyof RootStackParamList>(
  name: T,
  params?: RootStackParamList[T]
): void {
  if (navigationRef.isReady()) {
    // @ts-expect-error - navigation types are complex
    navigationRef.dispatch(StackActions.push(name, params));
  }
}

/**
 * Replace current screen
 */
export function replace<T extends keyof RootStackParamList>(
  name: T,
  params?: RootStackParamList[T]
): void {
  if (navigationRef.isReady()) {
    // @ts-expect-error - navigation types are complex
    navigationRef.dispatch(StackActions.replace(name, params));
  }
}

/**
 * Pop to top of stack
 */
export function popToTop(): void {
  if (navigationRef.isReady()) {
    navigationRef.dispatch(StackActions.popToTop());
  }
}

/**
 * Get current route name
 */
export function getCurrentRouteName(): string | undefined {
  if (navigationRef.isReady()) {
    return navigationRef.getCurrentRoute()?.name;
  }
  return undefined;
}

/**
 * Get current route params
 */
export function getCurrentRouteParams(): Record<string, unknown> | undefined {
  if (navigationRef.isReady()) {
    return navigationRef.getCurrentRoute()?.params as Record<string, unknown> | undefined;
  }
  return undefined;
}

/**
 * Navigation service object for easy imports
 */
export const NavigationService = {
  navigate,
  goBack,
  resetRoot,
  push,
  replace,
  popToTop,
  getCurrentRouteName,
  getCurrentRouteParams,
  isReady: isNavigationReady,
  ref: navigationRef,
};
