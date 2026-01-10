import { device, element, by, expect, waitFor } from 'detox';
import { TestIDs, typeInInput, tapElement, waitForElement } from './utils/testUtils';

describe('Vehicle Management Flow', () => {
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
    // Navigate to profile/vehicles
    await tapElement(TestIDs.profileTab);
    await tapElement('manage-vehicles-button');
  });

  describe('Vehicle List', () => {
    it('should display vehicle list', async () => {
      await expect(element(by.id(TestIDs.vehicleList))).toBeVisible();
    });

    it('should display add vehicle button', async () => {
      await expect(element(by.id(TestIDs.addVehicleButton))).toBeVisible();
    });

    it('should show empty state when no vehicles', async () => {
      // This test assumes no vehicles initially
      await waitFor(element(by.id('empty-vehicle-state')))
        .toBeVisible()
        .withTimeout(3000);
    });
  });

  describe('Add Vehicle', () => {
    beforeEach(async () => {
      await tapElement(TestIDs.addVehicleButton);
    });

    it('should display plate number input', async () => {
      await expect(element(by.id(TestIDs.plateNumberInput))).toBeVisible();
    });

    it('should display vehicle type selection', async () => {
      await expect(element(by.id('vehicle-type-selection'))).toBeVisible();
    });

    it('should validate Malaysian plate format', async () => {
      // Invalid format
      await typeInInput(TestIDs.plateNumberInput, 'INVALID');
      await tapElement(TestIDs.confirmButton);

      await waitFor(element(by.id(TestIDs.errorMessage)))
        .toBeVisible()
        .withTimeout(3000);
    });

    it('should accept valid Malaysian plate format', async () => {
      // Valid format: ABC 1234
      await typeInInput(TestIDs.plateNumberInput, 'WXY 1234');
      await tapElement('vehicle-type-car');
      await tapElement(TestIDs.confirmButton);

      // Should navigate back to list
      await waitFor(element(by.id(TestIDs.vehicleList)))
        .toBeVisible()
        .withTimeout(5000);
    });

    it('should auto-format plate number', async () => {
      await typeInInput(TestIDs.plateNumberInput, 'wxy1234');

      // Should be formatted to uppercase
      await expect(element(by.id(TestIDs.plateNumberInput))).toHaveText(
        'WXY 1234'
      );
    });
  });

  describe('Vehicle Details', () => {
    beforeEach(async () => {
      // Add a vehicle first
      await tapElement(TestIDs.addVehicleButton);
      await typeInInput(TestIDs.plateNumberInput, 'ABC 9999');
      await tapElement('vehicle-type-car');
      await tapElement(TestIDs.confirmButton);
      await waitForElement(TestIDs.vehicleList);

      // Tap on the vehicle
      await element(by.id('vehicle-item')).atIndex(0).tap();
    });

    it('should display vehicle details', async () => {
      await expect(element(by.text('ABC 9999'))).toBeVisible();
    });

    it('should display parking history', async () => {
      await expect(element(by.id('vehicle-parking-history'))).toBeVisible();
    });

    it('should allow setting as default', async () => {
      await tapElement('set-default-button');

      await waitFor(element(by.id('default-badge')))
        .toBeVisible()
        .withTimeout(3000);
    });

    it('should allow editing nickname', async () => {
      await tapElement('edit-nickname-button');
      await typeInInput('nickname-input', 'My Car');
      await tapElement(TestIDs.confirmButton);

      await expect(element(by.text('My Car'))).toBeVisible();
    });
  });

  describe('Delete Vehicle', () => {
    beforeEach(async () => {
      // Navigate to vehicle details
      await element(by.id('vehicle-item')).atIndex(0).tap();
    });

    it('should show confirmation dialog', async () => {
      await tapElement('delete-vehicle-button');

      await expect(element(by.text('Delete Vehicle?'))).toBeVisible();
      await expect(element(by.id(TestIDs.confirmButton))).toBeVisible();
      await expect(element(by.id(TestIDs.cancelButton))).toBeVisible();
    });

    it('should cancel deletion', async () => {
      await tapElement('delete-vehicle-button');
      await tapElement(TestIDs.cancelButton);

      // Should still be on details screen
      await expect(element(by.id('vehicle-details-screen'))).toBeVisible();
    });

    it('should confirm deletion', async () => {
      await tapElement('delete-vehicle-button');
      await tapElement(TestIDs.confirmButton);

      // Should navigate back to list
      await waitFor(element(by.id(TestIDs.vehicleList)))
        .toBeVisible()
        .withTimeout(5000);
    });
  });
});
