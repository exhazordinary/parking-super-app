import { AppRegistry } from 'react-native';
import { ScriptManager, Script } from '@callstack/repack/client';
import App from './src/App';
import { name as appName } from './app.json';

// Configure script manager for module federation
ScriptManager.shared.addResolver(async (scriptId, caller) => {
  // Map module names to their remote URLs
  const resolveURL = (id) => {
    // In development, modules are served from the dev server
    if (__DEV__) {
      return `http://localhost:9000/${id}.container.bundle`;
    }

    // In production, modules could be served from a CDN or bundled
    // You can customize this based on your deployment strategy
    return `https://cdn.yourapp.com/modules/${id}.container.bundle`;
  };

  // Module federation remotes
  const remoteModules = {
    auth_module: resolveURL('auth_module'),
    wallet_module: resolveURL('wallet_module'),
    parking_module: resolveURL('parking_module'),
    vehicle_module: resolveURL('vehicle_module'),
    provider_module: resolveURL('provider_module'),
    notification_module: resolveURL('notification_module'),
  };

  if (scriptId in remoteModules) {
    return {
      url: Script.getRemoteURL(remoteModules[scriptId]),
      cache: !__DEV__,
    };
  }

  return undefined;
});

AppRegistry.registerComponent(appName, () => App);
