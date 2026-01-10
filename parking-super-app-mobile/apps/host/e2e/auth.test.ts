import { device, element, by, expect, waitFor } from 'detox';
import { TestIDs, typeInInput, tapElement, waitForElement } from './utils/testUtils';

describe('Authentication Flow', () => {
  beforeAll(async () => {
    await device.launchApp({ newInstance: true });
  });

  beforeEach(async () => {
    await device.reloadReactNative();
  });

  describe('Login Screen', () => {
    it('should display login screen on app launch', async () => {
      await expect(element(by.id(TestIDs.loginScreen))).toBeVisible();
    });

    it('should display phone input field', async () => {
      await expect(element(by.id(TestIDs.phoneInput))).toBeVisible();
    });

    it('should display login button', async () => {
      await expect(element(by.id(TestIDs.loginButton))).toBeVisible();
    });

    it('should display link to register', async () => {
      await expect(element(by.id(TestIDs.registerLink))).toBeVisible();
    });

    it('should validate phone number format', async () => {
      await typeInInput(TestIDs.phoneInput, '123');
      await tapElement(TestIDs.loginButton);

      // Should show validation error
      await waitFor(element(by.id(TestIDs.errorMessage)))
        .toBeVisible()
        .withTimeout(3000);
    });

    it('should accept valid Malaysian phone number', async () => {
      await typeInInput(TestIDs.phoneInput, '60123456789');
      await tapElement(TestIDs.loginButton);

      // Should navigate to OTP screen
      await waitFor(element(by.id(TestIDs.otpScreen)))
        .toBeVisible()
        .withTimeout(5000);
    });
  });

  describe('Registration Screen', () => {
    beforeEach(async () => {
      await tapElement(TestIDs.registerLink);
      await waitForElement(TestIDs.registerScreen);
    });

    it('should display registration form', async () => {
      await expect(element(by.id(TestIDs.nameInput))).toBeVisible();
      await expect(element(by.id(TestIDs.phoneInput))).toBeVisible();
      await expect(element(by.id(TestIDs.emailInput))).toBeVisible();
    });

    it('should validate required fields', async () => {
      await tapElement(TestIDs.registerButton);

      // Should show validation error
      await waitFor(element(by.id(TestIDs.errorMessage)))
        .toBeVisible()
        .withTimeout(3000);
    });

    it('should register with valid data', async () => {
      await typeInInput(TestIDs.nameInput, 'John Doe');
      await typeInInput(TestIDs.phoneInput, '60198765432');
      await typeInInput(TestIDs.emailInput, 'john@example.com');
      await tapElement(TestIDs.registerButton);

      // Should navigate to OTP screen
      await waitFor(element(by.id(TestIDs.otpScreen)))
        .toBeVisible()
        .withTimeout(5000);
    });
  });

  describe('OTP Verification', () => {
    beforeEach(async () => {
      // Navigate to OTP screen via login
      await typeInInput(TestIDs.phoneInput, '60123456789');
      await tapElement(TestIDs.loginButton);
      await waitForElement(TestIDs.otpScreen);
    });

    it('should display OTP input field', async () => {
      await expect(element(by.id(TestIDs.otpInput))).toBeVisible();
    });

    it('should display verify button', async () => {
      await expect(element(by.id(TestIDs.verifyButton))).toBeVisible();
    });

    it('should display resend button', async () => {
      await expect(element(by.id(TestIDs.resendButton))).toBeVisible();
    });

    it('should verify with valid OTP', async () => {
      // In test environment, use mock OTP '123456'
      await typeInInput(TestIDs.otpInput, '123456');
      await tapElement(TestIDs.verifyButton);

      // Should navigate to home screen
      await waitFor(element(by.id(TestIDs.homeScreen)))
        .toBeVisible()
        .withTimeout(5000);
    });

    it('should show error for invalid OTP', async () => {
      await typeInInput(TestIDs.otpInput, '000000');
      await tapElement(TestIDs.verifyButton);

      // Should show error message
      await waitFor(element(by.id(TestIDs.errorMessage)))
        .toBeVisible()
        .withTimeout(3000);
    });
  });
});
