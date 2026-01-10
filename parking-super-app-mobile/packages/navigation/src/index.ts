// Types
export * from './types';

// Deep linking
export { linking } from './linking';

// Navigation service
export {
  navigationRef,
  NavigationService,
  navigate,
  goBack,
  resetRoot,
  push,
  replace,
  popToTop,
  getCurrentRouteName,
  getCurrentRouteParams,
  isNavigationReady,
} from './service';
