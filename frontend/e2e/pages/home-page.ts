import type { Page, Locator } from '@playwright/test';

export class HomePage {
  constructor(private page: Page) {}

  goto() {
    return this.page.goto('/');
  }

  getEventTypeCards(): Locator {
    return this.page.locator('[data-testid="event-type-card"]');
  }

  clickFirstEventType() {
    return this.getEventTypeCards().first().click();
  }

  isEmpty() {
    return this.page.getByText('Нет доступных типов событий').isVisible();
  }
}
