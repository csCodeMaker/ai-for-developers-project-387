import type { Page, Locator } from '@playwright/test';

export class ConfirmationPage {
  constructor(private page: Page) {}

  getBookingId(): Locator {
    return this.page.locator('[data-testid="booking-id"]');
  }

  getSuccessMessage(): Locator {
    return this.page.getByText('Бронь подтверждена!');
  }

  getBackToHomeButton(): Locator {
    return this.page.getByText('На главную');
  }

  async clickBackToHome() {
    await this.getBackToHomeButton().click();
  }
}
