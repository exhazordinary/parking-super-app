/**
 * Type declarations for federated modules
 * These are loaded at runtime via Module Federation
 */

declare module 'auth_module/AuthNavigator' {
  import { FC } from 'react';

  interface AuthNavigatorProps {
    initialScreen?: 'Login' | 'Register' | 'Profile';
  }

  const AuthNavigator: FC<AuthNavigatorProps>;
  export default AuthNavigator;
}

declare module 'wallet_module/WalletNavigator' {
  import { FC } from 'react';

  const WalletNavigator: FC;
  export default WalletNavigator;
}

declare module 'parking_module/ParkingNavigator' {
  import { FC } from 'react';

  const ParkingNavigator: FC;
  export default ParkingNavigator;
}

declare module 'vehicle_module/VehicleNavigator' {
  import { FC } from 'react';

  const VehicleNavigator: FC;
  export default VehicleNavigator;
}

declare module 'provider_module/ProviderNavigator' {
  import { FC } from 'react';

  const ProviderNavigator: FC;
  export default ProviderNavigator;
}

declare module 'notification_module/NotificationNavigator' {
  import { FC } from 'react';

  const NotificationNavigator: FC;
  export default NotificationNavigator;
}
