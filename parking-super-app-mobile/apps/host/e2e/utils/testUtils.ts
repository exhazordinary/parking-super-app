import { device, element, by, expect, waitFor } from 'detox';

/**
 * Test utilities for Detox E2E tests
 */

/**
 * Wait for an element to be visible with custom timeout
 */
export async function waitForElement(
  testID: string,
  timeout: number = 5000
): Promise<Detox.IndexableNativeElement> {
  const el = element(by.id(testID));
  await waitFor(el).toBeVisible().withTimeout(timeout);
  return el;
}

/**
 * Type text into an input field
 */
export async function typeInInput(testID: string, text: string): Promise<void> {
  const input = await waitForElement(testID);
  await input.clearText();
  await input.typeText(text);
}

/**
 * Tap an element by test ID
 */
export async function tapElement(testID: string): Promise<void> {
  const el = await waitForElement(testID);
  await el.tap();
}

/**
 * Check if an element contains specific text
 */
export async function expectText(testID: string, text: string): Promise<void> {
  const el = await waitForElement(testID);
  await expect(el).toHaveText(text);
}

/**
 * Check if an element is visible
 */
export async function expectVisible(testID: string): Promise<void> {
  const el = element(by.id(testID));
  await expect(el).toBeVisible();
}

/**
 * Check if an element is not visible
 */
export async function expectNotVisible(testID: string): Promise<void> {
  const el = element(by.id(testID));
  await expect(el).not.toBeVisible();
}

/**
 * Scroll to an element
 */
export async function scrollToElement(
  testID: string,
  scrollViewID: string,
  direction: 'up' | 'down' | 'left' | 'right' = 'down',
  pixels: number = 200
): Promise<void> {
  await waitFor(element(by.id(testID)))
    .toBeVisible()
    .whileElement(by.id(scrollViewID))
    .scroll(pixels, direction);
}

/**
 * Swipe on an element
 */
export async function swipeElement(
  testID: string,
  direction: 'up' | 'down' | 'left' | 'right'
): Promise<void> {
  const el = await waitForElement(testID);
  await el.swipe(direction);
}

/**
 * Reload React Native
 */
export async function reloadReactNative(): Promise<void> {
  await device.reloadReactNative();
}

/**
 * Take a screenshot
 */
export async function takeScreenshot(name: string): Promise<void> {
  await device.takeScreenshot(name);
}

/**
 * Test IDs for commonly used elements
 */
export const TestIDs = {
  // Auth screens
  loginScreen: 'login-screen',
  phoneInput: 'phone-input',
  loginButton: 'login-button',
  registerLink: 'register-link',

  registerScreen: 'register-screen',
  nameInput: 'name-input',
  emailInput: 'email-input',
  registerButton: 'register-button',

  otpScreen: 'otp-screen',
  otpInput: 'otp-input',
  verifyButton: 'verify-button',
  resendButton: 'resend-button',

  // Main screens
  homeScreen: 'home-screen',
  walletTab: 'wallet-tab',
  parkingTab: 'parking-tab',
  profileTab: 'profile-tab',

  // Wallet screens
  walletScreen: 'wallet-screen',
  balanceText: 'balance-text',
  topUpButton: 'top-up-button',
  historyButton: 'history-button',

  // Parking screens
  findParkingButton: 'find-parking-button',
  parkingMap: 'parking-map',
  parkingSpotCard: 'parking-spot-card',
  startSessionButton: 'start-session-button',

  // Vehicle screens
  vehicleList: 'vehicle-list',
  addVehicleButton: 'add-vehicle-button',
  plateNumberInput: 'plate-number-input',

  // Common
  loadingIndicator: 'loading-indicator',
  errorMessage: 'error-message',
  backButton: 'back-button',
  confirmButton: 'confirm-button',
  cancelButton: 'cancel-button',
};
