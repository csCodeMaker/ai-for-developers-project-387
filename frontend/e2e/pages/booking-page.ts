import type { Page, Locator } from '@playwright/test';

export class BookingPage {
  constructor(private page: Page) {}

  async goto(eventTypeId: string) {
    await this.page.goto(`/book/${eventTypeId}`);
  }

  getDatePicker(): Locator {
    return this.page.locator('[data-testid="date-picker"]');
  }

  getSlots(): Locator {
    return this.page.locator('[data-testid="slot"]');
  }

  getBusySlots(): Locator {
    return this.page.locator('[data-testid="slot"][data-busy="true"]');
  }

  getFreeSlots(): Locator {
    return this.page.locator('[data-testid="slot"][data-busy="false"]');
  }

  clickSlot(index: number) {
    return this.getFreeSlots().nth(index).click();
  }

  getNameInput(): Locator {
    return this.page.locator('[data-testid="guest-name"]');
  }

  getEmailInput(): Locator {
    return this.page.locator('[data-testid="guest-email"]');
  }

  getSubmitButton(): Locator {
    return this.page.locator('[data-testid="book-submit"]');
  }

  async fillName(name: string) {
    await this.getNameInput().fill(name);
  }

  async fillEmail(email: string) {
    await this.getEmailInput().fill(email);
  }

  async submit() {
    await this.getSubmitButton().click();
  }

  getErrorMessage(): Locator {
    return this.page.locator('[data-testid="booking-error"]');
  }

  getNoSlotsMessage(): Locator {
    return this.page.getByText('Нет свободных слотов на эту дату');
  }
}
