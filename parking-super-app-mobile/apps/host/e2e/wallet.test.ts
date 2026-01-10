import { device, element, by, expect, waitFor } from 'detox';
import { TestIDs, typeInInput, tapElement, waitForElement } from './utils/testUtils';

describe('Wallet Flow', () => {
  beforeAll(async () => {
    await device.launchApp({ newInstance: true });
    // Login first (assuming test user exists)
    await typeInInput(TestIDs.phoneInput, '60123456789');
    await tapElement(TestIDs.loginButton);
    await waitForElement(TestIDs.otpScreen);
    await typeInInput(TestIDs.otpInput, '123456');
    await tapElement(TestIDs.verifyButton);
    await waitForElement(TestIDs.homeScreen);
  });

  beforeEach(async () => {
    await device.reloadReactNative();
    // Wait for home screen to load
    await waitForElement(TestIDs.homeScreen, 10000);
  });

  describe('Wallet Home', () => {
    beforeEach(async () => {
      await tapElement(TestIDs.walletTab);
      await waitForElement(TestIDs.walletScreen);
    });

    it('should display wallet balance', async () => {
      await expect(element(by.id(TestIDs.balanceText))).toBeVisible();
    });

    it('should display top-up button', async () => {
      await expect(element(by.id(TestIDs.topUpButton))).toBeVisible();
    });

    it('should display transaction history button', async () => {
      await expect(element(by.id(TestIDs.historyButton))).toBeVisible();
    });
  });

  describe('Top Up Flow', () => {
    beforeEach(async () => {
      await tapElement(TestIDs.walletTab);
      await waitForElement(TestIDs.walletScreen);
      await tapElement(TestIDs.topUpButton);
    });

    it('should display predefined amounts', async () => {
      await expect(element(by.text('RM 10'))).toBeVisible();
      await expect(element(by.text('RM 20'))).toBeVisible();
      await expect(element(by.text('RM 50'))).toBeVisible();
      await expect(element(by.text('RM 100'))).toBeVisible();
    });

    it('should allow custom amount entry', async () => {
      await tapElement('custom-amount-input');
      await element(by.id('custom-amount-input')).typeText('75');
      await expect(element(by.text('RM 75'))).toBeVisible();
    });

    it('should navigate to payment on confirm', async () => {
      await tapElement('amount-50');
      await tapElement(TestIDs.confirmButton);

      // Should show payment options
      await waitFor(element(by.text('Select Payment Method')))
        .toBeVisible()
        .withTimeout(3000);
    });
  });

  describe('Transaction History', () => {
    beforeEach(async () => {
      await tapElement(TestIDs.walletTab);
      await waitForElement(TestIDs.walletScreen);
      await tapElement(TestIDs.historyButton);
    });

    it('should display transaction list', async () => {
      await waitFor(element(by.id('transaction-list')))
        .toBeVisible()
        .withTimeout(5000);
    });

    it('should display transaction details on tap', async () => {
      // Tap first transaction
      await element(by.id('transaction-item')).atIndex(0).tap();

      // Should show transaction details modal
      await waitFor(element(by.id('transaction-details-modal')))
        .toBeVisible()
        .withTimeout(3000);
    });

    it('should allow filtering by date', async () => {
      await tapElement('filter-button');
      await tapElement('filter-this-month');

      // Should update list
      await waitFor(element(by.id('transaction-list')))
        .toBeVisible()
        .withTimeout(3000);
    });
  });
});
