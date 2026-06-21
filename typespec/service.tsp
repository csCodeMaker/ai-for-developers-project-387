import "@typespec/http";
import "@typespec/rest";

using TypeSpec.Http;
using TypeSpec.Rest;

@service({
  title: "Booking Service",
})
@server("http://localhost:3000", "Dev server")
namespace BookingService;

// ─── Guest API ───────────────────────────────────────────────

@route("api/event-types")
interface GuestEventTypes {
  /** Только active типы (isDisabled = false) */
  @get list(): EventType[];
}

@route("api/event-types/{eventTypeId}/slots")
interface GuestSlots {
  /** Слоты на дату в окне +14 дней (с признаком занятости) */
  @get list(
    @path eventTypeId: string,
    @query date: plainDate,
  ): Slot[] | ErrorResponse;
}

@route("api/bookings")
interface GuestBookings {
  /** Создать бронь. endTime вычисляется сервером по eventType.duration */
  @post create(@body body: CreateBookingRequest): Booking | ErrorResponse;
}

// ─── Admin API ───────────────────────────────────────────────

@route("api/admin/owner")
interface AdminOwner {
  /** Получить профиль владельца */
  @get get(): Owner | ErrorResponse;

  /** Обновить профиль владельца */
  @put update(@body body: Owner): Owner | ErrorResponse;
}

@route("api/admin/event-types")
interface AdminEventTypes {
  /** Все типы событий (включая disabled) */
  @get list(): EventType[];

  /** Создать новый тип события */
  @post create(@body body: CreateEventTypeRequest): EventType | ErrorResponse;
}

@route("api/admin/event-types/{id}")
interface AdminEventTypeById {
  /** Удалить (пометить isDisabled = true) */
  @delete remove(@path id: string): void | ErrorResponse;

  /** Обновить тип события */
  @put update(
    @path id: string,
    @body body: CreateEventTypeRequest,
  ): EventType | ErrorResponse;
}

@route("api/admin/bookings")
interface AdminBookings {
  /** Все предстоящие бронирования (включая disabled event-types) */
  @get list(): Booking[];
}
