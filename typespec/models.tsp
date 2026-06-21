/** Владелец календаря — единственный предустановленный профиль */
model Owner {
  id: string;
  name: string;
  email: string;
  description: string;
  timeZone: string;
}

/** Тип события, создаваемый владельцем */
model EventType {
  id: string;
  title: string;
  description: string;
  duration: int32;
  isDisabled: boolean;
}

/** Бронирование, созданное гостем */
model Booking {
  id: string;
  eventTypeId: string;
  guestName: string;
  guestEmail: string;
  startTime: utcDateTime;
  endTime: utcDateTime;
  createdAt: utcDateTime;
}

/** Слот 30 минут (endTime = startTime + eventType.duration) */
model Slot {
  startTime: utcDateTime;
  endTime: utcDateTime;
  isBusy: boolean;
}

/** Единый формат ошибки */
model ErrorResponse {
  code: string;
  message: string;
}

/** Запрос на создание брони (endTime вычисляется сервером) */
model CreateBookingRequest {
  eventTypeId: string;
  guestName: string;
  guestEmail: string;
  startTime: utcDateTime;
}

/** Запрос на создание/обновление типа события */
model CreateEventTypeRequest {
  title: string;
  description: string;
  duration: int32;
}
