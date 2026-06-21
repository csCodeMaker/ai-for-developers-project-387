import { test, expect } from '@playwright/test';
import { HomePage } from '../pages/home-page';
import { BookingPage } from '../pages/booking-page';
import { ConfirmationPage } from '../pages/confirmation-page';
import { API_BASE } from '../fixtures/test-data';

test.describe.configure({ mode: 'serial' });

test.describe('Гость — просмотр типов событий (US-1)', () => {
  test('видит карточки активных типов с названием, описанием и длительностью', async ({ page }) => {
    const home = new HomePage(page);
    await home.goto();

    await expect(home.getEventTypeCards().first()).toBeVisible({ timeout: 10000 });
    await expect(page.getByText('30-минутный звонок')).toBeVisible();
    await expect(page.getByText('Быстрый созвон на 30 минут')).toBeVisible();
    await expect(page.getByText('30 мин', { exact: true })).toBeVisible();
  });
});

test.describe('Гость — выбор даты и слотов (US-2)', () => {
  test('выбирает дату и видит свободные слоты', async ({ page }) => {
    const home = new HomePage(page);
    await home.goto();
    await page.getByRole('button', { name: 'Записаться' }).click();
    await page.waitForURL(/\/book\//);

    await expect(page.locator('[data-testid="date-picker"]')).toBeVisible({ timeout: 10000 });

    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);

    const responsePromise = page.waitForResponse(
      (res) => res.url().includes('/slots') && res.status() === 200,
    );
    await page.locator('[data-testid="date-picker"] button').filter({ hasText: String(tomorrow.getDate()) }).first().click();
    await responsePromise;

    const booking = new BookingPage(page);
    await expect(booking.getSlots().first()).toBeVisible({ timeout: 5000 });
    const count = await booking.getSlots().count();
    expect(count).toBeGreaterThan(0);
  });

  test('занятые слоты отображаются как заблокированные', async ({ page }) => {
    const res = await page.request.get(`${API_BASE}/api/admin/event-types`);
    const eventTypes = await res.json();
    const eventTypeId = eventTypes[0].id;

    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    const dateStr = tomorrow.toISOString().slice(0, 10);
    const slotStart = `${dateStr}T08:00:00.000Z`;

    await page.request.post(`${API_BASE}/api/bookings`, {
      data: {
        eventTypeId,
        guestName: 'Тест',
        guestEmail: 'test@test.com',
        startTime: slotStart,
      },
    });

    const home = new HomePage(page);
    await home.goto();
    await page.getByRole('button', { name: 'Записаться' }).click();
    await page.waitForURL(/\/book\//);

    const responsePromise = page.waitForResponse(
      (res) => res.url().includes('/slots') && res.status() === 200,
    );
    await page.locator('[data-testid="date-picker"] button').filter({ hasText: String(tomorrow.getDate()) }).first().click();
    await responsePromise;

    const booking = new BookingPage(page);
    await expect(booking.getBusySlots().first()).toBeVisible({ timeout: 5000 });
  });
});

test.describe('Гость — создание бронирования (US-3)', () => {
  test('заполняет форму и успешно бронирует слот', async ({ page }) => {
    const home = new HomePage(page);
    await home.goto();
    await page.getByRole('button', { name: 'Записаться' }).click();
    await page.waitForURL(/\/book\//);

    const booking = new BookingPage(page);

    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);

    const responsePromise = page.waitForResponse(
      (res) => res.url().includes('/slots') && res.status() === 200,
    );
    const dayBtn = booking.getDatePicker().getByRole('button').filter({ hasText: String(tomorrow.getDate()) }).first();
    await expect(dayBtn).toBeVisible();
    await dayBtn.click();
    await responsePromise;

    await expect(booking.getFreeSlots().first()).toBeVisible({ timeout: 5000 });
    await booking.clickSlot(0);

    await booking.fillName('Иван Петров');
    await booking.fillEmail('ivan@example.com');

    await booking.submit();

    await expect(page).toHaveURL(/\/booking\//);
    const confirmation = new ConfirmationPage(page);
    await expect(confirmation.getSuccessMessage()).toBeVisible();
  });
});

test.describe('Гость — подтверждение бронирования (US-4)', () => {
  test('видит ID бронирования и кнопку возврата на главную', async ({ page }) => {
    const home = new HomePage(page);
    await home.goto();
    await page.getByRole('button', { name: 'Записаться' }).click();
    await page.waitForURL(/\/book\//);

    const booking = new BookingPage(page);

    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);

    const responsePromise = page.waitForResponse(
      (res) => res.url().includes('/slots') && res.status() === 200,
    );
    const dayBtn = booking.getDatePicker().getByRole('button').filter({ hasText: String(tomorrow.getDate()) }).first();
    await expect(dayBtn).toBeVisible();
    await dayBtn.click();
    await responsePromise;

    await expect(booking.getFreeSlots().first()).toBeVisible({ timeout: 5000 });
    await booking.clickSlot(0);

    await booking.fillName('Мария Сидорова');
    await booking.fillEmail('maria@example.com');
    await booking.submit();

    await expect(page).toHaveURL(/\/booking\//);
    const confirmation = new ConfirmationPage(page);
    await expect(confirmation.getBookingId()).toBeVisible();
    await expect(confirmation.getBookingId()).toContainText('#');

    await confirmation.clickBackToHome();
    await expect(page).toHaveURL('/');
  });
});

test.describe('Гость — ошибка при занятом слоте (US-5)', () => {
  test('получает ошибку 409 при попытке забронировать занятый слот', async ({ page }) => {
    const res = await page.request.get(`${API_BASE}/api/admin/event-types`);
    const eventTypes = await res.json();
    const eventTypeId = eventTypes[0].id;

    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);

    const home = new HomePage(page);
    await home.goto();
    await page.getByRole('button', { name: 'Записаться' }).click();
    await page.waitForURL(/\/book\//);

    const booking = new BookingPage(page);

    const slotsPromise = page.waitForResponse(
      (res) => res.url().includes('/slots') && res.status() === 200,
    );
    const dayBtn = booking.getDatePicker().getByRole('button').filter({ hasText: String(tomorrow.getDate()) }).first();
    await expect(dayBtn).toBeVisible();
    await dayBtn.click();
    await slotsPromise;

    const freeSlots = booking.getFreeSlots();
    await expect(freeSlots.first()).toBeVisible({ timeout: 5000 });

    const slotStartTime = await freeSlots.first().getAttribute('data-start-time');
    expect(slotStartTime).toBeTruthy();

    await freeSlots.first().click();

    await booking.fillName('Второй');
    await booking.fillEmail('second@test.com');

    // Другой пользователь бронирует тот же слот через API
    await page.request.post(`${API_BASE}/api/bookings`, {
      data: {
        eventTypeId,
        guestName: 'Первый',
        guestEmail: 'first@test.com',
        startTime: slotStartTime,
      },
    });

    const submitPromise = page.waitForResponse(
      (res) => res.url().includes('/api/bookings') && res.status() === 409,
    );
    await booking.submit();
    await submitPromise;

    await expect(booking.getErrorMessage()).toBeVisible({ timeout: 5000 });
  });

  test('получает ошибку при запросе слотов отключённого типа', async ({ page }) => {
    const booking = new BookingPage(page);
    await booking.goto('00000000-0000-0000-0000-000000000000');

    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);

    const dayBtn = booking.getDatePicker().getByRole('button').filter({ hasText: String(tomorrow.getDate()) }).first();
    await expect(dayBtn).toBeVisible();
    await dayBtn.click();

    await expect(page.getByRole('alert')).toBeVisible({ timeout: 5000 });
  });
});
