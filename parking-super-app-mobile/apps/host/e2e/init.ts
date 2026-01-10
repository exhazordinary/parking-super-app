import { device } from 'detox';

/**
 * Detox test setup file
 * This runs before all tests
 */

beforeAll(async () => {
  // Install the app if not already installed
  await device.launchApp({ newInstance: true });
});

afterAll(async () => {
  // Cleanup after tests
  await device.terminateApp();
});

/**
 * Test configuration
 */
export const TestConfig = {
  // API endpoints for test environment
  apiBaseUrl: 'http://localhost:3000',

  // Test user credentials
  testUser: {
    phone: '60123456789',
    otp: '123456',
  },

  // Timeouts
  defaultTimeout: 5000,
  longTimeout: 15000,

  // Test data
  testVehicle: {
    plateNumber: 'ABC 1234',
    type: 'car',
  },
};
