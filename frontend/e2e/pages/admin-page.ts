import type { Page, Locator } from '@playwright/test';

export class AdminPage {
  constructor(private page: Page) {}

  goto() {
    return this.page.goto('/admin');
  }

  getBookingsTable(): Locator {
    return this.page.locator('[data-testid="bookings-table"]');
  }

  getBookingRows(): Locator {
    return this.page.locator('[data-testid="booking-row"]');
  }

  getEventTypeList(): Locator {
    return this.page.locator('[data-testid="event-type-list"]');
  }

  getEventTypeRows(): Locator {
    return this.page.locator('[data-testid="event-type-row"]');
  }

  getCreateEventTypeButton(): Locator {
    return this.page.locator('[data-testid="create-event-type"]');
  }

  getEventTypeTitleInput(): Locator {
    return this.page.locator('[data-testid="event-type-title"]');
  }

  getEventTypeDescriptionInput(): Locator {
    return this.page.locator('[data-testid="event-type-description"]');
  }

  getEventTypeDurationInput(): Locator {
    return this.page.locator('[data-testid="event-type-duration"]');
  }

  getSaveEventTypeButton(): Locator {
    return this.page.locator('[data-testid="save-event-type"]');
  }

  getEditEventTypeButton(): Locator {
    return this.page.locator('[data-testid="edit-event-type"]');
  }

  getDeleteEventTypeButton(): Locator {
    return this.page.locator('[data-testid="delete-event-type"]');
  }

  getConfirmDeleteButton(): Locator {
    return this.page.locator('[data-testid="confirm-delete"]');
  }

  getOwnerForm(): Locator {
    return this.page.locator('[data-testid="owner-form"]');
  }

  getOwnerNameInput(): Locator {
    return this.page.locator('[data-testid="owner-name"]');
  }

  getOwnerEmailInput(): Locator {
    return this.page.locator('[data-testid="owner-email"]');
  }

  getOwnerDescriptionInput(): Locator {
    return this.page.locator('[data-testid="owner-description"]');
  }

  getOwnerTimeZoneInput(): Locator {
    return this.page.locator('[data-testid="owner-timezone"]');
  }

  getSaveOwnerButton(): Locator {
    return this.page.locator('[data-testid="save-owner"]');
  }

  getEmptyBookingsMessage(): Locator {
    return this.page.getByText('Нет предстоящих бронирований');
  }

  getDisabledBadge(): Locator {
    return this.page.getByText('Отключён');
  }

  async createEventType(title: string, description: string, duration: number) {
    await this.getCreateEventTypeButton().click();
    await this.getEventTypeTitleInput().fill(title);
    await this.getEventTypeDescriptionInput().fill(description);
    await this.getEventTypeDurationInput().fill(String(duration));
    await this.getSaveEventTypeButton().click();
  }
}
