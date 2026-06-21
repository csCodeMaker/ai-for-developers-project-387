import { test, expect } from '@playwright/test';
import { AdminPage } from '../pages/admin-page';
import { API_BASE, SEED_OWNER } from '../fixtures/test-data';

test.describe.configure({ mode: 'serial' });

test.describe('Администратор — просмотр бронирований (US-6)', () => {
  test('видит бронирование после того, как гость записался', async ({ page }) => {
    const res = await page.request.get(`${API_BASE}/api/admin/event-types`);
    const eventTypes = await res.json();
    const et = eventTypes[0];

    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    const dateStr = tomorrow.toISOString().slice(0, 10);
    const slotStart = `${dateStr}T10:00:00.000Z`;

    await page.request.post(`${API_BASE}/api/bookings`, {
      data: {
        eventTypeId: et.id,
        guestName: 'Гость Тестов',
        guestEmail: 'guest@test.com',
        startTime: slotStart,
      },
    });

    const admin = new AdminPage(page);
    await admin.goto();
    await expect(admin.getBookingsTable()).toBeVisible({ timeout: 10000 });
    await expect(admin.getBookingRows().first()).toBeVisible();

    await expect(page.getByText('Гость Тестов')).toBeVisible();
    await expect(page.getByText('guest@test.com')).toBeVisible();
  });
});

test.describe('Администратор — создание типа события (US-7)', () => {
  test('создаёт новый тип и видит его в списке', async ({ page }) => {
    const admin = new AdminPage(page);
    await admin.goto();

    await admin.createEventType('Тестовый звонок', 'Описание тестового звонка', 45);

    await expect(page.getByText('Тестовый звонок')).toBeVisible();
    await expect(page.getByText('45 мин')).toBeVisible();
    await expect(page.getByRole('cell', { name: 'Активен' }).first()).toBeVisible();
  });
});

test.describe('Администратор — редактирование типа события (US-8)', () => {
  test('изменяет название типа события', async ({ page }) => {
    const admin = new AdminPage(page);
    await admin.goto();
    await admin.getEditEventTypeButton().first().click();

    await admin.getEventTypeTitleInput().fill('Изменённый звонок');
    await admin.getSaveEventTypeButton().click();

    await expect(page.getByText('Изменённый звонок')).toBeVisible();
  });
});

test.describe('Администратор — отключение типа события (US-9)', () => {
  test('отключает тип и проверяет статус', async ({ page }) => {
    const admin = new AdminPage(page);
    await admin.goto();

    const disabledCountBefore = await page.getByRole('cell', { name: 'Отключён' }).count();

    await expect(admin.getDeleteEventTypeButton().first()).toBeVisible({ timeout: 10000 });
    await expect(admin.getEditEventTypeButton().first()).toBeVisible();

    await admin.getDeleteEventTypeButton().first().click();

    await expect(page.getByRole('cell', { name: 'Отключён' })).toHaveCount(disabledCountBefore + 1);
  });

  test('отключённый тип не показывается на главной', async ({ page }) => {
    const res = await page.request.get(`${API_BASE}/api/admin/event-types`);
    const eventTypes: Array<{ id: string; title: string }> = await res.json();

    for (const et of eventTypes) {
      await page.request.delete(`${API_BASE}/api/admin/event-types/${et.id}`);
    }

    await page.goto('/');
    await expect(page.locator('[data-testid="event-type-card"]')).toHaveCount(0);
  });
});

test.describe('Администратор — редактирование профиля владельца (US-10)', () => {
  test('обновляет имя и email владельца', async ({ page }) => {
    const admin = new AdminPage(page);
    await admin.goto();

    await expect(admin.getOwnerNameInput()).toBeVisible({ timeout: 10000 });
    await expect(admin.getOwnerNameInput()).toHaveValue(SEED_OWNER.name);

    await admin.getOwnerNameInput().fill('Новый Владелец');
    await admin.getOwnerEmailInput().fill('new@owner.com');
    await admin.getSaveOwnerButton().click();

    await page.reload();
    await expect(admin.getOwnerNameInput()).toHaveValue('Новый Владелец');
    await expect(admin.getOwnerEmailInput()).toHaveValue('new@owner.com');
  });
});
