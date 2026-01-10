import React, { useEffect, ReactNode } from 'react';
import { useAuthStore } from '@parking/auth';

interface AppProvidersProps {
  children: ReactNode;
}

export function AppProviders({ children }: AppProvidersProps): React.JSX.Element {
  const { initialize, isInitialized } = useAuthStore();

  useEffect(() => {
    // Initialize auth state on app start
    if (!isInitialized) {
      initialize();
    }
  }, [initialize, isInitialized]);

  return <>{children}</>;
}
