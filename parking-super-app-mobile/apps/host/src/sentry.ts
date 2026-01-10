import * as Sentry from '@sentry/react-native';

/**
 * Sentry configuration for error tracking
 *
 * IMPORTANT: Replace SENTRY_DSN with your actual DSN
 * Get your DSN from: https://sentry.io/settings/projects/{project}/keys/
 */

const SENTRY_DSN = ''; // TODO: Add your Sentry DSN here

export function initializeSentry(): void {
  if (!SENTRY_DSN) {
    console.warn('Sentry DSN not configured. Error tracking disabled.');
    return;
  }

  Sentry.init({
    dsn: SENTRY_DSN,

    // Enable automatic performance monitoring
    tracesSampleRate: __DEV__ ? 1.0 : 0.2,

    // Enable session tracking
    enableAutoSessionTracking: true,
    sessionTrackingIntervalMillis: 30000,

    // Configure release info
    release: 'parking-app@1.0.0',
    dist: '1',

    // Environment
    environment: __DEV__ ? 'development' : 'production',

    // Disable in development
    enabled: !__DEV__,

    // Configure integrations
    integrations: [
      new Sentry.ReactNativeTracing({
        tracingOrigins: ['localhost', /^\/api/],
        routingInstrumentation: Sentry.reactNativeNavigationIntegration,
      }),
    ],

    // Filter out sensitive data
    beforeSend(event) {
      // Remove sensitive headers
      if (event.request?.headers) {
        delete event.request.headers['Authorization'];
        delete event.request.headers['X-User-ID'];
      }
      return event;
    },

    // Configure breadcrumbs
    beforeBreadcrumb(breadcrumb) {
      // Filter out noisy breadcrumbs
      if (breadcrumb.category === 'xhr' && breadcrumb.data?.url?.includes('/health')) {
        return null;
      }
      return breadcrumb;
    },
  });
}

/**
 * Set user context for error tracking
 */
export function setSentryUser(user: { id: string; phone?: string } | null): void {
  if (user) {
    Sentry.setUser({ id: user.id, username: user.phone });
  } else {
    Sentry.setUser(null);
  }
}

/**
 * Add custom breadcrumb
 */
export function addBreadcrumb(
  category: string,
  message: string,
  data?: Record<string, unknown>
): void {
  Sentry.addBreadcrumb({
    category,
    message,
    data,
    level: 'info',
  });
}

/**
 * Capture exception with additional context
 */
export function captureException(
  error: Error,
  context?: Record<string, unknown>
): void {
  Sentry.withScope((scope) => {
    if (context) {
      scope.setExtras(context);
    }
    Sentry.captureException(error);
  });
}

/**
 * Capture message
 */
export function captureMessage(
  message: string,
  level: Sentry.SeverityLevel = 'info'
): void {
  Sentry.captureMessage(message, level);
}

/**
 * Start a transaction for performance monitoring
 */
export function startTransaction(
  name: string,
  op: string
): Sentry.Transaction {
  return Sentry.startTransaction({ name, op });
}

export { Sentry };
