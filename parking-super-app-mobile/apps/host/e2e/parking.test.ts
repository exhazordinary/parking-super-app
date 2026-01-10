import { device, element, by, expect, waitFor } from 'detox';
import { TestIDs, typeInInput, tapElement, waitForElement } from './utils/testUtils';

describe('Parking Flow', () => {
  beforeAll(async () => {
    await device.launchApp({ newInstance: true });
    // Login first
    await typeInInput(TestIDs.phoneInput, '60123456789');
    await tapElement(TestIDs.loginButton);
    await waitForElement(TestIDs.otpScreen);
    await typeInInput(TestIDs.otpInput, '123456');
    await tapElement(TestIDs.verifyButton);
    await waitForElement(TestIDs.homeScreen);
  });

  beforeEach(async () => {
    await device.reloadReactNative();
    await waitForElement(TestIDs.homeScreen, 10000);
  });

  describe('Parking Home', () => {
    beforeEach(async () => {
      await tapElement(TestIDs.parkingTab);
    });

    it('should display find parking button', async () => {
      await expect(element(by.id(TestIDs.findParkingButton))).toBeVisible();
    });

    it('should display recent parking sessions', async () => {
      await waitFor(element(by.id('recent-sessions-list')))
        .toBeVisible()
        .withTimeout(5000);
    });
  });

  describe('Find Parking', () => {
    beforeEach(async () => {
      await tapElement(TestIDs.parkingTab);
      await tapElement(TestIDs.findParkingButton);
    });

    it('should display parking map', async () => {
      await waitFor(element(by.id(TestIDs.parkingMap)))
        .toBeVisible()
        .withTimeout(5000);
    });

    it('should display search input', async () => {
      await expect(element(by.id('search-location-input'))).toBeVisible();
    });

    it('should search for parking spots', async () => {
      await typeInInput('search-location-input', 'KLCC');
      await tapElement('search-button');

      // Should show parking spots
      await waitFor(element(by.id(TestIDs.parkingSpotCard)))
        .toBeVisible()
        .withTimeout(5000);
    });

    it('should display parking spot details on tap', async () => {
      await typeInInput('search-location-input', 'KLCC');
      await tapElement('search-button');
      await waitForElement(TestIDs.parkingSpotCard);
      await tapElement(TestIDs.parkingSpotCard);

      // Should show spot details
      await expect(element(by.text('Available Spots'))).toBeVisible();
      await expect(element(by.text('Hourly Rate'))).toBeVisible();
    });
  });

  describe('Start Parking Session', () => {
    beforeEach(async () => {
      await tapElement(TestIDs.parkingTab);
      await tapElement(TestIDs.findParkingButton);
      // Search and select a parking spot
      await typeInInput('search-location-input', 'KLCC');
      await tapElement('search-button');
      await waitForElement(TestIDs.parkingSpotCard);
      await tapElement(TestIDs.parkingSpotCard);
    });

    it('should display vehicle selection', async () => {
      await tapElement(TestIDs.startSessionButton);
      await expect(element(by.id('vehicle-selection'))).toBeVisible();
    });

    it('should display duration selection', async () => {
      await tapElement(TestIDs.startSessionButton);
      await tapElement('vehicle-item');
      await expect(element(by.id('duration-selection'))).toBeVisible();
    });

    it('should show payment summary', async () => {
      await tapElement(TestIDs.startSessionButton);
      await tapElement('vehicle-item');
      await tapElement('duration-2hours');

      await expect(element(by.id('payment-summary'))).toBeVisible();
      await expect(element(by.text('Total'))).toBeVisible();
    });

    it('should start session successfully', async () => {
      await tapElement(TestIDs.startSessionButton);
      await tapElement('vehicle-item');
      await tapElement('duration-2hours');
      await tapElement(TestIDs.confirmButton);

      // Should show active session
      await waitFor(element(by.id('active-session-screen')))
        .toBeVisible()
        .withTimeout(5000);
    });
  });

  describe('Active Session', () => {
    // Note: This assumes an active session exists
    beforeEach(async () => {
      await tapElement(TestIDs.parkingTab);
      // Tap on active session banner if visible
      try {
        await element(by.id('active-session-banner')).tap();
      } catch {
        // No active session, skip
      }
    });

    it('should display remaining time', async () => {
      await expect(element(by.id('remaining-time'))).toBeVisible();
    });

    it('should display extend option', async () => {
      await expect(element(by.id('extend-button'))).toBeVisible();
    });

    it('should display end session option', async () => {
      await expect(element(by.id('end-session-button'))).toBeVisible();
    });

    it('should extend session', async () => {
      await tapElement('extend-button');
      await tapElement('extend-1hour');
      await tapElement(TestIDs.confirmButton);

      // Should show updated time
      await waitFor(element(by.id('extension-success')))
        .toBeVisible()
        .withTimeout(3000);
    });
  });
});
